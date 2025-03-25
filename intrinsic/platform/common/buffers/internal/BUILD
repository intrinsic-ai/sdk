# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:defs.bzl", "cc_library")

package(default_visibility = [
    "//intrinsic/platform/common/buffers:__subpackages__",
])

cc_library(
    name = "event_fd",
    srcs = ["event_fd.cc"],
    hdrs = ["event_fd.h"],
    deps = ["@com_google_absl//absl/log:check"],
)
