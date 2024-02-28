package bootstrap

import (
	"context"
	"fmt"
	"sync"

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

	var wg sync.WaitGroup
	wg.Add(2)

	var pkgsExist, pisExist bool
	var pkgsErr, pisErr error

	go func() {
		defer wg.Done()
		pkgsExist, pkgsErr = crdExists(ctx, cs, "packages")
	}()

	go func() {
		defer wg.Done()
		pisExist, pisErr = crdExists(ctx, cs, "packageinfos")
	}()

	wg.Wait()

	if pkgsErr != nil {
		return false, pkgsErr
	}
	if pisErr != nil {
		return false, pisErr
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
