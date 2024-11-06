package responder

type Navbar struct {
	ActiveItem string
}

type VersionDetails struct {
	MismatchWarning     bool
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
	TemplateName       string
	TemplateData       any
}
