package client

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/adapter"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/repo/client/auth"
	corev1 "k8s.io/api/core/v1"
)

type defaultClientsetClient struct {
	adapter.PackageClientAdapter
	adapter.KubernetesClientAdapter
}

type defaultClientset struct {
	client                  defaultClientsetClient
	clients                 map[string]repoClientWithState
	repoWithNameMutex       sync.Mutex
	repoMutex               sync.Mutex
	maxCacheAge             time.Duration
	clientInfoCheckInterval time.Duration
}

var _ RepoClientset = &defaultClientset{}

type repoClientWithState struct {
	client              RepoClient
	lastCheckedRepoSpec time.Time
	repo                v1alpha1.PackageRepository
}

func (s *repoClientWithState) lastCheckedRepoSpecAgo(maxAge time.Duration) bool {
	return s.lastCheckedRepoSpec.Add(maxAge).After(time.Now())
}

func (s *repoClientWithState) checkRepoSpec(repo v1alpha1.PackageRepository) bool {
	return reflect.DeepEqual(s.repo.Spec, repo.Spec)
}

func NewClientset(pkgClient adapter.PackageClientAdapter, k8sClient adapter.KubernetesClientAdapter) RepoClientset {
	return NewClientsetWithMaxCacheAge(pkgClient, k8sClient, 30*time.Second, 5*time.Minute)
}

func NewClientsetWithMaxCacheAge(pkgClient adapter.PackageClientAdapter, k8sClient adapter.KubernetesClientAdapter,
	clientInfoCheckInterval time.Duration,
	maxCacheAge time.Duration) RepoClientset {
	return &defaultClientset{
		client:                  defaultClientsetClient{pkgClient, k8sClient},
		clients:                 make(map[string]repoClientWithState),
		maxCacheAge:             maxCacheAge,
		clientInfoCheckInterval: clientInfoCheckInterval,
	}
}

// ForPackage implements RepoClientset.
func (d *defaultClientset) ForPackage(pkg ctrlpkg.Package) RepoClient {
	return d.ForRepoWithName(pkg.GetSpec().PackageInfo.RepositoryName)
}

// ForRepo implements RepoClientset.
func (d *defaultClientset) ForRepoWithName(name string) RepoClient {
	d.repoWithNameMutex.Lock()
	defer d.repoWithNameMutex.Unlock()
	if clientState, ok := d.clients[name]; ok && clientState.lastCheckedRepoSpecAgo(d.clientInfoCheckInterval) {
		return clientState.client
	}
	if len(name) > 0 {
		if repo, err := d.client.GetPackageRepository(context.TODO(), name); err != nil {
			return &errorclient{err: err}
		} else {
			return d.ForRepo(*repo)
		}
	} else {
		return d.Default()
	}
}

// Default implements RepoClientset.
func (d *defaultClientset) Default() RepoClient {
	if repos, err := d.client.ListPackageRepositories(context.TODO()); err != nil {
		return &errorclient{err: err}
	} else {
		for _, repo := range repos.Items {
			if repo.IsDefaultRepository() {
				return d.ForRepo(repo)
			}
		}
		return &errorclient{err: errors.New("default repository not found")}
	}
}

// ForRepo implements RepoClientset.
func (d *defaultClientset) ForRepo(repo v1alpha1.PackageRepository) RepoClient {
	d.repoMutex.Lock()
	defer d.repoMutex.Unlock()
	if clientState, ok := d.clients[repo.Name]; ok && clientState.checkRepoSpec(repo) {
		clientState.lastCheckedRepoSpec = time.Now()
		return clientState.client
	} else {
		if auth, err := d.newAuthenticator(repo); err != nil {
			return &errorclient{err: fmt.Errorf("invalid auth config: %w", err)}
		} else {
			client := New(repo.Spec.Url, auth, d.maxCacheAge)
			d.clients[repo.Name] = repoClientWithState{
				client:              client,
				lastCheckedRepoSpec: time.Now(),
				repo:                repo,
			}
			return client
		}
	}
}

func (d *defaultClientset) newAuthenticator(repo v1alpha1.PackageRepository) (auth.Authenticator, error) {
	if repo.Spec.Auth != nil {
		if repo.Spec.Auth.Basic != nil {
			var user, pass string
			var userSecret *corev1.Secret
			if repo.Spec.Auth.Basic.Username != nil {
				user = *repo.Spec.Auth.Basic.Username
			} else if repo.Spec.Auth.Basic.UsernameSecretRef != nil {
				if s, err := d.client.GetSecret(context.TODO(),
					repo.Spec.Auth.Basic.UsernameSecretRef.Name, "glasskube-system"); err != nil {
					return nil, fmt.Errorf("cannot get username: %w", err)
				} else {
					userSecret = s
				}
				if u, err := getKeyFromSecret(userSecret, repo.Spec.Auth.Basic.UsernameSecretRef.Key); err != nil {
					return nil, fmt.Errorf("cannot get username: %w", err)
				} else {
					user = u
				}
			}
			if repo.Spec.Auth.Basic.Password != nil {
				pass = *repo.Spec.Auth.Basic.Password
			} else if repo.Spec.Auth.Basic.PasswordSecretRef != nil {
				passSecret := userSecret
				if passSecret == nil || passSecret.Name != repo.Spec.Auth.Basic.PasswordSecretRef.Name {
					if s, err := d.client.GetSecret(context.TODO(),
						repo.Spec.Auth.Basic.PasswordSecretRef.Name, "glasskube-system"); err != nil {
						return nil, fmt.Errorf("cannot get password: %w", err)
					} else {
						passSecret = s
					}
					if p, err := getKeyFromSecret(passSecret, repo.Spec.Auth.Basic.PasswordSecretRef.Key); err != nil {
						return nil, fmt.Errorf("cannot get password: %w", err)
					} else {
						pass = p
					}
				}
			}
			return auth.Basic(user, pass), nil
		} else if repo.Spec.Auth.Bearer != nil {
			var token string
			if repo.Spec.Auth.Bearer.Token != nil {
				token = *repo.Spec.Auth.Bearer.Token
			} else if repo.Spec.Auth.Bearer.TokenSecretRef != nil {
				if tokenSecret, err := d.client.GetSecret(context.TODO(),
					repo.Spec.Auth.Bearer.TokenSecretRef.Name, "glasskube-system"); err != nil {
					return nil, fmt.Errorf("cannot get bearer token: %w", err)
				} else if t, err := getKeyFromSecret(tokenSecret, repo.Spec.Auth.Bearer.TokenSecretRef.Key); err != nil {
					return nil, fmt.Errorf("cannot get bearer token: %w", err)
				} else {
					token = t
				}
			}
			return auth.Bearer(token), nil
		}
	}
	return auth.Noop(), nil
}

// Meta implements RepoClientset.
func (d *defaultClientset) Meta() RepoMetaclient {
	return metaclient{clientset: d}
}

func getKeyFromSecret(secret *corev1.Secret, key string) (string, error) {
	if enc, ok := secret.Data[key]; ok {
		var dec []byte
		if _, err := base64.StdEncoding.Decode(dec, enc); err != nil {
			return "", err
		} else {
			return string(dec), nil
		}
	} else {
		return "", fmt.Errorf("%v has no key %v", secret.Name, key)
	}
}
