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

load("@aspect_rules_js//npm:defs.bzl", "npm_package", "stamped_package_json")
load("@bazel_skylib//rules:copy_file.bzl", "copy_file")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")

package(default_visibility = ["//visibility:public"])

copy_file(
    name = "license",
    src = "//intrinsic/frontend/sdk:LICENSE",
    out = "LICENSE",
)

stamped_package_json(
    name = "package",
    stamp_var = "STABLE_SEMVER",
)

npm_package(
    name = "types_package",
    srcs = [
        ":license",
        ":package",
        "//intrinsic/assets/proto/v1:asset_configuration_jspb_proto",
        "//intrinsic/executive/proto:blackboard_service_jspb_proto",
        "//intrinsic/executive/proto:executive_service_jspb_proto",
        "//intrinsic/executive/proto:proto_builder_jspb_proto",
        "//intrinsic/frontend/solution_service/proto:solution_service_jspb_proto",
        "//intrinsic/geometry/proto:geometry_service_jspb_proto",
        "//intrinsic/httpjson/proto/v1alpha:controlplane_jspb_proto",
        "//intrinsic/httpjson/proto/v1alpha:dataplane_jspb_proto",
        "//intrinsic/httpjson/proto/v1alpha:signalling_jspb_proto",
        "//intrinsic/httpjson/proto/v1alpha:webpubsub_jspb_proto",
        "//intrinsic/perception/proto/v1:train_service_jspb_proto",
        "//intrinsic/platform/pubsub/adapters:pubsub_jspb_proto",
        "//intrinsic/proto_tools/proto:proto_registry_jspb_proto",
        "//intrinsic/scene/proto/v1:scene_object_import_jspb_proto",
        "//intrinsic/world/proto:object_world_service_jspb_proto",
    ],
    include_external_repositories = ["*"],
    include_srcs_packages = ["**"],
    package = "@intrinsic-ai/types",
    publishable = True,
)

alias(
    name = "pkg",
    actual = ":types_package",
)

pkg_tar(
    name = "intrinsic-ai-types",
    srcs = [":types_package"],
    extension = "tgz",
    package_dir = "package",
    strip_prefix = "types_package",
)

alias(
    name = "pack",
    actual = ":intrinsic-ai-types",
)
