# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:defs.bzl", "cc_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

cc_library(
    name = "builtins",
    hdrs = ["builtins.h"],
)

cc_library(
    name = "id_types",
    hdrs = ["id_types.h"],
    deps = ["//intrinsic/production/external/intops:strong_int"],
)

cc_library(
    name = "slot_part_map",
    srcs = ["slot_part_map.cc"],
    hdrs = ["slot_part_map.h"],
    deps = [
        "//intrinsic/icon/proto/v1:types_cc_proto",
        "@com_google_absl//absl/container:btree",
        "@com_google_absl//absl/container:flat_hash_map",
    ],
)

cc_library(
    name = "part_properties",
    srcs = ["part_properties.cc"],
    hdrs = ["part_properties.h"],
    deps = [
        "//intrinsic/icon/proto/v1:service_cc_proto",
        "//intrinsic/production/external/intops:strong_int",
        "@com_google_absl//absl/container:flat_hash_map",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/status:statusor",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/time",
    ],
)

cc_library(
    name = "state_variable_path_constants",
    hdrs = [
        "state_variable_path_constants.h",
    ],
)

cc_library(
    name = "state_variable_path_util",
    srcs = [
        "state_variable_path_util.cc",
    ],
    hdrs = [
        "state_variable_path_util.h",
    ],
    deps = [
        ":state_variable_path_constants",
        "@com_google_absl//absl/algorithm:container",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/types:span",
    ],
)
