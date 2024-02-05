package requeue

import (
	"context"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	RequeueDuration      = 60 * time.Second
	ErrorRequeueDuration = 30 * time.Second
)

func Always(ctx context.Context, err error) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	if err != nil {
		log.Info("error during reconciliation: " + err.Error())
		return requeueAfter(ErrorRequeueDuration), nil
	}
	log.V(1).Info("reconciliation finished")
	return requeueAfter(RequeueDuration), nil
}

func requeueAfter(duration time.Duration) ctrl.Result {
	return ctrl.Result{RequeueAfter: duration}
}
