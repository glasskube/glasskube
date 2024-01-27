package bootstrap

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"io"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
	"net/http"
	"strings"
	"time"
)

type BootstrapClient struct {
	url          string
	version      string
	clientConfig *rest.Config
}

const INSTALLMESSAGE = `
## Installing GLASSKUBE ##
🧊 The missing Package Manager for Kubernetes 📦
`

func NewBootstrapClient(version string, kubeconfig string, url string) (*BootstrapClient, error) {
	config, err := initKubeConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	if url == "" {
		url = fmt.Sprintf("https://github.com/glasskube/glasskube/releases/download/v%s/manifest.yaml", version)
	}

	return &BootstrapClient{
		url:          url,
		version:      version,
		clientConfig: config,
	}, nil
}

func (c *BootstrapClient) Bootstrap() error {
	fmt.Println(INSTALLMESSAGE)

	statusMessage("Fetching Glasskube manifest from "+c.url, true)
	manifests, err := readManifest(c.url)
	if err != nil {
		statusMessage("Couldn't fetch Glasskube manifests", false)
		return err
	}
	statusMessage("Successfully fetched Glasskube manifests", true)

	statusMessage("Applying Glasskube manifests", true)
	err = c.applyManifests(manifests)
	if err != nil {
		statusMessage("Couldn't apply manifests", false)
	}
	statusMessage("Glasskube is successfully installed.", true)
	return nil
}

func (c *BootstrapClient) applyManifests(objs []unstructured.Unstructured) error {
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

	var bar *progressbar.ProgressBar

	bar = progressbar.Default(int64(len(objs)))
	bar.Describe("Applying manifests")

	for _, obj := range objs {
		gvk := obj.GroupVersionKind()
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return err
		}

		bar.Describe("Applying " + obj.GetName() + " (" + obj.GetKind() + ")")
		err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
			_, err = dynamicClient.Resource(mapping.Resource).Namespace(obj.GetNamespace()).Create(context.Background(), &obj, metav1.CreateOptions{})
			if err != nil {
				if errors.IsAlreadyExists(err) {
					return nil
				}
				return err
			}
			return nil
		})

		if obj.GetKind() == "Deployment" {
			bar.Describe("Checking Status of " + obj.GetName() + " (" + obj.GetKind() + ")")
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
	bar.Finish()

	return nil
}

func (c *BootstrapClient) checkWorkloadReady(namespace string, workloadName string, workloadType string, timeout time.Duration) error {
	dynamicClient, err := dynamic.NewForConfig(c.clientConfig)
	if err != nil {
		return err
	}

	var workloadRes schema.GroupVersionResource

	switch workloadType {
	case "Deployment":
		workloadRes = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	case "DaemonSet":
		workloadRes = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "daemonsets"}
	case "StatefulSet":
		workloadRes = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"}
	default:
		return fmt.Errorf("unsupported workload type: %s", workloadType)
	}

	timeoutCh := time.After(timeout)
	tick := time.Tick(5 * time.Second)

	for {
		select {
		case <-timeoutCh:
			return fmt.Errorf("%s is not ready within the specified timeout", workloadType)
		case <-tick:
			workload, err := dynamicClient.Resource(workloadRes).Namespace(namespace).Get(context.Background(), workloadName, metav1.GetOptions{})
			if err != nil {
				return err
			}

			status := workload.Object["status"].(map[string]interface{})
			var ready bool

			switch workloadType {
			case "Deployment":
				availableReplicas := status["availableReplicas"]
				replicas := status["replicas"]
				ready = availableReplicas == replicas
			case "DaemonSet":
				numberReady := status["numberReady"]
				desiredNumberScheduled := status["desiredNumberScheduled"]
				ready = numberReady == desiredNumberScheduled
			case "StatefulSet":
				readyReplicas := status["readyReplicas"]
				replicas := status["replicas"]
				ready = readyReplicas == replicas
			}

			if ready {
				return nil
			}
		}
	}
}

func initKubeConfig(kubeconfig string) (*rest.Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if kubeconfig != "" {
		loadingRules.ExplicitPath = kubeconfig
	}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	return clientConfig.ClientConfig()
}

func readManifest(url string) ([]unstructured.Unstructured, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fileContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	manifests := strings.Split(string(fileContent), "---")
	var objs []unstructured.Unstructured

	for _, manifest := range manifests {
		var obj unstructured.Unstructured
		manifestReader := strings.NewReader(manifest)
		decoder := yaml.NewYAMLOrJSONDecoder(manifestReader, 4096)
		err := decoder.Decode(&obj)
		if err != nil {
			return nil, err
		}
		objs = append(objs, obj)
	}

	return objs, nil
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
