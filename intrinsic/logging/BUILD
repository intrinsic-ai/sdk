# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

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
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/random",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
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
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
        "@abseil-cpp//absl/types:span",
        "@com_github_grpc_grpc//:grpc++",
    ],
)
