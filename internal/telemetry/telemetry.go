package telemetry

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/denisbrodbeck/machineid"
	"k8s.io/apimachinery/pkg/util/uuid"
)

var apiKey = "phc_EloQUW6cgfbTc0pI9c5CXElhQ4gVGRoBsrUAoakJVoQ" // TODO ??
const TelemetryIdAnnotation = "packages.glasskube.dev/telemetry-id"
const TelemetryEnabledAnnotation = "packages.glasskube.dev/telemetry-enabled"

func getMachineId() string {
	id, err := machineid.ProtectedID("glasskube")
	if err != nil {
		return fmt.Sprintf("fallback-id-%v-%v", rand.Int(), err.Error())
	}
	return id
}

func extractFromArgs(args []string) ([]string, []string) {
	var arguments = make([]string, 0)
	flags := make([]string, 0)
	if len(args) > 0 {
		noMoreArguments := false
		for _, arg := range args {
			if strings.HasPrefix(arg, "-") {
				noMoreArguments = true
				flags = append(flags, arg)
			} else if !noMoreArguments {
				arguments = append(arguments, arg)
			}
		}
	}
	return arguments, flags
}

// TODO
/*
func ForWeb() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			defer func() {
				if client := getClient(); client != nil {
					client.posthog.Enqueue(posthog.Capture{
						DistinctId: id,
						Type:       "web",
						Event:      "invoke_endpoint",
						Properties: map[string]any{
							"$current_url": r.URL.String(),
							"method":       r.Method,
							"path":         r.URL,
						},
					})
				}
			}()
		})
	}
}
*/

func IsEnabled(annotations map[string]string) (bool, string) {
	enabledAnnotationStr := annotations[TelemetryEnabledAnnotation]
	enabledAnnotation, _ := strconv.ParseBool(enabledAnnotationStr)
	if enabledAnnotation {
		return enabledAnnotation, annotations[TelemetryIdAnnotation]
	} else {
		return false, ""
	}
}

func UpdateAnnotations(currentAnnotations map[string]string, disabled bool) {
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
