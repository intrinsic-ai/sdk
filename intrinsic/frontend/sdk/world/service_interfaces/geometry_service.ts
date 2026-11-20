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

import type {RenderableWithMetadata} from '@intrinsic-ai/types/geometry/proto/geometry_service_types_pb';

import type {RemoteShape} from '../model/geometry';

/**
 * The subset of the geometry service that the world viewer depends on.
 */
export interface GeometryService {
  /** Gets the renderable for the given remote shape. */
  getRenderable(shape: RemoteShape): Promise<RenderableWithMetadata>;
}
