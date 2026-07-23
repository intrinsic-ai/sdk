# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(
    default_visibility =
        [
            "//intrinsic/scene:__subpackages__",
            "//intrinsic/simulation:__subpackages__",
            "//intrinsic/world/conversion/sdf:__subpackages__",
        ],
)

cc_library(
    name = "custom_tags",
    hdrs = ["custom_tags.h"],
    deps = ["@abseil-cpp//absl/strings:string_view"],
)

cc_library(
    name = "separators",
    hdrs = ["separators.h"],
    deps = ["@abseil-cpp//absl/strings:string_view"],
)

cc_library(
    name = "sdf_sensor_pose",
    srcs = ["sdf_sensor_pose.cc"],
    hdrs = ["sdf_sensor_pose.h"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/math:pose3",
        "@sdformat",
    ],
)

cc_library(
    name = "xml_utils",
    srcs = ["xml_utils.cc"],
    hdrs = ["xml_utils.h"],
    deps = [
        "//intrinsic/util/status:status_builder",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
        "@sdformat",
        "@tinyxml2",
    ],
)

cc_library(
    name = "single_element_loaders",
    testonly = 1,
    srcs = ["single_element_loaders.cc"],
    hdrs = ["single_element_loaders.h"],
    deps = [
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
        "@sdformat",
    ],
)
