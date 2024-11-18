package sse

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	"net/http"
	"reflect"

	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/web/sse/refresh"
)

type Broadcaster struct {
	sseHub *sseHub
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		sseHub: newHub(),
	}
}

func (b *Broadcaster) Run(stopCh chan struct{}) {
	b.sseHub.run(stopCh)
}

func (b *Broadcaster) Handler(w http.ResponseWriter, r *http.Request) {
	b.sseHub.handler(w)
}

func (b *Broadcaster) UpdatesAvailable(headerOnly refresh.RefreshTriggerHeaderOnly, pkgs ...ctrlpkg.Package) {
	pkgsOverviewDone := false
	clpkgsOverviewDone := false
	for _, pkg := range pkgs {
		b.sseHub.broadcast <- &sse{
			event: refresh.GetPackageRefreshDetailId(pkg, headerOnly),
		}

		// for each package scope, the overview trigger should sent at most once
		if pkg.IsNamespaceScoped() {
			if pkgsOverviewDone {
				continue
			}
			b.sseHub.broadcast <- &sse{
				event: refresh.RefreshPackageOverview,
			}
			pkgsOverviewDone = true
		} else {
			if clpkgsOverviewDone {
				continue
			}
			b.sseHub.broadcast <- &sse{
				event: refresh.RefreshClusterPackageOverview,
			}
			clpkgsOverviewDone = true
		}
	}
}

func (b *Broadcaster) UpdatesAvailableForPackage(oldPkg ctrlpkg.Package, newPkg ctrlpkg.Package) {
	// TODO rename and separate to three functions: package added, changed and removed
	oldPkgCopy := DeepCopyPackage(oldPkg)
	newPkgCopy := DeepCopyPackage(newPkg)
	if oldPkgCopy != nil && !oldPkgCopy.IsNil() && newPkgCopy != nil && !newPkgCopy.IsNil() {
		if !reflect.DeepEqual(oldPkgCopy.GetSpec(), newPkgCopy.GetSpec()) {
			b.UpdatesAvailable(refresh.RefreshTriggerAll, newPkgCopy)
		} else if !reflect.DeepEqual(oldPkgCopy.GetAnnotations(), newPkgCopy.GetAnnotations()) ||
			!reflect.DeepEqual(oldPkgCopy.GetLabels(), newPkgCopy.GetLabels()) {
			b.UpdatesAvailable(refresh.RefreshTriggerAll, newPkgCopy)
		} else if !reflect.DeepEqual(oldPkgCopy.GetStatus(), newPkgCopy.GetStatus()) {
			b.UpdatesAvailable(refresh.RefreshTriggerHeader, newPkgCopy)
		} else if !newPkgCopy.GetDeletionTimestamp().IsZero() {
			b.UpdatesAvailable(refresh.RefreshTriggerAll, newPkgCopy)
		}
	} else if oldPkgCopy != nil && !oldPkgCopy.IsNil() {
		b.UpdatesAvailable(refresh.RefreshTriggerAll, oldPkgCopy)
	} else if newPkgCopy != nil && !newPkgCopy.IsNil() {
		b.UpdatesAvailable(refresh.RefreshTriggerAll, newPkgCopy)
	}
}

func DeepCopyPackage(pkg ctrlpkg.Package) ctrlpkg.Package {
	if pkg == nil || pkg.IsNil() {
		return pkg
	}
	if nsPkg, ok := pkg.(*v1alpha1.Package); ok {
		return nsPkg.DeepCopy()
	} else if clPkg, ok := pkg.(*v1alpha1.ClusterPackage); ok {
		return clPkg.DeepCopy()
	} else {
		panic("unsupported package type")
	}
}
