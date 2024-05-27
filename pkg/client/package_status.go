package client

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/pkg/condition"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PackageStatus struct {
	Status  string
	Reason  string
	Message string
}

func GetStatus(status *v1alpha1.PackageStatus) *PackageStatus {
	readyCnd := meta.FindStatusCondition((*status).Conditions, string(condition.Ready))
	if readyCnd != nil && readyCnd.Status == metav1.ConditionTrue {
		return newPackageStatus(readyCnd)
	}
	failedCnd := meta.FindStatusCondition((*status).Conditions, string(condition.Failed))
	if failedCnd != nil && failedCnd.Status == metav1.ConditionTrue {
		return newPackageStatus(failedCnd)
	}
	return nil
}

func GetStatusOrPending(pkg *v1alpha1.Package) *PackageStatus {
	if pkg != nil {
		if !pkg.DeletionTimestamp.IsZero() {
			return NewUninstallingStatus()
		}
		if status := GetStatus(&pkg.Status); status != nil {
			return status
		} else {
			return NewPendingStatus()
		}
	} else {
		return nil
	}
}

func newPackageStatus(cnd *metav1.Condition) *PackageStatus {
	return &PackageStatus{
		Status:  cnd.Type,
		Reason:  cnd.Reason,
		Message: cnd.Message,
	}
}

func NewPendingStatus() *PackageStatus {

	return &PackageStatus{Status: "Pending"}
}

func NewUninstallingStatus() *PackageStatus {
	return &PackageStatus{Status: "Uninstalling"}
}
