package bootstrap

import (
	"errors"
	"fmt"
)

type BootstrapType string

const (
	BootstrapTypeSlim BootstrapType = "slim"
	BootstrapTypeAio  BootstrapType = "aio"
)

func (t *BootstrapType) String() string {
	return string(*t)
}

func (t *BootstrapType) Set(v string) error {
	switch v {
	case string(BootstrapTypeSlim), string(BootstrapTypeAio):
		*t = BootstrapType(v)
		return nil
	default:
		return errors.New(`must be one of "aio", "slim"`)
	}
}

func (e *BootstrapType) Type() string {
	return fmt.Sprintf("(%v|%v)", BootstrapTypeAio, BootstrapTypeSlim)
}
