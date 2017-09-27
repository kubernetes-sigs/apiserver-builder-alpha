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
        "conntrack.go",
        "server.go",
    ],
    importpath = "k8s.io/kubernetes/cmd/kube-proxy/app",
    tags = ["automanaged"],
    deps = [
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/github.com/prometheus/client_golang/prometheus:go_default_library",
        "//vendor/github.com/spf13/cobra:go_default_library",
        "//vendor/github.com/spf13/pflag:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/serializer:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/serializer/json:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/net:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        "//vendor/k8s.io/apiserver/pkg/server/healthz:go_default_library",
        "//vendor/k8s.io/apiserver/pkg/util/feature:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes/typed/core/v1:go_default_library",
        "//vendor/k8s.io/client-go/tools/clientcmd:go_default_library",
        "//vendor/k8s.io/client-go/tools/clientcmd/api:go_default_library",
        "//vendor/k8s.io/client-go/tools/record:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/api:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/apis/componentconfig:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/apis/componentconfig/v1alpha1:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/client/informers/informers_generated/internalversion:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/kubectl/cmd/util:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/kubelet/qos:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/master/ports:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/proxy:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/proxy/config:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/proxy/healthcheck:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/proxy/iptables:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/proxy/userspace:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/proxy/winuserspace:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/configz:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/dbus:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/exec:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/iptables:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/mount:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/netsh:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/node:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/oom:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/pointer:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/resourcecontainer:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/sysctl:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/version/verflag:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["server_test.go"],
    importpath = "k8s.io/kubernetes/cmd/kube-proxy/app",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/diff:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/api:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/apis/componentconfig:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/apis/componentconfig/v1alpha1:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/configz:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/iptables:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/pointer:go_default_library",
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
