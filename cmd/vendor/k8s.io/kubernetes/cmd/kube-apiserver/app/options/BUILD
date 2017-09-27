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
        "options.go",
        "validation.go",
    ],
    importpath = "k8s.io/kubernetes/cmd/kube-apiserver/app/options",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/github.com/spf13/pflag:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/net:go_default_library",
        "//cmd/vendor/k8s.io/apiserver/pkg/server/options:go_default_library",
        "//cmd/vendor/k8s.io/apiserver/pkg/storage/storagebackend:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api/validation:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/features:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubeapiserver/options:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/client:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/master/ports:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["options_test.go"],
    importpath = "k8s.io/kubernetes/cmd/kube-apiserver/app/options",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = ["//cmd/vendor/github.com/spf13/pflag:go_default_library"],
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
