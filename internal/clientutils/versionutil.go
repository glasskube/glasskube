package clientutils

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/client-go/listers/apps/v1"

	"k8s.io/client-go/rest"

	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/google/go-containerregistry/pkg/name"
	"k8s.io/client-go/kubernetes"
)

const (
	namespace      = "glasskube-system"
	deploymentName = "glasskube-controller-manager"
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

	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return getVersionOfDeployment(deployment)
}

func GetPackageOperatorVersionForLister(deploymentLister *appsv1.DeploymentLister) (string, error) {
	deployment, err := (*deploymentLister).Deployments(namespace).Get(deploymentName)
	if err != nil {
		return "", err
	}
	return getVersionOfDeployment(deployment)
}

func getVersionOfDeployment(deployment *v1.Deployment) (string, error) {
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
