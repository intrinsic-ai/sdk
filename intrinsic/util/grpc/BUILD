# Copyright 2023 Intrinsic Innovation LLC

"""Helper utilities for applications using grpc."""

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_python//python:defs.bzl", "py_library", "py_test")
load("//bazel:go_macros.bzl", "go_library", "go_test")

package(
    default_visibility = ["//visibility:public"],
)

cc_library(
    name = "channel_interface",
    srcs = ["channel_interface.cc"],
    hdrs = ["channel_interface.h"],
    deps = [
        ":grpc",
        "@com_github_grpc_grpc//:grpc++",
    ],
)

cc_library(
    name = "channel",
    srcs = ["channel.cc"],
    hdrs = ["channel.h"],
    deps = [
        ":channel_interface",
        ":connection_params",
        ":grpc",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
    ],
)

cc_library(
    name = "connection_params",
    srcs = ["connection_params.cc"],
    hdrs = ["connection_params.h"],
    deps = [
        "@abseil-cpp//absl/strings:str_format",
    ],
)

cc_library(
    name = "connection_cache",
    hdrs = ["connection_cache.h"],
    deps = [
        ":channel",
        ":connection_params",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/synchronization",
    ],
)

py_library(
    name = "connection",
    srcs = [
        "connection.py",
    ],
)

cc_library(
    name = "grpc",
    srcs = ["grpc.cc"],
    hdrs = ["grpc.h"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":limits",
        "//intrinsic/icon/release:grpc_time_support",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/thread",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
        "@com_github_grpc_grpc//:grpc_security_base",
        "@com_github_grpc_grpc//src/proto/grpc/health/v1:health_cc_grpc",
        "@com_github_grpc_grpc//src/proto/grpc/health/v1:health_cc_proto",
    ],
)

cc_library(
    name = "limits",
    hdrs = ["limits.h"],
    visibility = ["//intrinsic:public_api_users"],
)

py_library(
    name = "limits_py",
    srcs = ["limits.py"],
    visibility = ["//intrinsic:public_api_users"],
)

go_library(
    name = "grpclimits",
    srcs = ["limits.go"],
    visibility = ["//intrinsic:public_api_users"],
)

py_library(
    name = "interceptor",
    srcs = ["interceptor.py"],
    deps = [
        "//intrinsic/util:decorators",
        requirement("grpcio"),
    ],
)

py_library(
    name = "error_handling",
    srcs = ["error_handling.py"],
    deps = [
        "@ai_intrinsic_sdks_pip_deps//retrying",
        requirement("grpcio"),
        "@com_google_googleapis//google/rpc:status_py_proto",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_test(
    name = "error_handling_test",
    srcs = ["error_handling_test.py"],
    deps = [
        ":error_handling",
        "//intrinsic/skills/proto:error_py_pb2",
        "//intrinsic/util/status:extended_status_py_pb2",
        requirement("grpcio"),
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_googleapis//google/rpc:status_py_proto",
    ],
)

go_library(
    name = "statusutil",
    srcs = ["status_util.go"],
    deps = [
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "statusutil_test",
    srcs = ["status_util_test.go"],
    library = ":statusutil",
    deps = [
        "//intrinsic/testing:grpctest",
        "//intrinsic/util/grpc/testing:ping_go_grpc_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
        "@com_github_pkg_errors//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//credentials/local:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//types/known/emptypb",
    ],
)

py_library(
    name = "grpc_utils",
    srcs = ["grpc_utils.py"],
    deps = [
        ":interceptor",
        "//intrinsic/config:environments_py",
        "//intrinsic/kubernetes/acl/ipcidentity:ipcidentity_py",
        "//intrinsic/kubernetes/acl/ipcidentity/internal:metadata",
        requirement("grpcio"),
    ],
)

py_test(
    name = "grpc_utils_test",
    srcs = ["grpc_utils_test.py"],
    deps = [
        ":grpc_utils",
        "//intrinsic/kubernetes/acl/ipcidentity:ipcidentity_py",
        requirement("grpcio"),
        "@com_google_absl_py//absl/testing:absltest",
    ],
)
