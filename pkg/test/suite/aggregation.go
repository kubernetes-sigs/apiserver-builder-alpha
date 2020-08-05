package suite

import (
	"context"
	"fmt"
	"k8s.io/utils/pointer"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/onsi/gomega/gexec"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	apiregistrationv1client "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/typed/apiregistration/v1"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/testing_frameworks/integration/addr"
)

const (
	envLocalAPIServerBin = "TEST_ASSET_LOCAL_APISERVER"
)

type Environment struct {
	KubeAPIServerEnvironment envtest.Environment

	AggregatedAPIServerBuildArgs []string

	AggregatedAPIServerFlags        []string
	AggregatedAPIServerBinaryPath   string
	AggregatedAPIServerSession      *gexec.Session
	AggregatedAPIServerSecurePort   int
	AggregatedAPIServerInsecurePort int
	LoopbackClientConfig            *rest.Config
}

func NewDefaultTestingEnvironment() *Environment {
	securePort, _, _ := addr.Suggest()
	insecurePort, _, _ := addr.Suggest()
	return &Environment{
		AggregatedAPIServerSecurePort:   securePort,
		AggregatedAPIServerInsecurePort: insecurePort,
	}
}

func (e *Environment) StartLocalKubeAPIServer() error {
	var err error
	if e.LoopbackClientConfig, err = e.KubeAPIServerEnvironment.Start(); err != nil {
		return err
	}
	return nil
}

func (e *Environment) StartLocalAggregatedAPIServer(group, version string) error {
	if err := e.initAPIAggregationEnvironment(); err != nil {
		return err
	}
	if err := e.startAggregatedAPIServer(); err != nil {
		return err
	}
	if err := e.installAggregatedAPIServer(group, version); err != nil {
		return err
	}
	return nil
}

func (e *Environment) initAPIAggregationEnvironment() (err error) {
	e.AggregatedAPIServerFlags = append(e.AggregatedAPIServerFlags,
		"--etcd-servers="+e.KubeAPIServerEnvironment.ControlPlane.APIServer.EtcdURL.String(),
		"--cert-dir="+e.KubeAPIServerEnvironment.ControlPlane.APIServer.CertDir,
		"--delegated-auth=false",
		"--secure-port="+strconv.Itoa(e.AggregatedAPIServerSecurePort),
		"--insecure-port="+strconv.Itoa(e.AggregatedAPIServerInsecurePort),
		"--bind-address=127.0.0.1",
		"--insecure-bind-address=127.0.0.1")
	return nil
}

func (e *Environment) buildAggregatedAPIServer() (err error) {
	compiledPath := os.Getenv(envLocalAPIServerBin)
	if len(compiledPath) == 0 {
		// Compiling aggregated apiserver binary
		compiledPath, err = gexec.Build("../../../cmd/apiserver/main.go")
		if err != nil {
			return err
		}
	}

	e.AggregatedAPIServerBinaryPath = compiledPath
	e.AggregatedAPIServerFlags = []string{
		"--etcd-servers=" + e.KubeAPIServerEnvironment.ControlPlane.APIServer.EtcdURL.String(),
		"--cert-dir=" + e.KubeAPIServerEnvironment.ControlPlane.APIServer.CertDir,
		"--delegated-auth=false",
		"--secure-port=" + strconv.Itoa(e.AggregatedAPIServerSecurePort),
		"--insecure-port=" + strconv.Itoa(e.AggregatedAPIServerInsecurePort),
		"--bind-address=127.0.0.1",
		"--insecure-bind-address=127.0.0.1",
	}
	return err
}

func (e *Environment) installAggregatedAPIServer(group, version string) (err error) {
	serviceName := "aggregation-service"
	namespace := "default"

	corev1Client := corev1client.NewForConfigOrDie(e.LoopbackClientConfig)
	if _, err := corev1Client.Services(namespace).Create(context.TODO(), &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeExternalName,
			// NOTE: It turns out that it's a valid DNS1123 host name
			ExternalName: "127.0.0.1",
			Ports: []corev1.ServicePort{
				{
					Name:       "aggregation",
					Port:       443,
					TargetPort: intstr.FromInt(e.AggregatedAPIServerSecurePort),
				},
			},
		},
	}, metav1.CreateOptions{}); err != nil {
		return err
	}

	apiserviceClient := apiregistrationv1client.NewForConfigOrDie(e.LoopbackClientConfig)
	if _, err := apiserviceClient.APIServices().Create(context.TODO(), &apiregistrationv1.APIService{
		ObjectMeta: metav1.ObjectMeta{
			Name: version + "." + group,
		},
		Spec: apiregistrationv1.APIServiceSpec{
			Group:                 group,
			Version:               version,
			VersionPriority:       1000,
			GroupPriorityMinimum:  100,
			InsecureSkipTLSVerify: true,
			Service: &apiregistrationv1.ServiceReference{
				Namespace: namespace,
				Name:      serviceName,
				Port:      pointer.Int32Ptr(int32(e.AggregatedAPIServerSecurePort)),
			},
		},
	}, metav1.CreateOptions{}); err != nil {
		return err
	}
	if err := wait.PollImmediate(500*time.Millisecond, 5*time.Second, func() (done bool, err error) {
		resp, err := http.Get("http://" + e.KubeAPIServerEnvironment.ControlPlane.APIURL().Host + "/apis/" + group + "/" + version)
		return err == nil && resp.StatusCode == http.StatusOK, nil
	}); err != nil {
		return fmt.Errorf("failed installing aggregated api: %v", err)
	}
	return nil
}

func (e *Environment) startAggregatedAPIServer() (err error) {
	if err := e.buildAggregatedAPIServer(); err != nil {
		return err
	}

	cmd := exec.Command(e.AggregatedAPIServerBinaryPath, e.AggregatedAPIServerFlags...)
	session, err := gexec.Start(cmd, os.Stdout, os.Stderr)
	if err != nil {
		return err
	}

	if err := wait.PollImmediate(500*time.Millisecond, 5*time.Second, func() (done bool, err error) {
		var healthCheckErr error
		_, healthCheckErr = http.Get("http://" + net.JoinHostPort("127.0.0.1", strconv.Itoa(e.AggregatedAPIServerInsecurePort)) + "/healthz")
		return healthCheckErr == nil, nil
	}); err != nil {
		return fmt.Errorf("failed starting aggregated apiserver: %v", err)
	}
	e.AggregatedAPIServerSession = session
	return nil
}

func (e *Environment) StopLocalKubeAPIServer() (err error) {
	return e.KubeAPIServerEnvironment.Stop()
}

func (e *Environment) StopLocalAggregatedAPIServer() (err error) {
	if e.AggregatedAPIServerSession != nil {
		if err := wait.PollImmediate(100*time.Millisecond, 5*time.Second, func() (done bool, err error) {
			<-e.AggregatedAPIServerSession.Kill().Exited
			ln, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(e.AggregatedAPIServerSecurePort))
			defer func() {
				if ln != nil {
					ln.Close()
				}
			}()
			return err == nil, nil
		}); err != nil {
			return fmt.Errorf("port %v didn't released: %v", e.AggregatedAPIServerSecurePort, err)
		}
	}
	gexec.CleanupBuildArtifacts()
	return nil
}
