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
        "loadbalancer.go",
        "port_allocator.go",
        "proxier.go",
        "proxysocket.go",
        "rlimit.go",
        "roundrobin.go",
        "udp_server.go",
    ] + select({
        "@io_bazel_rules_go//go/platform:windows_amd64": [
            "rlimit_windows.go",
        ],
        "//conditions:default": [],
    }),
    importpath = "k8s.io/kubernetes/pkg/proxy/userspace",
    tags = ["automanaged"],
    deps = [
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/errors:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/net:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/sets:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/api:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/api/helper:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/proxy:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/proxy/util:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/exec:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/iptables:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/slice:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "port_allocator_test.go",
        "proxier_test.go",
        "roundrobin_test.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/proxy/userspace",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/net:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/runtime:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/api:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/proxy:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/exec:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/iptables/testing:go_default_library",
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
