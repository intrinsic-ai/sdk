# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "large_mesh",
    srcs = ["large_mesh.cc"],
    hdrs = ["large_mesh.h"],
    deps = [
        "//intrinsic/geometry/api:axis_aligned_bounding_box_3d",
        "//intrinsic/geometry/api:compute_axis_aligned_bounding_box_3d",
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/hashing",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
    ],
)

cc_test(
    name = "large_mesh_test",
    srcs = ["large_mesh_test.cc"],
    deps = [
        ":large_mesh",
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/shapes",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings:str_format",
    ],
)

cc_library(
    name = "scene_object_validation",
    srcs = ["scene_object_validation.cc"],
    hdrs = ["scene_object_validation.h"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/kinematics/types:cartesian_limits",
        "//intrinsic/math:pose3",
        "//intrinsic/math:proto_conversion",
        "//intrinsic/scene:constants",
        "//intrinsic/scene/conversion:object_properties_conversion",
        "//intrinsic/scene/proto/v1:collision_rules_cc_proto",
        "//intrinsic/scene/proto/v1:entity_cc_proto",
        "//intrinsic/scene/proto/v1:object_properties_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/proto:kinematics_component_cc_proto",
        "@abseil-cpp//absl/algorithm:container",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/container:flat_hash_set",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "scene_object_validate_geo",
    srcs = ["scene_object_validate_geo.cc"],
    hdrs = ["scene_object_validate_geo.h"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/api:renderable",
        "//intrinsic/geometry/api:validate_mesh",
        "//intrinsic/geometry/proto/v1:geometric_transform_cc_proto",
        "//intrinsic/geometry/proto/v1:geometry_cc_proto",
        "//intrinsic/geometry/proto/v1:transformed_geometry_cc_proto",
        "//intrinsic/geometry/storage:geometry_deserializer",
        "//intrinsic/math:proto_conversion",
        "//intrinsic/scene/proto/v1:entity_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/proto:geometry_component_cc_proto",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
    ],
)

cc_test(
    name = "scene_object_validation_test",
    srcs = ["scene_object_validation_test.cc"],
    deps = [
        ":scene_object_validation",
        "//intrinsic/scene:constants",
        "//intrinsic/scene/proto/v1:collision_rules_cc_proto",
        "//intrinsic/scene/proto/v1:entity_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/base:no_destructor",
        "@abseil-cpp//absl/container:flat_hash_set",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings:string_view",
        "@com_google_fuzztest//fuzztest",
    ],
)

cc_test(
    name = "scene_object_validate_geo_test",
    srcs = ["scene_object_validate_geo_test.cc"],
    deps = [
        ":scene_object_validate_geo",
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/storage:geometry_deserializer",
        "//intrinsic/geometry/testing:glb_utils",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/testing:gtest_wrapper_main",
    ],
)
