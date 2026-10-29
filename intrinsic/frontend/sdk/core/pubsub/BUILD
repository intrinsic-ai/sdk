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

load("@aspect_rules_ts//ts:defs.bzl", "ts_project")

package(default_visibility = ["//visibility:public"])

ts_project(
    name = "lib",
    srcs = [
        "src/chunking.ts",
        "src/client.ts",
        "src/control_plane.ts",
        "src/data_plane.ts",
        "src/diagnostics.ts",
        "src/errors.ts",
        "src/index.ts",
        "src/publisher.ts",
        "src/signaling_controller.ts",
        "src/subscriber.ts",
        "src/types.ts",
    ],
    declaration = True,
    declaration_map = True,
    out_dir = "dist",
    root_dir = "src",
    source_map = True,
    transpiler = "tsc",
    tsconfig = "//intrinsic/frontend/sdk/core:tsconfig",
    deps = [
        "//intrinsic/frontend/sdk/core:node_modules/@bufbuild/protobuf",
        "//intrinsic/frontend/sdk/core:node_modules/@intrinsic-ai/types",
    ],
)
