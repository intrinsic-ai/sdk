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

import {clone, create} from '@bufbuild/protobuf';
import {cyrb53} from '@intrinsic-ai/core/hash';
import {assumeExhaustive} from '@intrinsic-ai/core/typescript';
import {
  type Entity as EntityProto,
  EntitySchema,
} from '@intrinsic-ai/types/scene/proto/v1/entity_pb';
import type {SceneObject as SceneObjectProto} from '@intrinsic-ai/types/scene/proto/v1/scene_object_pb';
import {KinematicsComponentSchema} from '@intrinsic-ai/types/world/proto/kinematics_component_pb';
import * as three from 'three';

import {
  EntityData,
  FrameEntityData,
  GeometryData,
  JointEntityData,
  LinkEntityData,
  MotionType,
  SensorEntityData,
} from '../../model/entity';
import {
  EntityId,
  NULL_ENTITY_ID,
  NULL_OBJECT_ID,
  ObjectId,
  ROOT_ENTITY_ID,
  ROOT_OBJECT_ID,
  numberToStringIdForEntity,
} from '../../model/world_ids';
import {
  KinematicWorldObjectData,
  WorldObject,
  WorldObjectData,
} from '../../model/world_object';
import {
  EntityDataInterface,
  WorldObjectDataInterface,
} from '../../model/world_object_builder';
import {WorldObjectBuilder} from '../../model/world_object_builder';
import {
  geometryComponentToEntityData,
  kinematicsComponentToEntityData,
  physicsComponentToLinkPhysics,
  sensorComponentToEntityData,
} from '../proto/entity_data_proto';
import {poseProtoToMatrix} from '../proto/pose_proto';

function entityNameToEntityId(
  name: string,
  entityNameToIdMap: ReadonlyMap<string, EntityId>,
): EntityId {
  return entityNameToIdMap.get(name) ?? NULL_ENTITY_ID;
}

function entityNameToEntityStringId(
  name: string,
  entityNameToIdMap: ReadonlyMap<string, EntityId>,
): string {
  const id = entityNameToEntityId(name, entityNameToIdMap);
  return numberToStringIdForEntity(id);
}

interface EntitySource extends EntityDataInterface {
  proto: EntityProto;
  getParentId: () => string;
  sceneObjectId: ObjectId;
  entityNameToIdMap: Map<string, EntityId>;
}

interface WorldObjectSource extends WorldObjectDataInterface<EntitySource> {
  proto: SceneObjectProto;
  objectId: number;
  entityNameToIdMap: Map<string, EntityId>;
  getEntitiesMap: () => {
    get: (id: string) => EntitySource | undefined;
    keys: () => Iterable<string>;
    values: () => Iterable<EntitySource>;
  };
}

function getWorldObjectParamsFromProto(
  source: WorldObjectSource,
): WorldObjectData {
  const proto = source.proto;
  const rootEntity = proto.entities
    .filter((e) => e.parentName === '')
    .map((e) => e.name);
  if (rootEntity.length === 0) {
    throw new Error(
      `Scene object '${proto.name}' does not have a root entity.`,
    );
  }
  if (rootEntity.length > 1) {
    throw new Error(
      `Scene object '${proto.name}' has more than one root entity: ${rootEntity.join(', ')}`,
    );
  }
  const rootEntityId = entityNameToEntityId(
    rootEntity[0],
    source.entityNameToIdMap,
  );

  // NOTE - A SceneObject is converted to a single root WorldObject, so the parent entity id
  // and object id should be null here.
  return {
    id: source.objectId,
    name: proto.name,
    parentEntityId: NULL_ENTITY_ID,
    parentObjectId: NULL_OBJECT_ID,
    parentTThis: new three.Matrix4(),
    isStatic: proto.simulationSpec?.isStatic ?? false,
    isDisabled: proto.simulationSpec?.disabled ?? false,
    resourceName: undefined,
    userDataMap: new Map(Object.entries(proto.userData)),
    rootEntityId,
    children: [], // SceneObject does not have children.
  };
}

function getKinematicWorldObjectParamsFromProto(
  source: WorldObjectSource,
): KinematicWorldObjectData {
  const data = getWorldObjectParamsFromProto(source);
  const proto = source.proto;

  // For SceneObject, we return the list of joints in the order they are defined
  // in the scene object proto. For the purpose of frontend visualization, this
  // should not matter.
  const jointOrdering = proto.entities
    .filter((e) => e.entityType.case === 'joint')
    .map((e) => e.name);

  // Construct the joint positions list using the order from `jointOrdering`.
  const namedConfigurations = new Map<string, number[]>();
  for (const config of proto.properties?.kinematics?.namedConfigurations ??
    []) {
    const configName = config.name;
    if (!configName) {
      console.debug(`Named configuration does not have a name. Skipping.`);
      continue;
    }
    const configJointPositions = new Map(Object.entries(config.jointPositions));
    const jointPositionsList = [];
    for (const jointName of jointOrdering) {
      const position = configJointPositions.get(jointName);
      if (!position) {
        console.debug(
          `Joint '${jointName}' is not found in the named configuration '${configName}'.`,
        );
        break;
      }
      jointPositionsList.push(position);
    }

    // Only add the named configuration if it has the same number of joint
    // positions as the joint ordering.
    if (jointPositionsList.length === jointOrdering.length) {
      namedConfigurations.set(configName, jointPositionsList);
    } else {
      console.debug(
        `Named configuration '${configName}' does not have the same number of joint positions as the joint ordering. Skipping.`,
      );
    }
  }

  const isoFlangeFrames = proto.entities
    .filter((e) => e.entityType.case === 'frame' && e.name === 'flange')
    .map((e) => e.name);

  return {
    ...data,
    namedConfigurations,
    jointOrdering,
    isoFlangeFrames,
    payload: undefined,
  };
}

function getFrameEntityData(source: EntitySource): FrameEntityData {
  const data = getEntityParamsForObjectWorldProto(source);
  return {
    ...data,
    isAttachment:
      (source.proto.entityType.case === 'frame'
        ? source.proto.entityType.value
        : undefined
      )?.isAttachmentFrame ?? false,
  };
}

function getEntityParamsForObjectWorldProto(source: EntitySource): EntityData {
  const parentName = source.proto.parentName;
  const poseProto = source.proto.parentTThis;
  const name = source.proto.name;
  return {
    id: entityNameToEntityId(name, source.entityNameToIdMap),
    name,
    parentObjectId: source.sceneObjectId,
    parentTThis: poseProto ? poseProtoToMatrix(poseProto) : new three.Matrix4(),
    parentId: parentName
      ? entityNameToEntityId(parentName, source.entityNameToIdMap)
      : NULL_ENTITY_ID,
    resourceName: undefined,
  };
}

function getJointEntityData(source: EntitySource): JointEntityData {
  const data = getEntityParamsForObjectWorldProto(source);
  const entityProto = source.proto;
  const jointProto =
    entityProto.entityType.case === 'joint'
      ? entityProto.entityType.value
      : undefined;
  if (!jointProto) {
    throw new Error(
      `Cannot get entity data for non-joint entity ${source.proto.name}`,
    );
  }
  const kinematics = jointProto.kinematicsComponent;
  const entityData = kinematicsComponentToEntityData(
    data,
    kinematics ?? create(KinematicsComponentSchema),
  );

  // Joint entity is required to have a kinematics component. Since this does
  // not, let's return sane defaults with UNDEFINED motion type so that client
  // code can handle this as needed.
  if (!kinematics) {
    entityData.motionType = MotionType.UNDEFINED;
  }

  return entityData;
}

function getLinkEntityData(source: EntitySource): LinkEntityData {
  const data = getEntityParamsForObjectWorldProto(source);
  const linkProto =
    source.proto.entityType.case === 'link'
      ? source.proto.entityType.value
      : undefined;
  if (!linkProto) {
    throw new Error(
      `Cannot get entity data for non-link entity ${source.proto.name}`,
    );
  }
  const physics = physicsComponentToLinkPhysics(linkProto.physicsComponent);
  const geometryComponent = linkProto.geometryComponent;
  if (geometryComponent) {
    return {
      ...geometryComponentToEntityData(data, geometryComponent),
      physics,
    };
  }

  return {
    ...data,
    geometries: new Map<string, Map<string, GeometryData>>(),
    physics,
  };
}

function getSensorEntityData(source: EntitySource): SensorEntityData {
  const data = getEntityParamsForObjectWorldProto(source);
  const sensorProto =
    source.proto.entityType.case === 'sensor'
      ? source.proto.entityType.value
      : undefined;
  if (!sensorProto) {
    throw new Error(
      `Cannot get entity data for non-sensor entity ${source.proto.name}`,
    );
  }

  const sensorComponent = sensorProto.sensorComponent;
  if (sensorComponent) {
    return sensorComponentToEntityData(data, sensorComponent);
  }

  return {
    ...data,
    sensor: undefined,
  };
}

function getWorldObjectParamsForProto(
  source: WorldObjectSource,
): WorldObjectData {
  const isKinematic =
    source.proto.entities.filter((e) => e.entityType.case === 'joint').length >
    0;
  if (isKinematic) {
    return getKinematicWorldObjectParamsFromProto(source);
  }
  return getWorldObjectParamsFromProto(source);
}

/**
 * Builds a world object from the given SceneObject proto. The SceneObject proto
 * will be converted to a WorldObject tree with a single, root WorldObject. Each
 * Entity within the SceneObject will be converted to a corresponding Entity in
 * the WorldObject.
 *
 * @param proto The scene object proto.
 */
export function sceneObjectToWorldObject(
  proto: SceneObjectProto,
): WorldObject | undefined {
  if (proto.entities.length === 0) {
    throw new Error(`Scene object '${proto.name}' is blank (has no entities).`);
  }

  const worldObjectDataProvider = {
    getData: (source: WorldObjectSource) =>
      getWorldObjectParamsForProto(source),
  };
  const entityDataProvider = {
    getData: (source: EntitySource) => {
      const entityProto = source.proto;
      const entityTypeCase = entityProto.entityType.case;
      switch (entityTypeCase) {
        case 'frame':
          return getFrameEntityData(source);
        case 'joint':
          return getJointEntityData(source);
        case 'link':
          return getLinkEntityData(source);
        case 'sensor':
          return getSensorEntityData(source);
        case undefined:
          break;
        default:
          assumeExhaustive(entityTypeCase);
          break;
      }

      console.debug(
        `Entity '${entityProto.name}' is not a frame, joint, link, or sensor. Defaulting to a frame`,
      );
      return {
        ...getEntityParamsForObjectWorldProto(source),
        isAttachment: false,
      };
    },
  };
  // Generate a mapping from entity name to entity id.
  const entityNameToIdMap = new Map<string, number>();
  const maxEntitiesPerObject = 100000;
  if (maxEntitiesPerObject < proto.entities.length) {
    throw new Error(
      `Scene object '${proto.name}' has more entities than the maximum ${maxEntitiesPerObject} supported.`,
    );
  }

  const rootObjectId = ROOT_OBJECT_ID;
  const usedEntityIds = new Set<number>([NULL_ENTITY_ID, ROOT_ENTITY_ID]);
  proto.entities.forEach((e) => {
    if (e.parentName === '') {
      entityNameToIdMap.set(e.name, ROOT_ENTITY_ID);
    } else {
      // We use a hash of the entity name to keep the IDs stable when entities
      // are added, removed, or reparented. Note that an entity's id will still
      // change if it is renamed. It's safe to use the entity name here because
      // they are required to be unique within a SceneObject. Collisions should
      // not realistically happen (cyrb53 has 53-bit output), but we check just in case.
      const id = cyrb53(e.name);
      if (usedEntityIds.has(id)) {
        throw new Error(
          `Entity ID collision detected for entity '${e.name}' with ID ${id}`,
        );
      }
      usedEntityIds.add(id);
      entityNameToIdMap.set(e.name, id);
    }
  });

  const objMap = new WorldObjectBuilder<EntitySource, WorldObjectSource>(
    worldObjectDataProvider,
    entityDataProvider,
  ).build([
    {
      proto,
      objectId: rootObjectId,
      entityNameToIdMap,
      getEntitiesMap: (): Map<string, EntitySource> => {
        const entityAdapterMap = new Map<string, EntitySource>();
        const sceneObjectId = rootObjectId;
        for (const entity of proto.entities) {
          const entityStringId: string = entityNameToEntityStringId(
            entity.name,
            entityNameToIdMap,
          );
          entityAdapterMap.set(entityStringId, {
            proto: clone(EntitySchema, entity),
            getParentId: () =>
              entity.parentName
                ? entityNameToEntityStringId(
                    entity.parentName,
                    entityNameToIdMap,
                  )
                : '',
            sceneObjectId,
            entityNameToIdMap,
          });
        }
        return entityAdapterMap;
      },
    },
  ]);

  return objMap.values().next().value;
}
