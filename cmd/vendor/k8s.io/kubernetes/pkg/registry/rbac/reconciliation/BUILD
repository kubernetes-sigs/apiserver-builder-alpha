package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_test(
    name = "go_default_test",
    srcs = [
        "reconcile_role_test.go",
        "reconcile_rolebindings_test.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/registry/rbac/reconciliation",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api/helper:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/apis/rbac:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "clusterrole_interfaces.go",
        "clusterrolebinding_interfaces.go",
        "reconcile_role.go",
        "reconcile_rolebindings.go",
        "role_interfaces.go",
        "rolebinding_interfaces.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/registry/rbac/reconciliation",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/apimachinery/pkg/api/errors:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/apis/rbac:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/typed/core/internalversion:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/typed/rbac/internalversion:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/registry/rbac/validation:go_default_library",
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
