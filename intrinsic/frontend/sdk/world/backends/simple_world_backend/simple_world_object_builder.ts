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

import {ROOT_OBJECT_ID, numberToStringIdForEntity} from '../../model/world_ids';
import {
  KinematicWorldObjectData,
  WorldObject,
  WorldObjectData,
} from '../../model/world_object';
import {
  EntityDataInterface,
  TypedEntityData,
  WorldObjectBuilder,
  WorldObjectDataInterface,
} from '../../model/world_object_builder';

interface SimpleEntitySource extends EntityDataInterface {
  data: TypedEntityData;
}

interface SimpleWorldObjectSource extends WorldObjectDataInterface<SimpleEntitySource> {
  data: WorldObjectData | KinematicWorldObjectData;
}

/**
 * Builds a tree of WorldObjects from object and entity data and returns the root WorldObject.
 *
 * Usage example:
 * ```ts
 * const rootObjectData: WorldObjectData = {
 *   id: ROOT_OBJECT_ID,
 *   name: 'root',
 *   parentEntityId: NULL_ENTITY_ID,
 *   rootEntityId: ROOT_ENTITY_ID,
 *   parentObjectId: NULL_OBJECT_ID,
 *   parentTThis: new three.Matrix4(),
 *   children: [],
 * };
 *
 * const rootEntityData: LinkEntityData = {
 *   id: ROOT_ENTITY_ID,
 *   name: 'Root',
 *   parentObjectId: ROOT_OBJECT_ID,
 *   parentId: NULL_ENTITY_ID,
 *   parentTThis: new three.Matrix4(),
 *   geometries: new Map(),
 * };
 *
 * const rootWorldObject = buildWorldObjectTree([rootObjectData], [rootEntityData]);
 * ```
 */
export function buildWorldObjectTree(
  objects: Array<WorldObjectData | KinematicWorldObjectData>,
  entities: TypedEntityData[],
): WorldObject {
  const objectSources: SimpleWorldObjectSource[] = objects.map((obj) => {
    const objectEntities = entities.filter((e) => e.parentObjectId === obj.id);
    const entityMap = new Map<string, SimpleEntitySource>();
    for (const entity of objectEntities) {
      const key = numberToStringIdForEntity(entity.id);
      entityMap.set(key, {
        data: entity,
        getParentId: () => numberToStringIdForEntity(entity.parentId),
      });
    }

    return {
      data: obj,
      getEntitiesMap: () => ({
        get: (id: string) => entityMap.get(id),
        keys: () => entityMap.keys(),
        values: () => entityMap.values(),
      }),
    };
  });

  const objectDataProvider = {
    getData: (source: SimpleWorldObjectSource) => source.data,
  };

  const entityDataProvider = {
    getData: (source: SimpleEntitySource) => source.data,
  };

  const objMap = new WorldObjectBuilder<
    SimpleEntitySource,
    SimpleWorldObjectSource
  >(objectDataProvider, entityDataProvider).build(objectSources);

  const rootObj = objMap.get(ROOT_OBJECT_ID);
  if (!rootObj) {
    throw new Error(
      `Could not find root object with ID ${ROOT_OBJECT_ID} when building world tree`,
    );
  }
  return rootObj;
}
