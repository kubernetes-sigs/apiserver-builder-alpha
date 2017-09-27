package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["goroutinemap.go"],
    importpath = "k8s.io/kubernetes/pkg/util/goroutinemap",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/github.com/golang/glog:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/runtime:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/util/goroutinemap/exponentialbackoff:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["goroutinemap_test.go"],
    importpath = "k8s.io/kubernetes/pkg/util/goroutinemap",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = ["//cmd/vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library"],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//pkg/util/goroutinemap/exponentialbackoff:all-srcs",
    ],
    tags = ["automanaged"],
)
