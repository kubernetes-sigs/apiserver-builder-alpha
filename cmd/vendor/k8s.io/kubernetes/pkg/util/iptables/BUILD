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
        "doc.go",
        "iptables.go",
        "iptables_unsupported.go",
        "save_restore.go",
    ] + select({
        "@io_bazel_rules_go//go/platform:linux_amd64": [
            "iptables_linux.go",
        ],
        "//conditions:default": [],
    }),
    importpath = "k8s.io/kubernetes/pkg/util/iptables",
    tags = ["automanaged"],
    deps = [
        "//vendor/github.com/godbus/dbus:go_default_library",
        "//cmd/vendor/github.com/golang/glog:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/sets:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/util/dbus:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/util/exec:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/util/version:go_default_library",
    ] + select({
        "@io_bazel_rules_go//go/platform:linux_amd64": [
            "//vendor/golang.org/x/sys/unix:go_default_library",
            "//cmd/vendor/k8s.io/apimachinery/pkg/util/errors:go_default_library",
            "//cmd/vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        ],
        "//conditions:default": [],
    }),
)

go_test(
    name = "go_default_test",
    srcs = ["iptables_test.go"],
    importpath = "k8s.io/kubernetes/pkg/util/iptables",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/sets:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/util/dbus:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/util/exec:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//pkg/util/iptables/testing:all-srcs",
    ],
    tags = ["automanaged"],
)
