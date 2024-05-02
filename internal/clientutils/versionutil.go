package clientutils

import (
	"context"

	"k8s.io/client-go/rest"

	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPackageOperatorVersion(ctx context.Context) (string, error) {
	config := clicontext.ConfigFromContext(ctx)
	return GetPackageOperatorVersionForConfig(config, ctx)
}

func GetPackageOperatorVersionForConfig(config *rest.Config, ctx context.Context) (string, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	namespace := "glasskube-system"
	deploymentName := "glasskube-controller-manager"
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, deploymentName, v1.GetOptions{})
	if err != nil {
		return "", err
	}

	containers := deployment.Spec.Template.Spec.Containers
	for _, container := range containers {
		if container.Name == "manager" {
			ref, err := name.ParseReference(container.Image)
			if err != nil {
				return "", err
			}
			return ref.Identifier(), nil
		}
	}
	return "", nil
}
