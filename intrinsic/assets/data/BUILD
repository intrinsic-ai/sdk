# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

# Code for working with the Data asset type.

go_library(
    name = "dataassetsservice",
    srcs = ["data_assets_service.go"],
    visibility = ["//intrinsic/kubernetes/workcell_spec:__pkg__"],
    deps = [
        ":utils",
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:pageutils",
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/data/proto/v1:data_assets_go_grpc_proto",
        "//intrinsic/assets/data/proto/v1:referenced_data_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/resources/service:resourcetyperuntime",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
    ],
)

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
        "@com_github_minio_highwayhash//:go_default_library",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)
