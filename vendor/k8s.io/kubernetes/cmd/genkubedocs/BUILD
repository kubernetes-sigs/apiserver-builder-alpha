package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "genkubedocs",
    importpath = "k8s.io/kubernetes/cmd/genkubedocs",
    library = ":go_default_library",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["gen_kube_docs.go"],
    importpath = "k8s.io/kubernetes/cmd/genkubedocs",
    tags = ["automanaged"],
    deps = [
        "//vendor/github.com/spf13/cobra/doc:go_default_library",
        "//vendor/k8s.io/kubernetes/cmd/cloud-controller-manager/app:go_default_library",
        "//vendor/k8s.io/kubernetes/cmd/genutils:go_default_library",
        "//vendor/k8s.io/kubernetes/cmd/kube-apiserver/app:go_default_library",
        "//vendor/k8s.io/kubernetes/cmd/kube-controller-manager/app:go_default_library",
        "//vendor/k8s.io/kubernetes/cmd/kube-proxy/app:go_default_library",
        "//vendor/k8s.io/kubernetes/cmd/kubelet/app:go_default_library",
        "//vendor/k8s.io/kubernetes/plugin/cmd/kube-scheduler/app:go_default_library",
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
