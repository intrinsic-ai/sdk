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

package(default_visibility = ["//intrinsic/apps/common/resources:__subpackages__"])

py_library(
    name = "proto_registry_client",
    srcs = ["proto_registry_client.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/proto_tools/proto:proto_registry_py_pb2",
        "//intrinsic/proto_tools/proto:proto_registry_py_pb2_grpc",
        requirement("grpcio"),
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_test(
    name = "proto_registry_client_test",
    srcs = ["proto_registry_client_test.py"],
    data = [
        "//intrinsic/proto_tools/registry/test_data:descriptor_test_message_proto_descriptor_set",
    ],
    deps = [
        "//intrinsic/proto_tools/proto:proto_registry_py_pb2",
        "//intrinsic/proto_tools/proto:proto_registry_py_pb2_grpc",
        "//intrinsic/proto_tools/registry:proto_registry_client",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_protobuf//:protobuf_python",
    ],
)
