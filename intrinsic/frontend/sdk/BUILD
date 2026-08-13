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

load("@aspect_rules_js//js:proto.bzl", "js_proto_toolchain")
load("@aspect_rules_ts//ts:defs.bzl", "ts_config")
load("@npm//:defs.bzl", "npm_link_all_packages")
load("@npm//intrinsic/frontend/sdk:@bufbuild/protoc-gen-es/package_json.bzl", "bin")

package(default_visibility = ["//visibility:public"])

exports_files(
    srcs = ["LICENSE"],
    visibility = ["//intrinsic/frontend/sdk:__subpackages__"],
)

npm_link_all_packages(name = "node_modules")

bin.protoc_gen_es_binary(
    name = "protoc_gen_es",
)

js_proto_toolchain(
    name = "gen_es_toolchain",
    out_dts_extension = "_pb.d.ts",
    out_js_extension = "_pb.js",
    plugin_bin = ":protoc_gen_es",
    plugin_name = "es",
    plugin_options = [
        "keep_empty_files=true",
        "target=js+dts",
    ],
    runtime = "//intrinsic/frontend/sdk:node_modules/@bufbuild/protobuf",
    visibility = ["//visibility:public"],
)

ts_config(
    name = "tsconfig",
    src = "tsconfig.json",
    visibility = ["//visibility:public"],
)
