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
 * @file Conversions from protobuf-es CollisionSettings protos to native
 * collision rules.
 */
import type {
  CollisionSettings,
  ObjectOrEntityReference,
} from '@intrinsic-ai/types/world/proto/collision_settings_pb';

import type {WorldCollisionRule} from '../../model/collision_rules';

function objectOrEntityReferenceToId(
  ref: ObjectOrEntityReference,
): string | undefined {
  switch (ref.type.case) {
    case 'object':
      return ref.type.value.objectReference.case === 'id'
        ? ref.type.value.objectReference.value
        : '';
    case 'entity':
      return ref.type.value.entityReference.case === 'id'
        ? ref.type.value.entityReference.value
        : '';
    default:
      return undefined;
  }
}

function objectOrEntityReferencesToIds(
  refs: readonly ObjectOrEntityReference[],
): string[] {
  return refs
    .map(objectOrEntityReferenceToId)
    .filter((id): id is string => id !== undefined);
}

/**
 * Converts a CollisionSettings proto to a list of native WorldCollisionRules.
 *
 * Each proto rule can refer to two groups of one or more objects each (details
 * are described in the CollisionSettings proto). To simplify UX displays and
 * rule modifications, the result only contains 1:1 or 1:anything rules, so a
 * single proto rule may expand into several WorldCollisionRules.
 *
 * References that cannot be converted to an ID are skipped.
 */
export function collisionSettingsToWorldCollisionRules(
  proto: CollisionSettings,
): WorldCollisionRule[] {
  const rules: WorldCollisionRule[] = [];
  for (const rule of proto.collisionRules) {
    const leftIds = objectOrEntityReferencesToIds(rule.left);
    const rightIds = objectOrEntityReferencesToIds(rule.right);
    const action = rule.collisionAction?.action;
    const excluded = action?.case === 'isExcluded' ? action.value : false;
    const margin = action?.case === 'margin' ? action.value.hardMargin : 0;
    for (const leftId of leftIds) {
      // Only an empty right list means "anything". If the right list had
      // entries that could not be converted, don't widen the rule.
      if (rule.right.length === 0) {
        rules.push({leftId, excluded, margin});
      } else {
        for (const rightId of rightIds) {
          rules.push({leftId, rightId, excluded, margin});
        }
      }
    }
  }
  return rules;
}
