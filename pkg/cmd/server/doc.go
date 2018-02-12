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

package server

import (
	_ "github.com/go-openapi/analysis"
	_ "github.com/go-openapi/errors"
	_ "github.com/go-openapi/jsonpointer"
	_ "github.com/go-openapi/jsonreference"
	_ "github.com/go-openapi/loads"
	_ "github.com/go-openapi/runtime"
	_ "github.com/go-openapi/spec"
	_ "github.com/go-openapi/strfmt"
	_ "github.com/go-openapi/swag"
	_ "github.com/go-openapi/validate"
	_ "github.com/mailru/easyjson/jwriter"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)
