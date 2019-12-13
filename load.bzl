load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def repositories():
    http_archive(
        name = "io_k8s_repo_infra",
        sha256 = "5ee2a8e306af0aaf2844b5e2c79b5f3f53fc9ce3532233f0615b8d0265902b2a",
        strip_prefix = "repo-infra-0.0.1-alpha.1",
        urls = [
            "https://github.com/kubernetes/repo-infra/archive/v0.0.1-alpha.1.tar.gz",
        ],
    )
