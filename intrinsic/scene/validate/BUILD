# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

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
