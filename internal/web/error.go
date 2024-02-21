package web

import "k8s.io/client-go/tools/clientcmd"

type ServerConfigError interface {
	error
	KubeconfigDefaultLocation() string
	KubeconfigMissing() bool
	BootstrapMissing() bool
}

type defaultKubeConfigAccessor struct{}

func (defaultKubeConfigAccessor) KubeconfigDefaultLocation() string {
	return clientcmd.RecommendedHomeFile
}

type wrappedErr struct {
	Cause error
}

func (err wrappedErr) Error() string {
	return err.Cause.Error()
}

func (err wrappedErr) Unwrap() error {
	return err.Cause
}

type bootstrapErr struct {
	defaultKubeConfigAccessor
	wrappedErr
}

func (bootstrapErr) BootstrapMissing() bool {
	return true
}

func (bootstrapErr) KubeconfigMissing() bool {
	return false
}

func bootstrapError(err error) ServerConfigError {
	return &bootstrapErr{wrappedErr: wrappedErr{Cause: err}}
}

type kubeconfigErr struct {
	defaultKubeConfigAccessor
	wrappedErr
}

func (kubeconfigErr) BootstrapMissing() bool {
	return false
}

func (err kubeconfigErr) KubeconfigMissing() bool {
	return clientcmd.IsEmptyConfig(err.Cause)
}

func kubeconfigError(err error) ServerConfigError {
	return &kubeconfigErr{wrappedErr: wrappedErr{Cause: err}}
}
