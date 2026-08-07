# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "collision_surface_conversion",
    srcs = ["collision_surface_conversion.cc"],
    hdrs = ["collision_surface_conversion.h"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/math:proto_conversion",
        "//intrinsic/scene/sdf:sdf_util",
        "//intrinsic/scene/sdf:xml_utils",
        "//intrinsic/simulation/gazebo:type_conversion",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/proto:physics_component_cc_proto",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@gz-math",
        "@sdformat",
    ],
)

cc_library(
    name = "geometry_component_from_sdf",
    srcs = ["geometry_component_from_sdf.cc"],
    hdrs = ["geometry_component_from_sdf.h"],
    deps = [
        "//intrinsic/geometry/api:affine_transform_of_geometry",
        "//intrinsic/geometry/api:io",
        "//intrinsic/geometry/api:material",
        "//intrinsic/scene/sdf:sdf_path_resolver",
        "//intrinsic/scene/sdf:sdf_util",
        "//intrinsic/util/status:status_builder",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world:geometry_types",
        "//intrinsic/world/component:geometry_component",
        "//intrinsic/world/proto:geometry_component_cc_proto",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@sdformat",
    ],
)

cc_library(
    name = "joint_conversion",
    srcs = ["joint_conversion.cc"],
    hdrs = ["joint_conversion.h"],
    deps = [
        ":kinematics_component_from_sdf",
        ":sensor_conversion",
        "//intrinsic/math:pose3",
        "//intrinsic/scene/sdf:sdf_util",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/component:kinematics_component",
        "//intrinsic/world/proto:kinematics_component_cc_proto",
        "@abseil-cpp//absl/status:statusor",
        "@sdformat",
    ],
)

cc_library(
    name = "kinematics_component_from_sdf",
    srcs = ["kinematics_component_from_sdf.cc"],
    hdrs = ["kinematics_component_from_sdf.h"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/scene:constants",
        "//intrinsic/scene/sdf:custom_tags",
        "//intrinsic/scene/sdf:sdf_util",
        "//intrinsic/simulation/gazebo:type_conversion",
        "//intrinsic/util/status:ret_check",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/component:kinematics_component",
        "//intrinsic/world/hashing",
        "//intrinsic/world/proto:kinematics_component_cc_proto",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@gz-math",
        "@sdformat",
    ],
)

cc_library(
    name = "physics_component_from_sdf",
    srcs = ["physics_component_from_sdf.cc"],
    hdrs = ["physics_component_from_sdf.h"],
    deps = [
        ":collision_surface_conversion",
        "//intrinsic/eigenmath",
        "//intrinsic/math:pose3",
        "//intrinsic/simulation/gazebo:type_conversion",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/component:physics_component",
        "//intrinsic/world/proto:physics_component_cc_proto",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status:statusor",
        "@gz-math",
        "@sdformat",
    ],
)

cc_library(
    name = "sensor_conversion",
    srcs = ["sensor_conversion.cc"],
    hdrs = ["sensor_conversion.h"],
    deps = [
        "//intrinsic/math:pose3",
        "//intrinsic/scene/sdf:sdf_sensor_pose",
        "//intrinsic/scene/sdf:sdf_util",
        "//intrinsic/scene/sdf:xml_utils",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/component:sensor_component",
        "//intrinsic/world/hashing",
        "//intrinsic/world/proto:sensor_component_cc_proto",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@sdformat",
    ],
)
