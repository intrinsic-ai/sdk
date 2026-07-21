# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")
load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "proto_source_code_info_transitive_descriptor_set")

package(default_visibility = ["//intrinsic/tools/inbuild:__subpackages__"])

proto_library(
    name = "custom_data_proto",
    srcs = ["custom_data.proto"],
    deps = ["//intrinsic/assets/data/proto/v1:referenced_data_proto"],
)

go_proto_library(
    name = "custom_data_go_proto",
    importpath = "intrinsic/tools/inbuild/cmd/data/test_data/custom_data_go_proto",
    protos = [":custom_data_proto"],
    deps = ["//intrinsic/assets/data/proto/v1:referenced_data_go_proto"],
)

proto_source_code_info_transitive_descriptor_set(
    name = "custom_data_fds",
    testonly = True,
    deps = [":custom_data_proto"],
)
