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

import {DEFAULT_CLIENT_REASSEMBLY_TIMEOUT_MS} from './chunking';
import {WebSocketControlPlane} from './control_plane';
import {WebRtcDataPlane} from './data_plane';
import {DiagnosticRouter, EventEmitter} from './diagnostics';
import {PubSubGatewayError} from './errors';
import {PublisherHandle} from './publisher';
import {SignalingController} from './signaling_controller';
import {SubscriptionHandle} from './subscriber';
import type {
  ConnectionState,
  MessageMetadata,
  ObservableObserver,
  ObservableSubscription,
  PubSubClientOptions,
  PublishOptions,
  Publisher,
  SubscribeOptions,
  Subscription,
} from './types';

/** Timeout duration in milliseconds for fetching transient TURN relay credentials (5 seconds). */
const TURN_CREDENTIAL_FETCH_TIMEOUT_MS = 5_000;

/** Returns the base document URI or window location for resolving relative SDK URLs. */
function getBrowserBaseUrl(): string {
  // Respect explicit HTML <base href="..."> tags if defined in the DOM document.
  if (typeof document !== 'undefined' && document.baseURI) {
    return document.baseURI;
  }
  // Fall back to current window location for standard browser environments.
  if (typeof window !== 'undefined' && window.location?.href) {
    return window.location.href;
  }
  // Default to localhost fallback for non-browser runtimes (e.g. Node.js or SSR unit tests).
  return 'http://localhost/';
}

/** Parses a target URL string relative to the environment base URL and validates its scheme. */
function parseAndValidateUrl(targetUrl: string): URL {
  const url = new URL(targetUrl, getBrowserBaseUrl());
  if (!['ws:', 'wss:', 'http:', 'https:'].includes(url.protocol)) {
    throw new Error(
      `Unsupported URL protocol scheme "${url.protocol}". Only ws:, wss:, http:, and https: are allowed.`,
    );
  }
  return url;
}

/** Resolves a WebSocket control URL (wss:// or ws://). */
function resolveWebSocketUrl(targetUrl: string): string {
  const url = parseAndValidateUrl(targetUrl);

  const isLoopback =
    url.hostname === 'localhost' ||
    url.hostname === '127.0.0.1' ||
    url.hostname === '[::1]' ||
    url.hostname.endsWith('.localhost');
  const isSecureContext =
    (typeof window !== 'undefined' && window.location?.protocol === 'https:') ||
    (typeof document !== 'undefined' &&
      document.baseURI?.startsWith('https://'));

  if (url.protocol === 'https:') {
    url.protocol = 'wss:';
  } else if (url.protocol === 'http:' || url.protocol === 'ws:') {
    url.protocol = isSecureContext || !isLoopback ? 'wss:' : 'ws:';
  }

  return url.toString();
}

/** Resolves an HTTP REST RPC URL by replacing the control-connection path. */
function resolveHttpRpcUrl(controlUrl: string, rpcUrl: string): string {
  const url = parseAndValidateUrl(controlUrl);

  // Map WebSocket protocol schemes to HTTP schemes
  if (url.protocol === 'wss:') {
    url.protocol = 'https:';
  } else if (url.protocol === 'ws:') {
    url.protocol = 'http:';
  }

  const basePath = url.pathname.replace(
    /\/webpubsub\/v1alpha\/control-connection\/?$/,
    '',
  );
  url.pathname = `${basePath}${rpcUrl}`;
  url.search = '';
  return url.toString();
}

/**
 * Client interface for real-time publishing and subscribing to structured system topics.
 * Manages underlying connection negotiation, stream routing, and automatic re-connection.
 */
export class PubSubClient {
  private pc?: RTCPeerConnection;
  private dataPlane?: WebRtcDataPlane;
  private connectPromise?: Promise<void>;
  private readonly diag = new DiagnosticRouter();
  private readonly stateChanges = new EventEmitter<ConnectionState>();
  private readonly openEvents = new EventEmitter<void>();
  private readonly closeEvents = new EventEmitter<void>();

  private readonly url: string;
  private readonly explicitIceServers?: readonly RTCIceServer[];
  private readonly chunkReassemblyTimeoutMs: number;
  private readonly controlPlane: WebSocketControlPlane;
  private readonly signaling: SignalingController;

  constructor(options: PubSubClientOptions = {}) {
    this.url = options.url || '/webpubsub/v1alpha/control-connection';
    this.explicitIceServers = options.iceServers;
    this.chunkReassemblyTimeoutMs =
      options.chunkReassemblyTimeoutMs || DEFAULT_CLIENT_REASSEMBLY_TIMEOUT_MS;

    this.controlPlane = new WebSocketControlPlane(this.diag, () => {
      this.cleanup();
      this.closeEvents.emit();
    });

    this.signaling = new SignalingController(
      this.controlPlane,
      this.diag,
      (state) => this.stateChanges.emit(state),
    );
    this.controlPlane.setSignalingReceiver(this.signaling);
  }

  /** Registers a callback for network connection state transitions. Returns an unsubscribe handle. */
  onStateChange(listener: (state: ConnectionState) => void): () => void {
    return this.stateChanges.subscribe(listener);
  }

  /** Registers a callback for operational streaming errors. Returns an unsubscribe handle. */
  onError(listener: (error: PubSubGatewayError | Error) => void): () => void {
    return this.diag.onError(listener as (error: unknown) => void);
  }

  /** Registers a callback for non-fatal warning messages. Returns an unsubscribe handle. */
  onWarning(listener: (message: string) => void): () => void {
    return this.diag.onWarning(listener);
  }

  /** Registers a callback for verbose debugging logs. Returns an unsubscribe handle. */
  onDebug(listener: (message: string) => void): () => void {
    return this.diag.onDebug(listener);
  }

  /** Registers a callback invoked when the client connects and becomes ready for streaming. Returns an unsubscribe handle. */
  onOpen(listener: () => void): () => void {
    return this.openEvents.subscribe(listener);
  }

  /** Registers a callback invoked when the client disconnects or terminates. Returns an unsubscribe handle. */
  onClose(listener: () => void): () => void {
    return this.closeEvents.subscribe(listener);
  }

  /**
   * Establishes a real-time connection with the pubsub service.
   * Resolves when the client is fully connected and ready to publish or subscribe.
   *
   * @throws {Error} If connection negotiation or authentication times out or fails.
   */
  async connect(): Promise<void> {
    if (this.pc || this.dataPlane) {
      return;
    }
    if (this.connectPromise) {
      return this.connectPromise;
    }
    this.connectPromise = this.performConnect();
    return this.connectPromise;
  }

  /** Negotiates TURN server credentials, WebSocket control connection, and initial WebRTC session transport. */
  private async performConnect(): Promise<void> {
    try {
      const iceServers = await this.getIceServers(this.url);
      const wsUrl = resolveWebSocketUrl(this.url);
      await this.controlPlane.connect(wsUrl);

      this.pc = new RTCPeerConnection({iceServers: [...iceServers]});
      const pc = this.pc;

      this.dataPlane = new WebRtcDataPlane(
        this.diag,
        (label: string) => pc.createDataChannel(label, {ordered: true}),
        this.chunkReassemblyTimeoutMs,
      );
      const dp = this.dataPlane;

      pc.ondatachannel = (event) => dp.attachDataChannel(event.channel);

      await this.signaling.negotiate(pc, () =>
        dp.getOrCreateDataChannel('default'),
      );
      this.openEvents.emit();
    } catch (err: unknown) {
      this.close();
      throw err;
    } finally {
      this.connectPromise = undefined;
    }
  }

  /** Resolves ICE configuration: uses explicit servers if configured, otherwise fetches TURN with STUN fallback. */
  private async getIceServers(
    controlPath: string,
  ): Promise<readonly RTCIceServer[]> {
    if (this.explicitIceServers && this.explicitIceServers.length > 0) {
      return this.explicitIceServers;
    }

    const defaultStun: RTCIceServer[] = [
      {urls: 'stun:stun.l.google.com:19302'},
    ];

    const turnServer = await this.fetchTurnCredentials(controlPath);
    if (turnServer) {
      return [...defaultStun, turnServer];
    }
    return defaultStun;
  }

  /** Queries the REST gateway to generate transient TURN relay credentials. */
  private async fetchTurnCredentials(
    controlPath: string,
  ): Promise<RTCIceServer | null> {
    if (typeof fetch === 'undefined') {
      return null;
    }
    try {
      const httpUrl = resolveHttpRpcUrl(
        controlPath,
        '/webpubsub/v1alpha:generateTurnCredentials',
      );
      this.diag.debug(
        `[TURN] Fetching transient TURN credentials from ${httpUrl}...`,
      );
      const res = await fetch(httpUrl, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: '{}',
        signal: AbortSignal.timeout(TURN_CREDENTIAL_FETCH_TIMEOUT_MS),
      });
      if (res.ok) {
        const data = await res.json();
        if (data.uris && Array.isArray(data.uris) && data.uris.length > 0) {
          this.diag.debug(
            `[TURN] Successfully resolved TURN relay servers alongside default STUN.`,
          );
          return {
            urls: data.uris,
            username: data.username || undefined,
            credential: data.password || undefined,
          };
        }
      }
      this.diag.debug(
        `[TURN] Transient TURN credential endpoint returned HTTP ${res.status}. Defaulting to STUN only.`,
      );
    } catch (err: unknown) {
      this.diag.debug(
        `[TURN] Unable to fetch TURN credentials (${err instanceof Error ? err.message : err}). Defaulting to STUN only.`,
      );
    }
    return null;
  }

  // --- Public API: Subscribe, Publish, Streams & Observables ---

  /**
   * Subscribes to a topic and streams received messages directly to the provided callback functions.
   *
   * @param options Target topic identifier, optional protobuf schema for automatic type validation, and rate-limiting rules.
   * @param onMessage Callback invoked whenever a new message frame arrives along with timestamps and metadata.
   * @param onError Optional callback invoked if streaming encounters a transport or decoding failure.
   * @returns A promise resolving to an active subscription handle allowing programmatic teardown via `.unsubscribe()`.
   * @throws {Error} If the client is not connected or already subscribed to this topic.
   */
  async subscribe<Desc extends DescMessage = DescMessage>(
    options: SubscribeOptions<Desc>,
    onMessage: (message: MessageShape<Desc>, meta: MessageMetadata) => void,
    onError?: (error: unknown) => void,
  ): Promise<Subscription<Desc>> {
    if (!this.dataPlane) {
      throw new Error(
        'PubSubClient is not connected. Call connect() before subscribing to topics.',
      );
    }
    const dp = this.dataPlane;
    const dcLabel = options.dataChannel || 'default';
    await dp.getOrCreateDataChannel(dcLabel);

    if (dp.hasReceiver(options.topic)) {
      throw new Error(
        `Already subscribed to topic "${options.topic}" on this client instance`,
      );
    }

    const sub = new SubscriptionHandle(
      this.controlPlane,
      dp,
      options.topic,
      options.schema,
      dcLabel,
      onMessage,
      onError,
      options.rateLimit,
    );
    await sub.open();
    return sub;
  }

  /**
   * Creates an active sender handle for publishing structured messages to a topic.
   *
   * @param options Target topic identifier string and protobuf schema definition.
   * @returns A promise resolving to a handle supporting `.publish()` and `.close()`.
   * @throws {Error} If the client is not connected or publication negotiation fails.
   */
  async createPublisher<Desc extends DescMessage = DescMessage>(
    options: PublishOptions<Desc>,
  ): Promise<Publisher<Desc>> {
    if (!this.dataPlane) {
      throw new Error(
        'PubSubClient is not connected. Call connect() before creating publishers.',
      );
    }
    const dp = this.dataPlane;
    const dcLabel = options.dataChannel || 'default';
    await dp.getOrCreateDataChannel(dcLabel);

    const existing = dp.getSender<Publisher<Desc>>(options.topic);
    if (existing) {
      return existing;
    }

    const pub = new PublisherHandle(
      this.controlPlane,
      dp,
      options.topic,
      options.schema,
      dcLabel,
    );
    await pub.advertise();
    return pub;
  }

  /**
   * Subscribes to a topic as an async iterable stream (`for await...of`).
   * Exiting the loop via break, return, or throwing an error automatically closes the subscription and releases resources.
   *
   * @param options Target topic identifier and optional protobuf schema definition.
   * @returns An async iterable yielding `{message, meta}` pairs as frames arrive.
   */
  subscribeStream<Desc extends DescMessage = DescMessage>(
    options: SubscribeOptions<Desc>,
  ): AsyncIterable<{message: MessageShape<Desc>; meta: MessageMetadata}> {
    // use arrow function closure to capture this.subscribe to avoid no-this-alias warning
    const subscribe = (
      onMessage: (message: MessageShape<Desc>, meta: MessageMetadata) => void,
      onError: (err: unknown) => void,
    ) => this.subscribe(options, onMessage, onError);

    return {
      async *[Symbol.asyncIterator]() {
        const events = new AsyncEventQueue<{
          message: MessageShape<Desc>;
          meta: MessageMetadata;
        }>();

        const sub = await subscribe(
          (message, meta) => events.push({message, meta}),
          (err) => events.error(err),
        );

        try {
          yield* events;
        } finally {
          events.close();
          await sub.unsubscribe();
        }
      },
    };
  }

  /**
   * Subscribes to a topic as a reactive Observable stream (RxJS compatible).
   *
   * @param options Target topic identifier and optional protobuf schema definition.
   * @returns An object supporting `.subscribe({ next, error, complete })` that returns an `.unsubscribe()` handle.
   */
  observe<Desc extends DescMessage = DescMessage>(
    options: SubscribeOptions<Desc>,
  ): {
    subscribe(
      observer: ObservableObserver<{
        message: MessageShape<Desc>;
        meta: MessageMetadata;
      }>,
    ): ObservableSubscription;
  } {
    return {
      subscribe: (observer) => {
        let activeSub: Subscription<Desc> | undefined;
        let cancelled = false;

        this.subscribe(
          options,
          (message, meta) => {
            if (!cancelled) observer.next({message, meta});
          },
          (err) => {
            if (!cancelled) observer.error?.(err);
          },
        )
          .then((sub) => {
            if (cancelled) {
              void sub.unsubscribe();
            } else {
              activeSub = sub;
            }
          })
          .catch((err: unknown) => {
            if (!cancelled) observer.error?.(err);
          });

        return {
          unsubscribe: () => {
            cancelled = true;
            if (activeSub) void activeSub.unsubscribe();
          },
        };
      },
    };
  }

  /**
   * Terminates all active publishers and subscriptions, drops open connections, and releases client resources.
   */
  close(): void {
    this.connectPromise = undefined;
    if (this.dataPlane) {
      this.dataPlane
        .closeAll()
        .catch((err: unknown) =>
          this.diag.warn(
            `Error closing data plane during client close: ${err instanceof Error ? err.message : String(err)}`,
          ),
        );
      this.dataPlane = undefined;
    }
    this.signaling.close();
    this.controlPlane.close();
    if (this.pc) {
      try {
        this.pc.close();
      } catch (err: unknown) {
        this.diag.debug(
          `Error closing PeerConnection: ${err instanceof Error ? err.message : err}`,
        );
      }
      this.pc = undefined;
    }
  }

  private cleanup(): void {
    this.connectPromise = undefined;
    if (this.dataPlane) {
      this.dataPlane.dispose();
      this.dataPlane = undefined;
    }
    this.signaling.close();
    if (this.pc) {
      try {
        this.pc.close();
      } catch (err: unknown) {
        console.warn('Error disposing PeerConnection during cleanup:', err);
      }
      this.pc = undefined;
    }
  }
}

/** Maximum number of unconsumed events permitted in an AsyncEventQueue buffer before evicting the oldest items (1,000 items). */
const MAX_ASYNC_EVENT_QUEUE_SIZE = 1_000;

/**
 * Asynchronous FIFO event queue that bridges push-based event callbacks to pull-based async iteration.
 */
class AsyncEventQueue<T> implements AsyncIterable<T> {
  private readonly queue: T[] = [];
  private resolveNext?: () => void;
  private finished = false;
  private streamError?: unknown;

  /** Enqueues an item and immediately notifies waiting stream iterators, evicting the oldest item if capacity is exceeded. */
  push(item: T): void {
    if (this.finished) return;
    if (this.queue.length >= MAX_ASYNC_EVENT_QUEUE_SIZE) {
      this.queue.shift();
    }
    this.queue.push(item);
    this.notify();
  }

  /** Terminates future iteration with the specified error. */
  error(err: unknown): void {
    this.streamError = err;
    this.close();
  }

  /** Marks the stream as complete and releases waiting iterators. */
  close(): void {
    if (this.finished) return;
    this.finished = true;
    this.notify();
  }

  private notify(): void {
    if (this.resolveNext) {
      this.resolveNext();
      this.resolveNext = undefined;
    }
  }

  /** Yields queued items sequentially as they arrive. */
  async *[Symbol.asyncIterator](): AsyncGenerator<T> {
    while (!this.finished || this.queue.length > 0) {
      if (this.queue.length === 0) {
        await new Promise<void>((resolve) => {
          this.resolveNext = resolve;
        });
      }
      if (this.streamError) throw this.streamError;
      const items = this.queue.splice(0, this.queue.length);
      for (const item of items) {
        yield item;
      }
    }
  }
}
