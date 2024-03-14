package cliutils

import (
	"context"
	"fmt"
	"os"

	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/google/go-containerregistry/pkg/name"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPackageOperatorVersion(ctx context.Context) (string, error) {
	config := client.ConfigFromContext(ctx)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	namespace := "glasskube-system"
	deploymentName := "glasskube-controller-manager"
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
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

func CheckPackageOperatorVersion(ctx context.Context) error {
	operatorVersion, err := GetPackageOperatorVersion(ctx)
	if err != nil {
		return err
	}
	if operatorVersion[1:] != config.Version {
		fmt.Fprintf(os.Stderr, "❗ Glasskube PackageOperator needs to be updated: %s -> %s\n", operatorVersion[1:], config.Version)
		fmt.Fprintf(os.Stderr, "💡 Please run `glasskube bootstrap` again to update Glasskube PackageOperator\n")
	}
	return nil
}
