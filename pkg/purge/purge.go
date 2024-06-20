package purge

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/util"
	"github.com/glasskube/glasskube/pkg/statuswriter"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

type Purger struct {
	clientConfig *rest.Config
	mapper       meta.RESTMapper
	client       dynamic.Interface
	status       statuswriter.StatusWriter
}

func NewPurger(config *rest.Config) *Purger {
	return &Purger{clientConfig: config, status: statuswriter.Noop()}
}

func (c *Purger) WithStatusWriter(sw statuswriter.StatusWriter) *Purger {
	c.status = sw
	return c
}

func (c *Purger) initRestMapper() error {
	if discoveryClient, err := discovery.NewDiscoveryClientForConfig(c.clientConfig); err != nil {
		return err
	} else if groupResources, err := restmapper.GetAPIGroupResources(discoveryClient); err != nil {
		return err
	} else {
		c.mapper = restmapper.NewDiscoveryRESTMapper(groupResources)
		return nil
	}
}

func (c *Purger) Purge(ctx context.Context) error {
	c.status.Start()
	defer c.status.Stop()
	if err := c.initRestMapper(); err != nil {
		return err
	}

	if client, err := dynamic.NewForConfig(c.clientConfig); err != nil {
		return err
	} else {
		c.client = client
	}

	c.status.SetStatus("Starting purge process")

	operatorVersion, err := clientutils.GetPackageOperatorVersion(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to check package operator version: %v\n", err)
	}

	manifestUrl := fmt.Sprintf("https://github.com/glasskube/glasskube/releases/download/%v/manifest-%s.yaml",
		operatorVersion, "slim")

	c.status.SetStatus("Fetching Glasskube manifest from " + manifestUrl)
	manifests, err := clientutils.FetchResources(manifestUrl)
	if err != nil {
		return fmt.Errorf("Couldn't fetch Glasskube manifests: %w", err)
	}

	c.status.SetStatus("Deleting Glasskube resources")
	if err := c.purgeManifests(ctx, manifests); err != nil {
		return fmt.Errorf("an error occurred during purge: %w", err)
	}

	return nil
}

func (c *Purger) purgeManifests(ctx context.Context, objs []unstructured.Unstructured) error {
	slices.Reverse(objs)
	for _, obj := range objs {
		gvk := obj.GroupVersionKind()
		mapping, err := c.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return fmt.Errorf("could not get restmapping for %v %v: %w", obj.GetKind(), obj.GetName(), err)
		}

		c.status.SetStatus(fmt.Sprintf("Deleting %v (%v)", obj.GetName(), obj.GetKind()))
		err = c.client.Resource(mapping.Resource).Namespace(obj.GetNamespace()).
			Delete(ctx, obj.GetName(), metav1.DeleteOptions{PropagationPolicy: util.Pointer(metav1.DeletePropagationForeground)})
		if err != nil && !errors.IsNotFound(err) {
			return fmt.Errorf("could not delete %v %v: %w", obj.GetKind(), obj.GetName(), err)
		} else if errors.IsNotFound(err) {
			fmt.Printf("Resource %v %v not found, skipping.\n", obj.GetKind(), obj.GetName())
		}
	}

	return nil
}
