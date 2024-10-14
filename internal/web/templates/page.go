package templates

type Navbar struct {
	ActiveItem string
}

type VersionDetails struct {
	OperatorVersion     string
	ClientVersion       string
	NeedsOperatorUpdate bool
}

type Page struct {
	Navbar             Navbar
	VersionDetails     VersionDetails
	CurrentContext     string
	GitopsMode         bool
	Error              error
	CacheBustingString string
	CloudId            string
}
