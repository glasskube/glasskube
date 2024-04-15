package telemetry

const (
	apiKey   = "phc_EloQUW6cgfbTc0pI9c5CXElhQ4gVGRoBsrUAoakJVoQ" // TODO ??
	endpoint = "https://eu.posthog.com"
)

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
