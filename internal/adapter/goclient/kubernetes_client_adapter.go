package goclient

import (
	"context"

	"github.com/glasskube/glasskube/internal/adapter"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type clientSetKubernetesClientAdapter struct {
	clientset kubernetes.Clientset
}

// GetConfigMap implements adapter.KubernetesClientAdapter.
func (c *clientSetKubernetesClientAdapter) GetConfigMap(ctx context.Context, name string, namespace string) (
	*corev1.ConfigMap,
	error,
) {
	return c.clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
}

// GetSecret implements adapter.KubernetesClientAdapter.
func (c *clientSetKubernetesClientAdapter) GetSecret(ctx context.Context, name string, namespace string) (
	*corev1.Secret,
	error,
) {
	return c.clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
}

func NewKubernetesClientAdapter(clientset kubernetes.Clientset) adapter.KubernetesClientAdapter {
	return &clientSetKubernetesClientAdapter{clientset: clientset}
}
