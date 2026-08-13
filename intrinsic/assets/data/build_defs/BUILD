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
    name = "data_bzl",
    srcs = ["data.bzl"],
    visibility = ["//visibility:private"],
    deps = [
        "//intrinsic/assets/build_defs:asset_bzl",
        "//intrinsic/util/proto/build_defs:descriptor_set",
        "@bazel_lib//lib:paths.bzl",
        "@bazel_skylib//lib:paths",
        "@bazel_skylib//rules:common_settings",
    ],
)

go_library(
    name = "datagen",
    srcs = ["datagen.go"],
    importpath = "intrinsic/assets/data/build_defs/datagen",
    visibility = [
        "//intrinsic/assets/data:__subpackages__",
        "//intrinsic/tools/inbuild:__subpackages__",
    ],
    deps = [
        "//intrinsic/assets:referenceddata",
        "//intrinsic/assets/data:databundle",
        "//intrinsic/assets/data:datavalidate",
        "//intrinsic/assets/data:utils",
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/data/proto/v1:data_manifest_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/util/proto:protoio",
        "//intrinsic/util/proto:registryutil",
        "@org_golang_google_protobuf//reflect/protodesc:go_default_library",
        "@org_golang_google_protobuf//types/descriptorpb:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)

go_binary(
    name = "datagen_main",
    srcs = ["datagen_main.go"],
    importpath = "intrinsic/assets/data/build_defs/datagen_main",
    visibility = ["//visibility:private"],
    deps = [
        ":datagen",
        "//intrinsic/production:intrinsic",
        "//intrinsic/util:flag",
        "//intrinsic/util/proto:registryutil",
        "@com_github_golang_glog//:go_default_library",
    ],
)
