package velero

type Velero struct {
	Name             string
	Version          string
	AWSVersion       string
	AzureVersion     string
	GCPVersion       string
	LVPVersion       string
	KurlUtilsVersion string
}

func New(version, AWSVersion, azureVersion, GVPVersion, LVPVersion, kurlUtilsVersion string) *Velero {
	return &Velero{
		"velero",
		version,
		AWSVersion,
		azureVersion,
		GVPVersion,
		LVPVersion,
		kurlUtilsVersion,
	}
}
