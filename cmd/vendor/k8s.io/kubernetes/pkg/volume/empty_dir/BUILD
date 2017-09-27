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
        "doc.go",
        "empty_dir.go",
        "empty_dir_unsupported.go",
    ] + select({
        "@io_bazel_rules_go//go/platform:linux_amd64": [
            "empty_dir_linux.go",
        ],
        "//conditions:default": [],
    }),
    importpath = "k8s.io/kubernetes/pkg/volume/empty_dir",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/github.com/golang/glog:go_default_library",
        "//cmd/vendor/k8s.io/api/core/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/util/mount:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/util/strings:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/volume:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/volume/util:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = select({
        "@io_bazel_rules_go//go/platform:linux_amd64": [
            "empty_dir_test.go",
        ],
        "//conditions:default": [],
    }),
    importpath = "k8s.io/kubernetes/pkg/volume/empty_dir",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = select({
        "@io_bazel_rules_go//go/platform:linux_amd64": [
            "//cmd/vendor/k8s.io/api/core/v1:go_default_library",
            "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
            "//cmd/vendor/k8s.io/apimachinery/pkg/types:go_default_library",
            "//cmd/vendor/k8s.io/client-go/util/testing:go_default_library",
            "//cmd/vendor/k8s.io/kubernetes/pkg/util/mount:go_default_library",
            "//cmd/vendor/k8s.io/kubernetes/pkg/volume:go_default_library",
            "//cmd/vendor/k8s.io/kubernetes/pkg/volume/testing:go_default_library",
            "//cmd/vendor/k8s.io/kubernetes/pkg/volume/util:go_default_library",
        ],
        "//conditions:default": [],
    }),
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
