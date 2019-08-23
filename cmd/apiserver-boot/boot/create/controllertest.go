package create

import (
	"path/filepath"
	"sigs.k8s.io/kubebuilder/pkg/scaffold/input"
	"sigs.k8s.io/kubebuilder/pkg/scaffold/resource"
	"strings"
)

// Test scaffolds a Controller Test
type Test struct {
	input.Input

	// Resource is the Resource to make the Controller for
	Resource *resource.Resource

	// ResourcePackage is the package of the Resource
	ResourcePackage string
}

// GetInput implements input.File
func (a *Test) GetInput() (input.Input, error) {
	if a.Path == "" {
		a.Path = filepath.Join("pkg", "controller",
			strings.ToLower(a.Resource.Kind), strings.ToLower(a.Resource.Kind)+"_controller_test.go")
	}

	// Use the k8s.io/api package for core resources
	coreGroups := map[string]string{
		"apps":                  "",
		"admissionregistration": "k8s.io",
		"apiextensions":         "k8s.io",
		"authentication":        "k8s.io",
		"autoscaling":           "",
		"batch":                 "",
		"certificates":          "k8s.io",
		"core":                  "",
		"extensions":            "",
		"metrics":               "k8s.io",
		"policy":                "",
		"rbac.authorization":    "k8s.io",
		"storage":               "k8s.io",
	}

	a.ResourcePackage, _ = getResourceInfo(coreGroups, a.Resource, a.Input)

	a.TemplateBody = controllerTestTemplate
	a.Input.IfExistsAction = input.Error
	return a.Input, nil
}

var controllerTestTemplate = `{{ .Boilerplate }}

package {{ lower .Resource.Kind }}

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
	"golang.org/x/net/context"
	{{ if .Resource.CreateExampleReconcileBody -}}
	appsv1 "k8s.io/api/apps/v1"
	{{ end -}}
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	{{ .Resource.Group}}{{ .Resource.Version }} "{{ .ResourcePackage }}/{{ .Resource.Group}}/{{ .Resource.Version }}"
)

var c client.Client

var expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: "foo"{{ if .Resource.Namespaced }}, Namespace: "default"{{end}}}}
{{ if .Resource.CreateExampleReconcileBody }}var depKey = types.NamespacedName{Name: "foo-deployment", Namespace: "default"}
{{ end }}
const timeout = time.Second * 5

func TestReconcile(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	instance := &{{ .Resource.Group }}{{ .Resource.Version }}.{{ .Resource.Kind }}{ObjectMeta: metav1.ObjectMeta{Name: "foo"{{ if .Resource.Namespaced }}, Namespace: "default"{{end}}}}

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{MetricsBindAddress:"0"})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c = mgr.GetClient()

	recFn, requests := SetupTestReconcile(newReconciler(mgr))
	g.Expect(add(mgr, recFn)).To(gomega.Succeed())

	stopMgr, mgrStopped := StartTestManager(mgr, g)

	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	// Create the {{ .Resource.Kind }} object and expect the Reconcile{{ if .Resource.CreateExampleReconcileBody }} and Deployment to be created{{ end }}
	err = c.Create(context.TODO(), instance)
	// The instance object may not be a valid object because it might be missing some required fields.
	// Please modify the instance object by adding required fields and then remove the following if statement.
	if apierrors.IsInvalid(err) {
		t.Logf("failed to create object, got an invalid object error: %v", err)
		return
	}
	g.Expect(err).NotTo(gomega.HaveOccurred())
	defer c.Delete(context.TODO(), instance)
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))
{{ if .Resource.CreateExampleReconcileBody }}
	deploy := &appsv1.Deployment{}
	g.Eventually(func() error { return c.Get(context.TODO(), depKey, deploy) }, timeout).
		Should(gomega.Succeed())

	// Delete the Deployment and expect Reconcile to be called for Deployment deletion
	g.Expect(c.Delete(context.TODO(), deploy)).To(gomega.Succeed())
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))
	g.Eventually(func() error { return c.Get(context.TODO(), depKey, deploy) }, timeout).
		Should(gomega.Succeed())

	// Manually delete Deployment since GC isn't enabled in the test control plane
	g.Eventually(func() error { return c.Delete(context.TODO(), deploy) }, timeout).
	    Should(gomega.MatchError("deployments.apps \"foo-deployment\" not found"))
{{ end }}
}
`
