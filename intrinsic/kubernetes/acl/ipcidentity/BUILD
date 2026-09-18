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
load("@rules_python//python:defs.bzl", "py_library", "py_test")

py_library(
    name = "ipcidentity_py",
    srcs = ["ipcidentity.py"],
    visibility = ["//intrinsic/util/grpc:__pkg__"],
    deps = [
        "//intrinsic/config:environments_py",
        "//intrinsic/kubernetes/accounts/service/api/tokens/v2:tokens_py_pb2",
        "//intrinsic/kubernetes/accounts/service/api/tokens/v2:tokens_py_pb2_grpc",
        "//intrinsic/kubernetes/acl/ipcidentity/internal:metadata",
        "//intrinsic/kubernetes/acl/py:identity",
        "//intrinsic/kubernetes/acl/py:jwt",
        requirement("grpcio"),
    ],
)

py_test(
    name = "ipcidentity_py_test",
    srcs = ["ipcidentity_py_test.py"],
    deps = [
        ":ipcidentity_py",
        "//intrinsic/kubernetes/accounts/service/api/tokens/v2:tokens_py_pb2",
        "//intrinsic/kubernetes/acl/ipcidentity/internal:metadata",
        "//intrinsic/kubernetes/acl/py:identity",
        requirement("grpcio"),
        "@com_google_absl_py//absl/testing:absltest",
    ],
)
