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
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@intrinsic_js_rules//:defs.bzl", "js_library")
load("//bazel:go_macros.bzl", "go_proto_library")

proto_library(
    name = "hardware_device_manifest_proto",
    srcs = ["hardware_device_manifest.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/data/proto/v1:data_asset_proto",
        "//intrinsic/assets/proto:asset_tag_proto",
        "//intrinsic/assets/proto:documentation_proto",
        "//intrinsic/assets/proto:id_proto",
        "//intrinsic/assets/proto:vendor_proto",
        "//intrinsic/assets/proto/v1:asset_graph_proto",
        "//intrinsic/assets/proto/v1:reference_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_proto",
        "//intrinsic/assets/services/proto:service_manifest_proto",
    ],
)

go_proto_library(
    name = "hardware_device_manifest_go_proto",
    importpath = "intrinsic/assets/hardware_devices/proto/v1/hardware_device_manifest_go_proto",
    protos = [":hardware_device_manifest_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/proto:asset_tag_go_proto",
        "//intrinsic/assets/proto:documentation_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/assets/proto/v1:asset_graph_go_proto",
        "//intrinsic/assets/proto/v1:reference_go_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
    ],
)

py_proto_library(
    name = "hardware_device_manifest_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":hardware_device_manifest_proto"],
)

cc_proto_library(
    name = "hardware_device_manifest_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":hardware_device_manifest_proto"],
)

js_library(
    name = "hardware_device_manifest_jspb_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":hardware_device_manifest_proto"],
)
