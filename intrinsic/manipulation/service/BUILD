# Copyright 2023 Intrinsic Innovation LLC

# Services for the manipulation package.

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")

package(
    default_visibility = [
        "//intrinsic/logging:__subpackages__",
        "//intrinsic/manipulation:__subpackages__",
    ],
)

proto_library(
    name = "grasp_annotator_service_proto",
    srcs = ["grasp_annotator_service.proto"],
    deps = [
        "//intrinsic/manipulation/grasping:grasp_annotations_proto",
        "//intrinsic/manipulation/grasping:grasp_annotator_proto",
    ],
)

py_proto_library(
    name = "grasp_annotator_service_py_pb2",
    deps = [":grasp_annotator_service_proto"],
)

py_grpc_library(
    name = "grasp_annotator_service_py_pb2_grpc",
    srcs = [":grasp_annotator_service_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":grasp_annotator_service_py_pb2"],
)

proto_library(
    name = "multi_axis_auto_annotator_service_proto",
    srcs = ["multi_axis_auto_annotator_service.proto"],
    deps = [
        "//intrinsic/geometry/proto:triangle_mesh_proto",
        "//intrinsic/manipulation/grasping:schmalz_multi_axis_grasp_proto",
    ],
)

py_proto_library(
    name = "multi_axis_auto_annotator_service_py_pb2",
    deps = [":multi_axis_auto_annotator_service_proto"],
)

py_grpc_library(
    name = "multi_axis_auto_annotator_service_py_pb2_grpc",
    srcs = [":multi_axis_auto_annotator_service_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":multi_axis_auto_annotator_service_py_pb2"],
)
