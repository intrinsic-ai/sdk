# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "shared_memory_lockstep",
    srcs = [
        "shared_memory_lockstep.cc",
    ],
    hdrs = [
        "shared_memory_lockstep.h",
    ],
    deps = [
        "//intrinsic/icon/interprocess/shared_memory_manager",
        "//intrinsic/icon/interprocess/shared_memory_manager:domain_socket_utils",
        "//intrinsic/icon/interprocess/shared_memory_manager:memory_segment",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/thread:lockstep",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)
