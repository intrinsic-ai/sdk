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

import {type Registry, create, toJsonString} from '@bufbuild/protobuf';
import {
  base64Decode as base64ToBytes,
  base64Encode as bytesToBase64,
} from '@bufbuild/protobuf/wire';
import {
  type DataPacket,
  DataPacketChunkSchema,
  DataPacketSchema,
} from '@intrinsic-ai/types/httpjson/proto/v1alpha/dataplane_pb';

/** Shared module-level codec instances to avoid heap object allocation on chunking loop. */
const sharedTextDecoder = new TextDecoder();
const sharedTextEncoder = new TextEncoder();

let nextChunkId = 0;

/** Generates a unique chunk ID using crypto.randomUUID if available, or an incrementing counter. */
function generateUniqueID(): string {
  if (
    typeof crypto !== 'undefined' &&
    typeof crypto.randomUUID === 'function'
  ) {
    return crypto.randomUUID();
  }
  return `chunk-${++nextChunkId}`;
}

/** Maximum byte size allowed for a single WebRTC Data Channel packet (16 KiB). */
const MAX_CHUNK_SIZE_BYTES = 16 * 1024;

/** Maximum number of chunks permitted in a single message series (~50 MB max payload). */
const MAX_CHUNKS_PER_MESSAGE = 3200;

/** Maximum ASCII string length allowed for a chunk session ID. */
const MAX_CHUNK_ID_LENGTH = 64;

/** Maximum number of concurrent chunk reassemblies allowed in buffer to prevent client memory exhaustion. */
const MAX_CONCURRENT_ASSEMBLIES = 256;

/**
 * Default client-side chunk reassembly timeout in milliseconds (30 seconds).
 * Relaxed compared to the server's 3-second limit to tolerate background browser tab throttling.
 */
export const DEFAULT_CLIENT_REASSEMBLY_TIMEOUT_MS = 30_000;

/** Background evaluation interval in milliseconds for sweeping expired reassembly buffers (1 second). */
const REASSEMBLY_SWEEP_INTERVAL_MS = 1_000;

/** Safe Base64 character segment size to guarantee serialized chunk wrapper stays under 16 KiB. */
const CHUNK_SEGMENT_CHAR_LENGTH = 15_000;

/** Partitions large message payloads into sequential transmission frames. */
export class DataPacketChunker {
  /**
   * Partitions a packet into an array of serialized JSON strings, automatically splitting payloads exceeding MTU size.
   */
  chunk(packet: DataPacket, registry?: Registry): string[] {
    const jsonStr = toJsonString(DataPacketSchema, packet, {registry});
    const serializedBytes = sharedTextEncoder.encode(jsonStr);

    if (serializedBytes.byteLength <= MAX_CHUNK_SIZE_BYTES) {
      return [jsonStr];
    }

    const b64Data = bytesToBase64(serializedBytes);
    const totalChunks = Math.ceil(b64Data.length / CHUNK_SEGMENT_CHAR_LENGTH);

    if (totalChunks > MAX_CHUNKS_PER_MESSAGE) {
      throw new Error(
        `Message size too large to transmit over WebRTC: ${totalChunks} chunks exceeds maximum limit of ${MAX_CHUNKS_PER_MESSAGE}. Reduce payload size or partition data across multiple publications.`,
      );
    }

    const id = generateUniqueID();
    if (id.length > MAX_CHUNK_ID_LENGTH) {
      throw new Error(
        `Generated chunk ID "${id}" (${id.length} characters) exceeds limit of ${MAX_CHUNK_ID_LENGTH} characters`,
      );
    }

    const chunks: string[] = [];
    for (let i = 0; i < totalChunks; i++) {
      const start = i * CHUNK_SEGMENT_CHAR_LENGTH;
      const segment = b64Data.slice(start, start + CHUNK_SEGMENT_CHAR_LENGTH);

      const chunkData = create(DataPacketChunkSchema, {
        id,
        data: segment,
        num: i + 1,
        total: totalChunks,
      });

      const chunkPacket = create(DataPacketSchema, {
        payload: {
          case: 'chunk',
          value: chunkData,
        },
      });

      chunks.push(toJsonString(DataPacketSchema, chunkPacket, {registry}));
    }

    return chunks;
  }
}

interface PendingBuffer {
  chunks: Map<number, string>;
  total: number;
  timestamp: number;
}

/** Structural interface representing an incoming or reassembled DataChannel JSON frame. */
export interface DataChannelFrame {
  chunk?: {
    id?: string;
    total: number;
    num: number;
    data: string;
  };
  error?: {
    topic?: string;
    code?: number;
    message?: string;
    [key: string]: unknown;
  };
  message?: {
    topic?: string;
    payload?: unknown;
    [key: string]: unknown;
  };
  [key: string]: unknown;
}

/** Safely invokes .unref() on active Node.js timers without throwing in browser DOM runtimes. */
function unrefIfNodejs(timer: unknown): void {
  if (
    typeof timer === 'object' &&
    timer !== null &&
    'unref' in timer &&
    typeof timer.unref === 'function'
  ) {
    timer.unref();
  }
}

/** Aggregates incoming packet fragments and returns completed message frames. */
export class DataPacketAssembler {
  private readonly buffers = new Map<string, PendingBuffer>();
  private sweepTimer?: ReturnType<typeof setInterval>;
  private readonly timeoutMs: number;
  private readonly onWarning?: (msg: string) => void;

  constructor(
    timeoutMs = DEFAULT_CLIENT_REASSEMBLY_TIMEOUT_MS,
    onWarning?: (msg: string) => void,
  ) {
    this.timeoutMs = timeoutMs;
    this.onWarning = onWarning;
  }

  /** Starts a background interval timer to purge partial packet assemblies exceeding the timeout. */
  private startSweepTimer(): void {
    if (this.sweepTimer) {
      return;
    }
    this.sweepTimer = setInterval(() => {
      const now = Date.now();
      for (const [id, buffer] of this.buffers.entries()) {
        if (now - buffer.timestamp > this.timeoutMs) {
          this.buffers.delete(id);
          this.onWarning?.(
            `Dropped incomplete DataChannel packet due to client chunk reassembly timeout [ID: ${id}, timeout: ${this.timeoutMs}ms]`,
          );
        }
      }
      if (this.buffers.size === 0) {
        this.stopSweepTimer();
      }
    }, REASSEMBLY_SWEEP_INTERVAL_MS);
    unrefIfNodejs(this.sweepTimer);
  }

  /** Stops the background sweep interval timer if active. */
  private stopSweepTimer(): void {
    if (this.sweepTimer) {
      clearInterval(this.sweepTimer);
      this.sweepTimer = undefined;
    }
  }

  /** Evicts the oldest pending packet buffer when reaching maximum concurrent assembly limits. */
  private evictOldestPendingPacket(): void {
    let oldestId: string | null = null;
    let oldestTimestamp = Infinity;
    for (const [id, buf] of this.buffers.entries()) {
      if (buf.timestamp < oldestTimestamp) {
        oldestTimestamp = buf.timestamp;
        oldestId = id;
      }
    }
    if (oldestId !== null) {
      this.buffers.delete(oldestId);
      if (this.buffers.size === 0) {
        this.stopSweepTimer();
      }
      this.onWarning?.(
        `Evicted oldest incomplete chunk assembly [ID: ${oldestId}] due to concurrent assembly buffer limit`,
      );
    }
  }

  /** Validates chunk sequence numbers, ID format, and count boundaries. */
  private validateChunk(
    chunk: NonNullable<DataChannelFrame['chunk']>,
  ): NonNullable<DataChannelFrame['chunk']> & {id: string} {
    if (!chunk.id) {
      throw new Error('Chunk frame is missing a required chunk ID');
    }
    if (!Number.isInteger(chunk.num) || !Number.isInteger(chunk.total)) {
      throw new Error(
        `Chunk numbers must be valid integers (received num: ${chunk.num}, total: ${chunk.total})`,
      );
    }

    if (chunk.id.length > MAX_CHUNK_ID_LENGTH) {
      throw new Error(
        `Chunk ID length ${chunk.id.length} exceeds limit of ${MAX_CHUNK_ID_LENGTH}`,
      );
    }
    if (chunk.total > MAX_CHUNKS_PER_MESSAGE) {
      throw new Error(
        `Chunk total ${chunk.total} exceeds maximum limit of ${MAX_CHUNKS_PER_MESSAGE}`,
      );
    }
    if (chunk.num < 1 || chunk.num > chunk.total) {
      throw new Error(
        `Chunk sequence number ${chunk.num} is out of range [1, ${chunk.total}]`,
      );
    }
    // Safe type assertion: chunk.id has been validated as a non-empty string above.
    return chunk as NonNullable<DataChannelFrame['chunk']> & {id: string};
  }

  /**
   * Processes an incoming frame, returning the fully assembled payload once all fragments arrive, or null if pending.
   */
  assemble(packet: DataChannelFrame): DataChannelFrame | null {
    if (packet.chunk === undefined) {
      return packet;
    }

    const chunk = this.validateChunk(packet.chunk);

    let buffer = this.buffers.get(chunk.id);
    if (!buffer) {
      if (this.buffers.size >= MAX_CONCURRENT_ASSEMBLIES) {
        this.evictOldestPendingPacket();
      }
      buffer = {
        chunks: new Map<number, string>(),
        total: chunk.total,
        timestamp: Date.now(),
      };
      this.buffers.set(chunk.id, buffer);
      this.startSweepTimer();
    } else if (buffer.total !== chunk.total) {
      throw new Error(
        `Chunk total ${chunk.total} does not match previous total ${buffer.total} for ID ${chunk.id}`,
      );
    }

    buffer.timestamp = Date.now();
    buffer.chunks.set(chunk.num, chunk.data);

    if (buffer.chunks.size === buffer.total) {
      let unifiedBase64 = '';
      for (let i = 1; i <= buffer.total; i++) {
        const segment = buffer.chunks.get(i);
        if (segment === undefined) {
          throw new Error(
            `Missing chunk sequence number ${i} during reassembly for message ID "${chunk.id}"`,
          );
        }
        unifiedBase64 += segment;
      }
      this.buffers.delete(chunk.id);
      if (this.buffers.size === 0) {
        this.stopSweepTimer();
      }

      const decodedBytes = base64ToBytes(unifiedBase64);
      const jsonStr = sharedTextDecoder.decode(decodedBytes);

      let reassemble: DataChannelFrame;
      try {
        // Safe type assertion: parsed JSON string represents a DataChannelFrame object structure.
        reassemble = JSON.parse(jsonStr) as DataChannelFrame;
      } catch (err: unknown) {
        throw new Error(
          `Failed to parse reassembled JSON payload for message ID "${chunk.id}"`,
          {cause: err},
        );
      }

      if (reassemble.chunk !== undefined) {
        throw new Error('Nested packet chunking is not permitted');
      }
      return reassemble;
    }

    return null;
  }

  /** Stops background cleanup timers and discards any buffered packet fragments. */
  dispose(): void {
    this.stopSweepTimer();
    this.buffers.clear();
  }
}
