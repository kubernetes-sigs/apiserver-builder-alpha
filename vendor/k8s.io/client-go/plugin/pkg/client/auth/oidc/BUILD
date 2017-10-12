package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_test(
    name = "go_default_test",
    srcs = ["oidc_test.go"],
    importpath = "k8s.io/client-go/plugin/pkg/client/auth/oidc",
    library = ":go_default_library",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["oidc.go"],
    importpath = "k8s.io/client-go/plugin/pkg/client/auth/oidc",
    tags = ["automanaged"],
    deps = [
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/golang.org/x/oauth2:go_default_library",
        "//vendor/k8s.io/client-go/rest:go_default_library",
    ],
)
