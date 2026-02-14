# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@rules_python//python:defs.bzl", "py_binary")
load("//bazel:go_macros.bzl", "go_proto_library")
load("//bazel:python_oci_image.bzl", "python_oci_image")

package(default_visibility = ["//intrinsic/tools/inbuild:__subpackages__"])

proto_library(
    name = "example_service_proto",
    srcs = ["example_service.proto"],
)

py_proto_library(
    name = "example_service_py_pb2",
    deps = [":example_service_proto"],
)

go_proto_library(
    name = "example_service_go_proto",
    importpath = "intrinsic/tools/inbuild/cmd/service/test_data/example_service_go_proto",
    protos = [":example_service_proto"],
)

py_binary(
    name = "example_service_main_py",
    srcs = ["example_service_main.py"],
    main = "example_service_main.py",
    deps = [
        ":example_service_py_pb2",
        "//intrinsic/resources/proto:runtime_context_py_pb2",
        "@com_google_absl_py//absl:app",
        "@com_google_absl_py//absl/flags",
        "@com_google_absl_py//absl/logging",
    ],
)

python_oci_image(
    name = "example_service_py_image",
    base = "@distroless_base",
    binary = ":example_service_main_py",
    data_path = "/",
)

filegroup(
    name = "example_service_files",
    srcs = [
        "example_service.manifest.pbtxt",
        "example_service.proto",
        "example_service_default_config.textproto",
        "example_service_main.py",
    ],
)
