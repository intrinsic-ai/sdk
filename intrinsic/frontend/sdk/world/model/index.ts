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
 * @file Public API of `@intrinsic-ai/world/model`.
 */
export * from './collision_rules';
export * from './entity';
export * from './geometry';
export * from './joint_limits';
export * from './material_properties';
export * from './physics';
export * from './wait_for_world_model';
export * from './world_backend';
export * from './world_ids';
export * from './world_lookup';
export * from './world_model';
export * from './world_object';
export * from './world_object_builder';

// `entity` and `joint_limits` both export a `JointLimits` type. Explicit
// re-exports take precedence over the `export *` above and resolve the
// ambiguity.
export {JointLimits} from './joint_limits';
export type {JointLimits as EntityJointLimits} from './entity';
