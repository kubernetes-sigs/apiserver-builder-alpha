package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["clientauth.go"],
    importpath = "k8s.io/client-go/tools/auth",
    tags = ["automanaged"],
    deps = ["//vendor/k8s.io/client-go/rest:go_default_library"],
)

go_test(
    name = "go_default_xtest",
    srcs = ["clientauth_test.go"],
    importpath = "k8s.io/client-go/tools/auth_test",
    tags = ["automanaged"],
    deps = [":go_default_library"],
)
