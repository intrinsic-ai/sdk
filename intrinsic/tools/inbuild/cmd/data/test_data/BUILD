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
load("//bazel:go_macros.bzl", "go_proto_library")
load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "proto_source_code_info_transitive_descriptor_set")

package(default_visibility = ["//intrinsic/tools/inbuild:__subpackages__"])

proto_library(
    name = "custom_data_proto",
    srcs = ["custom_data.proto"],
    deps = ["//intrinsic/assets/data/proto/v1:referenced_data_proto"],
)

go_proto_library(
    name = "custom_data_go_proto",
    importpath = "intrinsic/tools/inbuild/cmd/data/test_data/custom_data_go_proto",
    protos = [":custom_data_proto"],
    deps = ["//intrinsic/assets/data/proto/v1:referenced_data_go_proto"],
)

proto_source_code_info_transitive_descriptor_set(
    name = "custom_data_fds",
    testonly = True,
    deps = [":custom_data_proto"],
)
