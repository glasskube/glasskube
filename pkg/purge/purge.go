package purge

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/httperror"
	"github.com/glasskube/glasskube/internal/releaseinfo"
	"github.com/glasskube/glasskube/internal/telemetry"
	"github.com/glasskube/glasskube/internal/telemetry/annotations"
	"github.com/glasskube/glasskube/pkg/bootstrap"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/schollz/progressbar/v3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

type PurgeClient struct {
	clientConfig *rest.Config
	mapper       meta.RESTMapper
	client       dynamic.Interface
}

type PurgeOptions struct {
	Type                    bootstrap.BootstrapType
	Url                     string
	Latest                  bool
	CreateDefaultRepository bool
	DisableTelemetry        bool
	Force                   bool
}

func DefaultPurgeOptions() PurgeOptions {
	return PurgeOptions{Type: bootstrap.BootstrapTypeAio, Latest: config.IsDevBuild(), CreateDefaultRepository: true}
}

func NewPurgeClient(config *rest.Config) *PurgeClient {
	return &PurgeClient{clientConfig: config}
}

func (c *PurgeClient) initRestMapper() error {
	if discoveryClient, err := discovery.NewDiscoveryClientForConfig(c.clientConfig); err != nil {
		return err
	} else if groupResources, err := restmapper.GetAPIGroupResources(discoveryClient); err != nil {
		return err
	} else {
		c.mapper = restmapper.NewDiscoveryRESTMapper(groupResources)
		return nil
	}
}

func (c *PurgeClient) Purge(ctx context.Context, options PurgeOptions) error {
	if err := c.initRestMapper(); err != nil {
		return err
	}

	if client, err := dynamic.NewForConfig(c.clientConfig); err != nil {
		return err
	} else {
		c.client = client
	}

	fmt.Println("Starting purge process")
	start := time.Now()

	if options.Url == "" {
		version := config.Version
		if options.Latest {
			if releaseInfo, err := releaseinfo.FetchLatestRelease(); err != nil {
				if httperror.Is(err, http.StatusServiceUnavailable) || httperror.IsTimeoutError(err) {
					telemetry.BootstrapFailure(time.Since(start))
					return fmt.Errorf("network connectivity error, check your network: %w", err)
				}
				telemetry.BootstrapFailure(time.Since(start))
				return fmt.Errorf("could not determine latest version: %w", err)
			} else {
				version = releaseInfo.Version
			}
		}
		options.Url = fmt.Sprintf("https://github.com/glasskube/glasskube/releases/download/v%v/manifest-%v.yaml",
			version, options.Type)
	}

	statusMessage("Fetching Glasskube manifest from "+options.Url, true)
	manifests, err := clientutils.FetchResources(options.Url)
	if err != nil {
		statusMessage("Couldn't fetch Glasskube manifests", false)
		return fmt.Errorf("couldn't fetch Glasskube manifests: %w", err)
	}

	statusMessage("Validating existing installation", true)

	if err = c.preprocessManifests(ctx, manifests, &options); err != nil {
		telemetry.BootstrapFailure(time.Since(start))
		statusMessage(fmt.Sprintf("Couldn't prepare manifests: %v", err), false)
		if !options.Force {
			return err
		} else {
			statusMessage("Attempting to force bootstrap anyways (Force option is enabled)", true)
		}
	}

	if options.CreateDefaultRepository {
		manifests = append(manifests, bootstrap.DefaultRepository())
	}

	statusMessage("Deleting Glasskube resources", true)
	if err := c.purgeManifests(ctx, manifests); err != nil {
		return fmt.Errorf("an error occurred during purge: %w", err)
	}

	elapsed := time.Since(start)
	c.handleTelemetry(options.DisableTelemetry, elapsed)

	statusMessage(fmt.Sprintf("Glasskube successfully purged! (took %v)", elapsed.Round(time.Second)), true)
	return nil
}

func (c *PurgeClient) purgeManifests(ctx context.Context, objs []unstructured.Unstructured) error {
	bar := progressbar.Default(int64(len(objs)), "Deleting manifests")
	progressbar.OptionClearOnFinish()(bar)
	progressbar.OptionOnCompletion(nil)(bar)
	defer func(bar *progressbar.ProgressBar) { _ = bar.Exit() }(bar)

	for _, obj := range objs {
		gvk := obj.GroupVersionKind()
		mapping, err := c.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return fmt.Errorf("could not get restmapping for %v %v: %w", obj.GetKind(), obj.GetName(), err)
		}

		bar.Describe(fmt.Sprintf("Deleting %v (%v)", obj.GetName(), obj.GetKind()))
		err = c.client.Resource(mapping.Resource).Namespace(obj.GetNamespace()).
			Delete(ctx, obj.GetName(), metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			return fmt.Errorf("could not delete %v %v: %w", obj.GetKind(), obj.GetName(), err)
		} else if errors.IsNotFound(err) {
			fmt.Printf("Resource %v %v not found, skipping.\n", obj.GetKind(), obj.GetName())
		}

		_ = bar.Add(1)
	}

	return nil
}

func (c *PurgeClient) preprocessManifests(
	ctx context.Context,
	objs []unstructured.Unstructured,
	options *PurgeOptions,
) error {
	var compositeErr error
	for _, obj := range objs {
		gvk := obj.GroupVersionKind()
		mapping, err := c.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return err
		}
		existing, err := c.client.Resource(mapping.Resource).Namespace(obj.GetNamespace()).
			Get(ctx, obj.GetName(), metav1.GetOptions{})
		if err != nil && !apierrors.IsNotFound(err) {
			return err
		} else if err == nil {
			if _, ok := existing.GetAnnotations()["kubectl.kubernetes.io/last-applied-configuration"]; ok {
				multierr.AppendInto(&compositeErr,
					fmt.Errorf("%v %v has been modified with kubectl", obj.GetKind(), obj.GetName()))
			}
		}

		if obj.GetKind() == "Namespace" && obj.GetName() == "glasskube-system" {
			nsAnnotations := obj.GetAnnotations()
			if nsAnnotations == nil {
				nsAnnotations = make(map[string]string, 1)
			}
			if existing != nil {
				existingAnnotations := existing.GetAnnotations()
				nsAnnotations[annotations.TelemetryEnabledAnnotation] = existingAnnotations[annotations.TelemetryEnabledAnnotation]
				nsAnnotations[annotations.TelemetryIdAnnotation] = existingAnnotations[annotations.TelemetryIdAnnotation]
			}
			annotations.UpdateTelemetryAnnotations(nsAnnotations, options.DisableTelemetry)
			obj.SetAnnotations(nsAnnotations)
			options.DisableTelemetry = !annotations.IsTelemetryEnabled(nsAnnotations)
		}
	}

	if compositeErr != nil {
		compositeErr = fmt.Errorf("unsupported installation: %w", compositeErr)
	}
	return compositeErr
}

func (c *PurgeClient) handleTelemetry(disabled bool, elapsed time.Duration) {
	if !disabled {
		statusMessage("Telemetry is enabled for this cluster – "+
			"Run \"glasskube telemetry status\" for more info.", true)
		telemetry.BootstrapSuccess(elapsed)
	}
}

func statusMessage(input string, success bool) {
	if success {
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Println(green("* " + input))
	} else {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Println(red("* " + input))
	}
}
