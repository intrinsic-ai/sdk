# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = [
    "//visibility:public",
])

cc_library(
    name = "gzfile",
    srcs = ["gzfile.cc"],
    hdrs = [
        "chunk_entry.h",
        "gzfile.h",
    ],
    visibility = [
        "//intrinsic:intrinsic_internal",
    ],
    deps = [
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/types:optional",
        "@libzip",
        "@or_tools//ortools/base:file",
        "@or_tools//ortools/base:path",
        "@or_tools//ortools/base:strong_int",
    ],
)
