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

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "hardware_module_config_proto",
    srcs = ["hardware_module_config.proto"],
    deps = [
        "//intrinsic/assets/proto:field_metadata_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_proto",
        "//intrinsic/icon/hardware_modules/sim_bus:sim_bus_hardware_module_proto",
        "@com_google_protobuf//:any_proto",
    ],
)

proto_library(
    name = "hardware_module_inspection_proto",
    srcs = ["hardware_module_inspection.proto"],
    deps = [
        "//intrinsic/assets/services/proto/v1:service_inspection_proto",
        "//intrinsic/icon/proto:safety_status_proto",
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)
