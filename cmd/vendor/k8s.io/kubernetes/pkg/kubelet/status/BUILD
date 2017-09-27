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
        "generate.go",
        "status_manager.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/kubelet/status",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/github.com/golang/glog:go_default_library",
        "//cmd/vendor/k8s.io/api/core/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/api/equality:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/api/errors:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/diff:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        "//cmd/vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//cmd/vendor/k8s.io/client-go/kubernetes/scheme:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api/v1/pod:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/container:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/pod:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/types:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/util/format:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "generate_test.go",
        "status_manager_test.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/kubelet/status",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
        "//cmd/vendor/k8s.io/api/core/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/api/errors:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
        "//cmd/vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//cmd/vendor/k8s.io/client-go/kubernetes/fake:go_default_library",
        "//cmd/vendor/k8s.io/client-go/testing:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api/v1/pod:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/configmap:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/container:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/pod:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/pod/testing:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/secret:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/status/testing:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/kubelet/types:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//pkg/kubelet/status/testing:all-srcs",
    ],
    tags = ["automanaged"],
)
