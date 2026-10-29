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

import type {ControlPlane, SignalingReceiver} from './control_plane';
import type {DiagnosticRouter} from './diagnostics';
import type {ConnectionState} from './types';

/** Timeout duration for WebRTC ICE candidate negotiation and connection establishment (15 seconds). */
const ICE_NEGOTIATION_TIMEOUT_MS = 15_000;

/**
 * Orchestrates real-time session negotiation and translates network connectivity changes into public state events.
 * @internal
 */
export class SignalingController implements SignalingReceiver {
  private pc?: RTCPeerConnection;
  private signalingQueue: Promise<void> = Promise.resolve();
  private readonly controlPlane: ControlPlane;
  private readonly diag: DiagnosticRouter;
  private readonly onStateChange: (state: ConnectionState) => void;
  private currentIceState?: RTCIceConnectionState;
  private pendingReject?: (err: Error) => void;

  constructor(
    controlPlane: ControlPlane,
    diag: DiagnosticRouter,
    onStateChange: (state: ConnectionState) => void,
  ) {
    this.controlPlane = controlPlane;
    this.diag = diag;
    this.onStateChange = onStateChange;
  }

  /** Initiates session description negotiation and resolves once peer connectivity is established. */
  async negotiate(
    pc: RTCPeerConnection,
    ensureDefaultChannel: () => Promise<unknown>,
  ): Promise<void> {
    this.pc = pc;
    let stateCallback: ((state: RTCIceConnectionState) => void) | undefined;

    pc.onicecandidate = (event) => {
      if (event.candidate && event.candidate.candidate !== '') {
        this.controlPlane.sendIceCandidate(event.candidate).catch((err) => {
          this.diag.warn(
            `Failed to send ICE candidate: ${err instanceof Error ? err.message : err}`,
          );
        });
      }
    };

    pc.oniceconnectionstatechange = () => {
      const state = pc.iceConnectionState;
      this.currentIceState = state;
      this.emitUserStatus(state);
      if (stateCallback) {
        stateCallback(state);
      }
    };

    this.diag.debug(
      `Generating WebRTC SDP Offer with initial data channel "default"...`,
    );
    const pendingChannel = ensureDefaultChannel();

    const offer = await pc.createOffer();
    if (!this.pc) {
      throw new Error('SignalingController was closed during negotiation');
    }
    await pc.setLocalDescription(offer);
    if (!this.pc) {
      throw new Error('SignalingController was closed during negotiation');
    }
    this.diag.debug(
      `SDP Offer generated (${offer.sdp?.length || 0} chars). Sending sessionDescription packet over WebSocket...`,
    );

    await this.controlPlane.sendSdpOffer(offer.sdp || '');
    if (!this.pc) {
      throw new Error('SignalingController was closed during negotiation');
    }
    this.diag.debug(
      `SDP Offer transmitted! Awaiting remote SDP Answer and ICE connectivity...`,
    );

    await new Promise<void>((resolve, reject) => {
      if (!this.pc) {
        return reject(
          new Error('SignalingController was closed during negotiation'),
        );
      }
      if (
        this.currentIceState === 'connected' ||
        this.currentIceState === 'completed'
      ) {
        return resolve();
      }

      const timer = setTimeout(() => {
        stateCallback = undefined;
        this.pendingReject = undefined;
        reject(
          new Error(
            `WebRTC ICE negotiation timed out after ${ICE_NEGOTIATION_TIMEOUT_MS}ms. Verify STUN/TURN server configuration and network candidate reachability.`,
          ),
        );
      }, ICE_NEGOTIATION_TIMEOUT_MS);

      this.pendingReject = (err: Error) => {
        clearTimeout(timer);
        stateCallback = undefined;
        this.pendingReject = undefined;
        reject(err);
      };

      stateCallback = (state: RTCIceConnectionState) => {
        if (state === 'connected' || state === 'completed') {
          clearTimeout(timer);
          stateCallback = undefined;
          this.pendingReject = undefined;
          resolve();
        } else if (state === 'failed' || state === 'closed') {
          clearTimeout(timer);
          stateCallback = undefined;
          this.pendingReject = undefined;
          reject(new Error(`ICE connection failed with state "${state}"`));
        }
      };
    });

    if (!this.pc) {
      throw new Error('SignalingController was closed during negotiation');
    }
    await pendingChannel;
  }

  /** Releases peer connection references and ceases further signaling event processing. */
  close(): void {
    const reject = this.pendingReject;
    this.pendingReject = undefined;
    this.pc = undefined;
    if (reject) {
      reject(new Error('SignalingController was closed during negotiation'));
    }
  }

  /** Translates internal WebRTC ICE connectivity states into high-level user status events. */
  private emitUserStatus(state: RTCIceConnectionState): void {
    let status: ConnectionState = 'disconnected';
    switch (state) {
      case 'new':
      case 'checking':
        status = 'connecting';
        break;
      case 'connected':
      case 'completed':
        status = 'connected';
        break;
      case 'disconnected':
        status = 'reconnecting';
        break;
      case 'failed':
      case 'closed':
        status = 'failed';
        break;
    }
    this.onStateChange(status);
  }

  /** Applies a remote session description answer received from the signaling server. */
  onRemoteSessionDescription(type: RTCSdpType, sdp: string): void {
    if ((type === 'answer' || type === 'pranswer') && this.pc) {
      const desc: RTCSessionDescriptionInit = {type, sdp};
      this.signalingQueue = this.signalingQueue.then(async () => {
        if (!this.pc) {
          return;
        }
        try {
          await this.pc.setRemoteDescription(desc);
          this.diag.debug(
            `[WebSocket Rx] Successfully applied remote SDP description!`,
          );
        } catch (err: unknown) {
          this.diag.error(
            new Error('Failed to set remote SDP description', {cause: err}),
          );
        }
      });
    }
  }

  /** Registers a remote network ICE candidate received from the signaling server. */
  onRemoteIceCandidate(candidate: RTCIceCandidateInit): void {
    if (candidate.candidate && this.pc) {
      this.signalingQueue = this.signalingQueue.then(async () => {
        if (!this.pc) {
          return;
        }
        try {
          await this.pc.addIceCandidate(candidate);
        } catch (err: unknown) {
          this.diag.warn(
            `Failed to add remote ICE candidate: ${err instanceof Error ? err.message : String(err)}. Verify network connectivity.`,
          );
        }
      });
    }
  }
}
