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

package(default_visibility = ["//intrinsic:public_api_users"])

py_library(
    name = "behavior_tree_visitor",
    srcs = ["behavior_tree_visitor.py"],
    deps = ["//intrinsic/executive/proto:behavior_tree_py_pb2"],
)

py_test(
    name = "behavior_tree_visitor_test",
    srcs = ["behavior_tree_visitor_test.py"],
    deps = [
        ":behavior_tree_visitor",
        "//intrinsic/executive/proto:behavior_tree_py_pb2",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
        "@com_google_protobuf//:protobuf_python",
    ],
)
