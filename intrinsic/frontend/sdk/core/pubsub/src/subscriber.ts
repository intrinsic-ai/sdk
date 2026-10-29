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
} from '@bufbuild/protobuf';
import {type Timestamp, timestampDate} from '@bufbuild/protobuf/wkt';
import type {PubSubMessage} from '@intrinsic-ai/types/httpjson/proto/v1alpha/dataplane_pb';

import type {ControlPlane} from './control_plane';
import type {DataPlane, TopicReceiver} from './data_plane';
import {GrpcStatus, PubSubGatewayError} from './errors';
import type {MessageMetadata, RateLimitConfig, Subscription} from './types';

/** Shared module-level TextDecoder to avoid heap allocations on hot data path. */
const sharedTextDecoder = new TextDecoder();

/** Pure function to intelligently decode Uint8Array or JSON-formatted string payloads without throwing. */
function tryUnwrapJson(val: unknown): unknown {
  if (val instanceof Uint8Array) {
    try {
      return JSON.parse(sharedTextDecoder.decode(val));
    } catch {
      return val;
    }
  }
  if (typeof val === 'string' && (val.startsWith('{') || val.startsWith('['))) {
    try {
      return JSON.parse(val);
    } catch {
      return val;
    }
  }
  return val;
}

/** Pure function to normalize arbitrary Protobuf or JSON timestamps into standard runtime representations. */
function normalizeTimestamp(input: unknown): number | Date | undefined {
  if (typeof input === 'number') {
    return input;
  }
  if (typeof input === 'string') {
    return new Date(input);
  }
  if (typeof input === 'object' && input !== null && 'seconds' in input) {
    // Safe assertion: Objects possessing a seconds property align with Protobuf Timestamp layout.
    try {
      return timestampDate(input as Timestamp);
    } catch (err: unknown) {
      console.debug('Failed to convert Protobuf Timestamp to Date:', err);
      return undefined;
    }
  }
  return undefined;
}

/**
 * Active subscription handle that receives streaming message frames.
 */
export class SubscriptionHandle<Desc extends DescMessage = DescMessage>
  implements Subscription<Desc>, TopicReceiver
{
  private state: 'initializing' | 'open' | 'closed' = 'initializing';
  private openPromise?: Promise<void>;
  private readonly expectedTypeSuffix?: string;

  private readonly controlPlane: ControlPlane;
  private readonly dataPlane: DataPlane;
  readonly topic: string;
  readonly schema: Desc | undefined;
  readonly dataChannel: string;
  private readonly onMessage: (
    message: MessageShape<Desc>,
    meta: MessageMetadata,
  ) => void;
  private readonly onError?: (error: unknown) => void;
  private readonly rateLimit?: RateLimitConfig;

  constructor(
    controlPlane: ControlPlane,
    dataPlane: DataPlane,
    topic: string,
    schema: Desc | undefined,
    dataChannel: string,
    onMessage: (message: MessageShape<Desc>, meta: MessageMetadata) => void,
    onError?: (error: unknown) => void,
    rateLimit?: RateLimitConfig,
  ) {
    this.controlPlane = controlPlane;
    this.dataPlane = dataPlane;
    this.topic = topic;
    this.schema = schema;
    this.dataChannel = dataChannel;
    this.onMessage = onMessage;
    this.onError = onError;
    this.rateLimit = rateLimit;

    if (this.schema?.typeName) {
      this.expectedTypeSuffix = `/${this.schema.typeName}`;
    }
  }

  /**
   * Registers topic subscription interest with the service.
   * @throws {Error} If the subscription handle is closed during negotiation.
   */
  async open(): Promise<void> {
    if (this.state === 'closed') {
      throw new Error(
        `Cannot open subscription for topic "${this.topic}": subscription handle is closed. Call client.subscribe() to create a new subscription instance.`,
      );
    }
    if (this.state === 'open') {
      return;
    }
    if (this.openPromise) {
      return this.openPromise;
    }
    this.openPromise = this.performOpen();
    return this.openPromise;
  }

  /** Registers receiver references and sends a subscribe command over the control link. */
  private async performOpen(): Promise<void> {
    try {
      this.dataPlane.registerReceiver(this.topic, this);
      await this.controlPlane.subscribeTopic(
        this.topic,
        this.dataChannel,
        this.schema?.typeName || undefined,
        this.rateLimit,
      );
      if (this.state === 'closed') {
        try {
          await this.controlPlane.unsubscribeTopic(
            this.topic,
            this.schema?.typeName || undefined,
          );
        } catch (err: unknown) {
          console.warn(
            `Failed to unsubscribe topic "${this.topic}" during rollback:`,
            err,
          );
        }
        throw new Error(
          `Subscription handle for topic "${this.topic}" was closed during open negotiation. Create a new subscription instance using client.subscribe().`,
        );
      }
      this.state = 'open';
    } catch (err: unknown) {
      this.state = 'closed';
      this.dataPlane.unregisterReceiver(this.topic);
      throw err;
    } finally {
      this.openPromise = undefined;
    }
  }

  /** Dispatches an error event to the registered listener or standard error log. */
  private handleError(error: unknown): void {
    if (this.onError) {
      this.onError(error);
    } else {
      console.error(
        `Unhandled error in subscription for topic "${this.topic}":`,
        error,
      );
    }
  }

  /** Unpacks, validates, and dispatches an incoming streaming frame to registered message callbacks. */
  handleMessage(msg: PubSubMessage): void {
    if (this.state === 'closed') {
      return;
    }
    const pkt = msg.pubsubPacket;
    if (!pkt || !pkt.payload) {
      this.handleError(
        new PubSubGatewayError(
          GrpcStatus.INVALID_ARGUMENT,
          `Received empty or invalid PubSubPacket payload for topic "${this.topic}"`,
          this.topic,
          this.dataChannel,
        ),
      );
      return;
    }

    const anyPayload = pkt.payload;
    const typeUrl = String(
      anyPayload.typeUrl ||
        (anyPayload as Record<string, unknown>)['@type'] ||
        '',
    );
    if (
      this.schema?.typeName &&
      typeUrl !== this.schema.typeName &&
      (!this.expectedTypeSuffix || !typeUrl.endsWith(this.expectedTypeSuffix))
    ) {
      const err = new PubSubGatewayError(
        GrpcStatus.INVALID_ARGUMENT,
        `Type mismatch drop on topic "${this.topic}": expected ${this.schema.typeName}, received ${typeUrl}`,
        this.topic,
        this.dataChannel,
      );
      this.handleError(err);
      return;
    }

    let unpacked: MessageShape<Desc> | unknown;
    const publishTime = normalizeTimestamp(pkt.publishTime);

    try {
      if (!this.schema) {
        // Schemaless JSON / Raw mode
        unpacked = tryUnwrapJson(anyPayload);
      } else {
        // Canonical ProtoJSON mode
        unpacked = fromJson(this.schema, anyPayload as unknown as JsonValue, {
          ignoreUnknownFields: true,
        }) as MessageShape<Desc>;
      }
    } catch (err: unknown) {
      this.handleError(err);
      return;
    }

    const meta: MessageMetadata = {
      topic: this.topic,
      publishTime,
      traceId: pkt.traceId ? String(pkt.traceId) : undefined,
      spanId: pkt.spanId ? String(pkt.spanId) : undefined,
    };

    try {
      this.onMessage(unpacked as MessageShape<Desc>, meta);
    } catch (err: unknown) {
      this.handleError(err);
    }
  }

  async close(): Promise<void> {
    return this.unsubscribe();
  }

  /** Stops message reception and releases server-side streaming resources. */
  async unsubscribe(): Promise<void> {
    if (this.state === 'closed') return;
    const wasOpen = this.state === 'open';
    this.state = 'closed';
    this.dataPlane.unregisterReceiver(this.topic);
    if (this.openPromise) {
      try {
        await this.openPromise;
      } catch (err: unknown) {
        console.warn(
          `Error awaiting subscription open during unsubscribe for topic "${this.topic}":`,
          err,
        );
      }
      return;
    }
    try {
      if (wasOpen) {
        await this.controlPlane.unsubscribeTopic(
          this.topic,
          this.schema?.typeName || undefined,
        );
      }
    } catch (err: unknown) {
      console.warn(
        `Error while sending unsubscribe for topic "${this.topic}":`,
        err,
      );
    }
  }
}
