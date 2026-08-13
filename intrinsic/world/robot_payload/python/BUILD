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

load("@rules_python//python:defs.bzl", "py_library", "py_test")

package(default_visibility = ["//visibility:public"])

py_library(
    name = "robot_payload",
    srcs = ["robot_payload.py"],
    deps = [
        "//intrinsic/math/python:data_types",
        "//intrinsic/math/python:proto_conversion",
        "//intrinsic/world/proto:robot_payload_py_pb2",
        "@ai_intrinsic_sdks_pip_deps//numpy",
    ],
)

py_test(
    name = "robot_payload_test",
    srcs = ["robot_payload_test.py"],
    deps = [
        ":robot_payload",
        "//intrinsic/math/python:data_types",
        "//intrinsic/math/python:proto_conversion",
        "//intrinsic/solutions/testing:compare",
        "//intrinsic/world/proto:robot_payload_py_pb2",
        "@ai_intrinsic_sdks_pip_deps//numpy",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
    ],
)
