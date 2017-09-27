package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "clientset_generated.go",
        "doc.go",
        "register.go",
    ],
    importpath = "k8s.io/kubernetes/staging/src/k8s.io/kube-aggregator/pkg/client/clientset_generated/internalclientset/fake",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/runtime/serializer:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/watch:go_default_library",
        "//cmd/vendor/k8s.io/client-go/discovery:go_default_library",
        "//cmd/vendor/k8s.io/client-go/discovery/fake:go_default_library",
        "//cmd/vendor/k8s.io/client-go/testing:go_default_library",
        "//cmd/vendor/k8s.io/kube-aggregator/pkg/apis/apiregistration:go_default_library",
        "//cmd/vendor/k8s.io/kube-aggregator/pkg/client/clientset_generated/internalclientset:go_default_library",
        "//cmd/vendor/k8s.io/kube-aggregator/pkg/client/clientset_generated/internalclientset/typed/apiregistration/internalversion:go_default_library",
        "//cmd/vendor/k8s.io/kube-aggregator/pkg/client/clientset_generated/internalclientset/typed/apiregistration/internalversion/fake:go_default_library",
    ],
)
