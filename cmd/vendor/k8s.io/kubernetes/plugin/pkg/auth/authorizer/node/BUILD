package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_test(
    name = "go_default_test",
    srcs = ["node_authorizer_test.go"],
    importpath = "k8s.io/kubernetes/plugin/pkg/auth/authorizer/node",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//cmd/vendor/k8s.io/apiserver/pkg/authentication/user:go_default_library",
        "//cmd/vendor/k8s.io/apiserver/pkg/authorization/authorizer:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/auth/nodeidentifier:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/plugin/pkg/auth/authorizer/rbac/bootstrappolicy:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "graph.go",
        "graph_populator.go",
        "node_authorizer.go",
    ],
    importpath = "k8s.io/kubernetes/plugin/pkg/auth/authorizer/node",
    tags = ["automanaged"],
    deps = [
        "//cmd/vendor/github.com/golang/glog:go_default_library",
        "//cmd/vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
        "//cmd/vendor/k8s.io/apiserver/pkg/authorization/authorizer:go_default_library",
        "//cmd/vendor/k8s.io/client-go/tools/cache:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api/persistentvolume:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/api/pod:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/apis/rbac:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/auth/nodeidentifier:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/pkg/client/informers/informers_generated/internalversion/core/internalversion:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/plugin/pkg/auth/authorizer/rbac:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/third_party/forked/gonum/graph:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/third_party/forked/gonum/graph/simple:go_default_library",
        "//cmd/vendor/k8s.io/kubernetes/third_party/forked/gonum/graph/traverse:go_default_library",
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
