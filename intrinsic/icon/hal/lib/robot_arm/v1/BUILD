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
    "//intrinsic/icon/hal/lib/robot_arm:__subpackages__",
    "//intrinsic/icon/hardware_modules:__subpackages__",
    "//sdk:__subpackages__",
])

proto_library(
    name = "robot_arm_bus_component_config_proto",
    srcs = ["robot_arm_bus_component_config.proto"],
)

cc_proto_library(
    name = "robot_arm_device_config_cc_proto",
    deps = [":robot_arm_bus_component_config_proto"],
)

go_proto_library(
    name = "robot_arm_bus_component_config_go_proto",
    importpath = "intrinsic/icon/hal/lib/robot_arm/v1/robot_arm_bus_component_config_go_proto",
    protos = [":robot_arm_bus_component_config_proto"],
)
