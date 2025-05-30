# Copyright 2023 Intrinsic Innovation LLC

load("@rules_python//python:defs.bzl", "py_library")
load("@rules_python//python:packaging.bzl", "py_package", "py_wheel")

package(default_visibility = ["//visibility:public"])

py_library(
    name = "skill_interface",
    srcs = ["skill_interface.py"],
    deps = [
        ":execute_context",
        ":execute_request",
        ":get_footprint_context",
        ":get_footprint_request",
        ":preview_context",
        ":preview_request",
        "//intrinsic/skills/proto:footprint_py_pb2",
        "//intrinsic/skills/proto:skill_service_py_pb2",
        "//intrinsic/util/status:status_exception",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_library(
    name = "skill_interface_utils",
    srcs = ["skill_interface_utils.py"],
    deps = [
        ":skill_interface",
        "//intrinsic/resources/proto:resource_handle_py_pb2",
        "//intrinsic/skills/internal:execute_context_impl_py",
        "//intrinsic/skills/internal:preview_context_impl_py",
    ],
)

py_library(
    name = "get_footprint_context",
    srcs = ["get_footprint_context.py"],
    deps = [
        "//intrinsic/geometry/proto:geometry_service_py_pb2_grpc",
        "//intrinsic/motion_planning:motion_planner_client_py",
        "//intrinsic/world/python:object_world_client",
        "//intrinsic/world/python:object_world_ids",
        "//intrinsic/world/python:object_world_resources",
    ],
)

py_library(
    name = "get_footprint_request",
    srcs = ["get_footprint_request.py"],
    deps = ["@com_google_protobuf//:protobuf_python"],
)

py_library(
    name = "basic_compute_context",
    srcs = ["basic_compute_context.py"],
    deps = [
        "//intrinsic/world/python:object_world_client",
    ],
)

py_library(
    name = "execute_context",
    srcs = ["execute_context.py"],
    deps = [
        ":skill_canceller",
        ":skill_logging_context",
        "//intrinsic/geometry/proto:geometry_service_py_pb2_grpc",
        "//intrinsic/motion_planning:motion_planner_client_py",
        "//intrinsic/resources/proto:resource_handle_py_pb2",
        "//intrinsic/world/python:object_world_client",
    ],
)

py_library(
    name = "execute_request",
    srcs = ["execute_request.py"],
    deps = ["@com_google_protobuf//:protobuf_python"],
)

py_library(
    name = "preview_context",
    srcs = ["preview_context.py"],
    deps = [
        ":skill_canceller",
        ":skill_logging_context",
        "//intrinsic/geometry/proto:geometry_service_py_pb2_grpc",
        "//intrinsic/motion_planning:motion_planner_client_py",
        "//intrinsic/world/proto:object_world_updates_py_pb2",
        "//intrinsic/world/python:object_world_client",
        "//intrinsic/world/python:object_world_ids",
        "//intrinsic/world/python:object_world_resources",
    ],
)

py_library(
    name = "preview_request",
    srcs = ["preview_request.py"],
    deps = ["@com_google_protobuf//:protobuf_python"],
)

py_library(
    name = "skill_canceller",
    srcs = ["skill_canceller.py"],
)

py_library(
    name = "proto_utils",
    srcs = ["proto_utils.py"],
    deps = ["@com_google_protobuf//:protobuf_python"],
)

py_library(
    name = "skill_logging_context",
    srcs = ["skill_logging_context.py"],
    deps = [
        "//intrinsic/logging/proto:context_py_pb2",
    ],
)

py_package(
    name = "skills_py_packages",
    packages = [
        "intrinsic",
    ],
    deps = [
        ":proto_utils",
        ":skill_interface",
        "//intrinsic/icon/proto:icon_protos",
        "//intrinsic/icon/python:icon",
        "//intrinsic/math/proto:math_protos",
        "//intrinsic/skills/proto:skills_all_protos",
        "//intrinsic/util:decorators",
        "//intrinsic/world/proto:world_protos",
    ],
)

py_wheel(
    name = "skills_wheel",
    distribution = "skills",
    requires = [
        "grpcio",
        "google-api-core",
        "numpy",
        "retrying",
    ],
    version = "0.0.1",
    deps = [":skills_py_packages"],
)
