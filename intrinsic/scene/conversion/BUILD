# Copyright 2023 Intrinsic Innovation LLC

# Libraries for conversion to Scene Object data structures.

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "object_properties_conversion",
    srcs = ["object_properties_conversion.cc"],
    hdrs = ["object_properties_conversion.h"],
    visibility = [
        "//intrinsic/scene/sdf:__subpackages__",
        "//intrinsic/scene/util:__subpackages__",
        "//intrinsic/scene/validate:__subpackages__",
        "//intrinsic/world:__subpackages__",
    ],
    deps = [
        "//intrinsic/icon/proto:eigen_conversion",
        "//intrinsic/kinematics/types:cartesian_limits",
        "//intrinsic/scene/proto/v1:object_properties_cc_proto",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "scene_object_model_utils",
    srcs = ["scene_object_model_utils.cc"],
    hdrs = ["scene_object_model_utils.h"],
    deps = [
        "//intrinsic/scene/proto/v1:entity_cc_proto",
        "//intrinsic/world/proto:kinematics_component_cc_proto",
        "@abseil-cpp//absl/container:flat_hash_set",
    ],
)

cc_library(
    name = "user_data",
    srcs = ["user_data.cc"],
    hdrs = ["user_data.h"],
    deps = [
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/util/proto:descriptors",
        "//intrinsic/util/proto:dynamic_message_parser",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@com_google_protobuf//:protobuf",
    ],
)
