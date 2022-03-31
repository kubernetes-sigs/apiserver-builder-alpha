package show

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/printers"
	"os"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
)

var showResourceCmd = &cobra.Command{
	Use:     "resource",
	Short:   "Show the status of an aggregated resource.",
	Long:    "Show the status of an aggregated resource.",
	Example: `apiserver-boot show resource foo`,
	Run:     RunShowResource,
}

var clientFactory genericclioptions.RESTClientGetter
var apiVersion string
var streams = genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

func AddShowResource(cmd *cobra.Command) {
	cmd.AddCommand(showResourceCmd)

	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	kubeConfigFlags.AddFlags(showResourceCmd.Flags())

	clientFactory = kubeConfigFlags
	showResourceCmd.Flags().StringVar(&apiVersion, "api-version", "",
		"The apiVersion of the showing resource.")
}

func RunShowResource(cmd *cobra.Command, args []string) {
	kubeClientConfig, err := clientFactory.ToRESTConfig()
	if err != nil {
		klog.Error("")
		return
	}
	kubeClient, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		klog.Error("")
		return
	}
	kubeAggregatorClient, err := clientset.NewForConfig(kubeClientConfig)
	if err != nil {
		klog.Error("")
		return
	}

	var apiResources []metav1.APIResource
	if len(apiVersion) > 0 {
		apiResourceList, err := kubeClient.Discovery().ServerResourcesForGroupVersion(apiVersion)
		if err != nil {
			klog.Fatal("")
		}
		apiResources = apiResourceList.APIResources
	} else {
		apiResourceLists, err := kubeClient.Discovery().ServerPreferredResources()
		if err != nil {
			klog.Fatal("")
		}
		for _, apiResourceList := range apiResourceLists {
			gv, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
			if err != nil {
				klog.Error(err)
			}
			for _, apiResource := range apiResourceList.APIResources {
				apiResource := apiResource
				apiResource.Group = gv.Group
				apiResource.Version = gv.Version
				apiResources = append(apiResources, apiResource)
			}
		}
	}

	for _, resourceName := range args {
		apiResource := find(resourceName, apiResources)
		if apiResource == nil {
			klog.Infof("No such resource %v", resourceName)
			continue
		}
		apiServiceName := apiResource.Version + "." + apiResource.Group
		apiService, err := kubeAggregatorClient.ApiregistrationV1().
			APIServices().
			Get(context.TODO(), apiServiceName, metav1.GetOptions{})
		if err != nil {
			klog.Error(err)
			continue
		}
		printer := printers.NewTablePrinter(printers.PrintOptions{})
		if err := printer.PrintObj(apiService, streams.Out); err != nil {
			klog.Error(err)
			continue
		}
	}
}

func find(resourceName string, apiResources []metav1.APIResource) *metav1.APIResource {
	for _, apiResource := range apiResources {
		if apiResource.Name == resourceName {
			return &apiResource
		}
		if apiResource.SingularName == resourceName {
			return &apiResource
		}
		if apiResource.Kind == resourceName {
			return &apiResource
		}
	}
	return nil
}
