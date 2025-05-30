# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@rules_python//python:defs.bzl", "py_library", "py_test")

package(default_visibility = ["//visibility:public"])

py_library(
    name = "object_world_resources",
    srcs = ["object_world_resources.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":object_world_ids",
        "//intrinsic/icon/proto:cart_space_py_pb2",
        "//intrinsic/kinematics/types:joint_limits_py_pb2",
        "//intrinsic/math/python:data_types",
        "//intrinsic/math/python:proto_conversion",
        "//intrinsic/util/grpc:error_handling",
        "//intrinsic/world/proto:object_world_refs_py_pb2",
        "//intrinsic/world/proto:object_world_service_py_pb2",
        "//intrinsic/world/proto:object_world_service_py_pb2_grpc",
        "//intrinsic/world/proto:object_world_updates_py_pb2",
        "//intrinsic/world/robot_payload/python:robot_payload",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_test(
    name = "object_world_resources_test",
    srcs = ["object_world_resources_test.py"],
    deps = [
        ":object_world_ids",
        ":object_world_resources",
        "//intrinsic/icon/proto:cart_space_py_pb2",
        "//intrinsic/kinematics/types:joint_limits_py_pb2",
        "//intrinsic/math/python:data_types",
        "//intrinsic/solutions/testing:compare",
        "//intrinsic/world/proto:object_world_refs_py_pb2",
        "//intrinsic/world/proto:object_world_service_py_pb2",
        "@ai_intrinsic_sdks_pip_deps//numpy",
        requirement("grpcio"),
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_library(
    name = "object_world_ids",
    srcs = ["object_world_ids.py"],
    visibility = ["//intrinsic:public_api_users"],
)

py_test(
    name = "object_world_ids_test",
    srcs = ["object_world_ids_test.py"],
    deps = [
        ":object_world_ids",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

py_library(
    name = "object_world_client",
    srcs = ["object_world_client.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":object_world_ids",
        ":object_world_resources",
        "//intrinsic/geometry/proto:geometry_service_py_pb2",
        "//intrinsic/geometry/proto:geometry_service_py_pb2_grpc",
        "//intrinsic/geometry/proto:geometry_storage_refs_py_pb2",
        "//intrinsic/icon/equipment:icon_equipment_py_pb2",
        "//intrinsic/icon/proto:cart_space_py_pb2",
        "//intrinsic/kinematics/types:joint_limits_py_pb2",
        "//intrinsic/math/python:data_types",
        "//intrinsic/math/python:proto_conversion",
        "//intrinsic/resources/proto:resource_handle_py_pb2",
        "//intrinsic/scene/product/proto:product_world_object_data_py_pb2",
        "//intrinsic/scene/proto/v1:scene_object_py_pb2",
        "//intrinsic/util/grpc:error_handling",
        "//intrinsic/world/proto:geometry_component_py_pb2",
        "//intrinsic/world/proto:object_world_refs_py_pb2",
        "//intrinsic/world/proto:object_world_service_py_pb2",
        "//intrinsic/world/proto:object_world_service_py_pb2_grpc",
        "//intrinsic/world/proto:object_world_updates_py_pb2",
        "//intrinsic/world/robot_payload/python:robot_payload",
        requirement("grpcio"),
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_test(
    name = "object_world_client_external_test",
    srcs = ["object_world_client_external_test.py"],
    deps = [
        ":object_world_client",
        ":object_world_ids",
        "//intrinsic/scene/proto/v1:scene_object_py_pb2",
        "//intrinsic/world/proto:geometry_component_py_pb2",
        "//intrinsic/world/proto:object_world_service_py_pb2",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_protobuf//:protobuf_python",
    ],
)
