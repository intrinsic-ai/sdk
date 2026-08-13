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

"""Provides a client for using the AssetConfigurationService."""

from __future__ import annotations

import warnings

from google.protobuf import any_pb2
import grpc

from intrinsic.assets.proto.v1 import asset_configuration_pb2
from intrinsic.assets.proto.v1 import asset_configuration_pb2_grpc
from intrinsic.util.grpc import error_handling


class AssetConfigurationClient:
  """Client for the AssetConfigurationService."""

  _stub: asset_configuration_pb2_grpc.AssetConfigurationServiceStub

  def __init__(
      self, stub: asset_configuration_pb2_grpc.AssetConfigurationServiceStub
  ):
    self._stub = stub

  @classmethod
  def from_channel(cls, grpc_channel: grpc.Channel) -> AssetConfigurationClient:
    return cls(
        asset_configuration_pb2_grpc.AssetConfigurationServiceStub(grpc_channel)
    )

  def recommend_asset_configuration(
      self,
      name: str,
      input_configuration: any_pb2.Any | None = None,
  ) -> asset_configuration_pb2.RecommendAssetConfigurationResponse:
    request = asset_configuration_pb2.RecommendAssetConfigurationRequest(
        name=name, input_configuration=input_configuration
    )
    try:
      return self._stub.RecommendAssetConfiguration(request)
    except grpc.RpcError as e:
      if error_handling.is_unavailable_grpc_status(e):
        warnings.warn(
            "Failed to get asset recommendation for asset: "
            f"{name}. Returning input configuration instead.",
            RuntimeWarning,
        )
        return asset_configuration_pb2.RecommendAssetConfigurationResponse(
            config=input_configuration
        )
      raise

  def get_asset_recommendation_info(
      self, name: str
  ) -> asset_configuration_pb2.AssetRecommendationInfo:
    request = asset_configuration_pb2.GetAssetRecommendationInfoRequest(
        name=name
    )
    try:
      return self._stub.GetAssetRecommendationInfo(request)
    except grpc.RpcError as e:
      if error_handling.is_unavailable_grpc_status(e):
        warnings.warn(
            f"Failed to get asset recommendation info for asset: {name}",
            RuntimeWarning,
        )
        return asset_configuration_pb2.AssetRecommendationInfo(
            name=name, has_recommendation=False
        )
      raise
