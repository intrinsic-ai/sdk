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
 * Standard gRPC status codes returned by the HTTP Gateway control and data planes.
 */
export enum GrpcStatus {
  OK = 0,
  CANCELLED,
  UNKNOWN,
  INVALID_ARGUMENT,
  DEADLINE_EXCEEDED,
  NOT_FOUND,
  ALREADY_EXISTS,
  PERMISSION_DENIED,
  RESOURCE_EXHAUSTED,
  FAILED_PRECONDITION,
  ABORTED,
  OUT_OF_RANGE,
  UNIMPLEMENTED,
  INTERNAL,
  UNAVAILABLE,
  DATA_LOSS,
  UNAUTHENTICATED,
}

/** Converts a numeric gRPC status code into its human-readable name string. */
function translateGrpcCode(code: number | GrpcStatus): string {
  const numCode = Number(code);
  if (
    Number.isInteger(numCode) &&
    Object.prototype.hasOwnProperty.call(GrpcStatus, numCode)
  ) {
    return GrpcStatus[numCode];
  }
  return `UNKNOWN_CODE_${code}`;
}

/**
 * Structured error representing operational failures reported by the pubsub service.
 * Includes explicit gRPC error codes, human-readable status names, and target topic context.
 */
export class PubSubGatewayError extends Error {
  readonly code: number | GrpcStatus;
  readonly topic?: string;
  readonly dataChannel?: string;
  readonly grpcStatus: string;

  constructor(
    code: number | GrpcStatus,
    message: string,
    topic?: string,
    dataChannel?: string,
    options?: ErrorOptions,
  ) {
    const status = translateGrpcCode(code);
    const topicCtx = topic ? ` (Topic: "${topic}")` : '';
    const channelCtx = dataChannel ? ` [Channel: "${dataChannel}"]` : '';
    super(`[${code} - ${status}]${topicCtx}${channelCtx}: ${message}`, options);
    this.name = 'PubSubGatewayError';
    this.code = code;
    this.topic = topic;
    this.dataChannel = dataChannel;
    this.grpcStatus = status;
    Object.setPrototypeOf(this, new.target.prototype);
  }
}
