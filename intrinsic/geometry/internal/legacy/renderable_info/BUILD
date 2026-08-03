# Copyright 2023 Intrinsic Innovation LLC

""" Intrinsic package geometry/legacy/renderable_info. """

load("@rules_cc//cc:cc_library.bzl", "cc_library")

cc_library(
    name = "renderable_info",
    srcs = ["renderable_info.cc"],
    hdrs = ["renderable_info.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    visibility = [
        "//intrinsic/geometry:__subpackages__",
    ],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/internal/util:eigen_serialization",
        "//intrinsic/marshal:riegeli_coder",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@com_google_riegeli//riegeli/records:record_reader",
        "@com_google_riegeli//riegeli/records:record_writer",
    ],
)
