# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@rules_python//python:defs.bzl", "py_library", "py_test")

package(default_visibility = ["//intrinsic/apps/common/resources:__subpackages__"])

py_library(
    name = "proto_registry_client",
    srcs = ["proto_registry_client.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/proto_tools/proto:proto_registry_py_pb2",
        "//intrinsic/proto_tools/proto:proto_registry_py_pb2_grpc",
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
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_protobuf//:protobuf_python",
    ],
)
