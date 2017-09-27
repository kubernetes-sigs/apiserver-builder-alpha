package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["validation.go"],
    importpath = "k8s.io/kubernetes/staging/src/k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/validation",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/apiextensions-apiserver/pkg/apis/apiextensions:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/api/validation:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/validation:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/validation/field:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["validation_test.go"],
    importpath = "k8s.io/kubernetes/staging/src/k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/validation",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/apiextensions-apiserver/pkg/apis/apiextensions:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/validation/field:go_default_library",
    ],
)
