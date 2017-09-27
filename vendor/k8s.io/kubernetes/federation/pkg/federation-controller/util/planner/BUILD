package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["planner.go"],
    importpath = "k8s.io/kubernetes/federation/pkg/federation-controller/util/planner",
    tags = ["automanaged"],
    deps = ["//vendor/k8s.io/kubernetes/federation/apis/federation:go_default_library"],
)

go_test(
    name = "go_default_test",
    srcs = ["planner_test.go"],
    importpath = "k8s.io/kubernetes/federation/pkg/federation-controller/util/planner",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
        "//vendor/k8s.io/kubernetes/federation/apis/federation:go_default_library",
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
