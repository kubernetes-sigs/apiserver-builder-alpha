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
        "errors_test.go",
        "path_test.go",
    ],
    importpath = "k8s.io/apimachinery/pkg/util/validation/field",
    library = ":go_default_library",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "errors.go",
        "path.go",
    ],
    importpath = "k8s.io/apimachinery/pkg/util/validation/field",
    tags = ["automanaged"],
    deps = [
        "//vendor/k8s.io/apimachinery/pkg/util/errors:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/sets:go_default_library",
    ],
)
