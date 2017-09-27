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
        "byte.go",
        "doc.go",
        "empty.go",
        "int.go",
        "int64.go",
        "string.go",
    ],
    importpath = "k8s.io/apimachinery/pkg/util/sets",
    tags = ["automanaged"],
)

go_test(
    name = "go_default_test",
    srcs = ["set_test.go"],
    importpath = "k8s.io/apimachinery/pkg/util/sets",
    library = ":go_default_library",
    tags = ["automanaged"],
)
