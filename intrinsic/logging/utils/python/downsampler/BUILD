# Copyright 2026 Intrinsic Innovation LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@pybind11_bazel//:build_defs.bzl", "pybind_extension")
load("@rules_python//python:defs.bzl", "py_test")

pybind_extension(
    name = "downsampler",
    srcs = [
        "downsampler.cc",
    ],
    deps = [
        "//intrinsic/logging/proto:downsampler_cc_proto",
        "//intrinsic/logging/proto:logger_service_cc_proto",
        "//intrinsic/logging/utils/downsampler",
        "//intrinsic/logging/utils/downsampler:proto_conversion",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/time",
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
