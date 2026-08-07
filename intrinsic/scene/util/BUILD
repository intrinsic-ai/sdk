# Copyright 2023 Intrinsic Innovation LLC

load("@pybind11_bazel//:build_defs.bzl", "pybind_extension")
load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")
load("@rules_python//python:defs.bzl", "py_test")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "scene_object_gzf",
    srcs = ["scene_object_gzf.cc"],
    hdrs = ["scene_object_gzf.h"],
    deps = [
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/world/gzfile",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//src/google/protobuf/io",
    ],
)

cc_test(
    name = "scene_object_gzf_test",
    srcs = ["scene_object_gzf_test.cc"],
    deps = [
        ":scene_object_gzf",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "//intrinsic/world/gzfile",
        "@abseil-cpp//absl/status",
        "@or_tools//ortools/base:path",
    ],
)

cc_library(
    name = "geometry_update_util",
    srcs = ["geometry_update_util.cc"],
    hdrs = ["geometry_update_util.h"],
    deps = [
        "//intrinsic/scene/proto/v1:scene_object_updates_cc_proto",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world:geometry_types",
        "//intrinsic/world/proto:geometry_component_cc_proto",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "scene_object_updates",
    srcs = [
        "scene_object_updates.cc",
        "scene_object_updates_internal.cc",
    ],
    hdrs = [
        "scene_object_updates.h",
        "scene_object_updates_internal.h",
    ],
    deps = [
        ":geometry_update_util",
        ":object_user_data",
        "//intrinsic/eigenmath",
        "//intrinsic/kinematics/types:joint_limits_cc_proto",
        "//intrinsic/math:pose3",
        "//intrinsic/math:proto_conversion",
        "//intrinsic/math/proto:pose_cc_proto",
        "//intrinsic/scene/proto/v1:collision_rules_cc_proto",
        "//intrinsic/scene/proto/v1:entity_cc_proto",
        "//intrinsic/scene/proto/v1:object_properties_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_updates_cc_proto",
        "//intrinsic/scene/validate:scene_object_validation",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world:geometry_types",
        "//intrinsic/world/proto:kinematics_component_cc_proto",
        "@abseil-cpp//absl/algorithm:container",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/container:flat_hash_set",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "scene_object_updates_test",
    srcs = ["scene_object_updates_test.cc"],
    deps = [
        ":scene_object_updates",
        "//intrinsic/scene/proto/v1:collision_rules_cc_proto",
        "//intrinsic/scene/proto/v1:object_properties_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_updates_cc_proto",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "geometry_update_util_test",
    srcs = ["geometry_update_util_test.cc"],
    deps = [
        ":geometry_update_util",
        "//intrinsic/scene/proto/v1:scene_object_updates_cc_proto",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "//intrinsic/world/proto:geometry_component_cc_proto",
        "@abseil-cpp//absl/status",
    ],
)

cc_library(
    name = "scene_object_creation",
    srcs = ["scene_object_creation.cc"],
    hdrs = ["scene_object_creation.h"],
    deps = [
        "//intrinsic/geometry/api:affine_transform_of_geometry",
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/proto:geometry_storage_refs_cc_proto",
        "//intrinsic/geometry/storage:geometry_serializer",
        "//intrinsic/kinematics/types:joint_limits",
        "//intrinsic/kinematics/types:joint_limits_cc_proto",
        "//intrinsic/scene/proto/v1:entity_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_updates_cc_proto",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world:geometry_types",
        "//intrinsic/world/component:geometry_component",
        "//intrinsic/world/proto:geometry_component_cc_proto",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
    ],
)

cc_test(
    name = "scene_object_creation_test",
    srcs = ["scene_object_creation_test.cc"],
    deps = [
        ":scene_object_creation",
        "//intrinsic/icon/utils:realtime_status_matchers",
        "//intrinsic/kinematics/types:joint_limits",
        "//intrinsic/kinematics/types:joint_limits_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_updates_cc_proto",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/util/testing:gtest_wrapper_main",
    ],
)

cc_library(
    name = "object_user_data",
    srcs = ["object_user_data.cc"],
    hdrs = ["object_user_data.h"],
    deps = [
        "//intrinsic/scene/conversion:user_data",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_updates_cc_proto",
        "//intrinsic/util/proto:descriptors",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "object_user_data_test",
    srcs = ["object_user_data_test.cc"],
    deps = [
        ":object_user_data",
        "//intrinsic/util/proto:descriptors",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "scene_object_from_asset",
    srcs = ["scene_object_from_asset.cc"],
    hdrs = ["scene_object_from_asset.h"],
    deps = [
        ":scene_object_updates",
        "//intrinsic/assets:id_utils",
        "//intrinsic/assets/proto:id_cc_proto",
        "//intrinsic/assets/proto:installed_assets_cc_grpc_proto",
        "//intrinsic/assets/proto:installed_assets_cc_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_cc_proto",
        "//intrinsic/connect/cc/grpc:channel",
        "//intrinsic/scene:constants",
        "//intrinsic/scene/config:scene_object_config",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_config_cc_proto",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/status:status_macros_grpc",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
    ],
)

cc_test(
    name = "scene_object_from_asset_test",
    srcs = ["scene_object_from_asset_test.cc"],
    deps = [
        ":scene_object_from_asset",
        "//intrinsic/assets/proto:installed_assets_cc_grpc_proto",
        "//intrinsic/assets/proto:installed_assets_cc_proto",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_protobuf//:protobuf",
    ],
)

pybind_extension(
    name = "scene_object_from_asset_py",
    srcs = ["scene_object_from_asset_py.cc"],
    deps = [
        ":scene_object_from_asset",
        "//intrinsic/assets/proto:id_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "@abseil-cpp//absl/strings",
        "@pybind11_abseil//pybind11_abseil:absl_casters",
        "@pybind11_abseil//pybind11_abseil:status_casters",
        "@pybind11_protobuf//pybind11_protobuf:native_proto_caster",
    ],
)

py_test(
    name = "scene_object_from_asset_py_test",
    srcs = ["scene_object_from_asset_py_test.py"],
    deps = [
        ":scene_object_from_asset_py",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_py_proto",
        "//intrinsic/scene/proto/v1:scene_object_config_py_pb2",
        "//intrinsic/scene/proto/v1:scene_object_py_pb2",
        "//intrinsic/scene/proto/v1:scene_object_updates_py_pb2",
        "//intrinsic/util/status:status_assertions",
        "@com_google_absl_py//absl/testing:absltest",
        "@pybind11_abseil//pybind11_abseil:import_status_module",
    ],
)

cc_library(
    name = "scene_object_geo_refs",
    srcs = ["scene_object_geo_refs.cc"],
    hdrs = ["scene_object_geo_refs.h"],
    deps = [
        "//intrinsic/geometry/proto:geometry_storage_refs_cc_proto",
        "//intrinsic/geometry/proto/v1:geometry_cc_proto",
        "//intrinsic/geometry/proto/v1:geometry_storage_refs_cc_proto",
        "//intrinsic/geometry/proto/v1:transformed_geometry_cc_proto",
        "//intrinsic/scene/proto/v1:entity_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/world/proto:geometry_component_cc_proto",
        "@abseil-cpp//absl/container:flat_hash_set",
        "@abseil-cpp//absl/strings:string_view",
    ],
)

cc_test(
    name = "scene_object_geo_refs_test",
    srcs = ["scene_object_geo_refs_test.cc"],
    deps = [
        ":scene_object_geo_refs",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/util/testing:gtest_wrapper_main",
    ],
)
