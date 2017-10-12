package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["json.go"],
    importpath = "k8s.io/apimachinery/pkg/util/json",
    tags = ["automanaged"],
)

go_test(
    name = "go_default_test",
    srcs = ["json_test.go"],
    importpath = "k8s.io/apimachinery/pkg/util/json",
    library = ":go_default_library",
    tags = ["automanaged"],
)
