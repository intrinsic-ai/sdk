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

import type {DescMessage, MessageShape} from '@bufbuild/protobuf';

/** Connection settings for PubSubClient. */
export interface PubSubClientOptions {
  /**
   * Target endpoint for the HTTP gateway control plane WebSocket connection.
   * If relative, resolved against document.baseURI or window.location.
   * To conform with transport encryption security rules,
   * remote non-loopback endpoints and secure browser contexts automatically enforce encrypted WebSockets (`wss://`).
   * @default '/webpubsub/v1alpha/control-connection'
   */
  url?: string;
  /** Optional array of RTCIceServer configurations for NAT traversal and relay (STUN/TURN). */
  iceServers?: readonly RTCIceServer[];
  /**
   * Client-side chunk reassembly timeout in milliseconds.
   * Configures predictable client tolerance for delayed packet segments.
   * @default 30000
   */
  chunkReassemblyTimeoutMs?: number;
}

/** Metadata accompanied with every received topic message. */
export interface MessageMetadata {
  /** Identifier of the Intrinsic topic. */
  readonly topic: string;
  /** Timestamp at which the message was published (epoch milliseconds number or Date). */
  readonly publishTime?: number | Date;
  /** OpenCensus trace ID context, if present. */
  readonly traceId?: string;
  /** OpenCensus span ID context, if present. */
  readonly spanId?: string;
}

/** Rate limiting configuration using server-side token buckets. */
export interface RateLimitConfig {
  /** Steady-state message volume limit per second. */
  messagesPerSecond: number;
  /** Maximum instantaneous message burst allowance (token bucket capacity). */
  tokenBucketCapacity: number;
}

/** Options for initiating a topic subscription. */
export interface SubscribeOptions<Desc extends DescMessage = DescMessage> {
  /** Topic identifier string in the system. */
  topic: string;
  /** \@bufbuild/protobuf schema definition for automatic Any unwrapping and type verification. Optional for schemaless JSON mode. */
  schema?: Desc;
  /**
   * WebRTC Data Channel label to route updates over.
   * If omitted, defaults to a shared Data Channel named "default".
   * @default 'default'
   */
  dataChannel?: string;
  /** Optional server-side throttling rules. */
  rateLimit?: RateLimitConfig;
}

/** Active subscription handle to manage ongoing topic reception. */
export interface Subscription<Desc extends DescMessage = DescMessage> {
  readonly topic: string;
  readonly dataChannel: string;
  readonly schema?: Desc;
  /**
   * Stops message reception and releases server-side streaming resources.
   * @returns A promise resolving when unsubscription completes.
   */
  unsubscribe(): Promise<void>;
}

/** Options for advertising and publishing to a topic. */
export interface PublishOptions<Desc extends DescMessage = DescMessage> {
  /** Target topic identifier string. */
  topic: string;
  /** \@bufbuild/protobuf schema definition for payload serialization and type safety. */
  schema: Desc;
  /**
   * WebRTC Data Channel label for transmission.
   * If omitted, defaults to a shared Data Channel named "default".
   * @default 'default'
   */
  dataChannel?: string;
}

/** Active publisher handle to send strongly-typed real-time messages. */
export interface Publisher<Desc extends DescMessage = DescMessage> {
  readonly topic: string;
  readonly dataChannel: string;
  /**
   * Transmits a message payload to all active subscribers on the topic.
   *
   * @param payload The structured data or JSON message payload to send.
   * @throws {Error} If the publisher is closed or serialization fails.
   */
  publish(payload: MessageShape<Desc>): Promise<void>;
  /**
   * Stops publication and deregisters this sender handle from the topic.
   *
   * @returns A promise resolving when topic unadvertisement completes.
   */
  close(): Promise<void>;
}

/** High-level connection state of the pubsub client, completely agnostic of transport protocol details. */
export type ConnectionState =
  'disconnected' | 'connecting' | 'connected' | 'reconnecting' | 'failed';

/** Minimal observable subscriber interface compatibility for UI reactive adaptations. */
export interface ObservableObserver<T> {
  /** Called whenever a new data frame arrives on the stream. */
  next(value: T): void;
  /** Called if streaming encounters an irrecoverable operational error. */
  error?(err: unknown): void;
  /** Called when streaming concludes or is explicitly terminated. */
  complete?(): void;
}

/** Minimal reactive stream subscription handle. */
export interface ObservableSubscription {
  /** Terminates message reception for this reactive subscription. */
  unsubscribe(): void;
}
