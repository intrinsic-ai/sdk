# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "data",
    srcs = ["data.go"],
    visibility = ["//intrinsic/tools/inctl:__pkg__"],
    deps = [
        ":install",
        ":release",
        "//intrinsic/tools/inctl/cmd:root",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_library(
    name = "install",
    srcs = ["install.go"],
    deps = [
        "//intrinsic/assets:bundleio",
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets:idutils",
        "//intrinsic/assets/proto:installed_assets_go_grpc_proto",
        "@com_github_pkg_errors//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_grpc//status:go_default_library",
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
