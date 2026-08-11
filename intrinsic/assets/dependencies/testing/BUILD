# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

proto_library(
    name = "test_configs_proto",
    srcs = ["test_configs.proto"],
    visibility = ["//intrinsic/assets/dependencies:__pkg__"],
    deps = [
        "//intrinsic/assets/proto:field_metadata_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_proto",
        "@com_google_protobuf//:any_proto",
    ],
)

cc_proto_library(
    name = "test_configs_cc_proto",
    visibility = ["//intrinsic/assets/dependencies:__pkg__"],
    deps = [":test_configs_proto"],
)

go_proto_library(
    name = "test_configs_go_proto",
    testonly = True,
    importpath = "intrinsic/assets/dependencies/testing/test_configs_go_proto",
    protos = [":test_configs_proto"],
    visibility = [
        "//intrinsic/assets:__subpackages__",
    ],
    deps = [
        "//intrinsic/assets/proto:field_metadata_go_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_go_proto",
    ],
)

proto_library(
    name = "multi_data_dependency_config_proto",
    srcs = ["multi_data_dependency_config.proto"],
    visibility = ["//intrinsic/assets/dependencies:__pkg__"],
    deps = [
        "//intrinsic/assets/proto:field_metadata_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_proto",
    ],
)

go_proto_library(
    name = "multi_data_dependency_config_go_proto",
    testonly = True,
    importpath = "intrinsic/assets/dependencies/testing/multi_data_dependency_config_go_proto",
    protos = [":multi_data_dependency_config_proto"],
    visibility = [
        "//intrinsic/assets:__subpackages__",
    ],
    deps = [
        "//intrinsic/assets/proto:field_metadata_go_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_go_proto",
    ],
)
