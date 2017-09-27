package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["interface.go"],
    importpath = "k8s.io/kubernetes/pkg/client/informers/informers_generated/internalversion/networking",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/kubernetes/pkg/client/informers/informers_generated/internalversion/internalinterfaces:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/client/informers/informers_generated/internalversion/networking/internalversion:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//pkg/client/informers/informers_generated/internalversion/networking/internalversion:all-srcs",
    ],
    tags = ["automanaged"],
)
