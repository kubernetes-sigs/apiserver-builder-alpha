package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "kubemark",
    importpath = "k8s.io/kubernetes/cmd/kubemark",
    library = ":go_default_library",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["hollow-node.go"],
    importpath = "k8s.io/kubernetes/cmd/kubemark",
    tags = ["automanaged"],
    deps = [
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/github.com/spf13/pflag:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/sets:go_default_library",
        "//vendor/k8s.io/apiserver/pkg/util/flag:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//vendor/k8s.io/client-go/rest:go_default_library",
        "//vendor/k8s.io/client-go/tools/clientcmd:go_default_library",
        "//vendor/k8s.io/client-go/tools/record:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/api:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/client/metrics/prometheus:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/kubelet/cadvisor/testing:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/kubelet/cm:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/kubelet/dockershim/libdocker:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/kubemark:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/exec:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/iptables/testing:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/sysctl/testing:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/version/prometheus:go_default_library",
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
