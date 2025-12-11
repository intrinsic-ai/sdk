# Copyright 2023 Intrinsic Innovation LLC

import threading
from absl.testing import absltest
from absl import logging

from intrinsic.platform.pubsub.python import pubsub
from intrinsic.platform.common.proto import test_pb2


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
