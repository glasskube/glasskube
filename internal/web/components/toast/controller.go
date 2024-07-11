package toast

type ToastInput struct {
	Message     string
	Dismissible bool
	CssClass    string
}

func ForToast(err error, cssClass string, dismissible bool) ToastInput {
	return ToastInput{
		Message:     err.Error(),
		Dismissible: dismissible,
		CssClass:    cssClass,
	}
}
