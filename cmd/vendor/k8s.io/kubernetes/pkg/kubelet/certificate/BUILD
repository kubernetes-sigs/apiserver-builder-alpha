package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "certificate_manager.go",
        "certificate_store.go",
        "kubelet.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/kubelet/certificate",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/github.com/golang/glog:go_default_library",
        "//cmd/vendor/k8s.io/api/certificates/v1beta1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/fields:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/watch:go_default_library",
        "//cmd/vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//cmd/vendor/k8s.io/client-go/kubernetes/typed/certificates/v1beta1:go_default_library",
        "//cmd/vendor/k8s.io/client-go/util/cert:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/apis/componentconfig:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/util/file:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "certificate_manager_test.go",
        "certificate_store_test.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/kubelet/certificate",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/api/certificates/v1beta1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/watch:go_default_library",
        "//cmd/vendor/k8s.io/client-go/kubernetes/typed/certificates/v1beta1:go_default_library",
        "//cmd/vendor/k8s.io/client-go/util/cert:go_default_library",
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
