// +build deps

package kots

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/replicatedhq/kots/pkg/upstream"
	upstreamtypes "github.com/replicatedhq/kots/pkg/upstream/types"
	"github.com/replicatedhq/r8d-installer/pkg/utils"
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
// https://github.com/replicatedhq/kots/blob/436b9212ef59494e91d4dde03f165665ec105237/cmd/kots/cli/admin-console-generate-manifests.go
func (k *KOTS) GetManifests() (string, error) {

	tmpDir, err := os.MkdirTemp("", "kots-manifest-")
	if err != nil {
		return "", errors.Wrap(err, "failed to create kots manifest temp dir")
	}

	options := upstreamtypes.WriteOptions{
		Namespace:      "default",
		SharedPassword: "password",
		// HTTPProxyEnvValue:    v.GetString("http-proxy"),
		// HTTPSProxyEnvValue:   v.GetString("https-proxy"),
		// NoProxyEnvValue:      v.GetString("no-proxy"),
		IncludeMinio:         false,
		IsMinimalRBAC:        false,
		AdditionalNamespaces: nil,
	}
	adminConsoleFiles, err := upstream.GenerateAdminConsoleFiles(tmpDir, options)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate admin console files")
	}

	for _, file := range adminConsoleFiles {
		fileRenderPath := filepath.Join(tmpDir, file.Path)
		d, _ := filepath.Split(fileRenderPath)
		if _, err := os.Stat(d); os.IsNotExist(err) {
			if err := os.MkdirAll(d, 0744); err != nil {
				return "", errors.Wrap(err, "failed to mkdir")
			}
		}

		if err := ioutil.WriteFile(fileRenderPath, file.Content, 0644); err != nil {
			return "", errors.Wrapf(err, "failed to write file %s", fileRenderPath)
		}
	}

	fmt.Printf("Admin Console manifests created in %s", filepath.Join(tmpDir, "admin-console"))

	return "", nil
}

// GetManifest returns a file path to the compressed airgap images for KOT.
// It's assume that the caller owns the file after calling.
// KOTS provides the pre-built compressed archive, but it's organized weirdly and there are unneeded images,
// so we'll pull images from the ENV and package them outselves
// TODO (dans): this is a hack. We should just publish the required bundle in KOTS, or at least the manifest for r8d-installer
func (k *KOTS) GetImageArchive() (string, error) {
	imageEnvContext, err := utils.GetSourceFileFromGithubRelease("replicatedhq", "kots", k.GetVersion(), ".image.env")
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
	ourDir, err := os.MkdirTemp("", "kots-image-archive-")
	if err != nil {
		return "", errors.Wrap(err, "failed to create temp dir for kots image archive")
	}
	destPath := path.Join(ourDir, "kots-image-archive.tar.zst")
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

// GetBinaries returns a file path to the compressed airgap images for KOTS.
// This will be the kubectl plugin.
func (k *KOTS) GetBinaries() ([]string, error) {
	cli, err := utils.DownloadAssetFromGithubRelease("replicatedhq", "kots", k.GetVersion(), "kots_linux_amd64.tar.gz")
	if err != nil {
		return nil, errors.Wrap(err, "failed to download kots_linux_amd64.tar.gz")
	}

	binaries := []string{
		cli,
	}
	return binaries, nil
}
