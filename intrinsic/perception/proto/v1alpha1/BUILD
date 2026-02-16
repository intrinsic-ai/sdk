# Copyright 2023 Intrinsic Innovation LLC

# Automatic symmetry detection.
# go/intrinsic-tdd-symmetry-detection

load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//intrinsic:public_api_users"])

proto_library(
    name = "symmetry_service_proto",
    srcs = ["symmetry_service.proto"],
    deps = [
        "//intrinsic/geometry/proto/v1:transformed_geometry_proto",
        "//intrinsic/perception/proto/v1:symmetry_proto",
        "//intrinsic/scene/proto/v1:scene_object_proto",
        "@com_google_googleapis//google/longrunning:operations_proto",
        "@com_google_protobuf//:empty_proto",
    ],
)

cc_proto_library(
    name = "symmetry_service_cc_proto",
    deps = [":symmetry_service_proto"],
)

go_proto_library(
    name = "symmetry_service_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/perception/proto/v1alpha1/symmetry_service_go_proto",
    protos = [":symmetry_service_proto"],
    visibility = ["//intrinsic/frontend/go:__subpackages__"],
    deps = [
        "//intrinsic/geometry/proto/v1:transformed_geometry_go_proto",
        "//intrinsic/perception/proto/v1:symmetry_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_protobuf//types/known/emptypb",
    ],
)

cc_grpc_library(
    name = "symmetry_service_cc_grpc",
    srcs = [":symmetry_service_proto"],
    grpc_only = True,
    deps = [":symmetry_service_cc_proto"],
)
