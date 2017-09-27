package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["checks.go"],
    importpath = "k8s.io/kubernetes/cmd/kube-apiserver/app/preflight",
    tags = ["automanaged"],
)

go_test(
    name = "go_default_test",
    srcs = ["checks_test.go"],
    importpath = "k8s.io/kubernetes/cmd/kube-apiserver/app/preflight",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = ["//cmd/vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library"],
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
