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
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = [
    "//intrinsic/icon/control/parts/proto:__subpackages__",
    "//intrinsic/icon/testing:__subpackages__",
    "//visibility:public",
])

proto_library(
    name = "realtime_part_config_proto",
    srcs = ["realtime_part_config.proto"],
    visibility = [
        "//intrinsic:intrinsic_internal",
        "//intrinsic/icon/control/parts/proto:__subpackages__",
    ],
    deps = ["@com_google_protobuf//:any_proto"],
)

cc_proto_library(
    name = "realtime_part_config_cc_proto",
    visibility = [
        "//intrinsic:intrinsic_internal",
        "//intrinsic/icon/control/parts/proto:__subpackages__",
    ],
    deps = [":realtime_part_config_proto"],
)

go_proto_library(
    name = "realtime_part_config_go_proto",
    importpath = "intrinsic/icon/control/parts/proto/v1/realtime_part_config_go_proto",
    protos = [":realtime_part_config_proto"],
    visibility = [
        "//intrinsic:intrinsic_internal",
        "//intrinsic/icon/control/parts/proto:__subpackages__",
    ],
)

proto_library(
    name = "mode_of_safe_operation_limits_config_proto",
    srcs = ["mode_of_safe_operation_limits_config.proto"],
    visibility = [
        "//intrinsic:intrinsic_internal",
        "//intrinsic/icon/control/parts/proto:__subpackages__",
    ],
)

cc_proto_library(
    name = "mode_of_safe_operation_limits_config_cc_proto",
    deps = [":mode_of_safe_operation_limits_config_proto"],
)

go_proto_library(
    name = "mode_of_safe_operation_limits_config_go_proto",
    importpath = "intrinsic/icon/control/parts/proto/v1/mode_of_safe_operation_limits_config_go_proto",
    protos = [":mode_of_safe_operation_limits_config_proto"],
    visibility = [
        "//intrinsic:intrinsic_internal",
        "//intrinsic/icon/control/parts/proto:__subpackages__",
    ],
)
