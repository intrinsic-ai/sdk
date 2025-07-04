# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:defs.bzl", "cc_library")

# Component to log structured data with context.
package(default_visibility = ["//visibility:public"])

exports_files(
    srcs = [
        "logger.yaml",
        "logging_disk.yaml",
        "timescaledb.yaml",
    ],
    visibility = ["//intrinsic:__subpackages__"],
)

cc_library(
    name = "data_logger_client",
    srcs = ["data_logger_client.cc"],
    hdrs = ["data_logger_client.h"],
    deps = [
        ":structured_logging_client",
        "//intrinsic/logging/proto:log_item_cc_proto",
        "//intrinsic/logging/proto:logger_service_cc_grpc",
        "//intrinsic/logging/proto:logger_service_cc_proto",
        "//intrinsic/util/grpc",
        "//intrinsic/util/status:status_macros",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/log:check",
        "@com_google_absl//absl/random",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/time",
    ],
)

cc_library(
    name = "structured_logging_client",
    srcs = ["structured_logging_client.cc"],
    hdrs = ["structured_logging_client.h"],
    deps = [
        "//intrinsic/logging/proto:log_item_cc_proto",
        "//intrinsic/logging/proto:logger_service_cc_grpc",
        "//intrinsic/logging/proto:logger_service_cc_proto",
        "//intrinsic/util:proto_time",
        "//intrinsic/util/grpc",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/status:status_macros",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_absl//absl/container:flat_hash_map",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/status:statusor",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/time",
        "@com_google_absl//absl/types:span",
    ],
)
