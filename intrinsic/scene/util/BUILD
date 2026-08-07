# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

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
