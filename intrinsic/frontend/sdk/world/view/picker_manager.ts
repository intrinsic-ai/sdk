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

import {BehaviorSubject, Observable, Subject} from 'rxjs';

import {NULL_ENTITY_ID, NULL_OBJECT_ID} from '../model/world_ids';

/**
 * A custom user-defined selection type. Use this if you want to allow the user
 * to pick things other than objects/entities in your app and have it interact
 * properly with the rest of the picker state. For example, you may want to allow
 * the user to pick asset types, an item from a list, etc.
 *
 * The PickerManager handles keeping this selection mutually exclusive with picking
 * regular objects and entities. Prefer using this rather than extending PickerManager to support
 * custom selection types.
 *
 * Usage:
 *
 * ```ts
 * // Define your selection type(s), each with a unique `kind`.
 * class ServiceSelection implements CustomSelection {
 *   readonly kind = 'service';
 *
 *   constructor(readonly id: string) {}
 * }
 *
 * // Set the selection. This clears any picked objects.
 * pickerManager.setCustomSelection(new ServiceSelection('my-service'));
 *
 * // Check the current selection, or subscribe to changes. Picking an object
 * // clears the selection, which then emits undefined.
 * const current = pickerManager.getCustomSelection();
 * if (current instanceof ServiceSelection) {
 *   console.log(current.id);
 * }
 * pickerManager.newCustomSelection.subscribe((selection) => {
 *   const serviceId =
 *     selection instanceof ServiceSelection ? selection.id : undefined;
 * });
 *
 * // Clear the selection (together with any picked objects and entities).
 * pickerManager.clearAllPicks();
 * ```
 */
export interface CustomSelection {
  readonly kind: string;
}

/** Captures a change to the set of picked objects */
interface PickerState {
  enabled: boolean;
  pickedObjectIds: ReadonlySet<number>;
  pickedEntityId: number;
  hoverPickedEntityId: number;
  hoverPickedObjectId: number;
  deepPickedObjectId: number;
  customSelection?: CustomSelection;
}

/**
 * Deep picking is a mode where the user 'focuses' on certain objects
 * for closer inspection and editing of its entities. When in deep pick mode,
 * we disallow picking objects and only allow picking entities. Outside of
 * deep pick mode, we only allow picking objects and do not allow picking entities.
 * Custom selections may be made whether in deep pick mode or not.
 */
function applyPickMode(state: PickerState): PickerState {
  if (state.deepPickedObjectId !== NULL_OBJECT_ID) {
    // We do not allow picking objects in deep pick mode, so clear these.
    return {
      ...state,
      pickedObjectIds: new Set<number>(),
      hoverPickedObjectId: NULL_OBJECT_ID,
    };
  }
  // We do not allow picking entities outside of deep pick mode, so clear these.
  return {
    ...state,
    pickedEntityId: NULL_ENTITY_ID,
    hoverPickedEntityId: NULL_ENTITY_ID,
  };
}

/** The kinds of selection that are mutually exclusive. */
type SelectionKind = 'objects' | 'entity' | 'custom';

/**
 * Returns the kind of selection that the update picked, if any. Clearing a
 * pick (an empty set or null id) doesn't count as picking.
 */
function pickedSelectionKind(
  state: PickerState,
  update: Partial<PickerState>,
): SelectionKind | undefined {
  if (update.pickedObjectIds !== undefined && state.pickedObjectIds.size > 0) {
    return 'objects';
  }
  if (
    update.pickedEntityId !== undefined &&
    state.pickedEntityId !== NULL_ENTITY_ID
  ) {
    return 'entity';
  }
  if (update.customSelection !== undefined) {
    return 'custom';
  }
  return undefined;
}

/**
 * Object picks, entity picks and the custom selection are mutually exclusive:
 * whichever was picked most recently wins, and the others are cleared.
 */
function applySelectionExclusivity(
  state: PickerState,
  update: Partial<PickerState>,
): PickerState {
  const picked = pickedSelectionKind(state, update);
  if (picked === undefined) return state;
  return {
    ...state,
    pickedObjectIds:
      picked === 'objects' ? state.pickedObjectIds : new Set<number>(),
    pickedEntityId: picked === 'entity' ? state.pickedEntityId : NULL_ENTITY_ID,
    customSelection: picked === 'custom' ? state.customSelection : undefined,
  };
}

function setsEqual<T>(first: ReadonlySet<T>, second: ReadonlySet<T>): boolean {
  if (first.size !== second.size) return false;
  for (const item of first) {
    if (!second.has(item)) return false;
  }
  return true;
}

function emitIfChanged<T>(subject: Subject<T>, oldValue: T, newValue: T) {
  if (oldValue !== newValue) {
    subject.next(newValue);
  }
}

/**
 * The PickerManager centralizes selection state, selection subscriptions, and
 * selection logic for the FlowState application. The global picker uses the
 * pickerManager to apply changes to the selection lists and to notify
 * subscribers of selection changes.
 * This manager handles picking objects, entities (while an object is deep
 * picked) and custom selections (see `CustomSelection`). The three are mutually
 * exclusive. The logic resides here to ensure they cannot be picked at the same
 * time.
 *
 * Deep picking:
 * Certain components like the object instance editor and type editor
 * set the "deep-picked" object. Deep-picking loosely means that we want to
 * select/focus a specific object for more fine-grained editing (usually at the Entity level).
 * Different components may listen to changes to the deep-picked object and enable
 * special behaviors accordingly.
 */
export class PickerManager {
  private state: PickerState = {
    // Whether or not picking is enabled.
    enabled: true,
    pickedObjectIds: new Set<number>(),
    pickedEntityId: NULL_ENTITY_ID,
    hoverPickedEntityId: NULL_ENTITY_ID,
    hoverPickedObjectId: NULL_OBJECT_ID,
    deepPickedObjectId: NULL_OBJECT_ID,
    customSelection: undefined,
  };

  private readonly newPickSubject = new BehaviorSubject<ReadonlySet<number>>(
    this.state.pickedObjectIds,
  );
  private readonly newHoverPickSubject = new BehaviorSubject<number>(
    this.state.hoverPickedObjectId,
  );
  private readonly deepPickSubject = new BehaviorSubject<number>(
    this.state.deepPickedObjectId,
  );
  private readonly newEntityPickSubject = new BehaviorSubject<number>(
    this.state.pickedEntityId,
  );
  private readonly newHoverEntityPickSubject = new BehaviorSubject<number>(
    this.state.hoverPickedEntityId,
  );
  private readonly newCustomSelectionSubject = new BehaviorSubject<
    CustomSelection | undefined
  >(this.state.customSelection);

  // Observables that notify subscribers of distinct picking events.
  // Note newPick is only fired when only a single object is picked.
  readonly newPick: Observable<ReadonlySet<number>> = this.newPickSubject;
  readonly newHoverPick: Observable<number> = this.newHoverPickSubject;
  readonly newEntityPick: Observable<number> = this.newEntityPickSubject;
  readonly newHoverEntityPick: Observable<number> =
    this.newHoverEntityPickSubject;
  readonly deepPick: Observable<number> = this.deepPickSubject;
  readonly newCustomSelection: Observable<CustomSelection | undefined> =
    this.newCustomSelectionSubject;

  // If you were thinking of adding a constructor parameter here, you should
  // _probably_ reconsider.

  isEnabled(): boolean {
    return this.state.enabled;
  }

  getPickedObjectIds(): ReadonlySet<number> {
    return this.state.pickedObjectIds;
  }

  getPickedEntityId(): number {
    return this.state.pickedEntityId;
  }

  getHoverPickedEntityId(): number {
    return this.state.hoverPickedEntityId;
  }

  getHoverPickedObjectId(): number {
    return this.state.hoverPickedObjectId;
  }

  getDeepPickedObjectId(): number {
    return this.state.deepPickedObjectId;
  }

  getCustomSelection(): CustomSelection | undefined {
    return this.state.customSelection;
  }

  private updateState(update: Partial<PickerState>) {
    const oldState = this.state;
    const mergedState: PickerState = {...oldState, ...update};
    // If the picker is disabled, only allow changes to enabled state.
    if (!mergedState.enabled) {
      this.state = {...oldState, enabled: false};
      return;
    }

    const newState = applySelectionExclusivity(
      applyPickMode(mergedState),
      update,
    );

    this.state = newState;
    this.notifyChanges(oldState, newState);
  }

  /** Notifies the subscribers of each part of the state that changed. */
  private notifyChanges(oldState: PickerState, newState: PickerState) {
    if (!setsEqual(oldState.pickedObjectIds, newState.pickedObjectIds)) {
      this.newPickSubject.next(newState.pickedObjectIds);
    }
    emitIfChanged(
      this.newHoverPickSubject,
      oldState.hoverPickedObjectId,
      newState.hoverPickedObjectId,
    );
    emitIfChanged(
      this.newEntityPickSubject,
      oldState.pickedEntityId,
      newState.pickedEntityId,
    );
    emitIfChanged(
      this.newHoverEntityPickSubject,
      oldState.hoverPickedEntityId,
      newState.hoverPickedEntityId,
    );
    emitIfChanged(
      this.deepPickSubject,
      oldState.deepPickedObjectId,
      newState.deepPickedObjectId,
    );
    emitIfChanged(
      this.newCustomSelectionSubject,
      oldState.customSelection,
      newState.customSelection,
    );
  }

  setEnabled(enabled: boolean) {
    this.updateState({enabled});
  }

  /**
   * Adds the picked object id to the set of picked objects, and makes it the
   * picked object for single object picker observers.
   * @param oid Object id of the object to pick.
   */
  addPickedObjectId(oid: number) {
    const pickedObjectIds = new Set<number>(this.state.pickedObjectIds);
    if (oid && !pickedObjectIds.has(oid)) {
      pickedObjectIds.add(oid);
      this.updateState({pickedObjectIds});
    }
  }

  /**
   * Adds the picked object id to the set of picked objects, and if it was the
   * picked object for single object picker observers, another from the set is
   * picked instead.
   * @param oid Object id of the object to pick.
   */
  removePickedObjectId(oid: number) {
    // Ignore removal of objects that are not already picked.
    if (!this.state.pickedObjectIds.has(oid)) {
      return;
    }
    const pickedObjectIds = new Set<number>(this.state.pickedObjectIds);
    pickedObjectIds.delete(oid);
    this.updateState({pickedObjectIds});
  }

  /**
   * Clear picked entity and notify observers.
   */
  clearHoverPickedEntity() {
    this.setHoverPickedEntityId(NULL_ENTITY_ID);
  }

  /**
   * Set the id of the entity of the hovered entity
   */
  setHoverPickedEntityId(eid: number) {
    this.updateState({hoverPickedEntityId: eid});
  }

  /**
   * Clear the current picked entity and notify observers.
   */
  clearPickedEntity() {
    this.updateState({pickedEntityId: NULL_ENTITY_ID});
  }

  /**
   * Set the id of the current picked entity and notify observers.
   */
  setPickedEntityId(eid: number) {
    this.updateState({pickedEntityId: eid});
  }

  // An object can be deep picked via double clicking (both 3D and tree view).
  setDeepPickObjectId(ofid: number) {
    this.updateState({deepPickedObjectId: ofid});
  }

  clearHoverPicked() {
    this.updateState({
      hoverPickedObjectId: NULL_OBJECT_ID,
      hoverPickedEntityId: NULL_ENTITY_ID,
    });
  }

  clearDeepPick() {
    this.updateState({deepPickedObjectId: NULL_OBJECT_ID});
  }

  clearAllPickedObjects() {
    // We also clear the entity id since entities are part of objects
    this.updateState({
      pickedObjectIds: new Set<number>(),
      pickedEntityId: NULL_ENTITY_ID,
    });
  }

  /**
   * Set the id of object to be picked. This function can be used to directly
   * set a picked object. For example, if the app exposes a list/tree of
   * objects, it can send user selection here to pick a selected object.
   */
  setPickedObjectId(oid: number) {
    // Ignore if the object id is the reserved null id.
    if (oid === NULL_OBJECT_ID) {
      return;
    }

    this.updateState({pickedObjectIds: new Set<number>([oid])});
  }

  /**
   * Sets the custom selection (e.g. a service picked in an app panel).
   * This clears the picked objects.
   */
  setCustomSelection(selection: CustomSelection) {
    this.updateState({customSelection: selection});
  }

  /** Clears the picked objects, the picked entity and the custom selection. */
  clearAllPicks() {
    // We also clear the entity id since entities are part of objects
    this.updateState({
      pickedObjectIds: new Set<number>(),
      pickedEntityId: NULL_ENTITY_ID,
      customSelection: undefined,
    });
  }

  setHoverPickedObjectId(oid: number) {
    this.updateState({hoverPickedObjectId: oid});
  }
}
