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

import {NEVER, Observable, of} from 'rxjs';

import {WorldModel} from '../model/world_model';
import {WorldEditor} from './world_editor';

/**
 * A WorldEditor that ignores all mutations. Used for read-only worlds, e.g.
 * when no other WorldEditor implementation has been registered.
 */
export class NoopWorldEditor implements WorldEditor {
  readonly pending: Observable<boolean> = of(false);
  readonly acted: Observable<void> = NEVER;

  constructor(readonly worldModel: WorldModel) {}

  async createFrame() {}

  reparent() {}

  toggleAttachmentFrame() {}

  async batchUpdateEntityCollisionSetting() {}

  async updateTransform() {}

  async batchUpdateTransforms() {}

  updateJointValues() {}

  async batchDelete() {
    return false;
  }

  updateName() {}

  async updateWorldCollisionSettings() {}
}
