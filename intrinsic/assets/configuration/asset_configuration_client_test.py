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

"""Tests of AssetConfigurationClient."""

from unittest import mock

from absl.testing import absltest
from google.protobuf import any_pb2
import grpc

from intrinsic.assets.configuration import asset_configuration_client
from intrinsic.assets.proto.v1 import asset_configuration_pb2


class AssetConfigurationClientTest(absltest.TestCase):

  _stub: mock.MagicMock
  _client: asset_configuration_client.AssetConfigurationClient

  def setUp(self):
    super().setUp()
    self._stub = mock.MagicMock()
    self._client = asset_configuration_client.AssetConfigurationClient(
        self._stub
    )

  def test_init(self):
    self.assertEqual(self._client._stub, self._stub)

  def test_recommend_asset_configuration(self):
    expected_response = (
        asset_configuration_pb2.RecommendAssetConfigurationResponse()
    )
    self._stub.RecommendAssetConfiguration.return_value = expected_response

    input_config = any_pb2.Any()

    response = self._client.recommend_asset_configuration(
        name="test_asset", input_configuration=input_config
    )

    self.assertEqual(response, expected_response)
    self._stub.RecommendAssetConfiguration.assert_called_once_with(
        asset_configuration_pb2.RecommendAssetConfigurationRequest(
            name="test_asset", input_configuration=input_config
        )
    )

  def test_get_asset_recommendation_info(self):
    expected_response = asset_configuration_pb2.AssetRecommendationInfo(
        name="test_asset", has_recommendation=True
    )
    self._stub.GetAssetRecommendationInfo.return_value = expected_response

    response = self._client.get_asset_recommendation_info(name="test_asset")

    self.assertEqual(response, expected_response)
    self._stub.GetAssetRecommendationInfo.assert_called_once_with(
        asset_configuration_pb2.GetAssetRecommendationInfoRequest(
            name="test_asset"
        )
    )

  def test_from_channel(self):
    mock_channel = mock.MagicMock(spec=grpc.Channel)
    client = asset_configuration_client.AssetConfigurationClient.from_channel(
        mock_channel
    )
    self.assertIsInstance(
        client, asset_configuration_client.AssetConfigurationClient
    )
    self.assertIsNotNone(client._stub)


if __name__ == "__main__":
  absltest.main()
