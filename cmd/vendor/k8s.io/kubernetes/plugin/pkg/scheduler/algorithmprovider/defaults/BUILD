package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["defaults.go"],
    importpath = "k8s.io/kubernetes/plugin/pkg/scheduler/algorithmprovider/defaults",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/github.com/golang/glog:go_default_library",
        "//cmd/vendor/k8s.io/api/core/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/sets:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/cloudprovider/providers/aws:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/plugin/pkg/scheduler/algorithm:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/plugin/pkg/scheduler/algorithm/predicates:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/plugin/pkg/scheduler/algorithm/priorities:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/plugin/pkg/scheduler/core:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/plugin/pkg/scheduler/factory:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "compatibility_test.go",
        "defaults_test.go",
    ],
    importpath = "k8s.io/kubernetes/plugin/pkg/scheduler/algorithmprovider/defaults",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/api/core/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/sets:go_default_library",
        "//cmd/vendor/k8s.io/client-go/informers:go_default_library",
        "//cmd/vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//cmd/vendor/k8s.io/client-go/rest:go_default_library",
        "//cmd/vendor/k8s.io/client-go/util/testing:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api/install:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/plugin/pkg/scheduler/api:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/plugin/pkg/scheduler/api/latest:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/plugin/pkg/scheduler/factory:go_default_library",
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
