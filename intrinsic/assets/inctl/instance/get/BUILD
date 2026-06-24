# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "get",
    srcs = ["get.go"],
    importpath = "intrinsic/assets/inctl/instance/get",
    visibility = ["//intrinsic/assets/inctl/instance:__subpackages__"],
    deps = [
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets:typeutils",
        "//intrinsic/assets/inctl/instance/common",
        "//intrinsic/assets/proto/v1:asset_instances_go_proto",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:printer",
        "//intrinsic/util/proto:protoio",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)
