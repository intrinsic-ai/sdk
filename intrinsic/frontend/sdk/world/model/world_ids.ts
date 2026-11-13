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
 * @file Utility functions for working with stringly-typed IDs.
 */

/**
 * Type synonym for IDs from the entity world.
 */
export type EntityId = number;

/**
 * Type synonym for Object IDs from the object world.
 */
export type ObjectId = number;

/**
 * Type synonym for Frame IDs from the object world.
 */
export type FrameId = number;

/**
 * Type synonym for objects or frames
 */
export type ObjectOrFrameId = number; // ObjectId | FrameId;

/**
 * Id of a undefined entity.
 */
export const NULL_ENTITY_ID = 0;

/**
 * The ID of the root Entity of a world tree (meaning, the root Entity of the root
 * WorldObject of a world tree).
 *
 * The root Entity must be set to have this ID when the world tree is built.
 */
export const ROOT_ENTITY_ID = 1;

/**
 * Id of a undefined object.
 */
export const NULL_OBJECT_ID = 0;

/**
 * The ID of the root WorldObject of a world tree.
 *
 * The root WorldObject must be set to have this ID when the world tree is built.
 */
export const ROOT_OBJECT_ID = 1;

/**
 * Used to separate the object name and frame name in a transform node reference.
 */
export const REFERENCE_STRING_SEPARATOR = '>';

/**
 * Convert from the string based object id to a number based one. This assumes
 * that all object world resource ids are numbers in the form of a string with
 * one of several, underscore-based prefixes (e.g. 'oid_123', 'fid_123', ...),
 * except for the root object which has the id 'root'.
 *
 * The client kept this number based id due to historical reasons, this eases
 * the burdon for the migration from entity based world id to object world id.
 * The corresponding backend conversions can be found in
 * intrinsic/world/objects/object_world_creation_utils.h.
 */
export function stringToNumberId(sid: string): number {
  if (sid === 'root' || sid === 'eid_root') {
    return ROOT_ENTITY_ID;
  } else {
    // Ignore everything up to and including the last occurrence of '_'.
    const val = Number(sid.split('_').pop());
    return Number.isNaN(val) ? NULL_ENTITY_ID : val;
  }
}

/**
 * Return true if the specified identifier belongs to an object or frame.
 */
export function isObjectOrFrameIdentifier(sid: string): boolean {
  return sid === 'root' || sid.startsWith('ofid_');
}

/**
 * Convert from the number based object/frame id to string based for number ids
 * of objects and frames in the world. This is the inverse of stringToNumberId()
 * above for objects and frames.
 */
export function numberToStringIdForObjectOrFrame(id: ObjectOrFrameId) {
  // Object's with id 1 have a special name 'root', other wise they are all
  // numbers.
  if (id === ROOT_ENTITY_ID) return 'root';
  return `ofid_${id}`;
}

/**
 * Convert from the number based entity id to string based for number ids
 * of objects and frames in the world. This is the inverse of stringToNumberId()
 * above.
 */
export function numberToStringIdForEntity(id: EntityId) {
  if (id === ROOT_ENTITY_ID) return 'eid_root';
  return `eid_${id}`;
}
