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
  type MessageShape,
  create,
  fromJson,
  toJsonString,
} from '@bufbuild/protobuf';
import {
  AdvertisePublisherSchema,
  type ControlPacket,
  ControlPacketSchema,
  RateLimitSchema,
  SubscribeSchema,
  UnadvertisePublisherSchema,
  UnsubscribeSchema,
} from '@intrinsic-ai/types/httpjson/proto/v1alpha/controlplane_pb';
import {
  RtcIceCandidateSchema,
  RtcSessionDescriptionSchema,
  RtcSessionDescription_SdpType,
} from '@intrinsic-ai/types/httpjson/proto/v1alpha/signalling_pb';

import type {DiagnosticRouter} from './diagnostics';
import {GrpcStatus, PubSubGatewayError} from './errors';
import type {RateLimitConfig} from './types';

/** Shared module-level TextDecoder to avoid heap allocations on hot control signaling path. */
const sharedTextDecoder = new TextDecoder();

/** Default timeout in milliseconds for establishing a WebSocket control connection. */
const DEFAULT_WS_CONNECT_TIMEOUT_MS = 10_000;

/**
 * Receiver interface for incoming WebRTC signaling commands over the WebSocket control plane.
 * @internal
 */
export interface SignalingReceiver {
  /** Called when a remote session description arrives from the server. */
  onRemoteSessionDescription(type: RTCSdpType, sdp: string): void;
  /** Called when a remote network ICE candidate arrives from the server. */
  onRemoteIceCandidate(candidate: RTCIceCandidateInit): void;
}

/**
 * Interface representing the WebRTC Control Plane responsible for signaling
 * topic advertisements, subscriptions, and SDP/ICE commands over the WebSocket control link.
 * @internal
 */
export interface ControlPlane {
  /** Registers interest in receiving streaming frames for a topic over the designated channel. */
  subscribeTopic(
    topic: string,
    dataChannel: string,
    typeName?: string,
    rateLimit?: RateLimitConfig,
  ): Promise<void>;
  /** Terminates an active subscription for a topic on the control plane. */
  unsubscribeTopic(topic: string, typeName?: string): Promise<void>;
  /** Advertises intention to publish structured messages on a topic. */
  advertisePublisher(topic: string, typeName?: string): Promise<void>;
  /** Revokes active publication advertisement for a topic. */
  unadvertisePublisher(topic: string, typeName?: string): Promise<void>;
  /** Sends a local network ICE candidate to the server to assist connection negotiation. */
  sendIceCandidate(candidate: RTCIceCandidate): Promise<void>;
  /** Sends a local session description offer to initiate connection handshaking. */
  sendSdpOffer(sdp: string): Promise<void>;
}

/**
 * Manages the WebSocket connection used for topic registration, publication advertisement, and session handshaking.
 * @internal
 */
export class WebSocketControlPlane implements ControlPlane {
  private socket?: WebSocket;
  private connectPromise?: Promise<void>;
  private signalingReceiver?: SignalingReceiver;
  private readonly diag: DiagnosticRouter;
  private readonly onClose?: () => void;

  constructor(diag: DiagnosticRouter, onClose?: () => void) {
    this.diag = diag;
    this.onClose = onClose;
  }

  /** Registers the designated controller to receive incoming session signaling messages. */
  setSignalingReceiver(receiver: SignalingReceiver): void {
    this.signalingReceiver = receiver;
  }

  /** Opens the WebSocket connection to the server and resolves once established. */
  async connect(
    wsUrl: string,
    timeoutMs = DEFAULT_WS_CONNECT_TIMEOUT_MS,
  ): Promise<void> {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      return;
    }
    if (this.connectPromise) {
      return this.connectPromise;
    }
    this.connectPromise = this.performConnect(wsUrl, timeoutMs);
    return this.connectPromise;
  }

  /** Establishes physical WebSocket link and configures event routing upon open. */
  private async performConnect(
    wsUrl: string,
    timeoutMs: number,
  ): Promise<void> {
    try {
      this.diag.debug(`Connecting WebSocket to ${wsUrl}...`);
      const socket = new WebSocket(wsUrl);
      socket.binaryType = 'arraybuffer';
      this.socket = socket;

      socket.onmessage = (event) => this.handleMessage(event.data);

      await new Promise<void>((resolve, reject) => {
        const timer = setTimeout(() => {
          this.socket = undefined;
          try {
            socket.close();
          } catch (err: unknown) {
            this.diag.debug(
              `Error closing timed-out WebSocket: ${err instanceof Error ? err.message : err}`,
            );
          }
          reject(
            new Error(
              `WebSocket connection to ${wsUrl} timed out after ${timeoutMs}ms. Verify gateway availability and network status.`,
            ),
          );
        }, timeoutMs);

        socket.onclose = (event) => {
          clearTimeout(timer);
          this.socket = undefined;
          const reasonMsg = event.reason ? `: ${event.reason}` : '';
          reject(
            new Error(
              `WebSocket connection to ${wsUrl} closed during handshake (code: ${event.code}${reasonMsg})`,
            ),
          );
        };

        socket.onerror = (err) => {
          clearTimeout(timer);
          this.socket = undefined;
          reject(
            new Error(`WebSocket connection failed to ${wsUrl}`, {cause: err}),
          );
        };

        socket.onopen = () => {
          clearTimeout(timer);
          socket.onerror = (err) => {
            this.diag.error(
              new Error(`WebSocket error on ${wsUrl}`, {cause: err}),
            );
          };
          socket.onclose = () => {
            if (this.onClose) {
              this.onClose();
            }
          };
          this.diag.debug(
            `WebSocket connected successfully (readyState 1). Initializing RTCPeerConnection...`,
          );
          resolve();
        };
      });
    } finally {
      this.connectPromise = undefined;
    }
  }

  /** Terminate the WebSocket control connection and release underlying networking resources. */
  close(): void {
    this.connectPromise = undefined;
    if (this.socket) {
      try {
        this.socket.close();
      } catch (err: unknown) {
        this.diag.debug(
          `Error closing WebSocket: ${err instanceof Error ? err.message : err}`,
        );
      }
      this.socket = undefined;
    }
  }

  /** Serializes a structured control packet into JSON and transmits it across the open WebSocket. */
  private sendPacket(packet: ControlPacket): Promise<void> {
    if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
      return Promise.reject(
        new Error(
          `WebSocket control connection is not open (readyState: ${this.socket ? this.socket.readyState : 'UNINITIALIZED'}). Call connect() before sending control packets.`,
        ),
      );
    }
    const json = toJsonString(ControlPacketSchema, packet);
    this.diag.debug(
      `[WebSocket Tx JSON] Transmitting (${json.length} chars)...`,
    );
    this.socket.send(json);
    return Promise.resolve();
  }

  /** Registers interest in receiving streaming frames for a topic over the designated channel. */
  async subscribeTopic(
    topic: string,
    dataChannel: string,
    typeName?: string,
    rateLimit?: RateLimitConfig,
  ): Promise<void> {
    let rlProto: MessageShape<typeof RateLimitSchema> | undefined;
    if (rateLimit) {
      rlProto = create(RateLimitSchema, {
        messagesPerSecond: rateLimit.messagesPerSecond,
        tokenBucketCapacity: rateLimit.tokenBucketCapacity,
      });
    }

    await this.sendPacket(
      create(ControlPacketSchema, {
        payload: {
          case: 'subscribe',
          value: create(SubscribeSchema, {
            topic,
            type: typeName || undefined,
            dataChannel,
            rateLimit: rlProto,
          }),
        },
      }),
    );
  }

  /** Terminates an active subscription for a topic on the control plane. */
  async unsubscribeTopic(topic: string, typeName?: string): Promise<void> {
    await this.sendPacket(
      create(ControlPacketSchema, {
        payload: {
          case: 'unsubscribe',
          value: create(UnsubscribeSchema, {
            topic,
            type: typeName || undefined,
          }),
        },
      }),
    );
  }

  /** Advertises intention to publish structured messages on a topic. */
  async advertisePublisher(topic: string, typeName?: string): Promise<void> {
    await this.sendPacket(
      create(ControlPacketSchema, {
        payload: {
          case: 'advertisePublisher',
          value: create(AdvertisePublisherSchema, {
            topic,
            type: typeName || undefined,
          }),
        },
      }),
    );
  }

  /** Revokes active publication advertisement for a topic. */
  async unadvertisePublisher(topic: string, typeName?: string): Promise<void> {
    await this.sendPacket(
      create(ControlPacketSchema, {
        payload: {
          case: 'unadvertisePublisher',
          value: create(UnadvertisePublisherSchema, {
            topic,
            type: typeName || undefined,
          }),
        },
      }),
    );
  }

  /** Sends a local network ICE candidate to the server to assist connection negotiation. */
  async sendIceCandidate(candidate: RTCIceCandidate): Promise<void> {
    await this.sendPacket(
      create(ControlPacketSchema, {
        payload: {
          case: 'iceCandidate',
          value: create(RtcIceCandidateSchema, {
            candidate: candidate.candidate,
            sdpMid: candidate.sdpMid ?? undefined,
            sdpMLineIndex: candidate.sdpMLineIndex ?? undefined,
            usernameFragment: candidate.usernameFragment ?? undefined,
          }),
        },
      }),
    );
  }

  /** Sends a local session description offer to initiate connection handshaking. */
  async sendSdpOffer(sdp: string): Promise<void> {
    await this.sendPacket(
      create(ControlPacketSchema, {
        payload: {
          case: 'sessionDescription',
          value: create(RtcSessionDescriptionSchema, {
            type: RtcSessionDescription_SdpType.offer,
            sdp,
          }),
        },
      }),
    );
  }

  /** Parses an incoming WebSocket message, forwarding signaling packets to the registered receiver. */
  private handleMessage(rawData: ArrayBuffer | Uint8Array | string): void {
    let packet: MessageShape<typeof ControlPacketSchema>;
    try {
      let str: string;
      if (typeof rawData === 'string') {
        str = rawData;
      } else {
        const bytes =
          rawData instanceof Uint8Array
            ? rawData
            : // Safe type assertion: rawData is checked not to be a string above, so it is an ArrayBuffer.
              new Uint8Array(rawData as ArrayBuffer);
        str = sharedTextDecoder.decode(bytes);
      }
      // Enable ignoreUnknownFields for forward compatibility when backend services add new schema fields.
      packet = fromJson(ControlPacketSchema, JSON.parse(str), {
        ignoreUnknownFields: true,
      });
    } catch (err: unknown) {
      this.diag.error(
        new Error('Error parsing incoming WebSocket message', {cause: err}),
      );
      return;
    }

    const payload = packet.payload;
    if (!payload || !payload.case) return;
    this.diag.debug(
      `[WebSocket Rx] Received signaling packet: "${payload.case}"`,
    );

    switch (payload.case) {
      case 'sessionDescription': {
        const desc = payload.value;
        const typeStr = this.resolveSdpType(desc.type);
        this.diag.debug(
          `[WebSocket Rx] Setting remote description (${typeStr}, ${desc.sdp?.length || 0} chars)...`,
        );
        if (this.signalingReceiver) {
          this.signalingReceiver.onRemoteSessionDescription(
            typeStr,
            desc.sdp || '',
          );
        }
        break;
      }
      case 'iceCandidate': {
        const cand = payload.value;
        if (cand.candidate && this.signalingReceiver) {
          this.signalingReceiver.onRemoteIceCandidate({
            candidate: cand.candidate,
            sdpMid: cand.sdpMid ?? undefined,
            sdpMLineIndex: cand.sdpMLineIndex ?? undefined,
            usernameFragment: cand.usernameFragment ?? undefined,
          });
        }
        break;
      }
      case 'error': {
        const err = payload.value;
        this.diag.error(
          new PubSubGatewayError(
            err.code ?? GrpcStatus.UNKNOWN,
            err.message ||
              `Control plane error (code: ${err.code ?? GrpcStatus.UNKNOWN})`,
          ),
        );
        break;
      }
    }
  }

  /** Maps raw protobuf SDP enum values into standard browser RTCSdpType strings. */
  private resolveSdpType(type: RtcSessionDescription_SdpType): RTCSdpType {
    // Safe type assertion: casting to unknown allows checking numeric enum values and string variants at runtime.
    const rawType = type as unknown;
    if (
      type === RtcSessionDescription_SdpType.offer ||
      rawType === 1 ||
      rawType === 'offer'
    ) {
      return 'offer';
    }
    if (
      type === RtcSessionDescription_SdpType.pranswer ||
      rawType === 2 ||
      rawType === 'pranswer'
    ) {
      return 'pranswer';
    }
    if (
      type === RtcSessionDescription_SdpType.rollback ||
      rawType === 4 ||
      rawType === 'rollback'
    ) {
      return 'rollback';
    }
    return 'answer';
  }
}
