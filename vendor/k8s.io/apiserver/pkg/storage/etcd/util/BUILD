package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_test(
    name = "go_default_test",
    srcs = ["etcd_util_test.go"],
    importpath = "k8s.io/apiserver/pkg/storage/etcd/util",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//vendor/github.com/coreos/etcd/client:go_default_library",
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "etcd_util.go",
    ],
    importpath = "k8s.io/apiserver/pkg/storage/etcd/util",
    tags = ["automanaged"],
    deps = ["//vendor/github.com/coreos/etcd/client:go_default_library"],
)
