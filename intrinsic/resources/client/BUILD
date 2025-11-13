# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_python//python:defs.bzl", "py_library", "py_test")

package(default_visibility = ["//visibility:public"])

py_library(
    name = "resource_registry_client",
    srcs = ["resource_registry_client.py"],
    deps = [
        "//intrinsic/resources/proto:resource_handle_py_pb2",
        "//intrinsic/resources/proto:resource_registry_py_pb2",
        "//intrinsic/resources/proto:resource_registry_py_pb2_grpc",
        "//intrinsic/util/grpc:error_handling",
        requirement("grpcio"),
    ],
)

py_test(
    name = "resource_registry_client_test",
    srcs = ["resource_registry_client_test.py"],
    deps = [
        ":resource_registry_client",
        "//intrinsic/resources/proto:resource_handle_py_pb2",
        "//intrinsic/resources/proto:resource_registry_py_pb2",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

cc_library(
    name = "resource_registry_client_interface",
    hdrs = ["resource_registry_client_interface.h"],
    deps = [
        "//intrinsic/resources/proto:resource_registry_cc_proto",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:string_view",
    ],
)

cc_library(
    name = "resource_registry_client_cc",
    srcs = ["resource_registry_client.cc"],
    hdrs = ["resource_registry_client.h"],
    deps = [
        ":resource_registry_client_interface",
        "//intrinsic/connect/cc/grpc:channel",
        "//intrinsic/resources/proto:resource_registry_cc_grpc",
        "//intrinsic/resources/proto:resource_registry_cc_proto",
        "//intrinsic/util/grpc",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:string_view",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
    ],
)
