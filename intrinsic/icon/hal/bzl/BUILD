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
load("//bazel:go_macros.bzl", "go_binary", "go_test")

package(
    default_visibility = ["//visibility:private"],
)

bzl_library(
    name = "hardware_module_binary_bzl",
    srcs = ["hardware_module_binary.bzl"],
    deps = ["@rules_cc//cc:core_rules"],
)

bzl_library(
    name = "hardware_module_image_bzl",
    srcs = ["hardware_module_image.bzl"],
    deps = [
        ":hardware_module_binary_bzl",
        "//bazel:container_bzl",
    ],
)

bzl_library(
    name = "resources_bzl",
    srcs = ["resources.bzl"],
)

go_binary(
    name = "hal_manifest",
    srcs = ["hal_manifest.go"],
    embedsrcs = ["hal_service_manifest.textproto.tmpl"],
    importpath = "intrinsic/icon/hal/bzl/hal_manifest",
    visibility = ["//intrinsic:intrinsic_internal"],
    deps = [
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/production:intrinsic",
        "//intrinsic/util:flag",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_protobuf//encoding/prototext:go_default_library",
    ],
)

go_test(
    name = "hal_manifest_test",
    srcs = [
        "hal_manifest.go",
        "hal_manifest_test.go",
    ],
    embedsrcs = ["hal_service_manifest.textproto.tmpl"],
    importpath = "intrinsic/icon/hal/bzl/hal_manifest_test",
    deps = [
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/assets/services/proto:service_volume_go_proto",
        "//intrinsic/production:intrinsic",
        "//intrinsic/util:flag",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_protobuf//encoding/prototext:go_default_library",
    ],
)
