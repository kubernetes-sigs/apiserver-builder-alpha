licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "auth.go",
        "cani.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/kubectl/cmd/auth",
    tags = ["automanaged"],
    visibility = [
        "//build/visible_to:pkg_kubectl_cmd_auth_CONSUMERS",
    ],
    deps = [
        "//cmd/vendor/github.com/spf13/cobra:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/api/meta:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/apis/authorization:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/typed/authorization/internalversion:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubectl/cmd/templates:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubectl/cmd/util:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = [
        "//build/visible_to:pkg_kubectl_cmd_auth_CONSUMERS",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["cani_test.go"],
    importpath = "k8s.io/kubernetes/pkg/kubectl/cmd/auth",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/client-go/rest:go_default_library",
        "//cmd/vendor/k8s.io/client-go/rest/fake:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubectl/cmd/testing:go_default_library",
    ],
)
