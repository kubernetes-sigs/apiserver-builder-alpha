/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package build

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"

	"github.com/kubernetes-incubator/apiserver-builder/cmd/apiserver-boot/boot/util"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var Name, Namespace string
var Versions []schema.GroupVersion
var ResourceConfigDir string

var buildResourceConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Create kubernetes resource config files to launch the apiserver.",
	Long:  `Create kubernetes resource config files to launch the apiserver.`,
	Example: `
# Build yaml resource config into the config/ directory for running the apiserver and
# controller-manager as an aggregated service in a Kubernetes cluster
# Generates CA and apiserver certificates.
apiserver-boot build config --name nameofservice --namespace mysystemnamespace --image gcr.io/myrepo/myimage:mytag
`,
	Run: RunBuildResourceConfig,
}

func AddBuildResourceConfig(cmd *cobra.Command) {
	cmd.AddCommand(buildResourceConfigCmd)
	AddBuildResourceConfigFlags(buildResourceConfigCmd)
}

func AddBuildResourceConfigFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Name, "name", "", "")
	cmd.Flags().StringVar(&Namespace, "namespace", "", "")
	cmd.Flags().StringVar(&Image, "image", "", "name of the apiserver Image with tag")
	cmd.Flags().StringVar(&ResourceConfigDir, "output", "config", "directory to output resourceconfig")
}

func RunBuildResourceConfig(cmd *cobra.Command, args []string) {
	if len(Name) == 0 {
		log.Fatalf("must specify --name")
	}
	if len(Namespace) == 0 {
		log.Fatalf("must specify --namespace")
	}
	if len(Image) == 0 {
		log.Fatalf("Must specify --image")
	}
	util.GetDomain()

	if _, err := os.Stat("pkg"); err != nil {
		log.Fatalf("could not find 'pkg' directory.  must run apiserver-boot init before generating config")
	}

	createCerts()
	buildResourceConfig()
}

func getBase64(file string) string {
	out, err := exec.Command("bash", "-c",
		fmt.Sprintf("base64 %s | awk 'BEGIN{ORS=\"\";} {print}'", file)).CombinedOutput()
	if err != nil {
		log.Fatalf("Could not base64 encode file: %v", err)
	}
	return string(out)
}

func buildResourceConfig() {
	initVersionedApis()
	dir := filepath.Join(ResourceConfigDir, "certificates")

	a := resourceConfigTemplateArgs{
		Name:       Name,
		Namespace:  Namespace,
		Image:      Image,
		Domain:     util.Domain,
		Versions:   Versions,
		ClientKey:  getBase64(filepath.Join(dir, "apiserver.key")),
		CACert:     getBase64(filepath.Join(dir, "apiserver_ca.crt")),
		ClientCert: getBase64(filepath.Join(dir, "apiserver.crt")),
	}
	path := filepath.Join(ResourceConfigDir, "apiserver.yaml")

	created := util.WriteIfNotFound(path, "config-template", resourceConfigTemplate, a)
	if !created {
		log.Fatalf("Resource config already exists.")
	}
}

func createCerts() {
	dir := filepath.Join(ResourceConfigDir, "certificates")
	util.DoCmd("mkdir", "-p", dir)

	if _, err := os.Stat(filepath.Join(dir, "apiserver_ca.crt")); os.IsNotExist(err) {
		util.DoCmd("openssl", "req", "-x509",
			"-newkey", "rsa:2048",
			"-keyout", filepath.Join(dir, "apiserver_ca.key"),
			"-out", filepath.Join(dir, "apiserver_ca.crt"),
			"-days", "365",
			"-nodes",
			"-subj", fmt.Sprintf("/C=/ST=/L=/O=/OU=/CN=%s-certificate-authority", Name),
		)
	} else {
		log.Printf("Skipping generate CA cert.  File already exists.")
	}

	if _, err := os.Stat(filepath.Join(dir, "apiserver.csr")); os.IsNotExist(err) {
		// Use <service-Name>.<Namespace>.svc as the domain Name for the certificate
		util.DoCmd("openssl", "req",
			"-out", filepath.Join(dir, "apiserver.csr"),
			"-new",
			"-newkey", "rsa:2048",
			"-nodes",
			"-keyout", filepath.Join(dir, "apiserver.key"),
			"-subj", fmt.Sprintf("/C=/ST=/L=/O=/OU=/CN=%s.%s.svc", Name, Namespace),
		)
	} else {
		log.Printf("Skipping generate apiserver csr.  File already exists.")
	}

	if _, err := os.Stat(filepath.Join(dir, "apiserver.crt")); os.IsNotExist(err) {
		util.DoCmd("openssl", "x509", "-req",
			"-days", "365",
			"-in", filepath.Join(dir, "apiserver.csr"),
			"-CA", filepath.Join(dir, "apiserver_ca.crt"),
			"-CAkey", filepath.Join(dir, "apiserver_ca.key"),
			"-CAcreateserial",
			"-out", filepath.Join(dir, "apiserver.crt"),
		)
	} else {
		log.Printf("Skipping signing apiserver crt.  File already exists.")
	}
}

func initVersionedApis() {
	groups, err := ioutil.ReadDir(filepath.Join("pkg", "apis"))
	if err != nil {
		log.Fatalf("could not read pkg/apis directory to find api Versions")
	}
	log.Printf("Adding APIs:")
	for _, g := range groups {
		if g.IsDir() {
			versionFiles, err := ioutil.ReadDir(filepath.Join("pkg", "apis", g.Name()))
			if err != nil {
				log.Fatalf("could not read pkg/apis/%s directory to find api Versions", g.Name())
			}
			versionMatch := regexp.MustCompile("^v\\d+(alpha\\d+|beta\\d+)*$")
			for _, v := range versionFiles {
				if v.IsDir() && versionMatch.MatchString(v.Name()) {
					log.Printf("\t%s.%s", g.Name(), v.Name())
					Versions = append(Versions, schema.GroupVersion{
						Group:   g.Name(),
						Version: v.Name(),
					})
				}
			}
		}
	}
	u := map[string]bool{}
	for _, a := range versionedAPIs {
		u[path.Dir(a)] = true
	}
	for a, _ := range u {
		unversionedAPIs = append(unversionedAPIs, a)
	}
}

type resourceConfigTemplateArgs struct {
	Versions   []schema.GroupVersion
	CACert     string
	ClientCert string
	ClientKey  string
	Domain     string
	Name       string
	Namespace  string
	Image      string
}

var resourceConfigTemplate = `
{{ $config := . -}}
{{ range $api := .Versions -}}
apiVersion: apiregistration.k8s.io/v1beta1
kind: APIService
metadata:
  name: {{ $api.Version }}.{{ $api.Group }}.{{ $config.Domain }}
  labels:
    api: {{ $config.Name }}
    apiserver: "true"
spec:
  version: {{ $api.Version }}
  group: {{ $api.Group }}.{{ $config.Domain }}
  groupPriorityMinimum: 2000
  priority: 200
  service:
    name: {{ $config.Name }}
    namespace: {{ $config.Namespace }}
  versionPriority: 10
  caBundle: "{{ $config.CACert }}"
---
{{ end -}}
apiVersion: v1
kind: Service
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
  labels:
    api: {{.Name}}
    apiserver: "true"
spec:
  ports:
  - port: 443
    protocol: TCP
    targetPort: 443
  selector:
    api: {{ .Name }}
    apiserver: "true"
---
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
  labels:
    api: {{.Name}}
    apiserver: "true"
spec:
  replicas: 1
  template:
    metadata:
      labels:
        api: {{.Name}}
        apiserver: "true"
    spec:
      containers:
      - name: apiserver
        image: {{.Image}}
        volumeMounts:
        - name: apiserver-certs
          mountPath: /apiserver.local.config/certificates
          readOnly: true
        command:
        - "./apiserver"
        args:
        - "--etcd-servers=http://localhost:2379"
        - "--tls-cert-file=/apiserver.local.config/certificates/tls.crt"
        - "--tls-private-key-file=/apiserver.local.config/certificates/tls.key"
        - "--audit-log-path=-"
        - "--audit-log-maxage=0"
        - "--audit-log-maxbackup=0"
      - name: controller
        image: {{.Image}}
        command:
        - "./controller-manager"
        args:
      - name: etcd
        image: quay.io/coreos/etcd:v3.0.17
      volumes:
      - name: apiserver-certs
        secret:
          secretName: {{ .Name }}
---
apiVersion: v1
kind: Secret
type: kubernetes.io/tls
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
  labels:
    api: {{.Name}}
    apiserver: "true"
data:
  tls.crt: {{ .ClientCert }}
  tls.key: {{ .ClientKey }}
---
`
