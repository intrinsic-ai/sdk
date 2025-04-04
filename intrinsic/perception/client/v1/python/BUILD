# Copyright 2023 Intrinsic Innovation LLC

# Python wrapper for public API.

load("@rules_python//python:defs.bzl", "py_library", "py_test")

# Unless we have more mature tests, the library remains unreleased.
package(default_visibility = [
    "//visibility:public",
])

py_library(
    name = "image_utils",
    srcs = ["image_utils.py"],
    deps = [
        "//intrinsic/perception/proto/v1:dimensions_py_pb2",
        "//intrinsic/perception/proto/v1:image_buffer_py_pb2",
        "@ai_intrinsic_sdks_pip_deps//numpy",
        "@ai_intrinsic_sdks_pip_deps//pillow",
    ],
)

py_test(
    name = "image_utils_test",
    srcs = ["image_utils_test.py"],
    data = ["//intrinsic/perception/testing/test_data:dataset_images"],
    deps = [
        ":image_utils",
        "//intrinsic/perception/proto/v1:dimensions_py_pb2",
        "//intrinsic/perception/proto/v1:image_buffer_py_pb2",
        "//intrinsic/util/path_resolver:path_resolver_py",
        "@ai_intrinsic_sdks_pip_deps//numpy",
        "@ai_intrinsic_sdks_pip_deps//pillow",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
    ],
)
