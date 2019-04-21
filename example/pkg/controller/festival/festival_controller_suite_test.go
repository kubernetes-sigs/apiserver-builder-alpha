package festival

import (
	stdlog "log"
	"os"
	"sync"
	"testing"

	"github.com/kubernetes-incubator/apiserver-builder-alpha/example/pkg/apis"
	"github.com/kubernetes-incubator/apiserver-builder-alpha/pkg/test/suite"
	"github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var cfg *rest.Config

func TestMain(m *testing.M) {
	t, err := suite.InstallLocalTestingAPIAggregationEnvironment("kingsport.k8s.io", "v1")
	if err != nil {
		stdlog.Fatal(err)
		return
	}
	apis.AddToScheme(scheme.Scheme)
	cfg = t.LoopbackClientConfig

	code := m.Run()

	stdlog.Print("stopping aggregated-apiserver..")
	if err := t.StopAggregatedAPIServer(); err != nil {
		stdlog.Fatal(err)
		return
	}
	stdlog.Print("stopping kube-apiserver..")
	if err := t.KubeAPIServerEnvironment.Stop(); err != nil {
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
