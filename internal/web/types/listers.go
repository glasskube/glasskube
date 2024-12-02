package types

import (
	appsv1 "k8s.io/client-go/listers/apps/v1"
	v1 "k8s.io/client-go/listers/core/v1"
)

type CoreListers struct {
	NamespaceLister  *v1.NamespaceLister
	ConfigMapLister  *v1.ConfigMapLister
	SecretLister     *v1.SecretLister
	DeploymentLister *appsv1.DeploymentLister
}
