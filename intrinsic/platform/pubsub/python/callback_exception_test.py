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

_LIVELINESS_KEYEXPR = 'key1'


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

  def test_exception_in_liveliness_subscription_callback(self):
    self.pubsub = pubsub.PubSub()
    event = threading.Event()

    self.pubsub.DeclareLivelinessToken(_LIVELINESS_KEYEXPR)

    def msg_callback(key: str, alive: bool):
      logging.info(
          'Liveliness callback called with key=%s, alive=%s', key, alive
      )
      event.set()
      raise RuntimeError('This should not cause any issues')

    sub = self.pubsub.CreateLivelinessSubscription(
        _LIVELINESS_KEYEXPR,
        True,
        msg_callback,
    )

    self.assertTrue(event.wait(timeout=5.0), 'Callback was not called in time')
    logging.info('Unsubscribing and dropping the liveliness token')
    sub.Unsubscribe()
    self.pubsub.DropLivelinessToken(_LIVELINESS_KEYEXPR)

  def test_exception_in_liveliness_get_callback(self):
    self.pubsub = pubsub.PubSub()
    event = threading.Event()

    self.pubsub.DeclareLivelinessToken(_LIVELINESS_KEYEXPR)

    def callback(key: str):
      logging.info('Liveliness get callback called with key=%s', key)
      raise RuntimeError('This should not cause any issues')

    def on_done(keyexpr: str):
      logging.info(
          'Liveliness on_done callback called with keyexpr=%s', keyexpr
      )
      event.set()

    self.pubsub.LivelinessGet(
        _LIVELINESS_KEYEXPR,
        callback=callback,
        on_done=on_done,
    )

    self.assertTrue(
        event.wait(timeout=5.0), 'on_done callback was not called in time'
    )
    logging.info('Dropping the liveliness token')
    self.pubsub.DropLivelinessToken(_LIVELINESS_KEYEXPR)

  def test_exception_in_liveliness_on_done_callback(self):
    self.pubsub = pubsub.PubSub()
    event = threading.Event()

    self.pubsub.DeclareLivelinessToken(_LIVELINESS_KEYEXPR)

    def callback(key: str):
      logging.info('Liveliness get callback called with key=%s', key)

    def on_done(keyexpr: str):
      logging.info(
          'Liveliness on_done callback called with keyexpr=%s', keyexpr
      )
      event.set()
      raise RuntimeError('This should not cause any issues')

    self.pubsub.LivelinessGet(
        _LIVELINESS_KEYEXPR,
        callback=callback,
        on_done=on_done,
    )

    self.assertTrue(
        event.wait(timeout=5.0), 'on_done callback was not called in time'
    )
    logging.info('Dropping the liveliness token')
    self.pubsub.DropLivelinessToken(_LIVELINESS_KEYEXPR)


if __name__ == '__main__':
  absltest.main()
