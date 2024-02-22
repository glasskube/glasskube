package web

import "k8s.io/client-go/tools/clientcmd"

type ServerConfigError interface {
	error
	KubeconfigDefaultLocation() string
	KubeconfigMissing() bool
	BootstrapMissing() bool
}

type kubeconfigDefaultLocationSupplier struct{}

func (kubeconfigDefaultLocationSupplier) KubeconfigDefaultLocation() string {
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
	kubeconfigDefaultLocationSupplier
	wrappedErr
}

func (bootstrapErr) BootstrapMissing() bool {
	return true
}

func (bootstrapErr) KubeconfigMissing() bool {
	return false
}

func newBootstrapErr(cause error) ServerConfigError {
	return &bootstrapErr{wrappedErr: wrappedErr{cause}}
}

type kubeconfigErr struct {
	kubeconfigDefaultLocationSupplier
	wrappedErr
}

func (kubeconfigErr) BootstrapMissing() bool {
	return false
}

func (err kubeconfigErr) KubeconfigMissing() bool {
	return clientcmd.IsEmptyConfig(err.Cause)
}

func newKubeconfigErr(cause error) ServerConfigError {
	return &kubeconfigErr{wrappedErr: wrappedErr{cause}}
}
