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

load("//bazel:go_macros.bzl", "go_library", "go_test")

go_library(
    name = "data",
    srcs = [
        "data.go",
    ],
    importpath = "intrinsic/assets/conversion/data",
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/data/proto/v1:data_manifest_go_proto",
    ],
)

go_test(
    name = "data_test",
    srcs = [
        "data_test.go",
    ],
    embed = [":data"],
    importpath = "intrinsic/assets/conversion/data_test",
    deps = [
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/data/proto/v1:data_manifest_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:documentation_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
        "@org_golang_google_protobuf//types/known/emptypb",
    ],
)

go_library(
    name = "idparser",
    srcs = [
        "id_parser.go",
    ],
    importpath = "intrinsic/assets/conversion/idparser",
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto/v1:asset_go_proto",
        "//intrinsic/assets/proto/v1:processed_asset_go_proto",
    ],
)

go_library(
    name = "sceneobject",
    srcs = [
        "sceneobject.go",
    ],
    importpath = "intrinsic/assets/conversion/sceneobject",
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
    ],
)

go_test(
    name = "sceneobject_test",
    srcs = [
        "sceneobject_test.go",
    ],
    embed = [":sceneobject"],
    importpath = "intrinsic/assets/conversion/sceneobject_test",
    deps = [
        "//intrinsic/assets/proto:documentation_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)
