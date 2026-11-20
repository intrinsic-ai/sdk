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

import {Observable} from 'rxjs';
import * as three from 'three';

import {GeometryShape, RemoteShape} from '../model/geometry';

/**
 * The AssetDepot is an abstract class for retrieving remote geometry for use
 * in the WorldViewer's 3D scene.
 *
 * It is meant to be flexible enough to support retrieving and loading any kind
 * of geometry, from any source, so long as the geometry can be converted to a
 * threejs object.
 */
export abstract class AssetDepot {
  /**
   * Emits whenever a new shape has been loaded into the depot.
   */
  abstract readonly shapeAdded: Observable<void>;

  /**
   * Retrieves or loads the 3D Three.js Object for the given geometry shape.
   *
   * In a typical implementation, the renderableRef and/or exactGeometryRef fields
   * of `shape` will contain some URI to a remote resource. This should be fetched
   * here, cached, and converted to a threejs object.
   */
  abstract getShape(shape: GeometryShape): Promise<three.Object3D>;

  /**
   * Returns true if the remote shape is already cached in the depot.
   */
  abstract hasShape(shape: RemoteShape): boolean;

  /**
   * Returns the current download and loading progress percentage (0 to 100)
   * of all currently pending operations.
   */
  abstract percentLoaded(): number;

  /**
   * Returns the total number of cached shapes stored in the depot.
   */
  abstract size(): number;
}
