package open

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/pkg/open"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/cache"
)

type openState struct {
	host            string
	stopCh          chan struct{}
	forwarders      map[string]*open.OpenResult
	forwardersMutex sync.Mutex
}

var state *openState

func Init(host string, stopCh chan struct{}) {
	if state != nil {
		panic("open already initialized")
	}

	state = &openState{
		forwarders: make(map[string]*open.OpenResult),
		host:       host,
		stopCh:     stopCh,
	}
}

func HandleOpen(ctx context.Context, pkg ctrlpkg.Package) error {
	pkgClient := clicontext.PackageClientFromContext(ctx)
	fwName := cache.NewObjectName(pkg.GetNamespace(), pkg.GetName()).String()
	state.forwardersMutex.Lock()
	defer state.forwardersMutex.Unlock()
	if result, ok := state.forwarders[fwName]; ok {
		result.WaitReady()
		_ = cliutils.OpenInBrowser(result.Url)
		return nil
	}

	result, err := open.NewOpener().Open(ctx, pkg, "", state.host, 0)
	if err != nil {
		return err
	} else {
		state.forwarders[fwName] = result
		result.WaitReady()
		_ = cliutils.OpenInBrowser(result.Url)

		go func() {
			ctx = context.WithoutCancel(ctx)
		resultLoop:
			for {
				select {
				case <-state.stopCh:
					break resultLoop
				case fwErr := <-result.Completion:
					// note: this does not happen in "realtime" (e.g. when the forwarded-to-pod is deleted), but only
					// the next time a connection on that port is requested, e.g. when the user reloads the forwarded
					// page or clicks open again – only then we will end up here.
					if fwErr != nil {
						fmt.Fprintf(os.Stderr, "forwarder %v completed with error: %v\n", fwName, fwErr)
					} else {
						fmt.Fprintf(os.Stderr, "forwarder %v completed without error\n", fwName)
					}
					// try to re-open if the package is still installed
					if pkg.IsNamespaceScoped() {
						var p v1alpha1.Package
						if err := pkgClient.Packages(pkg.GetNamespace()).Get(ctx, pkg.GetName(), &p); err != nil {
							if !apierrors.IsNotFound(err) {
								fmt.Fprintf(os.Stderr, "can not reopen %v: %v\n", fwName, err)
							}
							break resultLoop
						} else {
							pkg = &p
						}
					} else {
						var cp v1alpha1.ClusterPackage
						if err := pkgClient.ClusterPackages().Get(ctx, pkg.GetName(), &cp); err != nil {
							if !apierrors.IsNotFound(err) {
								fmt.Fprintf(os.Stderr, "can not reopen %v: %v\n", fwName, err)
							}
							break resultLoop
						} else {
							pkg = &cp
						}
					}
					result, err = open.NewOpener().Open(ctx, pkg, "", state.host, 0)
					state.forwardersMutex.Lock()
					if err != nil {
						fmt.Fprintf(os.Stderr, "failed to reopen forwarder %v: %v\n", fwName, err)
						delete(state.forwarders, fwName)
						state.forwardersMutex.Unlock()
						break resultLoop
					} else {
						fmt.Fprintf(os.Stderr, "reopened forwarder %v\n", fwName)
						state.forwarders[fwName] = result
						state.forwardersMutex.Unlock()
					}
				}
			}
		}()
	}
	return nil
}

func CloseForwarders(pkg ctrlpkg.Package) {
	fwName := cache.NewObjectName(pkg.GetNamespace(), pkg.GetName()).String()
	state.forwardersMutex.Lock()
	if result, ok := state.forwarders[fwName]; ok {
		result.Stop()
		delete(state.forwarders, fwName)
	}
	state.forwardersMutex.Unlock()
}
