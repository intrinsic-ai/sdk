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

import {Subject} from 'rxjs';
import * as three from 'three';

import {GeometryShape, RemoteShape} from '../../model/geometry';
import {AsyncCache} from '../../utils/async_cache';
import {AssetDepot} from '../asset_depot';
import {ShapeLoader} from './shape_loader';

/**
 * An `AssetDepot` that loads shapes with a `ShapeLoader`.
 *
 * BasicAssetDepot is meant to be a flexible starting point for implementing an
 * AssetDepot. It allows passing in a `ShapeLoader` that implements the actual
 * backend fetch of the geometry as well as loading the geometry to an Object3D
 * for threejs use.
 *
 * Shapes are identified by the key that the `ShapeLoader` returns for them.
 * Shapes with a key are cached, so each of them is only loaded once, and
 * concurrent requests for them share a single load. Shapes without a key are
 * loaded every time they are requested.
 */
export class BasicAssetDepot extends AssetDepot {
  private readonly shapes = new AsyncCache<string, three.Object3D>();
  private readonly shapeAddedSubject = new Subject<void>();
  readonly shapeAdded = this.shapeAddedSubject.asObservable();

  /** @param shapeLoader Identifies and loads shapes as threejs objects. */
  constructor(private readonly shapeLoader: ShapeLoader) {
    super();
  }

  async getShape(shape: GeometryShape): Promise<three.Object3D> {
    const key = this.shapeLoader.getShapeKey(shape);
    if (key === undefined) {
      return this.shapeLoader.loadShape(shape);
    }
    const isNewLoad = !this.shapes.has(key);
    const object = await this.shapes.getOrLoad(key, () =>
      this.shapeLoader.loadShape(shape),
    );
    if (isNewLoad) {
      this.shapeAddedSubject.next();
    }
    // Each caller gets its own copy of the cached object, as an object can only
    // be in one place in a scene graph.
    return object.clone();
  }

  hasShape(shape: RemoteShape): boolean {
    const key = this.shapeLoader.getShapeKey(shape);
    return key !== undefined && this.shapes.isLoaded(key);
  }

  percentLoaded(): number {
    if (this.shapes.size() === 0) return 100;
    return (100 * this.shapes.getNumLoaded()) / this.shapes.size();
  }

  size(): number {
    return this.shapes.getNumLoaded();
  }
}
