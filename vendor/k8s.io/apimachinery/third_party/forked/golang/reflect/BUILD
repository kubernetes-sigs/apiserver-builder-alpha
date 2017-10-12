package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_test(
    name = "go_default_test",
    srcs = ["deep_equal_test.go"],
    importpath = "k8s.io/apimachinery/third_party/forked/golang/reflect",
    library = ":go_default_library",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "deep_equal.go",
        "type.go",
    ],
    importpath = "k8s.io/apimachinery/third_party/forked/golang/reflect",
    tags = ["automanaged"],
)
