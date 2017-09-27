package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "custom_resource_definition.go",
        "etcd_failure.go",
        "framework.go",
        "garbage_collector.go",
        "initializers.go",
        "namespace.go",
        "table_conversion.go",
    ],
    importpath = "k8s.io/kubernetes/test/e2e/apimachinery",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/github.com/onsi/ginkgo:go_default_library",
        "//cmd/vendor/github.com/onsi/gomega:go_default_library",
        "//cmd/vendor/k8s.io/api/admissionregistration/v1alpha1:go_default_library",
        "//cmd/vendor/k8s.io/api/core/v1:go_default_library",
        "//cmd/vendor/k8s.io/api/extensions/v1beta1:go_default_library",
        "//cmd/vendor/k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1:go_default_library",
        "//cmd/vendor/k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset:go_default_library",
        "//cmd/vendor/k8s.io/apiextensions-apiserver/test/integration/testserver:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/api/errors:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1alpha1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/labels:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/intstr:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        "//cmd/vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api/v1/pod:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/client/retry:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/controller:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/printers:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/test/e2e/framework:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/test/e2e/metrics:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/test/e2e/workload:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/test/utils:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
)
