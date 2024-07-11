package toast

type ToastInput struct {
	Message     string
	Dismissible bool
	CssClass    string
}

func ForToast(err error, cssClass string) ToastInput {
	return ToastInput{
		Message:     err.Error(),
		Dismissible: false,
		CssClass:    cssClass,
	}
}
