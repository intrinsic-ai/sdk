# Copyright 2023 Intrinsic Innovation LLC

"""Fake implementation of the DataAssets service."""

from typing import Callable
from typing import NoReturn
from typing import Sequence

import grpc
from grpc.framework.foundation import logging_pool

from intrinsic.assets import id_utils
from intrinsic.assets.data.proto.v1 import data_asset_pb2
from intrinsic.assets.data.proto.v1 import data_assets_pb2
from intrinsic.assets.data.proto.v1 import data_assets_pb2_grpc

_default_page_size = 20

CleanupFunction = Callable[[], None]


class FakeDataAssetsService(data_assets_pb2_grpc.DataAssetsServicer):
  """Fake implementation of the DataAssets service."""

  _data_assets: dict[str, data_asset_pb2.DataAsset]

  @classmethod
  def start_server(
      cls, data_assets: Sequence[data_asset_pb2.DataAsset], port: int = 0
  ) -> tuple[data_assets_pb2_grpc.DataAssetsStub, CleanupFunction]:
    """Creates a FakeDataAssetsService.

    Args:
      data_assets: The data assets to serve.
      port: The port on which to listen, or 0 to pick an unused port.

    Returns:
      A tuple of the DataAssets stub and a cleanup function to close the
      channel and stop the server.
    """
    service = FakeDataAssetsService(data_assets)
    server = grpc.server(
        logging_pool.pool(max_workers=1),
        options=[("grpc.max_receive_message_length", -1)],
    )
    data_assets_pb2_grpc.add_DataAssetsServicer_to_server(service, server)
    port = server.add_secure_port(
        f"[::]:{port}", grpc.local_server_credentials()
    )
    server.start()

    channel = grpc.secure_channel(
        f"localhost:{port}", grpc.local_channel_credentials()
    )
    stub = data_assets_pb2_grpc.DataAssetsStub(channel)

    def cleanup():
      channel.close()
      server.stop(None)

    return stub, cleanup

  def __init__(self, data_assets: Sequence[data_asset_pb2.DataAsset]):
    self._data_assets = {}

    for data_asset in data_assets:
      asset_id = id_utils.id_from_proto(data_asset.metadata.id_version.id)
      if asset_id in self._data_assets:
        raise ValueError(f"Duplicate data asset ID: {asset_id}")

      self._data_assets[asset_id] = data_asset

  def ListDataAssets(
      self,
      request: data_assets_pb2.ListDataAssetsRequest,
      context: grpc.ServicerContext,
  ) -> data_assets_pb2.ListDataAssetsResponse:
    """Lists data assets."""
    filtered_assets = []
    for data_asset in self._data_assets.values():
      if request.strict_filter is not None and request.strict_filter.proto_name:
        type_url = data_asset.data.type_url
        if type_url.split("/")[-1] != request.strict_filter.proto_name:
          continue
      filtered_assets.append(data_asset)

    # Sort by ID for consistent pagination.
    filtered_assets.sort(
        key=lambda asset: id_utils.id_from_proto(asset.metadata.id_version.id)
    )

    # Determine the start of the page.
    offset = 0
    if request.page_token:
      asset_found = False
      for i, asset in enumerate(filtered_assets):
        if (
            id_utils.id_from_proto(asset.metadata.id_version.id)
            == request.page_token
        ):
          asset_found = True
          offset = i
          break
      if not asset_found:
        _abort_with_status(
            context=context,
            code=grpc.StatusCode.INVALID_ARGUMENT,
            message=f"Invalid page token: {request.page_token}",
        )

    page_size = (
        request.page_size if request.page_size > 0 else _default_page_size
    )
    last_index = min(offset + page_size, len(filtered_assets)) - 1

    next_page_token = ""
    if last_index < len(filtered_assets) - 1:
      next_page_token = id_utils.id_from_proto(
          filtered_assets[last_index + 1].metadata.id_version.id
      )

    return data_assets_pb2.ListDataAssetsResponse(
        data_assets=filtered_assets[offset : last_index + 1],
        next_page_token=next_page_token,
    )

  def ListDataAssetMetadata(
      self,
      request: data_assets_pb2.ListDataAssetMetadataRequest,
      context: grpc.ServicerContext,
  ) -> data_assets_pb2.ListDataAssetMetadataResponse:
    """Lists data asset metadata."""
    list_request = data_assets_pb2.ListDataAssetsRequest(
        strict_filter=request.strict_filter,
        page_size=request.page_size,
        page_token=request.page_token,
    )
    list_response = self.ListDataAssets(list_request, context)

    return data_assets_pb2.ListDataAssetMetadataResponse(
        metadata=[asset.metadata for asset in list_response.data_assets],
        next_page_token=list_response.next_page_token,
    )

  def GetDataAsset(
      self,
      request: data_assets_pb2.GetDataAssetRequest,
      context: grpc.ServicerContext,
  ) -> data_asset_pb2.DataAsset:
    """Gets a data asset."""
    try:
      return self._data_assets[id_utils.id_from_proto(request.id)]
    except KeyError as e:
      _abort_with_status(
          context=context,
          code=grpc.StatusCode.NOT_FOUND,
          message=f"Data asset not found: {e}",
      )

  def StreamReferencedData(
      self,
      request: data_assets_pb2.StreamReferencedDataRequest,
      context: grpc.ServicerContext,
  ) -> data_assets_pb2.StreamReferencedDataResponse:
    """Streams referenced data."""
    del request  # Unused.
    _abort_with_status(
        context=context,
        code=grpc.StatusCode.UNIMPLEMENTED,
        message=(
            "StreamReferencedData is not implemented in FakeDataAssetsService."
        ),
    )


def _abort_with_status(
    context: grpc.ServicerContext,
    code: grpc.StatusCode,
    message: str,
) -> NoReturn:
  context.abort(code, message)

  # This will never be raised, but we need it to satisfy static type checking,
  raise AssertionError("This error should not have been raised.")
