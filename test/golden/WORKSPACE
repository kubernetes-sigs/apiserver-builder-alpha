
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_k8s_repo_infra",
    sha256 = "5ff82744aad79b92b3963a26d779164d26b906aee0b177d66658be2c7a83617f",
    strip_prefix = "repo-infra-0.1.2",
    urls = [
        "https://github.com/kubernetes/repo-infra/archive/v0.1.2.tar.gz",
    ],
)

load("@io_k8s_repo_infra//:load.bzl", _repo_infra_repos = "repositories")

_repo_infra_repos()

load("@io_k8s_repo_infra//:repos.bzl", "configure")

# use k8s.io/repo-infra to configure go and bazel
# default minimum_bazel_version is 0.29.1
configure(
    go_version = "1.15",
    rbe_name = None,
)
