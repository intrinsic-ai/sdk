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

/**
 * @file Util file to gather colors shared between components.
 */
import * as three from 'three';

/** Intrinsic color for x axis used by frames and view widget */
export const X_AXIS = new three.Color(0xdd3344);
/** Intrinsic color for y axis used by frames and view widget */
export const Y_AXIS = new three.Color(0x228800);
/** Intrinsic color for z axis used by frames and view widget */
export const Z_AXIS = new three.Color(0x3377cc);

/** Intrinsic color for neutral color used by frames and view widget */
export const NEUTRAL_70 = new three.Color(0xa1a1a1);

/** Intrinsic color for neutral color used by frames and view widget */
export const NEUTRAL_80 = new three.Color(0xbfbfbf);

/** Intrinsic color for neutral color used by frames and view widget */
export const NEUTRAL_99 = new three.Color(0xfbfbfb);

/** Full white */
export const WHITE = new three.Color(0xffffff);

/** Intrinsic color for the world view background in light mode. */
export const VIEWER_BACKGROUND_COLOR_LIGHT = new three.Color(0xf9f9f9);
/** Intrinsic color for the world view background in dark mode. */
export const VIEWER_BACKGROUND_COLOR_DARK = new three.Color(0x343434);

/** Intrinsic color for collider previews */
export const COLLIDER_PREVIEW_COLOR = new three.Color(0x7e4afc);

/** Intrinsic color for world overlay diffs */
export const OVERLAY_DIFF = new three.Color(0xffba05);

/** Intrinsic primary light from colors.scss */
export const PRIMARY_LIGHT = new three.Color(0x5f2ddf);

/** Intrinsic color for success status. */
export const STATUS_COLOR_SUCCESS = new three.Color(0x086238);

/** Intrinsic color for error status. */
export const STATUS_COLOR_ERROR = new three.Color(0xb01b1d);

/** Intrinsic color for warning status. */
export const STATUS_COLOR_WARNING = new three.Color(0xdca91d);
