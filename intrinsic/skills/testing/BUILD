# Copyright 2023 Intrinsic Innovation LLC

# Testing utilities for skills.

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@rules_cc//cc:defs.bzl", "cc_library", "cc_test")
load("@rules_python//python:defs.bzl", "py_library", "py_test")
load("//bazel:go_macros.bzl", "go_proto_library")
load("//intrinsic/skills/build_defs:skill.bzl", "py_skill", "skill_manifest")

package(
    default_testonly = True,
)

proto_library(
    name = "echo_skill_proto",
    srcs = ["echo_skill.proto"],
)

cc_proto_library(
    name = "echo_skill_cc_proto",
    visibility = [
        "//intrinsic/executive/engine:__pkg__",
        "//intrinsic/skills/cc:__pkg__",
        "//intrinsic/skills/internal:__pkg__",
    ],
    deps = [":echo_skill_proto"],
)

py_proto_library(
    name = "echo_skill_py_pb2",
    visibility = [
        "//intrinsic/skills/internal:__pkg__",
        "//intrinsic/skills/python:__pkg__",
    ],
    deps = [":echo_skill_proto"],
)

py_library(
    name = "echo_skill_py",
    srcs = ["echo_skill.py"],
    visibility = ["//intrinsic/skills/internal:__pkg__"],
    deps = [
        ":echo_skill_py_pb2",
        "//intrinsic/skills/python:skill_interface",
        "//intrinsic/skills/python:skill_interface_utils",
        "//intrinsic/util:decorators",
    ],
)

py_test(
    name = "echo_skill_py_test",
    srcs = ["echo_skill_py_test.py"],
    data = [":echo_skill_py_manifest"],
    deps = [
        ":echo_skill_py",
        ":echo_skill_py_pb2",
        ":skill_test_utils",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

cc_library(
    name = "echo_skill",
    srcs = ["echo_skill.cc"],
    hdrs = ["echo_skill.h"],
    visibility = ["//intrinsic/skills/internal:__pkg__"],
    deps = [
        ":echo_skill_cc_proto",
        "//intrinsic/skills/cc:skill_interface",
        "//intrinsic/skills/cc:skill_interface_utils",
        "//intrinsic/skills/proto:equipment_cc_proto",
        "//intrinsic/skills/proto:skill_service_cc_proto",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status:statusor",
        "@com_google_protobuf//:protobuf",
    ],
)

skill_manifest(
    name = "echo_skill_py_manifest",
    src = "echo_skill_manifest_py.textproto",
    visibility = [
        "//intrinsic/icon/release:__pkg__",
        "//intrinsic/skills/generator:__pkg__",
        "//intrinsic/skills/internal:__pkg__",
    ],
    deps = [":echo_skill_proto"],
)

py_skill(
    name = "echo_skill_py_skill",
    manifest = ":echo_skill_py_manifest",
    deps = [
        ":echo_skill_py",
        ":echo_skill_py_pb2",
    ],
)

proto_library(
    name = "no_op_skill_proto",
    testonly = False,
    srcs = ["no_op_skill.proto"],
    visibility = ["//intrinsic/skills/build_defs/tests:__pkg__"],
)

cc_proto_library(
    name = "no_op_skill_cc_proto",
    testonly = False,
    visibility = [
        "//intrinsic/executive/engine:__pkg__",
        "//intrinsic/skills/build_defs/tests:__pkg__",
        "//intrinsic/skills/internal:__pkg__",
    ],
    deps = [":no_op_skill_proto"],
)

py_proto_library(
    name = "no_op_skill_py_pb2",
    visibility = [
        "//intrinsic/skills/build_defs/tests:__pkg__",
        "//intrinsic/skills/internal:__pkg__",
    ],
    deps = [":no_op_skill_proto"],
)

go_proto_library(
    name = "no_op_skill_go_proto",
    visibility = [
        "//intrinsic/skills/build_defs:__pkg__",
    ],
    deps = [":no_op_skill_proto"],
)

py_library(
    name = "no_op_skill_py",
    srcs = ["no_op_skill.py"],
    visibility = [
        "//intrinsic/skills/build_defs/tests:__pkg__",
        "//intrinsic/skills/internal:__pkg__",
    ],
    deps = [
        ":no_op_skill_py_pb2",
        "//intrinsic/skills/proto:footprint_py_pb2",
        "//intrinsic/skills/python:skill_interface",
        "//intrinsic/skills/python:skill_interface_utils",
        "//intrinsic/util:decorators",
    ],
)

cc_library(
    name = "no_op_skill_cc",
    testonly = False,
    srcs = ["no_op_skill.cc"],
    hdrs = ["no_op_skill.h"],
    visibility = [
        "//intrinsic/skills/build_defs/tests:__pkg__",
        "//intrinsic/skills/internal:__pkg__",
    ],
    deps = [
        ":no_op_skill_cc_proto",
        "//intrinsic/skills/cc:skill_interface",
        "//intrinsic/skills/cc:skill_interface_utils",
        "//intrinsic/skills/proto:equipment_cc_proto",
        "//intrinsic/skills/proto:skill_service_cc_proto",
        "@abseil-cpp//absl/status:statusor",
        "@com_google_protobuf//:protobuf",
    ],
)

py_library(
    name = "map_skill_repository_py",
    testonly = False,
    srcs = ["map_skill_repository.py"],
    visibility = ["//intrinsic/skills/internal:__pkg__"],
    deps = [
        "//intrinsic/skills/internal:runtime_data_py",
        "//intrinsic/skills/internal:single_skill_factory_py",
        "//intrinsic/skills/internal:skill_repository_py",
        "//intrinsic/skills/python:skill_interface",
        "//intrinsic/util:decorators",
    ],
)

py_library(
    name = "skill_test_utils",
    srcs = ["skill_test_utils.py"],
    visibility = [
        "//intrinsic:public_api_users",
    ],
    deps = [
        "//intrinsic/assets:id_utils_py",
        "//intrinsic/geometry/proto:geometry_service_py_pb2_grpc",
        "//intrinsic/icon/release:file_helpers_py",
        "//intrinsic/logging/proto:context_py_pb2",
        "//intrinsic/motion_planning:motion_planner_client_py",
        "//intrinsic/motion_planning/proto/v1:motion_planner_service_py_pb2_grpc",
        "//intrinsic/resources/proto:resource_handle_py_pb2",
        "//intrinsic/skills/internal:execute_context_impl_py",
        "//intrinsic/skills/internal:get_footprint_context_impl_py",
        "//intrinsic/skills/internal:preview_context_impl_py",
        "//intrinsic/skills/proto:skill_manifest_py_pb2",
        "//intrinsic/skills/python:execute_request",
        "//intrinsic/skills/python:get_footprint_request",
        "//intrinsic/skills/python:preview_request",
        "//intrinsic/skills/python:skill_canceller",
        "//intrinsic/skills/python:skill_interface",
        "//intrinsic/skills/python:skill_logging_context",
        "//intrinsic/util/path_resolver:path_resolver_py",
        "//intrinsic/world/proto:object_world_service_py_pb2_grpc",
        "//intrinsic/world/python:object_world_client",
        requirement("grpcio"),
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_test(
    name = "skill_test_utils_py_test",
    srcs = ["skill_test_utils_test.py"],
    data = [
        ":echo_skill_py_manifest",
    ],
    main = "skill_test_utils_test.py",
    deps = [
        ":echo_skill_py_pb2",
        ":skill_test_utils",
        "//intrinsic/assets:id_utils_py",
        "//intrinsic/assets/services/examples/calcserver:calc_server_py",
        "//intrinsic/assets/services/examples/calcserver:calc_server_py_pb2",
        "//intrinsic/assets/services/examples/calcserver:calc_server_py_pb2_grpc",
        "//intrinsic/logging/proto:log_item_py_pb2",
        "//intrinsic/skills/proto:skill_manifest_py_pb2",
        "//intrinsic/skills/python:skill_logging_context",
        "//intrinsic/solutions/testing:compare",
        requirement("grpcio"),
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
    ],
)

cc_library(
    name = "skill_test_utils_cc",
    testonly = True,
    srcs = ["skill_test_utils.cc"],
    hdrs = ["skill_test_utils.h"],
    visibility = [
        "//intrinsic:public_api_users",
    ],
    deps = [
        "//intrinsic/motion_planning:motion_planner_client",
        "//intrinsic/motion_planning/proto/v1:motion_planner_service_cc_grpc_proto",
        "//intrinsic/skills/cc:equipment_pack",
        "//intrinsic/skills/cc:execute_request",
        "//intrinsic/skills/cc:skill_canceller",
        "//intrinsic/skills/cc:skill_interface",
        "//intrinsic/skills/cc:skill_logging_context",
        "//intrinsic/skills/internal:execute_context_impl",
        "//intrinsic/skills/internal:get_footprint_context_impl",
        "//intrinsic/skills/internal:preview_context_impl",
        "//intrinsic/skills/proto:skill_manifest_cc_proto",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/testing:gtest_wrapper",  # buildcleaner: keep
        "//intrinsic/world/objects:object_world_client",
        "//intrinsic/world/proto:object_world_service_cc_grpc_proto",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
        "@com_github_grpc_grpc//:grpc_security_base",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "skill_test_utils_test",
    srcs = ["skill_test_utils_test.cc"],
    deps = [
        ":echo_skill",
        ":echo_skill_cc_proto",
        ":skill_test_utils_cc",
        "//intrinsic/assets:id_utils",
        "//intrinsic/logging/proto:context_cc_proto",
        "//intrinsic/motion_planning:motion_planner_client",
        "//intrinsic/motion_planning/proto/v1:motion_planner_service_cc_grpc_proto",
        "//intrinsic/motion_planning/proto/v1:motion_planner_service_cc_proto",
        "//intrinsic/skills/cc:equipment_pack",
        "//intrinsic/skills/cc:skill_canceller",
        "//intrinsic/skills/cc:skill_interface",
        "//intrinsic/skills/cc:skill_logging_context",
        "//intrinsic/skills/proto:skill_manifest_cc_proto",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "//intrinsic/world/objects:object_world_client",
        "//intrinsic/world/proto:object_world_service_cc_grpc_proto",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
    ],
)
