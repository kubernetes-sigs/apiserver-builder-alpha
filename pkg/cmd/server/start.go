/*
Copyright 2016 The Kubernetes Authors.

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

package server

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"k8s.io/apimachinery/pkg/runtime/schema"
	genericapifilters "k8s.io/apiserver/pkg/endpoints/filters"
	openapinamer "k8s.io/apiserver/pkg/endpoints/openapi"
	"k8s.io/apiserver/pkg/features"
	"k8s.io/apiserver/pkg/server"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericfilters "k8s.io/apiserver/pkg/server/filters"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/apiserver/pkg/util/feature"
	"k8s.io/apiserver/pkg/util/logs"
	"k8s.io/klog"
	openapi "k8s.io/kube-openapi/pkg/common"

	"github.com/kubernetes-incubator/apiserver-builder-alpha/pkg/apiserver"
	"github.com/kubernetes-incubator/apiserver-builder-alpha/pkg/builders"
	"github.com/kubernetes-incubator/apiserver-builder-alpha/pkg/validators"
)

var GetOpenApiDefinition openapi.GetOpenAPIDefinitions

type ServerOptions struct {
	RecommendedOptions     *genericoptions.RecommendedOptions
	APIBuilders            []*builders.APIGroupBuilder
	InsecureServingOptions *genericoptions.DeprecatedInsecureServingOptionsWithLoopback

	PrintBearerToken bool
	PrintOpenapi     bool
	RunDelegatedAuth bool
	BearerToken      string
	Kubeconfig       string
	PostStartHooks   []PostStartHook
}

type PostStartHook struct {
	Fn   genericapiserver.PostStartHookFunc
	Name string
}

// StartApiServer starts an apiserver hosting the provider apis and openapi definitions.
func StartApiServer(etcdPath string, apis []*builders.APIGroupBuilder, openapidefs openapi.GetOpenAPIDefinitions, title, version string, tweakConfigFuncs ...func(apiServer *apiserver.Config) error) {
	logs.InitLogs()
	defer logs.FlushLogs()

	GetOpenApiDefinition = openapidefs

	signalCh := genericapiserver.SetupSignalHandler()
	// To disable providers, manually specify the list provided by getKnownProviders()
	cmd, _ := NewCommandStartServer(etcdPath, os.Stdout, os.Stderr, apis, signalCh, title, version, tweakConfigFuncs...)

	cmd.Flags().AddFlagSet(pflag.CommandLine)
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

func NewServerOptions(etcdPath string, out, errOut io.Writer, b []*builders.APIGroupBuilder) *ServerOptions {
	versions := []schema.GroupVersion{}
	for _, b := range b {
		versions = append(versions, b.GetLegacyCodec()...)
	}

	o := &ServerOptions{
		RecommendedOptions: genericoptions.NewRecommendedOptions(etcdPath, builders.Codecs.LegacyCodec(versions...)),
		APIBuilders:        b,
		RunDelegatedAuth:   true,
	}
	o.RecommendedOptions.SecureServing.BindPort = 443

	o.RecommendedOptions.Authorization.RemoteKubeConfigFileOptional = true
	o.RecommendedOptions.Authentication.RemoteKubeConfigFileOptional = true
	o.InsecureServingOptions = func() *genericoptions.DeprecatedInsecureServingOptionsWithLoopback {
		o := genericoptions.DeprecatedInsecureServingOptions{}
		return o.WithLoopback()
	}()

	return o
}

// NewCommandStartMaster provides a CLI handler for 'start master' command
func NewCommandStartServer(etcdPath string, out, errOut io.Writer, builders []*builders.APIGroupBuilder,
	stopCh <-chan struct{}, title, version string, tweakConfigFuncs ...func(apiServer *apiserver.Config) error) (*cobra.Command, *ServerOptions) {
	o := NewServerOptions(etcdPath, out, errOut, builders)

	// Support overrides
	cmd := &cobra.Command{
		Short: "Launch an API server",
		Long:  "Launch an API server",
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.RunServer(stopCh, title, version, tweakConfigFuncs...); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&o.PrintBearerToken, "print-bearer-token", false,
		"Print a curl command with the bearer token to test the server")
	flags.BoolVar(&o.PrintOpenapi, "print-openapi", false,
		"Print the openapi json and exit")
	flags.BoolVar(&o.RunDelegatedAuth, "delegated-auth", true,
		"Setup delegated auth")
	//flags.StringVar(&o.Kubeconfig, "kubeconfig", "", "Kubeconfig of apiserver to talk to.")
	o.RecommendedOptions.AddFlags(flags)
	o.InsecureServingOptions.AddFlags(flags)
	feature.DefaultFeatureGate.AddFlag(flags)

	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)
	flags.AddGoFlagSet(klogFlags)

	// Sync the glog and klog flags.
	klogFlags.VisitAll(func(f *flag.Flag) {
		goFlag := flag.CommandLine.Lookup(f.Name)
		if goFlag != nil {
			goFlag.Value.Set(f.Value.String())
		}
	})

	return cmd, o
}

func (o ServerOptions) Validate(args []string) error {
	return nil
}

func (o *ServerOptions) Complete() error {
	return nil
}

func applyOptions(config *genericapiserver.Config, applyTo ...func(*genericapiserver.Config) error) error {
	for _, fn := range applyTo {
		if err := fn(config); err != nil {
			return err
		}
	}
	return nil
}

func (o ServerOptions) Config(tweakConfigFuncs ...func(config *apiserver.Config) error) (*apiserver.Config, error) {
	// switching pagination according to the feature-gate
	o.RecommendedOptions.Etcd.StorageConfig.Paging = feature.DefaultFeatureGate.Enabled(features.APIListChunking)

	// TODO have a "real" external address
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts(
		"localhost", nil, nil); err != nil {

		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	serverConfig := genericapiserver.NewConfig(builders.Codecs)

	err := applyOptions(
		serverConfig,
		o.RecommendedOptions.Etcd.ApplyTo,
		func(cfg *genericapiserver.Config) error {
			return o.RecommendedOptions.SecureServing.ApplyTo(&cfg.SecureServing, &cfg.LoopbackClientConfig)
		},
		o.RecommendedOptions.Audit.ApplyTo,
		o.RecommendedOptions.Features.ApplyTo,
	)
	if err != nil {
		return nil, err
	}

	//if serverConfig.SharedInformerFactory == nil && len(o.Kubeconfig) > 0 {
	//	path, _ := filepath.Abs(o.Kubeconfig)
	//	klog.Infof("Creating shared informer factory from kubeconfig %s", path)
	//	config, err := clientcmd.BuildConfigFromFlags("", o.Kubeconfig)
	//	clientset, err := kubernetes.NewForConfig(config)
	//	if err != nil {
	//		klog.Errorf("Couldn't create clientset due to %v. SharedInformerFactory will not be set.", err)
	//		return nil, err
	//	}
	//	serverConfig.SharedInformerFactory = informers.NewSharedInformerFactory(clientset, 10*time.Minute)
	//}

	if o.RunDelegatedAuth {
		err := applyOptions(
			serverConfig,
			func(cfg *genericapiserver.Config) error {
				return o.RecommendedOptions.Authentication.ApplyTo(&cfg.Authentication, cfg.SecureServing, cfg.OpenAPIConfig)
			},
			func(cfg *genericapiserver.Config) error {
				return o.RecommendedOptions.Authorization.ApplyTo(&cfg.Authorization)
			},
		)
		if err != nil {
			return nil, err
		}
	}

	var insecureServingInfo *genericapiserver.DeprecatedInsecureServingInfo
	if err := o.InsecureServingOptions.ApplyTo(&insecureServingInfo, &serverConfig.LoopbackClientConfig); err != nil {
		return nil, err
	}

	config := &apiserver.Config{
		GenericConfig:       serverConfig,
		InsecureServingInfo: insecureServingInfo,
	}
	for _, tweakConfigFunc := range tweakConfigFuncs {
		if err := tweakConfigFunc(config); err != nil {
			return nil, err
		}
	}
	return config, nil
}

func (o *ServerOptions) RunServer(stopCh <-chan struct{}, title, version string, tweakConfigFuncs ...func(apiserver *apiserver.Config) error) error {
	config, err := o.Config(tweakConfigFuncs...)
	if err != nil {
		return err
	}

	if o.PrintBearerToken {
		klog.Infof("Serving on loopback...")
		klog.Infof("\n\n********************************\nTo test the server run:\n"+
			"curl -k -H \"Authorization: Bearer %s\" %s\n********************************\n\n",
			config.GenericConfig.LoopbackClientConfig.BearerToken,
			config.GenericConfig.LoopbackClientConfig.Host)
	}
	o.BearerToken = config.GenericConfig.LoopbackClientConfig.BearerToken

	for _, provider := range o.APIBuilders {
		config.AddApi(provider)
	}

	config.Init()

	config.GenericConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(GetOpenApiDefinition, openapinamer.NewDefinitionNamer(builders.Scheme))
	config.GenericConfig.OpenAPIConfig.Info.Title = title
	config.GenericConfig.OpenAPIConfig.Info.Version = version

	genericServer, err := config.Complete().New()
	if err != nil {
		return err
	}

	for _, h := range o.PostStartHooks {
		if err := genericServer.GenericAPIServer.AddPostStartHook(h.Name, h.Fn); err != nil {
			return err
		}
	}

	s := genericServer.GenericAPIServer.PrepareRun()
	err = validators.OpenAPI.SetSchema(readOpenapi(config.GenericConfig.LoopbackClientConfig.BearerToken, genericServer.GenericAPIServer.Handler))
	if o.PrintOpenapi {
		fmt.Printf("%s", validators.OpenAPI.OpenApi)
		os.Exit(0)
	}
	if err != nil {
		return err
	}

	if config.InsecureServingInfo != nil {
		c := config.GenericConfig
		handler := s.GenericAPIServer.UnprotectedHandler()
		handler = genericapifilters.WithAudit(handler, c.AuditBackend, c.AuditPolicyChecker, c.LongRunningFunc)
		handler = genericapifilters.WithAuthentication(handler, server.InsecureSuperuser{}, nil)
		handler = genericfilters.WithCORS(handler, c.CorsAllowedOriginList, nil, nil, nil, "true")
		handler = genericfilters.WithTimeoutForNonLongRunningRequests(handler, c.LongRunningFunc, c.RequestTimeout)
		handler = genericfilters.WithMaxInFlightLimit(handler, c.MaxRequestsInFlight, c.MaxMutatingRequestsInFlight, c.LongRunningFunc)
		handler = genericapifilters.WithRequestInfo(handler, server.NewRequestInfoResolver(c))
		handler = genericfilters.WithPanicRecovery(handler)
		if err := config.InsecureServingInfo.Serve(handler, config.GenericConfig.RequestTimeout, stopCh); err != nil {
			return err
		}
	}

	return s.Run(stopCh)
}

func readOpenapi(bearerToken string, handler *genericapiserver.APIServerHandler) string {
	req, err := http.NewRequest("GET", "/swagger.json", nil)
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", bearerToken))
	if err != nil {
		panic(fmt.Errorf("Could not create openapi request %v", err))
	}
	resp := &BufferedResponse{}
	handler.ServeHTTP(resp, req)
	return resp.String()
}

type BufferedResponse struct {
	bytes.Buffer
}

func (BufferedResponse) Header() http.Header { return http.Header{} }
func (BufferedResponse) WriteHeader(int)     {}
