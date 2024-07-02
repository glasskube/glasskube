package web

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type cacheVerifier struct {
	lister     verifyLister
	diffMap    sync.Map
	restConfig *rest.Config
	stopCh     chan struct{}
	errCh      chan error
}

type diffMapEntry struct {
	cachedResourceVersion    string
	nonCachedResourceVersion string
}

type cacheVerificationError struct {
	itemName                 string
	cachedResourceVersion    string
	nonCachedResourceVersion string
}

func (e *cacheVerificationError) Error() string {
	return fmt.Sprintf("cache verification for %s failed two times in a row "+
		"(cached resource version is '%s', but expected '%s')", e.itemName, e.cachedResourceVersion, e.nonCachedResourceVersion)
}

func newVerifier(restConfig *rest.Config, lister verifyLister) *cacheVerifier {
	return &cacheVerifier{
		restConfig: restConfig,
		stopCh:     make(chan struct{}),
		errCh:      make(chan error, 1),
		lister:     lister,
	}
}

func (verifier *cacheVerifier) listItems(ctx context.Context, client client.PackageV1Alpha1Client) ([]metav1.Object, error) {
	return verifier.lister(ctx, client)
}

type verifyLister = func(ctx context.Context, client client.PackageV1Alpha1Client) ([]metav1.Object, error)

var clusterPackageVerifyLister = func(ctx context.Context, client client.PackageV1Alpha1Client) ([]metav1.Object, error) {
	var ls v1alpha1.ClusterPackageList
	err := client.ClusterPackages().GetAll(ctx, &ls)
	target := make([]metav1.Object, len(ls.Items))
	for idx, item := range ls.Items {
		target[idx] = item.GetObjectMeta()
	}
	return target, err
}

var packageVerifyLister = func(ctx context.Context, client client.PackageV1Alpha1Client) ([]metav1.Object, error) {
	var ls v1alpha1.PackageList
	err := client.Packages("").GetAll(ctx, &ls)
	target := make([]metav1.Object, len(ls.Items))
	for idx, item := range ls.Items {
		target[idx] = item.GetObjectMeta()
	}
	return target, err
}

var packageInfoVerifyLister = func(ctx context.Context, client client.PackageV1Alpha1Client) ([]metav1.Object, error) {
	var ls v1alpha1.PackageInfoList
	err := client.PackageInfos().GetAll(ctx, &ls)
	target := make([]metav1.Object, len(ls.Items))
	for idx, item := range ls.Items {
		target[idx] = item.GetObjectMeta()
	}
	return target, err
}

var packageRepoVerifyLister = func(ctx context.Context, client client.PackageV1Alpha1Client) ([]metav1.Object, error) {
	var ls v1alpha1.PackageRepositoryList
	err := client.PackageRepositories().GetAll(ctx, &ls)
	target := make([]metav1.Object, len(ls.Items))
	for idx, item := range ls.Items {
		target[idx] = item.GetObjectMeta()
	}
	return target, err
}

func (verifier *cacheVerifier) start(ctx context.Context, cachedClient client.PackageV1Alpha1Client, minInterval float32) {
	// setup client without cache
	nonCachedClient := client.NewOrDie(verifier.restConfig)

	var diffErr *cacheVerificationError

endlessLoop:
	for {
		time.Sleep(time.Second * time.Duration(minInterval*(rand.Float32()+1))) // between [minInterval, 2*minInterval) seconds

		if nonCachedItems, err := verifier.listItems(ctx, nonCachedClient); err != nil {
			fmt.Fprintf(os.Stderr, "CACHEVERIFIER: failed to get actual items: %v\n", err)
		} else if cachedItems, err := verifier.listItems(ctx, cachedClient); err != nil {
			fmt.Fprintf(os.Stderr, "CACHEVERIFIER: failed get cached items: %v\n", err)
		} else {
		outerLoop:
			for _, nonCachedItem := range nonCachedItems {
				for _, cachedItem := range cachedItems {
					if nonCachedItem.GetNamespace() == cachedItem.GetNamespace() &&
						nonCachedItem.GetName() == cachedItem.GetName() {
						if nonCachedItem.GetResourceVersion() != cachedItem.GetResourceVersion() {
							fmt.Fprintf(os.Stderr, "CACHEVERIFIER: resource versions of %s differ, need to check\n",
								cachedItem.GetName())
							if err := verifier.handleDiff(cachedItem, nonCachedItem); errors.As(err, &diffErr) {
								break endlessLoop
							} else if err != nil {
								fmt.Fprintf(os.Stderr, "CACHEVERIFIER: failed to handle diffs: %v\n", err)
							}
						} else {
							// previously stored inconsistency of this item can now be deleted
							verifier.diffMap.Delete(cache.MetaObjectToName(cachedItem).String())
						}
						continue outerLoop
					}
				}
				fmt.Fprintf(os.Stderr, "CACHEVERIFIER: non-cached item %s not found in cache, need to check\n",
					nonCachedItem.GetName())
				// non cached package found, that is not in cache (yet?)
				if err := verifier.handleDiff(nil, nonCachedItem); errors.As(err, &diffErr) {
					break endlessLoop
				} else if err != nil {
					fmt.Fprintf(os.Stderr, "CACHEVERIFIER: failed to handle diffs: %v\n", err)
				}
			}

			// check the other way around (i.e. whether there are cached items, that are not in the non-cached list anymore)
			for _, cachedClPkg := range cachedItems {
				found := false
				for _, nonCachedClPkg := range nonCachedItems {
					if nonCachedClPkg.GetNamespace() == cachedClPkg.GetNamespace() &&
						nonCachedClPkg.GetName() == cachedClPkg.GetName() {
						found = true
						break
					}
				}
				if !found {
					fmt.Fprintf(os.Stderr, "CACHEVERIFIER: cached item %s not found in non-cached list, need to check\n",
						cachedClPkg.GetName())
					// when here, there exists a cached item, that is not in the non-cached list
					if err := verifier.handleDiff(cachedClPkg, nil); errors.As(err, &diffErr) {
						break endlessLoop
					} else if err != nil {
						fmt.Fprintf(os.Stderr, "CACHEVERIFIER: failed to handle diffs: %v\n", err)
					}
				}
			}
		}
	}
	verifier.closeAndReset(diffErr)
}

func (verifier *cacheVerifier) closeAndReset(err error) {
	close(verifier.stopCh)
	verifier.diffMap = sync.Map{}
	verifier.stopCh = make(chan struct{})
	verifier.errCh <- err
}

func (verifier *cacheVerifier) handleDiff(cachedObj metav1.Object, nonCachedObj metav1.Object) error {
	var diffMapKey string
	var cachedResourceVersion string
	var nonCachedResourceVersion string
	if cachedObj != nil {
		diffMapKey = cache.MetaObjectToName(cachedObj).String()
		cachedResourceVersion = cachedObj.GetResourceVersion()
	}
	if nonCachedObj != nil {
		diffMapKey = cache.MetaObjectToName(nonCachedObj).String()
		nonCachedResourceVersion = nonCachedObj.GetResourceVersion()
	}

	updatedEntry := diffMapEntry{
		cachedResourceVersion:    cachedResourceVersion,
		nonCachedResourceVersion: nonCachedResourceVersion,
	}

	if entry, exists := verifier.diffMap.LoadOrStore(diffMapKey, updatedEntry); exists {
		if lastDiff, ok := entry.(diffMapEntry); ok {
			if lastDiff.cachedResourceVersion == cachedResourceVersion {
				// if cached item has not been updated since last time, we are probably stuck
				return &cacheVerificationError{
					itemName:                 diffMapKey,
					cachedResourceVersion:    cachedResourceVersion,
					nonCachedResourceVersion: nonCachedResourceVersion,
				}
			} else {
				// cached item has been updated since last time, but is not yet update, which is okay
				// we consider this as not being stuck (yet), but mark for check in next round
				verifier.diffMap.Store(diffMapKey, updatedEntry)
			}
		} else {
			return errors.New("unexpected cache type in diffMap")
		}
	}

	return nil
}
