package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["spdy.go"],
    importpath = "k8s.io/kubernetes/staging/src/k8s.io/client-go/transport/spdy",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/httpstream:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/httpstream/spdy:go_default_library",
        "//cmd/vendor/k8s.io/client-go/rest:go_default_library",
    ],
)
