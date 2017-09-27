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
        "admission_failure_handler_stub.go",
        "doc.go",
        "fake_handler_runner.go",
        "handlers.go",
        "interfaces.go",
        "predicate.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/kubelet/lifecycle",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/github.com/golang/glog:go_default_library",
        "//cmd/vendor/k8s.io/api/core/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/intstr:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/container:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/types:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/util/format:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/security/apparmor:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/plugin/pkg/scheduler/algorithm:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/plugin/pkg/scheduler/algorithm/predicates:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/plugin/pkg/scheduler/schedulercache:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["handlers_test.go"],
    importpath = "k8s.io/kubernetes/pkg/kubelet/lifecycle",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/api/core/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/intstr:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/container:go_default_library",
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
