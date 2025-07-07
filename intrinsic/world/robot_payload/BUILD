# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:defs.bzl", "cc_library")

package(default_visibility = [
    "//visibility:public",
])

cc_library(
    name = "robot_payload",
    srcs = ["robot_payload.cc"],
    hdrs = ["robot_payload.h"],
    deps = [
        ":robot_payload_base",
        "//intrinsic/eigenmath",
        "//intrinsic/kinematics:validate_link_parameters",
        "//intrinsic/math:almost_equals",
        "//intrinsic/math:inertia_utils",
        "//intrinsic/math:pose3",
        "//intrinsic/math:proto_conversion",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/proto:robot_payload_cc_proto",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "robot_payload_base",
    srcs = ["robot_payload_base.cc"],
    hdrs = ["robot_payload_base.h"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/math:almost_equals",
        "//intrinsic/math:pose3",
        "//intrinsic/world/proto:robot_payload_cc_proto",
    ],
)
