package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["storage_apps.go"],
    importpath = "k8s.io/kubernetes/pkg/registry/apps/rest",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/api/apps/v1beta1:go_default_library",
        "//cmd/vendor/k8s.io/apiserver/pkg/registry/generic:go_default_library",
        "//cmd/vendor/k8s.io/apiserver/pkg/registry/rest:go_default_library",
        "//cmd/vendor/k8s.io/apiserver/pkg/server:go_default_library",
        "//cmd/vendor/k8s.io/apiserver/pkg/server/storage:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/apis/apps:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/registry/apps/controllerrevision/storage:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/registry/apps/statefulset/storage:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/registry/extensions/deployment/storage:go_default_library",
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
