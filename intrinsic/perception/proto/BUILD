# Copyright 2023 Intrinsic Innovation LLC

# This package contains public protos that are exposed to clients of Intrinsic perception.

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

proto_library(
    name = "camera_params_proto",
    srcs = ["camera_params.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":distortion_params_proto",
        ":intrinsic_params_proto",
    ],
)

py_proto_library(
    name = "camera_params_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_params_proto"],
)

proto_library(
    name = "dimensions_proto",
    srcs = ["dimensions.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

go_proto_library(
    name = "dimensions_go_proto",
    importpath = "intrinsic/perception/proto/dimensions_go_proto",
    protos = [":dimensions_proto"],
    visibility = ["//intrinsic:public_api_users"],
)

proto_library(
    name = "distortion_params_proto",
    srcs = ["distortion_params.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

cc_proto_library(
    name = "distortion_params_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":distortion_params_proto"],
)

py_proto_library(
    name = "distortion_params_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":distortion_params_proto"],
)

go_proto_library(
    name = "distortion_params_go_proto",
    importpath = "intrinsic/perception/proto/distortion_params_go_proto",
    protos = [":distortion_params_proto"],
    visibility = ["//intrinsic:public_api_users"],
)

proto_library(
    name = "post_processing_proto",
    srcs = ["post_processing.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":dimensions_proto",
        ":image_buffer_proto",
    ],
)

py_proto_library(
    name = "post_processing_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":post_processing_proto"],
)

proto_library(
    name = "intrinsic_params_proto",
    srcs = ["intrinsic_params.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":dimensions_proto"],
)

cc_proto_library(
    name = "intrinsic_params_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":intrinsic_params_proto"],
)

py_proto_library(
    name = "intrinsic_params_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":intrinsic_params_proto"],
)

go_proto_library(
    name = "intrinsic_params_go_proto",
    importpath = "intrinsic/perception/proto/intrinsic_params_go_proto",
    protos = [":intrinsic_params_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":dimensions_go_proto"],
)

proto_library(
    name = "pose_estimation_result",
    srcs = ["pose_estimation_result.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":image_buffer_proto",
        "//intrinsic/math/proto:pose_proto",
        "@com_google_protobuf//:any_proto",
    ],
)

proto_library(
    name = "hand_eye_calibration",
    srcs = ["hand_eye_calibration.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":intrinsic_calibration_proto",
        ":pattern_detection_result_proto",
        "//intrinsic/math/proto:pose_proto",
    ],
)

proto_library(
    name = "camera_to_robot_calibration_with_metadata",
    srcs = ["camera_to_robot_calibration_with_metadata.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":hand_eye_calibration",
    ],
)

proto_library(
    name = "image_buffer_proto",
    srcs = ["image_buffer.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":dimensions_proto"],
)

proto_library(
    name = "camera_drivers_proto",
    srcs = ["camera_drivers.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

cc_proto_library(
    name = "camera_drivers_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_drivers_proto"],
)

py_proto_library(
    name = "camera_drivers_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_drivers_proto"],
)

go_proto_library(
    name = "camera_drivers_go_proto",
    importpath = "intrinsic/perception/proto/camera_drivers_go_proto",
    protos = [":camera_drivers_proto"],
    visibility = ["//intrinsic:public_api_users"],
)

proto_library(
    name = "camera_identifier_proto",
    srcs = ["camera_identifier.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":camera_drivers_proto",
        "@com_google_protobuf//:any_proto",
    ],
)

cc_proto_library(
    name = "camera_identifier_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_identifier_proto"],
)

py_proto_library(
    name = "camera_identifier_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_identifier_proto"],
)

go_proto_library(
    name = "camera_identifier_go_proto",
    importpath = "intrinsic/perception/proto/camera_identifier_go_proto",
    protos = [":camera_identifier_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_drivers_go_proto"],
)

proto_library(
    name = "camera_config",
    srcs = ["camera_config.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":camera_identifier_proto",
        ":camera_settings_proto",
        ":distortion_params_proto",
        ":intrinsic_params_proto",
        ":sensor_config",
    ],
)

py_proto_library(
    name = "camera_config_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_config"],
)

cc_proto_library(
    name = "camera_config_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["camera_config"],
)

go_proto_library(
    name = "camera_config_go_proto",
    importpath = "intrinsic/perception/proto/camera_config_go_proto",
    protos = [":camera_config"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":camera_identifier_go_proto",
        ":camera_settings_go_proto",
        ":distortion_params_go_proto",
        ":intrinsic_params_go_proto",
        ":sensor_config_go_proto",
    ],
)

proto_library(
    name = "sensor_config",
    srcs = ["sensor_config.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":camera_params_proto",
        "//intrinsic/math/proto:pose_proto",
    ],
)

py_proto_library(
    name = "sensor_config_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":sensor_config"],
)

cc_proto_library(
    name = "sensor_config_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["sensor_config"],
)

go_proto_library(
    name = "sensor_config_go_proto",
    importpath = "intrinsic/perception/proto/sensor_config_go_proto",
    protos = [":sensor_config"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":camera_params_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
    ],
)

proto_library(
    name = "sensor_image",
    srcs = ["sensor_image.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":image_buffer_proto",
        ":sensor_config",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

py_proto_library(
    name = "sensor_image_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":sensor_image"],
)

cc_proto_library(
    name = "sensor_image_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["sensor_image"],
)

proto_library(
    name = "capture_result_proto",
    srcs = ["capture_result.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":sensor_image",
        "@com_google_protobuf//:duration_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

cc_proto_library(
    name = "capture_result_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":capture_result_proto"],
)

py_proto_library(
    name = "capture_result_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":capture_result_proto"],
)

proto_library(
    name = "capture_data_proto",
    srcs = ["capture_data.proto"],
    deprecation = "This proto is deprecated. Please use //intrinsic/perception/proto/v1:capture_data_proto instead.",
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/platform/pubsub:storage_location",
    ],
)

cc_proto_library(
    name = "capture_data_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":capture_data_proto"],
)

py_proto_library(
    name = "capture_data_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":capture_data_proto"],
)

proto_library(
    name = "camera_settings_proto",
    srcs = ["camera_settings.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/perception/proto/v1:settings_proto",
        "@com_google_protobuf//:empty_proto",
    ],
)

py_proto_library(
    name = "camera_settings_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_settings_proto"],
)

cc_proto_library(
    name = "camera_settings_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_settings_proto"],
)

go_proto_library(
    name = "camera_settings_go_proto",
    importpath = "intrinsic/perception/proto/camera_settings_go_proto",
    protos = [":camera_settings_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = ["//intrinsic/perception/proto/v1:settings_go_proto"],
)

py_proto_library(
    name = "image_buffer_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":image_buffer_proto"],
)

go_proto_library(
    name = "image_buffer_go_proto",
    importpath = "intrinsic/perception/proto/image_buffer_go_proto",
    protos = [":image_buffer_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":dimensions_go_proto"],
)

proto_library(
    name = "intrinsic_calibration_proto",
    srcs = ["intrinsic_calibration.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":camera_params_proto",
        ":dimensions_proto",
        ":pattern_detection_result_proto",
    ],
)

py_proto_library(
    name = "intrinsic_calibration_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":intrinsic_calibration_proto"],
)

proto_library(
    name = "pattern_detection_result_proto",
    srcs = ["pattern_detection_result.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":vector_proto"],
)

py_proto_library(
    name = "pattern_detection_result_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":pattern_detection_result_proto"],
)

proto_library(
    name = "vector_proto",
    srcs = ["vector.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

py_proto_library(
    name = "vector_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":vector_proto"],
)

cc_proto_library(
    name = "camera_params_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_params_proto"],
)

go_proto_library(
    name = "camera_params_go_proto",
    importpath = "intrinsic/perception/proto/camera_params_go_proto",
    protos = [":camera_params_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":distortion_params_go_proto",
        ":intrinsic_params_go_proto",
    ],
)

proto_library(
    name = "charuco_pattern_proto",
    srcs = ["charuco_pattern.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

cc_proto_library(
    name = "charuco_pattern_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":charuco_pattern_proto"],
)

py_proto_library(
    name = "charuco_pattern_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":charuco_pattern_proto"],
)

proto_library(
    name = "pattern_detection_config",
    srcs = ["pattern_detection_config.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":charuco_pattern_proto",
    ],
)

cc_proto_library(
    name = "pattern_detection_config_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["pattern_detection_config"],
)

py_proto_library(
    name = "pattern_detection_config_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":pattern_detection_config"],
)

go_proto_library(
    name = "pattern_detection_result_go_proto",
    importpath = "intrinsic/perception/proto/pattern_detection_result_go_proto",
    protos = [":pattern_detection_result_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":vector_go_proto"],
)

go_proto_library(
    name = "pose_estimation_result_go_proto",
    importpath = "intrinsic/perception/proto/pose_estimation_result_go_proto",
    protos = [":pose_estimation_result"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":image_buffer_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
    ],
)

py_proto_library(
    name = "dimensions_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":dimensions_proto"],
)

cc_proto_library(
    name = "post_processing_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["post_processing_proto"],
)

cc_proto_library(
    name = "hand_eye_calibration_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["hand_eye_calibration"],
)

py_proto_library(
    name = "hand_eye_calibration_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":hand_eye_calibration"],
)

go_proto_library(
    name = "hand_eye_calibration_go_proto",
    importpath = "intrinsic/perception/proto/hand_eye_calibration_go_proto",
    protos = [":hand_eye_calibration"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":intrinsic_calibration_go_proto",
        "//intrinsic/math/proto:point_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/math/proto:quaternion_go_proto",
        "//intrinsic/perception/proto:pattern_detection_result_go_proto",
    ],
)

cc_proto_library(
    name = "camera_to_robot_calibration_with_metadata_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_to_robot_calibration_with_metadata"],
)

py_proto_library(
    name = "camera_to_robot_calibration_with_metadata_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_to_robot_calibration_with_metadata"],
)

go_proto_library(
    name = "camera_to_robot_calibration_with_metadata_go_proto",
    importpath = "intrinsic/perception/proto/camera_to_robot_calibration_with_metadata_go_proto",
    protos = [":camera_to_robot_calibration_with_metadata"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":hand_eye_calibration_go_proto",
    ],
)

go_proto_library(
    name = "intrinsic_calibration_go_proto",
    importpath = "intrinsic/perception/proto/intrinsic_calibration_go_proto",
    protos = [":intrinsic_calibration_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/perception/proto:camera_params_go_proto",
        "//intrinsic/perception/proto:dimensions_go_proto",
        "//intrinsic/perception/proto:pattern_detection_result_go_proto",
    ],
)

go_proto_library(
    name = "vector_go_proto",
    importpath = "intrinsic/perception/proto/vector_go_proto",
    protos = [":vector_proto"],
    visibility = ["//intrinsic:public_api_users"],
)

proto_library(
    name = "pose_estimator_id_proto",
    srcs = ["pose_estimator_id.proto"],
    deprecation = "This proto is deprecated. Please use //intrinsic/perception/proto/v1:pose_estimator_id_proto instead.",
    visibility = ["//intrinsic:public_api_users"],
)

cc_proto_library(
    name = "pose_estimator_id_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":pose_estimator_id_proto"],
)

py_proto_library(
    name = "pose_estimator_id_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":pose_estimator_id_proto"],
)

go_proto_library(
    name = "pose_estimator_id_go_proto",
    importpath = "intrinsic/perception/proto/pose_estimator_id_go_proto",
    protos = [":pose_estimator_id_proto"],
    visibility = ["//intrinsic:public_api_users"],
)

proto_library(
    name = "image_processing_proto",
    srcs = ["image_processing.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

go_proto_library(
    name = "image_processing_go_proto",
    importpath = "intrinsic/perception/proto/image_processing_go_proto",
    protos = [":image_processing_proto"],
    visibility = ["//intrinsic:public_api_users"],
)

proto_library(
    name = "point_cloud_proto",
    srcs = ["point_cloud.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/perception/proto/v1:sensor_image_proto",
    ],
)

proto_library(
    name = "pose_estimate_in_root_proto",
    srcs = ["pose_estimate_in_root.proto"],
    deprecation = "This proto is deprecated. Please use //intrinsic/perception/proto/v1:pose_estimate_in_root_proto instead.",
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/math/proto:pose_proto",
    ],
)

cc_proto_library(
    name = "pose_estimate_in_root_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":pose_estimate_in_root_proto"],
)

py_proto_library(
    name = "pose_estimate_in_root_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":pose_estimate_in_root_proto"],
)

go_proto_library(
    name = "pose_estimate_in_root_go_proto",
    importpath = "intrinsic/perception/proto/pose_estimate_in_root_go_proto",
    protos = [":pose_estimate_in_root_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/math/proto:pose_go_proto",
    ],
)

py_proto_library(
    name = "point_cloud_proto_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":point_cloud_proto"],
)

cc_proto_library(
    name = "point_cloud_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":point_cloud_proto"],
)

go_proto_library(
    name = "point_cloud_go_proto",
    importpath = "intrinsic/perception/proto/point_cloud_go_proto",
    protos = [":point_cloud_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/perception/proto/v1:sensor_image_go_proto",
    ],
)

cc_proto_library(
    name = "image_processing_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":image_processing_proto"],
)

py_proto_library(
    name = "image_processing_proto_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":image_processing_proto"],
)
