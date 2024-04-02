package alert

type alertInput struct {
	Message     string
	Dismissible bool
	Type        string
}

func ForAlert(err error, alertType string) alertInput {
	return alertInput{
		Message:     err.Error(),
		Dismissible: false,
		Type:        alertType,
	}
}
