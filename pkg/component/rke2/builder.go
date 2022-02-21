// +build deps

package rke2

// GetVersion returns the version of RKE2 used in this binary
func (r *RKE2) GetVersion() string {
	return r.Version
}

// GetVersion returns name for printing
func (r *RKE2) GetName() string {
	return r.Name
}

// GetManifest returns an empty string because RKE2 does not require a manifest
func (r *RKE2) GetManifest() (string, error) {
	return "", nil
}

// GetManifest returns a file path to the compressed airgap images for RKE2.
// It's assume that the caller owns the file after calling.
func (r *RKE2) GetImageArchive() (string, error) {
	return "", nil
}

// GetBinaries returns a file path to the compressed airgap images for RKE2.
// It's assume that the caller owns the file after calling.
func (r *RKE2) GetBinaries() ([]string, error) {
	return []string{}, nil
}
