package types

type Navbar struct {
	ActiveItem string
}

type VersionDetails struct {
	MismatchWarning     bool
	OperatorVersion     string
	ClientVersion       string
	NeedsOperatorUpdate bool
}

type TemplateContextData struct {
	Navbar             Navbar
	VersionDetails     VersionDetails
	CurrentContext     string
	GitopsMode         bool
	Error              error
	CacheBustingString string
	CloudId            string
	TemplateName       string
}

type TemplateContextHolder struct {
	Ctx TemplateContextData
}

func (p *TemplateContextHolder) SetContextData(ctx TemplateContextData) {
	p.Ctx = ctx
}

type ContextInjectable interface {
	SetContextData(TemplateContextData)
}
