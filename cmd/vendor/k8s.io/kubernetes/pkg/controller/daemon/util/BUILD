package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["daemonset_util.go"],
    importpath = "k8s.io/kubernetes/pkg/controller/daemon/util",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/api/core/v1:go_default_library",
        "//cmd/vendor/k8s.io/api/extensions/v1beta1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//cmd/vendor/k8s.io/client-go/kubernetes/scheme:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api/v1/helper:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api/v1/pod:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/util/labels:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/plugin/pkg/scheduler/algorithm:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = ["daemonset_util_test.go"],
    importpath = "k8s.io/kubernetes/pkg/controller/daemon/util",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/api/core/v1:go_default_library",
        "//cmd/vendor/k8s.io/api/extensions/v1beta1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api/testapi:go_default_library",
    ],
)
