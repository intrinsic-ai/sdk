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
