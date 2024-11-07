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

type TemplateContext struct {
	Navbar             Navbar
	VersionDetails     VersionDetails
	CurrentContext     string
	GitopsMode         bool
	Error              error
	CacheBustingString string
	CloudId            string
	TemplateName       string
}

type TemplateData = any

type Page struct {
	TemplateContext
	TemplateData
}
