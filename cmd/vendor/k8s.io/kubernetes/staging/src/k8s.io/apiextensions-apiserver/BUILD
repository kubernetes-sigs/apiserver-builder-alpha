package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "apiextensions-apiserver",
    importpath = "k8s.io/kubernetes/staging/src/k8s.io/apiextensions-apiserver",
    library = ":go_default_library",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "k8s.io/kubernetes/staging/src/k8s.io/apiextensions-apiserver",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/apiextensions-apiserver/pkg/cmd/server:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        "//cmd/vendor/k8s.io/apiserver/pkg/util/logs:go_default_library",
    ],
)
