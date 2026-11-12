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

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@rules_python//python:defs.bzl", "py_binary", "py_library", "py_test")

package(default_visibility = ["//visibility:public"])

py_library(
    name = "solution_diff_lib",
    srcs = ["solution_diff.py"],
    deps = [
        "//intrinsic/solution_versions/proto/v1:snapshot_source_py_pb2",
        "//intrinsic/solution_versions/proto/v1:solution_version_service_pb2_grpc",
        "//intrinsic/solution_versions/proto/v1:solution_version_service_py_pb2",
        "//intrinsic/util/grpc:auth_py",
        "//intrinsic/util/grpc:dialerutil",
        "@ai_intrinsic_sdks_pip_deps//protobuf",
        requirement("grpcio"),
        "@com_google_absl_py//absl:app",
        "@com_google_absl_py//absl/flags",
    ],
)

py_binary(
    name = "solution_diff",
    srcs = ["solution_diff.py"],
    deps = [":solution_diff_lib"],
)

py_test(
    name = "solution_diff_test",
    srcs = ["solution_diff_test.py"],
    deps = [
        ":solution_diff_lib",
        "//intrinsic/assets/proto/v1:solution_py_pb2",
        "//intrinsic/solution_versions/proto/v1:branch_py_pb2",
        "//intrinsic/solution_versions/proto/v1:snapshot_py_pb2",
        "//intrinsic/solution_versions/proto/v1:solution_version_service_py_pb2",
        "//intrinsic/util/grpc:auth_py",
        "//intrinsic/util/grpc:dialerutil",
        "@ai_intrinsic_sdks_pip_deps//protobuf",
        requirement("grpcio"),
        "@com_google_absl_py//absl:app",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:flagsaver",
    ],
)
