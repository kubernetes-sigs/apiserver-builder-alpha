package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["hairpin.go"],
    importpath = "k8s.io/kubernetes/pkg/kubelet/network/hairpin",
    tags = ["automanaged"],
    deps = [
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/exec:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["hairpin_test.go"],
    importpath = "k8s.io/kubernetes/pkg/kubelet/network/hairpin",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = ["//vendor/k8s.io/kubernetes/pkg/util/exec:go_default_library"],
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
