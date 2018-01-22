package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_test(
    name = "go_default_test",
    srcs = ["clusterroles_test.go"],
    importpath = "k8s.io/kubernetes/cmd/kubeadm/app/phases/apiconfig",
    library = ":go_default_library",
)

go_library(
    name = "go_default_library",
    srcs = ["clusterroles.go"],
    importpath = "k8s.io/kubernetes/cmd/kubeadm/app/phases/apiconfig",
    deps = [
        "//vendor/k8s.io/api/rbac/v1beta1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/api/errors:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//vendor/k8s.io/kubernetes/cmd/kubeadm/app/constants:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/version:go_default_library",
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
