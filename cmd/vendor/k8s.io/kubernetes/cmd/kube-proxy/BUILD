package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)
load("//pkg/version:def.bzl", "version_x_defs")

go_binary(
    name = "kube-proxy",
    gc_linkopts = [
        "-linkmode",
        "external",
        "-extldflags",
        "-static",
    ],
    importpath = "k8s.io/kubernetes/cmd/kube-proxy",
    library = ":go_default_library",
    tags = ["automanaged"],
    x_defs = version_x_defs(),
)

go_library(
    name = "go_default_library",
    srcs = ["proxy.go"],
    importpath = "k8s.io/kubernetes/cmd/kube-proxy",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/github.com/spf13/pflag:go_default_library",
        "//cmd/vendor/k8s.io/apiserver/pkg/util/flag:go_default_library",
        "//cmd/vendor/k8s.io/apiserver/pkg/util/logs:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/cmd/kube-proxy/app:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/client/metrics/prometheus:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//cmd/kube-proxy/app:all-srcs",
    ],
    tags = ["automanaged"],
)
