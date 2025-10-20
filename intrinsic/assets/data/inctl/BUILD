# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "data",
    srcs = ["data.go"],
    visibility = ["//intrinsic/tools/inctl:__pkg__"],
    deps = [
        ":release",
        "//intrinsic/assets/inctl:install",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:cobrautil",
    ],
)

go_library(
    name = "release",
    srcs = ["release.go"],
    deps = [
        "//intrinsic/assets:bundleio",
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets:idutils",
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_grpc_proto",
        "//intrinsic/assets/catalog/proto/v1:release_metadata_go_proto",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_spf13_cobra//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
