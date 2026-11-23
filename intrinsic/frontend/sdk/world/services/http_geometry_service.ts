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

import {create} from '@bufbuild/protobuf';
import {GetRenderableRequestSchema} from '@intrinsic-ai/types/geometry/proto/geometry_service_pb';
import {
  type RenderableWithMetadata,
  RenderableWithMetadataSchema,
} from '@intrinsic-ai/types/geometry/proto/geometry_service_types_pb';

import type {RemoteShape} from '../model/geometry';
import type {GeometryService} from '../service_interfaces/geometry_service';
import type {HttpTransport} from './http_transport';
import {decodeResponse, postProto} from './proto_json';

/**
 * A `GeometryService` that uses the platform's geometry REST API.
 *
 * TODO(b/589802342) - This service class lives here temporarily until we figure
 * out a better longer-term place for it.
 */
export class HttpGeometryService implements GeometryService {
  constructor(private readonly transport: HttpTransport) {}

  /** Gets the renderable for the given remote shape. */
  async getRenderable(shape: RemoteShape): Promise<RenderableWithMetadata> {
    const request = create(GetRenderableRequestSchema, {
      geometryStorageRefs: {
        renderableRef: shape.renderableRef,
        exactGeometryRef: shape.exactGeometryRef,
        keepRenderable: shape.keepRenderable,
      },
    });
    const response = await postProto(
      this.transport,
      'api/geometries/renderables',
      GetRenderableRequestSchema,
      request,
    );
    return decodeResponse(response, RenderableWithMetadataSchema);
  }
}
