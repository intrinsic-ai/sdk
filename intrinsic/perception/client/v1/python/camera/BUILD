# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@rules_python//python:defs.bzl", "py_library")

py_library(
    name = "_camera_utils",
    srcs = ["_camera_utils.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/perception/proto:camera_config_py_pb2",
        "//intrinsic/perception/proto:camera_identifier_py_pb2",
        "//intrinsic/perception/proto:camera_params_py_pb2",
        "//intrinsic/perception/proto:camera_settings_py_pb2",
        "//intrinsic/perception/proto:dimensions_py_pb2",
        "//intrinsic/perception/proto:distortion_params_py_pb2",
        "//intrinsic/perception/proto:intrinsic_params_py_pb2",
        "//intrinsic/perception/proto:sensor_config_py_pb2",
        "//intrinsic/perception/proto/v1:camera_config_py_pb2",
        "//intrinsic/perception/proto/v1:camera_identifier_py_pb2",
        "//intrinsic/perception/proto/v1:camera_params_py_pb2",
        "//intrinsic/perception/proto/v1:camera_settings_py_pb2",
        "//intrinsic/perception/proto/v1:dimensions_py_pb2",
        "//intrinsic/perception/proto/v1:distortion_params_py_pb2",
        "//intrinsic/perception/proto/v1:intrinsic_params_py_pb2",
        "//intrinsic/perception/proto/v1:sensor_config_py_pb2",
        "//intrinsic/resources/proto:resource_handle_py_pb2",
        "//intrinsic/skills/python:proto_utils",
        "@ai_intrinsic_sdks_pip_deps//numpy",
        requirement("grpcio"),
        "@com_google_absl_py//absl/logging",
    ],
)

py_library(
    name = "data_classes",
    srcs = ["data_classes.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":_camera_utils",
        "//intrinsic/math/python:pose3",
        "//intrinsic/math/python:proto_conversion",
        "//intrinsic/perception/client/v1/python:image_utils",
        "//intrinsic/perception/proto/v1:camera_config_py_pb2",
        "//intrinsic/perception/proto/v1:camera_params_py_pb2",
        "//intrinsic/perception/proto/v1:camera_server_py_pb2",
        "//intrinsic/perception/proto/v1:capture_result_py_pb2",
        "//intrinsic/perception/proto/v1:sensor_config_py_pb2",
        "//intrinsic/perception/proto/v1:sensor_image_py_pb2",
        "@ai_intrinsic_sdks_pip_deps//numpy",
    ],
)

py_library(
    name = "camera_client",
    srcs = ["camera_client.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/perception/proto/v1:camera_config_py_pb2",
        "//intrinsic/perception/proto/v1:camera_identifier_py_pb2",
        "//intrinsic/perception/proto/v1:camera_server_py_pb2",
        "//intrinsic/perception/proto/v1:camera_server_py_pb2_grpc",
        "//intrinsic/perception/proto/v1:camera_settings_py_pb2",
        "//intrinsic/perception/proto/v1:capture_result_py_pb2",
        "//intrinsic/util/grpc:connection",
        "//intrinsic/util/grpc:error_handling",
        "//intrinsic/util/grpc:interceptor",
        requirement("grpcio"),
    ],
)

py_library(
    name = "cameras",
    srcs = ["cameras.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":_camera_utils",
        ":camera_client",
        ":data_classes",
        "//intrinsic/math/python:pose3",
        "//intrinsic/perception/proto/v1:settings_py_pb2",
        "//intrinsic/resources/client:resource_registry_client",
        "//intrinsic/resources/proto:resource_handle_py_pb2",
        "//intrinsic/skills/proto:equipment_py_pb2",
        "//intrinsic/skills/python:skill_interface",
        "//intrinsic/util/grpc:connection",
        "//intrinsic/world/python:object_world_client",
        "//intrinsic/world/python:object_world_resources",
        "@ai_intrinsic_sdks_pip_deps//numpy",
        requirement("grpcio"),
        "@com_google_absl_py//absl/logging",
        "@com_google_protobuf//:protobuf_python",
    ],
)
