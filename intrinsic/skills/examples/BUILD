# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")
load("//intrinsic/skills/build_defs:skill.bzl", "cc_skill", "skill_manifest")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "calculate_skill_proto",
    srcs = ["calculate_skill.proto"],
    deps = [
        "//intrinsic/assets/proto:field_metadata_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_proto",
        "//intrinsic/assets/services/examples/calcserver:calc_server_proto",
    ],
)

cc_proto_library(
    name = "calculate_skill_cc_proto",
    deps = [":calculate_skill_proto"],
)

cc_library(
    name = "calculate_skill_cc",
    srcs = ["calculate_skill.cc"],
    hdrs = ["calculate_skill.h"],
    deps = [
        ":calculate_skill_cc_proto",
        "//intrinsic/assets/dependencies:utils_cc",
        "//intrinsic/assets/proto/v1:resolved_dependency_cc_proto",
        "//intrinsic/assets/services/examples/calcserver:calc_server_cc_grpc_proto",
        "//intrinsic/assets/services/examples/calcserver:calc_server_cc_proto",
        "//intrinsic/skills/cc:skill_interface",
        "//intrinsic/skills/proto:skill_service_cc_proto",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status:statusor",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "calculate_skill_test",
    srcs = ["calculate_skill_test.cc"],
    deps = [
        ":calculate_skill_cc",
        ":calculate_skill_cc_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_cc_proto",
        "//intrinsic/assets/services/examples/calcserver:calc_server",
        "//intrinsic/assets/services/examples/calcserver:calc_server_cc_proto",
        "//intrinsic/skills/cc:skill_interface",
        "//intrinsic/skills/testing:skill_test_utils_cc",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@com_google_protobuf//:protobuf_lite",
    ],
)

skill_manifest(
    name = "calculate_skill_manifest",
    src = "calculate_skill_manifest.textproto",
    deps = [":calculate_skill_proto"],
)

cc_skill(
    name = "calculate_skill",
    manifest = ":calculate_skill_manifest",
    deps = [
        ":calculate_skill_cc",
        ":calculate_skill_cc_proto",
    ],
)
