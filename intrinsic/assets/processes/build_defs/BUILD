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

package(default_visibility = ["//visibility:private"])

bzl_library(
    name = "processes_bzl",
    srcs = ["processes.bzl"],
    deps = [
        "//intrinsic/assets/build_defs:asset_bzl",
        "//intrinsic/util/proto/build_defs:descriptor_set",
        "@bazel_lib//lib:paths.bzl",
    ],
)

go_library(
    name = "processgen",
    srcs = ["processgen.go"],
    importpath = "intrinsic/assets/processes/build_defs/processgen",
    deps = [
        "//intrinsic/assets/processes:processbundle",
        "//intrinsic/assets/processes/proto:process_manifest_go_proto",
        "//intrinsic/executive/proto:behavior_tree_go_proto",
        "//intrinsic/util/proto:protoio",
        "//intrinsic/util/proto:registryutil",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
        "@org_golang_google_protobuf//types/descriptorpb:go_default_library",
    ],
)

go_binary(
    name = "processgen_main",
    srcs = ["processgen_main.go"],
    importpath = "intrinsic/assets/processes/build_defs/processgen_main",
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":processgen",
        "//intrinsic/production:intrinsic",
        "//intrinsic/util:flag",
        "@com_github_golang_glog//:go_default_library",
    ],
)
