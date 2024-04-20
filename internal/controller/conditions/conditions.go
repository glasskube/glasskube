package conditions

import (
	"context"
	"fmt"

	"github.com/glasskube/glasskube/internal/telemetry"
	"github.com/glasskube/glasskube/pkg/condition"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// SetInitialAndUpdate initializes the Conditions slice and updates the resource if the conditions slice is nil or empty.
// Types Ready and Failed are initialized with Status Unknown.
func SetInitialAndUpdate(ctx context.Context, client client.Client, obj client.Object, objConditions *[]metav1.Condition) error {
	if objConditions == nil || len(*objConditions) == 0 {
		_, err := SetUnknownAndUpdate(ctx, client, obj, objConditions, condition.Reconciling, "Starting reconciliation")
		return err
	}
	return nil
}

func SetUnknown(ctx context.Context, objConditions *[]metav1.Condition, reason condition.Reason, message string) bool {
	log := log.FromContext(ctx)
	log.V(1).Info("set condition to unknown: " + message)
	return setStatusConditions(objConditions,
		metav1.Condition{Type: string(condition.Ready), Status: metav1.ConditionUnknown, Reason: string(reason), Message: message},
		metav1.Condition{Type: string(condition.Failed), Status: metav1.ConditionUnknown, Reason: string(reason), Message: message},
	)
}

func SetUnknownAndUpdate(ctx context.Context, client client.Client, obj client.Object, objConditions *[]metav1.Condition, reason condition.Reason, message string) (bool, error) {
	if SetUnknown(ctx, objConditions, reason, message) {
		return true, updateAfterConditionsChanged(ctx, client, obj)
	}
	return false, nil
}

// SetReady sets the Ready condition to Status=True and the Failed condition to Status=False.
func SetReady(ctx context.Context, recorder record.EventRecorder, obj client.Object, objConditions *[]metav1.Condition, reason condition.Reason, message string) bool {
	log := log.FromContext(ctx)
	log.V(1).Info("set condition to ready: " + message)
	recorder.Event(obj, "Normal", string(reason), message)
	changed := setStatusConditions(objConditions,
		metav1.Condition{Type: string(condition.Ready), Status: metav1.ConditionTrue, Reason: string(reason), Message: message},
		metav1.Condition{Type: string(condition.Failed), Status: metav1.ConditionFalse, Reason: string(reason), Message: message},
	)
	if changed {
		telemetry.ForOperator().OnEvent(obj, condition.Ready, reason)
	}
	return changed
}

func SetReadyAndUpdate(ctx context.Context, client client.Client, recorder record.EventRecorder, obj client.Object, objConditions *[]metav1.Condition, reason condition.Reason, message string) error {
	if SetReady(ctx, recorder, obj, objConditions, reason, message) {
		return updateAfterConditionsChanged(ctx, client, obj)
	}
	return nil
}

func SetFailed(ctx context.Context, recorder record.EventRecorder, obj client.Object, objConditions *[]metav1.Condition, reason condition.Reason, message string) (bool, bool) {
	log := log.FromContext(ctx)
	log.V(1).Info("set condition to failed: " + message)
	recorder.Event(obj, "Warning", string(reason), message)
	changed := setStatusConditions(objConditions,
		metav1.Condition{Type: string(condition.Ready), Status: metav1.ConditionFalse, Reason: string(reason), Message: message},
		metav1.Condition{Type: string(condition.Failed), Status: metav1.ConditionTrue, Reason: string(reason), Message: message},
	)
	if changed {
		telemetry.ForOperator().OnEvent(obj, condition.Failed, reason)
	}
	return changed, reason.Recoverable()
}

// SetFailedAndUpdate sets the Ready condition to Status=False and the Failed condition to Status=True, then updates the resource if the conditions have changed.
func SetFailedAndUpdate(ctx context.Context, client client.Client, recorder record.EventRecorder, obj client.Object, objConditions *[]metav1.Condition, reason condition.Reason, message string) error {
	needsUpdate, _ := SetFailed(ctx, recorder, obj, objConditions, reason, message)
	if needsUpdate {
		return updateAfterConditionsChanged(ctx, client, obj)
	}

	return nil
}

func updateAfterConditionsChanged(ctx context.Context, cl client.Client, obj client.Object) error {
	log := log.FromContext(ctx)
	log.V(1).Info("Updating status after conditions changed")
	if err := cl.Status().Update(ctx, obj); err != nil {
		return fmt.Errorf("could not set conditions: failed to update object status: %w", err)
	}
	return nil
}

func setStatusConditions(statusConditions *[]metav1.Condition, newConditions ...metav1.Condition) bool {
	needsUpdate := false
	for _, condition := range newConditions {
		changed := meta.SetStatusCondition(statusConditions, condition)
		needsUpdate = changed || needsUpdate
	}
	return needsUpdate
}
