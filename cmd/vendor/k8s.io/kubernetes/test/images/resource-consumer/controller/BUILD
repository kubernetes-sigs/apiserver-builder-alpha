package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "controller",
    importpath = "k8s.io/kubernetes/test/images/resource-consumer/controller",
    library = ":go_default_library",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["controller.go"],
    importpath = "k8s.io/kubernetes/test/images/resource-consumer/controller",
    tags = ["automanaged"],
    deps = ["//cmd/vendor/k8s.io/kubernetes/test/images/resource-consumer/common:go_default_library"],
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
