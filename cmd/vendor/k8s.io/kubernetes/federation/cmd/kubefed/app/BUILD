package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["kubefed.go"],
    importpath = "k8s.io/kubernetes/federation/cmd/kubefed/app",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/kubernetes/federation/pkg/kubefed:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/client/metrics/prometheus:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubectl/cmd/util:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/util/logs:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/version:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/version/prometheus:go_default_library",
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
