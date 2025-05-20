# Copyright 2023 Intrinsic Innovation LLC

load("@rules_python//python:defs.bzl", "py_library", "py_test")

package(default_visibility = [
    "//intrinsic/kubernetes/acl/ipcidentity:__subpackages__",
    "//intrinsic/util/grpc:__subpackages__",
])

py_library(
    name = "metadata",
    srcs = ["metadata.py"],
)

py_test(
    name = "metadata_test",
    srcs = ["metadata_test.py"],
    deps = [
        ":metadata",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)
