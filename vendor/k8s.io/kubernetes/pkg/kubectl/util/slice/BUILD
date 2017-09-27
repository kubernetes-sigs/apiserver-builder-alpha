package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["slice.go"],
    importpath = "k8s.io/kubernetes/pkg/kubectl/util/slice",
    tags = ["automanaged"],
)

go_test(
    name = "go_default_test",
    srcs = ["slice_test.go"],
    importpath = "k8s.io/kubernetes/pkg/kubectl/util/slice",
    library = ":go_default_library",
    tags = ["automanaged"],
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
