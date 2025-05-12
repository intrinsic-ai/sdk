# Copyright 2023 Intrinsic Innovation LLC

# Python wrapper for public API.

load("@rules_python//python:defs.bzl", "py_library")

py_library(
    name = "image_utils",
    srcs = ["image_utils.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/perception/proto:dimensions_py_pb2",
        "//intrinsic/perception/proto:image_buffer_py_pb2",
        "@ai_intrinsic_sdks_pip_deps//numpy",
        "@ai_intrinsic_sdks_pip_deps//pillow",
    ],
)
