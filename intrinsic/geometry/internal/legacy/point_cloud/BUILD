# Copyright 2023 Intrinsic Innovation LLC

""" Intrinsic package geometry/legacy/point_cloud. """

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//intrinsic/geometry:__subpackages__"])

exports_files(
    srcs = [
        "test_data/Raw_IPS.ply",
        "test_data/example_point_cloud.ply",
    ],
    visibility = [
        "//intrinsic/geometry:__subpackages__",
        "//intrinsic/world/service/benchmarks:__subpackages__",
    ],
)

cc_library(
    name = "point_cloud_riegeli_coder",
    srcs = ["point_cloud_riegeli_coder.cc"],
    hdrs = ["point_cloud_riegeli_coder.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/proto/v1:point_cloud_cc_proto",
        "//intrinsic/geometry/shapes",
        "//intrinsic/marshal:riegeli_proto_coder",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)
