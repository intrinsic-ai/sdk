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

import {
  type JsonValue,
  type MessageInitShape,
  create,
  isMessage,
} from '@bufbuild/protobuf';
import type {
  JointLimits,
  JointLimitsUpdate,
} from '@intrinsic-ai/types/kinematics/types/joint_limits_pb';
import type {Pose} from '@intrinsic-ai/types/math/proto/pose_pb';
import type {CollisionActionSchema} from '@intrinsic-ai/types/world/proto/collision_action_pb';
import type {
  CollisionSettings_CollisionRuleSchema,
  ObjectOrEntityReferenceSchema,
} from '@intrinsic-ai/types/world/proto/collision_settings_pb';
import {
  type FrameReference,
  FrameReferenceSchema,
  type ObjectEntityFilterSchema,
  type ObjectReference,
  type ObjectReferenceSchema,
  type TransformNodeReferenceSchema,
} from '@intrinsic-ai/types/world/proto/object_world_refs_pb';
import {
  type Frame,
  FrameSchema,
  ObjectSchema,
  UpdateWorldResourcesRequestSchema,
  type World,
  type Object$ as WorldObject,
  WorldSchema,
} from '@intrinsic-ai/types/world/proto/object_world_service_pb';
import {
  CreateFrameRequestSchema,
  DeleteFrameRequestSchema,
  DeleteObjectRequestSchema,
  type ObjectWorldUpdate,
  ObjectWorldUpdateSchema,
  ReparentFrameRequestSchema,
  UpdateCollisionSettingsRequestSchema,
  UpdateFrameNameRequestSchema,
  UpdateFramePropertiesRequestSchema,
  UpdateObjectJointsRequestSchema,
  UpdateObjectNameRequestSchema,
  UpdateTransformRequestSchema,
} from '@intrinsic-ai/types/world/proto/object_world_updates_pb';

import type {WorldCollisionRule} from '../model/collision_rules';
import {
  type ObjectId,
  ROOT_OBJECT_ID,
  isObjectOrFrameIdentifier,
  numberToStringIdForObjectOrFrame,
} from '../model/world_ids';
import {
  DeleteObjectPolicy,
  type ObjectWorldService,
} from '../service_interfaces/object_world_service';
import type {HttpTransport} from './http_transport';
import {decodeResponse, getProto, postProto} from './proto_json';

/**
 * An `ObjectWorldService` that uses the platform's object world REST API.
 *
 * Object and frame ids are the numeric ids used by the world model (see
 * `world_ids`).
 *
 * TODO(b/589802342) - This service class lives here temporarily until we figure
 * out a better longer-term place for it.
 */
export class HttpObjectWorldService implements ObjectWorldService {
  constructor(private readonly transport: HttpTransport) {}

  // World

  /**
   * Gets the full world, including all objects and collision settings. Only
   * the world's collision rule set is requested, since entity collision
   * exclusions are not displayed or editable in the world viewer.
   *
   * TODO(b/589812258) - Currently, all user-data is stripped by Objects, to
   * avoid decoding failures.
   */
  async getWorld(worldId: string): Promise<World> {
    return getProto(
      this.transport,
      `api/objectworlds/${worldId}?view=WORLD_VIEW_FULL` +
        `&collision_settings_view=COLLISION_SETTINGS_VIEW_RULE_SET`,
      WorldSchema,
      worldWithoutUserData,
    );
  }

  // Transforms

  /**
   * Updates the transform of `id` (or of `updateId`, if given) such that the
   * pose of `id` relative to `refId` is `relativePose`.
   */
  async updateTransform(
    worldId: string,
    id: number,
    refId: number,
    relativePose: Pose,
    updateId?: number,
  ): Promise<void> {
    const request = create(UpdateTransformRequestSchema, {
      nodeA: transformNodeRef(refId),
      nodeB: transformNodeRef(id),
      nodeToUpdate: transformNodeRef(updateId ?? id),
      aTB: relativePose,
    });
    await postProto(
      this.transport,
      `api/objectworlds/${worldId}:updateTransform`,
      UpdateTransformRequestSchema,
      request,
    );
  }

  /** Sets the world poses of several objects or frames, keyed by id. */
  async batchUpdateTransform(
    worldId: string,
    worldPoses: Map<number, Pose>,
  ): Promise<void> {
    const updates = [...worldPoses].map(([id, worldPose]) =>
      updateWorldPose(id, worldPose),
    );
    await this.batchUpdateResources(worldId, updates);
  }

  // Objects

  /** Renames an object. */
  async updateObjectName(
    worldId: string,
    oId: number,
    name: string,
  ): Promise<WorldObject> {
    const request = create(UpdateObjectNameRequestSchema, {
      name,
      nameIsGlobalAlias: true,
    });
    const response = await postProto(
      this.transport,
      `api/objectworlds/${worldId}/objects/${idString(oId)}:updateName`,
      UpdateObjectNameRequestSchema,
      request,
    );
    return decodeResponse(response, ObjectSchema, objectWithoutUserData);
  }

  /**
   * Moves an object to a new parent object, optionally attaching it to the
   * parent entity named `targetEntityName` and setting its world pose.
   *
   * If `targetEntityName` is not given, the object is attached to the final
   * entity of the new parent.
   */
  async updateObjectParentAndTransform(
    worldId: string,
    id: ObjectId,
    pId: ObjectId,
    targetEntityName?: string,
    worldPose?: Pose,
  ): Promise<void> {
    const updates: ObjectWorldUpdate[] = [
      create(ObjectWorldUpdateSchema, {
        update: {
          case: 'reparentObject',
          value: {
            object: objectRef(id),
            attachTo: {
              case: 'parentObject',
              value: {
                reference: objectRef(pId),
                entityFilter: parentEntityFilter(targetEntityName),
              },
            },
          },
        },
      }),
    ];
    if (worldPose) {
      updates.push(updateWorldPose(id, worldPose));
    }
    await this.batchUpdateResources(worldId, updates);
  }

  /** Deletes an object, handling its children according to `policy`. */
  async deleteObject(
    worldId: string,
    oId: number,
    policy: DeleteObjectPolicy,
  ): Promise<void> {
    const {force, reparentChildObjects} = deleteObjectPolicyFields(policy);
    const params: string[] = [];
    if (force) {
      params.push('force');
    }
    if (reparentChildObjects) {
      params.push('reparentChildObjects');
    }
    const query = params.length > 0 ? `?${params.join('&')}` : '';
    await postProto(
      this.transport,
      `api/objectworlds/${worldId}/objects/${idString(oId)}:delete${query}`,
      DeleteObjectRequestSchema,
      create(DeleteObjectRequestSchema),
    );
  }

  // Frames

  /** Creates a new frame under the given parent frame or object. */
  async createFrame(
    worldId: string,
    parentReference: FrameReference | ObjectReference,
    name: string,
    relativePose: Pose,
  ): Promise<Frame> {
    const request = create(CreateFrameRequestSchema, {
      newFrameName: name,
      parentTNewFrame: relativePose,
      attachTo: isMessage(parentReference, FrameReferenceSchema)
        ? {case: 'parentFrame', value: parentReference}
        : {
            case: 'parentObjectWithFilter',
            value: {reference: parentReference},
          },
    });
    const response = await postProto(
      this.transport,
      `api/objectworlds/${worldId}/frames`,
      CreateFrameRequestSchema,
      request,
    );
    return decodeResponse(response, FrameSchema);
  }

  /** Renames a frame. */
  async updateFrameName(
    worldId: string,
    fId: number,
    name: string,
  ): Promise<Frame> {
    const request = create(UpdateFrameNameRequestSchema, {name});
    const response = await postProto(
      this.transport,
      `api/objectworlds/${worldId}/frames/${idString(fId)}:updateName`,
      UpdateFrameNameRequestSchema,
      request,
    );
    return decodeResponse(response, FrameSchema);
  }

  /** Updates the properties of a frame. */
  async updateFrameProps(
    worldId: string,
    fId: number,
    isAttachmentFrame: boolean,
  ): Promise<Frame> {
    const request = create(UpdateFramePropertiesRequestSchema, {
      isAttachmentFrame,
    });
    const response = await postProto(
      this.transport,
      `api/objectworlds/${worldId}/frames/${idString(fId)}:updateProps`,
      UpdateFramePropertiesRequestSchema,
      request,
    );
    return decodeResponse(response, FrameSchema);
  }

  /**
   * Moves a frame to the parent object `pId`, optionally attaching it to the
   * parent entity named `entityTarget`.
   *
   * If `entityTarget` is not given, the frame is attached to the final entity
   * of the new parent.
   */
  async updateFrameParent(
    worldId: string,
    fId: number,
    pId: number,
    entityTarget?: string,
  ): Promise<void> {
    const request = create(ReparentFrameRequestSchema, {
      attachTo: {
        case: 'parentObject',
        value: {
          reference: objectRef(pId),
          entityFilter: parentEntityFilter(entityTarget),
        },
      },
    });
    await postProto(
      this.transport,
      `api/objectworlds/${worldId}/frames/${idString(fId)}:reparent`,
      ReparentFrameRequestSchema,
      request,
    );
  }

  /** Deletes a frame. If `force` is set, its child frames are deleted too. */
  async deleteFrame(
    worldId: string,
    fId: number,
    force: boolean,
  ): Promise<void> {
    const query = force ? '?force' : '';
    await postProto(
      this.transport,
      `api/objectworlds/${worldId}/frames/${idString(fId)}:delete${query}`,
      DeleteFrameRequestSchema,
      create(DeleteFrameRequestSchema),
    );
  }

  // Joints

  /**
   * Sets the joint positions of a robot. Resolves to the joint positions after
   * the update.
   */
  async updateJoints(
    worldId: string,
    robotId: number,
    jvals: number[],
  ): Promise<readonly number[]> {
    const object = await this.updateObjectJoints(worldId, robotId, {
      jointPositions: jvals,
    });
    return object.kinematicObjectComponent?.jointPositions ?? [];
  }

  /**
   * Updates the joint application limits of a robot. Resolves to the limits
   * after the update, if any.
   */
  async updateJointApplicationLimits(
    worldId: string,
    robotId: number,
    newLimits: JointLimitsUpdate,
  ): Promise<JointLimits | undefined> {
    const object = await this.updateObjectJoints(worldId, robotId, {
      jointApplicationLimits: newLimits,
    });
    return object.kinematicObjectComponent?.jointApplicationLimits;
  }

  // Collisions

  /** Replaces the collision rules of the world. */
  async updateWorldCollisionSettings(
    worldId: string,
    collisionRules: WorldCollisionRule[],
  ): Promise<void> {
    const request = create(UpdateCollisionSettingsRequestSchema, {
      worldId,
      type: {
        case: 'collisionSettings',
        value: {
          disableCollisionChecking: false,
          collisionRules: collisionRules.map(toCollisionRule),
        },
      },
    });
    await postProto(
      this.transport,
      `api/objectworlds/${worldId}:updateCollisionSettings`,
      UpdateCollisionSettingsRequestSchema,
      request,
    );
  }

  private async updateObjectJoints(
    worldId: string,
    robotId: number,
    init: MessageInitShape<typeof UpdateObjectJointsRequestSchema>,
  ): Promise<WorldObject> {
    const request = create(UpdateObjectJointsRequestSchema, init);
    const response = await postProto(
      this.transport,
      `api/objectworlds/${worldId}/objects/${idString(robotId)}:updateJoints?view=FULL`,
      UpdateObjectJointsRequestSchema,
      request,
    );
    return decodeResponse(response, ObjectSchema, objectWithoutUserData);
  }

  private async batchUpdateResources(
    worldId: string,
    updates: ObjectWorldUpdate[],
  ): Promise<void> {
    const request = create(UpdateWorldResourcesRequestSchema, {
      worldUpdates: {updates},
    });
    await postProto(
      this.transport,
      `api/objectworlds/${worldId}:batchUpdateResources`,
      UpdateWorldResourcesRequestSchema,
      request,
    );
  }
}

/**
 * Removes the user data from the JSON of an object.
 *
 * User data is a map of `google.protobuf.Any` messages, which can only be
 * decoded from JSON if all of their types are known.
 *
 * TODO(b/589812258) - We currently remove the user-data until we figure out a
 * better long-term solution for decoding it.
 */
function objectWithoutUserData(object: JsonValue): JsonValue {
  if (isJsonObject(object) && isJsonObject(object['objectComponent'])) {
    delete object['objectComponent']['userData'];
  }
  return object;
}

/** Removes the user data from the JSON of all objects in a world. */
function worldWithoutUserData(world: JsonValue): JsonValue {
  if (isJsonObject(world) && Array.isArray(world['objects'])) {
    world['objects'].forEach(objectWithoutUserData);
  }
  return world;
}

function isJsonObject(json: JsonValue | undefined): json is {
  [key: string]: JsonValue;
} {
  return typeof json === 'object' && json !== null && !Array.isArray(json);
}

function idString(id: number): string {
  return numberToStringIdForObjectOrFrame(id);
}

function transformNodeRef(
  id: number,
): MessageInitShape<typeof TransformNodeReferenceSchema> {
  return {transformNodeReference: {case: 'id', value: idString(id)}};
}

function objectRef(id: number): MessageInitShape<typeof ObjectReferenceSchema> {
  return {objectReference: {case: 'id', value: idString(id)}};
}

/**
 * The entity filter selecting the entity of a parent object to attach to: the
 * entity named `entityName` if given, otherwise the parent's final entity.
 */
function parentEntityFilter(
  entityName?: string,
): MessageInitShape<typeof ObjectEntityFilterSchema> {
  if (entityName !== undefined) {
    return {entityNames: [entityName]};
  }
  return {includeFinalEntity: true, includeBaseEntity: false};
}

/** An update setting the world pose of the object or frame `id`. */
function updateWorldPose(id: number, worldPose: Pose): ObjectWorldUpdate {
  return create(ObjectWorldUpdateSchema, {
    update: {
      case: 'updateTransform',
      value: {
        nodeA: transformNodeRef(ROOT_OBJECT_ID),
        nodeB: transformNodeRef(id),
        nodeToUpdate: transformNodeRef(id),
        aTB: worldPose,
      },
    },
  });
}

function toCollisionRule(
  rule: WorldCollisionRule,
): MessageInitShape<typeof CollisionSettings_CollisionRuleSchema> {
  // `isExcluded` and `margin` are alternatives of the same oneof, so a margin
  // (if set) takes precedence over the exclusion.
  const collisionAction: MessageInitShape<typeof CollisionActionSchema> =
    rule.margin
      ? {action: {case: 'margin', value: {hardMargin: rule.margin}}}
      : {action: {case: 'isExcluded', value: rule.excluded}};
  return {
    collisionAction,
    left: [objectOrEntityRef(rule.leftId)],
    right: rule.rightId !== undefined ? [objectOrEntityRef(rule.rightId)] : [],
  };
}

function objectOrEntityRef(
  id: string,
): MessageInitShape<typeof ObjectOrEntityReferenceSchema> {
  if (isObjectOrFrameIdentifier(id)) {
    return {
      type: {case: 'object', value: {objectReference: {case: 'id', value: id}}},
    };
  }
  return {
    type: {case: 'entity', value: {entityReference: {case: 'id', value: id}}},
  };
}

/**
 * The `DeleteObjectRequest` fields that express a `DeleteObjectPolicy`.
 *
 * The world service rejects a request that sets both, so the two flags are
 * always derived together from a single policy.
 */
function deleteObjectPolicyFields(policy: DeleteObjectPolicy): {
  force: boolean;
  reparentChildObjects: boolean;
} {
  switch (policy) {
    case DeleteObjectPolicy.DeleteIfNoChildren:
      return {force: false, reparentChildObjects: false};
    case DeleteObjectPolicy.DeleteRecursively:
      return {force: true, reparentChildObjects: false};
    case DeleteObjectPolicy.DeleteAndReparentChildren:
      return {force: false, reparentChildObjects: true};
  }
}
