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

from unittest import mock

from absl.testing import absltest

from intrinsic.logging.proto import pubsub_listener_service_pb2
from intrinsic.solutions import logging_pubsub_listener


class PubsubListeningTest(absltest.TestCase):

  def setUp(self):
    super().setUp()
    self.stub = mock.MagicMock()
    self.ps_listening = logging_pubsub_listener.LoggingPubsubListenerClient(
        self.stub
    )

  def test_get_logging_topics(self):
    test_subscription = pubsub_listener_service_pb2.Subscription(
        topic_expr="test_topic"
    )
    return_value = pubsub_listener_service_pb2.GetSubscriptionResponse()
    return_value.subscription.append(test_subscription)
    self.stub.GetTopicSubscriptions.return_value = return_value

    result = self.ps_listening.get_logging_topics()

    self.assertEqual(
        result[0].topic_expr,
        "test_topic",
    )

  def test_start_logging_topic(self):
    test_subscription = pubsub_listener_service_pb2.Subscription(
        topic_expr="test_topic"
    )

    result = self.ps_listening.start_logging_topic(test_subscription)

    self.assertEqual(self.stub.SetTopicSubscriptions.call_count, 1)
    self.assertEqual(
        type(self.stub.SetTopicSubscriptions.call_args.args[0]),
        pubsub_listener_service_pb2.SetSubscriptionRequest,
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0]
        .subscription[0]
        .topic_expr,
        "test_topic",
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0].allow,
        True,
    )
    self.assertIsNone(result)

  def test_start_logging_topic_with_expression(self):
    result = self.ps_listening.start_logging_topic_with_expression("test_topic")

    self.assertEqual(self.stub.SetTopicSubscriptions.call_count, 1)
    self.assertEqual(
        type(self.stub.SetTopicSubscriptions.call_args.args[0]),
        pubsub_listener_service_pb2.SetSubscriptionRequest,
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0]
        .subscription[0]
        .topic_expr,
        "test_topic",
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0].allow,
        True,
    )
    self.assertIsNone(result)

  def test_start_logging_topics(self):
    test_subscription_1 = pubsub_listener_service_pb2.Subscription(
        topic_expr="test_topic_1"
    )
    test_subscription_2 = pubsub_listener_service_pb2.Subscription(
        topic_expr="test_topic_2"
    )

    result = self.ps_listening.start_logging_topics(
        [test_subscription_1, test_subscription_2]
    )

    self.assertEqual(self.stub.SetTopicSubscriptions.call_count, 1)
    self.assertEqual(
        type(self.stub.SetTopicSubscriptions.call_args.args[0]),
        pubsub_listener_service_pb2.SetSubscriptionRequest,
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0]
        .subscription[0]
        .topic_expr,
        "test_topic_1",
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0]
        .subscription[1]
        .topic_expr,
        "test_topic_2",
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0].allow,
        True,
    )
    self.assertIsNone(result)

  def test_start_logging_topics_with_expressions(self):
    result = self.ps_listening.start_logging_topics_with_expressions(
        ["test_topic_1", "test_topic_2"]
    )

    self.assertEqual(self.stub.SetTopicSubscriptions.call_count, 1)
    self.assertEqual(
        type(self.stub.SetTopicSubscriptions.call_args.args[0]),
        pubsub_listener_service_pb2.SetSubscriptionRequest,
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0]
        .subscription[0]
        .topic_expr,
        "test_topic_1",
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0]
        .subscription[1]
        .topic_expr,
        "test_topic_2",
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0].allow,
        True,
    )
    self.assertIsNone(result)

  def test_stop_logging_topic(self):
    test_subscription = pubsub_listener_service_pb2.Subscription(
        topic_expr="test_topic"
    )

    result = self.ps_listening.stop_logging_topic(test_subscription)

    self.assertEqual(self.stub.SetTopicSubscriptions.call_count, 1)
    self.assertEqual(
        type(self.stub.SetTopicSubscriptions.call_args.args[0]),
        pubsub_listener_service_pb2.SetSubscriptionRequest,
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0]
        .subscription[0]
        .topic_expr,
        "test_topic",
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0].allow,
        False,
    )
    self.assertIsNone(result)

  def test_stop_logging_topic_with_expression(self):
    result = self.ps_listening.stop_logging_topic_with_expression("test_topic")

    self.assertEqual(self.stub.SetTopicSubscriptions.call_count, 1)
    self.assertEqual(
        type(self.stub.SetTopicSubscriptions.call_args.args[0]),
        pubsub_listener_service_pb2.SetSubscriptionRequest,
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0]
        .subscription[0]
        .topic_expr,
        "test_topic",
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0].allow,
        False,
    )
    self.assertIsNone(result)

  def test_stop_logging_topics(self):
    test_subscription_1 = pubsub_listener_service_pb2.Subscription(
        topic_expr="test_topic_1"
    )
    test_subscription_2 = pubsub_listener_service_pb2.Subscription(
        topic_expr="test_topic_2"
    )

    result = self.ps_listening.stop_logging_topics(
        [test_subscription_1, test_subscription_2]
    )

    self.assertEqual(self.stub.SetTopicSubscriptions.call_count, 1)
    self.assertEqual(
        type(self.stub.SetTopicSubscriptions.call_args.args[0]),
        pubsub_listener_service_pb2.SetSubscriptionRequest,
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0]
        .subscription[0]
        .topic_expr,
        "test_topic_1",
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0]
        .subscription[1]
        .topic_expr,
        "test_topic_2",
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0].allow,
        False,
    )
    self.assertIsNone(result)

  def test_stop_logging_topics_with_expressions(self):
    result = self.ps_listening.stop_logging_topics_with_expressions(
        ["test_topic_1", "test_topic_2"]
    )

    self.assertEqual(self.stub.SetTopicSubscriptions.call_count, 1)
    self.assertEqual(
        type(self.stub.SetTopicSubscriptions.call_args.args[0]),
        pubsub_listener_service_pb2.SetSubscriptionRequest,
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0]
        .subscription[0]
        .topic_expr,
        "test_topic_1",
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0]
        .subscription[1]
        .topic_expr,
        "test_topic_2",
    )
    self.assertEqual(
        self.stub.SetTopicSubscriptions.call_args.args[0].allow,
        False,
    )
    self.assertIsNone(result)


if __name__ == "__main__":
  absltest.main()
