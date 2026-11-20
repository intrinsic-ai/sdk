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

import {CollisionRules} from '../model/collision_rules';
import {Entity, FrameEntity} from '../model/entity';
import {FrameId} from '../model/world_ids';
import {WorldModel} from '../model/world_model';
import {WorldObject} from '../model/world_object';

/**
 * Centralizes the mutation operations on a WorldModel, e.g. for the world
 * viewer and its tools. Note that it is not necessarily the only way to mutate
 * a WorldModel.
 *
 * Implementations decide how the mutations are carried out (e.g., as
 * undoable commands). The owner of a workspace chooses its implementation when
 * creating it (see WorldWorkspaceService.createWorkspace).
 */
export interface WorldEditor {
  /** The model that this editor mutates. */
  readonly worldModel: WorldModel;

  /** Whether the editor has a mutation in progress. */
  readonly pending: Observable<boolean>;

  /**
   * Notifies subscribers when the editor acts, e.g. by adding, modifying or
   * removing objects or frames.
   */
  readonly acted: Observable<void>;

  /**
   * Creates a new frame under an object.
   *
   * @param oid The id of the object to create the frame under.
   * @param name The name of the new frame.
   * @param position The position of the new frame relative to the object.
   * @param quaternion The rotation of the new frame relative to the object.
   * @param focusOnNewFrame Whether to select the new frame once it's created.
   */
  createFrame(
    oid: number,
    name: string,
    position: three.Vector3,
    quaternion: three.Quaternion,
    focusOnNewFrame?: boolean,
  ): Promise<void>;

  /**
   * Reparents an object or frame under a different object or entity.
   *
   * @param subject The object or frame being reparented.
   * @param target The new parent.
   */
  reparent(
    subject: WorldObject | FrameEntity,
    target: WorldObject | Entity,
  ): void;

  /**
   * Toggles whether a frame is an attachment frame.
   *
   * @param fid The id of the frame to toggle.
   */
  toggleAttachmentFrame(fid: FrameId): void;

  /**
   * Updates the collision settings of several entities.
   *
   * @param ids The ids of the entities to update.
   * @param collisionMargin The collision margin to set for the entities.
   * @param collisionsEnabled Whether collisions are enabled for the entities.
   *     When false, collisionMargin has no effect.
   */
  batchUpdateEntityCollisionSetting(
    ids: string[],
    collisionMargin: number,
    collisionsEnabled: boolean,
  ): Promise<void>;

  /** Updates the transform of an object. */
  updateTransform(oid: number, transform: three.Matrix4): Promise<void>;

  /** Updates the transforms of several objects, keyed by object id. */
  batchUpdateTransforms(transforms: Map<number, three.Matrix4>): Promise<void>;

  /** Updates joint values of a robot, keyed by joint id. */
  updateJointValues(
    robotObjectId: number,
    jointValues: Map<number, number>,
  ): void;

  /**
   * Deletes several objects or frames.
   *
   * @returns whether the objects or frames were deleted.
   */
  batchDelete(ofids: number[]): Promise<boolean>;

  /**
   * Updates the name of an object or frame.
   *
   * @param id The id of the object or frame.
   * @param name The new name.
   */
  updateName(id: number, name: string): void;

  /** Replaces the world's collision rules. */
  updateWorldCollisionSettings(collisionRules: CollisionRules): Promise<void>;
}
