# gazelle:repository_macro repos.bzl%go_repositories
workspace(name = "io_k8s_sigs_apiserver_builder_alpha")

load("//:load.bzl", "repositories")

repositories()

load("@io_k8s_repo_infra//:load.bzl", _repo_infra_repos = "repositories")

_repo_infra_repos()

load("@io_k8s_repo_infra//:repos.bzl", "configure")

# use k8s.io/repo-infra to configure go and bazel
# default minimum_bazel_version is 0.29.1
configure(
    go_version = "1.13",
    rbe_name = None,
)

load("//:repos.bzl", "go_repositories")

# load go dependencies imported from cmd/go.mod
go_repositories()
