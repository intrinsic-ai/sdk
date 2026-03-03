# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@rules_cc//cc:cc_binary.bzl", "cc_binary")
load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")
load("@rules_python//python:defs.bzl", "py_library", "py_test")
load("//bazel:cc_oci_image.bzl", "cc_oci_image")
load("//intrinsic/util/proto/build_defs:proto_data.bzl", "proto_data")

package(default_visibility = ["//intrinsic/apps/common/resources:__subpackages__"])

proto_data(
    name = "status_specs_message",
    src = "//intrinsic/proto_tools/registry:proto_registry.status_specs.pbtxt",
    proto_deps = [
        "//intrinsic/assets/proto:status_spec_proto",
    ],
    proto_name = "intrinsic_proto.assets.StatusSpecs",
)

cc_library(
    name = "extended_status_codes",
    hdrs = ["extended_status_codes.h"],
)

cc_library(
    name = "file_descriptor_set_name",
    srcs = ["file_descriptor_set_name.cc"],
    hdrs = ["file_descriptor_set_name.h"],
    deps = [
        "//intrinsic/proto_tools/proto:proto_registry_internal_cc_proto",
        "//intrinsic/util:string_type",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
    ],
)

cc_test(
    name = "file_descriptor_set_name_test",
    srcs = ["file_descriptor_set_name_test.cc"],
    deps = [
        ":file_descriptor_set_name",
        "//intrinsic/util/testing:gtest_wrapper_main",
    ],
)

cc_library(
    name = "proto_registry_service",
    srcs = ["proto_registry_service.cc"],
    hdrs = ["proto_registry_service.h"],
    deps = [
        ":file_descriptor_set_name",
        ":resolver",
        ":resolver_asset",
        "//intrinsic/icon/release:source_location",
        "//intrinsic/proto_tools/compatibility:proto_compat",
        "//intrinsic/proto_tools/proto:proto_registry_cc_grpc_proto",
        "//intrinsic/proto_tools/proto:proto_registry_cc_proto",
        "//intrinsic/util:lru_cache",
        "//intrinsic/util/proto:descriptor_pools",
        "//intrinsic/util/proto:descriptors",
        "//intrinsic/util/proto:parsed_type_url",
        "//intrinsic/util/proto:type_url",
        "//intrinsic/util/status:extended_status_cc_proto",
        "//intrinsic/util/status:get_extended_status",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/status:status_macros_grpc",
        "//intrinsic/util/status:status_specs",
        "//intrinsic/util/thread",
        "//intrinsic/util/thread:stop_token",
        "//third_party/imported/cpp_libraries/clock",
        "@abseil-cpp//absl/algorithm:container",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/base:nullability",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "proto_registry_service_test",
    srcs = ["proto_registry_service_test.cc"],
    deps = [
        ":file_descriptor_set_name",
        ":proto_registry_service",
        ":resolver",
        ":resolver_mock",
        ":resolver_test_fixture",
        "//intrinsic/executive/clips_cpp:extended_status_matchers",
        "//intrinsic/proto_tools/proto:proto_registry_cc_proto",
        "//intrinsic/util/proto:compile_utils",
        "//intrinsic/util/proto:dynamic_message_parser",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/util/proto:type_url",
        "//intrinsic/util/status:extended_status_cc_proto",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/status:status_specs",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "//third_party/imported/cpp_libraries/clock",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "resolver",
    srcs = [
        "resolver.cc",
        "resolver_skill.cc",
    ],
    hdrs = [
        "resolver.h",
        "resolver_skill.h",
    ],
    deps = [
        ":file_descriptor_set_name",
        "//intrinsic/assets:id_utils",
        "//intrinsic/skills/internal:skill_registry_client_interface",
        "//intrinsic/skills/proto:skills_cc_proto",
        "//intrinsic/util/proto:parsed_type_url",
        "//intrinsic/util/proto:type_url",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/status:status_specs",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/container:flat_hash_set",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "resolver_test",
    srcs = ["resolver_test.cc"],
    deps = [
        ":resolver",
        ":resolver_mock",
        ":resolver_test_fixture",
        "//intrinsic/util/proto:parsed_type_url",
        "//intrinsic/util/proto:type_url",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "resolver_asset",
    srcs = ["resolver_asset.cc"],
    hdrs = ["resolver_asset.h"],
    deps = [
        ":file_descriptor_set_name",
        ":resolver",
        "//intrinsic/assets:id_utils",
        "//intrinsic/assets/proto:installed_assets_cc_grpc_proto",
        "//intrinsic/assets/proto:installed_assets_cc_proto",
        "//intrinsic/assets/proto:view_cc_proto",
        "//intrinsic/skills/proto:skills_cc_proto",
        "//intrinsic/util/grpc",
        "//intrinsic/util/proto:parsed_type_url",
        "//intrinsic/util/proto:type_url",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/status:status_specs",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/container:flat_hash_set",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "resolver_asset_test",
    srcs = ["resolver_asset_test.cc"],
    deps = [
        ":resolver",
        ":resolver_asset",
        ":resolver_test_fixture",
        "//intrinsic/assets/proto:installed_assets_cc_grpc_proto",
        "//intrinsic/assets/proto:installed_assets_cc_proto",
        "//intrinsic/executive/clips_cpp:extended_status_matchers",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/util/proto:parsed_type_url",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "resolver_skill_test",
    srcs = ["resolver_skill_test.cc"],
    deps = [
        ":resolver",
        ":resolver_test_fixture",
        "//intrinsic/skills/proto:skills_cc_proto",
        "//intrinsic/skills/testing:mock_skill_registry_client",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/util/proto:parsed_type_url",
        "//intrinsic/util/proto:type_url",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "resolver_mock",
    testonly = True,
    hdrs = ["resolver_mock.h"],
    deps = [
        ":file_descriptor_set_name",
        ":resolver",
        ":resolver_asset",
        "//intrinsic/util/proto:parsed_type_url",
        "//intrinsic/util/testing:gtest_wrapper",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/status:statusor",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "resolver_test_fixture",
    testonly = True,
    hdrs = ["resolver_test_fixture.h"],
    data = [
        ":status_specs_message",
    ],
    deps = [
        ":extended_status_codes",
        "//intrinsic/util/path_resolver",
        "//intrinsic/util/status:status_specs",
        "//intrinsic/util/testing:gtest_wrapper",
    ],
)

cc_library(
    name = "registry_pubsub",
    srcs = ["registry_pubsub.cc"],
    hdrs = ["registry_pubsub.h"],
    deps = [
        ":resolver",
        "//intrinsic/platform/pubsub",
        "//intrinsic/platform/pubsub:queryable",
        "//intrinsic/proto_tools/proto:proto_registry_pubsub_cc_proto",
        "//intrinsic/util/proto:parsed_type_url",
        "//intrinsic/util/proto:type_url",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/status:status_specs",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/container:flat_hash_set",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/synchronization",
        "@boringssl//:crypto",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "registry_pubsub_test",
    srcs = ["registry_pubsub_test.cc"],
    # PubSub code cannot run on TAP in asan and msan configurations
    tags = [
        "noasan",
        "nomsan",
    ],
    deps = [
        ":registry_pubsub",
        ":resolver",
        "//intrinsic/platform/pubsub",
        "//intrinsic/proto_tools/proto:proto_registry_pubsub_cc_proto",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/util/proto:parsed_type_url",
        "//intrinsic/util/testing:gtest_wrapper_main",
    ],
)

cc_binary(
    name = "proto-registry",
    srcs = ["proto_registry_main.cc"],
    data = [
        ":status_specs_message",
    ],
    deps = [
        ":extended_status_codes",
        ":proto_registry_service",
        ":registry_pubsub",
        ":resolver",
        ":resolver_asset",
        "//intrinsic/assets/proto:installed_assets_cc_grpc_proto",
        "//intrinsic/connect/cc/grpc:channel",
        "//intrinsic/icon/release/portable:init_intrinsic",
        "//intrinsic/skills/internal:skill_registry_client",
        "//intrinsic/stats:opencensus",
        "//intrinsic/util/path_resolver",
        "//intrinsic/util/status:return",
        "//intrinsic/util/status:status_builder",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/status:status_specs",
        "//intrinsic/util/thread",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
        "@com_github_grpc_grpc//:grpc_security_base",
    ],
)

# Proto Registry service
cc_oci_image(
    name = "proto-registry-image",
    binary = ":proto-registry",
    data_path = "/",
)

py_library(
    name = "proto_registry_client",
    srcs = ["proto_registry_client.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/proto_tools/proto:proto_registry_py_pb2",
        "//intrinsic/proto_tools/proto:proto_registry_py_pb2_grpc",
        "//intrinsic/util/proto:descriptors_py",
        requirement("grpcio"),
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_test(
    name = "proto_registry_client_test",
    srcs = ["proto_registry_client_test.py"],
    data = [
        "//intrinsic/proto_tools/registry/test_data:descriptor_test_message_proto_descriptor_set",
    ],
    deps = [
        ":proto_registry_client",
        "//intrinsic/proto_tools/proto:proto_registry_py_pb2",
        "//intrinsic/proto_tools/proto:proto_registry_py_pb2_grpc",
        "//intrinsic/util/path_resolver:path_resolver_py",
        "//intrinsic/util/proto:descriptors_py",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_protobuf//:protobuf_python",
    ],
)
