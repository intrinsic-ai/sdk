# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "eigen_serialization",
    srcs = [
        "eigen_serialization.cc",
    ],
    hdrs = [
        "eigen_serialization.h",
    ],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/eigenmath:pose2",
        "//intrinsic/marshal:riegeli_coder",
        "//intrinsic/math:pose3",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_riegeli//riegeli/records:record_reader",
        "@com_google_riegeli//riegeli/records:record_writer",
    ],
)
