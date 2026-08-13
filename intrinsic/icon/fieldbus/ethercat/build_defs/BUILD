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

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

package(default_visibility = ["//visibility:public"])

go_binary(
    name = "eni_datagen",
    embed = [":build_defs_lib"],
    importpath = "intrinsic/icon/fieldbus/ethercat/build_defs/eni_datagen",
)

filegroup(
    name = "bzl",
    srcs = glob(["*.bzl"]),
)

go_library(
    name = "build_defs_lib",
    srcs = ["eni_datagen.go"],
    importpath = "intrinsic/icon/fieldbus/ethercat/build_defs/build_defs_lib",
    deps = [
        "//intrinsic/assets/data/proto/v1:data_manifest_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/icon/fieldbus/ethercat/device_service/v1:eni_go_proto",
        "@org_golang_google_protobuf//encoding/prototext",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)
