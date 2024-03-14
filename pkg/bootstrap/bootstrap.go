package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/glasskube/glasskube/internal/config"
	"github.com/glasskube/glasskube/internal/constants"
	"github.com/glasskube/glasskube/internal/releaseinfo"

	"github.com/fatih/color"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/schollz/progressbar/v3"
	appsv1 "k8s.io/api/apps/v1"
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
}

type BootstrapOptions struct {
	Type   BootstrapType
	Url    string
	Latest bool
}

func DefaultOptions() BootstrapOptions {
	return BootstrapOptions{Type: BootstrapTypeAio, Latest: config.IsDevBuild()}
}

const installMessage = `
## Installing GLASSKUBE ##
🧊 The missing Package Manager for Kubernetes 📦`

func NewBootstrapClient(config *rest.Config) *BootstrapClient {
	return &BootstrapClient{clientConfig: config}
}

func (c *BootstrapClient) Bootstrap(ctx context.Context, options BootstrapOptions) error {
	fmt.Println(installMessage)

	if options.Url == "" {
		version := config.Version
		if options.Latest {
			if releaseInfo, err := releaseinfo.FetchLatestRelease(); err != nil {
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
		return err
	}
	statusMessage("Successfully fetched Glasskube manifests", true)

	statusMessage("Applying Glasskube manifests", true)
	err = c.applyManifests(ctx, manifests)
	if err != nil {
		statusMessage(fmt.Sprintf("Couldn't apply manifests: %v", err), false)
	}
	statusMessage("Glasskube is successfully installed.", true)
	return nil
}

func (c *BootstrapClient) applyManifests(ctx context.Context, objs *[]unstructured.Unstructured) error {
	dynamicClient, err := dynamic.NewForConfig(c.clientConfig)
	if err != nil {
		return err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(c.clientConfig)
	if err != nil {
		return err
	}

	groupResources, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		return err
	}

	mapper := restmapper.NewDiscoveryRESTMapper(groupResources)

	bar := progressbar.Default(int64(len(*objs)), "Applying manifests")
	progressbar.OptionClearOnFinish()(bar)
	progressbar.OptionOnCompletion(nil)(bar)
	defer func(bar *progressbar.ProgressBar) {
		_ = bar.Exit()
	}(bar)

	for _, obj := range *objs {
		gvk := obj.GroupVersionKind()
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return err
		}

		bar.Describe(fmt.Sprintf("Applying %v (%v)", obj.GetName(), obj.GetKind()))
		err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
			_, err = dynamicClient.Resource(mapping.Resource).Namespace(obj.GetNamespace()).
				Apply(ctx, obj.GetName(), &obj, metav1.ApplyOptions{Force: true, FieldManager: "glasskube"})
			return err
		})

		if obj.GetKind() == constants.Deployment {
			bar.Describe(fmt.Sprintf("Checking Status of %v (%v)", obj.GetName(), obj.GetKind()))
			err = c.checkWorkloadReady(obj.GetNamespace(), obj.GetName(), obj.GetKind(), 5*time.Minute)
			if err != nil {
				return err
			}
		}
		err = bar.Add(1)
		if err != nil {
			return err
		}
	}

	return nil
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
			ready = availableReplicas == replicas
		case constants.DaemonSet:
			numberReady := status["numberReady"]
			desiredNumberScheduled := status["desiredNumberScheduled"]
			ready = numberReady == desiredNumberScheduled
		case constants.StatefulSet:
			readyReplicas := status["readyReplicas"]
			replicas := status["replicas"]
			ready = readyReplicas == replicas
		}

		return ready, nil
	}

	if ok, err := checkReady(); err != nil {
		return err
	} else if ok {
		return nil
	}

	timeoutCh := time.After(timeout)
	tick := time.NewTimer(5 * time.Second)
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
		fmt.Println(green("* " + input))
	} else {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Println(red("* " + input))
	}
}
