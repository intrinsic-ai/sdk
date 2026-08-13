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

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@intrinsic_js_rules//:defs.bzl", "js_library")

package(
    default_visibility = [
        "//intrinsic/robot_definitions:__subpackages__",
        "//visibility:public",
    ],
)

proto_library(
    name = "digital_input_output_proto",
    srcs = ["digital_input_output.proto"],
)

cc_proto_library(
    name = "digital_input_output_cc_proto",
    deps = [":digital_input_output_proto"],
)

js_library(
    name = "digital_input_output_jspb_proto",
    deps = [":digital_input_output_proto"],
)
