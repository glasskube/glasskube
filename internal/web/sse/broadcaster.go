package sse

import (
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

func (b *Broadcaster) Run() {
	b.sseHub.run()
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
	if oldPkg != nil && !oldPkg.IsNil() && newPkg != nil && !newPkg.IsNil() {
		if !reflect.DeepEqual(oldPkg.GetSpec(), newPkg.GetSpec()) {
			b.UpdatesAvailable(refresh.RefreshTriggerAll, newPkg)
		} else if !reflect.DeepEqual(oldPkg.GetStatus(), newPkg.GetStatus()) {
			b.UpdatesAvailable(refresh.RefreshTriggerHeader, newPkg)
		} else if !reflect.DeepEqual(oldPkg.GetAnnotations(), newPkg.GetAnnotations()) ||
			!reflect.DeepEqual(oldPkg.GetLabels(), newPkg.GetLabels()) {
			b.UpdatesAvailable(refresh.RefreshTriggerHeader, newPkg)
		} else if !newPkg.GetDeletionTimestamp().IsZero() {
			b.UpdatesAvailable(refresh.RefreshTriggerAll, newPkg)
		}
	} else if oldPkg != nil && !oldPkg.IsNil() {
		b.UpdatesAvailable(refresh.RefreshTriggerAll, oldPkg)
	} else if newPkg != nil && !newPkg.IsNil() {
		b.UpdatesAvailable(refresh.RefreshTriggerAll, newPkg)
	}
}
