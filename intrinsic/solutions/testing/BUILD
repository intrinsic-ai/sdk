# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@rules_python//python:defs.bzl", "py_library", "py_test")
load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "proto_source_code_info_transitive_descriptor_set")

package(
    default_testonly = True,
    default_visibility = ["//visibility:public"],
)

proto_library(
    name = "test_skill_params_proto",
    testonly = True,
    srcs = ["test_skill_params.proto"],
    deps = [
        "//intrinsic/executive/proto:test_message_proto",
        "//intrinsic/icon/proto:joint_space_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/motion_planning/proto:motion_target_proto",
        "//intrinsic/perception/proto/v1:pose_estimator_id_proto",
        "//intrinsic/world/proto:collision_settings_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
        "//intrinsic/world/proto:robot_payload_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_proto",
        "@com_google_protobuf//:duration_proto",
    ],
)

py_proto_library(
    name = "test_skill_params_py_pb2",
    testonly = True,
    deps = [":test_skill_params_proto"],
)

proto_source_code_info_transitive_descriptor_set(
    name = "test_skill_params_proto_descriptors",
    testonly = True,
    visibility = [
        "//intrinsic/executive/clips_cpp:__subpackages__",
        "//intrinsic/solutions:__subpackages__",
    ],
    deps = [":test_skill_params_proto"],
)

py_library(
    name = "compare",
    srcs = ["compare.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_test(
    name = "compare_test",
    srcs = ["compare_test.py"],
    deps = [
        ":compare",
        "//intrinsic/executive/proto:test_message_py_pb2",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
    ],
)

py_library(
    name = "skill_test_utils",
    testonly = True,
    srcs = ["skill_test_utils.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":test_skill_params_py_pb2",
        "//intrinsic/resources/client:resource_registry_client",
        "//intrinsic/resources/proto:resource_handle_py_pb2",
        "//intrinsic/resources/proto:resource_registry_py_pb2",
        "//intrinsic/skills/client:skill_registry_client",
        "//intrinsic/skills/proto:skill_registry_py_pb2",
        "//intrinsic/skills/proto:skills_py_pb2",
        "@com_google_absl_py//absl/flags",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_protobuf//:protobuf_python",
    ],
)
