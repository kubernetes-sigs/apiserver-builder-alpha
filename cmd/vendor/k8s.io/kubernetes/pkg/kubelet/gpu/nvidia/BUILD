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
        "helpers.go",
        "nvidia_gpu_manager.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/kubelet/gpu/nvidia",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/github.com/golang/glog:go_default_library",
        "//cmd/vendor/k8s.io/api/core/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/api/resource:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/sets:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/dockershim/libdocker:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/gpu:go_default_library",
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
    srcs = ["nvidia_gpu_manager_test.go"],
    importpath = "k8s.io/kubernetes/pkg/kubelet/gpu/nvidia",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
        "//cmd/vendor/k8s.io/api/core/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/api/resource:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/sets:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/uuid:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/dockershim/libdocker:go_default_library",
    ],
)
