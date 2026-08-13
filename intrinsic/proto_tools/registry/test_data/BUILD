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

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "proto_source_code_info_transitive_descriptor_set")

package(default_visibility = ["//intrinsic/proto_tools/registry:__subpackages__"])

proto_library(
    name = "descriptor_test_message_proto",
    testonly = True,
    srcs = ["descriptor_test_message.proto"],
    deps = ["@com_google_protobuf//:timestamp_proto"],
)

proto_source_code_info_transitive_descriptor_set(
    name = "descriptor_test_message_proto_descriptor_set",
    testonly = True,
    deps = [":descriptor_test_message_proto"],
)
