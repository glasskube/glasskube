package client

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/adapter"
	"github.com/glasskube/glasskube/internal/maputils"
	"github.com/glasskube/glasskube/internal/repo/types"
	"github.com/glasskube/glasskube/internal/semver"
	"github.com/glasskube/glasskube/internal/util"
	"go.uber.org/multierr"
	corev1 "k8s.io/api/core/v1"
)

type defaultClientsetClient struct {
	adapter.PackageClientAdapter
	adapter.KubernetesClientAdapter
}

type defaultClientset struct {
	client            defaultClientsetClient
	clients           map[string]RepoClient
	repoWithNameMutex sync.Mutex
	repoMutex         sync.Mutex
	maxCacheAge       time.Duration
}

var _ RepoClientset = &defaultClientset{}
var _ RepoAggregator = &defaultClientset{}

func NewClientset(pkgClient adapter.PackageClientAdapter, k8sClient adapter.KubernetesClientAdapter) RepoClientset {
	return NewClientsetWithMaxCacheAge(pkgClient, k8sClient, 5*time.Minute)
}

func NewClientsetWithMaxCacheAge(pkgClient adapter.PackageClientAdapter, k8sClient adapter.KubernetesClientAdapter,
	maxCacheAge time.Duration) RepoClientset {
	return &defaultClientset{
		client:      defaultClientsetClient{pkgClient, k8sClient},
		maxCacheAge: maxCacheAge,
		clients:     make(map[string]RepoClient),
	}
}

// ForPackage implements RepoClientset.
func (d *defaultClientset) ForPackage(pkg v1alpha1.Package) RepoClient {
	return d.ForRepoWithName(pkg.Spec.PackageInfo.RepositoryName)
}

// ForRepo implements RepoClientset.
func (d *defaultClientset) ForRepoWithName(name string) RepoClient {
	d.repoWithNameMutex.Lock()
	defer d.repoWithNameMutex.Unlock()
	if client, ok := d.clients[name]; ok {
		// TODO: update client details if older than maxCacheAge
		return client
	}
	if len(name) > 0 {
		if repo, err := d.client.GetPackageRepository(context.TODO(), name); err != nil {
			return &errorclient{err}
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
		return &errorclient{err}
	} else {
		for _, repo := range repos.Items {
			if repo.IsDefaultRepository() {
				return d.ForRepo(repo)
			}
		}
		return &errorclient{errors.New("default repository not found")}
	}
}

// ForRepo implements RepoClientset.
func (d *defaultClientset) ForRepo(repo v1alpha1.PackageRepository) RepoClient {
	d.repoMutex.Lock()
	defer d.repoMutex.Unlock()
	if client, ok := d.clients[repo.Name]; ok {
		// TODO: update client details if older than maxCacheAge
		return client
	} else {
		if headers, err := d.getAuthHeaders(repo); err != nil {
			return &errorclient{fmt.Errorf("invalid auth config: %w", err)}
		} else {
			client := New(repo.Spec.Url, headers, d.maxCacheAge)
			d.clients[repo.Name] = client
			return client
		}
	}
}

func (d *defaultClientset) getAuthHeaders(repo v1alpha1.PackageRepository) (http.Header, error) {
	headers := http.Header{}
	if repo.Spec.Auth != nil {
		if repo.Spec.Auth.Basic != nil {
			user := repo.Spec.Auth.Basic.Username
			var userSecret *corev1.Secret
			if len(user) == 0 {
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
			pass := repo.Spec.Auth.Basic.Password
			if len(pass) == 0 {
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
			userpass := fmt.Sprintf("%v:%v", user, pass)
			userpassEncoded := base64.StdEncoding.EncodeToString([]byte(userpass))
			headers.Set("Authorization", fmt.Sprintf("Basic %v", userpassEncoded))
		} else if repo.Spec.Auth.Bearer != nil {
			token := repo.Spec.Auth.Bearer.Token
			if len(token) == 0 {
				if tokenSecret, err := d.client.GetSecret(context.TODO(),
					repo.Spec.Auth.Bearer.TokenSecretRef.Name, "glasskube-system"); err != nil {
					return nil, fmt.Errorf("cannot get bearer token: %w", err)
				} else if t, err := getKeyFromSecret(tokenSecret, repo.Spec.Auth.Bearer.TokenSecretRef.Key); err != nil {
					return nil, fmt.Errorf("cannot get bearer token: %w", err)
				} else {
					token = t
				}
			}
			headers.Set("Authorization", fmt.Sprintf("Bearer %v", token))
		}
	}
	return headers, nil
}

// Aggregate implements RepoClientset.
func (d *defaultClientset) Aggregate() RepoAggregator {
	return d
}

// FetchPackageRepoIndex implements RepoAggregator.
func (d *defaultClientset) FetchPackageRepoIndex(target *types.PackageRepoIndex) error {
	if repoList, err := d.client.ListPackageRepositories(context.TODO()); err != nil {
		return err
	} else {
		var compositeErr error
		indexMap := make(map[string]types.PackageRepoIndexItem)
		util.SortBy(repoList.Items, func(repo v1alpha1.PackageRepository) string { return repo.Name })
		slices.Reverse(repoList.Items)
		for _, repo := range repoList.Items {
			var index types.PackageRepoIndex
			if err := d.ForRepo(repo).FetchPackageRepoIndex(&index); err != nil {
				multierr.AppendInto(&compositeErr, err)
			} else {
				for _, item := range index.Packages {
					if _, ok := indexMap[item.Name]; !ok || !repo.IsDefaultRepository() {
						indexMap[item.Name] = item
					}
				}
			}
		}
		*target = types.PackageRepoIndex{
			Packages: make([]types.PackageRepoIndexItem, len(indexMap)),
		}
		for i, name := range maputils.KeysSorted(indexMap) {
			target.Packages[i] = indexMap[name]
		}
		return compositeErr
	}
}

// GetReposForPackage implements RepoAggregator.
func (d *defaultClientset) GetReposForPackage(name string) ([]v1alpha1.PackageRepository, error) {
	if repoList, err := d.client.ListPackageRepositories(context.TODO()); err != nil {
		return nil, err
	} else {
		var result []v1alpha1.PackageRepository
		for _, repo := range repoList.Items {
			var index types.PackageRepoIndex
			if err := d.ForRepo(repo).FetchPackageRepoIndex(&index); err != nil {
				return nil, err
			}
			if slices.ContainsFunc(index.Packages, func(item types.PackageRepoIndexItem) bool { return item.Name == name }) {
				result = append(result, repo)
			}
		}
		return result, nil
	}
}

// GetLatestVersion implements RepoAggregator.
func (d *defaultClientset) GetLatestVersion(pkgName string) (string, error) {
	if repoList, err := d.client.ListPackageRepositories(context.TODO()); err != nil {
		return "", err
	} else {
		var latest string
		for _, repo := range repoList.Items {
			var index types.PackageIndex
			if err := d.ForRepo(repo).FetchPackageIndex(pkgName, &index); err != nil {
				return "", err
			}
			if latest == "" || semver.IsUpgradable(latest, index.LatestVersion) {
				latest = index.LatestVersion
			}
		}
		return latest, nil
	}
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
