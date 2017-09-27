licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["kubectl.go"],
    importpath = "k8s.io/kubernetes/cmd/kubectl/app",
    tags = ["automanaged"],
    visibility = [
        "//build/visible_to:cmd_kubectl_app_CONSUMERS",
    ],
    deps = [
        "//cmd/vendor/k8s.io/client-go/plugin/pkg/client/auth:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/client/metrics/prometheus:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubectl/cmd:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubectl/cmd/util:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/util/logs:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/version/prometheus:go_default_library",
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
        "//build/visible_to:cmd_kubectl_app_CONSUMERS",
    ],
)
