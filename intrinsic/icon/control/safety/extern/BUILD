# Copyright 2023 Intrinsic Innovation LLC

load("//intrinsic/platform:flatbuffers.bzl", "cc_flatbuffers_library", "flatbuffers_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

flatbuffers_library(
    name = "safety_status_fbs",
    srcs = [
        "safety_status.fbs",
    ],
)

cc_flatbuffers_library(
    name = "safety_status_fbs_cc",
    visibility = ["//visibility:public"],
    deps = [":safety_status_fbs"],
)
