# Copyright 2023 Intrinsic Innovation LLC

# Python ACl helpers

load("@rules_python//python:defs.bzl", "py_library", "py_test")

package(default_visibility = [
    "//visibility:public",
])

py_library(
    name = "jwt",
    srcs = ["jwt.py"],
)

py_test(
    name = "jwt_test",
    srcs = ["jwt_test.py"],
    deps = [
        ":jwt",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)
