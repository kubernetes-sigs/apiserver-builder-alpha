package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "genyaml",
    importpath = "k8s.io/kubernetes/cmd/genyaml",
    library = ":go_default_library",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["gen_kubectl_yaml.go"],
    importpath = "k8s.io/kubernetes/cmd/genyaml",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/github.com/spf13/cobra:go_default_library",
        "//cmd/vendor/github.com/spf13/pflag:go_default_library",
        "//vendor/gopkg.in/yaml.v2:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/cmd/genutils:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubectl/cmd:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubectl/cmd/util:go_default_library",
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
