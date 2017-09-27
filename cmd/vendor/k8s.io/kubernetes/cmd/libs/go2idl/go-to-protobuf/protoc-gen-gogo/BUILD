package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "protoc-gen-gogo",
    importpath = "k8s.io/kubernetes/cmd/libs/go2idl/go-to-protobuf/protoc-gen-gogo",
    library = ":go_default_library",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "k8s.io/kubernetes/cmd/libs/go2idl/go-to-protobuf/protoc-gen-gogo",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/github.com/gogo/protobuf/gogoproto:go_default_library",
        "//cmd/vendor/github.com/gogo/protobuf/proto:go_default_library",
        "//cmd/vendor/github.com/gogo/protobuf/sortkeys:go_default_library",
        "//cmd/vendor/github.com/gogo/protobuf/vanity/command:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
)
