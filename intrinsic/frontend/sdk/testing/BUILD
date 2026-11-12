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
load("@aspect_rules_ts//ts:defs.bzl", "ts_project")
load("@npm//:defs.bzl", "npm_link_all_packages")

# This is a workspace package of internal test helpers. It is not published to
# npm. Other packages in the workspace that need access to test utilities (such
# as `clean_state.ts`) should import it using workspace links. Note that this
# package is specifically for test utils that should not be exposed to the npm
# package users. For things like spies, mocks, etc that are useful for
# npm-package users, the test code should live in a `/testing` subdirectory of
# the relevant package, such as `core/testing` and export the test functions
# that it should make available.

npm_link_all_packages(name = "node_modules")

ts_project(
    name = "lib",
    srcs = [
        "clean_state.ts",
        "defer_promise.ts",
        "index.ts",
        "observable.ts",
        "parameterize.ts",
    ],
    declaration = True,
    declaration_map = True,
    out_dir = "dist",
    source_map = True,
    transpiler = "tsc",
    tsconfig = "//intrinsic/frontend/sdk:tsconfig",
    visibility = ["//intrinsic:internal_api_users"],
    deps = [":node_modules/rxjs"],
)

# Linked into other workspace packages (e.g. `"@intrinsic-ai/testing":
# "workspace:*"`) by rules_js, which uses the `pkg` target of first-party
# workspace packages. Never published.
npm_package(
    name = "pkg",
    srcs = [
        "package.json",
        ":lib",
    ],
    package = "@intrinsic-ai/testing",
    publishable = False,
    visibility = ["//intrinsic:internal_api_users"],
)
