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
