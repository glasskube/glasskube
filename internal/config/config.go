package config

const defaultVersion = "dev"

var (
	Kubeconfig     string
	NonInteractive bool
	Version        = defaultVersion
	Commit         = "none"
	Date           = "unknown"
)

func IsDevBuild() bool {
	return Version == defaultVersion
}
