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
  Entity,
  FrameEntity,
  FrameEntityData,
  JointEntity,
  JointEntityData,
  LinkEntity,
  LinkEntityData,
  MotionType,
  SensorEntity,
  SensorEntityData,
} from './entity';
import {
  EntityId,
  NULL_OBJECT_ID,
  ObjectId,
  numberToStringIdForEntity,
  stringToNumberId,
} from './world_ids';
import {
  KinematicWorldObject,
  KinematicWorldObjectData,
  WorldObject,
  WorldObjectData,
} from './world_object';

declare interface DataProvider<DataSource, DataType> {
  getData: (source: DataSource) => DataType;
}

/**
 * WorldObjectDataInterface is the interface to instantiate a WorldObject.
 */
export declare interface WorldObjectDataInterface<EntityDataType> {
  getEntitiesMap: () => {
    get: (id: string) => EntityDataType | undefined;
    keys: () => Iterable<string>;
    values: () => Iterable<EntityDataType>;
  };
}

/**
 * EntityDataInterface is the interface to instantiate an Entity within the
 * WorldObject. This includes frames, joints, links, and sensors.
 */
export declare interface EntityDataInterface {
  getParentId: () => string;
}

function isKinematicData(
  data: WorldObjectData | KinematicWorldObjectData,
): data is KinematicWorldObjectData {
  return 'jointOrdering' in data && data.jointOrdering.length > 0;
}

/**
 * Union type representing all specialized entity data types within a WorldObject.
 */
export type TypedEntityData =
  JointEntityData | LinkEntityData | FrameEntityData | SensorEntityData;

function isJointEntityData(data: TypedEntityData): data is JointEntityData {
  return 'motionType' in data && data.motionType !== MotionType.UNDEFINED;
}

function isLinkEntityData(data: TypedEntityData): data is LinkEntityData {
  return 'geometries' in data;
}

function isSensorEntityData(data: TypedEntityData): data is SensorEntityData {
  return 'sensor' in data;
}

/**
 * A WorldObjectBuilder is a class that constructs a world model from a
 * "flattened" data source.
 *
 * There are two primitives in generating world objects:
 *   - Objects
 *   - Entities (links, joints, frames and sensors)
 *
 * For each primitive, it is expected that we will need to specify both: a data
 * source, which is the specific data instance of the flattened source, and a
 * data provider, which will convert the data source into one of the
 * primitive data types (e.g. WorldObjectData, LinkEntityData, etc).
 *
 * One example when using the `proto_world_object_builder`, is that the data
 * source for world objects is the `Object` jspb proto from the world service,
 * and the data provider is an inline function that converts this proto to a
 * `WorldObjectData` object as defined in `world_object.ts`. Please see
 * `proto_world_object_builder.ts` for more examples.
 **/
export class WorldObjectBuilder<
  E extends EntityDataInterface,
  O extends WorldObjectDataInterface<E>,
> {
  constructor(
    private readonly objectDataProvider: DataProvider<
      O,
      WorldObjectData | KinematicWorldObjectData
    >,
    private readonly entityDataProvider: DataProvider<E, TypedEntityData>,
  ) {}

  /**
   * Given a list of object sources, uses the information within them to build the
   * tree of objects, entities, and frames that contain them. It is expected that
   * there is only one root entity among the members of the list. Throws an error
   * if the world structure is ill-formed.
   */
  build(objectSources: readonly O[]): Map<ObjectId, WorldObject> {
    const result = new Map<ObjectId, WorldObject>();

    // Find all objects and their children
    const objectSourceMap = new Map<ObjectId, O>();
    const numUnprocessedChildren = new Map<ObjectId, number>();
    for (const obj of objectSources) {
      const objData = this.objectDataProvider.getData(obj);
      const oid = objData.id;
      objectSourceMap.set(oid, obj);
      numUnprocessedChildren.set(oid, objData.children.length);
    }

    const childrenToProcess = [...numUnprocessedChildren.keys()].filter(
      (oid) => numUnprocessedChildren.get(oid) === 0,
    );

    // Go bottom up, for each object without children, and then remove
    // it from the list to process.
    while (childrenToProcess.length > 0) {
      const nextObjectId = childrenToProcess.shift();
      if (nextObjectId === NULL_OBJECT_ID) continue;
      if (!nextObjectId) {
        throw new Error(
          'Internal error: Could not find childless object when building world tree!',
        );
      }

      // Collect all children
      const nextObjectChildren = [...result.values()].filter(
        (wo) => wo.getParentObjectId() === nextObjectId,
      );

      const nextObjectSource = objectSourceMap.get(nextObjectId);
      if (!nextObjectSource) {
        throw new Error(
          `Internal error: Could not find object source for id ${nextObjectId}`,
        );
      }
      const nextObject = this.createObjectFromWorldObjectSource(
        nextObjectSource,
        nextObjectChildren,
      );
      result.set(nextObjectId, nextObject);

      // Remove object from parent list of children to process.
      const nextObjectParentId = nextObject.getParentObjectId();
      if (nextObjectParentId !== undefined) {
        const nChildren = numUnprocessedChildren.get(nextObjectParentId);
        if (nChildren) {
          if (nChildren === 1) {
            childrenToProcess.push(nextObjectParentId);
          }
          numUnprocessedChildren.set(nextObjectParentId, nChildren - 1);
        }
      }
    }

    return result;
  }

  private createObjectFromWorldObjectSource(
    worldObjectSource: O,
    childObjects: WorldObject[],
  ): WorldObject {
    const data = this.objectDataProvider.getData(worldObjectSource);
    const newObj = (() => {
      if (isKinematicData(data)) {
        return new KinematicWorldObject(data, childObjects);
      }
      return new WorldObject(data, childObjects);
    })();
    this.processChildrenForWorldObjectData(
      newObj,
      worldObjectSource,
      childObjects,
    );
    if (isKinematicData(data) && newObj instanceof KinematicWorldObject) {
      this.processKinematicsForObject(newObj, data);
    }
    return newObj;
  }

  /**
   * Creates a new Entity for the given EntityProto.
   *
   * @param source The data source for this entity.
   * @param parentObject The parent object for this entity.
   *     CAUTION: might not be fully initialized -- best to avoid using it here.
   * @param children The children for this entity.
   */
  private createEntityFromWorldObjectData(
    source: E,
    parentObject: WorldObject,
    children: Entity[],
  ): Entity {
    const data = this.entityDataProvider.getData(source);
    if (isJointEntityData(data)) {
      return new JointEntity(data, parentObject, children);
    } else if (isLinkEntityData(data)) {
      return new LinkEntity(data, parentObject, children);
    } else if (isSensorEntityData(data)) {
      return new SensorEntity(data, parentObject, children);
    } else {
      // data instanceof FrameEntityData
      return new FrameEntity(data, parentObject, children);
    }
  }

  private processKinematicsForObject(
    obj: KinematicWorldObject,
    data: KinematicWorldObjectData,
  ) {
    // Set joint references. To do that, ensure that the joints array is empty
    // to start.
    const objJoints = obj.getJointsInternal();
    if (!objJoints) return;

    objJoints.length = 0;
    for (const jointName of data.jointOrdering) {
      const jointEntity = [...obj.getEntities()].find(
        (j) => j.name === jointName,
      );
      if (!jointEntity?.isJoint()) {
        throw new Error(
          `Internal error: kinematics data specified joint entity with name '${jointName}', but no such joint was found.`,
        );
      }

      objJoints.push(jointEntity);
    }

    // Set flange frame references
    const objFlangeFrames = obj.getFlangeFramesInternal();
    if (!objFlangeFrames) return;

    for (const frameName of data.isoFlangeFrames) {
      const frame = obj.getFrameByName(frameName);
      if (!frame) {
        throw new Error(
          `Internal error: kinematics data specified iso frame with name '${frameName}', but no such frame was found.`,
        );
      }

      objFlangeFrames.push(frame);
    }

    const objLinks = obj.getLinksInternal();
    if (objLinks) {
      // Breadth-first traversal over robot object entities to collect links in
      // order.
      objLinks.length = 0;
      const rootEntity = obj.getRootEntity();
      const entities: Entity[] = rootEntity ? [rootEntity] : [];
      while (entities.length > 0) {
        const entity = entities.shift();
        if (!entity || entity.parentObject.id !== obj.id) continue;

        if (entity.isLink()) {
          objLinks.push(entity);
        }

        entities.push(...entity.getChildren());
      }
    }
  }

  private processChildrenForWorldObjectData(
    obj: WorldObject,
    source: O,
    childObjects: WorldObject[],
  ) {
    // Handle entities. We need to find all entities that either have no child entities, or the
    // children are in the child objects, and then work our way towards the
    // root entity.

    // The entities to process are those that have no current dependencies on
    // unprocessed entities. So, we keep a map of entities in this and how many
    // children they have.
    const numUnprocessedChildren = new Map<EntityId, number>();
    for (const sid of source.getEntitiesMap().keys()) {
      const eid = stringToNumberId(sid);
      numUnprocessedChildren.set(eid, 0);
    }
    for (const e of source.getEntitiesMap().values()) {
      const pid = stringToNumberId(e.getParentId());
      const nChildren = numUnprocessedChildren.get(pid);
      if (nChildren !== undefined) {
        numUnprocessedChildren.set(pid, nChildren + 1);
      }
    }
    const toProcessEntities = [...numUnprocessedChildren.keys()].filter(
      (eid) => numUnprocessedChildren.get(eid) === 0,
    );

    // Build initial map of child entities: all root entities for the child objects.
    const possibleChildEntities = new Map<EntityId, Entity>();
    for (const child of obj.getChildren()) {
      const childRoot = child.getRootEntity();
      possibleChildEntities.set(childRoot.id, childRoot);
    }

    // Process all of the entities here
    while (toProcessEntities.length > 0) {
      const nextEntityId = toProcessEntities.shift();
      if (nextEntityId === undefined) break;

      // Find all of the children of this entity.
      const childEntities = [...possibleChildEntities.values()].filter(
        (child) => child.getParentId() === nextEntityId,
      );

      // Process this entity.
      const entitySource = source
        .getEntitiesMap()
        .get(numberToStringIdForEntity(nextEntityId));
      if (!entitySource) {
        throw new Error(`Entity proto for entity ${nextEntityId} not found`);
      }

      const nextEntity = this.createEntityFromWorldObjectData(
        entitySource,
        obj,
        childEntities,
      );

      // Update all of our tracking.
      obj.getEntitiesInternal().set(nextEntityId, nextEntity);
      const nextEntityParentId = nextEntity.getParentId();
      if (nextEntityParentId === undefined) {
        // Processing root entity -- this should only happen when we've finished
        // all of the other entities.
        if (toProcessEntities.length > 0) {
          throw new Error(
            `Internal error: Reached root entity (${nextEntityId}) for object ${obj.id}, but not all children were processed.`,
          );
        }
        break;
      }

      const nChildren = numUnprocessedChildren.get(nextEntityParentId);
      if (nChildren) {
        // If this was the last child, mark the parent to process.
        if (nChildren === 1) {
          toProcessEntities.push(nextEntityParentId);
        }
        numUnprocessedChildren.set(nextEntityParentId, nChildren - 1);
      }
      possibleChildEntities.set(nextEntityId, nextEntity);
    }

    const rootEntityId = this.objectDataProvider.getData(source).rootEntityId;
    const rootEntity = obj.getEntitiesInternal().get(rootEntityId);
    if (!rootEntity) {
      throw new Error(
        `Internal error: Root entity (${rootEntityId}) for object ${obj.id} not found`,
      );
    }
    // See the requirements documented on `WorldObject`.
    if (!(rootEntity instanceof LinkEntity)) {
      throw new Error(
        `The root entity '${rootEntity.name}' (${rootEntityId}) of object '${obj.name}' (${obj.id}) ` +
          `is a ${rootEntity.getEntityType()}, but the root entity of an object must be a link.`,
      );
    }
    obj.setRootEntityInternal(rootEntity);
  }
}
