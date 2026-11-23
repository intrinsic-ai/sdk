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
  type DescMessage,
  type JsonValue,
  type MessageShape,
  fromJson,
  toJson,
} from '@bufbuild/protobuf';

import type {HttpTransport} from './http_transport';

/**
 * Sends a GET request and decodes the JSON response as a `schema` message.
 * Rejects if the response status is not 2xx.
 */
export async function getProto<Desc extends DescMessage>(
  transport: HttpTransport,
  url: string,
  schema: Desc,
  prepareJson?: (json: JsonValue) => JsonValue,
): Promise<MessageShape<Desc>> {
  const response = await transport.fetch(url);
  if (!response.ok) {
    throw await responseError(url, response);
  }
  return decodeResponse(response, schema, prepareJson);
}

/**
 * Sends a POST request with `message` as its JSON body. Rejects if the
 * response status is not 2xx.
 */
export async function postProto<Desc extends DescMessage>(
  transport: HttpTransport,
  url: string,
  schema: Desc,
  message: MessageShape<Desc>,
): Promise<Response> {
  const response = await transport.fetch(url, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(toJson(schema, message)),
  });
  if (!response.ok) {
    throw await responseError(url, response);
  }
  return response;
}

/**
 * Decodes the JSON body of `response` as a `schema` message. Unknown fields
 * are ignored so that newer backends remain compatible, and an empty body
 * decodes to the default message.
 *
 * TODO(b/589812258) - Note that no type registry is used, so the JSON must not
 * contain any `google.protobuf.Any` messages. Use `prepareJson` to remove them
 * before decoding.
 */
export async function decodeResponse<Desc extends DescMessage>(
  response: Response,
  schema: Desc,
  prepareJson: (json: JsonValue) => JsonValue = (json) => json,
): Promise<MessageShape<Desc>> {
  const text = await response.text();
  const json: JsonValue = text.trim() ? JSON.parse(text) : {};
  return fromJson(schema, prepareJson(json), {ignoreUnknownFields: true});
}

async function responseError(url: string, response: Response): Promise<Error> {
  return new Error(
    `Request to ${url} failed with status ${response.status} ` +
      `${response.statusText}: ${await response.text()}`,
  );
}
