package utils

import (
	"archive/tar"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	"github.com/pkg/errors"
	"github.com/valyala/gozstd"
)

type manifestJSON []map[string]interface{}

type repositoryJSON map[string]interface{}

func CreateArchive(name string, images []string) (string, error) {

	dir, err := os.MkdirTemp("", "r8d-archive-")
	if err != nil {
		return "", errors.Wrap(err, "failed to create temp dir")
	}
	defer os.RemoveAll(dir)

	archives := []string{}
	for _, image := range images {
		archivePath, err := createDockerArchive(dir, image)
		if err != nil {
			return "", errors.Wrapf(err, "failed to create docker archive for %s", image)
		}
		archives = append(archives, archivePath)
	}
	// TODO (dans): remove the individual archives

	// extract manifests & merge json files
	tmpDir, err := os.MkdirTemp("", "r8d-merge-archive-")
	if err != nil {
		return "", errors.Wrap(err, "failed to create temp dir")
	}
	defer os.RemoveAll(tmpDir)

	manifestData, repositoriesData, err := extractAndMerge(tmpDir, archives)
	if err != nil {
		return "", errors.Wrap(err, "failed to extract and merge")
	}

	// Write manifest to disk
	manifest, err := os.Create(path.Join(dir, "manifest.json"))
	if err != nil {
		return "", errors.Wrap(err, "failed to create temp manifest file")
	}
	manifest.Write(manifestData)
	if err = manifest.Close(); err != nil {
		return "", errors.Wrap(err, "failed to close temp manifest file")
	}
	defer os.Remove(manifest.Name())

	// Write repositories to disk
	repositories, err := os.Create(path.Join(dir, "repositories"))
	if err != nil {
		return "", errors.Wrap(err, "failed to create temp manifest file")
	}
	repositories.Write(repositoriesData)
	if err = repositories.Close(); err != nil {
		return "", errors.Wrap(err, "failed to close temp repositories file")
	}
	defer os.Remove(repositories.Name())

	// create new archive
	file, err := os.CreateTemp("", "r8d-image-archive-*.tar.zst")
	if err != nil {
		return "", errors.Wrap(err, "failed to create temp file")
	}
	defer file.Close()

	zw := gozstd.NewWriter(file)
	defer zw.Close()

	tw := tar.NewWriter(zw)
	defer tw.Close()

	// add manifest and repository
	if err = addFileToTar(tw, manifest.Name()); err != nil {
		return "", errors.Wrap(err, "failed to add manifest to tar")
	}

	if err = addFileToTar(tw, repositories.Name()); err != nil {
		return "", errors.Wrap(err, "failed to add repositories to tar")
	}

	// Add all the other files
	if err = CopyDirToTar(tw, tmpDir); err != nil {
		return "", errors.Wrap(err, "failed to copy temp dir to tar")
	}

	return file.Name(), nil
}

func createDockerArchive(dir, image string) (string, error) {
	// Good examples of using the container/image library:
	// https://iximiuz.com/en/posts/working-with-container-images-in-go/
	// https://github.com/containers/skopeo/blob/main/cmd/skopeo/copy.go
	// https://github.com/containers/skopeo/blob/main/docs/skopeo-copy.1.md

	srcName := fmt.Sprintf("docker://%s", image)

	srcRef, err := alltransports.ParseImageName(srcName)
	if err != nil {
		return "", errors.Wrapf(err, "invalid source image name %s", srcName)
	}

	policy := &signature.Policy{Default: []signature.PolicyRequirement{signature.NewPRInsecureAcceptAnything()}}
	if err != nil {
		return "", errors.Wrap(err, "failed to get default policy")
	}
	policyCtx, err := signature.NewPolicyContext(policy)
	if err != nil {
		return "", errors.Wrap(err, "failed to get new policy context")
	}
	defer policyCtx.Destroy()

	imageName := strings.Split(image, ":")[0]
	imageName = path.Base(imageName)
	archivePath := path.Join(dir, fmt.Sprintf("%s.tar", imageName))

	// Make sure to use the full path here https://github.com/containers/skopeo/issues/730
	dstName := fmt.Sprintf("docker-archive:%s:%s", archivePath, image)

	dstRef, err := alltransports.ParseImageName(dstName)
	if err != nil {
		return "", errors.Wrapf(err, "invalid destination image name %s", dstName)
	}

	_, err = copy.Image(
		context.Background(),
		policyCtx,
		dstRef,
		srcRef,
		&copy.Options{
			SourceCtx: &types.SystemContext{
				ArchitectureChoice: "amd64",
				OSChoice:           "linux",
			},
			DestinationCtx: &types.SystemContext{
				ArchitectureChoice: "amd64",
				OSChoice:           "linux",
			},
		},
	)
	if err != nil {
		return "", errors.Wrapf(err, "failed to copy image %s to %s", srcName, dstName)
	}

	return archivePath, nil
}

// extractAndMerge extracts each of the tar archives into the same temporary directory
// for merging. It does NOT copy over the manifest.json or repositories, instead
// returning them on success. These can be added first to the new, composite archive.
// https://github.com/rancher/rke2/blob/master/scripts/package-images#L10-L12
func extractAndMerge(tmpDir string, archives []string) ([]byte, []byte, error) {
	manifests := [][]byte{}
	repositories := [][]byte{}

	for _, archive := range archives {
		manifest, repository, err := extractAndFilter(tmpDir, archive)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to extract manifest from %s", archive)
		}
		manifests = append(manifests, manifest)
		repositories = append(repositories, repository)
	}

	var mergedManifestJSON manifestJSON
	for _, file := range manifests {
		var manifest manifestJSON
		err := json.Unmarshal(file, &manifest)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to unmarshal manifest %s", file)
		}
		mergedManifestJSON = append(mergedManifestJSON, manifest...)
	}

	mergedManifest, err := json.Marshal(mergedManifestJSON)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to marshal merged manifest %s", mergedManifestJSON)
	}

	mergedRepositoriesJSON := repositoryJSON{}

	for _, file := range repositories {
		var repository repositoryJSON
		err := json.Unmarshal(file, &repository)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to unmarshal manifest %s", file)
		}
		for k, v := range repository {
			mergedRepositoriesJSON[k] = v
		}
	}

	mergedRepositories, err := json.Marshal(mergedRepositoriesJSON)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to marshal merged manifest %s", mergedRepositoriesJSON)
	}

	return mergedManifest, mergedRepositories, nil
}

func extractAndFilter(tmpDir string, archive string) ([]byte, []byte, error) {
	reader, err := os.Open(archive)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to open archive %s", archive)
	}
	defer reader.Close()

	tr := tar.NewReader(reader)

	var manifest, repositories []byte
	var hasManifest, hasRepositories bool

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to read next header from archive %s", archive)
		}

		if hdr.Name == "manifest.json" {
			hasManifest = true
			manifest, err = ioutil.ReadAll(tr)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "failed to read manifest from archive %s", archive)
			}
			continue
		}

		if hdr.Name == "repositories" {
			hasRepositories = true
			repositories, err = ioutil.ReadAll(tr)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "failed to read repositores from archive %s", archive)
			}
			continue
		}

		if hdr.Typeflag == tar.TypeDir {
			fmt.Println("MAKING DIR", hdr.Name)
			os.MkdirAll(path.Join(tmpDir, hdr.Name), 0755)
			continue
		}

		if hdr.Typeflag == tar.TypeReg {
			outpath := path.Join(tmpDir, hdr.Name)
			if err := os.MkdirAll(path.Dir(outpath), 0755); err != nil {
				return nil, nil, errors.Wrapf(err, "failed to create directory %s", outpath)
			}

			file, err := os.OpenFile(outpath, os.O_CREATE|os.O_WRONLY, os.FileMode(hdr.Mode))
			if err != nil {
				return nil, nil, errors.Wrapf(err, "failed to open file %s", hdr.Name)
			}
			defer file.Close()

			if _, err := io.Copy(file, tr); err != nil {
				return nil, nil, errors.Wrapf(err, "failed to copy file %s", hdr.Name)
			}
		}
	}

	if !hasManifest || !hasRepositories {
		return nil, nil, errors.New("manifest and/or repositories not found in archive")
	}

	return manifest, repositories, nil
}

func addFileToTar(tw *tar.Writer, file string) error {
	fi, err := os.Stat(file)
	if err != nil {
		return errors.Wrapf(err, "failed to stat file %s", file)
	}

	hdr, err := tar.FileInfoHeader(fi, fi.Name())
	if err != nil {
		return errors.Wrapf(err, "failed to create header for file %s", file)
	}

	if err := tw.WriteHeader(hdr); err != nil {
		return errors.Wrapf(err, "failed to write header for file %s", file)
	}

	f, err := os.Open(file)
	if err != nil {
		return errors.Wrapf(err, "failed to open file %s", file)
	}
	defer f.Close()

	if _, err := io.Copy(tw, f); err != nil {
		return errors.Wrapf(err, "failed to copy file %s", file)
	}

	return nil
}
