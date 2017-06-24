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

package boot

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os/exec"
)

var name, namespace string
var versions []schema.GroupVersion

var buildResourceConfigCmd = &cobra.Command{
	Use:   "build-resource-config",
	Short: "Create kubernetes resource config files to launch the apiserver.",
	Long:  `Create kubernetes resource config files to launch the apiserver.`,
	Run:   RunBuildResourceConfig,
}

func AddBuildResourceConfig(cmd *cobra.Command) {
	cmd.AddCommand(buildResourceConfigCmd)
	buildResourceConfigCmd.Flags().StringVar(&domain, "domain", "", "api groups domain")
	buildResourceConfigCmd.Flags().StringVar(&name, "name", "", "")
	buildResourceConfigCmd.Flags().StringVar(&namespace, "namespace", "", "")
	buildResourceConfigCmd.Flags().StringVar(&image, "image", "", "name of the apiserver image with tag")
}

func RunBuildResourceConfig(cmd *cobra.Command, args []string) {
	if len(name) == 0 {
		log.Fatalf("must specify --name")
	}
	if len(namespace) == 0 {
		log.Fatalf("must specify --namespace")
	}
	if len(image) == 0 {
		log.Fatalf("Must specify --image")
	}
	if len(domain) == 0 {
		domain = getDomain()
	}

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
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	initVersionedApis()
	log.Printf("versions %+v", versions)
	a := resourceConfigTemplateArgs{
		Name:       name,
		Namespace:  namespace,
		Image:      image,
		Domain:     domain,
		Versions:   versions,
		ClientKey:  getBase64(filepath.Join("bin", "certificates", "apiserver.key")),
		CACert:     getBase64(filepath.Join("bin", "certificates", "apiserver_ca.crt")),
		ClientCert: getBase64(filepath.Join("bin", "certificates", "apiserver.crt")),
	}

	root, err := os.Executable()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	root = filepath.Dir(root)

	doCmd("mkdir", "-p", filepath.Join(dir, "config"))

	path := filepath.Join(dir, "config", "apiserver.yaml")
	created := writeIfNotFound(path, "config-template", resourceConfigTemplate, a)
	if !created && !ignoreExists {
		log.Fatalf("Resource config already exists.")
	}
}

func createCerts() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dir = filepath.Join(dir, "bin", "certificates")

	doCmd("mkdir", "-p", dir)

	if _, err := os.Stat(filepath.Join(dir, "apiserver_ca.crt")); os.IsNotExist(err) {
		doCmd("openssl", "req", "-x509",
			"-newkey", "rsa:2048",
			"-keyout", filepath.Join(dir, "apiserver_ca.key"),
			"-out", filepath.Join(dir, "apiserver_ca.crt"),
			"-days", "365",
			"-nodes",
			"-subj", fmt.Sprintf("/C=/ST=/L=/O=/OU=/CN=%s-certificate-authority", name),
		)
	} else {
		log.Printf("Skipping generate CA cert.  File already exists.")
	}

	if _, err := os.Stat(filepath.Join(dir, "apiserver.csr")); os.IsNotExist(err) {
		// Use <service-name>.<namespace>.svc as the domain name for the certificate
		doCmd("openssl", "req",
			"-out", filepath.Join(dir, "apiserver.csr"),
			"-new",
			"-newkey", "rsa:2048",
			"-nodes",
			"-keyout", filepath.Join(dir, "apiserver.key"),
			"-subj", fmt.Sprintf("/C=/ST=/L=/O=/OU=/CN=%s.%s.svc", name, namespace),
		)
	} else {
		log.Printf("Skipping generate apiserver csr.  File already exists.")
	}

	if _, err := os.Stat(filepath.Join(dir, "apiserver.crt")); os.IsNotExist(err) {
		doCmd("openssl", "x509", "-req",
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
	if len(versionedAPIs) == 0 {
		groups, err := ioutil.ReadDir(filepath.Join("pkg", "apis"))
		if err != nil {
			log.Fatalf("could not read pkg/apis directory to find api versions")
		}
		for _, g := range groups {
			if g.IsDir() {
				versionFiles, err := ioutil.ReadDir(filepath.Join("pkg", "apis", g.Name()))
				if err != nil {
					log.Fatalf("could not read pkg/apis/%s directory to find api versions", g.Name())
				}
				versionMatch := regexp.MustCompile("^v\\d+(alpha\\d+|beta\\d+)*$")
				for _, v := range versionFiles {
					if v.IsDir() && versionMatch.MatchString(v.Name()) {
						versions = append(versions, schema.GroupVersion{
							Group:   g.Name(),
							Version: v.Name(),
						})
					}
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
        - "--print-bearer-token"
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
