# Copyright 2023 Intrinsic Innovation LLC

"""Provides a client for using the InstalledAssets service.

The provided client class is for use in the Solution Building Library.
"""

from __future__ import annotations

from google.longrunning import operations_pb2
from google.longrunning import operations_pb2_grpc
from google.protobuf import duration_pb2
from google.rpc import status_pb2
import grpc
from grpc_status import rpc_status
from intrinsic.assets import id_utils
from intrinsic.assets.proto import asset_type_pb2
from intrinsic.assets.proto import id_pb2
from intrinsic.assets.proto import installed_assets_pb2
from intrinsic.assets.proto import installed_assets_pb2_grpc
from intrinsic.assets.proto import view_pb2
from intrinsic.util.grpc import error_handling

_MAX_PAGE_SIZE = 200
_WAIT_OPERATION_TIMEOUT = duration_pb2.Duration(seconds=10)


def _to_id_proto(id: str | id_pb2.Id) -> id_pb2.Id:  # pylint: disable=redefined-builtin
  """Converts a string or id_pb2.Id to an id_pb2.Id proto."""
  if isinstance(id, str):
    return id_pb2.Id(
        package=id_utils.package_from(id), name=id_utils.name_from(id)
    )
  elif isinstance(id, id_pb2.Id):
    return id
  else:
    raise TypeError(f'Unsupported type for id: {type(id)}')


class OperationError(RuntimeError):
  """An error resulting from a long-running installed assets operation."""

  _code: grpc.StatusCode
  _status_proto: status_pb2.Status

  def __init__(self, status_proto: status_pb2.Status):
    status = rpc_status.to_status(status_proto)
    self._code = status.code
    self._status_proto = status_proto

  def __str__(self) -> str:
    return (
        f'<OperationError:\n\tcode = {self._code}\n\tdetails ='
        f' {self._status_proto.message}\n>'
    )

  def code(self) -> grpc.StatusCode:
    return self._code

  def message(self) -> str:
    return self._status_proto.message


class InstalledAssetsClient:
  """Client for the InstalledAssets service."""

  _installed_assets: installed_assets_pb2_grpc.InstalledAssetsStub
  _operations: operations_pb2_grpc.OperationsStub

  def __init__(
      self,
      installed_assets: installed_assets_pb2_grpc.InstalledAssetsStub,
      operations: operations_pb2_grpc.OperationsStub,
  ):
    """Constructs a new InstalledAssetsClient object.

    Args:
      installed_assets: The gRPC stub to be used for communication with the
        installed assets service.
      operations: The gRPC stub to be used for communication with the operations
        service.
    """
    self._installed_assets = installed_assets
    self._operations = operations

  @classmethod
  def from_channel(cls, grpc_channel: grpc.Channel) -> InstalledAssetsClient:
    """Create a new InstalledAssetsClient from a gRPC channel.

    Args:
      grpc_channel: Channel to the InstalledAssets service.

    Returns:
      A newly created instance of the InstalledAssetsClient class.
    """
    return cls(
        installed_assets_pb2_grpc.InstalledAssetsStub(grpc_channel),
        operations_pb2_grpc.OperationsStub(grpc_channel),
    )

  def get_installed_asset(
      self,
      id: str | id_pb2.Id,  # pylint: disable=redefined-builtin
      view: view_pb2.AssetViewType = view_pb2.AssetViewType.ASSET_VIEW_TYPE_FULL,
  ) -> installed_assets_pb2.InstalledAsset:
    """Calls the GetInstalledAsset method of the InstalledAssets service."""
    return self._get_installed_asset_with_retry(
        installed_assets_pb2.GetInstalledAssetRequest(
            id=_to_id_proto(id), view=view
        )
    )

  def list_all_installed_assets(
      self,
      asset_types: list[asset_type_pb2.AssetType] | None = None,
      view: view_pb2.AssetViewType = view_pb2.AssetViewType.ASSET_VIEW_TYPE_BASIC,
  ) -> list[installed_assets_pb2.InstalledAsset]:
    """Calls the ListInstalledAssets method of the InstalledAssets service.

    Handles pagination and returns the combined result of all pages in a single
    list.

    Args:
      asset_types: The asset types to filter by.
      view: The view of the assets to return.

    Returns:
      The list of all installed assets matching the filter.
    """

    next_page_token = None

    result: list[installed_assets_pb2.InstalledAsset] = []
    while True:
      response = self._list_installed_assets_with_retry(
          installed_assets_pb2.ListInstalledAssetsRequest(
              strict_filter=installed_assets_pb2.ListInstalledAssetsRequest.Filter(
                  asset_types=asset_types
              ),
              page_size=_MAX_PAGE_SIZE,
              page_token=next_page_token,
              view=view,
          )
      )
      result.extend(response.installed_assets)
      if not response.next_page_token:
        break
      next_page_token = response.next_page_token

    return result

  def create_installed_asset(
      self,
      asset: installed_assets_pb2.CreateInstalledAssetRequest.Asset,
      update_policy: installed_assets_pb2.UpdatePolicy = installed_assets_pb2.UpdatePolicy.UPDATE_POLICY_UNSPECIFIED,
  ) -> installed_assets_pb2.InstalledAsset:
    """Calls the CreateInstalledAsset method of the InstalledAssets service.

    Blocks until the installation operation has completed.

    Args:
      asset: The asset to create.
      update_policy: The update policy to use.

    Returns:
      The installed asset that was created. Note that in this asset only the
      metadata is set.

    Raises:
      OperationError: If the installation operation fails.
    """
    operation = self._create_installed_asset_with_retry(
        installed_assets_pb2.CreateInstalledAssetRequest(
            asset=asset, policy=update_policy
        ),
    )

    while not operation.done:
      operation = self._wait_operation_with_retry(
          operations_pb2.WaitOperationRequest(
              name=operation.name, timeout=_WAIT_OPERATION_TIMEOUT
          )
      )
      if not operation.done:
        print('Waiting for operation to create installed asset to finish...')

    if operation.HasField('error'):
      raise OperationError(operation.error)

    installed_asset = installed_assets_pb2.InstalledAsset()
    operation.response.Unpack(installed_asset)
    return installed_asset

  def delete_installed_asset(
      self,
      id: str | id_pb2.Id,  # pylint: disable=redefined-builtin
      delete_policy: installed_assets_pb2.DeletePolicy = installed_assets_pb2.DeletePolicy.POLICY_UNSPECIFIED,
  ):
    """Calls the DeleteInstalledAsset method of the InstalledAssets service.

    Blocks until the deletion operation has completed.

    Args:
      id: The id of the asset to delete.
      delete_policy: The delete policy to use.

    Raises:
      OperationError: If the installation operation fails.
    """
    operation = self._delete_installed_asset_with_retry(
        installed_assets_pb2.DeleteInstalledAssetRequest(
            asset=_to_id_proto(id), policy=delete_policy
        )
    )

    while not operation.done:
      operation = self._wait_operation_with_retry(
          operations_pb2.WaitOperationRequest(
              name=operation.name, timeout=_WAIT_OPERATION_TIMEOUT
          )
      )
      if not operation.done:
        print('Waiting for operation to delete installed asset to finish...')

    if operation.HasField('error'):
      raise OperationError(operation.error)

  @error_handling.retry_on_grpc_unavailable
  def _get_installed_asset_with_retry(
      self,
      request: installed_assets_pb2.GetInstalledAssetRequest,
  ) -> installed_assets_pb2.InstalledAsset:
    return self._installed_assets.GetInstalledAsset(request)

  @error_handling.retry_on_grpc_unavailable
  def _list_installed_assets_with_retry(
      self,
      request: installed_assets_pb2.ListInstalledAssetsRequest,
  ) -> installed_assets_pb2.ListInstalledAssetsResponse:
    return self._installed_assets.ListInstalledAssets(request)

  @error_handling.retry_on_grpc_unavailable
  def _create_installed_asset_with_retry(
      self,
      request: installed_assets_pb2.CreateInstalledAssetRequest,
  ) -> operations_pb2.Operation:
    return self._installed_assets.CreateInstalledAsset(request)

  @error_handling.retry_on_grpc_unavailable
  def _delete_installed_asset_with_retry(
      self,
      request: installed_assets_pb2.DeleteInstalledAssetRequest,
  ) -> operations_pb2.Operation:
    return self._installed_assets.DeleteInstalledAsset(request)

  @error_handling.retry_on_grpc_unavailable
  def _wait_operation_with_retry(
      self,
      request: operations_pb2.WaitOperationRequest,
  ) -> operations_pb2.Operation:
    return self._operations.WaitOperation(request)
