/*
Copyright 2019 The Kubernetes Authors.

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

package main

import (
	"os"
	"strconv"

	"k8s.io/klog/v2"
	mysqlv1 "sigs.k8s.io/apiserver-builder-alpha/example/kine/pkg/apis/mysql/v1"
	"sigs.k8s.io/apiserver-runtime/pkg/builder"
	"sigs.k8s.io/apiserver-runtime/pkg/experimental/storage/mysql"
)

func main() {
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort, _ := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	mysqlUser := os.Getenv("MYSQL_USERNAME")
	mysqlPasswd := os.Getenv("MYSQL_PASSWORD")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	err := builder.APIServer.
		WithResourceAndStorage(&mysqlv1.Tiger{}, mysql.NewMysqlStorageProvider(
			mysqlHost,
			int32(mysqlPort),
			mysqlUser,
			mysqlPasswd,
			mysqlDatabase,
		)). // namespaced resource
		WithoutEtcd().
		WithLocalDebugExtension().
		Execute()
	if err != nil {
		klog.Fatal(err)
	}
}
