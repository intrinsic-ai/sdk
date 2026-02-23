# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "eigenmath",
    hdrs = [
        "types.h",
    ],
    deps = ["@com_gitlab_libeigen_eigen//:eigen"],
)

cc_library(
    name = "rotation_utils",
    hdrs = ["rotation_utils.h"],
    deps = [
        ":eigenmath",
        "@abseil-cpp//absl/log:check",
        "@com_gitlab_libeigen_eigen//:eigen",
    ],
)

cc_library(
    name = "so3",
    hdrs = ["so3.h"],
    deps = [
        ":eigenmath",
        ":rotation_utils",
        "//intrinsic/icon/testing:realtime_annotations",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/util/status:status_builder",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@com_gitlab_libeigen_eigen//:eigen",
    ],
)

cc_library(
    name = "clamp",
    hdrs = [
        "clamp.h",
    ],
    deps = [
        "@abseil-cpp//absl/base:core_headers",
    ],
)
