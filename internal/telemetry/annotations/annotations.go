package annotations

import (
	"strconv"

	"k8s.io/apimachinery/pkg/util/uuid"
)

const (
	TelemetryIdAnnotation      = "packages.glasskube.dev/telemetry-id"
	TelemetryEnabledAnnotation = "packages.glasskube.dev/telemetry-enabled"
	GitopsModeEnabled          = "packages.glasskube.dev/gitops-mode-enabled"
)

func IsTelemetryEnabled(annotations map[string]string) bool {
	enabledAnnotation, _ := strconv.ParseBool(annotations[TelemetryEnabledAnnotation])
	return enabledAnnotation
}

func UpdateTelemetryAnnotations(currentAnnotations map[string]string, disabled bool) {
	enabledAnnotationStr, hasAnnotation := currentAnnotations[TelemetryEnabledAnnotation]
	enabledAnnotation, _ := strconv.ParseBool(enabledAnnotationStr)

	if disabled {
		currentAnnotations[TelemetryEnabledAnnotation] = strconv.FormatBool(false)
	} else if !hasAnnotation || enabledAnnotation {
		currentAnnotations[TelemetryEnabledAnnotation] = strconv.FormatBool(true)
		if _, hasId := currentAnnotations[TelemetryIdAnnotation]; !hasId {
			currentAnnotations[TelemetryIdAnnotation] = string(uuid.NewUUID())
		}
	}
	/*
	 * !hasAnnotation && !disabled: create annotation + ID if not exists
	 * !hasAnnotation && disabled: create annotation and set to false
	 * hasAnnotation && enabled && !disabled: create annotation + ID if not exists
	 * hasAnnotation && disabled && !disabled: do not change anything (update of a telemetry-disabled installation)
	 * hasAnnotation && enabled && disabled: set annotation to false
	 * hasAnnotation && disabled && disabled: do not change anything (update of a telemetry-disabled installation with explicit --disable too)
	 */
}

func IsGitopsModeEnabled(annotations map[string]string) bool {
	enabledAnnotation, _ := strconv.ParseBool(annotations[GitopsModeEnabled])
	return enabledAnnotation
}
