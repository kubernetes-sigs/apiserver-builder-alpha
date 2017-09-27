package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["install.go"],
    importpath = "k8s.io/kubernetes/plugin/pkg/admission/resourcequota/apis/resourcequota/install",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/apimachinery/pkg/apimachinery/announced:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/apimachinery/registered:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/plugin/pkg/admission/resourcequota/apis/resourcequota:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/plugin/pkg/admission/resourcequota/apis/resourcequota/v1alpha1:go_default_library",
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
