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

import type {ReadonlyPose} from '@intrinsic-ai/core/math';
import type {SceneObjectInstanceUpdate} from '@intrinsic-ai/types/scene/proto/v1/scene_object_updates_pb';
import type {WorldStateUpdate} from '@intrinsic-ai/types/world/proto/world_state_pb';
import type {Observable} from 'rxjs';

import type {CollisionPair} from '../../model/world_backend';

/**
 * A snapshot of the entity poses and joint positions of a world. It doesn't
 * describe the whole world, e.g. its objects and their structure.
 */
export type {WorldStateUpdate};

/**
 * Streams live data about worlds from the backend, e.g. over a pubsub
 * connection.
 */
export interface WorldStreamingSource {
  /**
   * Returns a stream of the state of the world with the given id. Each value is
   * the latest state of the world, e.g. its entity poses and joint positions.
   *
   * The stream may error or complete, e.g. when the connection drops. Callers
   * are expected to resubscribe.
   */
  subscribeToWorldState(worldId: string): Observable<WorldStateUpdate>;
}

/**
 * Applies updates (e.g. geometry or simulation properties) to a scene object
 * instance.
 */
export interface SceneObjectUpdater {
  /**
   * Applies `update` to the scene object instance with the given name.
   * Resolves once the update has been applied.
   */
  updateSceneObject(
    name: string,
    update: SceneObjectInstanceUpdate,
  ): Promise<void>;
}

/**
 * Deletes asset instances. In the ObjectWorldService backend, every world
 * object is an asset instance. So to delete them we need to be able to delete
 * asset instances (which also removes the object from the world).
 */
export interface AssetInstanceDeleter {
  /**
   * Deletes the asset instance with the given name from the given world.
   * Resolves once the deletion is complete, and rejects if it fails.
   */
  deleteAssetInstance(name: string, worldId: string): Promise<void>;
}

/** Checks worlds for collisions between entities. */
export interface CollisionChecker {
  /**
   * Checks the world with the given id for collisions. Resolves to the pairs of
   * entities that are in collision.
   */
  checkCollisions(worldId: string): Promise<readonly CollisionPair[]>;
}

/** Computes inverse kinematics (IK) solutions for robots in worlds. */
export interface IkSolver {
  /**
   * Gets an IK solution for the given tip frame of a robot and target pose.
   * Resolves to a map from joint entity id to joint value.
   *
   * @param worldId The id of the world that the robot belongs to.
   * @param tipId The id of the tip frame on the robot to solve IK for.
   * @param targetPose Target world pose of the tip frame.
   */
  getIkSolution(
    worldId: string,
    tipId: number,
    targetPose: ReadonlyPose,
  ): Promise<Map<number, number>>;
}
