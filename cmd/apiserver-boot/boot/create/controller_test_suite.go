package create

import (
	"path/filepath"
	"strings"

	"sigs.k8s.io/kubebuilder/pkg/scaffold/input"
	"sigs.k8s.io/kubebuilder/pkg/scaffold/resource"
)

// SuiteTest scaffolds a SuiteTest
type SuiteTest struct {
	input.Input

	// Resource is the Resource to make the Controller for
	Resource *resource.Resource
}

// GetInput implements input.File
func (a *SuiteTest) GetInput() (input.Input, error) {
	if a.Path == "" {
		a.Path = filepath.Join("pkg", "controller",
			strings.ToLower(a.Resource.Kind), strings.ToLower(a.Resource.Kind)+"_controller_suite_test.go")
	}
	a.TemplateBody = controllerSuiteTestTemplate
	return a.Input, nil
}

var controllerSuiteTestTemplate = `{{ .Boilerplate }}

package {{ lower .Resource.Kind }}

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	stdlog "k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/apiserver-builder-alpha/pkg/test/suite"
	"{{ .Repo }}/pkg/apis"
)

var cfg *rest.Config

func TestMain(m *testing.M) {

	env := suite.NewDefaultTestingEnvironment()
	if err := env.StartLocalKubeAPIServer(); err != nil {
		stdlog.Fatal(err)
		return
	}
	if err := env.StartLocalAggregatedAPIServer("{{ .Resource.Group }}.{{ .Domain }}", "{{ .Resource.Version }}"); err != nil {
		stdlog.Fatal(err)
		return
	}

	apis.AddToScheme(scheme.Scheme)
	cfg = env.LoopbackClientConfig

	code := m.Run()

	stdlog.Info("stopping aggregated-apiserver..")
	if err := env.StopLocalAggregatedAPIServer(); err != nil {
		stdlog.Fatal(err)
		return
	}
	stdlog.Info("stopping kube-apiserver..")
	if err := env.StopLocalKubeAPIServer(); err != nil {
		stdlog.Fatal(err)
		return
	}

	os.Exit(code)
}

// SetupTestReconcile returns a reconcile.Reconcile implementation that delegates to inner and
// writes the request to requests after Reconcile is finished.
func SetupTestReconcile(inner reconcile.Reconciler) (reconcile.Reconciler, chan reconcile.Request) {
	requests := make(chan reconcile.Request)
	fn := reconcile.Func(func(req reconcile.Request) (reconcile.Result, error) {
		result, err := inner.Reconcile(req)
		requests <- req
		return result, err
	})
	return fn, requests
}

// StartTestManager adds recFn
func StartTestManager(mgr manager.Manager, g *gomega.GomegaWithT) (chan struct{}, *sync.WaitGroup) {
	stop := make(chan struct{})
	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		g.Expect(mgr.Start(stop)).NotTo(gomega.HaveOccurred())
		wg.Done()
	}()
	return stop, wg
}
`
