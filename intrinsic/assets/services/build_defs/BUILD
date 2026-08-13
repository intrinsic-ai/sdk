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

package(
    default_visibility = ["//visibility:public"],
)

bzl_library(
    name = "services_bzl",
    srcs = ["services.bzl"],
    visibility = ["//visibility:private"],
    deps = [
        "//intrinsic/assets/build_defs:asset_bzl",
        "//intrinsic/util/proto/build_defs:descriptor_set",
        "@bazel_lib//lib:paths.bzl",
    ],
)

go_library(
    name = "servicegen",
    srcs = ["servicegen.go"],
    importpath = "intrinsic/assets/services/build_defs/servicegen",
    deps = [
        "//intrinsic/assets/services:servicebundle",
        "//intrinsic/assets/services:servicefix",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/util/proto:protoio",
        "//intrinsic/util/proto:registryutil",
        "//intrinsic/util/proto:sourcecodeinfoview",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_protobuf//types/descriptorpb:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)

go_binary(
    name = "servicegen_main",
    srcs = ["servicegen_main.go"],
    importpath = "intrinsic/assets/services/build_defs/servicegen_main",
    deps = [
        ":servicegen",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/production:intrinsic",
        "//intrinsic/util:flag",
        "//intrinsic/util/proto:protoio",
        "@com_github_golang_glog//:go_default_library",
    ],
)
