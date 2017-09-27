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
        "attacher.go",
        "aws_ebs.go",
        "aws_util.go",
        "doc.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/volume/aws_ebs",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/github.com/golang/glog:go_default_library",
        "//cmd/vendor/k8s.io/api/core/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/api/resource:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/cloudprovider:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/cloudprovider/providers/aws:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/util/exec:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/util/mount:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/util/strings:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/volume:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/volume/util:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/volume/util/volumehelper:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "attacher_test.go",
        "aws_ebs_test.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/volume/aws_ebs",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/github.com/golang/glog:go_default_library",
        "//cmd/vendor/k8s.io/api/core/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//cmd/vendor/k8s.io/client-go/kubernetes/fake:go_default_library",
        "//cmd/vendor/k8s.io/client-go/util/testing:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/cloudprovider/providers/aws:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/util/mount:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/volume:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/volume/testing:go_default_library",
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
