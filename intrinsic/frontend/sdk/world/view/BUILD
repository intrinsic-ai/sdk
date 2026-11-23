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
        "camera.ts",
        "colors.ts",
        "debug.ts",
        "frame_tracker.ts",
        "geometry_loader_service.ts",
        "geometry_service_shape_loader.ts",
        "index.ts",
        "intr_orbit_controls.ts",
        "intr_transform_controls.ts",
        "materials.ts",
        "object_state_manager.ts",
        "outline_renderer.ts",
        "phantom.ts",
        "picker_manager.ts",
        "renderer.ts",
        "rigid_object_target.ts",
        "scene_updater.ts",
        "shapes.ts",
        "signals.ts",
        "snapping.ts",
        "transform_target.ts",
        "transformer.ts",
        "user_input_mode.ts",
        "user_input_mode_manager.ts",
        "view_controller.ts",
        "view_widget.ts",
        "world_renderer_service.ts",
        "world_root_node.ts",
        "world_scene.ts",
        "world_scene_nodes.ts",
        "world_viewer.ts",
        "world_viewer_options.ts",
        "world_viewer_plugin.ts",
        "world_viewer_styles.ts",
        "world_workspace_service.ts",
    ],
    declaration = True,
    declaration_map = True,
    source_map = True,
    transpiler = "tsc",
    tsconfig = "//intrinsic/frontend/sdk:tsconfig",
    visibility = [
        "//intrinsic/frontend/sdk/world:__pkg__",
    ],
    deps = [
        "//intrinsic/frontend/sdk/world:node_modules/@angular/core",
        "//intrinsic/frontend/sdk/world:node_modules/@intrinsic-ai/angular",
        "//intrinsic/frontend/sdk/world:node_modules/@intrinsic-ai/core",
        "//intrinsic/frontend/sdk/world:node_modules/@intrinsic-ai/types",
        "//intrinsic/frontend/sdk/world:node_modules/@types/three",
        "//intrinsic/frontend/sdk/world:node_modules/rxjs",
        "//intrinsic/frontend/sdk/world:node_modules/three",
        "//intrinsic/frontend/sdk/world:node_modules/three-mesh-bvh",
        "//intrinsic/frontend/sdk/world/editor:lib",
        "//intrinsic/frontend/sdk/world/loaders:lib",
        "//intrinsic/frontend/sdk/world/loaders/basic_asset_depot:lib",
        "//intrinsic/frontend/sdk/world/model:lib",
        "//intrinsic/frontend/sdk/world/service_interfaces:lib",
        "//intrinsic/frontend/sdk/world/utils:lib",
    ],
)
