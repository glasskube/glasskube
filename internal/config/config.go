package config

const defaultVersion = "dev"

var (
	Kubeconfig string
	Version    = defaultVersion
	Commit     = "none"
	Date       = "unknown"
)

func IsDevBuild() bool {
	return Version == defaultVersion
}
