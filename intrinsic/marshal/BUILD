# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "riegeli_proto_coder",
    hdrs = ["riegeli_proto_coder.h"],
    visibility = [
        "//intrinsic:intrinsic_internal",
    ],
    deps = [
        ":riegeli_coder",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log:check",
        "@com_google_protobuf//:protobuf",
        "@com_google_protobuf//src/google/protobuf/io",
        "@com_google_riegeli//riegeli/messages:serialize_message",
        "@com_google_riegeli//riegeli/records:record_reader",
        "@com_google_riegeli//riegeli/records:record_writer",
    ],
)

cc_library(
    name = "riegeli_coder",
    srcs = ["riegeli_coder.cc"],
    hdrs = ["riegeli_coder.h"],
    visibility = [
        "//intrinsic:intrinsic_internal",
    ],
    deps = [
        "//intrinsic/util:demangle",
        "//intrinsic/util/status:status_cc_proto",
        "//intrinsic/util/status:status_conversion_proto",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_riegeli//riegeli/bytes:string_reader",
        "@com_google_riegeli//riegeli/bytes:string_writer",
        "@com_google_riegeli//riegeli/endian:endian_reading",
        "@com_google_riegeli//riegeli/endian:endian_writing",
        "@com_google_riegeli//riegeli/records:record_reader",
        "@com_google_riegeli//riegeli/records:record_writer",
    ],
)
