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
 * @file Utility functions and classes for debugging the world viewer.
 */
import * as three from 'three';

import {isMesh} from '../utils/three_typechecks';
import {IntrVisualMesh} from './shapes';

/**
 * A class to manage and display debug objects in the scene.
 * Use to add boxes, lines, and other shapes to the scene while debugging.
 */
export class DebugObjects {
  private debugObjects: three.Object3D[] = [];

  constructor(private readonly scene: three.Scene) {}

  /**
   * Adds a visual representation of a Box3 to the debug objects.
   * @param box3 The bounding box to visualize.
   * @param color The color of the debug box (optional, defaults to red).
   */
  addBox(box3: three.Box3, color: three.Color = new three.Color(0xff0000)) {
    const boxHelper = new three.Box3Helper(box3, color);
    this.scene.add(boxHelper);
    this.debugObjects.push(boxHelper);
  }

  addBoxesRecursive(
    obj: three.Object3D,
    color: three.Color = new three.Color(0xff0000),
  ) {
    obj.traverse((child) => {
      if (isMesh(child)) {
        const box = new three.Box3().setFromObject(child);
        this.addBox(box, color);
      }
    });
  }

  highlightMeshes(obj: three.Object3D) {
    obj.traverse((child) => {
      if (child instanceof IntrVisualMesh) {
        child.enableColorHighlight(getRandomColor(obj.id));
      }
    });
  }

  /**
   * Removes all debug objects from the scene.
   */
  clear() {
    this.debugObjects.forEach((obj) => {
      if (obj.parent) {
        obj.parent.remove(obj);
      }
      // Assuming debug objects will have geometry and material like Meshes
      if ((obj as three.Mesh).geometry) {
        (obj as three.Mesh).geometry.dispose();
      }
      if ((obj as three.Mesh).material) {
        const material = (obj as three.Mesh).material;
        if (Array.isArray(material)) {
          material.forEach((m) => {
            m.dispose();
          });
        } else {
          material.dispose();
        }
      }
    });
    this.debugObjects = [];
  }
}

/**
 * Generates a random HSL color based on a given integer id.
 * @param id The integer used to generate the color.
 * @returns A three.Color object.
 */
export function getRandomColor(id: number): three.Color {
  // Generate a hue based on the id.
  const hue = ((id * 7) % 360) / 360;
  return new three.Color().setHSL(hue, 0.7, 0.5);
}
