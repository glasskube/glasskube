package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/constants"
	"github.com/glasskube/glasskube/internal/httperror"
	"github.com/glasskube/glasskube/internal/releaseinfo"
	"github.com/glasskube/glasskube/internal/telemetry"
	"github.com/glasskube/glasskube/internal/telemetry/annotations"
	"github.com/glasskube/glasskube/internal/util"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/multierr"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/util/retry"
)

type BootstrapClient struct {
	clientConfig *rest.Config
	mapper       meta.RESTMapper
	client       dynamic.Interface
}

type BootstrapOptions struct {
	Type                    BootstrapType
	Url                     string
	Latest                  bool
	DisableTelemetry        bool
	Force                   bool
	CreateDefaultRepository bool
	DryRun                  bool
}

func DefaultOptions() BootstrapOptions {
	return BootstrapOptions{
		Type:                    BootstrapTypeAio,
		Latest:                  config.IsDevBuild(),
		CreateDefaultRepository: true,
		DryRun:                  false,
	}
}

const installMessage = `
## Installing GLASSKUBE ##
🧊 The next generation Package Manager for Kubernetes 📦`

func NewBootstrapClient(config *rest.Config) *BootstrapClient {
	return &BootstrapClient{clientConfig: config}
}

func (c *BootstrapClient) initRestMapper() error {
	if discoveryClient, err := discovery.NewDiscoveryClientForConfig(c.clientConfig); err != nil {
		return err
	} else if groupResources, err := restmapper.GetAPIGroupResources(discoveryClient); err != nil {
		return err
	} else {
		c.mapper = restmapper.NewDiscoveryRESTMapper(groupResources)
		return nil
	}
}

func (c *BootstrapClient) Bootstrap(
	ctx context.Context,
	options BootstrapOptions,
) ([]unstructured.Unstructured, error) {
	telemetry.BootstrapAttempt()

	if err := c.initRestMapper(); err != nil {
		return nil, err
	}

	if client, err := dynamic.NewForConfig(c.clientConfig); err != nil {
		return nil, err
	} else {
		c.client = client
	}

	start := time.Now()

	if options.Url == "" {
		version := config.Version
		if options.Latest {
			if releaseInfo, err := releaseinfo.FetchLatestRelease(); err != nil {
				if httperror.Is(err, http.StatusServiceUnavailable) || httperror.IsTimeoutError(err) {
					telemetry.BootstrapFailure(time.Since(start))
					return nil, fmt.Errorf("network connectivity error, check your network: %w", err)
				}
				telemetry.BootstrapFailure(time.Since(start))
				return nil, fmt.Errorf("could not determine latest version: %w", err)
			} else {
				version = releaseInfo.Version
			}
		}
		options.Url = fmt.Sprintf("https://github.com/glasskube/glasskube/releases/download/v%v/manifest-%v.yaml",
			version, options.Type)
	}

	fmt.Fprintln(os.Stderr, installMessage)

	statusMessage("Fetching Glasskube manifest from "+options.Url, true)
	manifests, err := clientutils.FetchResources(options.Url)
	if err != nil {
		statusMessage("Couldn't fetch Glasskube manifests", false)
		telemetry.BootstrapFailure(time.Since(start))
		return nil, err
	}

	statusMessage("Validating existing installation", true)

	if err = c.preprocessManifests(ctx, manifests, &options); err != nil {
		telemetry.BootstrapFailure(time.Since(start))
		statusMessage(fmt.Sprintf("Couldn't prepare manifests: %v", err), false)
		if !options.Force {
			return nil, err
		} else {
			statusMessage("Attempting to force bootstrap anyways (Force option is enabled)", true)
		}
	}

	if options.CreateDefaultRepository {
		manifests = append(manifests, defaultRepository())
	}

	statusMessage("Applying Glasskube manifests", true)

	if err = c.applyManifests(ctx, manifests, options); err != nil {
		telemetry.BootstrapFailure(time.Since(start))
		statusMessage(fmt.Sprintf("Couldn't apply manifests: %v", err), false)
		return nil, err
	}

	elapsed := time.Since(start)
	c.handleTelemetry(options.DisableTelemetry, elapsed)

	statusMessage(fmt.Sprintf("Glasskube successfully installed! (took %v)", elapsed.Round(time.Second)), true)
	return manifests, nil
}

func (c *BootstrapClient) preprocessManifests(
	ctx context.Context,
	objs []unstructured.Unstructured,
	options *BootstrapOptions,
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

func (c *BootstrapClient) applyManifests(
	ctx context.Context,
	objs []unstructured.Unstructured,
	options BootstrapOptions,
) error {
	bar := progressbar.Default(int64(len(objs)), "Applying manifests")
	progressbar.OptionClearOnFinish()(bar)
	progressbar.OptionOnCompletion(nil)(bar)
	defer func(bar *progressbar.ProgressBar) { _ = bar.Exit() }(bar)

	var checkWorkloads []*unstructured.Unstructured
	for i, obj := range objs {
		gvk := obj.GroupVersionKind()
		mapping, err := c.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			if options.DryRun {
				continue
			} else {
				return fmt.Errorf("could not get restmapping for %v %v: %w", obj.GetKind(), obj.GetName(), err)
			}
		}

		bar.Describe(fmt.Sprintf("Applying %v (%v)", obj.GetName(), obj.GetKind()))
		if obj.GetKind() == "Job" {
			err := c.client.Resource(mapping.Resource).Namespace(obj.GetNamespace()).
				Delete(ctx, obj.GetName(), getDeleteOptions(options))
			if err != nil && !apierrors.IsNotFound(err) {
				return err
			}
		}

		if err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
			_, err = c.client.Resource(mapping.Resource).Namespace(obj.GetNamespace()).
				Apply(ctx, obj.GetName(), &obj, getApplyOptions(options))
			return err
		}); err != nil && (!options.DryRun || !apierrors.IsNotFound(err)) {
			return err
		}

		if obj.GetKind() == constants.Deployment && !options.DryRun {
			checkWorkloads = append(checkWorkloads, &objs[i])
			bar.ChangeMax(bar.GetMax() + 1)
		} else if obj.GetKind() == "CustomResourceDefinition" {
			// The RESTMapping must be re-created after applying a CRD, so we can create resources of that kind immediately.
			if err := c.initRestMapper(); err != nil {
				return err
			}
		}

		_ = bar.Add(1)
	}

	for _, obj := range checkWorkloads {
		bar.Describe(fmt.Sprintf("Checking Status of %v (%v)", obj.GetName(), obj.GetKind()))
		if err := c.checkWorkloadReady(obj.GetNamespace(), obj.GetName(), obj.GetKind(), 5*time.Minute); err != nil {
			return err
		}
		_ = bar.Add(1)
	}

	return nil
}

func getApplyOptions(options BootstrapOptions) metav1.ApplyOptions {
	applyOptions := metav1.ApplyOptions{Force: true, FieldManager: "glasskube"}
	if options.DryRun {
		applyOptions.DryRun = []string{metav1.DryRunAll}
	}
	return applyOptions
}

func getDeleteOptions(options BootstrapOptions) metav1.DeleteOptions {
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: util.Pointer(metav1.DeletePropagationBackground),
	}
	if options.DryRun {
		deleteOptions.DryRun = []string{metav1.DryRunAll}
	}
	return deleteOptions
}

func (c *BootstrapClient) handleTelemetry(disabled bool, elapsed time.Duration) {
	if !disabled {
		statusMessage("Telemetry is enabled for this cluster – "+
			"Run \"glasskube telemetry status\" for more info.", true)
		telemetry.BootstrapSuccess(elapsed)
	}
}

func defaultRepository() unstructured.Unstructured {
	repo := v1alpha1.PackageRepository{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1alpha1.GroupVersion.String(),
			Kind:       "PackageRepository",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "glasskube",
		},
		Spec: v1alpha1.PackageRepositorySpec{
			Url: constants.DefaultRepoUrl,
		},
	}
	repo.SetDefaultRepository()
	var repoUnstructured unstructured.Unstructured
	if err := json.Unmarshal(util.Must(json.Marshal(repo)), &repoUnstructured); err != nil {
		panic(err)
	}
	return repoUnstructured
}

func (c *BootstrapClient) checkWorkloadReady(
	namespace string,
	workloadName string,
	workloadType string,
	timeout time.Duration,
) error {
	dynamicClient, err := dynamic.NewForConfig(c.clientConfig)
	if err != nil {
		return err
	}

	var workloadRes schema.GroupVersionResource

	switch workloadType {
	case constants.Deployment:
		workloadRes = appsv1.SchemeGroupVersion.WithResource("deployments")
	case constants.DaemonSet:
		workloadRes = appsv1.SchemeGroupVersion.WithResource("daemonsets")
	case constants.StatefulSet:
		workloadRes = appsv1.SchemeGroupVersion.WithResource("statefulsets")
	default:
		return fmt.Errorf("unsupported workload type: %s", workloadType)
	}

	checkReady := func() (bool, error) {
		workload, err := dynamicClient.
			Resource(workloadRes).
			Namespace(namespace).
			Get(context.Background(), workloadName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		status := workload.Object["status"].(map[string]interface{})
		var ready bool

		switch workloadType {
		case constants.Deployment:
			availableReplicas := status["availableReplicas"]
			replicas := status["replicas"]
			ready = availableReplicas != nil && availableReplicas == replicas
		case constants.DaemonSet:
			numberReady := status["numberReady"]
			desiredNumberScheduled := status["desiredNumberScheduled"]
			ready = numberReady != nil && numberReady == desiredNumberScheduled
		case constants.StatefulSet:
			readyReplicas := status["readyReplicas"]
			replicas := status["replicas"]
			ready = readyReplicas != nil && readyReplicas == replicas
		}

		return ready, nil
	}

	if ok, err := checkReady(); err != nil {
		return err
	} else if ok {
		return nil
	}

	timeoutCh := time.After(timeout)
	tick := time.NewTicker(2 * time.Second)
	tickC := tick.C
	defer tick.Stop()

	for {
		select {
		case <-timeoutCh:
			return fmt.Errorf("%s is not ready within the specified timeout", workloadType)
		case <-tickC:
			if ok, err := checkReady(); err != nil {
				return err
			} else if ok {
				return nil
			}
		}
	}
}

func statusMessage(input string, success bool) {
	if success {
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Fprintln(os.Stderr, green("* "+input))
	} else {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Fprintln(os.Stderr, red("* "+input))
	}
}
