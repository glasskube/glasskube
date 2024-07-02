package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/spf13/cobra"
)

type repoAuthType string

func (t *repoAuthType) String() string {
	return string(*t)
}

func (t *repoAuthType) Set(v string) error {
	switch v {
	case string(repoNoAuth), string(repoBasicAuth), string(repoBearerAuth):
		*t = repoAuthType(v)
		return nil
	default:
		return errors.New(`must be one of "none", "basic", "bearer"`)
	}
}

func (e *repoAuthType) Type() string {
	return fmt.Sprintf("(%v|%v|%v)", repoNoAuth, repoBasicAuth, repoBearerAuth)
}

const (
	repoBasicAuth  repoAuthType = "basic"
	repoBearerAuth repoAuthType = "bearer"
	repoNoAuth     repoAuthType = "none"
)

type repoOptions struct {
	Default  bool
	Auth     repoAuthType
	Username string
	Password string
	Token    string
	Url      string
}

func (opts *repoOptions) BindToCmdFlags(cmd *cobra.Command, update bool) {
	cmd.Flags().BoolVar(&opts.Default, "default", opts.Default, "use this repository as default")
	cmd.Flags().Var(&opts.Auth, "auth", "type of authentication")
	if update {
		cmd.Flags().StringVar(&opts.Url, "url", opts.Url, "new url for the repository")
	}
	cmd.Flags().StringVar(&opts.Username, "username", opts.Username, "username for basic authentication")
	cmd.Flags().StringVar(&opts.Password, "password", opts.Password, "password for basic authentication")
	cmd.Flags().StringVar(&opts.Token, "token", opts.Token, "token for bearer authentication")
	cmd.MarkFlagsMutuallyExclusive("username", "token")
	cmd.MarkFlagsMutuallyExclusive("password", "token")
}

func (opts *repoOptions) Normalize() error {
	if len(opts.Url) > 0 {
		if _, err := url.ParseRequestURI(opts.Url); err != nil {
			return fmt.Errorf("use a valid URL for the package repository (got %v)", opts.Url)
		}
	}

	if len(opts.Username) > 0 || len(opts.Password) > 0 {
		if opts.Auth == repoNoAuth || opts.Auth == repoBearerAuth {
			return fmt.Errorf("username/password only applies to basic authentication (got %v)", opts.Auth)
		} else {
			opts.Auth = repoBasicAuth
		}
	} else if len(opts.Token) > 0 {
		if opts.Auth == repoNoAuth || opts.Auth == repoBasicAuth {
			return fmt.Errorf("token only applies to bearer authentication (got %v)", opts.Auth)
		} else {
			opts.Auth = repoBearerAuth
		}
	}
	return nil
}

func (opts *repoOptions) SetAuth() *v1alpha1.PackageRepositoryAuthSpec {
	switch opts.Auth {
	case repoBasicAuth:
		if len(opts.Username) == 0 {
			fmt.Fprintln(os.Stderr, "Basic authentication was requested. Please enter a username:")
			for {
				username := cliutils.GetInputStr("username")
				if len(username) > 0 {
					opts.Username = username
					break
				}
			}
		}
		if len(opts.Password) == 0 {
			fmt.Fprintln(os.Stderr, "Basic authentication was requested. Please enter a password:")
			for {
				password := cliutils.GetInputStr("password")
				if len(password) > 0 {
					opts.Password = password
					break
				}
			}
		}
		return &v1alpha1.PackageRepositoryAuthSpec{
			Basic: &v1alpha1.PackageRepositoryBasicAuthSpec{
				Username: &opts.Username,
				Password: &opts.Password,
			},
		}
	case repoBearerAuth:
		if len(opts.Token) == 0 {
			fmt.Fprintln(os.Stderr, "Bearer authentication was requested. Please enter a token:")
			for {
				token := cliutils.GetInputStr("token")
				if len(token) > 0 {
					opts.Token = token
					break
				}
			}
		}
		return &v1alpha1.PackageRepositoryAuthSpec{
			Bearer: &v1alpha1.PackageRepositoryBearerAuthSpec{
				Token: &opts.Token,
			},
		}
	}
	return nil
}
