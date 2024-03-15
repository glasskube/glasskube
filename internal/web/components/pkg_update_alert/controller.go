package pkg_update_alert

import (
	"html/template"
	"io"
)

const TemplateId = "pkg-update-alert"

type pkgUpdateAlertInput struct {
	UpdatesAvailable bool
}

func Render(w io.Writer, tmpl *template.Template, updatesAvailabele bool) error {
	return tmpl.ExecuteTemplate(w, TemplateId, &pkgUpdateAlertInput{UpdatesAvailable: updatesAvailabele})
}

func ForPkgUpdateAlert(data map[string]any) *pkgUpdateAlertInput {
	return &pkgUpdateAlertInput{UpdatesAvailable: data["UpdatesAvailable"].(bool)}
}
