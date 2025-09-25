# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

# Code for working with Data Assets.

go_library(
    name = "utils",
    srcs = ["utils.go"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/data/proto/v1:referenced_data_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/util/proto:registryutil",
        "//intrinsic/util/proto:walkmessages",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)

go_library(
    name = "datamanifest",
    srcs = ["datamanifest.go"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:metadatautils",
        "//intrinsic/assets/data/proto/v1:data_manifest_go_proto",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
    ],
)

go_test(
    name = "datamanifest_test",
    srcs = ["datamanifest_test.go"],
    library = ":datamanifest",
    deps = [
        "//intrinsic/assets/data/proto/v1:data_manifest_go_proto",
        "//intrinsic/assets/data/proto/v1:referenced_data_struct_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/util/proto:descriptor",
        "//intrinsic/util/proto:registryutil",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)
