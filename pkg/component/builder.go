// +build deps

package component

type Builder interface {
	GetVersion() string
	GetName() string
	GetManifest() (string, error)
	GetImageArchive() (string, error)
	GetBinaries() ([]string, error)
}
