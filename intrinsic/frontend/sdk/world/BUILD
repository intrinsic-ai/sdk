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

load("@aspect_rules_js//npm:defs.bzl", "npm_package")
load("@bazel_skylib//rules:copy_file.bzl", "copy_file")
load("@npm//:defs.bzl", "npm_link_all_packages")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")
load("//intrinsic/frontend/sdk:stamped_package_json.bzl", "stamped_package_json")

copy_file(
    name = "license",
    src = "//intrinsic/frontend/sdk:LICENSE",
    out = "LICENSE",
)

stamped_package_json(
    name = "package",
    stamp_dependencies = [
        "@intrinsic-ai/core",
    ],
    stamp_var = "STABLE_SEMVER",
)

npm_link_all_packages(name = "node_modules")

# Not publishable yet. See the "private" field in package.json.
npm_package(
    name = "world_package",
    srcs = [
        ":license",
        ":package",
        "//intrinsic/frontend/sdk/world/utils:lib",
    ],
    package = "@intrinsic-ai/world",
    publishable = False,
)

alias(
    name = "pkg",
    actual = ":world_package",
    visibility = ["//intrinsic:internal_api_users"],
)

pkg_tar(
    name = "intrinsic-ai-world",
    srcs = [":world_package"],
    extension = "tgz",
    package_dir = "package",
    strip_prefix = "world_package",
)

alias(
    name = "pack",
    actual = ":intrinsic-ai-world",
    visibility = ["//intrinsic:internal_api_users"],
)
