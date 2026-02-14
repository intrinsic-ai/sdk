# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_python//python:defs.bzl", "py_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "channel_factory",
    srcs = ["channel_factory.cc"],
    hdrs = ["channel_factory.h"],
    deps = [
        "//intrinsic/connect/cc/grpc:channel",
        "//intrinsic/util/grpc:channel",
        "//intrinsic/util/grpc:channel_interface",
        "//intrinsic/util/grpc:connection_params",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/time",
    ],
)

cc_library(
    name = "equipment_utils",
    srcs = ["equipment_utils.cc"],
    hdrs = ["equipment_utils.h"],
    deps = [
        ":channel_factory",
        ":icon_equipment_cc_proto",
        "//intrinsic/connect/cc/grpc:channel",
        "//intrinsic/skills/cc:equipment_pack",
        "//intrinsic/skills/cc:skill_utils",
        "//intrinsic/skills/proto:equipment_cc_proto",
        "//intrinsic/util/grpc:channel_interface",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
    ],
)

proto_library(
    name = "icon_equipment_proto",
    srcs = ["icon_equipment.proto"],
    deps = [
        ":force_control_settings_proto",
    ],
)

cc_proto_library(
    name = "icon_equipment_cc_proto",
    deps = [":icon_equipment_proto"],
)

py_proto_library(
    name = "icon_equipment_py_pb2",
    deps = [":icon_equipment_proto"],
)

proto_library(
    name = "force_control_settings_proto",
    srcs = ["force_control_settings.proto"],
    deps = ["//intrinsic/icon/proto:cart_space_proto"],
)

go_proto_library(
    name = "force_control_settings_go_proto",
    importpath = "intrinsic/icon/equipment/force_control_settings_go_proto",
    protos = [":force_control_settings_proto"],
    deps = [
        "//intrinsic/icon/proto:cart_space_go_proto",
    ],
)

cc_proto_library(
    name = "force_control_settings_cc_proto",
    deps = [":force_control_settings_proto"],
)

py_proto_library(
    name = "force_control_settings_py_pb2",
    deps = [":force_control_settings_proto"],
)

py_library(
    name = "equipment_utils_py",
    srcs = ["equipment_utils.py"],
    deps = [
        ":icon_equipment_py_pb2",
        "//intrinsic/icon/python:icon",
        "//intrinsic/resources/proto:resource_handle_py_pb2",
        "//intrinsic/skills/proto:equipment_py_pb2",
        "//intrinsic/util/grpc:connection",
        "@com_google_absl_py//absl/logging",
    ],
)
