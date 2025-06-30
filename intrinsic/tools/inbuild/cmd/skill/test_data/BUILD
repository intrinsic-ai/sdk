# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@rules_python//python:defs.bzl", "py_binary")
load("//bazel:go_macros.bzl", "go_proto_library")
load("//bazel:python_oci_image.bzl", "python_oci_image")

package(default_visibility = ["//intrinsic/tools/inbuild:__subpackages__"])

proto_library(
    name = "example_skill_proto",
    srcs = ["example_skill.proto"],
)

py_proto_library(
    name = "example_skill_py_pb2",
    deps = [":example_skill_proto"],
)

go_proto_library(
    name = "example_skill_go_proto",
    deps = [":example_skill_proto"],
)

py_binary(
    name = "example_skill_main_py",
    srcs = ["example_skill_main.py"],
    main = "example_skill_main.py",
    deps = [
        "@com_google_absl_py//absl:app",
        "@com_google_absl_py//absl/logging",
    ],
)

python_oci_image(
    name = "example_skill_py_image",
    base = "@distroless_base",
    binary = ":example_skill_main_py",
    data_path = "/",
)

filegroup(
    name = "example_skill_files",
    srcs = [
        "example_skill.proto",
        "example_skill_cc.manifest.pbtxt",
        "example_skill_main.py",
        "example_skill_py.manifest.pbtxt",
    ],
)
