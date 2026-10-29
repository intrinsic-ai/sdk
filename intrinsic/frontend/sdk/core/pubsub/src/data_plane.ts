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

import {type Registry, create} from '@bufbuild/protobuf';
import {type Any, timestampNow} from '@bufbuild/protobuf/wkt';
import {
  DataPacketSchema,
  type PubSubMessage,
} from '@intrinsic-ai/types/httpjson/proto/v1alpha/dataplane_pb';

import {
  DEFAULT_CLIENT_REASSEMBLY_TIMEOUT_MS,
  DataPacketAssembler,
  DataPacketChunker,
} from './chunking';
import type {DiagnosticRouter} from './diagnostics';
import {GrpcStatus, PubSubGatewayError} from './errors';

/** Shared module-level TextDecoder to avoid heap allocations on hot data path. */
const sharedTextDecoder = new TextDecoder();

/** Timeout duration for waiting for a WebRTC DataChannel to reach open state (20 seconds). */
const DATA_CHANNEL_OPEN_TIMEOUT_MS = 20_000;

/**
 * Interface representing a subscriber capable of processing reassembled Data Channel frames.
 * @internal
 */
export interface TopicReceiver {
  /** Delivers a fully reassembled message payload to the receiver handle. */
  handleMessage(message: PubSubMessage): void;
  /** Terminates message reception and releases underlying receiver resources. */
  close(): Promise<void>;
}

/**
 * Interface representing a publisher handle registered on the Data Plane.
 * @internal
 */
export interface TopicSender {
  /** Terminates publication and releases underlying sender resources. */
  close(): Promise<void>;
}

/**
 * Interface representing the WebRTC Data Plane responsible for high-frequency pub/sub data delivery,
 * chunk partitioning, reassembly, and topic multiplexing over WebRTC Data Channels.
 * @internal
 */
export interface DataPlane {
  /** Registers an active topic receiver to handle incoming streaming frames. */
  registerReceiver(topic: string, receiver: TopicReceiver): void;
  /** Deregisters an active topic receiver from incoming frame delivery. */
  unregisterReceiver(topic: string): void;
  /** Checks whether an active receiver handle is currently registered for the topic. */
  hasReceiver(topic: string): boolean;
  /** Registers an active publisher handle for the topic. */
  registerSender(topic: string, sender: TopicSender): void;
  /** Deregisters an active publisher handle from the topic. */
  unregisterSender(topic: string): void;
  /** Returns an existing publisher handle registered for the topic, if present. */
  getSender<T extends TopicSender = TopicSender>(topic: string): T | undefined;
  /** Transmits a structured message payload across the designated channel, chunking if needed. */
  publish(
    dataChannel: string,
    topic: string,
    payload: Any,
    registry?: Registry,
  ): Promise<void>;
  /** Returns an open DataChannel for the specified label, creating and awaiting negotiation if necessary. */
  getOrCreateDataChannel(label: string): Promise<RTCDataChannel>;
  /** Registers and monitors an externally created or negotiated DataChannel instance. */
  attachDataChannel(dc: RTCDataChannel): void;
}

/**
 * Manages streaming channels, payload chunking, frame assembly, and message delivery to topic handles.
 * @internal
 */
export class WebRtcDataPlane implements DataPlane {
  private readonly channels = new Map<string, ManagedDataChannel>();
  private readonly receivers = new Map<string, TopicReceiver>();
  private readonly senders = new Map<string, TopicSender>();
  private readonly chunker: DataPacketChunker;
  private readonly assembler: DataPacketAssembler;
  private readonly diag: DiagnosticRouter;
  private readonly createChannel: (label: string) => RTCDataChannel;

  constructor(
    diag: DiagnosticRouter,
    createChannel: (label: string) => RTCDataChannel,
    chunkReassemblyTimeoutMs = DEFAULT_CLIENT_REASSEMBLY_TIMEOUT_MS,
  ) {
    this.diag = diag;
    this.createChannel = createChannel;
    this.chunker = new DataPacketChunker();
    this.assembler = new DataPacketAssembler(chunkReassemblyTimeoutMs, (msg) =>
      this.diag.warn(msg),
    );
  }

  /** Registers and monitors an externally created or negotiated DataChannel instance. */
  attachDataChannel(dc: RTCDataChannel): ManagedDataChannel {
    const existing = this.channels.get(dc.label);
    if (existing && existing.rawChannel !== dc && !existing.isClosingOrClosed) {
      return existing;
    }
    if (existing && existing.rawChannel !== dc) {
      existing.close();
      this.channels.delete(dc.label);
    }
    const managed = new ManagedDataChannel(
      dc,
      this.diag,
      (label, data) => this.handleIncomingFrame(label, data),
      (label) => {
        if (this.channels.get(label) === managed) {
          this.channels.delete(label);
        }
      },
    );
    this.channels.set(dc.label, managed);
    return managed;
  }

  /** Returns an open DataChannel for the specified label, creating and awaiting negotiation if necessary. */
  async getOrCreateDataChannel(label: string): Promise<RTCDataChannel> {
    let channel = this.channels.get(label);
    if (!channel || channel.isClosingOrClosed) {
      this.diag.debug(
        `[DataChannel] Initializing label "${label}" (readyState: connecting)`,
      );
      const dc = this.createChannel(label);
      channel = this.attachDataChannel(dc);
    }
    return channel.waitUntilReady();
  }

  /** Registers an active topic receiver to handle incoming streaming frames. */
  registerReceiver(topic: string, receiver: TopicReceiver): void {
    this.receivers.set(topic, receiver);
  }

  /** Deregisters an active topic receiver from incoming frame delivery. */
  unregisterReceiver(topic: string): void {
    this.receivers.delete(topic);
  }

  /** Checks whether an active receiver handle is currently registered for the topic. */
  hasReceiver(topic: string): boolean {
    return this.receivers.has(topic);
  }

  /** Registers an active publisher handle for the topic. */
  registerSender(topic: string, sender: TopicSender): void {
    this.senders.set(topic, sender);
  }

  /** Deregisters an active publisher handle from the topic. */
  unregisterSender(topic: string): void {
    this.senders.delete(topic);
  }

  /** Returns an existing publisher handle registered for the topic, if present. */
  getSender<T extends TopicSender = TopicSender>(topic: string): T | undefined {
    return this.senders.get(topic) as T | undefined;
  }

  /** Transmits a structured message payload across the designated channel, chunking if needed. */
  async publish(
    dataChannel: string,
    topic: string,
    payload: Any,
    registry?: Registry,
  ): Promise<void> {
    const dc = await this.getOrCreateDataChannel(dataChannel);
    const packet = create(DataPacketSchema, {
      payload: {
        case: 'message',
        value: {
          topic,
          pubsubPacket: {
            payload,
            publishTime: timestampNow(),
          },
        },
      },
    });

    const chunkStrings = this.chunker.chunk(packet, registry);
    for (const str of chunkStrings) {
      dc.send(str);
    }
  }

  /** Decodes, reassembles, and routes an incoming DataChannel frame to its target topic receiver. */
  private handleIncomingFrame(
    label: string,
    rawData: ArrayBuffer | Uint8Array | string,
  ): void {
    try {
      this.diag.debug(
        `[DataChannel Rx] Received frame on label "${label}" (${typeof rawData === 'string' ? rawData.length + ' chars' : (rawData as ArrayBuffer | Uint8Array).byteLength + ' bytes'})`,
      );
      let str: string;
      if (typeof rawData === 'string') {
        str = rawData;
      } else {
        str = sharedTextDecoder.decode(rawData);
      }
      const packet = JSON.parse(str);

      const assembled = this.assembler.assemble(packet);
      if (!assembled) {
        return;
      }

      if (assembled.error !== undefined) {
        const err = assembled.error;
        this.diag.error(
          new PubSubGatewayError(
            err.code ?? GrpcStatus.UNKNOWN,
            err.message ||
              `Data channel error on label "${label}" (code: ${err.code ?? GrpcStatus.UNKNOWN})`,
            undefined,
            label,
          ),
        );
      } else if (assembled.message !== undefined) {
        const msg = assembled.message;
        if (!msg.topic) {
          return;
        }
        const receiver = this.receivers.get(msg.topic);
        if (receiver) {
          receiver.handleMessage(msg as unknown as PubSubMessage);
        } else {
          this.diag.warn(
            `Received message for unsubscribed topic "${msg.topic}" on channel "${label}"`,
          );
        }
      }
    } catch (err: unknown) {
      this.diag.error(
        new Error(`Failed to process DataChannel message on "${label}"`, {
          cause: err,
        }),
      );
    }
  }

  /** Gracefully terminates all active topic receivers, senders, and DataChannels. */
  async closeAll(): Promise<void> {
    this.assembler.dispose();

    const activeReceivers = Array.from(this.receivers.values());
    this.receivers.clear();
    for (const receiver of activeReceivers) {
      try {
        await receiver.close();
      } catch (err: unknown) {
        this.diag.warn(
          `Error closing receiver during DataPlane close: ${err instanceof Error ? err.message : err}`,
        );
      }
    }

    const activeSenders = Array.from(this.senders.values());
    this.senders.clear();
    for (const sender of activeSenders) {
      try {
        await sender.close();
      } catch (err: unknown) {
        this.diag.warn(
          `Error closing sender during DataPlane close: ${err instanceof Error ? err.message : err}`,
        );
      }
    }

    const channels = Array.from(this.channels.values());
    this.channels.clear();
    for (const ch of channels) {
      try {
        ch.close();
      } catch (err: unknown) {
        this.diag.debug(
          `Error closing DataChannel "${ch.label}": ${err instanceof Error ? err.message : err}`,
        );
      }
    }
  }

  /** Synchronously drops all active receivers, senders, DataChannels, and assembly timers during teardown. */
  dispose(): void {
    this.assembler.dispose();
    this.receivers.clear();
    this.senders.clear();
    const channels = Array.from(this.channels.values());
    this.channels.clear();
    for (const ch of channels) {
      try {
        ch.close();
      } catch (err: unknown) {
        this.diag.warn(
          `Error disposing DataChannel "${ch.label}": ${err instanceof Error ? err.message : err}`,
        );
      }
    }
  }
}

/**
 * Awaits until an RTCDataChannel reaches the 'open' state, or rejects on error, close, or timeout.
 * Automatically cleans up temporary listeners and timers.
 */
function waitForDataChannelOpen(
  dc: RTCDataChannel,
  diag: DiagnosticRouter,
  timeoutMs = DATA_CHANNEL_OPEN_TIMEOUT_MS,
): Promise<RTCDataChannel> {
  if (dc.readyState === 'open') {
    diag.debug(
      `[DataChannel] Label "${dc.label}" already open upon admission!`,
    );
    return Promise.resolve(dc);
  }
  if (dc.readyState === 'closing' || dc.readyState === 'closed') {
    return Promise.reject(
      new Error(`DataChannel "${dc.label}" is already ${dc.readyState}`),
    );
  }

  return new Promise<RTCDataChannel>((resolve, reject) => {
    let timer: ReturnType<typeof setTimeout> | undefined;

    const cleanup = () => {
      if (timer) {
        clearTimeout(timer);
        timer = undefined;
      }
      dc.removeEventListener('open', handleOpen);
      dc.removeEventListener('close', handleClose);
      dc.removeEventListener('error', handleError);
    };

    const handleOpen = () => {
      cleanup();
      diag.debug(
        `[DataChannel] Label "${dc.label}" transition complete -> OPEN!`,
      );
      resolve(dc);
    };

    const handleClose = () => {
      cleanup();
      reject(new Error(`DataChannel "${dc.label}" closed before opening`));
    };

    const handleError = (err: Event) => {
      cleanup();
      reject(
        new Error(`DataChannel "${dc.label}" error before opening`, {
          cause: err,
        }),
      );
    };

    timer = setTimeout(() => {
      cleanup();
      try {
        dc.close();
      } catch (err: unknown) {
        diag.debug(
          `Error closing timed-out DataChannel "${dc.label}": ${err instanceof Error ? err.message : err}`,
        );
      }
      reject(
        new Error(
          `DataChannel "${dc.label}" failed to open within ${timeoutMs}ms timeout. Verify peer connection state.`,
        ),
      );
    }, timeoutMs);

    dc.addEventListener('open', handleOpen);
    dc.addEventListener('close', handleClose);
    dc.addEventListener('error', handleError);
  });
}

/**
 * Wraps an RTCDataChannel instance and encapsulates its lifecycle, frame routing,
 * error handling, and teardown.
 * @internal
 */
class ManagedDataChannel {
  readonly label: string;
  readonly rawChannel: RTCDataChannel;
  private readonly readyPromise: Promise<RTCDataChannel>;

  constructor(
    dc: RTCDataChannel,
    diag: DiagnosticRouter,
    onIncomingFrame: (
      label: string,
      rawData: ArrayBuffer | Uint8Array | string,
    ) => void,
    onClosed: (label: string) => void,
  ) {
    this.rawChannel = dc;
    this.label = dc.label;
    dc.binaryType = 'arraybuffer';

    dc.onmessage = (event) => onIncomingFrame(this.label, event.data);
    dc.onerror = (err) =>
      diag.error(
        new Error(`DataChannel error on label "${this.label}"`, {cause: err}),
      );
    dc.onclose = () => onClosed(this.label);

    this.readyPromise = waitForDataChannelOpen(dc, diag).catch((err) => {
      onClosed(this.label);
      throw err;
    });
  }

  get isClosingOrClosed(): boolean {
    return (
      this.rawChannel.readyState === 'closing' ||
      this.rawChannel.readyState === 'closed'
    );
  }

  waitUntilReady(): Promise<RTCDataChannel> {
    return this.readyPromise;
  }

  send(data: string): void {
    this.rawChannel.send(data);
  }

  close(): void {
    try {
      this.rawChannel.close();
    } catch {
      // Ignore close errors during explicit teardown
    }
  }
}
