# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_visibility = [
        "//intrinsic/platform/pubsub/connect/onprem/forwarding_service:__subpackages__",
    ],
)

proto_library(
    name = "forwarding_service_proto",
    srcs = ["forwarding_service.proto"],
    deps = [
        "//intrinsic/platform/pubsub/connect/cloud/proto/v1alpha1:endpoint_spec_proto",
    ],
)

py_proto_library(
    name = "forwarding_service_proto_py_pb2",
    visibility = [
        "//intrinsic/tools/inctl/cmd/pubsub:__pkg__",
    ],
    deps = [":forwarding_service_proto"],
)

go_proto_library(
    name = "forwarding_service_go_proto",
    importpath = "intrinsic/platform/pubsub/connect/onprem/forwarding_service/forwarding_service_go_proto",
    protos = [":forwarding_service_proto"],
    visibility = [
        "//intrinsic/platform/pubsub/connect/onprem/forwarding_service:__subpackages__",
        "//intrinsic/tools/inctl/cmd/pubsub:__pkg__",
    ],
    deps = [
        "//intrinsic/platform/pubsub/connect/cloud/proto/v1alpha1:endpoint_spec_go_proto",
    ],
)
