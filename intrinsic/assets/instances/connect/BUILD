# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")
load("@rules_python//python:defs.bzl", "py_library", "py_test")
load("//bazel:go_macros.bzl", "go_library", "go_test")

go_library(
    name = "connect",
    srcs = ["connect.go"],
    importpath = "intrinsic/assets/instances/connect",
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/proto/v1:grpc_connection_go_proto",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials/insecure:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
    ],
)

go_test(
    name = "connect_test",
    srcs = ["connect_test.go"],
    embed = [":connect"],
    importpath = "intrinsic/assets/instances/connect_test",
    deps = [
        "//intrinsic/assets/instances/connect/testing:test_service_go_proto",
        "//intrinsic/assets/proto/v1:grpc_connection_go_proto",
        "//intrinsic/testing:grpctest",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
    ],
)

cc_library(
    name = "connect_cc",
    srcs = ["connect.cc"],
    hdrs = ["connect.h"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/proto/v1:grpc_connection_cc_proto",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_github_grpc_grpc//:grpc++",
    ],
)

cc_test(
    name = "connect_cc_test",
    srcs = ["connect_test.cc"],
    deps = [
        ":connect_cc",
        "//intrinsic/assets/instances/connect/testing:test_service_cc_grpc_proto",
        "//intrinsic/assets/proto/v1:grpc_connection_cc_proto",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@com_github_grpc_grpc//:grpc++",
        "@com_github_grpc_grpc//:grpc_security_base",
    ],
)

py_library(
    name = "connect_py",
    srcs = ["connect.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/proto/v1:grpc_connection_py_pb2",
        "//intrinsic/util/grpc:interceptor",
        requirement("grpcio"),
        "@com_google_absl_py//absl/logging",
    ],
)

py_test(
    name = "connect_py_test",
    srcs = ["connect_test.py"],
    main = "connect_test.py",
    deps = [
        ":connect_py",
        "//intrinsic/assets/instances/connect/testing:test_service_py_grpc_proto",
        "//intrinsic/assets/instances/connect/testing:test_service_py_pb2",
        "//intrinsic/assets/proto/v1:grpc_connection_py_pb2",
        requirement("grpcio"),
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
    ],
)
