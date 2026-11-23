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

import {ShapeLoader} from '../loaders/basic_asset_depot/shape_loader';
import {GeometryShape, RemoteShape, isRemoteShape} from '../model/geometry';
import {GeometryService} from '../service_interfaces/geometry_service';
import {GeometryLoaderService} from './geometry_loader_service';

/**
 * A `ShapeLoader` that gets the renderables of remote shapes from the geometry
 * service (through a `GeometryService`) and converts them, and local
 * shapes, to threejs objects with a `GeometryLoaderService`.
 */
export class GeometryServiceShapeLoader implements ShapeLoader {
  constructor(
    private readonly geometryService: GeometryService,
    private readonly geometryLoaderService: GeometryLoaderService,
  ) {}

  /**
   * Remote shapes are identified by their renderable ref, or by their exact
   * geometry ref if they don't have a renderable yet. Local shapes have no
   * identity, so they aren't cached.
   */
  getShapeKey(shape: GeometryShape): string | undefined {
    if (!isRemoteShape(shape)) {
      return undefined;
    }
    if (shape.renderableRef !== '') {
      return shape.renderableRef;
    }
    return shape.exactGeometryRef;
  }

  async loadShape(shape: GeometryShape): Promise<three.Object3D> {
    if (!isRemoteShape(shape)) {
      return this.geometryLoaderService.loadGLTFModelFromLocalShape(shape);
    }
    return this.geometryLoaderService.loadGLTFModelFromMemory(
      await this.fetchGlb(shape),
    );
  }

  /** Fetches the renderable of `shape` as GLB (binary glTF) bytes. */
  private async fetchGlb(shape: RemoteShape): Promise<ArrayBuffer> {
    const renderable = await this.geometryService.getRenderable(shape);
    const glbBytes = renderable.renderable?.glbBytes;
    if (!glbBytes) {
      throw new Error(
        `Renderable for renderable ref '${shape.renderableRef}', exact ` +
          `geometry ref '${shape.exactGeometryRef}' contains no GLB bytes.`,
      );
    }
    // The bytes may be a view into a larger buffer, so copy exactly the GLB
    // bytes into a buffer of their own.
    return glbBytes.slice().buffer;
  }
}
