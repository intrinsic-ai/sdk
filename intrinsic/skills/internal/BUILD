# Copyright 2023 Intrinsic Innovation LLC

# This directory should contain code that is used internally to our platform.
# Some files may still be exported as part of the Skills SDK, but they should
# not needed to be included directly when authoring a Skill.

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_python//python:defs.bzl", "py_library", "py_test")
load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//intrinsic/skills:__subpackages__"])

py_library(
    name = "basic_compute_context_impl_py",
    srcs = ["basic_compute_context_impl.py"],
    visibility = [
        "//intrinsic/executive:__subpackages__",
        "//intrinsic/skills:__subpackages__",
    ],
    deps = [
        "//intrinsic/skills/python:basic_compute_context",
        "//intrinsic/world/python:object_world_client",
    ],
)

cc_library(
    name = "default_parameters",
    srcs = ["default_parameters.cc"],
    hdrs = ["default_parameters.h"],
    deps = [
        "//intrinsic/skills/proto:skills_cc_proto",
        "//intrinsic/util/proto:merge",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)

py_library(
    name = "default_parameters_py",
    srcs = ["default_parameters.py"],
    deps = [
        "//intrinsic/skills/internal:error_bindings",
        "@com_google_protobuf//:protobuf_python",
        "@pybind11_abseil//pybind11_abseil:import_status_module",
    ],
)

cc_library(
    name = "equipment_utilities",
    srcs = ["equipment_utilities.cc"],
    hdrs = ["equipment_utilities.h"],
    deps = [
        "//intrinsic/resources/proto:resource_handle_cc_proto",
        "//intrinsic/skills/proto:equipment_cc_proto",
        "//intrinsic/skills/proto:footprint_cc_proto",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "error_utils",
    srcs = ["error_utils.cc"],
    hdrs = ["error_utils.h"],
    deps = [
        ":runtime_data",
        "//intrinsic/assets/proto:status_spec_cc_proto",
        "//intrinsic/logging/proto:context_cc_proto",
        "//intrinsic/skills/proto:error_cc_proto",
        "//intrinsic/util:proto_time",
        "//intrinsic/util/proto:type_url",
        "//intrinsic/util/status:extended_status_cc_proto",
        "//intrinsic/util/status:get_extended_status",
        "//intrinsic/util/status:status_conversion_rpc",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:cord",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/strings:string_view",
        "@com_google_googleapis//google/rpc:status_cc_proto",
        "@com_google_protobuf//:protobuf",
    ],
)

py_library(
    name = "error_utils_py",
    srcs = ["error_utils.py"],
    deps = [
        "//intrinsic/skills/proto:error_py_pb2",
        "//intrinsic/util/grpc:error_handling",
        requirement("grpcio"),
        "@pybind11_abseil//pybind11_abseil:import_status_module",
    ],
)

py_library(
    name = "error_bindings",
    srcs = ["error_bindings.py"],
    deps = ["@pybind11_abseil//pybind11_abseil:import_status_module"],
)

py_test(
    name = "error_bindings_test",
    srcs = ["error_bindings_test.py"],
    deps = [
        ":error_bindings",
        "@com_google_absl_py//absl/testing:absltest",
        "@pybind11_abseil//pybind11_abseil:import_status_module",
    ],
)

cc_library(
    name = "get_footprint_context_impl",
    srcs = ["get_footprint_context_impl.cc"],
    hdrs = ["get_footprint_context_impl.h"],
    deps = [
        "//intrinsic/logging/proto:context_cc_proto",
        "//intrinsic/motion_planning:motion_planner_client",
        "//intrinsic/resources/proto:resource_handle_cc_proto",
        "//intrinsic/skills/cc:equipment_pack",
        "//intrinsic/skills/cc:skill_interface",
        "//intrinsic/skills/proto:equipment_cc_proto",
        "//intrinsic/skills/proto:footprint_cc_proto",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/objects:frame",
        "//intrinsic/world/objects:kinematic_object",
        "//intrinsic/world/objects:object_world_client",
        "//intrinsic/world/objects:object_world_ids",
        "//intrinsic/world/objects:world_object",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

py_library(
    name = "get_footprint_context_impl_py",
    srcs = ["get_footprint_context_impl.py"],
    deps = [
        "//intrinsic/geometry/proto:geometry_service_py_pb2_grpc",
        "//intrinsic/motion_planning:motion_planner_client_py",
        "//intrinsic/resources/proto:resource_handle_py_pb2",
        "//intrinsic/skills/python:get_footprint_context",
        "//intrinsic/world/python:object_world_client",
        "//intrinsic/world/python:object_world_ids",
        "//intrinsic/world/python:object_world_resources",
    ],
)

cc_library(
    name = "execute_context_impl",
    hdrs = ["execute_context_impl.h"],
    deps = [
        "//intrinsic/logging/proto:context_cc_proto",
        "//intrinsic/motion_planning:motion_planner_client",
        "//intrinsic/skills/cc:equipment_pack",
        "//intrinsic/skills/cc:skill_canceller",
        "//intrinsic/skills/cc:skill_interface",
        "//intrinsic/skills/cc:skill_logging_context",
        "//intrinsic/world/objects:object_world_client",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@com_github_grpc_grpc//:grpc++",
    ],
)

py_library(
    name = "execute_context_impl_py",
    srcs = ["execute_context_impl.py"],
    deps = [
        "//intrinsic/geometry/proto:geometry_service_py_pb2_grpc",
        "//intrinsic/motion_planning:motion_planner_client_py",
        "//intrinsic/resources/proto:resource_handle_py_pb2",
        "//intrinsic/skills/python:execute_context",
        "//intrinsic/skills/python:skill_canceller",
        "//intrinsic/skills/python:skill_logging_context",
        "//intrinsic/world/python:object_world_client",
    ],
)

cc_library(
    name = "execute_context_view",
    hdrs = ["execute_context_view.h"],
    deps = [
        "//intrinsic/logging/proto:context_cc_proto",
        "//intrinsic/motion_planning:motion_planner_client",
        "//intrinsic/skills/cc:equipment_pack",
        "//intrinsic/skills/cc:skill_canceller",
        "//intrinsic/skills/cc:skill_interface",
        "//intrinsic/skills/cc:skill_logging_context",
        "//intrinsic/world/objects:object_world_client",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "preview_context_impl",
    srcs = ["preview_context_impl.cc"],
    hdrs = ["preview_context_impl.h"],
    deps = [
        "//intrinsic/logging/proto:context_cc_proto",
        "//intrinsic/motion_planning:motion_planner_client",
        "//intrinsic/resources/proto:resource_handle_cc_proto",
        "//intrinsic/skills/cc:equipment_pack",
        "//intrinsic/skills/cc:preview_context",
        "//intrinsic/skills/cc:skill_canceller",
        "//intrinsic/skills/cc:skill_logging_context",
        "//intrinsic/skills/proto:equipment_cc_proto",
        "//intrinsic/skills/proto:prediction_cc_proto",
        "//intrinsic/skills/proto:skill_service_cc_proto",
        "//intrinsic/util:proto_time",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/objects:frame",
        "//intrinsic/world/objects:kinematic_object",
        "//intrinsic/world/objects:object_world_client",
        "//intrinsic/world/objects:object_world_ids",
        "//intrinsic/world/objects:world_object",
        "//intrinsic/world/proto:object_world_updates_cc_proto",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:string_view",
        "@abseil-cpp//absl/time",
        "@com_google_protobuf//:protobuf",
    ],
)

py_library(
    name = "preview_context_impl_py",
    srcs = ["preview_context_impl.py"],
    deps = [
        "//intrinsic/geometry/proto:geometry_service_py_pb2_grpc",
        "//intrinsic/motion_planning:motion_planner_client_py",
        "//intrinsic/resources/proto:resource_handle_py_pb2",
        "//intrinsic/skills/proto:prediction_py_pb2",
        "//intrinsic/skills/python:preview_context",
        "//intrinsic/skills/python:skill_canceller",
        "//intrinsic/skills/python:skill_logging_context",
        "//intrinsic/world/proto:object_world_updates_py_pb2",
        "//intrinsic/world/python:object_world_client",
        "//intrinsic/world/python:object_world_ids",
        "//intrinsic/world/python:object_world_resources",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_library(
    name = "skill_service_impl_py",
    srcs = ["skill_service_impl.py"],
    deps = [
        ":default_parameters_py",
        ":error_bindings",
        ":error_utils_py",
        ":execute_context_impl_py",
        ":get_footprint_context_impl_py",
        ":preview_context_impl_py",
        ":runtime_data_py",
        ":skill_repository_py",
        "//intrinsic/assets:id_utils_py",
        "//intrinsic/geometry/proto:geometry_service_py_pb2_grpc",
        "//intrinsic/logging/proto:context_py_pb2",
        "//intrinsic/motion_planning:motion_planner_client_py",
        "//intrinsic/motion_planning/proto/v1:motion_planner_service_py_pb2_grpc",
        "//intrinsic/skills/proto:error_py_pb2",
        "//intrinsic/skills/proto:footprint_py_pb2",
        "//intrinsic/skills/proto:prediction_py_pb2",
        "//intrinsic/skills/proto:skill_service_py_pb2",
        "//intrinsic/skills/proto:skill_service_py_pb2_grpc",
        "//intrinsic/skills/proto:skills_py_pb2",
        "//intrinsic/skills/python:proto_utils",
        "//intrinsic/skills/python:skill_canceller",
        "//intrinsic/skills/python:skill_interface",
        "//intrinsic/skills/python:skill_logging_context",
        "//intrinsic/util/status:extended_status_py_pb2",
        "//intrinsic/util/status:status_exception",
        "//intrinsic/world/proto:object_world_service_py_pb2_grpc",
        "//intrinsic/world/python:object_world_client",
        requirement("grpcio"),
        "@com_google_absl_py//absl/logging",
        "@com_google_googleapis//google/longrunning:operations_py_proto",
        "@com_google_googleapis//google/rpc:status_py_proto",
        "@com_google_protobuf//:protobuf_python",
        "@pybind11_abseil//pybind11_abseil:import_status_module",
    ],
)

py_library(
    name = "proto_utils",
    srcs = ["proto_utils.py"],
    visibility = [
        "//visibility:public",
    ],
    deps = [
        "//intrinsic/assets:id_utils_py",
        "//intrinsic/skills/proto:skill_manifest_py_pb2",
        "//intrinsic/skills/proto:skills_py_pb2",
        "//intrinsic/util/proto:source_code_info_view_py",
        "@com_google_protobuf//:protobuf_python",
    ],
)

cc_library(
    name = "skill_repository",
    hdrs = ["skill_repository.h"],
    deps = [
        ":runtime_data",
        "//intrinsic/skills/cc:skill_interface",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

py_library(
    name = "skill_repository_py",
    srcs = ["skill_repository.py"],
    deps = [
        ":runtime_data_py",
        "//intrinsic/skills/python:skill_interface",
    ],
)

cc_library(
    name = "runtime_data",
    srcs = ["runtime_data.cc"],
    hdrs = ["runtime_data.h"],
    visibility = ["//visibility:public"],
    deps = [
        "//intrinsic/assets/proto:status_spec_cc_proto",
        "//intrinsic/skills/cc:client_common",
        "//intrinsic/skills/proto:equipment_cc_proto",
        "//intrinsic/skills/proto:skill_service_config_cc_proto",
        "//intrinsic/skills/proto:skills_cc_proto",
        "//intrinsic/util:proto_time",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/algorithm:container",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/time",
        "@abseil-cpp//absl/types:span",
        "@com_google_protobuf//:protobuf",
    ],
)

py_library(
    name = "runtime_data_py",
    srcs = ["runtime_data.py"],
    visibility = ["//visibility:public"],
    deps = [
        "//intrinsic/assets/proto:status_spec_py_pb2",
        "//intrinsic/skills/proto:equipment_py_pb2",
        "//intrinsic/skills/proto:skill_service_config_py_pb2",
        "//intrinsic/skills/proto:skills_py_pb2",
        "@com_google_protobuf//:protobuf_python",
    ],
)

cc_library(
    name = "single_skill_factory",
    srcs = ["single_skill_factory.cc"],
    hdrs = ["single_skill_factory.h"],
    visibility = ["//visibility:public"],
    deps = [
        ":runtime_data",
        ":skill_repository",
        "//intrinsic/assets:id_utils",
        "//intrinsic/skills/cc:skill_interface",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/synchronization",
    ],
)

py_library(
    name = "single_skill_factory_py",
    srcs = ["single_skill_factory.py"],
    visibility = ["//visibility:public"],
    deps = [
        ":runtime_data_py",
        ":skill_repository_py",
        "//intrinsic/assets:id_utils_py",
        "//intrinsic/skills/python:skill_interface",
        "//intrinsic/util:decorators",
    ],
)

cc_library(
    name = "skill_init",
    srcs = ["skill_init.cc"],
    hdrs = ["skill_init.h"],
    visibility = ["//visibility:public"],
    deps = [
        ":skill_registry_client",
        ":skill_repository",
        ":skill_service_impl",
        "//intrinsic/logging:data_logger_client",
        "//intrinsic/motion_planning/proto/v1:motion_planner_service_cc_grpc_proto",
        "//intrinsic/skills/proto:skill_service_config_cc_proto",
        "//intrinsic/skills/proto:skills_cc_proto",
        "//intrinsic/util/grpc",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/proto:object_world_service_cc_grpc_proto",
        "@abseil-cpp//absl/algorithm:container",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
        "@com_github_grpc_grpc//:grpc_security_base",
        "@com_google_protobuf//:protobuf",
    ],
)

py_library(
    name = "skill_init_py",
    srcs = ["skill_init.py"],
    visibility = ["//visibility:public"],
    deps = [
        ":skill_repository_py",
        ":skill_service_impl_py",
        "//intrinsic/geometry/proto:geometry_service_py_pb2_grpc",
        "//intrinsic/motion_planning/proto/v1:motion_planner_service_py_pb2_grpc",
        "//intrinsic/skills/proto:skill_service_config_py_pb2",
        "//intrinsic/skills/proto:skill_service_py_pb2_grpc",
        "//intrinsic/world/proto:object_world_service_py_pb2_grpc",
        requirement("grpcio"),
        "@com_google_absl_py//absl/logging",
    ],
)

cc_library(
    name = "skill_service_config_utils",
    srcs = ["skill_service_config_utils.cc"],
    hdrs = ["skill_service_config_utils.h"],
    visibility = ["//visibility:public"],
    deps = [
        ":skill_proto_utils",
        "//intrinsic/assets/proto:id_cc_proto",
        "//intrinsic/assets/proto:status_spec_cc_proto",
        "//intrinsic/icon/release:file_helpers",
        "//intrinsic/skills/proto:skill_manifest_cc_proto",
        "//intrinsic/skills/proto:skill_service_config_cc_proto",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
        "@com_google_protobuf//:protobuf",
    ],
)

py_library(
    name = "skill_service_config_utils_py",
    srcs = ["skill_service_config_utils.py"],
    visibility = ["//visibility:public"],
    deps = [
        ":proto_utils",
        "//intrinsic/skills/proto:skill_manifest_py_pb2",
        "//intrinsic/skills/proto:skill_service_config_py_pb2",
        "@com_google_absl_py//absl/logging",
        "@com_google_protobuf//:protobuf_python",
    ],
)

cc_library(
    name = "skill_service_impl",
    srcs = ["skill_service_impl.cc"],
    hdrs = ["skill_service_impl.h"],
    deps = [
        ":equipment_utilities",
        ":error_utils",
        ":execute_context_impl",
        ":get_footprint_context_impl",
        ":preview_context_impl",
        ":runtime_data",
        ":skill_registry_client_interface",
        ":skill_repository",
        "//intrinsic/assets:id_utils",
        "//intrinsic/logging/proto:context_cc_proto",
        "//intrinsic/motion_planning:motion_planner_client",
        "//intrinsic/motion_planning/proto/v1:motion_planner_service_cc_grpc_proto",
        "//intrinsic/skills/cc:equipment_pack",
        "//intrinsic/skills/cc:skill_canceller",
        "//intrinsic/skills/cc:skill_interface",
        "//intrinsic/skills/cc:skill_logging_context",
        "//intrinsic/skills/proto:skill_service_cc_grpc_proto",
        "//intrinsic/skills/proto:skill_service_cc_proto",
        "//intrinsic/skills/proto:skills_cc_proto",
        "//intrinsic/util:proto_time",
        "//intrinsic/util/status:extended_status_cc_proto",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/status:status_conversion_rpc",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/status:status_macros_grpc",
        "//intrinsic/util/thread",
        "//intrinsic/world/objects:object_world_client",
        "//intrinsic/world/proto:object_world_service_cc_grpc_proto",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/functional:any_invocable",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_googleapis//google/longrunning:longrunning_cc_proto",
        "@com_google_googleapis//google/rpc:status_cc_proto",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "skill_registry_client",
    srcs = ["skill_registry_client.cc"],
    hdrs = ["skill_registry_client.h"],
    deps = [
        ":skill_registry_client_interface",
        "//intrinsic/skills/cc:client_common",
        "//intrinsic/skills/cc:equipment_pack",
        "//intrinsic/skills/internal/proto:behavior_tree_registry_internal_cc_grpc_proto",
        "//intrinsic/skills/internal/proto:behavior_tree_registry_internal_cc_proto",
        "//intrinsic/skills/internal/proto:skill_registry_internal_cc_grpc_proto",
        "//intrinsic/skills/internal/proto:skill_registry_internal_cc_proto",
        "//intrinsic/skills/proto:behavior_tree_registry_cc_grpc_proto",
        "//intrinsic/skills/proto:behavior_tree_registry_cc_proto",
        "//intrinsic/skills/proto:equipment_cc_proto",
        "//intrinsic/skills/proto:skill_registry_cc_grpc_proto",
        "//intrinsic/skills/proto:skill_registry_cc_proto",
        "//intrinsic/skills/proto:skill_registry_config_cc_proto",
        "//intrinsic/skills/proto:skills_cc_proto",
        "//intrinsic/util/grpc",
        "//intrinsic/util/status:annotate",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "skill_registry_client_interface",
    hdrs = ["skill_registry_client_interface.h"],
    deps = [
        "//intrinsic/skills/cc:equipment_pack",
        "//intrinsic/skills/proto:equipment_cc_proto",
        "//intrinsic/skills/proto:skill_registry_config_cc_proto",
        "//intrinsic/skills/proto:skills_cc_proto",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
    ],
)

cc_library(
    name = "skill_proto_utils",
    srcs = ["skill_proto_utils.cc"],
    hdrs = ["skill_proto_utils.h"],
    deps = [
        "//intrinsic/assets/proto:documentation_cc_proto",
        "//intrinsic/assets/proto:id_cc_proto",
        "//intrinsic/skills/proto:equipment_cc_proto",
        "//intrinsic/skills/proto:skill_manifest_cc_proto",
        "//intrinsic/skills/proto:skills_cc_proto",
        "//intrinsic/util/proto:source_code_info_view",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
        "@com_googlesource_code_re2//:re2",
    ],
)

go_library(
    name = "skillmanifest",
    srcs = ["skillmanifest.go"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:metadatautils",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "//intrinsic/util/proto:sourcecodeinfoview",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
    ],
)

go_test(
    name = "skillmanifest_test",
    srcs = ["skillmanifest_test.go"],
    data = [
        "//intrinsic/skills/build_defs/tests:no_op_skill_cc_manifest",
        "//intrinsic/skills/build_defs/tests:no_op_skill_py_manifest",
    ],
    library = ":skillmanifest",
    deps = [
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "//intrinsic/util/proto:protoio",
        "//intrinsic/util/proto:registryutil",
        "//intrinsic/util/testing:testio",
        "@org_golang_google_protobuf//proto",
    ],
)
