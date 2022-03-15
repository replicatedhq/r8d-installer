//go:build deps
// +build deps

package velero

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/replicatedhq/r8d-installer/pkg/utils"
)

const (
	owner = "vmware-tanzu"
	repo  = "velero"
)

// GetVersion returns the version of Velero used in this binary
func (v *Velero) GetVersion() string {
	return v.Version
}

// GetVersion returns "velero"
func (v *Velero) GetName() string {
	return v.Name
}

// GetManifest merges all the YAML files for Velero into a single YAML.
// Based on the `velero install --dry-run` command
func (v *Velero) GetManifests() (string, error) {

	yamlBuffer, err := v.getVeleroManifests()
	if err != nil {
		return "", errors.Wrap(err, "failed to get velero manifests")
	}

	// merge manifests
	manifest := path.Join(os.TempDir(), "velero.yaml")
	file, err := os.Create(manifest)
	if err != nil {
		return "", errors.Wrap(err, "failed to create velero.yaml")
	}
	defer file.Close()

	if _, err := file.Write(yamlBuffer); err != nil {
		return "", errors.Wrap(err, "failed to write to velero.yaml")
	}

	return manifest, nil
}

// GetManifest returns a file path to the compressed airgap images for Velero.
// It's assume that the caller owns the file after calling.
// It pulls the images from the manifest files.
func (v *Velero) GetImageArchive() (string, error) {

	yamlBuffer, err := v.getVeleroManifests()
	if err != nil {
		return "", errors.Wrap(err, "failed to get velero manifests")
	}

	images := utils.GetImages(string(yamlBuffer))

	archivePath, err := utils.CreateArchive(v.GetName(), images)
	if err != nil {
		return "", errors.Wrap(err, "failed to create image archive")
	}

	// Rename the generic archive
	destPath := path.Join(os.TempDir(), "velero-image-archive.tar.zst")
	if err = os.Rename(archivePath, destPath); err != nil {
		return "", errors.Wrap(err, "failed to rename velero image archive")
	}

	return destPath, nil
}

// getVeleroManifests returns the string representation of the Velero manifests
func (v *Velero) getVeleroManifests() ([]byte, error) {

	binaryName := fmt.Sprintf("velero-%s-%s-amd64.tar.gz", v.GetVersion(), runtime.GOOS)

	cli, err := utils.DownloadAssetFromGithubRelease(owner, repo, v.GetVersion(), binaryName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download %s", binaryName)
	}
	defer os.Remove(cli)

	// unpack into a temp dir
	tmpDir, err := os.MkdirTemp("", "velero-cli-")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temp dir for velero cli bundle")
	}
	defer os.RemoveAll(tmpDir)

	if err := utils.Untar(cli, tmpDir); err != nil {
		return nil, errors.Wrap(err, "failed to unpack velero cli bundle")
	}

	// exe cli
	yamlDir, err := os.MkdirTemp("", "velero-manifests-")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temp dir for velero manifests")
	}
	defer os.RemoveAll(yamlDir)

	subDir := fmt.Sprintf("velero-%s-%s-amd64", v.GetVersion(), runtime.GOOS)

	plugins := []string{
		fmt.Sprintf("velero/velero-plugin-for-aws:%s", v.AWSVersion),
		fmt.Sprintf("velero/velero-plugin-for-microsoft-azure:%s", v.AzureVersion),
		fmt.Sprintf("velero/velero-plugin-for-gcp:%s", v.GCPVersion),
		fmt.Sprintf("replicated/local-volume-provider:%s", v.LVPVersion),
		fmt.Sprintf("replicated/kurl-util:%s", v.KurlUtilsVersion),
	}

	// Don't install a BSL, because this should be managed by the admin console.
	args := []string{
		"install",
		"--dry-run",
		"--use-restic",
		"--use-volume-snapshots=false",
		"--namespace=velero",
		"--no-default-backup-location",
		"--no-secret",
		"--plugins",
		strings.Join(plugins, ","),
		"-oyaml",
	}
	cmd := exec.Command(path.Join(tmpDir, subDir, "velero"), args...)
	yamlBuffer, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to run %s", cmd.String())
	}
	return yamlBuffer, nil
}

// GetBinaries returns a file path to the compressed airgap images for Velero.
// This will be the kubectl plugin.
func (v *Velero) GetBinaries() ([]string, error) {

	binaryName := fmt.Sprintf("velero-%s-linux-amd64.tar.gz", v.GetVersion())

	cli, err := utils.DownloadAssetFromGithubRelease(owner, repo, v.GetVersion(), binaryName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download %s", binaryName)
	}

	binaries := []string{
		cli,
	}
	return binaries, nil
}
