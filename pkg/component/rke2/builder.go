// +build deps

package rke2

import (
	"github.com/pkg/errors"
	"github.com/replicatedhq/r8d-installer/pkg/utils"
)

// GetVersion returns the version of RKE2 used in this binary
func (r *RKE2) GetVersion() string {
	return r.Version
}

// GetVersion returns name for printing
func (r *RKE2) GetName() string {
	return r.Name
}

// GetManifest returns an empty string because RKE2 does not require a manifest
func (r *RKE2) GetManifests() (string, error) {
	return "", nil
}

// GetManifest returns a file path to the compressed airgap images for RKE2.
// It's assume that the caller owns the file after calling.
// RKE2 provides the pre-built compressed archive, so we'll just use that
func (r *RKE2) GetImageArchive() (string, error) {
	airgapImages, err := utils.DownloadAssetFromGithubRelease("rancher", "rke2", r.GetVersion(), "rke2-images-canal.linux-amd64.tar.zst")
	if err != nil {
		return "", errors.Wrap(err, "failed to download rke2-linux-amd64.tar.gz")
	}

	return airgapImages, nil
}

// GetBinaries returns a file path to the compressed airgap images for RKE2.
// It's assume that the caller owns the file after calling.
func (r *RKE2) GetBinaries() ([]string, error) {
	cli, err := utils.DownloadAssetFromGithubRelease("rancher", "rke2", r.GetVersion(), "rke2.linux-amd64.tar.gz")
	if err != nil {
		return nil, errors.Wrap(err, "failed to download rke2-linux-amd64.tar.gz")
	}

	binaries := []string{
		cli,
	}
	return binaries, nil
}
