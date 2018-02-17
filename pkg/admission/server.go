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

package admission

import (
	"flag"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"net/http"
)

var (
	cert = flag.String("tls-cert-file", "", ""+
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert).")
	key = flag.String("tls-private-key-file", "", ""+
		"File containing the default x509 private key matching --tls-cert-file.")

	Scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(Scheme)

	// Singleton
	Singleton = AdmissionServer{
		Handlers: map[string]handler{},
	}
)

type handler func(w http.ResponseWriter, r *http.Request)

type AdmissionServer struct {
	Handlers map[string]handler
}

func (as *AdmissionServer) Run() {
	server := &http.Server{
		Addr: ":443",
	}
	for add, fn := range as.Handlers {
		http.HandleFunc(add, fn)
	}
	server.ListenAndServeTLS(*cert, *key)
}
