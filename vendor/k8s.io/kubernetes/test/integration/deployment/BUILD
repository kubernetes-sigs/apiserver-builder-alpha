package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_test(
    name = "go_default_test",
    size = "large",
    srcs = [
        "deployment_test.go",
        "main_test.go",
    ],
    importpath = "k8s.io/kubernetes/test/integration/deployment",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/controller/deployment/util:go_default_library",
        "//vendor/k8s.io/kubernetes/test/integration/framework:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["util.go"],
    importpath = "k8s.io/kubernetes/test/integration/deployment",
    tags = ["automanaged"],
    deps = [
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/api/extensions/v1beta1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        "//vendor/k8s.io/client-go/informers:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//vendor/k8s.io/client-go/rest:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/api/v1/pod:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/controller/deployment:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/controller/replicaset:go_default_library",
        "//vendor/k8s.io/kubernetes/test/integration/framework:go_default_library",
        "//vendor/k8s.io/kubernetes/test/utils:go_default_library",
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
