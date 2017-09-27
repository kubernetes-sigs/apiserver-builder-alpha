package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_test(
    name = "go_default_test",
    srcs = ["group_version_test.go"],
    importpath = "k8s.io/apimachinery/pkg/runtime/schema",
    library = ":go_default_library",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "generated.pb.go",
        "group_version.go",
        "interfaces.go",
    ],
    importpath = "k8s.io/apimachinery/pkg/runtime/schema",
    tags = ["automanaged"],
    deps = ["//vendor/github.com/gogo/protobuf/proto:go_default_library"],
)

filegroup(
    name = "go_default_library_protos",
    srcs = ["generated.proto"],
    visibility = ["//visibility:public"],
)
