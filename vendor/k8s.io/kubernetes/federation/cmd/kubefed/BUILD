package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "kubefed",
    importpath = "k8s.io/kubernetes/federation/cmd/kubefed",
    library = ":go_default_library",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["kubefed.go"],
    importpath = "k8s.io/kubernetes/federation/cmd/kubefed",
    tags = ["automanaged"],
    deps = ["//vendor/k8s.io/kubernetes/federation/cmd/kubefed/app:go_default_library"],
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
        "//federation/cmd/kubefed/app:all-srcs",
    ],
    tags = ["automanaged"],
)
