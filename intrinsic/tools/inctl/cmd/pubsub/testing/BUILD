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

package(default_visibility = ["//visibility:public"])

go_library(
    name = "testing",
    srcs = [
        "fake_asset_catalog_server.go",
        "fake_asset_deployment_server.go",
        "fake_asset_instances_server.go",
        "fake_installed_assets_server.go",
        "fake_operations_server.go",
    ],
    importpath = "intrinsic/tools/inctl/cmd/pubsub/testing",
    deps = [
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_proto",
        "//intrinsic/assets/proto:asset_deployment_go_proto",
        "//intrinsic/assets/proto:installed_assets_go_proto",
        "//intrinsic/assets/proto/v1:asset_instances_go_proto",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
