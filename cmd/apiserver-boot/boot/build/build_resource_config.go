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
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog"
	"sigs.k8s.io/apiserver-builder-alpha/cmd/apiserver-boot/boot/util"
)

var Name, Namespace string
var Versions []schema.GroupVersion
var ResourceConfigDir string
var ControllerArgs []string
var ApiserverArgs []string
var ControllerSecret string
var ControllerSecretMount string
var ControllerSecretEnv []string
var ImagePullSecrets []string
var ServiceAccount string
var StorageClass string

var buildResourceConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Create kubernetes resource config files to launch the apiserver.",
	Long:  `Create kubernetes resource config files to launch the apiserver.`,
	Example: `
# Build yaml resource config into the config/ directory for running the apiserver and
# controller-manager as an aggregated service in a Kubernetes cluster as a container.
# Generates CA and apiserver certificates.
apiserver-boot build config --name nameofservice --namespace mysystemnamespace --image gcr.io/myrepo/myimage:mytag

# Build yaml resource config into the config/ directory for running the apiserver and
# controller-manager locally, but registered through aggregation into a local minikube cluster
# Generates CA and apiserver certificates.
apiserver-boot build config --name nameofservice --namespace mysystemnamespace --local-minikube
`,
	Run: RunBuildResourceConfig,
}

func AddBuildResourceConfig(cmd *cobra.Command) {
	cmd.AddCommand(buildResourceConfigCmd)
	AddBuildResourceConfigFlags(buildResourceConfigCmd)
}

func AddBuildResourceConfigFlags(cmd *cobra.Command) {
	cmd.Flags().StringSliceVar(&ControllerSecretEnv, "controller-env", []string{}, "")
	cmd.Flags().StringVar(&ControllerSecret, "controller-secret", "", "")
	cmd.Flags().StringVar(&ControllerSecretMount, "controller-secret-mount", "", "")
	cmd.Flags().StringSliceVar(&ControllerArgs, "controller-args", []string{}, "")
	cmd.Flags().StringSliceVar(&ApiserverArgs, "apiserver-args", []string{}, "")
	cmd.Flags().StringVar(&Name, "name", "", "")
	cmd.Flags().StringVar(&Namespace, "namespace", "", "")
	cmd.Flags().StringSliceVar(&ImagePullSecrets, "image-pull-secrets", []string{}, "List of secret names for docker registry")
	cmd.Flags().StringVar(&ServiceAccount, "service-account", "", "Name of service account that will be attached to deployed pod")
	cmd.Flags().StringVar(&Image, "image", "", "name of the apiserver Image with tag")
	cmd.Flags().StringVar(&ResourceConfigDir, "output", "config", "directory to output resourceconfig")
	cmd.Flags().StringVar(&StorageClass, "storage-class", "standard", "storageclass of which etcd is using to store data")
}

func RunBuildResourceConfig(cmd *cobra.Command, args []string) {
	if len(Name) == 0 {
		klog.Fatalf("must specify --name")
	}
	if len(Namespace) == 0 {
		klog.Fatalf("must specify --namespace")
	}
	if len(Image) == 0 {
		klog.Fatalf("Must specify --image")
	}
	util.GetDomain()

	if _, err := os.Stat("pkg"); err != nil {
		klog.Fatalf("could not find 'pkg' directory.  must run apiserver-boot init before generating config")
	}

	createCerts()
	buildResourceConfig()
}

func getBase64(file string) string {
	//out, err := exec.Command("bash", "-c",
	//	fmt.Sprintf("base64 %s | awk 'BEGIN{ORS=\"\";} {print}'", file)).CombinedOutput()
	//if err != nil {
	//	klog.Fatalf("Could not base64 encode file: %v", err)
	//}

	buff := bytes.Buffer{}
	enc := base64.NewEncoder(base64.StdEncoding, &buff)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		klog.Fatalf("Could not read file %s: %v", file, err)
	}

	_, err = enc.Write(data)
	if err != nil {
		klog.Fatalf("Could not write bytes: %v", err)
	}
	enc.Close()
	return buff.String()

	//if string(out) != buff.String() {
	//	fmt.Printf("\nNot Equal\n")
	//}
	//
	//return string(out)
}

func buildResourceConfig() {
	initVersionedApis()
	dir := filepath.Join(ResourceConfigDir, "certificates")

	created := util.WriteIfNotFound(
		filepath.Join(ResourceConfigDir, "apiservice.yaml"),
		"apiservice-config-template", apiserviceYamlTemplate, apiserviceYamlTemplateArgs{
			Name:      Name,
			Namespace: Namespace,
			Domain:    util.Domain,
			Versions:  Versions,
			CACert:    getBase64(filepath.Join(dir, "apiserver_ca.crt")),
		})
	if !created {
		klog.Warningf("Resource config already exists.")
	}

	// build apiserver yaml config
	created = util.WriteIfNotFound(
		filepath.Join(ResourceConfigDir, "aggregated-apiserver.yaml"),
		"apiserver-config-template", resourceConfigApiserverYaml, resourceConfigApiserverYamlArgs{
			Name:             Name,
			Namespace:        Namespace,
			Image:            Image,
			ApiserverArgs:    ApiserverArgs,
			ImagePullSecrets: ImagePullSecrets,
			ServiceAccount:   ServiceAccount,
			ClientKey:        getBase64(filepath.Join(dir, "apiserver.key")),
			ClientCert:       getBase64(filepath.Join(dir, "apiserver.crt")),
		})
	if !created {
		klog.Warningf("Aggregated Apiserver config already exists.")
	}

	// build controller yaml config
	created = util.WriteIfNotFound(
		filepath.Join(ResourceConfigDir, "controller-manager.yaml"),
		"controller-config-template", resourceConfigControllerYaml, resourceConfigControllerYamlArgs{
			Name:             Name,
			Namespace:        Namespace,
			Image:            Image,
			ControllerArgs:   ControllerArgs,
			ImagePullSecrets: ImagePullSecrets,
			ServiceAccount:   ServiceAccount,
		})
	if !created {
		klog.Warningf("Controller-manager config already exists.")
	}

	// build RBAC yaml config
	created = util.WriteIfNotFound(
		filepath.Join(ResourceConfigDir, "rbac.yaml"),
		"rbac-config-template", resourceConfigRBACYaml, resourceConfigRBACYamlArgs{
			Name:      Name,
			Namespace: Namespace,
			Domain:    util.Domain,
			Versions:  Versions,
		})
	if !created {
		klog.Warningf("RBAC config already exists.")
	}

	// build etcd yaml config
	created = util.WriteIfNotFound(
		filepath.Join(ResourceConfigDir, "etcd.yaml"),
		"etcd-config-template", etcdYaml, etcdYamlArgs{
			Namespace:    Namespace,
			StorageClass: StorageClass,
		})
	if !created {
		klog.Warningf("ETCD config already exists.")
	}
}

func createCerts() {
	dir := filepath.Join(ResourceConfigDir, "certificates")
	os.MkdirAll(dir, 0700)

	svrName := fmt.Sprintf("%s.%s.svc", Name, Namespace)

	if _, err := os.Stat(filepath.Join(dir, "apiserver_ca.crt")); os.IsNotExist(err) {
		util.DoCmd("openssl", "req", "-x509",
			"-newkey", "rsa:2048",
			"-addext", "basicConstraints=critical,CA:TRUE,pathlen:1",
			"-keyout", filepath.Join(dir, "apiserver_ca.key"),
			"-out", filepath.Join(dir, "apiserver_ca.crt"),
			"-days", "365",
			"-nodes",
			"-subj", fmt.Sprintf("/C=un/ST=st/L=l/O=o/OU=ou/CN=%s-certificate-authority", Name),
		)
	} else {
		klog.Infof("Skipping generate CA cert.  File already exists.")
	}

	caCert, caKey, err := util.TryLoadCertAndKeyFromDisk(dir, "apiserver_ca")
	if err != nil {
		klog.Fatal(err)
	}

	apiserverCert, apiserverKey, err := util.NewCertAndKey(caCert, caKey, util.Config{
		CommonName:   svrName,
		Organization: []string{},
		AltNames: util.AltNames{
			DNSNames: []string{
				"localhost",
				svrName,
			},
			IPs: []net.IP{
				net.ParseIP("127.0.0.1"),
			},
		},
		Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	})
	if err != nil {
		klog.Fatal(err)
	}

	apiserverCertData := util.EncodeCertPEM(apiserverCert)
	apiserverKeyData := util.EncodePrivateKeyPEM(apiserverKey)

	if err := ioutil.WriteFile(filepath.Join(dir, "apiserver.crt"), apiserverCertData, 0644); err != nil {
		klog.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join(dir, "apiserver.key"), apiserverKeyData, 0644); err != nil {
		klog.Fatal(err)
	}
}

func initVersionedApis() {
	groups, err := ioutil.ReadDir(filepath.Join("pkg", "apis"))
	if err != nil {
		klog.Fatalf("could not read pkg/apis directory to find api Versions")
	}
	klog.Infof("Adding APIs:")
	for _, g := range groups {
		if g.IsDir() {
			versionFiles, err := ioutil.ReadDir(filepath.Join("pkg", "apis", g.Name()))
			if err != nil {
				klog.Fatalf("could not read pkg/apis/%s directory to find api Versions", g.Name())
			}
			versionMatch := regexp.MustCompile("^v\\d+(alpha\\d+|beta\\d+)*$")
			for _, v := range versionFiles {
				if v.IsDir() && versionMatch.MatchString(v.Name()) {
					klog.Infof("\t%s.%s", g.Name(), v.Name())
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
	for a := range u {
		unversionedAPIs = append(unversionedAPIs, a)
	}
}

type resourceConfigApiserverYamlArgs struct {
	Name      string
	Namespace string

	Image            string
	ServiceAccount   string
	ImagePullSecrets []string
	ApiserverArgs    []string
	ClientCert       string
	ClientKey        string
}

var resourceConfigApiserverYaml = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Name}}-apiserver
  namespace: {{.Namespace}}
  labels:
    api: {{.Name}}
    apiserver: "true"
spec:
  selector:
    matchLabels:
      api: {{.Name}}
      apiserver: "true"
  replicas: 1
  template:
    metadata:
      labels:
        api: {{.Name}}
        apiserver: "true"
    spec:
      {{- if .ImagePullSecrets }}
      imagePullSecrets:
      {{range .ImagePullSecrets }}- name: {{.}}
      {{ end }}
      {{- end -}}
      {{- if .ServiceAccount }}
      serviceAccount: {{.ServiceAccount}}
      {{- end }}
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
        - "--etcd-servers=http://etcd-svc:2379"
        - "--tls-cert-file=/apiserver.local.config/certificates/tls.crt"
        - "--tls-private-key-file=/apiserver.local.config/certificates/tls.key"
        - "--audit-log-path=-"
        - "--audit-log-maxage=0"
        - "--audit-log-maxbackup=0"{{ range $arg := .ApiserverArgs }}
        - "{{ $arg }}"{{ end }}
        resources:
          requests:
            cpu: 100m
            memory: 20Mi
          limits:
            cpu: 100m
            memory: 30Mi
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
`

type resourceConfigControllerYamlArgs struct {
	Name      string
	Namespace string

	Image            string
	ServiceAccount   string
	ImagePullSecrets []string
	ControllerArgs   []string
}

var resourceConfigControllerYaml = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Name}}-controller
  namespace: {{.Namespace}}
  labels:
    api: {{.Name}}
    controller: "true"
spec:
  selector:
    matchLabels:
      api: {{.Name}}
      controller: "true"
  replicas: 1
  template:
    metadata:
      labels:
        api: {{.Name}}
        controller: "true"
    spec:
      {{- if .ImagePullSecrets }}
      imagePullSecrets:
      {{range .ImagePullSecrets }}- name: {{.}}
      {{ end }}
      {{- end -}}
      {{- if .ServiceAccount }}
      serviceAccount: {{.ServiceAccount}}
      {{- end }}
      containers:
      - name: controller
        image: {{.Image}}
        command:
        - "./controller-manager"
        args:{{ range $arg := .ControllerArgs }}
        - "{{ $arg }}"{{ end }}
        resources:
          requests:
            cpu: 100m
            memory: 200Mi
          limits:
            cpu: 100m
            memory: 300Mi
      volumes:
      - name: apiserver-certs
        secret:
          secretName: {{ .Name }}
`

type resourceConfigRBACYamlArgs struct {
	Name      string
	Namespace string
	Domain    string
	Versions  []schema.GroupVersion
}

var resourceConfigRBACYaml = `---
{{ $config := . -}}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: {{.Name}}-apiserver-auth-reader
rules:
  - apiGroups:
      - ""
    resourceNames:
      - extension-apiserver-authentication
    resources:
      - configmaps
    verbs:
      - get
      - list
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: {{.Name}}-apiserver-auth-reader
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: {{.Name}}-apiserver-auth-reader
subjects:
  - kind: ServiceAccount
    namespace: default
    name: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: {{.Name}}-apiserver-auth-delegator
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:auth-delegator
subjects:
  - kind: ServiceAccount
    namespace: default
    name: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: {{.Name}}-controller
  namespace: default
rules:
  - apiGroups:
{{- range $api := .Versions }}
      - '{{ $api.Group }}.{{ $config.Domain }}'
{{- end }}
    resources:
      - '*'
    verbs:
      - '*'
  - apiGroups:
      - ''
    resources:
      - 'configmaps'
      - 'namespaces'
    verbs:
      - 'get'
      - 'list'
      - 'watch'
  - apiGroups:
      - 'admissionregistration.k8s.io'
    resources:
      - '*'
    verbs:
      - 'list'
      - 'watch'
  - nonResourceURLs:
      - '*'
    verbs:
      - '*'

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: {{.Name}}-controller
  namespace: default
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: {{.Name}}-controller
subjects:
  - kind: ServiceAccount
    namespace: default
    name: default
`

type etcdYamlArgs struct {
	Namespace    string
	StorageClass string
}

var etcdYaml = `---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: etcd
  namespace: {{ .Namespace }}
spec:
  selector:
    matchLabels:
      app: etcd
  serviceName: "etcd"
  replicas: 1
  template:
    metadata:
      labels:
        app: etcd
    spec:
      terminationGracePeriodSeconds: 10
      containers:
      - name: etcd
        image: quay.io/coreos/etcd:latest
        imagePullPolicy: Always
        resources:
          requests:
            cpu: 100m
            memory: 20Mi
          limits:
            cpu: 100m
            memory: 30Mi
        env:
        - name: ETCD_DATA_DIR
          value: /etcd-data-dir
        command:
        - /usr/local/bin/etcd
        - --listen-client-urls
        - http://0.0.0.0:2379
        - --advertise-client-urls
        - http://localhost:2379
        ports:
        - containerPort: 2379
        volumeMounts:
        - name: etcd-data-dir
          mountPath: /etcd-data-dir
        readinessProbe:
          httpGet:
            port: 2379
            path: /health
          failureThreshold: 1
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 2
        livenessProbe:
          httpGet:
            port: 2379
            path: /health
          failureThreshold: 3
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 2
  volumeClaimTemplates:
  - metadata:
     name: etcd-data-dir
     annotations:
        volume.beta.kubernetes.io/storage-class: {{.StorageClass}}
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
         storage: 10Gi
---
apiVersion: v1
kind: Service
metadata:
  name: etcd-svc
  namespace: {{ .Namespace }}
  labels:
    app: etcd
spec:
  ports:
  - port: 2379
    name: etcd
    targetPort: 2379
  selector:
    app: etcd
`

type apiserviceYamlTemplateArgs struct {
	Versions  []schema.GroupVersion
	CACert    string
	Domain    string
	Name      string
	Namespace string
}

var apiserviceYamlTemplate = `
{{ $config := . -}}
{{ range $api := .Versions -}}
apiVersion: apiregistration.k8s.io/v1
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
  service:
    name: {{ $config.Name }}
    namespace: {{ $config.Namespace }}
  versionPriority: 10
  caBundle: "{{ $config.CACert }}"
---
{{ end -}}
`

var localConfigTemplate = `
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
  type: ExternalName
  externalName: "{{ .LocalIp }}"
  ports:
  - port: 443
    protocol: TCP
`
