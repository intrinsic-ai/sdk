# Copyright 2023 Intrinsic Innovation LLC

# This package holds wrappers that make ICON classes available to users of the RTCL C API
# (see parent directory).

load("@rules_cc//cc:defs.bzl", "cc_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

cc_library(
    name = "string_wrapper",
    srcs = ["string_wrapper.cc"],
    hdrs = ["string_wrapper.h"],
    deps = [
        "//intrinsic/icon/control/c_api:icon_c_api",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "streaming_io_wrapper",
    hdrs = ["streaming_io_wrapper.h"],
    deps = [
        ":string_wrapper",
        "//intrinsic/icon/control/c_api:convert_c_realtime_status",
        "//intrinsic/icon/control/c_api:icon_c_api",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_or",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)
