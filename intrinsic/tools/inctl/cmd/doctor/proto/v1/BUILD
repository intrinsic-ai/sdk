# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "report_proto",
    srcs = ["report.proto"],
)

go_proto_library(
    name = "report_go_proto",
    importpath = "intrinsic/tools/inctl/cmd/doctor/proto/v1/report_go_proto",
    protos = [":report_proto"],
)
