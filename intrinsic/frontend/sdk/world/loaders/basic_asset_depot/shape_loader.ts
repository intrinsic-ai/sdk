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

import * as three from 'three';

import {GeometryShape} from '../../model/geometry';

/** Loads geometry shapes as threejs objects, e.g. for a `BasicAssetDepot`. */
export interface ShapeLoader {
  /**
   * Returns the key that identifies `shape`, or undefined if `shape` shouldn't
   * be cached.
   *
   * Shapes with the same key are considered the same shape: a shape loaded for
   * one of them is cached and reused for the others. Shapes without a key are
   * loaded every time they are requested.
   */
  getShapeKey(shape: GeometryShape): string | undefined;

  /**
   * Loads `shape` as a threejs object. Rejects if the shape can't be loaded.
   */
  loadShape(shape: GeometryShape): Promise<three.Object3D>;
}
