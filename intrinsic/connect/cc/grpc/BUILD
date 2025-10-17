# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "channel",
    srcs = ["channel.cc"],
    hdrs = ["channel.h"],
    deps = [
        "//intrinsic/icon/release:grpc_time_support",
        "//intrinsic/util/grpc:limits",
        "//intrinsic/util/status:status_conversion_grpc",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/strings:string_view",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
        "@com_github_grpc_grpc//:grpc_security_base",
        "@com_github_grpc_grpc//src/proto/grpc/health/v1:health_cc_grpc",
        "@com_github_grpc_grpc//src/proto/grpc/health/v1:health_cc_proto",
    ],
)
