# Copyright 2026 Intrinsic Innovation LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import threading

from absl import logging
from absl.testing import absltest

from intrinsic.platform.common.proto import test_pb2
from intrinsic.platform.pubsub.python import pubsub


class SegfaultReproductionTest(absltest.TestCase):

  def test_callback_exception(self):
    self.pubsub = pubsub.PubSub()
    config = pubsub.TopicConfig()
    self.pub = self.pubsub.CreatePublisher('segfault_topic', config)
    event = threading.Event()

    def msg_callback(message):
      logging.info('Callback called, raising exception...')
      event.set()
      raise RuntimeError('This should not cause a segfault!')

    self.sub = self.pubsub.CreateSubscription(
        topic='segfault_topic',
        config=config,
        exemplar=test_pb2.TestMessageString(),
        msg_callback=msg_callback,
    )

    # Publish a message to trigger the callback
    self.pub.Publish(test_pb2.TestMessageString(data='trigger'))

    # Wait for the callback to execute
    self.assertTrue(event.wait(timeout=5.0), 'Callback was not called in time')


if __name__ == '__main__':
  absltest.main()
