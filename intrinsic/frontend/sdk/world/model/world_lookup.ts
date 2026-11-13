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

import * as three from 'three';

import {CameraSensor, Entity, FrameEntity, isCameraSensor} from './entity';
import {
  EntityId,
  FrameId,
  NULL_ENTITY_ID,
  NULL_OBJECT_ID,
  ObjectId,
  ObjectOrFrameId,
  ROOT_ENTITY_ID,
} from './world_ids';
import {WorldObject} from './world_object';

/**
 * WorldLookup is a helper class that, given a root WorldObject, constructs
 * lookup maps and implements helpful lookup functions.
 *
 * Changing the structure of the WorldObject tree (such as adding/removing entities) will
 * invalidate this object. You must construct a new `WorldLookup` in these cases.
 *
 * TODO(b/507855047) - There are a lot of redundant lookup fields here, like everything
 * related to Frames. Should be removed.
 */
export class WorldLookup {
  // Objects in the world from the ObjectWorld representation.
  // Keyed by unique Object.id
  readonly objects = new Map<ObjectId, WorldObject>();
  // Map from object/frame id to names.
  readonly objFrameNames = new Map<ObjectOrFrameId, string>();
  // Frames (frame entities) of all objects in the world, keyed by id.
  readonly frames = new Map<FrameId, FrameEntity>();
  // Map storing the set of frames under each object.
  readonly objFrames = new Map<ObjectId, Set<FrameId>>();

  // Map storing the set of entities within each object.
  readonly objEntities = new Map<ObjectId, Set<EntityId>>();
  // Map storing the root entity id of each object.
  readonly objRootEntity = new Map<ObjectId, EntityId>();
  // Map storing the object id for each entity id.
  readonly entityObjIds = new Map<EntityId, ObjectId>();

  // List of resource names associated with objects.
  readonly resourceNames = new Map<ObjectId, string>();

  // Entities in the world keyed by the entity id.
  readonly entities = new Map<EntityId, Entity>();
  readonly entityNames = new Map<EntityId, string>();
  readonly entityParentIds = new Map<EntityId, EntityId>();
  readonly entityChildIds = new Map<EntityId, Set<EntityId>>();
  // Local transform of each object, stored here for restoration operation.
  readonly entityInitialTrans = new Map<EntityId, three.Matrix4>();

  constructor(rootObject: WorldObject) {
    // Populate everything else by traversing object structure.
    const objectQueue = [rootObject];
    while (objectQueue.length > 0) {
      const obj = objectQueue.shift();
      if (!obj) {
        break;
      }
      const oid = obj.id;
      this.objects.set(oid, obj);

      // Set frames
      this.objFrames.set(oid, new Set([...obj.getFrames()].map((f) => f.id)));
      this.objFrameNames.set(oid, obj.name);
      for (const f of obj.getFrames()) {
        this.frames.set(f.id, f);
        this.entityObjIds.set(f.id, oid);
        this.objFrameNames.set(f.id, f.name);
      }

      // Set entities
      for (const e of obj.getEntities()) {
        this.entityObjIds.set(e.id, oid);
        this.entities.set(e.id, e);
      }

      this.objEntities.set(
        oid,
        new Set([...obj.getEntities()].map((e) => e.id)),
      );
      this.objRootEntity.set(oid, obj.getRootEntity().id);

      // Set resource names
      if (obj.isResource()) {
        this.resourceNames.set(oid, obj.getResourceName());
      }

      objectQueue.push(...obj.getChildren());
    }

    // Extract entity info and build the child/parent map of entities.
    for (const [eid, entity] of this.entities.entries()) {
      this.processEntity(eid, entity);
    }
  }

  private updateChildParentMap(eid: EntityId, pid: EntityId) {
    const oldPid = this.entityParentIds.get(eid);
    if (oldPid) {
      this.entityChildIds.get(oldPid)?.delete(eid);
    }
    let childIds = this.entityChildIds.get(pid);
    if (!childIds) {
      childIds = new Set<EntityId>();
      this.entityChildIds.set(pid, childIds);
    }
    childIds.add(eid);
    this.entityParentIds.set(eid, pid);
  }

  private processEntity(id: EntityId, entity: Entity) {
    const eid = Number(id);
    this.entityNames.set(eid, entity.name);
    this.entityInitialTrans.set(eid, new three.Matrix4());

    if (eid === ROOT_ENTITY_ID) {
      this.updateChildParentMap(eid, NULL_ENTITY_ID);
    } else {
      const parent = entity.getParent();
      const pid = parent ? parent.id : NULL_ENTITY_ID;
      this.updateChildParentMap(eid, pid);
      this.entityInitialTrans.set(eid, entity.parentTThis);
    }
  }

  /**
   * Traverse the model tree from a given entity node and return the list of
   * entity ids encountered during traversal. There is no traverse order for
   * siblings as they are not ordered originally.
   */
  traverseEntityChildren(eid: EntityId) {
    let res = [eid];
    const cids = this.entityChildIds.get(eid) ?? [];
    for (const cid of cids) {
      res.push(cid);
      res = res.concat(this.traverseEntityChildren(cid));
    }
    return res;
  }

  isEntityFrame(id: EntityId) {
    return this.frames.has(id);
  }

  /**
   * Returns true if the entity have any shape.
   */
  hasShape(eid: EntityId): boolean {
    const entity = this.entities.get(eid);
    if (!entity?.isLink()) return false;
    return (
      entity.getCollisionGeometries().length > 0 ||
      entity.getVisualGeometries().length > 0
    );
  }

  /**
   * @returns Whether entity eid1 is ancestor of eid2 in this model.
   *
   * NOTE: eid1 has to be ancestor of eid2, if eid1 equals eid2 this returns
   * false.
   */
  isAncestorEntity(eid1: EntityId, eid2: EntityId): boolean {
    if (!this.entities.has(eid1) || !this.entities.has(eid2)) {
      return false;
    }
    while ((eid2 = this.entityParentIds.get(eid2) ?? NULL_ENTITY_ID)) {
      if (eid1 === eid2) {
        return true;
      }
    }
    return false;
  }

  isObjectResource(oid: ObjectId): boolean {
    const obj = this.objects.get(oid);
    if (obj) return obj.isResource();

    return !!this.entities.get(oid)?.isResource();
  }

  // Returns the first camera sensor of the object.
  getObjectCamera(oid: ObjectId): CameraSensor | undefined {
    const eids = this.objEntities.get(oid);
    if (!eids) return undefined;
    for (const eid of eids.values()) {
      const sensor = this.entities.get(eid)?.asSensor()?.sensor;
      if (sensor && isCameraSensor(sensor)) {
        return sensor;
      }
    }
    return undefined;
  }

  // Returns true if any entity under the object have camera sensor component.
  // Note entities of child objects does not count, this only checks entities
  // that compose the object.
  isObjectCamera(oid: ObjectId): boolean {
    return !!this.getObjectCamera(oid);
  }

  // World objects have a parent object, but also specify the specific entity
  // they are parented to. This method returns the specific entity id, if one is
  // specified, or undefined.
  getObjectParentEntityId(oid: ObjectId): EntityId | undefined {
    return this.objects.get(oid)?.getRootEntity().getParent()?.id;
  }

  // Convert an input object/frame id to its corresponding entity id.
  toEntityId(id: ObjectOrFrameId): EntityId {
    if (this.isEntityFrame(id)) {
      return id;
    }
    const rootEntityId = this.objRootEntity.get(id);
    if (rootEntityId !== undefined) {
      return rootEntityId;
    }
    if (this.entities.has(id)) {
      return id;
    }
    // Input is unknown.
    return NULL_ENTITY_ID;
  }

  /**
   * Convert an input entity id to its corresponding object/frame id.
   */
  getEntityObjectId(eid: EntityId): ObjectId {
    if (this.objects.has(eid)) return eid;
    return this.entityObjIds.get(eid) ?? NULL_OBJECT_ID;
  }

  /**
   * Compute and return the world transform of the entity based on proto.
   * NOTE this is a potentially expensive operation.
   */
  getWorldTrans(eid?: EntityId): three.Matrix4 | undefined {
    const trans = new three.Matrix4();
    while (eid && eid !== NULL_ENTITY_ID) {
      const entity = this.entities.get(eid);
      if (!entity) return undefined;
      trans.premultiply(entity.parentTThis);
      eid = entity.getParent()?.id;
    }
    return trans;
  }
}
