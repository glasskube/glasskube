package bootstrap

import (
	"context"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func IsBootstrapped(ctx context.Context, cfg *rest.Config) (bool, error) {
	cs, err := clientset.NewForConfig(cfg)
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

func crdExists(ctx context.Context, clientset clientset.Interface, crdName string) (bool, error) {
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
