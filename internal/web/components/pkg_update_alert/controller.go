package pkg_update_alert

const TemplateId = "pkg-update-alert"

type pkgUpdateAlertInput struct {
	UpdatesAvailable bool
}

func ForPkgUpdateAlert(data map[string]any) *pkgUpdateAlertInput {
	return &pkgUpdateAlertInput{UpdatesAvailable: data["UpdatesAvailable"].(bool)}
}
