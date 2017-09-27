package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["interface.go"],
    importpath = "k8s.io/client-go/informers/apps",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/client-go/informers/apps/v1beta1:go_default_library",
        "//cmd/vendor/k8s.io/client-go/informers/internalinterfaces:go_default_library",
    ],
)
