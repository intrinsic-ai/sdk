# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library", "py_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

package(default_visibility = [
    "//intrinsic/assets:__subpackages__",
    "//intrinsic/frontend:__subpackages__",
    "//intrinsic/kubernetes:__subpackages__",
    "//intrinsic/skills/tools/skill:__subpackages__",
    "//intrinsic/solutions:__subpackages__",
    "//intrinsic/tools/inctl:__subpackages__",
])

proto_library(
    name = "clusterdiscovery_api_proto",
    srcs = ["clusterdiscovery_api.proto"],
)

go_proto_library(
    name = "clusterdiscovery_api_go_proto",
    deps = [":clusterdiscovery_api_proto"],
)

go_grpc_library(
    name = "clusterdiscovery_api_go_grpc_proto",
    srcs = [":clusterdiscovery_api_proto"],
    deps = [":clusterdiscovery_api_go_proto"],
)

proto_library(
    name = "clustermanager_api_proto",
    srcs = ["clustermanager_api.proto"],
    deps = [
        "//intrinsic/util/status:extended_status_proto",
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-openapiv2/options:options_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/api:field_behavior_proto",
        "@com_google_googleapis//google/longrunning:operations_proto",
        "@com_google_protobuf//:duration_proto",
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:field_mask_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

go_proto_library(
    name = "clustermanager_api_go_proto",
    go_deps = [
        "@org_golang_google_genproto_googleapis_api//annotations",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "//intrinsic/util/status:extended_status_go_proto",
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-openapiv2/options:options",
    ],
    deps = [":clustermanager_api_proto"],
)

go_grpc_library(
    name = "clustermanager_api_go_grpc_proto",
    srcs = [":clustermanager_api_proto"],
    deps = [
        "//intrinsic/util/status:extended_status_go_proto",
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-openapiv2/options",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)

proto_library(
    name = "solutiondiscovery_api_proto",
    srcs = ["solutiondiscovery_api.proto"],
    deps = ["//intrinsic/frontend/cloud/api/v1:clusterdiscovery_api_proto"],
)

py_proto_library(
    name = "solutiondiscovery_api_py_pb2",
    deps = [":solutiondiscovery_api_proto"],
)

py_grpc_library(
    name = "solutiondiscovery_api_py_pb2_grpc",
    srcs = [":solutiondiscovery_api_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":solutiondiscovery_api_py_pb2"],
)

go_proto_library(
    name = "solutiondiscovery_api_go_proto",
    go_deps = ["//intrinsic/frontend/cloud/api/v1:clusterdiscovery_api_go_grpc_proto"],
    deps = [
        ":solutiondiscovery_api_proto",
    ],
)

go_grpc_library(
    name = "solutiondiscovery_api_go_grpc_proto",
    srcs = [":solutiondiscovery_api_proto"],
    deps = [
        ":solutiondiscovery_api_go_proto",
        "//intrinsic/frontend/cloud/api/v1:clusterdiscovery_api_go_grpc_proto",
    ],
)
