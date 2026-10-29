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
  type MessageShape,
  type Registry,
  createRegistry,
} from '@bufbuild/protobuf';
import {anyPack} from '@bufbuild/protobuf/wkt';

import type {ControlPlane} from './control_plane';
import type {DataPlane, TopicSender} from './data_plane';
import type {Publisher} from './types';

/**
 * Active sender handle that transmits structured messages to a topic.
 */
export class PublisherHandle<Desc extends DescMessage = DescMessage>
  implements Publisher<Desc>, TopicSender
{
  private state: 'unadvertised' | 'advertised' | 'closed' = 'unadvertised';
  private advertisePromise?: Promise<void>;

  private readonly controlPlane: ControlPlane;
  private readonly dataPlane: DataPlane;
  readonly topic: string;
  readonly schema: Desc;
  private readonly registry: Registry;
  readonly dataChannel: string;

  constructor(
    controlPlane: ControlPlane,
    dataPlane: DataPlane,
    topic: string,
    schema: Desc,
    dataChannel: string,
  ) {
    this.controlPlane = controlPlane;
    this.dataPlane = dataPlane;
    this.topic = topic;
    this.schema = schema;
    this.registry = createRegistry(schema);
    this.dataChannel = dataChannel;
  }

  /** Verifies the handle is still open, throwing an informative error if already closed. */
  private ensureOpen(action: 'advertise publisher for' | 'publish to'): void {
    if (this.state === 'closed') {
      throw new Error(
        `Cannot ${action} topic "${this.topic}": publisher handle is closed. Call createPublisher() to create a new publisher instance.`,
      );
    }
  }

  /**
   * Registers publication intent with the server before streaming data frames.
   * @throws {Error} If the publisher handle is closed during negotiation.
   */
  async advertise(): Promise<void> {
    this.ensureOpen('advertise publisher for');
    if (this.state === 'advertised') {
      return;
    }
    if (this.advertisePromise) {
      return this.advertisePromise;
    }
    this.advertisePromise = this.performAdvertise();
    return this.advertisePromise;
  }

  /** Registers sender references and sends an advertise command over the control link. */
  private async performAdvertise(): Promise<void> {
    try {
      this.dataPlane.registerSender(this.topic, this);
      await this.controlPlane.advertisePublisher(
        this.topic,
        this.schema.typeName,
      );
      if (this.state === 'closed') {
        try {
          await this.controlPlane.unadvertisePublisher(
            this.topic,
            this.schema.typeName,
          );
        } catch (err: unknown) {
          console.warn(
            `Failed to unadvertise topic "${this.topic}" during rollback:`,
            err,
          );
        }
        throw new Error(
          `Publisher for topic "${this.topic}" was closed while advertising`,
        );
      }
      this.state = 'advertised';
    } catch (err: unknown) {
      this.state = 'closed';
      this.dataPlane.unregisterSender(this.topic);
      throw err;
    } finally {
      this.advertisePromise = undefined;
    }
  }

  /**
   * Transmits a structured message payload to all subscribers on the topic.
   *
   * @param payload The message content matching this handle's protobuf schema.
   * @throws {Error} If this publisher handle has been closed.
   */
  async publish(payload: MessageShape<Desc>): Promise<void> {
    this.ensureOpen('publish to');
    if (this.state !== 'advertised') {
      await this.advertise();
    }

    const anyPayload = anyPack(this.schema, payload);

    await this.dataPlane.publish(
      this.dataChannel,
      this.topic,
      anyPayload,
      this.registry,
    );
  }

  /**
   * Stops publication and deregisters this sender handle from the topic.
   */
  async close(): Promise<void> {
    if (this.state === 'closed') {
      return;
    }
    const wasAdvertised = this.state === 'advertised';
    this.state = 'closed';
    this.dataPlane.unregisterSender(this.topic);
    if (this.advertisePromise) {
      try {
        await this.advertisePromise;
      } catch (err: unknown) {
        console.warn(
          `Error awaiting advertising resolution during close for topic "${this.topic}":`,
          err,
        );
      }
      return;
    }
    try {
      if (wasAdvertised) {
        await this.controlPlane.unadvertisePublisher(
          this.topic,
          this.schema.typeName,
        );
      }
    } catch (err: unknown) {
      console.warn(
        `Error while sending unadvertise for topic "${this.topic}":`,
        err,
      );
    }
  }
}
