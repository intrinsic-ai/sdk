# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@rules_cc//cc:cc_binary.bzl", "cc_binary")
load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("//bazel:cc_oci_image.bzl", "cc_oci_image")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "kvstore_proto",
    srcs = ["kvstore.proto"],
    deps = ["@com_google_protobuf//:any_proto"],
)

exports_files(
    ["kvservice.yaml"],
    visibility = ["//intrinsic:intrinsic_internal"],
)

go_proto_library(
    name = "kvstore_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/platform/pubsub/kvstore_grpc/kvstore_go_proto",
    protos = [":kvstore_proto"],
)

cc_proto_library(
    name = "kvstore_cc_proto",
    deps = [":kvstore_proto"],
)

cc_grpc_library(
    name = "kvstore_cc_grpc_proto",
    srcs = [":kvstore_proto"],
    grpc_only = True,
    deps = [":kvstore_cc_proto"],
)

py_proto_library(
    name = "kvstore_py_pb2",
    deps = [":kvstore_proto"],
)

py_grpc_library(
    name = "kvstore_py_pb2_grpc",
    srcs = [":kvstore_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":kvstore_py_pb2"],
)

cc_library(
    name = "server_impl",
    srcs = ["server_impl.cc"],
    hdrs = ["server_impl.h"],
    deps = [
        ":kvstore_cc_grpc_proto",
        ":kvstore_cc_proto",
        "//intrinsic/platform/pubsub",
        "//intrinsic/platform/pubsub:kvstore",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/status:status_macros_grpc",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_binary(
    name = "server_main",
    srcs = ["server_main.cc"],
    deps = [
        ":server_impl",
        "//intrinsic/icon/release/portable:init_intrinsic",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/strings:str_format",
        "@com_github_grpc_grpc//:grpc++",
        "@com_github_grpc_grpc//:grpc_security_base",
    ],
)

cc_oci_image(
    name = "kvservice-image",
    binary = ":server_main",
    data_path = "/",
    visibility = ["//intrinsic:intrinsic_internal"],
)
