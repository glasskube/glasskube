package watch

import (
	"context"

	ownerutils "github.com/glasskube/glasskube/internal/controller/owners/utils"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func EnqueueRequestsFromOwnedResource(
	scheme *runtime.Scheme,
	targetLister PackageLister,
	ownedResourcesGetter ownedMapperFunc,
) handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
		log := ctrl.LoggerFrom(ctx)
		objRef, err := ownerutils.ToOwnedResourceRef(scheme, obj)
		if err != nil {
			log.Error(err, "could not map object in event handler")
			return nil
		}
		if pkgs, err := targetLister.ListPackages(ctx); err != nil {
			log.Error(err, "could not list packages event handler")
			return nil
		} else {
			var res []reconcile.Request
			for _, pkg := range pkgs {
				for _, ownedRef := range ownedResourcesGetter(pkg) {
					if ownerutils.RefersToSameResource(ownedRef, objRef) {
						res = append(res, reconcile.Request{NamespacedName: client.ObjectKeyFromObject(pkg)})
						break
					}
				}
			}
			return res
		}
	})
}
