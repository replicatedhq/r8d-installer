package kots

type KOTS struct {
	Name    string
	Version string
}

func New(version string) *KOTS {
	return &KOTS{
		"kots",
		version,
	}
}
