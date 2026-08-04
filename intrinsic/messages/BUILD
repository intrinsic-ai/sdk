# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

licenses(["unencumbered"])

proto_library(
    name = "messages_proto",
    srcs = [
        "geometry_types.proto",
    ],
)

go_proto_library(
    name = "messages_go_proto",
    importpath = "intrinsic/messages/messages_go_proto",
    protos = [":messages_proto"],
)

py_proto_library(
    name = "messages_py_pb2",
    deps = [":messages_proto"],
)

cc_proto_library(
    name = "messages_cc_proto",
    deps = [":messages_proto"],
)
