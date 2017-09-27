package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_test(
    name = "go_default_test",
    srcs = ["runtime_test.go"],
    importpath = "k8s.io/apimachinery/pkg/util/runtime",
    library = ":go_default_library",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["runtime.go"],
    importpath = "k8s.io/apimachinery/pkg/util/runtime",
    tags = ["automanaged"],
    deps = ["//vendor/github.com/golang/glog:go_default_library"],
)
