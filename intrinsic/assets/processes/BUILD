# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//intrinsic:public_api_users"])

go_library(
    name = "processutil",
    srcs = ["processutil.go"],
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

go_test(
    name = "processutil_test",
    srcs = ["processutil_test.go"],
    library = ":processutil",
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:metadatautils",
        "//intrinsic/assets/processes/proto:process_asset_go_proto",
        "//intrinsic/assets/processes/proto:process_manifest_go_proto",
        "//intrinsic/assets/proto:asset_tag_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:documentation_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/executive/proto:behavior_tree_go_proto",
        "//intrinsic/skills/proto:equipment_go_proto",
        "//intrinsic/skills/proto:skills_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)
