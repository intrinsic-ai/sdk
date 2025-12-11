# Copyright 2023 Intrinsic Innovation LLC

# Python wrapper for the intrinsic pubsub c++ library.

load("@pybind11_bazel//:build_defs.bzl", "pybind_extension")
load("@rules_python//python:defs.bzl", "py_test")

package(default_visibility = ["//visibility:public"])

pybind_extension(
    name = "pubsub",
    srcs = ["pubsub.cc"],
    deps = [
        "//intrinsic/platform/pubsub",
        "//intrinsic/platform/pubsub:kvstore",
        "//intrinsic/platform/pubsub:publisher",
        "//intrinsic/platform/pubsub:subscription",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
        "@com_google_protobuf//:protobuf",
        "@pybind11_abseil//pybind11_abseil:absl_casters",
        "@pybind11_abseil//pybind11_abseil:no_throw_status",
        "@pybind11_abseil//pybind11_abseil:status_casters",
        "@pybind11_protobuf//pybind11_protobuf:native_proto_caster",
    ],
)

py_test(
    name = "pubsub_test",
    srcs = ["pubsub_test.py"],
    tags = [
        "noasan",
        "nomsan",
    ],
    deps = [
        ":pubsub",
        "//intrinsic/platform/common/proto:test_cc_proto",
        "//intrinsic/platform/common/proto:test_py_pb2",
        "//intrinsic/solutions/testing:compare",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_test(
    name = "callback_exception_test",
    srcs = ["callback_exception_test.py"],
    deps = [
        ":pubsub",
        "//intrinsic/platform/common/proto:test_cc_proto",
        "//intrinsic/platform/common/proto:test_py_pb2",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_protobuf//:protobuf_python",
    ],
)
