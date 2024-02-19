package uninstall

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var deletePropagationForeground = metav1.DeletePropagationForeground

func Uninstall(pkgClient *client.PackageV1Alpha1Client, ctx context.Context, pkg *v1alpha1.Package) error {
	err := pkgClient.Packages().Delete(ctx, pkg, metav1.DeleteOptions{PropagationPolicy: &deletePropagationForeground})
	if err != nil {
		return err
	}
	return nil
}
