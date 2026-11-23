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

import {InjectionToken} from '@angular/core';
import type * as three from 'three';

/**
 * An environment map, which lights the scene and shows in reflections: the URL
 * of an equirectangular HDR (`.hdr`) image, or an equirectangular texture.
 */
export type EnvironmentMapSource = string | three.Texture;

/** Options that configure all world viewers in an application. */
export interface WorldViewerOptions {
  /** Shows a frames-per-second counter in each view. */
  readonly showFps: boolean;
  /** Disables WebGL rendering entirely, e.g. in environments without a GPU. */
  readonly suppressRendering: boolean;
  /** Fails instead of falling back to software (CPU) rendering. */
  readonly noSoftwareRendering: boolean;
  /**
   * The environment map of all views. If unset, the scene is lit by a generic
   * room environment.
   *
   * TODO(b/588515838): Eventually this will be settable per WorldViewer.
   */
  readonly environmentMap?: EnvironmentMapSource;
}

/** The default world viewer options. */
export const DEFAULT_WORLD_VIEWER_OPTIONS: WorldViewerOptions = {
  showFps: false,
  suppressRendering: false,
  noSoftwareRendering: false,
};

/**
 * Injection token for the world viewer options. Provide it at the application
 * level to override the defaults, e.g.:
 *
 * ```ts
 * providers: [
 *   {
 *     provide: WORLD_VIEWER_OPTIONS,
 *     useValue: {...DEFAULT_WORLD_VIEWER_OPTIONS, showFps: true},
 *   },
 * ]
 * ```
 */
export const WORLD_VIEWER_OPTIONS = new InjectionToken<WorldViewerOptions>(
  'WorldViewerOptions',
  {providedIn: 'root', factory: () => DEFAULT_WORLD_VIEWER_OPTIONS},
);
