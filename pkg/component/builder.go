// +build deps

package component

type Builder interface {
	// GetVersion returns the version tag associated with this component
	GetVersion() string

	// GetName returns the lower-camelcase name of this component
	GetName() string

	// GetManfest returns a file path to a single multi-doc YAML that encompasses all
	// manifests for this component.
	GetManifests() (string, error)

	// GetImageArchive returns a file path to the compressed airgap images for this component.
	// It's the component that owns downloading and archiving the images to a tar.zst.
	// It's assumed that the caller owns the files after calling.
	GetImageArchive() (string, error)

	// GetBinaries returns a list of file paths to the binaries for this component.
	// It's assumed that the caller owns the files after calling.
	GetBinaries() ([]string, error)
}
