load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_repositories():
    go_repository(
        name = "com_github_ajg_form",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/ajg/form",
        sum = "h1:zvClvFQwU++UpIUBGC8YmDlfhUrweEy1R1Fj1gu5iIM=",
        version = "v0.0.0-20160822230020-523a5da1a92f",
    )
    go_repository(
        name = "com_github_azure_go_ansiterm",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/Azure/go-ansiterm",
        sum = "h1:w+iIsaOQNcT7OZ575w+acHgRric5iCyQh+xv+KJ4HB8=",
        version = "v0.0.0-20170929234023-d6e3b3328b78",
    )
    go_repository(
        name = "com_github_azure_go_autorest",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/Azure/go-autorest",
        sum = "h1:viZ3tV5l4gE2Sw0xrasFHytCGtzYCrT+um/rrSQ1BfA=",
        version = "v11.1.2+incompatible",
    )
    go_repository(
        name = "com_github_beorn7_perks",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/beorn7/perks",
        sum = "h1:xJ4a3vCFaGF/jqvzLMYoU8P317H5OQ+Via4RmuPwCS0=",
        version = "v0.0.0-20180321164747-3a771d992973",
    )
    go_repository(
        name = "com_github_blang_semver",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/blang/semver",
        sum = "h1:CGxCgetQ64DKk7rdZ++Vfnb1+ogGNnB17OJKJXD2Cfs=",
        version = "v3.5.0+incompatible",
    )
    go_repository(
        name = "com_github_burntsushi_toml",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/BurntSushi/toml",
        sum = "h1:WXkYYl6Yr3qBf1K79EBnL4mak0OimBfB0XUf9Vl28OQ=",
        version = "v0.3.1",
    )
    go_repository(
        name = "com_github_cockroachdb_apd",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/cockroachdb/apd",
        sum = "h1:3LFP3629v+1aKXU5Q37mxmRxX/pIu1nijXydLShEq5I=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_cockroachdb_cockroach_go",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/cockroachdb/cockroach-go",
        sum = "h1:2zRrJWIt/f9c9HhNHAgrRgq0San5gRRUJTBXLkchal0=",
        version = "v0.0.0-20181001143604-e0a95dfd547c",
    )
    go_repository(
        name = "com_github_codegangsta_negroni",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/codegangsta/negroni",
        sum = "h1:+aYywywx4bnKXWvoWtRfJ91vC59NbEhEY03sZjQhbVY=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_coreos_bbolt",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/coreos/bbolt",
        sum = "h1:uTXKg9gY70s9jMAKdfljFQcuh4e/BXOM+V+d00KFj3A=",
        version = "v1.3.1-coreos.6",
    )
    go_repository(
        name = "com_github_coreos_etcd",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/coreos/etcd",
        sum = "h1:8F3hqu9fGYLBifCmRCJsicFqDx/D68Rt3q1JMazcgBQ=",
        version = "v3.3.13+incompatible",
    )
    go_repository(
        name = "com_github_coreos_go_oidc",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/coreos/go-oidc",
        sum = "h1:X+JQSgXg3CcxgcBoMAqU8NoS0fch8zHxjiKWcXclxaI=",
        version = "v0.0.0-20180117170138-065b426bd416",
    )
    go_repository(
        name = "com_github_coreos_go_semver",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/coreos/go-semver",
        sum = "h1:WqY2Kv7eI1jeoU3pC05YYK/kK4tdXyLzzaBzCR51r9M=",
        version = "v0.0.0-20180108230905-e214231b295a",
    )
    go_repository(
        name = "com_github_coreos_go_systemd",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/coreos/go-systemd",
        sum = "h1:u9SHYsPQNyt5tgDm3YN7+9dYrpK96E5wFilTFWIDZOM=",
        version = "v0.0.0-20180511133405-39ca1b05acc7",
    )
    go_repository(
        name = "com_github_coreos_pkg",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/coreos/pkg",
        sum = "h1:n2Ltr3SrfQlf/9nOna1DoGKxLx3qTSI8Ttl6Xrqp6mw=",
        version = "v0.0.0-20180108230652-97fdf19511ea",
    )
    go_repository(
        name = "com_github_davecgh_go_spew",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/davecgh/go-spew",
        sum = "h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_dgrijalva_jwt_go",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/dgrijalva/jwt-go",
        sum = "h1:7qlOGliEKZXTDg6OTjfoBKDXWrumCAMpl/TFQ4/5kLM=",
        version = "v3.2.0+incompatible",
    )
    go_repository(
        name = "com_github_docker_docker",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/docker/docker",
        sum = "h1:w3NnFcKR5241cfmQU5ZZAsf0xcpId6mWOupTvJlUX2U=",
        version = "v0.7.3-0.20190327010347-be7ac8be2ae0",
    )
    go_repository(
        name = "com_github_docker_spdystream",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/docker/spdystream",
        sum = "h1:cenwrSVm+Z7QLSV/BsnenAOcDXdX4cMv4wP0B/5QbPg=",
        version = "v0.0.0-20160310174837-449fdfce4d96",
    )
    go_repository(
        name = "com_github_dustin_go_humanize",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/dustin/go-humanize",
        sum = "h1:VSnTsYCnlFHaM2/igO1h6X3HA71jcobQuxemgkq4zYo=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_elazarl_goproxy",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/elazarl/goproxy",
        sum = "h1:p1yVGRW3nmb85p1Sh1ZJSDm4A4iKLS5QNbvUHMgGu/M=",
        version = "v0.0.0-20170405201442-c4fc26588b6e",
    )
    go_repository(
        name = "com_github_emicklei_go_restful",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/emicklei/go-restful",
        sum = "h1:H2pdYOb3KQ1/YsqVWoWNLQO+fusocsw354rqGTZtAgw=",
        version = "v0.0.0-20170410110728-ff4f55a20633",
    )
    go_repository(
        name = "com_github_evanphx_json_patch",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/evanphx/json-patch",
        sum = "h1:fUDGZCv/7iAN7u0puUVhvKCcsR6vRfwrJatElLBEf0I=",
        version = "v4.2.0+incompatible",
    )
    go_repository(
        name = "com_github_fatih_color",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/fatih/color",
        sum = "h1:DkWD4oS2D8LGGgTQ6IvwJJXSL5Vp2ffcQg58nFV38Ys=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_github_fatih_structs",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/fatih/structs",
        sum = "h1:Q7juDM0QtcnhCpeyLGQKyg4TOIghuNXrkL32pHAUMxo=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_fsnotify_fsnotify",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/fsnotify/fsnotify",
        sum = "h1:IXs+QLmnXW2CcXuY+8Mzv/fWEsPGWxqefPtCP5CnV9I=",
        version = "v1.4.7",
    )
    go_repository(
        name = "com_github_ghodss_yaml",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/ghodss/yaml",
        sum = "h1:bRzFpEzvausOAt4va+I/22BZ1vXDtERngp0BNYDKej0=",
        version = "v0.0.0-20180820084758-c7ce16629ff4",
    )
    go_repository(
        name = "com_github_go_openapi_jsonpointer",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/go-openapi/jsonpointer",
        sum = "h1:gihV7YNZK1iK6Tgwwsxo2rJbD1GTbdm72325Bq8FI3w=",
        version = "v0.19.3",
    )
    go_repository(
        name = "com_github_go_openapi_jsonreference",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/go-openapi/jsonreference",
        sum = "h1:o20suLFB4Ri0tuzpWtyHlh7E7HnkqTNLq6aR6WVNS1w=",
        version = "v0.19.2",
    )
    go_repository(
        name = "com_github_go_openapi_spec",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/go-openapi/spec",
        sum = "h1:0XRyw8kguri6Yw4SxhsQA/atC88yqrk0+G4YhI2wabc=",
        version = "v0.19.3",
    )
    go_repository(
        name = "com_github_go_openapi_swag",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/go-openapi/swag",
        sum = "h1:lTz6Ys4CmqqCQmZPBlbQENR1/GucA2bzYTE12Pw4tFY=",
        version = "v0.19.5",
    )
    go_repository(
        name = "com_github_go_sql_driver_mysql",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/go-sql-driver/mysql",
        sum = "h1:7LxgVwFb2hIQtMm87NdgAVfXjnt4OePseqT1tKx+opk=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_github_gobuffalo_buffalo",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/buffalo",
        sum = "h1:Fyn55HJULJpFPMUNx9lrPK31qvr37+bpNGFbpAOauGI=",
        version = "v0.13.0",
    )
    go_repository(
        name = "com_github_gobuffalo_buffalo_plugins",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/buffalo-plugins",
        sum = "h1:odC6PspRooQxDoBy1QvusRGLshbK/2uRkxkSB54nc30=",
        version = "v1.6.11",
    )
    go_repository(
        name = "com_github_gobuffalo_buffalo_pop",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/buffalo-pop",
        sum = "h1:8aXdBlo9MEKFGHyl489+28Jw7Ud59Th1U+5Ayu1wNL0=",
        version = "v1.0.5",
    )
    go_repository(
        name = "com_github_gobuffalo_envy",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/envy",
        sum = "h1:zvJZDJ3Aj42EDL5cjBcH9q4mWJHkoGyFzfmVP+weYxo=",
        version = "v1.6.10",
    )
    go_repository(
        name = "com_github_gobuffalo_events",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/events",
        sum = "h1:T9SXUVyO1kF61uns5a8cpqsSPj5txZgIplbOoNNltqk=",
        version = "v1.1.8",
    )
    go_repository(
        name = "com_github_gobuffalo_fizz",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/fizz",
        sum = "h1:JJOkmlStog5AiBL434UoGMJ896p3MnTnzedFVaZSF3k=",
        version = "v1.0.12",
    )
    go_repository(
        name = "com_github_gobuffalo_flect",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/flect",
        sum = "h1:nvA0/snr4wQeCwBYmrbftniJun/kxOjK/Pz3ivb7wis=",
        version = "v0.0.0-20181114183036-47375f6d8328",
    )
    go_repository(
        name = "com_github_gobuffalo_genny",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/genny",
        sum = "h1:RC+rmtL7i+1RXTg0451HLcei1aNs88k8BsbfS9gC3+M=",
        version = "v0.0.0-20181119162812-e8ff4adce8bb",
    )
    go_repository(
        name = "com_github_gobuffalo_github_flavored_markdown",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/github_flavored_markdown",
        sum = "h1:Vjvz4wqOnviiLEfTh5bh270b3lhpJiwwQEWOWmHMwY8=",
        version = "v1.0.7",
    )
    go_repository(
        name = "com_github_gobuffalo_httptest",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/httptest",
        sum = "h1:LWp2khlgA697h4BIYWW2aRxvB93jMnBrbakQ/r2KLzs=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_gobuffalo_licenser",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/licenser",
        sum = "h1:/iWW53J3dNmrjkG5C7IFHhy4SFbhikdAvGtFwvMnh/U=",
        version = "v0.0.0-20181109171355-91a2a7aac9a7",
    )
    go_repository(
        name = "com_github_gobuffalo_logger",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/logger",
        sum = "h1:Z/ppYX6EtPEysbW4VEGz2dO+4F4VTthWp2sWRUCANdU=",
        version = "v0.0.0-20181127160119-5b956e21995c",
    )
    go_repository(
        name = "com_github_gobuffalo_makr",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/makr",
        sum = "h1:lOlpv2iz0dNa4qse0ZYQgbtT+ybwVxWEAcOZbcPmeYc=",
        version = "v1.1.5",
    )
    go_repository(
        name = "com_github_gobuffalo_mapi",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/mapi",
        sum = "h1:JRuTiZzDEZhBHkFiHTxJkYRT6CbYuL0K/rn+1byJoEA=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_gobuffalo_meta",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/meta",
        sum = "h1:E3y9POfXjqW9KuXWF75Ny1zeSakkctzRAnjUSKK3xHE=",
        version = "v0.0.0-20181114191255-b130ebedd2f7",
    )
    go_repository(
        name = "com_github_gobuffalo_mw_basicauth",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/mw-basicauth",
        sum = "h1:bCqDBHnByenQitOtFdEtMvlWVgPwODrfZ+nVkgGoJZ8=",
        version = "v1.0.3",
    )
    go_repository(
        name = "com_github_gobuffalo_mw_contenttype",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/mw-contenttype",
        sum = "h1:SUFp8EbFjlKXkvqstoxPWx3nVPV3BSKZTswQNTZFaik=",
        version = "v0.0.0-20180802152300-74f5a47f4d56",
    )
    go_repository(
        name = "com_github_gobuffalo_mw_csrf",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/mw-csrf",
        sum = "h1:A13B4mhcFQcjPJ1GFBrh61B4Qo87fZa82FfTt9LX/QU=",
        version = "v0.0.0-20180802151833-446ff26e108b",
    )
    go_repository(
        name = "com_github_gobuffalo_mw_forcessl",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/mw-forcessl",
        sum = "h1:v94+IGhlBro0Lz1gOR3lrdAVSZ0mJF2NxsdppKd7FnI=",
        version = "v0.0.0-20180802152810-73921ae7a130",
    )
    go_repository(
        name = "com_github_gobuffalo_mw_i18n",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/mw-i18n",
        sum = "h1:pZhsgF8RXEngHdibuRNOXNk1pL0K9rFa5HOcvURNTQ4=",
        version = "v0.0.0-20180802152014-e3060b7e13d6",
    )
    go_repository(
        name = "com_github_gobuffalo_mw_paramlogger",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/mw-paramlogger",
        sum = "h1:TsmUXyHjj5ReuN1AJjEVukf72J6AfRTF2CfTEaqVLT8=",
        version = "v0.0.0-20181005191442-d6ee392ec72e",
    )
    go_repository(
        name = "com_github_gobuffalo_mw_tokenauth",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/mw-tokenauth",
        sum = "h1:nhPzONHNGlXZIMFfKm6cWpRSq5oTanRK1qBtfCPBFyE=",
        version = "v0.0.0-20181001105134-8545f626c189",
    )
    go_repository(
        name = "com_github_gobuffalo_packd",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/packd",
        sum = "h1:039sqjIDVgXXcB5hOIKaSBL1kZpQj9u8WLXuu2xOcX8=",
        version = "v0.0.0-20181124090624-311c6248e5fb",
    )
    go_repository(
        name = "com_github_gobuffalo_packr",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/packr",
        sum = "h1:p2ujcDJQp2QTiYWcI0ByHbr/gMoCouok6M0vXs/yTYQ=",
        version = "v1.21.0",
    )
    go_repository(
        name = "com_github_gobuffalo_packr_v2",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/packr/v2",
        sum = "h1:bnphKdb4fGtSq3muDvV0LP3d0cTixcGVPjCwQSfjFdA=",
        version = "v2.0.0-rc.8",
    )
    go_repository(
        name = "com_github_gobuffalo_plush",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/plush",
        sum = "h1:SJxTrFxizR3wCxb7GOYf1ZUjo3mfAczK4B/FYROZg00=",
        version = "v3.7.22+incompatible",
    )
    go_repository(
        name = "com_github_gobuffalo_pop",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/pop",
        sum = "h1:5v15ZgICK3MFTU90QRqCaqDUf4wcriIbws1hqpYL2Xo=",
        version = "v4.8.4+incompatible",
    )
    go_repository(
        name = "com_github_gobuffalo_release",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/release",
        sum = "h1:hSzdILdeu3utEwV+EEVLUq9iNPl+twNrdO7KVl+8W38=",
        version = "v1.0.72",
    )
    go_repository(
        name = "com_github_gobuffalo_shoulders",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/shoulders",
        sum = "h1:BqVJBUXlBWAf+WLhXijVk3SCpp75LXrVBiIkOCzZbNc=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_gobuffalo_syncx",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/syncx",
        sum = "h1:S5EeH1reN93KR0L6TQvkRpu9YggCYXrUqFh1iEgvdC0=",
        version = "v0.0.0-20181120194010-558ac7de985f",
    )
    go_repository(
        name = "com_github_gobuffalo_tags",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/tags",
        sum = "h1:zLkaontB8lWefU+DX38mzPLRKFGTJL8FKb9JnKMt0Z0=",
        version = "v2.0.11+incompatible",
    )
    go_repository(
        name = "com_github_gobuffalo_uuid",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/uuid",
        sum = "h1:c5uWRuEnYggYCrT9AJm0U2v1QTG7OVDAvxhj8tIV5Gc=",
        version = "v2.0.5+incompatible",
    )
    go_repository(
        name = "com_github_gobuffalo_validate",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/validate",
        sum = "h1:6f4JCEz11Zi6iIlexMv7Jz10RBPvgI795AOaubtCwTE=",
        version = "v2.0.3+incompatible",
    )
    go_repository(
        name = "com_github_gobuffalo_x",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gobuffalo/x",
        sum = "h1:N0iqtKwkicU8M2rLirTDJxdwuL8I2/8MjMlEayaNSgE=",
        version = "v0.0.0-20181007152206-913e47c59ca7",
    )
    go_repository(
        name = "com_github_gofrs_uuid",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gofrs/uuid",
        sum = "h1:q2rtkjaKT4YEr6E1kamy0Ha4RtepWlQBedyHx0uzKwA=",
        version = "v3.1.0+incompatible",
    )
    go_repository(
        name = "com_github_gogo_protobuf",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gogo/protobuf",
        sum = "h1:WSBJMqJbLxsn+bTCPyPYZfqHdJmc8MK4wrBjMft6BAM=",
        version = "v0.0.0-20171007142547-342cbe0a0415",
    )
    go_repository(
        name = "com_github_golang_groupcache",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/golang/groupcache",
        sum = "h1:LbsanbbD6LieFkXbj9YNNBupiGHJgFeLpO0j0Fza1h8=",
        version = "v0.0.0-20160516000752-02826c3e7903",
    )
    go_repository(
        name = "com_github_golang_protobuf",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/golang/protobuf",
        sum = "h1:P3YflyNX/ehuJFLhxviNdFxQPkGK5cDcApsge1SqnvM=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_google_btree",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/google/btree",
        sum = "h1:JHB7F/4TJCrYBW8+GZO8VkWDj1jxcWuCl6uxKODiyi4=",
        version = "v0.0.0-20160524151835-7d79101e329e",
    )
    go_repository(
        name = "com_github_google_go_cmp",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/google/go-cmp",
        sum = "h1:crn/baboCvb5fXaQ0IJ1SGTsTVrWpDsCWC8EGETZijY=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_google_gofuzz",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/google/gofuzz",
        sum = "h1:+RRA9JqSOZFfKrOeqr2z77+8R2RKyh8PG66dcu1V0ck=",
        version = "v0.0.0-20170612174753-24818f796faf",
    )
    go_repository(
        name = "com_github_google_uuid",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/google/uuid",
        sum = "h1:Gkbcsh/GbpXz7lPftLA3P6TYMwjCLYm83jiFQZF/3gY=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_googleapis_gnostic",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/googleapis/gnostic",
        sum = "h1:7XGaL1e6bYS1yIonGp9761ExpPPV1ui0SAC59Yube9k=",
        version = "v0.0.0-20170729233727-0c5108395e2d",
    )
    go_repository(
        name = "com_github_gophercloud_gophercloud",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gophercloud/gophercloud",
        sum = "h1:L9JPKrtsHMQ4VCRQfHvbbHBfB2Urn8xf6QZeXZ+OrN4=",
        version = "v0.0.0-20190126172459-c818fa66e4c8",
    )
    go_repository(
        name = "com_github_gorilla_context",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gorilla/context",
        sum = "h1:AWwleXJkX/nhcU9bZSnZoi3h/qGYqQAGhq6zZe/aQW8=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_gorilla_mux",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gorilla/mux",
        sum = "h1:Pgr17XVTNXAk3q/r4CpKzC5xBM/qW1uVLV+IhRZpIIk=",
        version = "v1.6.2",
    )
    go_repository(
        name = "com_github_gorilla_pat",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gorilla/pat",
        sum = "h1:LqbZZ9sNMWVjeXS4NN5oVvhMjDyLhmA1LG86oSo+IqY=",
        version = "v0.0.0-20180118222023-199c85a7f6d1",
    )
    go_repository(
        name = "com_github_gorilla_securecookie",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gorilla/securecookie",
        sum = "h1:miw7JPhV+b/lAHSXz4qd/nN9jRiAFV5FwjeKyCS8BvQ=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_gorilla_sessions",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gorilla/sessions",
        sum = "h1:uXoZdcdA5XdXF3QzuSlheVRUvjl+1rKY7zBXL68L9RU=",
        version = "v1.1.3",
    )
    go_repository(
        name = "com_github_gorilla_websocket",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gorilla/websocket",
        sum = "h1:Lh2aW+HnU2Nbe1gqD9SOJLJxW1jBMmQOktN2acDyJk8=",
        version = "v0.0.0-20170926233335-4201258b820c",
    )
    go_repository(
        name = "com_github_gregjones_httpcache",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/gregjones/httpcache",
        sum = "h1:6TSoaYExHper8PYsJu23GWVNOyYRCSnIFyxKgLSZ54w=",
        version = "v0.0.0-20170728041850-787624de3eb7",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_go_grpc_middleware",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/grpc-ecosystem/go-grpc-middleware",
        sum = "h1:lR9ssWAqp9qL0bALxqEEkuudiP1eweOdv9jsRK3e7lE=",
        version = "v0.0.0-20190222133341-cfaf5686ec79",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_go_grpc_prometheus",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/grpc-ecosystem/go-grpc-prometheus",
        sum = "h1:f5vL2EW5pL274ztMNnizZAEa457nKyKPEaN/sm/kdBk=",
        version = "v0.0.0-20170330212424-2500245aa611",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_grpc_gateway",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/grpc-ecosystem/grpc-gateway",
        sum = "h1:HJtP6RRwj2EpPCD/mhAWzSvLL/dFTdPm1UrWwanoFos=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_hashicorp_golang_lru",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/hashicorp/golang-lru",
        sum = "h1:CL2msUPvZTLb5O648aiLNJw3hnBxN2+1Jq8rCOH9wdo=",
        version = "v0.5.0",
    )
    go_repository(
        name = "com_github_hashicorp_hcl",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/hashicorp/hcl",
        sum = "h1:0Anlzjpi4vEasTeNFn2mLJgTSwt0+6sfsiTG8qcWGx4=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_hpcloud_tail",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/hpcloud/tail",
        sum = "h1:nfCOvKYfkgYP8hkirhJocXT2+zOD8yUNjXaWfTlyFKI=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_imdario_mergo",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/imdario/mergo",
        sum = "h1:JboBksRwiiAJWvIYJVo46AfV+IAIKZpfrSzVKj42R4Q=",
        version = "v0.3.5",
    )
    go_repository(
        name = "com_github_inconshreveable_mousetrap",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/inconshreveable/mousetrap",
        sum = "h1:Z8tu5sraLXCXIcARxBp/8cbvlwVa7Z1NHg9XEKhtSvM=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_jackc_fake",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/jackc/fake",
        sum = "h1:vr3AYkKovP8uR8AvSGGUK1IDqRa5lAAvEkZG1LKaCRc=",
        version = "v0.0.0-20150926172116-812a484cc733",
    )
    go_repository(
        name = "com_github_jackc_pgx",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/jackc/pgx",
        sum = "h1:0Vihzu20St42/UDsvZGdNE6jak7oi/UOeMzwMPHkgFY=",
        version = "v3.2.0+incompatible",
    )
    go_repository(
        name = "com_github_jmoiron_sqlx",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/jmoiron/sqlx",
        sum = "h1:5B0uxl2lzNRVkJVg+uGHxWtRt4C0Wjc6kJKo5XYx8xE=",
        version = "v0.0.0-20180614180643-0dae4fefe7c0",
    )
    go_repository(
        name = "com_github_joho_godotenv",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/joho/godotenv",
        sum = "h1:Zjp+RcGpHhGlrMbJzXTrZZPrWj+1vfm90La1wgB6Bhc=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_jonboulle_clockwork",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/jonboulle/clockwork",
        sum = "h1:XpRROA6ssPlTwJI8/pH+61uieOkcJhmAFz25cu0B94Y=",
        version = "v0.0.0-20141017032234-72f9bd7c4e0c",
    )
    go_repository(
        name = "com_github_json_iterator_go",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/json-iterator/go",
        sum = "h1:AHimNtVIpiBjPUhEF5KNCkrUyqTSA5zWUl8sQ2bfGBE=",
        version = "v0.0.0-20180701071628-ab8a2e0c74be",
    )
    go_repository(
        name = "com_github_karrick_godirwalk",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/karrick/godirwalk",
        sum = "h1:VbzFqwXwNbAZoA6W5odrLr+hKK197CcENcPh6E/gJ0M=",
        version = "v1.7.5",
    )
    go_repository(
        name = "com_github_kballard_go_shellquote",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/kballard/go-shellquote",
        sum = "h1:Z9n2FFNUXsshfwJMBgNA0RU6/i7WVaAegv3PtuIHPMs=",
        version = "v0.0.0-20180428030007-95032a82bc51",
    )
    go_repository(
        name = "com_github_konsorten_go_windows_terminal_sequences",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/konsorten/go-windows-terminal-sequences",
        sum = "h1:mweAR1A6xJ3oS2pRaGiHgQ4OO8tzTaLawm8vnODuwDk=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_kr_pretty",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/kr/pretty",
        sum = "h1:L/CwN0zerZDmRFUapSPitk6f+Q3+0za1rQkzVuMiMFI=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_kr_pty",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/kr/pty",
        sum = "h1:hyz3dwM5QLc1Rfoz4FuWJQG5BN7tc6K1MndAUnGpQr4=",
        version = "v1.1.5",
    )
    go_repository(
        name = "com_github_kr_text",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/kr/text",
        sum = "h1:45sCR5RtlFHMR4UwH9sdQ5TC8v0qDQCHnXt+kaKSTVE=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_lib_pq",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/lib/pq",
        sum = "h1:X5PMW56eZitiTeO7tKzZxFCSpbFZJtkMMooicw2us9A=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_magiconair_properties",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/magiconair/properties",
        sum = "h1:LLgXmsheXeRoUOBOjtwPQCWIYqM/LU1ayDtDePerRcY=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_github_mailru_easyjson",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/mailru/easyjson",
        sum = "h1:hB2xlXdHp/pmPZq0y3QnmWAArdw9PqbmotexnWx/FU8=",
        version = "v0.0.0-20190626092158-b2ccc519800e",
    )
    go_repository(
        name = "com_github_markbates_deplist",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/markbates/deplist",
        sum = "h1:BKTJDTV5EynLGvTyONdgYVvV34DWq20mJAfGBCP+AYs=",
        version = "v1.0.5",
    )
    go_repository(
        name = "com_github_markbates_going",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/markbates/going",
        sum = "h1:uNQHDDfMRNOUmuxDbPbvatyw4wr4UOSUZkGkdkcip1o=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_markbates_grift",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/markbates/grift",
        sum = "h1:JjTyhlgPtgEnyHNvVn5lk21zWQbWD3cGE0YdyvvbZYg=",
        version = "v1.0.4",
    )
    go_repository(
        name = "com_github_markbates_hmax",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/markbates/hmax",
        sum = "h1:yo2N0gBoCnUMKhV/VRLHomT6Y9wUm+oQQENuWJqCdlM=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_markbates_inflect",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/markbates/inflect",
        sum = "h1:5fh1gzTFhfae06u3hzHYO9xe3l3v3nW5Pwt3naLTP5g=",
        version = "v1.0.4",
    )
    go_repository(
        name = "com_github_markbates_oncer",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/markbates/oncer",
        sum = "h1:Mlji5gkcpzkqTROyE4ZxZ8hN7osunMb2RuGVrbvMvCc=",
        version = "v0.0.0-20181014194634-05fccaae8fc4",
    )
    go_repository(
        name = "com_github_markbates_refresh",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/markbates/refresh",
        sum = "h1:6EZ/vvVpWiam8OTIhrhfV9cVJR/NvScvcCiqosbTkbA=",
        version = "v1.4.10",
    )
    go_repository(
        name = "com_github_markbates_safe",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/markbates/safe",
        sum = "h1:yjZkbvRM6IzKj9tlu/zMJLS0n/V351OZWRnF3QfaUxI=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_markbates_sigtx",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/markbates/sigtx",
        sum = "h1:y/xtkBvNPRjD4KeEplf4w9rJVSc23/xl+jXYGowTwy0=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_markbates_willie",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/markbates/willie",
        sum = "h1:394PpHImWjScL9X2VRCDXJAcc77sHsSr3w3sOnL/DVc=",
        version = "v1.0.9",
    )
    go_repository(
        name = "com_github_masterminds_semver",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/Masterminds/semver",
        sum = "h1:WBLTQ37jOCzSLtXNdoo8bNM8876KhNqOKvrlGITgsTc=",
        version = "v1.4.2",
    )
    go_repository(
        name = "com_github_mattn_go_colorable",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/mattn/go-colorable",
        sum = "h1:UVL0vNpWh04HeJXV0KLcaT7r06gOH2l4OW6ddYRUIY4=",
        version = "v0.0.9",
    )
    go_repository(
        name = "com_github_mattn_go_isatty",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/mattn/go-isatty",
        sum = "h1:bnP0vzxcAdeI1zdubAl5PjU6zsERjGZb7raWodagDYs=",
        version = "v0.0.4",
    )
    go_repository(
        name = "com_github_mattn_go_sqlite3",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/mattn/go-sqlite3",
        sum = "h1:pDRiWfl+++eC2FEFRy6jXmQlvp4Yh3z1MJKg4UeYM/4=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_github_matttproud_golang_protobuf_extensions",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/matttproud/golang_protobuf_extensions",
        sum = "h1:4hp9jkHxhMHkqkrB3Ix0jegS5sx/RkqARlsWZ6pIwiU=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_microcosm_cc_bluemonday",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/microcosm-cc/bluemonday",
        sum = "h1:SIYunPjnlXcW+gVfvm0IlSeR5U3WZUOLfVmqg85Go44=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_mitchellh_go_homedir",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/mitchellh/go-homedir",
        sum = "h1:vKb8ShqSby24Yrqr/yDYkuFz8d0WUjys40rvnGC8aR0=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_mitchellh_mapstructure",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/mitchellh/mapstructure",
        sum = "h1:fmNYVwqnSfB9mZU6OS2O6GsXM+wcskZDuKQzvN1EDeE=",
        version = "v1.1.2",
    )
    go_repository(
        name = "com_github_modern_go_concurrent",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/modern-go/concurrent",
        sum = "h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=",
        version = "v0.0.0-20180306012644-bacd9c7ef1dd",
    )
    go_repository(
        name = "com_github_modern_go_reflect2",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/modern-go/reflect2",
        sum = "h1:9f412s+6RmYXLWZSEzVVgPGK7C2PphHj5RJrvfx9AWI=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_monoculum_formam",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/monoculum/formam",
        sum = "h1:FEJJhVHSH+Kyxa5qNe/7dprlZbFcj2TG51OWIouhwls=",
        version = "v0.0.0-20180901015400-4e68be1d79ba",
    )
    go_repository(
        name = "com_github_munnerz_goautoneg",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/munnerz/goautoneg",
        sum = "h1:7PxY7LVfSZm7PEeBTyK1rj1gABdCO2mbri6GKO1cMDs=",
        version = "v0.0.0-20120707110453-a547fc61f48d",
    )
    go_repository(
        name = "com_github_mxk_go_flowrate",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/mxk/go-flowrate",
        sum = "h1:y5//uYreIhSUg3J1GEMiLbxo1LJaP8RfCpH6pymGZus=",
        version = "v0.0.0-20140419014527-cca7078d478f",
    )
    go_repository(
        name = "com_github_natefinch_lumberjack",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/natefinch/lumberjack",
        sum = "h1:4QJd3OLAMgj7ph+yZTuX13Ld4UpgHp07nNdFX7mqFfM=",
        version = "v2.0.0+incompatible",
    )
    go_repository(
        name = "com_github_nicksnyder_go_i18n",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/nicksnyder/go-i18n",
        sum = "h1:5AzlPKvXBH4qBzmZ09Ua9Gipyruv6uApMcrNZdo96+Q=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_github_nytimes_gziphandler",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/NYTimes/gziphandler",
        sum = "h1:lsxEuwrXEAokXB9qhlbKWPpo3KMLZQ5WB5WLQRW1uq0=",
        version = "v0.0.0-20170623195520-56545f4a5d46",
    )
    go_repository(
        name = "com_github_onsi_ginkgo",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/onsi/ginkgo",
        sum = "h1:Ix8l273rp3QzYgXSR+c8d1fTG7UPgYkOSELPhiY/YGw=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_onsi_gomega",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/onsi/gomega",
        sum = "h1:3mYCb7aPxS/RU7TI1y4rkEn1oKmPRjNJLNEXgw7MH2I=",
        version = "v1.4.2",
    )
    go_repository(
        name = "com_github_pborman_uuid",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/pborman/uuid",
        sum = "h1:J7Q5mO4ysT1dv8hyrUGHb9+ooztCXu1D8MY8DZYsu3g=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_pelletier_go_toml",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/pelletier/go-toml",
        sum = "h1:T5zMGML61Wp+FlcbWjRDT7yAxhJNAiPPLOFECq181zc=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_peterbourgon_diskv",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/peterbourgon/diskv",
        sum = "h1:UBdAOUP5p4RWqPBg048CAvpKN+vxiaj6gdUUzhl4XmI=",
        version = "v2.0.1+incompatible",
    )
    go_repository(
        name = "com_github_pkg_errors",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/pkg/errors",
        sum = "h1:WdK/asTD0HN+q6hsWO3/vpuAkAr+tw6aNJNDFFf0+qw=",
        version = "v0.8.0",
    )
    go_repository(
        name = "com_github_pmezard_go_difflib",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/pmezard/go-difflib",
        sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_pquerna_cachecontrol",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/pquerna/cachecontrol",
        sum = "h1:0XM1XL/OFFJjXsYXlG30spTkV/E9+gmd5GD1w2HE8xM=",
        version = "v0.0.0-20171018203845-0dec1b30a021",
    )
    go_repository(
        name = "com_github_prometheus_client_golang",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/prometheus/client_golang",
        sum = "h1:awm861/B8OKDd2I/6o1dy3ra4BamzKhYOiGItCeZ740=",
        version = "v0.9.2",
    )
    go_repository(
        name = "com_github_prometheus_client_model",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/prometheus/client_model",
        sum = "h1:idejC8f05m9MGOsuEi1ATq9shN03HrxNkD/luQvxCv8=",
        version = "v0.0.0-20180712105110-5c3871d89910",
    )
    go_repository(
        name = "com_github_prometheus_common",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/prometheus/common",
        sum = "h1:PnBWHBf+6L0jOqq0gIVUe6Yk0/QMZ640k6NvkxcBf+8=",
        version = "v0.0.0-20181126121408-4724e9255275",
    )
    go_repository(
        name = "com_github_prometheus_procfs",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/prometheus/procfs",
        sum = "h1:9a8MnZMP0X2nLJdBg+pBmGgkJlSaKC2KaQmTCk1XDtE=",
        version = "v0.0.0-20181204211112-1dc9a6cbc91a",
    )
    go_repository(
        name = "com_github_puerkitobio_purell",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/PuerkitoBio/purell",
        sum = "h1:WEQqlqaGbrPkxLJWfBwQmfEAE1Z7ONdDLqrN38tNFfI=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_puerkitobio_urlesc",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/PuerkitoBio/urlesc",
        sum = "h1:d+Bc7a5rLufV/sSk/8dngufqelfh6jnri85riMAaF/M=",
        version = "v0.0.0-20170810143723-de5bf2ad4578",
    )
    go_repository(
        name = "com_github_rogpeppe_go_internal",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/rogpeppe/go-internal",
        sum = "h1:Pu6L/AHdrTeRz8kQZd/RRxDLs4Vpl720/UfJZrjR4BQ=",
        version = "v1.0.1-alpha.3",
    )
    go_repository(
        name = "com_github_satori_go_uuid",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/satori/go.uuid",
        sum = "h1:0uYX9dsZ2yD7q2RtLRtPSdGDWzjeM3TbMJP9utgA0ww=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_serenize_snaker",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/serenize/snaker",
        sum = "h1:ofR1ZdrNSkiWcMsRrubK9tb2/SlZVWttAfqUjJi6QYc=",
        version = "v0.0.0-20171204205717-a683aaf2d516",
    )
    go_repository(
        name = "com_github_sergi_go_diff",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/sergi/go-diff",
        sum = "h1:Kpca3qRNrduNnOQeazBd0ysaKrUJiIuISHxogkT9RPQ=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_shopspring_decimal",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/shopspring/decimal",
        sum = "h1:pntxY8Ary0t43dCZ5dqY4YTJCObLY1kIXl0uzMv+7DE=",
        version = "v0.0.0-20180709203117-cd690d0c9e24",
    )
    go_repository(
        name = "com_github_shurcool_go",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/shurcooL/go",
        sum = "h1:MZM7FHLqUHYI0Y/mQAt3d2aYa0SiNms/hFqC9qJYolM=",
        version = "v0.0.0-20180423040247-9e1955d9fb6e",
    )
    go_repository(
        name = "com_github_shurcool_go_goon",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/shurcooL/go-goon",
        sum = "h1:llrF3Fs4018ePo4+G/HV/uQUqEI1HMDjCeOf2V6puPc=",
        version = "v0.0.0-20170922171312-37c2f522c041",
    )
    go_repository(
        name = "com_github_shurcool_highlight_diff",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/shurcooL/highlight_diff",
        sum = "h1:vYEG87HxbU6dXj5npkeulCS96Dtz5xg3jcfCgpcvbIw=",
        version = "v0.0.0-20170515013008-09bb4053de1b",
    )
    go_repository(
        name = "com_github_shurcool_highlight_go",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/shurcooL/highlight_go",
        sum = "h1:xLQlo0Ghg8zBaQi+tjpK+z/WLjbg/BhAWP9pYgqo/LQ=",
        version = "v0.0.0-20170515013102-78fb10f4a5f8",
    )
    go_repository(
        name = "com_github_shurcool_octicon",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/shurcooL/octicon",
        sum = "h1:j3cAp1j8k/tSLaCcDiXIpVJ8FzSJ9g1eeOAPRJYM75k=",
        version = "v0.0.0-20180602230221-c42b0e3b24d9",
    )
    go_repository(
        name = "com_github_shurcool_sanitized_anchor_name",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/shurcooL/sanitized_anchor_name",
        sum = "h1:/vdW8Cb7EXrkqWGufVMES1OH2sU9gKVb2n9/1y5NMBY=",
        version = "v0.0.0-20170918181015-86672fcb3f95",
    )
    go_repository(
        name = "com_github_sirupsen_logrus",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/sirupsen/logrus",
        sum = "h1:juTguoYk5qI21pwyTXY3B3Y5cOTH3ZUyZCg1v/mihuo=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_soheilhy_cmux",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/soheilhy/cmux",
        sum = "h1:09wy7WZk4AqO03yH85Ex1X+Uo3vDsil3Fa9AgF8Emss=",
        version = "v0.1.3",
    )
    go_repository(
        name = "com_github_sourcegraph_annotate",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/sourcegraph/annotate",
        sum = "h1:yKm7XZV6j9Ev6lojP2XaIshpT4ymkqhMeSghO5Ps00E=",
        version = "v0.0.0-20160123013949-f4cad6c6324d",
    )
    go_repository(
        name = "com_github_sourcegraph_syntaxhighlight",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/sourcegraph/syntaxhighlight",
        sum = "h1:qpG93cPwA5f7s/ZPBJnGOYQNK/vKsaDaseuKT5Asee8=",
        version = "v0.0.0-20170531221838-bd320f5d308e",
    )
    go_repository(
        name = "com_github_spf13_afero",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/spf13/afero",
        sum = "h1:5jhuqJyZCZf2JRofRvN/nIFgIWNzPa3/Vz8mYylgbWc=",
        version = "v1.2.2",
    )
    go_repository(
        name = "com_github_spf13_cast",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/spf13/cast",
        sum = "h1:oget//CVOEoFewqQxwr0Ej5yjygnqGkvggSE/gB35Q8=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_spf13_cobra",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/spf13/cobra",
        sum = "h1:ZlrZ4XsMRm04Fr5pSFxBgfND2EBVa1nLpiy1stUsX/8=",
        version = "v0.0.3",
    )
    go_repository(
        name = "com_github_spf13_jwalterweatherman",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/spf13/jwalterweatherman",
        sum = "h1:XHEdyB+EcvlqZamSM4ZOMGlc93t6AcsBEu9Gc1vn7yk=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_spf13_pflag",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/spf13/pflag",
        sum = "h1:zPAT6CGy6wXeQ7NtTnaTerfKOsV6V6F8agHXFiazDkg=",
        version = "v1.0.3",
    )
    go_repository(
        name = "com_github_spf13_viper",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/spf13/viper",
        sum = "h1:bIcUwXqLseLF3BDAZduuNfekWG87ibtFxi59Bq+oI9M=",
        version = "v1.2.1",
    )
    go_repository(
        name = "com_github_stretchr_objx",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/stretchr/objx",
        sum = "h1:Hbg2NidpLE8veEBkEZTL3CvlkUIVzuU9jDplZO54c48=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_stretchr_testify",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/stretchr/testify",
        sum = "h1:TivCn/peBQ7UY8ooIcPgZFpTNSz0Q2U6UrFlUfqbe0Q=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_tmc_grpc_websocket_proxy",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/tmc/grpc-websocket-proxy",
        sum = "h1:ndzgwNDnKIqyCvHTXaCqh9KlOWKvBry6nuXMJmonVsE=",
        version = "v0.0.0-20170815181823-89b8d40f7ca8",
    )
    go_repository(
        name = "com_github_unrolled_secure",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/unrolled/secure",
        sum = "h1:ltz/eIXkYWdMCZbu3Rb+bUmWVTm5AqM0QM8o0uKir4U=",
        version = "v0.0.0-20181005190816-ff9db2ff917f",
    )
    go_repository(
        name = "com_github_xiang90_probing",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/xiang90/probing",
        sum = "h1:MPPkRncZLN9Kh4MEFmbnK4h3BD7AUmskWv2+EeZJCCs=",
        version = "v0.0.0-20160813154853-07dd2e8dfe18",
    )
    go_repository(
        name = "com_google_cloud_go",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "cloud.google.com/go",
        sum = "h1:eOI3/cP2VTU6uZLDYAoic+eyzzB9YyGmJ7eIjl8rOPg=",
        version = "v0.34.0",
    )
    go_repository(
        name = "in_gopkg_airbrake_gobrake_v2",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "gopkg.in/airbrake/gobrake.v2",
        sum = "h1:7z2uVWwn7oVeeugY1DtlPAy5H+KYgB1KeKTnqjNatLo=",
        version = "v2.0.9",
    )
    go_repository(
        name = "in_gopkg_alexcesaro_quotedprintable_v3",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "gopkg.in/alexcesaro/quotedprintable.v3",
        sum = "h1:2gGKlE2+asNV9m7xrywl36YYNnBG5ZQ0r/BOOxqPpmk=",
        version = "v3.0.0-20150716171945-2caba252f4dc",
    )
    go_repository(
        name = "in_gopkg_check_v1",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "gopkg.in/check.v1",
        sum = "h1:qIbj1fsPNlZgppZ+VLlY7N33q108Sa+fhmuc+sWQYwY=",
        version = "v1.0.0-20180628173108-788fd7840127",
    )
    go_repository(
        name = "in_gopkg_errgo_v2",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "gopkg.in/errgo.v2",
        sum = "h1:0vLT13EuvQ0hNvakwLuFZ/jYrLp5F3kcWHXdRggjCE8=",
        version = "v2.1.0",
    )
    go_repository(
        name = "in_gopkg_fsnotify_v1",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "gopkg.in/fsnotify.v1",
        sum = "h1:xOHLXZwVvI9hhs+cLKq5+I5onOuwQLhQwiu63xxlHs4=",
        version = "v1.4.7",
    )
    go_repository(
        name = "in_gopkg_gemnasium_logrus_airbrake_hook_v2",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "gopkg.in/gemnasium/logrus-airbrake-hook.v2",
        sum = "h1:OAj3g0cR6Dx/R07QgQe8wkA9RNjB2u4i700xBkIT4e0=",
        version = "v2.1.2",
    )
    go_repository(
        name = "in_gopkg_gomail_v2",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "gopkg.in/gomail.v2",
        sum = "h1:n7WqCuqOuCbNr617RXOY0AWRXxgwEyPp2z+p0+hgMuE=",
        version = "v2.0.0-20160411212932-81ebce5c23df",
    )
    go_repository(
        name = "in_gopkg_inf_v0",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "gopkg.in/inf.v0",
        sum = "h1:3zYtXIO92bvsdS3ggAdA8Gb4Azj0YU+TVY1uGYNFA8o=",
        version = "v0.9.0",
    )
    go_repository(
        name = "in_gopkg_mail_v2",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "gopkg.in/mail.v2",
        sum = "h1:a3llQg4+Czqaf+QH4diHuHiKv4j1abMwuRXwaRNHTPU=",
        version = "v2.0.0-20180731213649-a0242b2233b4",
    )
    go_repository(
        name = "in_gopkg_natefinch_lumberjack_v2",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "gopkg.in/natefinch/lumberjack.v2",
        sum = "h1:986b60BAz5vO2Vaf48yQaq+wb2bU4JsXxKu1+itW6x8=",
        version = "v2.0.0-20150622162204-20b71e5b60d7",
    )
    go_repository(
        name = "in_gopkg_square_go_jose_v2",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "gopkg.in/square/go-jose.v2",
        sum = "h1:ELQJ5WuT+ydETLCpWvAuw8iGBQRGoJq+A3RAbbAcZUY=",
        version = "v2.0.0-20180411045311-89060dee6a84",
    )
    go_repository(
        name = "in_gopkg_tomb_v1",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "gopkg.in/tomb.v1",
        sum = "h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=",
        version = "v1.0.0-20141024135613-dd632973f1e7",
    )
    go_repository(
        name = "in_gopkg_yaml_v1",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "gopkg.in/yaml.v1",
        sum = "h1:POO/ycCATvegFmVuPpQzZFJ+pGZeX22Ufu6fibxDVjU=",
        version = "v1.0.0-20140924161607-9f9df34309c0",
    )
    go_repository(
        name = "in_gopkg_yaml_v2",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "gopkg.in/yaml.v2",
        sum = "h1:ZCJp+EgiOT7lHqUV2J862kp8Qj64Jo6az82+3Td9dZw=",
        version = "v2.2.2",
    )
    go_repository(
        name = "io_k8s_api",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "k8s.io/api",
        sum = "h1:7Gz7/nQ7X2qmPXMyN0bNq7Zm9Uip+UnFuMZTd2l3vms=",
        version = "v0.0.0-20190819141258-3544db3b9e44",
    )
    go_repository(
        name = "io_k8s_apimachinery",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "k8s.io/apimachinery",
        sum = "h1:7Kns6qqhMAQWvGkxYOLSLRZ5hJO0/5pcE5lPGP2fxUw=",
        version = "v0.0.0-20190817020851-f2f3a405f61d",
    )
    go_repository(
        name = "io_k8s_apiserver",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "k8s.io/apiserver",
        sum = "h1:xgPuryD+iPG1FdQ0R6lKYvMwLnhPDi/1WDOMz4Y88bw=",
        version = "v0.0.0-20190819142446-92cc630367d0",
    )
    go_repository(
        name = "io_k8s_client_go",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "k8s.io/client-go",
        sum = "h1:w1BoabVnPpPqQCY3sHK4qVwa12Lk8ip1pKMR1C+qbdo=",
        version = "v0.0.0-20190819141724-e14f31a72a77",
    )
    go_repository(
        name = "io_k8s_component_base",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "k8s.io/component-base",
        sum = "h1:4C6bgyEgzfGDQEkyq/swmBelfEIH494iHGdZPUF0KO8=",
        version = "v0.0.0-20190819141909-f0f7c184477d",
    )
    go_repository(
        name = "io_k8s_gengo",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "k8s.io/gengo",
        sum = "h1:j4L8spMe0tFfBvvW6lrc0c+Ql8+nnkcV3RYfi3eSwGY=",
        version = "v0.0.0-20190907103519-ebc107f98eab",
    )
    go_repository(
        name = "io_k8s_klog",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "k8s.io/klog",
        sum = "h1:RVgyDHY/kFKtLqh67NvEWIgkMneNoIrdkN0CxDSQc68=",
        version = "v0.3.1",
    )
    go_repository(
        name = "io_k8s_kube_openapi",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "k8s.io/kube-openapi",
        sum = "h1:TRb4wNWoBVrH9plmkp2q86FIDppkbrEXdXlxU3a3BMI=",
        version = "v0.0.0-20190228160746-b3a7cee44a30",
    )
    go_repository(
        name = "io_k8s_sigs_controller_tools",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "sigs.k8s.io/controller-tools",
        sum = "h1:onJbaVxkLKbLhy3BcLdesjkP6y+6WFVzg1BJv28BpCo=",
        version = "v0.1.10",
    )
    go_repository(
        name = "io_k8s_sigs_kubebuilder",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "sigs.k8s.io/kubebuilder",
        sum = "h1:WBb/IUMWIhQbioR0dCgckr/DvF6LnS0DtVNU+psVX24=",
        version = "v0.0.0-20190320190143-2621a6fdb324",
    )
    go_repository(
        name = "io_k8s_sigs_structured_merge_diff",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "sigs.k8s.io/structured-merge-diff",
        sum = "h1:9r5DY45ef9LtcA6BnkhW8MPV7OKAfbf2AUwUhq3LeRk=",
        version = "v0.0.0-20190302045857-e85c7b244fd2",
    )
    go_repository(
        name = "io_k8s_sigs_yaml",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "sigs.k8s.io/yaml",
        sum = "h1:4A07+ZFc2wgJwo8YNlQpr1rVlgUDlxXHhPJciaPY5gs=",
        version = "v1.1.0",
    )
    go_repository(
        name = "io_k8s_utils",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "k8s.io/utils",
        sum = "h1:ElyM7RPonbKnQqOcw7dG2IK5uvQQn3b/WPHqD5mBvP4=",
        version = "v0.0.0-20190221042446-c2654d5206da",
    )
    go_repository(
        name = "org_golang_google_appengine",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "google.golang.org/appengine",
        sum = "h1:KxkO13IPW4Lslp2bz+KHP2E3gtFlrIGNThxkZQ3g+4c=",
        version = "v1.5.0",
    )
    go_repository(
        name = "org_golang_google_genproto",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "google.golang.org/genproto",
        sum = "h1:72GtwBPfq6av9X0Ru2HtAopsPW+d+vh1K1zaxanTdE8=",
        version = "v0.0.0-20170731182057-09f6ed296fc6",
    )
    go_repository(
        name = "org_golang_google_grpc",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "google.golang.org/grpc",
        sum = "h1:bHIbVsCwmvbArgCJmLdgOdHFXlKqTOVjbibbS19cXHc=",
        version = "v1.13.0",
    )
    go_repository(
        name = "org_golang_x_crypto",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "golang.org/x/crypto",
        sum = "h1:1wopBVtVdWnn03fZelqdXTqk7U7zPQCb+T4rbU9ZEoU=",
        version = "v0.0.0-20190611184440-5c40567a22f8",
    )
    go_repository(
        name = "org_golang_x_net",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "golang.org/x/net",
        sum = "h1:k7pJ2yAPLPgbskkFdhRCsA77k2fySZ1zf2zCjvQCiIM=",
        version = "v0.0.0-20190827160401-ba9fcec4b297",
    )
    go_repository(
        name = "org_golang_x_oauth2",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "golang.org/x/oauth2",
        sum = "h1:tImsplftrFpALCYumobsd0K86vlAs/eXGFms2txfJfA=",
        version = "v0.0.0-20190402181905-9f3314589c9a",
    )
    go_repository(
        name = "org_golang_x_sync",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "golang.org/x/sync",
        sum = "h1:8gQV6CLnAEikrhgkHFbMAEhagSSnXWGV915qUMm9mrU=",
        version = "v0.0.0-20190423024810-112230192c58",
    )
    go_repository(
        name = "org_golang_x_sys",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "golang.org/x/sys",
        sum = "h1:25KHgbfyiSm6vwQLbM3zZIe1v9p/3ea4Rz+nnM5K/i4=",
        version = "v0.0.0-20190616124812-15dcb6c0061f",
    )
    go_repository(
        name = "org_golang_x_text",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "golang.org/x/text",
        sum = "h1:tW2bmiBqwgJj/UpqtC8EpXEZVYOwU0yG4iWbprSVAcs=",
        version = "v0.3.2",
    )
    go_repository(
        name = "org_golang_x_time",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "golang.org/x/time",
        sum = "h1:TnM+PKb3ylGmZvyPXmo9m/wktg7Jn/a/fNmr33HSj8g=",
        version = "v0.0.0-20161028155119-f51c12702a4d",
    )
    go_repository(
        name = "org_golang_x_tools",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "golang.org/x/tools",
        sum = "h1:QjA/9ArTfVTLfEhClDCG7SGrZkZixxWpwNCDiwJfh88=",
        version = "v0.0.0-20190614205625-5aca471b1d59",
    )
    go_repository(
        name = "org_uber_go_atomic",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "go.uber.org/atomic",
        sum = "h1:nSQar3Y0E3VQF/VdZ8PTAilaXpER+d7ypdABCrpwMdg=",
        version = "v0.0.0-20181018215023-8dc6146f7569",
    )
    go_repository(
        name = "org_uber_go_multierr",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "go.uber.org/multierr",
        sum = "h1:shvkWr0NAZkg4nPuE3XrKP0VuBPijjk3TfX6Y6acFNg=",
        version = "v0.0.0-20180122172545-ddea229ff1df",
    )
    go_repository(
        name = "org_uber_go_zap",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "go.uber.org/zap",
        sum = "h1:Z2sc4+v0JHV6Mn4kX1f2a5nruNjmV+Th32sugE8zwz8=",
        version = "v0.0.0-20180814183419-67bc79d13d15",
    )
    go_repository(
        name = "tools_gotest",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "gotest.tools",
        sum = "h1:VsBPFP1AI068pPrMxtb/S8Zkgf9xEmTLJjfM+P5UIEo=",
        version = "v2.2.0+incompatible",
    )
    go_repository(
        name = "com_github_golang_glog",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/golang/glog",
        sum = "h1:VKtxabqXZkF25pY9ekfRL6a582T4P37/31XEstQ5p58=",
        version = "v0.0.0-20160126235308-23def4e6c14b",
    )
    go_repository(
        name = "com_github_burntsushi_xgb",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/BurntSushi/xgb",
        sum = "h1:1BDTz0u9nC3//pOCMdNH+CiXJVYJh5UQNCOBG7jbELc=",
        version = "v0.0.0-20160522181843-27f122750802",
    )
    go_repository(
        name = "com_github_remyoudompheng_bigfft",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/remyoudompheng/bigfft",
        sum = "h1:/NRJ5vAYoqz+7sG51ubIDHXeWO8DlTSrToPu6q11ziA=",
        version = "v0.0.0-20170806203942-52369c62f446",
    )
    go_repository(
        name = "io_k8s_code_generator",
        build_extra_args = [
            "--known_import=k8s.io/kubernetes",
        ],
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "k8s.io/code-generator",
        sum = "h1:24rP1wJ96ap3+lAZZTmJMJ0U783FLox4irw17QoxhDg=",
        version = "v0.15.8-beta.0",
    )
    go_repository(
        name = "org_golang_x_exp",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "golang.org/x/exp",
        sum = "h1:I6A9Ag9FpEKOjcKrRNjQkPHawoXIhKyTGfvvjFAiiAk=",
        version = "v0.0.0-20190312203227-4b39c73a6495",
    )
    go_repository(
        name = "org_golang_x_image",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "golang.org/x/image",
        sum = "h1:KYGJGHOQy8oSi1fDlSpcZF0+juKwk/hEMv5SiwHogR0=",
        version = "v0.0.0-20190227222117-0694c2d4d067",
    )
    go_repository(
        name = "org_golang_x_mobile",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "golang.org/x/mobile",
        sum = "h1:Tus/Y4w3V77xDsGwKUC8a/QrV7jScpU557J77lFffNs=",
        version = "v0.0.0-20190312151609-d3739f865fa6",
    )
    go_repository(
        name = "org_gonum_v1_gonum",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "gonum.org/v1/gonum",
        sum = "h1:OB/uP/Puiu5vS5QMRPrXCDWUPb+kt8f1KW8oQzFejQw=",
        version = "v0.0.0-20190331200053-3d26580ed485",
    )
    go_repository(
        name = "org_gonum_v1_netlib",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "gonum.org/v1/netlib",
        sum = "h1:jRyg0XfpwWlhEV8mDfdNGBeSJM2fuyh9Yjrnd8kF2Ts=",
        version = "v0.0.0-20190331212654-76723241ea4e",
    )
    go_repository(
        name = "org_modernc_cc",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "modernc.org/cc",
        sum = "h1:nPibNuDEx6tvYrUAtvDTTw98rx5juGsa5zuDnKwEEQQ=",
        version = "v1.0.0",
    )
    go_repository(
        name = "org_modernc_golex",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "modernc.org/golex",
        sum = "h1:wWpDlbK8ejRfSyi0frMyhilD3JBvtcx2AdGDnU+JtsE=",
        version = "v1.0.0",
    )
    go_repository(
        name = "org_modernc_mathutil",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "modernc.org/mathutil",
        sum = "h1:93vKjrJopTPrtTNpZ8XIovER7iCIH1QU7wNbOQXC60I=",
        version = "v1.0.0",
    )
    go_repository(
        name = "org_modernc_strutil",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "modernc.org/strutil",
        sum = "h1:XVFtQwFVwc02Wk+0L/Z/zDDXO81r5Lhe6iMKmGX3KhE=",
        version = "v1.0.0",
    )
    go_repository(
        name = "org_modernc_xc",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "modernc.org/xc",
        sum = "h1:7ccXrupWZIS3twbUGrtKmHS2DXY6xegFua+6O3xgAFU=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_kubernetes_incubator_reference_docs",
        build_extra_args = [
            "--known_import=k8s.io/kubernetes",
            "--exclude=gen-kubectldocs",
        ],
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/kubernetes-incubator/reference-docs",
        replace = "github.com/kubernetes-sigs/reference-docs",
        sum = "h1:5H0rNjhi9XVX98G0hC2ELjYJFbtBcfwQ61nsaop+6d0=",
        version = "v0.0.0-20170929004150-fcf65347b256",
    )
    go_repository(
        name = "com_github_asaskevich_govalidator",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/asaskevich/govalidator",
        sum = "h1:idn718Q4B6AGu/h5Sxe66HYVdqdGu2l9Iebqhi/AEoA=",
        version = "v0.0.0-20190424111038-f61b66f89f4a",
    )
    go_repository(
        name = "com_github_globalsign_mgo",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/globalsign/mgo",
        sum = "h1:D4uzjWwKYQ5XnAvUbuvHW93esHg7F8N/OYeBBcJoTr0=",
        version = "v0.0.0-20180905125535-1ca0a4f7cbcb",
    )
    go_repository(
        name = "com_github_go_openapi_analysis",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/go-openapi/analysis",
        sum = "h1:OcMMVVJBRiSsAkXQwSVL4sRrVaqeqPnOcy1Kw0JI47w=",
        version = "v0.19.7",
    )
    go_repository(
        name = "com_github_go_openapi_errors",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/go-openapi/errors",
        sum = "h1:a2kIyV3w+OS3S97zxUndRVD46+FhGOUBDFY7nmu4CsY=",
        version = "v0.19.2",
    )
    go_repository(
        name = "com_github_go_openapi_loads",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/go-openapi/loads",
        sum = "h1:wCOBNscACI8L93tt5tvB2zOMkJ098XCw3fP0BY2ybDA=",
        version = "v0.19.0",
    )
    go_repository(
        name = "com_github_go_openapi_strfmt",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/go-openapi/strfmt",
        sum = "h1:eRfyY5SkaNJCAwmmMcADjY31ow9+N7MCLW7oRkbsINA=",
        version = "v0.19.3",
    )
    go_repository(
        name = "com_github_go_stack_stack",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/go-stack/stack",
        sum = "h1:5SgMzNM5HxrEjV0ww2lTmX6E2Izsfxas4+YHWRs3Lsk=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_github_tidwall_pretty",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "github.com/tidwall/pretty",
        sum = "h1:HsD+QiTn7sK6flMKIvNmpqz1qrpP3Ps6jOKIKMooyg4=",
        version = "v1.0.0",
    )
    go_repository(
        name = "org_mongodb_go_mongo_driver",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
        importpath = "go.mongodb.org/mongo-driver",
        sum = "h1:Sq1fR+0c58RME5EoqKdjkiQAmPjmfHlZOoRI6fTUOcs=",
        version = "v1.1.1",
    )
