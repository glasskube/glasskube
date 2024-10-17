package toast

type ToastInput struct {
	Message     string
	Dismissible bool
	Severity    severity
}

func ForToast(err error, severity severity, dismissible bool) ToastInput {
	return ToastInput{
		Message:     err.Error(),
		Dismissible: dismissible,
		Severity:    severity,
	}
}
