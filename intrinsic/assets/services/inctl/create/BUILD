# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

go_library(
    name = "create",
    srcs = ["create.go"],
    embedsrcs = [
        "templates/BUILD_cc.template",  #keep
        "templates/BUILD_py.template",  #keep
        "templates/service_cc.template",  #keep
        "templates/service_manifest.template",  #keep
        "templates/config_proto.template",  #keep
        "templates/service_py.template",  #keep
        "templates/defaults.textproto.template",  #keep
    ],
    importpath = "intrinsic/assets/services/inctl/create/create",
    visibility = ["//intrinsic/assets/services/inctl:__pkg__"],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:buildozer",
        "//intrinsic/tools/inctl/util:printer",
        "//intrinsic/tools/inctl/util:templateutil",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_test(
    name = "create_test",
    srcs = ["create_test.go"],
    embed = [":create"],
    importpath = "intrinsic/assets/services/inctl/create/create_test",
    deps = [
        "@com_github_google_go_cmp//cmp:go_default_library",
    ],
)
