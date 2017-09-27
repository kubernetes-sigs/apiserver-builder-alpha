package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "namespace.go",
        "runtime.go",
        "whitelist.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/kubelet/sysctl",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/api/core/v1:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api/v1/helper:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api/validation:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/apis/extensions/validation:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/container:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/lifecycle:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "namespace_test.go",
        "whitelist_test.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/kubelet/sysctl",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = ["//cmd/vendor/k8s.io/api/core/v1:go_default_library"],
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
