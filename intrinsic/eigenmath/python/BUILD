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

load("@pybind11_bazel//:build_defs.bzl", "pybind_extension", "pybind_library")
load("@rules_python//python:defs.bzl", "py_test")

package(default_visibility = ["//visibility:public"])

pybind_library(
    name = "pose3_caster",
    hdrs = ["pose3_caster.h"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/math:pose3",
    ],
)

pybind_library(
    name = "quaternion_caster",
    hdrs = ["quaternion_caster.h"],
    deps = ["@com_gitlab_libeigen_eigen//:eigen"],
)

# For testing purpose only, do not use otherwise.
pybind_extension(
    name = "pose3_test_util",
    srcs = ["pose3_test_util.cc"],
    deps = [
        ":pose3_caster",
        "//intrinsic/eigenmath",
        "//intrinsic/eigenmath:so3",
        "//intrinsic/math:pose3",
    ],
)

py_test(
    name = "pose3_test",
    srcs = ["pose3_test.py"],
    deps = [
        ":pose3_test_util",
        "//intrinsic/math:pose3",
        "//intrinsic/math/python:data_types",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)
