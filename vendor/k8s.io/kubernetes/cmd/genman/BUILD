package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "genman",
    importpath = "k8s.io/kubernetes/cmd/genman",
    library = ":go_default_library",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["gen_kube_man.go"],
    importpath = "k8s.io/kubernetes/cmd/genman",
    tags = ["automanaged"],
    deps = [
        "//vendor/github.com/cpuguy83/go-md2man/md2man:go_default_library",
        "//vendor/github.com/spf13/cobra:go_default_library",
        "//vendor/github.com/spf13/pflag:go_default_library",
        "//vendor/k8s.io/kubernetes/cmd/cloud-controller-manager/app:go_default_library",
        "//vendor/k8s.io/kubernetes/cmd/genutils:go_default_library",
        "//vendor/k8s.io/kubernetes/cmd/kube-apiserver/app:go_default_library",
        "//vendor/k8s.io/kubernetes/cmd/kube-controller-manager/app:go_default_library",
        "//vendor/k8s.io/kubernetes/cmd/kube-proxy/app:go_default_library",
        "//vendor/k8s.io/kubernetes/cmd/kubelet/app:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/kubectl/cmd:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/kubectl/cmd/util:go_default_library",
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
