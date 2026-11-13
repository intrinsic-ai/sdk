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

ts_project(
    name = "lib",
    srcs = [
        "collision_rules.ts",
        "entity.ts",
        "geometry.ts",
        "index.ts",
        "joint_limits.ts",
        "material_properties.ts",
        "physics.ts",
        "wait_for_world_model.ts",
        "world_backend.ts",
        "world_ids.ts",
        "world_lookup.ts",
        "world_model.ts",
        "world_object.ts",
        "world_object_builder.ts",
    ],
    declaration = True,
    declaration_map = True,
    source_map = True,
    transpiler = "tsc",
    tsconfig = "//intrinsic/frontend/sdk:tsconfig",
    visibility = [
        "//intrinsic/frontend/sdk/world:__pkg__",
        "//intrinsic/frontend/sdk/world/backends/simple_world_backend:__pkg__",
    ],
    deps = [
        "//intrinsic/frontend/sdk/world:node_modules/@angular/core",
        "//intrinsic/frontend/sdk/world:node_modules/@intrinsic-ai/core",
        "//intrinsic/frontend/sdk/world:node_modules/@types/three",
        "//intrinsic/frontend/sdk/world:node_modules/rxjs",
        "//intrinsic/frontend/sdk/world:node_modules/three",
        "//intrinsic/frontend/sdk/world/utils:lib",
    ],
)
