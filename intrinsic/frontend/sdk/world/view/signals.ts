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

import {Signal, computed} from '@angular/core';

import {ObjectId} from '../model/world_ids';
import {WorldObject} from '../model/world_object';
import {WorldViewer} from './world_viewer';

/**
 * Returns a signal that produces the most recent object with the given ID.
 */
export function toObjectSignalForWorldViewer(
  worldViewer: Signal<WorldViewer>,
  objectId: Signal<ObjectId>,
): Signal<WorldObject | undefined> {
  return computed(() => {
    const worldModelSignal = worldViewer().getWorldModelSignal();
    const model = worldModelSignal();
    if (!model) {
      return undefined;
    }
    const rootObjectSignal = model.getRootObjectSignal();
    return rootObjectSignal()?.findObjectInSubtree(objectId());
  });
}
