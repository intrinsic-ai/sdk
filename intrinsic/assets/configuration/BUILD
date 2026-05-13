# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@rules_python//python:defs.bzl", "py_library", "py_test")

py_library(
    name = "asset_configuration_client",
    srcs = ["asset_configuration_client.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/proto/v1:asset_configuration_pb2_grpc",
        "//intrinsic/assets/proto/v1:asset_configuration_py_pb2",
        "//intrinsic/util/grpc:error_handling",
        requirement("grpcio"),
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_test(
    name = "asset_configuration_client_test",
    srcs = ["asset_configuration_client_test.py"],
    deps = [
        ":asset_configuration_client",
        "//intrinsic/assets/proto/v1:asset_configuration_py_pb2",
        requirement("grpcio"),
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_protobuf//:protobuf_python",
    ],
)
