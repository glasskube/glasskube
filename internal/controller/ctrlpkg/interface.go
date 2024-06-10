package ctrlpkg

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type PackageCommon interface {
	metav1.Object
	runtime.Object
	AutoUpdatesEnabled() bool
	SetAutoUpdatesEnabled(enabled bool)
	GetSpec() *v1alpha1.PackageSpec
	GetStatus() *v1alpha1.PackageStatus
}
