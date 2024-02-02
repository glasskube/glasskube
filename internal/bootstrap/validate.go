package bootstrap

import (
	"context"
	"fmt"
	"os"

	"github.com/glasskube/glasskube/api/v1alpha1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

var bootstrapMessage = `
Sorry, it seems Glasskube is not yet bootstrapped in your cluster!

As Glasskube is still in a technical preview phase, please execute the bootstrap command by yourself:

glasskube bootstrap

For further information on bootstrapping, please consult the docs: https://glasskube.dev/docs/getting-started/bootstrap
If you need any help or run into issues, don't hesitate to contact us:
Github: https://github.com/glasskube/glasskube
Discord: https://discord.gg/SxH6KUCGH7

`

func RequireBootstrapped(ctx context.Context, cfg *rest.Config) {
	ok, err := IsBootstrapped(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error validating Glasskube:\n\n%v\n", err)
		os.Exit(1)
	}
	if !ok {
		fmt.Fprint(os.Stderr, bootstrapMessage)
		os.Exit(1)
	}
}

func IsBootstrapped(ctx context.Context, cfg *rest.Config) (bool, error) {
	cs, err := apiextension.NewForConfig(cfg)
	if err != nil {
		return false, err
	}
	pkgsExist, err := crdExists(ctx, cs, "packages")
	if err != nil {
		return false, err
	}
	pisExist, err := crdExists(ctx, cs, "packageinfos")
	if err != nil {
		return false, err
	}
	return pkgsExist && pisExist, nil
}

func crdExists(ctx context.Context, clientset *apiextension.Clientset, crdName string) (bool, error) {
	_, err := clientset.ApiextensionsV1().
		CustomResourceDefinitions().
		Get(ctx, fmt.Sprintf("%s.%s", crdName, v1alpha1.GroupVersion.Group), v1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
