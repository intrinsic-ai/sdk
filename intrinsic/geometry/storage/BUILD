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
    name = "gzf_storage",
    srcs = ["gzf_storage.cc"],
    hdrs = ["gzf_storage.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        ":geometry_deserializer",
        ":geometry_library",
        ":geometry_serializer",
        "//intrinsic/geometry/api:exact_geometry",
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/api:geometry_fingerprint",
        "//intrinsic/geometry/api:geometry_options",
        "//intrinsic/geometry/api:io",
        "//intrinsic/geometry/api:renderable",
        "//intrinsic/geometry/api:renderable_generation",
        "//intrinsic/geometry/proto:geometry_cc_proto",
        "//intrinsic/geometry/proto:geometry_storage_refs_cc_proto",
        "//intrinsic/geometry/proto/v1:geometry_storage_refs_cc_proto",
        "//intrinsic/geometry/proto/v1:material_cc_proto",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/gzfile",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "in_memory_storage",
    srcs = ["in_memory_storage.cc"],
    hdrs = ["in_memory_storage.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        ":geometry_deserializer",
        ":geometry_library",
        ":geometry_serializer",
        "//intrinsic/geometry/api:exact_geometry",
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/api:geometry_fingerprint",
        "//intrinsic/geometry/api:geometry_options",
        "//intrinsic/geometry/api:renderable",
        "//intrinsic/geometry/api:renderable_generation",
        "//intrinsic/geometry/proto:geometry_storage_refs_cc_proto",
        "//intrinsic/geometry/proto/v1:geometry_storage_refs_cc_proto",
        "//intrinsic/geometry/proto/v1:material_cc_proto",
        "//intrinsic/util/status:status_builder",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/hashing",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/synchronization",
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
