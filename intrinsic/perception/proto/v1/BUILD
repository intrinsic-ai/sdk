# Copyright 2023 Intrinsic Innovation LLC

# Public API protos for the Intrinsic Perception API v1.

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")
load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "proto_source_code_info_transitive_descriptor_set")

proto_library(
    name = "camera_config_proto",
    srcs = ["camera_config.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":camera_identifier_proto",
        ":camera_settings_proto",
        ":sensor_config_proto",
    ],
)

cc_proto_library(
    name = "camera_config_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["camera_config_proto"],
)

py_proto_library(
    name = "camera_config_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_config_proto"],
)

go_proto_library(
    name = "camera_config_go_proto",
    go_deps = [
        ":camera_identifier_go_proto",
        ":camera_settings_go_proto",
        ":sensor_config_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_config_proto"],
)

proto_library(
    name = "camera_drivers_proto",
    srcs = ["camera_drivers.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

cc_proto_library(
    name = "camera_drivers_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["camera_drivers_proto"],
)

py_proto_library(
    name = "camera_drivers_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_drivers_proto"],
)

go_proto_library(
    name = "camera_drivers_go_proto",
    deps = [":camera_drivers_proto"],
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
    go_deps = [":camera_drivers_go_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_identifier_proto"],
)

proto_library(
    name = "camera_params_proto",
    srcs = ["camera_params.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":distortion_params_proto",
        ":intrinsic_params_proto",
    ],
)

cc_proto_library(
    name = "camera_params_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_params_proto"],
)

py_proto_library(
    name = "camera_params_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_params_proto"],
)

go_proto_library(
    name = "camera_params_go_proto",
    go_deps = [
        ":distortion_params_go_proto",
        ":intrinsic_params_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_params_proto"],
)

proto_library(
    name = "camera_service_proto",
    srcs = ["camera_service.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":camera_config_proto",
        ":camera_identifier_proto",
        ":camera_params_proto",
        ":camera_settings_proto",
        ":capture_result_proto",
        ":dimensions_proto",
        ":image_buffer_proto",
        ":post_processing_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/platform/pubsub:storage_location",
        "@com_google_protobuf//:duration_proto",
    ],
)

proto_source_code_info_transitive_descriptor_set(
    name = "camera_service_proto_descriptor_set",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_service_proto"],
)

cc_proto_library(
    name = "camera_service_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_service_proto"],
)

cc_grpc_library(
    name = "camera_service_cc_grpc",
    srcs = [":camera_service_proto"],
    grpc_only = True,
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_service_cc_proto"],
)

py_proto_library(
    name = "camera_service_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_service_proto"],
)

py_grpc_library(
    name = "camera_service_py_pb2_grpc",
    srcs = [":camera_service_proto"],
    grpc_library = requirement("grpcio"),
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_service_py_pb2"],
)

go_proto_library(
    name = "camera_service_go_proto",
    go_deps = [
        ":camera_config_go_proto",
        ":camera_identifier_go_proto",
        ":camera_params_go_proto",
        ":camera_settings_go_proto",
        ":capture_result_go_proto",
        ":dimensions_go_proto",
        ":image_buffer_go_proto",
        ":post_processing_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/platform/pubsub:storage_location_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_service_proto"],
)

go_grpc_library(
    name = "camera_service_go_grpc_proto",
    srcs = [":camera_service_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":camera_config_go_proto",
        ":camera_identifier_go_proto",
        ":camera_params_go_proto",
        ":camera_settings_go_proto",
        ":capture_result_go_proto",
        ":dimensions_go_proto",
        ":image_buffer_go_proto",
        ":post_processing_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/platform/pubsub:storage_location_go_proto",
    ],
)

proto_library(
    name = "camera_settings_proto",
    srcs = ["camera_settings.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":settings_proto",
        "@com_google_protobuf//:empty_proto",
    ],
)

cc_proto_library(
    name = "camera_settings_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_settings_proto"],
)

py_proto_library(
    name = "camera_settings_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_settings_proto"],
)

go_proto_library(
    name = "camera_settings_go_proto",
    go_deps = [":settings_go_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":camera_settings_proto"],
)

proto_library(
    name = "capture_data_proto",
    srcs = ["capture_data.proto"],
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

go_proto_library(
    name = "capture_data_go_proto",
    go_deps = [
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/platform/pubsub:storage_location_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":capture_data_proto"],
)

proto_library(
    name = "capture_result_proto",
    srcs = ["capture_result.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":sensor_image_proto",
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

go_proto_library(
    name = "capture_result_go_proto",
    go_deps = [
        ":sensor_image_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":capture_result_proto"],
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

go_proto_library(
    name = "charuco_pattern_go_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":charuco_pattern_proto"],
)

proto_library(
    name = "dimensions_proto",
    srcs = ["dimensions.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

cc_proto_library(
    name = "dimensions_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":dimensions_proto"],
)

py_proto_library(
    name = "dimensions_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":dimensions_proto"],
)

go_proto_library(
    name = "dimensions_go_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":dimensions_proto"],
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
    visibility = ["//intrinsic:public_api_users"],
    deps = [":distortion_params_proto"],
)

proto_library(
    name = "image_buffer_proto",
    srcs = ["image_buffer.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":dimensions_proto"],
)

cc_proto_library(
    name = "image_buffer_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["image_buffer_proto"],
)

py_proto_library(
    name = "image_buffer_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":image_buffer_proto"],
)

go_proto_library(
    name = "image_buffer_go_proto",
    go_deps = [":dimensions_go_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":image_buffer_proto"],
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
    go_deps = [":dimensions_go_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":intrinsic_params_proto"],
)

proto_library(
    name = "pattern_detection_config_proto",
    srcs = ["pattern_detection_config.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":charuco_pattern_proto"],
)

cc_proto_library(
    name = "pattern_detection_config_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":pattern_detection_config_proto"],
)

go_proto_library(
    name = "pattern_detection_config_go_proto",
    go_deps = [
        ":charuco_pattern_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":pattern_detection_config_proto"],
)

py_proto_library(
    name = "pattern_detection_config_py_pb2",
    deps = [":pattern_detection_config_proto"],
)

proto_library(
    name = "perception_model_proto",
    srcs = ["perception_model.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":pose_estimation_config",
        "//intrinsic/assets/data/proto/v1:referenced_data_proto",
    ],
)

go_proto_library(
    name = "perception_model_go_proto",
    go_deps = [
        ":pose_estimation_config_go_proto",
        "//intrinsic/assets/data/proto/v1:referenced_data_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":perception_model_proto"],
)

py_proto_library(
    name = "perception_model_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":perception_model_proto"],
)

cc_proto_library(
    name = "perception_model_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":perception_model_proto"],
)

proto_library(
    name = "pose_estimate_in_root_proto",
    srcs = ["pose_estimate_in_root.proto"],
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
    go_deps = [
        "//intrinsic/math/proto:pose_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":pose_estimate_in_root_proto"],
)

proto_library(
    name = "pose_estimation_config",
    srcs = ["pose_estimation_config.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":intrinsic_params_proto",
        ":target_proto",
        "@com_google_protobuf//:any_proto",
    ],
)

cc_proto_library(
    name = "pose_estimation_config_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["pose_estimation_config"],
)

py_proto_library(
    name = "pose_estimation_config_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":pose_estimation_config"],
)

go_proto_library(
    name = "pose_estimation_config_go_proto",
    go_deps = [
        ":intrinsic_params_go_proto",
        ":target_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":pose_estimation_config"],
)

proto_library(
    name = "pose_estimator_id_proto",
    srcs = ["pose_estimator_id.proto"],
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
    visibility = ["//intrinsic:public_api_users"],
    deps = [":pose_estimator_id_proto"],
)

proto_library(
    name = "pose_range_proto",
    srcs = ["pose_range.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

cc_proto_library(
    name = "pose_range_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":pose_range_proto"],
)

py_proto_library(
    name = "pose_range_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":pose_range_proto"],
)

go_proto_library(
    name = "pose_range_go_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":pose_range_proto"],
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

cc_proto_library(
    name = "post_processing_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["post_processing_proto"],
)

py_proto_library(
    name = "post_processing_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":post_processing_proto"],
)

go_proto_library(
    name = "post_processing_go_proto",
    go_deps = [
        ":dimensions_go_proto",
        ":image_buffer_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":post_processing_proto"],
)

proto_library(
    name = "sensor_config_proto",
    srcs = ["sensor_config.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":camera_params_proto",
        "//intrinsic/math/proto:pose_proto",
    ],
)

cc_proto_library(
    name = "sensor_config_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["sensor_config_proto"],
)

py_proto_library(
    name = "sensor_config_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":sensor_config_proto"],
)

go_proto_library(
    name = "sensor_config_go_proto",
    go_deps = [
        ":camera_params_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":sensor_config_proto"],
)

proto_library(
    name = "sensor_image_proto",
    srcs = ["sensor_image.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":image_buffer_proto",
        ":sensor_config_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

cc_proto_library(
    name = "sensor_image_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["sensor_image_proto"],
)

py_proto_library(
    name = "sensor_image_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":sensor_image_proto"],
)

go_proto_library(
    name = "sensor_image_go_proto",
    go_deps = [
        ":image_buffer_go_proto",
        ":sensor_config_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":sensor_image_proto"],
)

proto_library(
    name = "settings_proto",
    srcs = ["settings.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

cc_proto_library(
    name = "settings_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":settings_proto"],
)

py_proto_library(
    name = "settings_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":settings_proto"],
)

go_proto_library(
    name = "settings_go_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":settings_proto"],
)

proto_library(
    name = "symmetry_proto",
    srcs = ["symmetry.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/math/proto:vector3_proto",
    ],
)

cc_proto_library(
    name = "symmetry_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":symmetry_proto"],
)

py_proto_library(
    name = "symmetry_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":symmetry_proto"],
)

go_proto_library(
    name = "symmetry_go_proto",
    go_deps = ["//intrinsic/math/proto:vector3_go_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":symmetry_proto"],
)

proto_library(
    name = "target_proto",
    srcs = ["target.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":charuco_pattern_proto",
        ":pose_range_proto",
        ":symmetry_proto",
        "//intrinsic/math/proto:affine_proto",
        "//intrinsic/scene/proto/v1:scene_object_proto",
    ],
)

cc_proto_library(
    name = "target_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":target_proto"],
)

py_proto_library(
    name = "target_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":target_proto"],
)

go_proto_library(
    name = "target_go_proto",
    go_deps = [
        ":charuco_pattern_go_proto",
        ":pose_range_go_proto",
        ":symmetry_go_proto",
        "//intrinsic/math/proto:affine_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":target_proto"],
)

proto_library(
    name = "train_service_proto",
    srcs = ["train_service.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":pose_estimation_config",
        "//intrinsic/assets/proto:id_proto",
        "//intrinsic/util/status:extended_status_proto",
        "@com_google_googleapis//google/longrunning:operations_proto",
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

proto_source_code_info_transitive_descriptor_set(
    name = "train_service_proto_descriptor_set",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":train_service_proto"],
)

cc_proto_library(
    name = "train_service_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":train_service_proto"],
)

cc_grpc_library(
    name = "train_service_cc_grpc_proto",
    srcs = [":train_service_proto"],
    grpc_only = True,
    visibility = ["//intrinsic:public_api_users"],
    deps = [":train_service_cc_proto"],
)

go_proto_library(
    name = "train_service_go_proto",
    go_deps = [
        ":pose_estimation_config_go_proto",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/util/status:extended_status_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":train_service_proto"],
)

go_grpc_library(
    name = "train_service_go_grpc_proto",
    srcs = [":train_service_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":pose_estimation_config_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
        "//intrinsic/util/status:extended_status_go_proto",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_protobuf//types/known/anypb",
        "@org_golang_google_protobuf//types/known/durationpb",
        "@org_golang_google_protobuf//types/known/emptypb",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

py_proto_library(
    name = "train_service_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":train_service_proto"],
)

py_grpc_library(
    name = "train_service_py_pb2_grpc",
    srcs = [":train_service_proto"],
    grpc_library = requirement("grpcio"),
    visibility = ["//intrinsic:public_api_users"],
    deps = [":train_service_py_pb2"],
)
