package show

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/apiserver-builder-alpha/pkg/utils"
)

var showApiserverCmd = &cobra.Command{
	Use:     "apiserver",
	Short:   "Show the running detail of the aggregated apiserver.",
	Long:    "Show the running detail of the aggregated apiserver.",
	Example: `apiserver-boot show apiserver -n <pod namespace> <pod name>`,
	Run:     RunShowApiserver,
}

var port int32

func AddApiserver(cmd *cobra.Command) {
	cmd.AddCommand(showApiserverCmd)

	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	kubeConfigFlags.AddFlags(showApiserverCmd.Flags())

	clientFactory = kubeConfigFlags
	showApiserverCmd.Flags().Int32VarP(&port, "port", "p", 443,
		"The serving port of the target aggregated apiserver")
}

func ValidateApiserver(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("should at least provide one pod name")
	}
	return nil
}

func RunShowApiserver(cmd *cobra.Command, args []string) {
	if err := ValidateApiserver(args); err != nil {
		fmt.Fprintf(streams.ErrOut, "failed command validation: %v", err)
		return
	}

	kubeClientConfig, err := clientFactory.ToRESTConfig()
	if err != nil {
		fmt.Fprintf(streams.ErrOut, "Failed building kube client config: %v", err)
		return
	}
	kubeClient, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		fmt.Fprintf(streams.ErrOut, "Failed building kube client: %v", err)
		return
	}

	podName := args[0]
	podNamespace := corev1.NamespaceDefault
	if clientFactory.Namespace != nil {
		podNamespace = *clientFactory.Namespace
	}

	pod, err := kubeClient.CoreV1().Pods(podNamespace).
		Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		fmt.Fprintf(streams.ErrOut, "Failed getting pod %v/%v: %v", podNamespace, podName, err)
		return
	}

	proxy := utils.NewLocalProxy(kubeClientConfig, podNamespace, podName, port)
	stopFunc, err := proxy.Listen()
	if err != nil {
		fmt.Fprintf(streams.ErrOut, "Failed running local proxy to port %v: %v", port, err)
		return
	}
	defer stopFunc()

	podClientConfig := rest.CopyConfig(kubeClientConfig)
	podClientConfig.Insecure = true
	podClientConfig.CAData = nil
	podClientConfig.Host = net.JoinHostPort("127.0.0.1", strconv.Itoa(int(port)))

	podClient, err := kubernetes.NewForConfig(podClientConfig)
	if err != nil {
		fmt.Fprintf(streams.ErrOut, "Failed building aggregated apiserver client: %v", err)
		return
	}

	resultData, err := podClient.RESTClient().Get().
		AbsPath("/healthz").DoRaw(context.TODO())
	if err != nil {
		fmt.Fprintf(streams.ErrOut, "Failed calling aggregated apiserver: %v", err)
		return
	}
	isHealthy := string(resultData) == "ok"
	// group -> version -> resource
	apiVersions := make(map[string]map[string][]string)
	_, apiResourceLists, err := podClient.Discovery().ServerGroupsAndResources()
	if err != nil {
		fmt.Fprintf(streams.ErrOut, "Failed calling api-discovery upon the aggregated apiserver: %v", err)
		return
	}
	for _, apiResourceList := range apiResourceLists {
		if apiVersions[apiResourceList.GroupVersion] == nil {
			apiVersions[apiResourceList.GroupVersion] = make(map[string][]string)
		}
		for _, apiResource := range apiResourceList.APIResources {
			apiVersions[apiResourceList.GroupVersion][apiResource.Name] = apiResource.Verbs
		}
	}

	metricsData, err := podClient.RESTClient().Get().
		AbsPath("/metrics").DoRaw(context.TODO())

	prefixWriter := utils.NewPrefixWriter(streams.Out)
	prefixWriter.Write(utils.LEVEL_0, "Pod Namespace: %v\n", podNamespace)
	prefixWriter.Write(utils.LEVEL_0, "Pod Name: %v\n", podName)
	prefixWriter.Write(utils.LEVEL_0, "Healthiness: %v\n", strconv.FormatBool(isHealthy))
	prefixWriter.Write(utils.LEVEL_0, "Resources:\n")
	for apiVersion, resources := range apiVersions {
		prefixWriter.Write(utils.LEVEL_1, "APIVersion: %v\n", apiVersion)
		for resource, verbs := range resources {
			prefixWriter.Write(utils.LEVEL_2, "Resource: %v %v\n", resource, verbs)
		}
	}
	prefixWriter.Write(utils.LEVEL_0, "Metrics:\n")
	memReq := "<none>"
	memLimit := "<none>"
	rss, _ := parseInUsedMemory(string(metricsData))
	for _, c := range pod.Spec.Containers {
		if assignedMemReq, ok := c.Resources.Requests[corev1.ResourceMemory]; ok {
			memReq = assignedMemReq.String()
		}
		if assignedMemLimits, ok := c.Resources.Requests[corev1.ResourceMemory]; ok {
			memLimit = assignedMemLimits.String()
		}
	}
	prefixWriter.Write(utils.LEVEL_1, "Memory: [%v|%v|%v]\n", memReq, rss, memLimit)
}

const (
	inUsedMemoryMetricsName = "process_resident_memory_bytes"
)

func parseInUsedMemory(metricsData string) (string, bool) {
	lines := strings.Split(metricsData, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, inUsedMemoryMetricsName) {
			bytesStr := strings.TrimSpace(strings.TrimPrefix(line, inUsedMemoryMetricsName))
			q, err := resource.ParseQuantity(bytesStr)
			if err != nil {
				return "", false
			}
			rssInMega := q.ScaledValue(resource.Mega)
			return strconv.Itoa(int(rssInMega)) + "Mb", true
		}
	}
	return "", false
}
