# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")
load("@rules_python//python:defs.bzl", "py_library", "py_test")
load("//intrinsic/skills/build_defs:skill.bzl", "cc_skill", "py_skill", "skill_manifest")

package(default_visibility = [
    "//visibility:public",
])

cc_library(
    name = "say_skill_cc",
    srcs = ["say_skill.cc"],
    hdrs = ["say_skill.h"],
    deps = [
        ":say_skill_cc_proto",
        "//intrinsic/skills/cc:skill_interface",
        "//intrinsic/skills/cc:skill_interface_utils",
        "//intrinsic/skills/proto:footprint_cc_proto",
        "//intrinsic/skills/proto:skill_service_cc_proto",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/time",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "say_skill_test",
    srcs = ["say_skill_test.cc"],
    deps = [
        ":say_skill_cc",
        ":say_skill_cc_proto",
        "//intrinsic/skills/cc:skill_canceller",
        "//intrinsic/skills/proto:skill_manifest_cc_proto",
        "//intrinsic/skills/testing:skill_test_utils_cc",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "//intrinsic/util/thread",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/time",
    ],
)

skill_manifest(
    name = "say_skill_cc_manifest",
    src = "say_skill_cc_manifest.textproto",
    deps = [":say_skill_proto"],
)

cc_skill(
    name = "say_skill",
    manifest = ":say_skill_cc_manifest",
    deps = [
        ":say_skill_cc",
        ":say_skill_cc_proto",
    ],
)

proto_library(
    name = "say_skill_proto",
    srcs = ["say_skill.proto"],
    deps = ["//intrinsic/assets/proto:field_metadata_proto"],
)

cc_proto_library(
    name = "say_skill_cc_proto",
    deps = [":say_skill_proto"],
)

py_proto_library(
    name = "say_skill_proto_py_pb2",
    deps = [":say_skill_proto"],
)

py_library(
    name = "pysay_skill_py",
    srcs = ["pysay_skill.py"],
    deps = [
        ":say_skill_proto_py_pb2",
        "//intrinsic/skills/proto:footprint_py_pb2",
        "//intrinsic/skills/python:skill_interface",
        "//intrinsic/skills/python:skill_interface_utils",
        "//intrinsic/util:decorators",
        "@com_google_absl_py//absl/logging",
    ],
)

skill_manifest(
    name = "pysay_skill_manifest",
    src = "pysay_skill_manifest.textproto",
    deps = [":say_skill_proto"],
)

py_skill(
    name = "pysay_skill",
    manifest = ":pysay_skill_manifest",
    deps = [
        "pysay_skill_py",
        ":say_skill_proto_py_pb2",
    ],
)

proto_library(
    name = "demo_skill_proto",
    srcs = ["demo_skill.proto"],
    deps = ["//intrinsic/motion_planning/proto/v1:motion_specification_proto"],
)

cc_proto_library(
    name = "demo_skill_cc_proto",
    deps = [":demo_skill_proto"],
)

py_proto_library(
    name = "demo_skill_proto_py_pb2",
    deps = [":demo_skill_proto"],
)

cc_library(
    name = "demo_skill_cc",
    srcs = ["demo_skill.cc"],
    hdrs = ["demo_skill.h"],
    deps = [
        ":demo_skill_cc_proto",
        "//intrinsic/icon/actions:trajectory_tracking_action_info",
        "//intrinsic/icon/cc_client:client_utils",
        "//intrinsic/icon/cc_client:condition",
        "//intrinsic/icon/cc_client:session",
        "//intrinsic/icon/common:id_types",
        "//intrinsic/icon/equipment:channel_factory",
        "//intrinsic/icon/equipment:equipment_utils",
        "//intrinsic/icon/proto:joint_space_cc_proto",
        "//intrinsic/motion_planning:motion_planner_client",
        "//intrinsic/motion_planning/skills:motion_planning_util",
        "//intrinsic/resources/proto:resource_handle_cc_proto",
        "//intrinsic/skills/cc:equipment_pack",
        "//intrinsic/skills/cc:skill_interface",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/objects:kinematic_object",
        "//intrinsic/world/objects:object_world_client",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
        "@abseil-cpp//absl/types:span",
        "@com_google_protobuf//:protobuf",
    ],
)

skill_manifest(
    name = "demo_skill_cc_manifest",
    src = "demo_skill_cc_manifest.textproto",
    incompatible_disallow_manifest_dependencies = False,
    deps = [":demo_skill_proto"],
)

cc_skill(
    name = "demo_skill",
    manifest = ":demo_skill_cc_manifest",
    deps = [
        ":demo_skill_cc",
        ":demo_skill_cc_proto",
    ],
)

py_library(
    name = "pydemo_skill_py",
    srcs = ["pydemo_skill.py"],
    deps = [
        ":demo_skill_proto_py_pb2",
        "//intrinsic/icon/equipment:equipment_utils_py",
        "//intrinsic/icon/proto:joint_space_py_pb2",
        "//intrinsic/icon/python:create_action_utils",
        "//intrinsic/icon/python:icon",
        "//intrinsic/motion_planning/proto/v1:robot_specification_py_pb2",
        "//intrinsic/skills/python:skill_interface",
        "//intrinsic/util:decorators",
        "//intrinsic/world/proto:object_world_refs_py_pb2",
    ],
)

skill_manifest(
    name = "pydemo_skill_manifest",
    src = "pydemo_skill_manifest.textproto",
    incompatible_disallow_manifest_dependencies = False,
    deps = [
        ":demo_skill_proto",
    ],
)

py_skill(
    name = "pydemo_skill",
    manifest = ":pydemo_skill_manifest",
    deps = [
        ":demo_skill_proto_py_pb2",
        ":pydemo_skill_py",
    ],
)

cc_library(
    name = "adder_skill_cc",
    srcs = ["adder_skill.cc"],
    hdrs = ["adder_skill.h"],
    deps = [
        ":adder_skill_cc_proto",
        "//intrinsic/assets/services/examples/calcserver:calc_server_cc_grpc_proto",
        "//intrinsic/assets/services/examples/calcserver:calc_server_cc_proto",
        "//intrinsic/resources/proto:resource_handle_cc_proto",
        "//intrinsic/skills/cc:equipment_pack",
        "//intrinsic/skills/cc:skill_interface",
        "//intrinsic/skills/proto:equipment_cc_proto",
        "//intrinsic/skills/proto:skill_service_cc_proto",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status:statusor",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_protobuf//:protobuf",
    ],
)

skill_manifest(
    name = "adder_skill_manifest",
    src = "adder_skill_manifest.textproto",
    incompatible_disallow_manifest_dependencies = False,
    deps = [":adder_skill_proto"],
)

cc_skill(
    name = "adder_skill",
    manifest = ":adder_skill_manifest",
    deps = [
        ":adder_skill_cc",
        ":adder_skill_cc_proto",
    ],
)

cc_test(
    name = "adder_skill_test",
    srcs = ["adder_skill_test.cc"],
    deps = [
        ":adder_skill_cc",
        ":adder_skill_cc_proto",
        "//intrinsic/assets/services/examples/calcserver:calc_server",
        "//intrinsic/assets/services/examples/calcserver:calc_server_cc_proto",
        "//intrinsic/skills/cc:equipment_pack",
        "//intrinsic/skills/cc:skill_interface",
        "//intrinsic/skills/testing:skill_test_utils_cc",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@com_google_protobuf//:protobuf",
        "@com_google_protobuf//:protobuf_lite",
    ],
)

proto_library(
    name = "adder_skill_proto",
    srcs = ["adder_skill.proto"],
)

cc_proto_library(
    name = "adder_skill_cc_proto",
    deps = [":adder_skill_proto"],
)

proto_library(
    name = "calculate_skill_proto",
    srcs = ["calculate_skill.proto"],
    deps = [
        "//intrinsic/assets/proto:field_metadata_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_proto",
        "//intrinsic/assets/services/examples/calcserver:calc_server_proto",
    ],
)

cc_proto_library(
    name = "calculate_skill_cc_proto",
    deps = [":calculate_skill_proto"],
)

cc_library(
    name = "calculate_skill_cc",
    srcs = ["calculate_skill.cc"],
    hdrs = ["calculate_skill.h"],
    deps = [
        ":calculate_skill_cc_proto",
        "//intrinsic/assets/dependencies:utils_cc",
        "//intrinsic/assets/proto/v1:resolved_dependency_cc_proto",
        "//intrinsic/assets/services/examples/calcserver:calc_server_cc_grpc_proto",
        "//intrinsic/assets/services/examples/calcserver:calc_server_cc_proto",
        "//intrinsic/skills/cc:skill_interface",
        "//intrinsic/skills/proto:skill_service_cc_proto",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status:statusor",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "calculate_skill_test",
    srcs = ["calculate_skill_test.cc"],
    deps = [
        ":calculate_skill_cc",
        ":calculate_skill_cc_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_cc_proto",
        "//intrinsic/assets/services/examples/calcserver:calc_server",
        "//intrinsic/assets/services/examples/calcserver:calc_server_cc_proto",
        "//intrinsic/skills/cc:skill_interface",
        "//intrinsic/skills/testing:skill_test_utils_cc",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@com_google_protobuf//:protobuf_lite",
    ],
)

skill_manifest(
    name = "calculate_skill_manifest",
    src = "calculate_skill_manifest.textproto",
    deps = [":calculate_skill_proto"],
)

cc_skill(
    name = "calculate_skill",
    manifest = ":calculate_skill_manifest",
    deps = [
        ":calculate_skill_cc",
        ":calculate_skill_cc_proto",
    ],
)

# [START use_world_all_rules]
# [START use_world_proto_rule]
proto_library(
    name = "use_world_proto",
    srcs = ["use_world.proto"],
    deps = ["//intrinsic/world/proto:object_world_refs_proto"],
)
# [END use_world_proto_rule]

py_proto_library(
    name = "use_world_proto_py_pb2",
    deps = [":use_world_proto"],
)

py_library(
    name = "use_world",
    srcs = ["use_world.py"],
    deps = [
        ":use_world_proto_py_pb2",
        "//intrinsic/math/python:data_types",
        "//intrinsic/skills/python:skill_interface",
        "//intrinsic/util:decorators",
        "//intrinsic/world/python:object_world_client",
        "//intrinsic/world/python:object_world_resources",
    ],
)

skill_manifest(
    name = "use_world_skill_manifest",
    src = "use_world_skill_manifest.textproto",
    incompatible_disallow_manifest_dependencies = False,
    deps = [":use_world_proto"],
)

py_skill(
    name = "use_world_skill",
    manifest = ":use_world_skill_manifest",
    deps = [
        ":use_world",
        ":use_world_proto_py_pb2",
    ],
)
# [END use_world_all_rules]

py_test(
    name = "use_world_test",
    srcs = ["use_world_test.py"],
    data = [
        ":use_world_skill_manifest",
        "//intrinsic/skills/apps/test_data:ur5e_world",
    ],
    deps = [
        ":use_world",
        ":use_world_proto_py_pb2",
        "//intrinsic/resources/proto:resource_handle_py_pb2",
        "//intrinsic/skills/proto:equipment_py_pb2",
        "//intrinsic/skills/testing:skill_test_utils",
        "//intrinsic/util/path_resolver:path_resolver_py",
        "//intrinsic/world/proto:object_world_service_py_pb2_grpc",
        "//intrinsic/world/python:object_world_client",
        "//intrinsic/world/python:world",
        "//intrinsic/world/python/testing:fake_world_service",
        "//intrinsic/world/python/testing:object_world_test_utils",
        requirement("grpcio"),
        "@com_google_absl_py//absl/testing:absltest",
    ],
)
