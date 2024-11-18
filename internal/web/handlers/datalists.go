package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/glasskube/glasskube/internal/web/components"
	opts "github.com/glasskube/glasskube/internal/web/options"
	"github.com/glasskube/glasskube/internal/web/responder"
)

// GetNamesDatalist is a GET endpoint returning an html datalist, containing options depending on the given valueName,
// kind of reference and namespace. It is only usable for ConfigMap and Secret refs, since packages don't have a
// namespace. In case the refKind is ConfigMap, the datalist contains the config maps of the given namespace; in case
// the refKind is Secret, the datalist contains the secrets of the given namespace; if no namespace is given or an
// error occurs, an empty datalist is returned
func GetNamesDatalist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	valueName := r.PathValue("valueName")
	refKind := r.FormValue(refKindKey)
	id := r.FormValue("id")
	nsKey := formKey(valueName, namespaceKey)
	namespace := r.Form.Get(nsKey)
	var options []string
	if refKind == refKindConfigMap {
		if opts, err := opts.GetConfigMapNameOptions(ctx, namespace); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get config map name options: %v\n", err)
		} else {
			options = opts
		}
	} else if refKind == refKindSecret {
		if opts, err := opts.GetSecretNameOptions(ctx, namespace); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get secret name options: %v\n", err)
		} else {
			options = opts
		}
	}
	responder.SendComponent(w, r, "components/datalist",
		responder.RawTemplate(components.DatalistInput{
			Options: options,
			Id:      id,
		}))
}

func GetKeysDatalist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	valueName := r.PathValue("valueName")
	refKind := r.FormValue(refKindKey)
	nsKey := formKey(valueName, namespaceKey)
	nameKey := formKey(valueName, nameKey)
	pkgKey := formKey(valueName, packageKey)
	namespace := r.Form.Get(nsKey)
	name := r.Form.Get(nameKey)
	pkg := r.Form.Get(pkgKey)
	var options []string
	var err error
	if refKind == refKindConfigMap {
		if options, err = opts.GetConfigMapKeyOptions(ctx, namespace, name); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get ConfigMap key options of %v in %v: %v\n", name, namespace, err)
		}
	} else if refKind == refKindSecret {
		if options, err = opts.GetSecretKeyOptions(ctx, namespace, name); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get Secret key options of %v in %v: %v\n", name, namespace, err)
		}
	} else if refKind == refKindPackage {
		if options, err = opts.GetPackageValuesOptions(ctx, pkg); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get package value options of %v: %v\n", pkg, err)
		}
	}
	responder.SendComponent(w, r, "components/datalist",
		responder.RawTemplate(components.DatalistInput{
			Id:      r.FormValue("id"),
			Options: options,
		}))
}
