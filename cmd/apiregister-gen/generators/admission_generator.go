package generators

import (
	"fmt"
	"io"
	"k8s.io/klog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"k8s.io/gengo/generator"
)

var _ generator.Generator = &apiGenerator{}

type admissionGenerator struct {
	generator.DefaultGen
	projectRootPath string
	admissionKinds  []string
}

func CreateAdmissionGenerator(apis *APIs, filename string, projectRootPath string, outputBase string) generator.Generator {
	admissionKinds := []string{}
	// filter out those resources created w/ `--admission-controller` flag
	for _, group := range apis.Groups {
		for _, version := range group.Versions {
			for _, resource := range version.Resources {
				resourceAdmissionControllerPkg := filepath.Join(outputBase, projectRootPath, "plugin", "admission", strings.ToLower(resource.Kind))
				// if "<repo>/plugin/admission" package is present in the project, add it to the generated installation function
				if _, err := os.Stat(resourceAdmissionControllerPkg); err == nil {
					admissionKinds = append(admissionKinds, resource.Kind)
					klog.V(5).Infof("found existing admission controller for resource: %v/%v", resource.Group, resource.Kind)
				}
			}
		}
	}

	return &admissionGenerator{
		generator.DefaultGen{OptionalName: filename},
		projectRootPath,
		admissionKinds,
	}
}

func (d *admissionGenerator) Imports(c *generator.Context) []string {
	imports := []string{
		"sigs.k8s.io/apiserver-builder-alpha/pkg/cmd/server",
		"k8s.io/client-go/rest",
		`genericserver "k8s.io/apiserver/pkg/server"`,
		"k8s.io/apiserver/pkg/admission",
	}
	for _, kind := range d.admissionKinds {
		imports = append(imports, fmt.Sprintf(
			`. "%s/plugin/admission/%s"`, d.projectRootPath, strings.ToLower(kind)))
	}
	imports = append(imports,
		fmt.Sprintf(`aggregatedclientset "%s/pkg/client/clientset_generated/clientset"`, d.projectRootPath))
	imports = append(imports,
		fmt.Sprintf(`aggregatedinformerfactory "%s/pkg/client/informers_generated/externalversions"`, d.projectRootPath))
	imports = append(imports,
		fmt.Sprintf(`initializer "%s/plugin/admission"`, d.projectRootPath))
	return imports
}

type AdmissionGeneratorParam struct {
	Kind string
}

func (d *admissionGenerator) Finalize(context *generator.Context, w io.Writer) error {
	if len(d.admissionKinds) == 0 {
		return nil
	}

	temp := template.Must(template.New("admission-install-template").Parse(AdmissionsInstallTemplate))
	return temp.Execute(w, &struct {
		Admissions []string
	}{
		Admissions: d.admissionKinds,
	})
}

var AdmissionsInstallTemplate = `
func init() {
	server.AggregatedAdmissionInitializerGetter = GetAggregatedResourceAdmissionControllerInitializer
{{ range .Admissions -}}
	server.AggregatedAdmissionPlugins["{{.}}"] = New{{.}}Plugin()
{{ end }}
}

func GetAggregatedResourceAdmissionControllerInitializer(config *rest.Config) (admission.PluginInitializer, genericserver.PostStartHookFunc) {
	// init aggregated resource clients
	aggregatedResourceClient := aggregatedclientset.NewForConfigOrDie(config)
	aggregatedInformerFactory := aggregatedinformerfactory.NewSharedInformerFactory(aggregatedResourceClient, 0)
	aggregatedResourceInitializer := initializer.New(aggregatedResourceClient, aggregatedInformerFactory)

	return aggregatedResourceInitializer, func(context genericserver.PostStartHookContext) error {
		aggregatedInformerFactory.Start(context.StopCh)
		return nil
	}
}
`
