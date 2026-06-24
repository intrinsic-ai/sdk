# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "list",
    srcs = ["list.go"],
    importpath = "intrinsic/assets/inctl/instance/list",
    visibility = ["//intrinsic/assets/inctl/instance:__subpackages__"],
    deps = [
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:tagutils",
        "//intrinsic/assets:typeutils",
        "//intrinsic/assets/inctl/instance/common",
        "//intrinsic/assets/proto:asset_tag_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto/v1:asset_instances_go_proto",
        "//intrinsic/assets/proto/v1:dependency_go_proto",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)
