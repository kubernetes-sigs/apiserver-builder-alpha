package show

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	"k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	"sigs.k8s.io/apiserver-builder-alpha/pkg/utils"
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

func Validate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("should at least provide one resource name")
	}
	return nil
}

func RunShowResource(cmd *cobra.Command, args []string) {
	if err := Validate(args); err != nil {
		fmt.Fprintf(streams.ErrOut, "failed command validation: %v", err)
	}

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

	resourceName := args[0]
	apiResource := find(resourceName, apiResources)
	if apiResource == nil {
		fmt.Fprintf(streams.ErrOut, "No such resource found: %v", resourceName)
		return
	}
	apiServiceName := apiResource.Version + "." + apiResource.Group
	apiService, err := kubeAggregatorClient.ApiregistrationV1().
		APIServices().
		Get(context.TODO(), apiServiceName, metav1.GetOptions{})
	if err != nil {
		fmt.Fprintf(streams.ErrOut, "APIService %q not found are the resource %v", apiServiceName, resourceName)
		return
	}

	serviceNamespace := apiService.Spec.Service.Namespace
	serviceName := apiService.Spec.Service.Name

	service, err := kubeClient.CoreV1().Services(serviceNamespace).
		Get(context.TODO(), serviceName, metav1.GetOptions{})
	selector := labels.NewSelector()
	for k, v := range service.Spec.Selector {
		req, err := labels.NewRequirement(k, selection.Equals, []string{v})
		if err != nil {
			fmt.Fprintf(streams.ErrOut, "failed building label selector: %v", err)
			return
		}
		selector = selector.Add(*req)
	}
	podList, err := kubeClient.CoreV1().Pods(serviceNamespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		fmt.Fprintf(streams.ErrOut, "failed listing pods: %v", err)
		return
	}

	prefixWriter := utils.NewPrefixWriter(streams.Out)
	prefixWriter.Write(utils.LEVEL_0, "Resource Name: %v\n", resourceName)
	prefixWriter.Write(utils.LEVEL_0, "APIService:\n")
	prefixWriter.Write(utils.LEVEL_1, "Name: %v\n", apiService.Name)
	prefixWriter.Write(utils.LEVEL_1, "Service: %v/%v\n", apiService.Spec.Service.Namespace, apiService.Spec.Service.Name)
	available := ""
	for _, cond := range apiService.Status.Conditions {
		if cond.Type == apiregistrationv1.Available {
			available = string(cond.Status)
		}
	}
	prefixWriter.Write(utils.LEVEL_1, "Available: %v\n", available)
	prefixWriter.Write(utils.LEVEL_0, "Pods:\n")
	for _, pod := range podList.Items {
		printPod(kubeClient, &pod, prefixWriter)
	}

}

func printPod(kubeClient kubernetes.Interface, pod *corev1.Pod, prefixWriter utils.PrefixWriter) {
	prefixWriter.Write(utils.LEVEL_1, "Name: %v\n", pod.Name)
	prefixWriter.Write(utils.LEVEL_2, "State: %v\n", pod.Status.Phase)
	secretDir, certFile, isSelfSigned := isSelfSignedServerCertificate(pod)
	if isSelfSigned {
		prefixWriter.Write(utils.LEVEL_2, "Certificate: %v\n", color.RedString("<self-signed>"))
		return
	}
	secretName, found := findMountedSecretFromDir(pod, secretDir)
	if !found {
		prefixWriter.Write(utils.LEVEL_2, "Certificate: %v\n", color.RedString("<not-found>"))
		return
	}
	secret, err := kubeClient.CoreV1().Secrets(pod.Namespace).
		Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		fmt.Fprintf(streams.ErrOut, "failed getting secret %v: %v", secretName, err)
		return
	}
	certData, ok := secret.Data[certFile]
	if !ok {
		klog.Errorf(`no key "tls.crt" found in the secret %v`, secretName)
		return
	}
	notAfter, err := parseCertificateExpiry(certData)
	if err != nil {
		klog.Errorf(`failed parsing certificate in the secret %v`, secretName)
		return
	}
	prefixWriter.Write(utils.LEVEL_2, "Certificate:\n")
	if time.Now().After(*notAfter) {
		prefixWriter.Write(utils.LEVEL_3, "Not After: %v\n", color.RedString(notAfter.String()))
	} else {
		prefixWriter.Write(utils.LEVEL_3, "Not After: %v\n", color.GreenString(notAfter.String()))
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

const (
	certFileArg         = "--tls-cert-file"
	keyFileArg          = "--tls-private-key-file"
	certDirFlag         = "--cert-dir"
	certDirFlagFilename = "apiserver.crt"
)

func isSelfSignedServerCertificate(pod *corev1.Pod) (string, string, bool) {
	for _, c := range pod.Spec.Containers {
		hasCertFlag := false
		hasKeyFlag := false
		for _, arg := range c.Args {
			hasCertFlag = hasCertFlag || strings.HasPrefix(arg, certFileArg)
			hasKeyFlag = hasKeyFlag || strings.HasPrefix(arg, keyFileArg)
		}
		if hasCertFlag && hasKeyFlag {
			certFilePath := extractArgValue(c.Args, certFileArg)
			certDir := filepath.Dir(certFilePath)
			certFileName := filepath.Base(certFilePath)
			return certDir, certFileName, false
		}
		hasCertDirFlag := false
		for _, arg := range c.Args {
			hasCertDirFlag = hasCertDirFlag || strings.HasPrefix(arg, certDirFlag)
		}
		if hasCertDirFlag {
			certDir := extractArgValue(c.Args, certDirFlag)
			return certDir, certDirFlagFilename, false
		}
	}
	return "", "", true
}

func extractArgValue(args []string, argKey string) string {
	for i, arg := range args {
		if arg == argKey {
			return args[i+1]
		}
		if strings.HasPrefix(arg, argKey+"=") {
			return strings.TrimPrefix(arg, argKey+"=")
		}
		if strings.HasPrefix(arg, argKey) {
			return strings.TrimPrefix(arg, argKey)
		}
	}
	return ""
}

func findMountedSecretFromDir(pod *corev1.Pod, path string) (string, bool) {
	for _, c := range pod.Spec.Containers {
		for _, m := range c.VolumeMounts {
			dirname := filepath.Dir(m.MountPath)
			if strings.HasPrefix(path, dirname) {
				volumeName := m.Name
				for _, v := range pod.Spec.Volumes {
					if v.Name == volumeName {
						return v.Secret.SecretName, true
					}
				}
			}
		}
	}
	return "", false
}

func parseCertificateExpiry(pemCertData []byte) (*time.Time, error) {
	b, _ := pem.Decode(pemCertData)
	x509Cert, err := x509.ParseCertificate(b.Bytes)
	if err != nil {
		return nil, err
	}
	return &x509Cert.NotAfter, nil
}
