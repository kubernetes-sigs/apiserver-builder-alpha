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

    # built-in versions of package rules will no longer be addressed
    # see https://docs.bazel.build/versions/1.1.0/be/pkg.html#deprecated
    http_archive(
        name = "rules_pkg",
        urls = ["https://github.com/bazelbuild/rules_pkg/releases/download/0.2.4/rules_pkg-0.2.4.tar.gz"],
        sha256 = "4ba8f4ab0ff85f2484287ab06c0d871dcb31cc54d439457d28fd4ae14b18450a",
    )
