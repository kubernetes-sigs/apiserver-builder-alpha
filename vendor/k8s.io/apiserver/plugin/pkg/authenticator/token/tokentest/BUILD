package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["tokentest.go"],
    importpath = "k8s.io/apiserver/plugin/pkg/authenticator/token/tokentest",
    tags = ["automanaged"],
    deps = ["//vendor/k8s.io/apiserver/pkg/authentication/user:go_default_library"],
)
