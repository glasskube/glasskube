package controllerruntime

import (
	"context"

	"github.com/glasskube/glasskube/internal/adapter"
	corev1 "k8s.io/api/core/v1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ctrlKubernetsClientAdapter struct {
	client ctrlclient.Client
}

// GetConfigMap implements adapter.KubernetesClientAdapter.
func (c *ctrlKubernetsClientAdapter) GetConfigMap(ctx context.Context, name string, namespace string) (
	*corev1.ConfigMap,
	error,
) {
	var cm corev1.ConfigMap
	if err := c.client.Get(ctx, ctrlclient.ObjectKey{Name: name, Namespace: namespace}, &cm); err != nil {
		return nil, err
	} else {
		return &cm, nil
	}
}

// GetSecret implements adapter.KubernetesClientAdapter.
func (c *ctrlKubernetsClientAdapter) GetSecret(ctx context.Context, name string, namespace string) (
	*corev1.Secret,
	error,
) {
	var s corev1.Secret
	if err := c.client.Get(ctx, ctrlclient.ObjectKey{Name: name, Namespace: namespace}, &s); err != nil {
		return nil, err
	} else {
		return &s, nil
	}
}

func NewKubernetesClientAdapter(client ctrlclient.Client) adapter.KubernetesClientAdapter {
	return &ctrlKubernetsClientAdapter{client: client}
}
