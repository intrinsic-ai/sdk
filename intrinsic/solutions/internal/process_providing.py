# Copyright 2023 Intrinsic Innovation LLC

"""Provides behavior trees from a solution."""

import dataclasses
from typing import Iterable, Iterator
import warnings
from google.longrunning import operations_pb2
from google.longrunning import operations_pb2_grpc
from google.protobuf import duration_pb2
from google.rpc import code_pb2
import grpc
from intrinsic.assets import id_utils
from intrinsic.assets.processes.proto import process_asset_pb2
from intrinsic.assets.proto import asset_type_pb2
from intrinsic.assets.proto import id_pb2
from intrinsic.assets.proto import installed_assets_pb2
from intrinsic.assets.proto import installed_assets_pb2_grpc
from intrinsic.assets.proto import metadata_pb2
from intrinsic.assets.proto import view_pb2
from intrinsic.executive.proto import behavior_tree_pb2
from intrinsic.frontend.solution_service.proto import solution_service_pb2
from intrinsic.frontend.solution_service.proto import solution_service_pb2_grpc
from intrinsic.solutions import behavior_tree
from intrinsic.solutions import providers
from intrinsic.util.grpc import error_handling


_INSTALLED_ASSETS_MAX_PAGE_SIZE = 200
_SOLUTION_SERVICE_MAX_PAGE_SIZE = 50
_WAIT_OPERATION_TIMEOUT = duration_pb2.Duration(seconds=10)


@dataclasses.dataclass
class _Process:
  metadata_proto: metadata_pb2.Metadata | None
  behavior_tree_proto: behavior_tree_pb2.BehaviorTree

  def create_behavior_tree(self) -> behavior_tree.BehaviorTree:
    if self.metadata_proto is None:
      return behavior_tree.BehaviorTree.create_from_proto(
          self.behavior_tree_proto
      )
    else:
      return behavior_tree.BehaviorTree.create_from_proto(
          process_asset_pb2.ProcessAsset(
              metadata=self.metadata_proto,
              behavior_tree=self.behavior_tree_proto,
          )
      )


class Processes(providers.ProcessProvider):
  """Provides the processes (= behavior trees) from a solution."""

  _solution: solution_service_pb2_grpc.SolutionServiceStub
  _installed_assets: installed_assets_pb2_grpc.InstalledAssetsStub
  _operations: operations_pb2_grpc.OperationsStub

  def __init__(
      self,
      solution: solution_service_pb2_grpc.SolutionServiceStub,
      installed_assets: installed_assets_pb2_grpc.InstalledAssetsStub,
      operations: operations_pb2_grpc.OperationsStub,
  ):
    self._solution = solution
    self._installed_assets = installed_assets
    self._operations = operations

  def keys(self) -> Iterable[str]:
    return self._list_all_processes(keys_only=True).keys()

  def items(self) -> Iterable[tuple[str, behavior_tree.BehaviorTree]]:
    for id, process in self._list_all_processes(keys_only=False).items():  # pylint: disable=redefined-builtin
      if process is not None:
        yield id, process.create_behavior_tree()

  def values(self) -> Iterable[behavior_tree.BehaviorTree]:
    for process in self._list_all_processes(keys_only=False).values():
      if process is not None:
        yield process.create_behavior_tree()

  def __iter__(self) -> Iterator[str]:
    return self._list_all_processes(keys_only=True).keys().__iter__()

  def __contains__(self, identifier: str) -> bool:
    return self._get_process(identifier) is not None

  def __getitem__(self, identifier: str) -> behavior_tree.BehaviorTree:
    process = self._get_process(identifier)
    if process is None:
      raise KeyError(f'Process "{identifier}" not found')
    return process.create_behavior_tree()

  def __setitem__(self, identifier: str, value: behavior_tree.BehaviorTree):
    warnings.warn(
        '__setitem__ is deprecated. Please use Solution.processes.save()'
        ' instead.',
        DeprecationWarning,
        stacklevel=2,
    )
    if not isinstance(value, behavior_tree.BehaviorTree):
      raise TypeError(f'Expected a BehaviorTree, got {type(value)}.')
    if value.asset_metadata_proto is not None:
      raise ValueError(
          'BehaviorTree represents a Process asset and must be saved using'
          ' Solution.processes.save().'
      )

    # Update the legacy process
    value.name = identifier
    self._save_legacy_process(value)

  def save(self, bt: behavior_tree.BehaviorTree):
    if bt.asset_metadata_proto is None:
      self._save_legacy_process(bt)
    else:
      self._save_process_asset(bt)

  def __delitem__(self, identifier: str):
    try:
      self._solution.DeleteBehaviorTree(
          solution_service_pb2.DeleteBehaviorTreeRequest(name=identifier)
      )
    except Exception as e:
      raise KeyError(f"Failed to delete behavior tree '{identifier}'") from e

  def _get_process(self, identifier: str) -> _Process | None:
    if id_utils.is_id(identifier):
      process_asset = self._get_process_asset(identifier)
      if process_asset is not None:
        return process_asset
    # Fallback: Always try the legacy lookup. Even if it looks like an asset id
    # it can be a behavior tree name like "my_tree.bt.pb".
    return self._get_legacy_process(identifier)

  @error_handling.retry_on_grpc_unavailable
  def _get_process_asset(self, identifier: str) -> _Process | None:
    try:
      response = self._installed_assets.GetInstalledAsset(
          installed_assets_pb2.GetInstalledAssetRequest(
              id=id_pb2.Id(
                  package=id_utils.package_from(identifier),
                  name=id_utils.name_from(identifier),
              ),
              view=view_pb2.AssetViewType.ASSET_VIEW_TYPE_FULL,
          )
      )
      return _Process(
          metadata_proto=response.deployment_data.process.process.metadata,
          behavior_tree_proto=(
              response.deployment_data.process.process.behavior_tree
          ),
      )
    except grpc.RpcError as e:
      if hasattr(e, 'code') and e.code() == grpc.StatusCode.NOT_FOUND:
        return None
      raise e

  @error_handling.retry_on_grpc_unavailable
  def _get_legacy_process(self, identifier: str) -> _Process | None:
    try:
      bt_proto = self._solution.GetBehaviorTree(
          solution_service_pb2.GetBehaviorTreeRequest(name=identifier)
      )
      return _Process(metadata_proto=None, behavior_tree_proto=bt_proto)
    except grpc.RpcError as e:
      if hasattr(e, 'code') and e.code() == grpc.StatusCode.NOT_FOUND:
        return None
      raise e

  # Returns a dict with None values if keys_only is True. This is faster because
  # we need to request less data from the backends.
  def _list_all_processes(
      self,
      *,
      keys_only: bool,
  ) -> dict[str, _Process | None]:
    process_assets = self._list_all_process_assets(keys_only=keys_only)
    legacy_processes = self._list_all_legacy_processes(keys_only=keys_only)

    # "Concatenate" process assets (first) and legacy processes (second). In
    # case of a collision between asset id and legacy behavior tree name the
    # asset takes precedence. Note that we have to do that manually since
    # neither "a | b" nor "b | a" is what we want (one gets the order wrong, the
    # other gets the precedence wrong).
    combined = process_assets
    for k, v in legacy_processes.items():
      if k not in combined:
        combined[k] = v
    return combined

  @error_handling.retry_on_grpc_unavailable
  def _list_all_process_assets(
      self,
      *,
      keys_only: bool,
  ) -> dict[str, _Process | None]:
    view = (
        view_pb2.AssetViewType.ASSET_VIEW_TYPE_BASIC
        if keys_only
        else view_pb2.AssetViewType.ASSET_VIEW_TYPE_FULL
    )
    next_page_token = None

    result: dict[str, _Process | None] = {}
    while True:
      response = self._installed_assets.ListInstalledAssets(
          installed_assets_pb2.ListInstalledAssetsRequest(
              strict_filter=installed_assets_pb2.ListInstalledAssetsRequest.Filter(
                  asset_types=[asset_type_pb2.AssetType.ASSET_TYPE_PROCESS]
              ),
              page_size=_INSTALLED_ASSETS_MAX_PAGE_SIZE,
              page_token=next_page_token,
              view=view,
          )
      )
      for installed_asset in response.installed_assets:
        id_str = id_utils.id_from_proto(installed_asset.metadata.id_version.id)
        if keys_only:
          result[id_str] = None
        else:
          result[id_str] = _Process(
              metadata_proto=(
                  installed_asset.deployment_data.process.process.metadata
              ),
              behavior_tree_proto=(
                  installed_asset.deployment_data.process.process.behavior_tree
              ),
          )
      if not response.next_page_token:
        break
      next_page_token = response.next_page_token

    return result

  @error_handling.retry_on_grpc_unavailable
  def _list_all_legacy_processes(
      self,
      *,
      keys_only: bool,
  ) -> dict[str, _Process | None]:
    view = (
        solution_service_pb2.BehaviorTreeView.BEHAVIOR_TREE_VIEW_BASIC
        if keys_only
        else solution_service_pb2.BehaviorTreeView.BEHAVIOR_TREE_VIEW_FULL
    )
    next_page_token = None

    result: dict[str, _Process | None] = {}
    while True:
      response = self._solution.ListBehaviorTrees(
          solution_service_pb2.ListBehaviorTreesRequest(
              page_size=_SOLUTION_SERVICE_MAX_PAGE_SIZE,
              page_token=next_page_token,
              view=view,
          )
      )
      for bt in response.behavior_trees:
        result[bt.name] = (
            None
            if keys_only
            else _Process(metadata_proto=None, behavior_tree_proto=bt)
        )
      if not response.next_page_token:
        break
      next_page_token = response.next_page_token

    return result

  @error_handling.retry_on_grpc_unavailable
  def _save_legacy_process(self, bt: behavior_tree.BehaviorTree):
    self._solution.UpdateBehaviorTree(
        solution_service_pb2.UpdateBehaviorTreeRequest(
            behavior_tree=bt.proto,
            allow_missing=True,
        )
    )

  @error_handling.retry_on_grpc_unavailable
  def _save_process_asset(self, bt: behavior_tree.BehaviorTree):
    # The installed assets service requires the version and output-only fields
    # to be unset in an installation request.

    # `asset_metadata_proto` returns a copy which we can safely mutate.
    metadata_for_saving = bt.asset_metadata_proto
    assert metadata_for_saving is not None
    metadata_for_saving.id_version.ClearField('version')
    metadata_for_saving.ClearField('file_descriptor_set')
    bt_for_saving = bt.proto
    bt_for_saving.description.ClearField('id_version')

    operation = self._installed_assets.CreateInstalledAsset(
        installed_assets_pb2.CreateInstalledAssetRequest(
            asset=installed_assets_pb2.CreateInstalledAssetRequest.Asset(
                process=process_asset_pb2.ProcessAsset(
                    metadata=metadata_for_saving,
                    behavior_tree=bt_for_saving,
                ),
            ),
        ),
    )

    while not operation.done:
      operation = self._operations.WaitOperation(
          operations_pb2.WaitOperationRequest(
              name=operation.name, timeout=_WAIT_OPERATION_TIMEOUT
          )
      )
      if not operation.done:
        print('Waiting for save operation to finish...')

    if operation.HasField('error'):
      # The installed assets service returns ALREADY_EXISTS if a Process asset
      # with exactly the same content already exists. In this case we consider
      # the save a success.
      if operation.error.code == code_pb2.ALREADY_EXISTS:
        return
      raise RuntimeError(
          'Operation to save Process asset failed with'
          f' {code_pb2.Code.Name(operation.error.code)}:'
          f' {operation.error.message}'
      )

    # Note that the installed asset in the operation response only contains the
    # metadata, not the behavior tree.
    saved_asset = installed_assets_pb2.InstalledAsset()
    operation.response.Unpack(saved_asset)

    # Update the BehaviorTree metadata. Effectively, this only changes the asset
    # version which gets generated by the installed assets service upon
    # installation.
    bt.asset_metadata_proto = saved_asset.metadata
