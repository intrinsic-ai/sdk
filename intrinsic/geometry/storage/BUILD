# Copyright 2023 Intrinsic Innovation LLC

""" Intrinsic package geometry/storage. """

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "geometry_deserializer",
    hdrs = ["geometry_deserializer.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/api:geometry_options",
        "//intrinsic/geometry/proto:geometry_storage_refs_cc_proto",
        "//intrinsic/geometry/proto/v1:geometry_storage_refs_cc_proto",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "geometry_serializer",
    hdrs = ["geometry_serializer.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/proto:geometry_storage_refs_cc_proto",
        "//intrinsic/geometry/proto/v1:geometry_storage_refs_cc_proto",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "geometry_library",
    hdrs = ["geometry_library.h"],
    deps = [
        ":geometry_deserializer",
        ":geometry_serializer",
    ],
)
