package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["desired_state_of_world_populator.go"],
    importpath = "k8s.io/kubernetes/pkg/controller/volume/attachdetach/populator",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/github.com/golang/glog:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/api/errors:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/labels:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/runtime:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        "//cmd/vendor/k8s.io/client-go/listers/core/v1:go_default_library",
        "//cmd/vendor/k8s.io/client-go/tools/cache:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/controller/volume/attachdetach/cache:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/controller/volume/attachdetach/util:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/volume:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/volume/util/volumehelper:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = ["desired_state_of_world_populator_test.go"],
    importpath = "k8s.io/kubernetes/pkg/controller/volume/attachdetach/populator",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/api/core/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//cmd/vendor/k8s.io/client-go/informers:go_default_library",
        "//cmd/vendor/k8s.io/client-go/kubernetes/fake:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/controller:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/controller/volume/attachdetach/cache:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/volume/testing:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/volume/util/volumehelper:go_default_library",
    ],
)
