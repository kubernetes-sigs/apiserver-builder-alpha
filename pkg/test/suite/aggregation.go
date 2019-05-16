package suite

import (
	"fmt"
	"github.com/onsi/gomega/gexec"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	apiregistrationv1client "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/typed/apiregistration/v1"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"strconv"
	"time"
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
	return &Environment{
		AggregatedAPIServerSecurePort:   443,
		AggregatedAPIServerInsecurePort: 8080,
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
	// Compiling aggregated apiserver binary
	binName := "aggregated-apiserver"
	binPath := filepath.Join(e.KubeAPIServerEnvironment.ControlPlane.APIServer.CertDir, binName)
	cmd := exec.Command("go",
		append(
			append(
				[]string{"build", "-o", binPath}, e.AggregatedAPIServerBuildArgs...),
			"../../../cmd/apiserver/main.go")...)
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return err
	}

	e.AggregatedAPIServerBinaryPath = binPath
	e.AggregatedAPIServerSecurePort = 443
	e.AggregatedAPIServerInsecurePort = 8080
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
	if _, err := corev1Client.Services(namespace).Create(&corev1.Service{
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
	}); err != nil {
		return err
	}

	apiserviceClient := apiregistrationv1client.NewForConfigOrDie(e.LoopbackClientConfig)
	if _, err := apiserviceClient.APIServices().Create(&apiregistrationv1.APIService{
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
			},
		},
	}); err != nil {
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
	}); err == nil {
		e.AggregatedAPIServerSession = session
		return nil
	}

	return fmt.Errorf("failed starting aggregated apiserver")
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
	return nil
}
