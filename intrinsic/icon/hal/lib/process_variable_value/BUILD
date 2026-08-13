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

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = [
    "//intrinsic/icon/hardware_modules:__subpackages__",
    "//sdk:__subpackages__",
])

cc_library(
    name = "process_variable_device",
    srcs = ["process_variable_device.cc"],
    hdrs = ["process_variable_device.h"],
    deps = [
        "//intrinsic/icon/hal:default_hardware_interfaces",
        "//intrinsic/icon/hal/lib/fieldbus:async_device_request",
        "//intrinsic/icon/hal/lib/fieldbus:bus_component",
        "//intrinsic/icon/hal/lib/fieldbus:bus_component_factory",
        "//intrinsic/icon/hal/lib/fieldbus:device_init_context",
        "//intrinsic/icon/hal/lib/fieldbus:process_variable",
        "//intrinsic/icon/hal/lib/fieldbus:variable_registry",
        "//intrinsic/icon/hal/lib/process_variable_value/v1:process_variable_config_cc_proto",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)
