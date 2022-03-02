// +build deps

package kots

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/replicatedhq/r8d-installer/pkg/utils"
)

const (
	owner = "replicatedhq"
	repo  = "kots"
)

// GetVersion returns the version of KOTS used in this binary
func (k *KOTS) GetVersion() string {
	return k.Version
}

// GetVersion returns name for printing
func (k *KOTS) GetName() string {
	return k.Name
}

// GetManifest merges all the YAML files for KOTS into a single YAML.
// Based on the `kot admin-console generate-manifest` command
func (k *KOTS) GetManifests() (string, error) {

	var binaryName string
	switch runtime.GOOS {
	case "darwin":
		binaryName = "kots_darwin_all.tar.gz"
	case "linux":
		binaryName = "kots_linux_amd64.tar.gz"
	}

	cli, err := utils.DownloadAssetFromGithubRelease("replicatedhq", "kots", k.GetVersion(), binaryName)
	if err != nil {
		return "", errors.Wrapf(err, "failed to download %s", binaryName)
	}
	defer os.Remove(cli)

	// unpack into a temp dir
	tmpDir, err := os.MkdirTemp("", "kots-cli-")
	if err != nil {
		return "", errors.Wrap(err, "failed to create temp dir for kots cli bundle")
	}
	defer os.RemoveAll(tmpDir)

	if err := utils.Untar(cli, tmpDir); err != nil {
		return "", errors.Wrap(err, "failed to unpack kots cli bundle")
	}

	// exe cli
	yamlDir, err := os.MkdirTemp("", "kots-manifests-")
	if err != nil {
		return "", errors.Wrap(err, "failed to create temp dir for kots manifests")
	}
	defer os.RemoveAll(yamlDir)

	args := []string{
		"admin-console",
		"generate-manifests",
		fmt.Sprintf("--rootdir=%s", yamlDir),
		"--with-minio=false",
		"--namespace=default",
		"--shared-password=password",
	}
	cmd := exec.Command(path.Join(tmpDir, "kots"), args...)
	if err = cmd.Run(); err != nil {
		return "", errors.Wrapf(err, "failed to run %s", cmd.String())
	}

	// merge manifest
	manifest := path.Join(os.TempDir(), "kots.yaml")
	os.Remove(manifest)
	if err := utils.MergeYAML(yamlDir, manifest); err != nil {
		return "", errors.Wrap(err, "failed to merge kots manifests")
	}

	return manifest, nil
}

// GetManifest returns a file path to the compressed airgap images for KOTS.
// It's assume that the caller owns the file after calling.
// KOTS provides the pre-built compressed archive, but it's organized weirdly and there are unneeded images,
// so we'll pull images from the ENV and package them outselves
// TODO (dans): this is a hack. We should just publish the required bundle in KOTS, or at least the manifest for r8d-installer
func (k *KOTS) GetImageArchive() (string, error) {
	imageEnvContext, err := utils.GetSourceFileFromGithubRelease(owner, repo, k.GetVersion(), ".image.env")
	if err != nil {
		return "", errors.Wrap(err, "failed to download .image.env file")
	}

	imageEnv, err := godotenv.Unmarshal(imageEnvContext)

	imageMap := map[string]string{}

	err = getImageTags(&imageEnv, &imageMap, k.GetVersion())
	if err != nil {
		return "", errors.Wrap(err, "failed to filter images")
	}

	// Save each of the images
	imageList := []string{}
	for image, tag := range imageMap {
		imageList = append(imageList, image+":"+tag)
	}

	archivePath, err := utils.CreateArchive(k.GetName(), imageList)
	if err != nil {
		return "", errors.Wrap(err, "failed to create image archive")
	}

	// Rename the generic archive
	destPath := path.Join(os.TempDir(), "kots-image-archive.tar.zst")
	if err = os.Rename(archivePath, destPath); err != nil {
		return "", errors.Wrap(err, "failed to rename kots image archive")
	}

	return destPath, nil
}

// getImageTags returns a filtered map of images to tags used for r8d-installer
// TODO (dans): remove hardcoded image tags
func getImageTags(input, output *map[string]string, tag string) error {

	if val, ok := (*input)["POSTGRES_ALPINE_TAG"]; ok {
		(*output)["docker.io/postgres"] = val
	} else {
		return errors.New("failed to find POSTGRES_ALPINE_TAG in .image.env")
	}

	if val, ok := (*input)["DEX_TAG"]; ok {
		(*output)["docker.io/dexidp/dex"] = val
	} else {
		return errors.New("failed to find POSTGRES_ALPINE_TAG in .image.env")
	}

	(*output)["docker.io/kotsadm/kotsadm"] = tag
	return nil
}

// GetBinaries returns a file path to the KOTS binary.
// This will be the kubectl plugin.
func (k *KOTS) GetBinaries() ([]string, error) {
	cli, err := utils.DownloadAssetFromGithubRelease(owner, repo, k.GetVersion(), "kots_linux_amd64.tar.gz")
	if err != nil {
		return nil, errors.Wrap(err, "failed to download kots_linux_amd64.tar.gz")
	}

	binaries := []string{
		cli,
	}
	return binaries, nil
}
