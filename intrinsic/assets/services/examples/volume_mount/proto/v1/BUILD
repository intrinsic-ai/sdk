# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

package(
    default_visibility = ["//intrinsic/assets/services/examples/volume_mount:__subpackages__"],
)

proto_library(
    name = "volume_mount_proto",
    srcs = ["volume_mount.proto"],
)

go_proto_library(
    name = "volume_mount_go_proto",
    protos = [":volume_mount_proto"],
)

go_grpc_library(
    name = "volume_mount_go_grpc_proto",
    protos = [":volume_mount_proto"],
    deps = [":volume_mount_go_proto"],
)
