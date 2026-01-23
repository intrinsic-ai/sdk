# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_python//python:defs.bzl", "py_library")
load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "fakedataassets",
    testonly = True,
    srcs = ["fake_data_assets.go"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/data/proto/v1:data_assets_go_grpc_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/testing:grpctest",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//credentials/local:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

cc_library(
    name = "fake_data_assets_cc",
    testonly = True,
    srcs = ["fake_data_assets.cc"],
    hdrs = ["fake_data_assets.h"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets:id_utils",
        "//intrinsic/assets/data/proto/v1:data_assets_cc_grpc_proto",
        "//intrinsic/assets/data/proto/v1:data_assets_cc_proto",
        "//intrinsic/util/status:ret_check",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/status:status_macros_grpc",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/types:span",
        "@com_github_grpc_grpc//:grpc++",
    ],
)

py_library(
    name = "fake_data_assets_py",
    testonly = True,
    srcs = ["fake_data_assets.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets:id_utils_py",
        "//intrinsic/assets/data/proto/v1:data_asset_py_pb2",
        "//intrinsic/assets/data/proto/v1:data_assets_py_pb2",
        "//intrinsic/assets/data/proto/v1:data_assets_py_pb2_grpc",
        requirement("grpcio"),
    ],
)

go_library(
    name = "databundle",
    srcs = ["databundle.go"],
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        ":datavalidate",
        ":utils",
        "//intrinsic/assets:ioutils",
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_grpc_proto",
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/data/proto/v1:referenced_data_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/util/archive:tartooling",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_google_safearchive//tar",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)

go_library(
    name = "datavalidate",
    srcs = ["datavalidate.go"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":utils",
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:metadatautils",
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/data/proto/v1:data_manifest_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/util/proto:registryutil",
        "@org_golang_google_protobuf//reflect/protodesc:go_default_library",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
    ],
)

go_library(
    name = "utils",
    srcs = ["utils.go"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/data/proto/v1:referenced_data_go_proto",
        "//intrinsic/util/proto:registryutil",
        "//intrinsic/util/proto:walkmessages",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)
