package podexec_test

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	podexecv1 "sigs.k8s.io/apiserver-builder-alpha/example/podexec/pkg/apis/podexec/v1"
)

func ExamplePodExecSimple() {

	restConfig, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	client := kubernetes.NewForConfigOrDie(restConfig)
	request := client.CoreV1().RESTClient().
		Post().
		AbsPath("/apis/podexec.example.com/v1/namespaces/default/pods/myapp-pod/exec").
		VersionedParams(&podexecv1.PodExec{
			Container: "myapp-container",
			Command:   []string{"echo", "yes!"},
			Stdout:    true,
			Stderr:    true,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(restConfig, "POST", request.URL())
	if err != nil {
		panic(err)
	}
	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: outBuf,
		Stderr: errBuf,
		Tty:    false,
	})
	if err != nil {
		panic(err)
	}

	responseBody, err := ioutil.ReadAll(outBuf)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(responseBody))

}
