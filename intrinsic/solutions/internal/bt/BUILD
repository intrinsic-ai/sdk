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

load("@rules_python//python:defs.bzl", "py_library")

package(default_visibility = ["//intrinsic/solutions:__subpackages__"])

py_library(
    name = "bt_util",
    srcs = ["bt_util.py"],
    deps = [
        "//intrinsic/util/proto:descriptors_py",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_library(
    name = "world_query",
    srcs = ["world_query.py"],
    deps = [
        "//intrinsic/executive/proto:any_with_assignments_py_pb2",
        "//intrinsic/executive/proto:world_query_py_pb2",
        "//intrinsic/solutions:blackboard_value",
        "//intrinsic/solutions:errors",
        "//intrinsic/solutions:utils",
        "//intrinsic/world/proto:object_world_refs_py_pb2",
        "//intrinsic/world/python:object_world_resources",
    ],
)
