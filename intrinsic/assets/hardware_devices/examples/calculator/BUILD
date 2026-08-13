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

load("//intrinsic/assets/data/build_defs:data.bzl", "intrinsic_data")
load("//intrinsic/assets/hardware_devices/build_defs:hardware_device.bzl", "intrinsic_hardware_device")

intrinsic_data(
    name = "calculator_config_data",
    manifest = "calculator_config_data_manifest.textproto",
    deps = [
        "//intrinsic/assets/services/examples/calcserver:calc_server_proto",
    ],
)

intrinsic_hardware_device(
    name = "calculator",
    assets = [
        ":calculator_config_data",
        "//intrinsic/assets/scene_objects/examples:calculator_scene_object_reference",
        "//intrinsic/assets/services/examples/calcserver:calculator_service",
        "//intrinsic/assets/services/examples/calcserver:units_conversion_data",
    ],
    manifest = "calculator_manifest.textproto",
)
