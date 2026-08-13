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

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "assetutils",
    srcs = ["assetutils.go"],
    importpath = "intrinsic/assets/catalog/assetutils",
    visibility = ["//intrinsic:intrinsic_internal"],
    deps = [
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_proto",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/kubernetes/workcell_spec/proto:image_go_proto",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//types/descriptorpb:go_default_library",
    ],
)

go_library(
    name = "releaseasset",
    srcs = ["releaseasset.go"],
    importpath = "intrinsic/assets/catalog/releaseasset",
    visibility = [
        "//intrinsic:intrinsic_internal",
    ],
    deps = [
        "//intrinsic/assets:bundle",
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:imagetransfer",
        "//intrinsic/assets:referenceddata",
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_proto",
        "//intrinsic/assets/catalog/proto/v1:release_metadata_go_proto",
        "//intrinsic/assets/proto/v1:asset_artifacts_go_proto",
        "//intrinsic/assets/scene_objects:gzfprocessor",
        "//intrinsic/assets/services:bundleimages",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
