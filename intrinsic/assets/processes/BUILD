# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//intrinsic:public_api_users"])

go_library(
    name = "processbundle",
    srcs = ["processbundle.go"],
    importpath = "intrinsic/assets/processes/processbundle",
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        ":processmanifest",
        ":processvalidate",
        "//intrinsic/assets:ioutils",
        "//intrinsic/assets/processes/proto:process_asset_go_proto",
        "//intrinsic/assets/processes/proto:process_manifest_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/util/archive:tartooling",
        "@com_github_google_safearchive//tar",
    ],
)

go_library(
    name = "processmanifest",
    srcs = ["processmanifest.go"],
    importpath = "intrinsic/assets/processes/processmanifest",
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/executive/proto:behavior_tree_go_proto",
        "//intrinsic/skills/proto:skills_go_proto",
    ],
)

go_library(
    name = "processvalidate",
    srcs = ["processvalidate.go"],
    importpath = "intrinsic/assets/processes/processvalidate",
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:metadatautils",
        "//intrinsic/assets/processes/proto:process_asset_go_proto",
        "//intrinsic/assets/processes/proto:process_manifest_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:documentation_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/executive/proto:behavior_tree_go_proto",
        "//intrinsic/skills/proto:skills_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
