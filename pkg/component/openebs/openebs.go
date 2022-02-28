package openebs

type OpenEBS struct {
	Name    string
	Version string
}

func New(version string) *OpenEBS {
	return &OpenEBS{
		"openebs",
		version,
	}
}
