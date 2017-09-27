licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)
load("//pkg/version:def.bzl", "version_x_defs")

go_binary(
    name = "kubectl",
    gc_linkopts = [
        "-linkmode",
        "external",
        "-extldflags",
        "-static",
    ],
    importpath = "k8s.io/kubernetes/cmd/kubectl",
    library = ":go_default_library",
    tags = ["automanaged"],
    visibility = [
        "//build/visible_to:COMMON_release",
    ],
    x_defs = version_x_defs(),
)

go_library(
    name = "go_default_library",
    srcs = ["kubectl.go"],
    importpath = "k8s.io/kubernetes/cmd/kubectl",
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
    deps = ["//cmd/vendor/k8s.io/kubernetes/cmd/kubectl/app:go_default_library"],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//cmd/kubectl/app:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = [
        "//build/visible_to:cmd_kubectl_CONSUMERS",
    ],
)
