package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_test(
    name = "go_default_test",
    srcs = ["storage_test.go"],
    importpath = "k8s.io/kubernetes/pkg/registry/apps/controllerrevision/storage",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/fields:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/labels:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//cmd/vendor/k8s.io/apiserver/pkg/registry/generic:go_default_library",
        "//cmd/vendor/k8s.io/apiserver/pkg/storage/etcd/testing:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/apis/apps:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/registry/registrytest:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["storage.go"],
    importpath = "k8s.io/kubernetes/pkg/registry/apps/controllerrevision/storage",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//cmd/vendor/k8s.io/apiserver/pkg/registry/generic:go_default_library",
        "//cmd/vendor/k8s.io/apiserver/pkg/registry/generic/registry:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/apis/apps:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/registry/apps/controllerrevision:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/registry/cachesize:go_default_library",
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
