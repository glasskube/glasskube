package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

type repoAddAuthType string

func (t *repoAddAuthType) String() string {
	return string(*t)
}

func (t *repoAddAuthType) Set(v string) error {
	switch v {
	case string(repoAddBasicAuth), string(repoAddBearerAuth):
		*t = repoAddAuthType(v)
		return nil
	default:
		return errors.New(`must be one of "basic", "bearer"`)
	}
}

func (e *repoAddAuthType) Type() string {
	return fmt.Sprintf("[%v|%v]", repoAddBasicAuth, repoAddBearerAuth)
}

const (
	repoAddBasicAuth  repoAddAuthType = "basic"
	repoAddBearerAuth repoAddAuthType = "bearer"
	repoAddNoAuth     repoAddAuthType = ""
)

type repoOptions struct {
	Default  bool
	Auth     repoAddAuthType
	Username string
	Password string
	Token    string
}

func (opts *repoOptions) BindToCmdFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&opts.Default, "default", opts.Default, "use this repository as default")
	cmd.Flags().Var(&opts.Auth, "auth", "type of authentication")
	cmd.Flags().StringVar(&opts.Username, "username", opts.Username, "username for basic authentication")
	cmd.Flags().StringVar(&opts.Password, "password", opts.Password, "password for basic authentication")
	cmd.Flags().StringVar(&opts.Token, "token", opts.Token, "token for bearer authentication")
	cmd.MarkFlagsMutuallyExclusive("username", "token")
	cmd.MarkFlagsMutuallyExclusive("password", "token")
}

func (opts *repoOptions) Normalize() error {
	if len(opts.Username) > 0 || len(opts.Password) > 0 {
		if opts.Auth != repoAddNoAuth {
			return fmt.Errorf("username/password only applies to basic authentication (got %v)", opts.Auth)
		} else {
			opts.Auth = repoAddBasicAuth
		}
	} else if len(opts.Token) > 0 {
		if opts.Auth != repoAddNoAuth {
			return fmt.Errorf("token only applies to bearer authentication (got %v)", opts.Auth)
		} else {
			opts.Auth = repoAddBearerAuth
		}
	}
	return nil
}
