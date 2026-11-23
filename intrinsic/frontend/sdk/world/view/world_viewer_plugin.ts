// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import type {WorldViewer} from './world_viewer';

/**
 * Interface for plugins that hook into the WorldViewer component. See
 * the `registerPlugin`/`unregisterPlugin` method on WorldViewer.
 *
 * Use WorldViewerPlugins for adding functionality like overlays, visual markers,
 * etc.
 */
export interface WorldViewerPlugin {
  /**
   * Called when the plugin is registered with a WorldViewer.
   */
  onRegister(worldViewer: WorldViewer): void;

  /**
   * Called when the plugin is unregistered from a WorldViewer.
   */
  onUnregister(worldViewer: WorldViewer): void;
}
