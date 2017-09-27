package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["reconciler.go"],
    importpath = "k8s.io/kubernetes/pkg/kubelet/volumemanager/reconciler",
    tags = ["automanaged"],
    deps = [
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//vendor/k8s.io/kubernetes/cmd/kubelet/app/options:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/kubelet/volumemanager/cache:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/file:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/goroutinemap/exponentialbackoff:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/mount:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/strings:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/volume:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/volume/util/nestedpendingoperations:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/volume/util/operationexecutor:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/volume/util/types:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/volume/util/volumehelper:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["reconciler_test.go"],
    importpath = "k8s.io/kubernetes/pkg/kubelet/volumemanager/reconciler",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes/fake:go_default_library",
        "//vendor/k8s.io/client-go/testing:go_default_library",
        "//vendor/k8s.io/client-go/tools/record:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/kubelet/volumemanager/cache:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/util/mount:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/volume:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/volume/testing:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/volume/util/operationexecutor:go_default_library",
        "//vendor/k8s.io/kubernetes/pkg/volume/util/volumehelper:go_default_library",
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
