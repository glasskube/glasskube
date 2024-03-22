package alert

type alertInput struct {
	Message     string
	Dismissible bool
}

func ForAlert(err error) alertInput {
	return alertInput{
		Message:     err.Error(),
		Dismissible: false,
	}
}
