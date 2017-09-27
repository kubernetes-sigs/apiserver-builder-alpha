package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_test(
    name = "go_default_test",
    srcs = ["log_test.go"],
    importpath = "k8s.io/apiserver/pkg/server/httplog",
    library = ":go_default_library",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "log.go",
    ],
    importpath = "k8s.io/apiserver/pkg/server/httplog",
    tags = ["automanaged"],
    deps = ["//vendor/github.com/golang/glog:go_default_library"],
)
