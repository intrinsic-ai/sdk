# Copyright 2023 Intrinsic Innovation LLC

# Python ACl helpers

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
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

py_library(
    name = "identity",
    srcs = ["identity.py"],
    deps = [
        ":jwt",
        requirement("grpcio"),
        "@com_google_absl_py//absl/logging",
    ],
)

py_test(
    name = "identity_test",
    srcs = ["identity_test.py"],
    deps = [
        ":identity",
        requirement("grpcio"),
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
    ],
)
