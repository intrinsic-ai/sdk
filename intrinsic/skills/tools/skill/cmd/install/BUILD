# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "install",
    srcs = ["install.go"],
    deps = [
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:imagetransfer",
        "//intrinsic/assets/proto:installed_assets_go_grpc_proto",
        "//intrinsic/assets/services:bundleimages",
        "//intrinsic/skills:skillbundle",
        "//intrinsic/skills/tools/skill/cmd:waitforskill",
        "//intrinsic/skills/tools/skill/cmd/directupload",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
