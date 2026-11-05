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
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")
load("//intrinsic/frontend/sdk:stamped_package_json.bzl", "stamped_package_json")

# intrinsic:!enable_npm_core:insert load("@npm//:defs.bzl", "npm_link_all_packages")
# intrinsic:!enable_npm_core:insert load("@aspect_rules_ts//ts:defs.bzl", "ts_config")

package(default_visibility = ["//visibility:public"])

copy_file(
    name = "license",
    src = "//intrinsic/frontend/sdk:LICENSE",
    out = "LICENSE",
)

stamped_package_json(
    name = "package",
    stamp_dependencies = [
        "@intrinsic-ai/types",
    ],
    stamp_var = "STABLE_SEMVER",
)

# intrinsic:!enable_npm_core:insert npm_link_all_packages(name = "node_modules")

# intrinsic:!enable_npm_core:insert ts_config(
# intrinsic:!enable_npm_core:insert     name = "tsconfig",
# intrinsic:!enable_npm_core:insert     src = "tsconfig.json",
# intrinsic:!enable_npm_core:insert     deps = ["//intrinsic/frontend/sdk:tsconfig"],
# intrinsic:!enable_npm_core:insert )

npm_package(
    name = "core_package",
    srcs = [
        ":license",
        ":package",
        # intrinsic:!enable_npm_core:insert "//intrinsic/frontend/sdk/core/pubsub:lib",
    ],
    package = "@intrinsic-ai/core",
    publishable = True,
)

alias(
    name = "pkg",
    actual = ":core_package",
)

pkg_tar(
    name = "intrinsic-ai-core",
    srcs = [":core_package"],
    extension = "tgz",
    package_dir = "package",
    strip_prefix = "core_package",
)

alias(
    name = "pack",
    actual = ":intrinsic-ai-core",
)
