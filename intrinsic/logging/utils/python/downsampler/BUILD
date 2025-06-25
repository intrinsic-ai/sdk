# Copyright 2023 Intrinsic Innovation LLC

load("@pybind11_bazel//:build_defs.bzl", "pybind_extension")
load("@rules_python//python:defs.bzl", "py_test")

pybind_extension(
    name = "downsampler",
    srcs = [
        "downsampler.cc",
    ],
    deps = [
        "//intrinsic/logging/proto:downsampler_cc_proto",
        "//intrinsic/logging/proto:log_item_cc_proto",
        "//intrinsic/logging/proto:logger_service_cc_proto",
        "//intrinsic/logging/utils/downsampler",
        "//intrinsic/logging/utils/downsampler:proto_conversion",
        "@pybind11_abseil//pybind11_abseil:absl_casters",
        "@pybind11_abseil//pybind11_abseil:status_casters",
        "@pybind11_protobuf//pybind11_protobuf:native_proto_caster",
    ],
)

py_test(
    name = "downsampler_test",
    srcs = ["downsampler_test.py"],
    deps = [
        ":downsampler",
        "//intrinsic/logging/proto:log_item_py_pb2",
        "@com_google_absl_py//absl/testing:absltest",
        "@pybind11_abseil//pybind11_abseil:import_status_module",
    ],
)

py_test(
    name = "proto_conversion_test",
    srcs = ["proto_conversion_test.py"],
    deps = [
        ":downsampler",
        "//intrinsic/logging/proto:downsampler_py_pb2",
        "@com_google_absl_py//absl/testing:absltest",
        "@pybind11_abseil//pybind11_abseil:import_status_module",
    ],
)
