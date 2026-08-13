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

load("@bazel_skylib//:bzl_library.bzl", "bzl_library")
load("//bazel:go_macros.bzl", "go_binary", "go_library")

bzl_library(
    name = "hardware_device_bzl",
    srcs = ["hardware_device.bzl"],
    visibility = ["//visibility:private"],
    deps = [
        "//intrinsic/assets/build_defs:asset_bzl",
        "@bazel_lib//lib:paths.bzl",
    ],
)

go_library(
    name = "hardwaredevicegen",
    srcs = ["hardwaredevicegen.go"],
    importpath = "intrinsic/assets/hardware_devices/build_defs/hardwaredevicegen",
    visibility = ["//visibility:private"],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets/build_defs:asset_go_proto",
        "//intrinsic/assets/hardware_devices:hardwaredevicebundle",
        "//intrinsic/assets/hardware_devices:hardwaredevicefix",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_go_proto",
        "//intrinsic/assets/proto/v1:reference_go_proto",
        "//intrinsic/util/proto:protoio",
    ],
)

go_binary(
    name = "hardwaredevicegen_main",
    srcs = ["hardwaredevicegen_main.go"],
    importpath = "intrinsic/assets/hardware_devices/build_defs/hardwaredevicegen_main",
    visibility = ["//visibility:private"],
    deps = [
        ":hardwaredevicegen",
        "//intrinsic/production:intrinsic",
        "//intrinsic/util:flag",
        "@com_github_golang_glog//:go_default_library",
    ],
)
