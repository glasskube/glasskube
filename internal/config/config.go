package config

var (
	Kubeconfig        string
	ForceUninstall    bool
	ListInstalledOnly bool
	Verbose           bool
	Version           = "dev"
	Commit            = "none"
	Date              = "unknown"
)
