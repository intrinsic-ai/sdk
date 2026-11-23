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

import {NgZone, isDevMode} from '@angular/core';
import * as three from 'three';
import {MeshoptDecoder} from 'three/examples/jsm/libs/meshopt_decoder.module.js';
import {DRACOLoader} from 'three/examples/jsm/loaders/DRACOLoader.js';
import {GLTFLoader} from 'three/examples/jsm/loaders/GLTFLoader.js';
import {MTLLoader} from 'three/examples/jsm/loaders/MTLLoader.js';
import {OBJLoader} from 'three/examples/jsm/loaders/OBJLoader.js';

import {LocalShape, Primitive} from '../model/geometry';
import {isMesh} from '../utils/three_typechecks';
import * as shapes from './shapes';

/**
 * Generally all loaders in Three.js assume resources are loaded over http. We
 * implement a custom texture loader similar to Three.js's TextureLoader, but
 * instead resolving textures from in-memory File handles.
 */
class LocalTextureLoader extends three.TextureLoader {
  private readonly promises: Array<Promise<void>>;
  constructor(private readonly textures?: Map<string, File>) {
    super();
    this.promises = [];
  }

  override load(
    url: string,
    onLoad?: (texture: three.Texture<HTMLImageElement>) => void,
    onProgress?: (event: ProgressEvent) => void,
    onError?: (err: unknown) => void,
  ): three.Texture<HTMLImageElement> {
    const filename = url.split(/\//).pop();
    const texture = new three.Texture<HTMLImageElement>();
    const textureFile = filename ? this.textures?.get(filename) : undefined;
    if (textureFile) {
      this.promises.push(
        new Promise<void>((resolve, reject) => {
          const reader = new FileReader();
          reader.onload = () => {
            const img = new Image();
            if (onError) {
              img.onerror = (e) => {
                onError(new ErrorEvent(`Failed to decode image ${e}`));
              };
              reject();
            }
            img.onload = () => {
              texture.image = img;
              texture.needsUpdate = true;
              texture.colorSpace = three.SRGBColorSpace;
              resolve();
            };
            img.src = reader.result as string;
          };
          if (onError) {
            reader.onerror = (e) => {
              onError(new ErrorEvent(`Failed to read file ${e}`));
            };
          }
          reader.readAsDataURL(textureFile);
        }),
      );
    } else {
      // Generate a stub in the event we don't have the requested texture.
      texture.image = new Image(1, 1);
    }
    return texture;
  }

  /*
   * Returns a promise that resolves when all queued up image decodes complete.
   * Must be called AFTER all texture objects are constructed.
   */
  async getResolved(): Promise<void> {
    return Promise.all(this.promises).then();
  }
}

/**
 * The world-loader Component provides utility methods for loading the world
 * model and renderables used by world-view component to draw the intrinsic world.
 *
 * Currently only support loading gltf/obj models, but can be easily extended to
 * support other 3D formats with other ThreesJs loaders or a custom loader.
 */
export class GeometryLoaderService {
  private readonly gltfLoader: GLTFLoader;

  /**
   * @param ngZone Zone that the loaders run outside of, so that parsing doesn't
   *     trigger change detection.
   */
  constructor(private readonly ngZone: NgZone) {
    this.gltfLoader = new GLTFLoader();
    const decoder = new DRACOLoader();
    // NOTE: this path need to be set correcty by the application at run time.
    decoder.setDecoderPath('../wasm/draco/');
    this.gltfLoader.setDRACOLoader(decoder);
    this.gltfLoader.setMeshoptDecoder(MeshoptDecoder);
  }

  // Load a glTF model from in-memory buffer.
  loadGLTFModelFromMemory(data: ArrayBuffer): Promise<three.Object3D> {
    return new Promise((resolve, reject) => {
      // Schedules asynchronous GTLF parsing to not trigger Angular's change
      // detection when the event finishes.
      this.ngZone.runOutsideAngular(() => {
        this.gltfLoader.parse(
          data,
          '',
          // Called when the resource is loaded
          (gltf) => {
            this.ngZone.run(() => {
              resolve(shapes.createIntrVisuals(gltf.scene));
            });
          },
          // Called when loading has errors
          (error) => {
            if (isDevMode()) {
              console.log(
                'An error happened while parsing gltf data from memory!',
              );
            }
            this.ngZone.run(() => {
              reject(error);
            });
          },
        );
      });
    });
  }

  // load a glTF model from a local shape.
  async loadGLTFModelFromLocalShape(
    shape: LocalShape,
  ): Promise<three.Object3D> {
    const renderable = shape.renderable;
    if (renderable) {
      return this.loadGLTFModelFromMemory(
        renderable.getGlbBytes().buffer.slice(),
      );
    }
    if (!('primitives' in shape.exactGeometry)) {
      throw new Error(
        `Unsupported local shape exact geometry type ${JSON.stringify(shape.exactGeometry)}.
        Only primitives are supported.`,
      );
    }
    // No need to run this synchronous CPU bound code outside Angular.
    const scene = new three.Group();
    for (const transformedPrimitive of shape.exactGeometry.primitives) {
      const primitiveMesh = threeObjectFromPrimitive(
        transformedPrimitive.primitive,
      );
      primitiveMesh.applyMatrix4(transformedPrimitive.refTShape);
      scene.add(primitiveMesh);
    }
    return shapes.createIntrVisuals(scene);
  }

  // Load an OBJ model from text.
  loadOBJModelFromMemory(
    data: string,
    mtlData?: string,
    textures?: Map<string, File>,
  ): Promise<three.Object3D> {
    return new Promise((resolve) => {
      this.ngZone.runOutsideAngular(() => {
        const manager = new three.LoadingManager();
        const objLoader = new OBJLoader();
        const texLoader = new LocalTextureLoader(textures);
        manager.addHandler(/.*/, texLoader);
        if (mtlData) {
          const mtlLoader = new MTLLoader(manager);
          const mtl = mtlData ? mtlLoader.parse(mtlData, '') : undefined;
          if (mtl) {
            objLoader.setMaterials(mtl);
          }
        }
        const obj = objLoader.parse(data);
        // Unfortunately, multiple popular exporters invert the opacity in the
        // MTL file. To be safe we disable any transparency.
        obj.traverse((child: three.Object3D) => {
          if (isMesh(child)) {
            if (Array.isArray(child.material)) {
              for (const mat of child.material) {
                mat.opacity = 1;
              }
            } else {
              child.material.opacity = 1;
            }
          }
        });
        void texLoader.getResolved().then(() => {
          this.ngZone.run(() => {
            resolve(obj);
          });
        });
      });
    });
  }
}

function threeObjectFromPrimitive(p: Primitive): three.Object3D {
  const mat = new three.MeshStandardMaterial();

  if (p.isBox()) {
    const geo = new three.BoxGeometry(p.size.x, p.size.y, p.size.z);
    return new three.Mesh(geo, mat);
  }
  if (p.isCapsule()) {
    const geo = new three.CapsuleGeometry(p.radius, p.length, 16, 32).rotateX(
      Math.PI / 2,
    );
    return new three.Mesh(geo, mat);
  }
  if (p.isCylinder()) {
    const geo = new three.CylinderGeometry(
      p.radius,
      p.radius,
      p.length,
      20,
    ).rotateX(Math.PI / 2);
    return new three.Mesh(geo, mat);
  }
  if (p.isEllipsoid()) {
    const geo = new three.OctahedronGeometry(1, 10);
    geo.scale(p.radii.x, p.radii.y, p.radii.z);
    return new three.Mesh(geo, mat);
  }
  if (p.isFrustum()) {
    const geo = shapes.frustumBufferGeometry(
      p.minZDistance,
      p.maxZDistance,
      p.xAngle,
      p.yAngle,
    );
    return new three.Mesh(geo, mat);
  }
  if (p.isSphere()) {
    const geo = new three.OctahedronGeometry(p.radius, 10);
    return new three.Mesh(geo, mat);
  }
  throw new Error('Unsupported primitive type');
}
