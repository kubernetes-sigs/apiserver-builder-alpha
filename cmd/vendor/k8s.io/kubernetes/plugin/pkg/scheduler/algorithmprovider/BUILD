package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["plugins.go"],
    importpath = "k8s.io/kubernetes/plugin/pkg/scheduler/algorithmprovider",
    tags = ["automanaged"],
    deps = ["//cmd/vendor/k8s.io/kubernetes/plugin/pkg/scheduler/algorithmprovider/defaults:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = ["plugins_test.go"],
    importpath = "k8s.io/kubernetes/plugin/pkg/scheduler/algorithmprovider",
    library = ":go_default_library",
    tags = [
        "automanaged",
    ],
    deps = ["//cmd/vendor/k8s.io/kubernetes/plugin/pkg/scheduler/factory:go_default_library"],
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
        "//plugin/pkg/scheduler/algorithmprovider/defaults:all-srcs",
    ],
    tags = ["automanaged"],
)
