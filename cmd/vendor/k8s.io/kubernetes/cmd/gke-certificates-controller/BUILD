package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "k8s.io/kubernetes/cmd/gke-certificates-controller",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/github.com/spf13/pflag:go_default_library",
        "//cmd/vendor/k8s.io/apiserver/pkg/util/flag:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/cmd/gke-certificates-controller/app:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/util/logs:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/version/verflag:go_default_library",
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
        "//cmd/gke-certificates-controller/app:all-srcs",
    ],
    tags = ["automanaged"],
)

go_binary(
    name = "gke-certificates-controller",
    importpath = "k8s.io/kubernetes/cmd/gke-certificates-controller",
    library = ":go_default_library",
    tags = ["automanaged"],
)
