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
 * @file This module uses object world protos to build the frontend
 * world model.
 */
import type {
  Entity,
  KinematicObjectComponent,
  Object$ as WorldObjectProto,
} from '@intrinsic-ai/types/world/proto/object_world_service_pb';
import type {RobotPayload} from '@intrinsic-ai/types/world/proto/robot_payload_pb';
import * as three from 'three';

import {EntityData, GeometryData} from '../../model/entity';
import * as worldIdUtils from '../../model/world_ids';
import {
  KinematicWorldObjectData,
  WorldObject,
  WorldObjectData,
} from '../../model/world_object';
import {WorldObjectBuilder} from '../../model/world_object_builder';
import {
  geometryComponentToEntityData,
  kinematicsComponentToEntityData,
  physicsComponentToLinkPhysics,
  sensorComponentToEntityData,
} from '../proto/entity_data_proto';
import {matrixdProtoToMatrix3} from '../proto/geometry_proto';
import {poseProtoToMatrix} from '../proto/pose_proto';

function getRobotPayloadFromProto(proto: RobotPayload) {
  const intertiaProto = proto.inertia;
  const inertia = intertiaProto
    ? matrixdProtoToMatrix3(intertiaProto)
    : new three.Matrix3();
  return {
    massKg: proto.massKg,
    tipTCoG: poseProtoToMatrix(proto.tipTCog),
    inertia,
  };
}

function getWorldObjectParamsFromProto(
  proto: WorldObjectProto,
): WorldObjectData {
  const id = worldIdUtils.stringToNumberId(proto.id);
  const name = proto.name;

  let parentEntityId = worldIdUtils.NULL_ENTITY_ID;
  const parentEntityProto = proto.parentEntity;
  if (parentEntityProto) {
    parentEntityId = worldIdUtils.stringToNumberId(
      parentEntityProto.entityReference.case === 'id'
        ? parentEntityProto.entityReference.value
        : '',
    );
  }

  let parentObjectId = worldIdUtils.NULL_OBJECT_ID;
  const parentObjectProto = proto.parent;
  if (parentObjectProto) {
    parentObjectId = worldIdUtils.stringToNumberId(parentObjectProto.id);
  }

  let parentTThis = new three.Matrix4();
  let resourceName: string | undefined;
  let isStatic = false;
  let isDisabled = false;
  let userDataMap = new Map<string, unknown>();
  const objComponentProto = proto.objectComponent;
  if (objComponentProto) {
    const ptt = objComponentProto.parentTThis;
    parentTThis = ptt ? poseProtoToMatrix(ptt) : new three.Matrix4();

    isStatic = !!objComponentProto.simulationComponent?.isStatic;
    isDisabled = !!objComponentProto.simulationComponent?.disabled;
    userDataMap = new Map(Object.entries(objComponentProto.userData));

  }

  const rootEntityId = worldIdUtils.stringToNumberId(proto.rootEntityId);
  const children = proto.children.map((c) => c.name);

  return {
    id,
    name,
    parentEntityId,
    parentObjectId,
    parentTThis,
    isStatic,
    isDisabled,
    resourceName,
    userDataMap,
    rootEntityId,
    children,
  };
}

function getKinematicWorldObjectParamsFromProto(
  proto: WorldObjectProto,
): KinematicWorldObjectData {
  const kc = proto.kinematicObjectComponent;
  if (!kc) {
    throw new Error(
      `Internal error: object '${proto.name}' has no kinematic object component.`,
    );
  }
  const data = getWorldObjectParamsFromProto(proto);
  const payloadProto = kc.mountedPayload;
  return {
    ...data,
    namedConfigurations: getNamedConfigurationsFromProto(kc),
    jointOrdering: getJointOrderingForObject(proto),
    isoFlangeFrames: kc.isoFlangeFrames.map((idName) => idName.name),
    payload: payloadProto ? getRobotPayloadFromProto(payloadProto) : undefined,
  };
}

function getNamedConfigurationsFromProto(
  kc: KinematicObjectComponent,
): Map<string, number[]> {
  const namedConfigs = new Map<string, number[]>();
  for (const config of kc.namedJointConfigurations) {
    namedConfigs.set(config.name, config.jointPositions.slice());
  }
  return namedConfigs;
}

function getJointOrderingForObject(obj: WorldObjectProto): string[] {
  const jointIdsToNames = new Map<string, string>();
  for (const [entityId, entity] of Object.entries(obj.entities)) {
    if (entity.kinematicsComponent) {
      jointIdsToNames.set(entityId, entity.name);
    }
  }
  const jointIdsInOrder = obj.kinematicObjectComponent?.jointEntityIds ?? [];
  return jointIdsInOrder.map((id) => jointIdsToNames.get(id) ?? '');
}

function getEntityParamsForObjectWorldProto(proto: Entity): EntityData {
  const objectProto = proto.object;
  if (!objectProto) {
    throw new Error(`Internal error: entity '${proto.name}' has no object.`);
  }
  const parentIdProto = proto.parentId;
  const poseProto = proto.parentTThis;
  let resourceName: string | undefined;
  return {
    id: worldIdUtils.stringToNumberId(proto.id),
    name: proto.name,
    parentObjectId: worldIdUtils.stringToNumberId(objectProto.id),
    parentTThis: poseProto ? poseProtoToMatrix(poseProto) : new three.Matrix4(),
    parentId: parentIdProto
      ? worldIdUtils.stringToNumberId(parentIdProto)
      : worldIdUtils.NULL_ENTITY_ID,
    resourceName,
  };
}

/** Adapts an `Entity` proto to the interface of `WorldObjectBuilder`. */
interface EntityAdapter {
  readonly proto: Entity;
  getParentId(): string;
}

/** Adapts an `Object` proto to the interface of `WorldObjectBuilder`. */
interface ObjectAdapter {
  readonly proto: WorldObjectProto;
  getEntitiesMap(): ReadonlyMap<string, EntityAdapter>;
}

function toObjectAdapter(proto: WorldObjectProto): ObjectAdapter {
  const entities = new Map<string, EntityAdapter>(
    Object.entries(proto.entities).map(([id, entity]) => [
      id,
      {proto: entity, getParentId: () => entity.parentId},
    ]),
  );
  return {proto, getEntitiesMap: () => entities};
}

/**
 * Builds a world model from the given list of world service protos and returns
 * the root object, if it exists by ID. If not, then it returns any other object
 * created. If objectProtos is empty, returns undefined.
 *
 * In an Object proto, we have a couple of structural requirements that get carried over
 * to the WorldObject built here:
 *
 * - The root Object contains one Entity, which is the root Entity of the scene. It is
 * an empty entity. This gets mapped to a root WorldObject with exactly one FrameEntity
 * (with `isAttachmentFrame`=false) as a child.
 */
export function buildFromProtoList(
  objectProtos: readonly WorldObjectProto[],
): WorldObject | undefined {
  const objectDataProvider = {
    getData: ({proto: objProto}: ObjectAdapter) => {
      const kc = objProto.kinematicObjectComponent;
      if (kc) {
        return getKinematicWorldObjectParamsFromProto(objProto);
      }
      return getWorldObjectParamsFromProto(objProto);
    },
  };

  const entityDataProvider = {
    getData: ({proto: entityProto}: EntityAdapter) => {
      const data = getEntityParamsForObjectWorldProto(entityProto);
      const kinematics = entityProto.kinematicsComponent;
      if (kinematics) {
        return kinematicsComponentToEntityData(data, kinematics);
      }

      const geometry = entityProto.geometryComponent;
      const physics = physicsComponentToLinkPhysics(
        entityProto.physicsComponent,
      );
      if (geometry) {
        return {
          ...geometryComponentToEntityData(data, geometry),
          physics,
        };
      }

      const sensorComponent = entityProto.sensorComponent;
      if (sensorComponent) {
        return sensorComponentToEntityData(data, sensorComponent);
      }

      if (physics) {
        return {
          ...data,
          geometries: new Map<string, Map<string, GeometryData>>(),
          physics,
        };
      }

      // NOTE - the `Entity` proto above does not have any `isAttachmentFrame` field
      // (for the frame case). We lookup the parent WorldObject then get the
      // matching `Frame` proto from its frames list to find the proper value.
      const parentObjectId = entityProto.object?.id;
      let isAttachment = false;
      if (parentObjectId) {
        const parentObjectProto = objectProtos.find(
          (obj) => obj.id === parentObjectId,
        );
        if (parentObjectProto) {
          const matchingFrameProto = parentObjectProto.frames.find(
            (f) => f.name === entityProto.name,
          );
          if (matchingFrameProto) {
            isAttachment = matchingFrameProto.isAttachmentFrame;
          }
        }
      }

      return {
        ...data,
        isAttachment,
      };
    },
  };

  const objMap = new WorldObjectBuilder<EntityAdapter, ObjectAdapter>(
    objectDataProvider,
    entityDataProvider,
  ).build(objectProtos.map(toObjectAdapter));

  const rootObj = objMap.get(worldIdUtils.ROOT_OBJECT_ID);
  if (rootObj) {
    return rootObj;
  }

  // Return any object.
  return objMap.values().next().value;
}
