package rke2

type RKE2 struct {
	Name    string
	Version string
}

func New(version string) *RKE2 {
	return &RKE2{
		"rke2",
		version,
	}
}
