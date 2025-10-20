# Copyright 2023 Intrinsic Innovation LLC

"""Provides behavior trees from a solution."""

from typing import Iterable, Iterator
import grpc
from intrinsic.assets import id_utils
from intrinsic.assets.processes.proto import process_asset_pb2
from intrinsic.assets.proto import asset_type_pb2
from intrinsic.assets.proto import id_pb2
from intrinsic.assets.proto import installed_assets_pb2
from intrinsic.assets.proto import installed_assets_pb2_grpc
from intrinsic.assets.proto import view_pb2
from intrinsic.executive.proto import behavior_tree_pb2
from intrinsic.frontend.solution_service.proto import solution_service_pb2
from intrinsic.frontend.solution_service.proto import solution_service_pb2_grpc
from intrinsic.solutions import behavior_tree
from intrinsic.solutions import providers


_INSTALLED_ASSETS_MAX_PAGE_SIZE = 200
_SOLUTION_SERVICE_MAX_PAGE_SIZE = 50


class Processes(providers.ProcessProvider):
  """Provides the processes (= behavior trees) from a solution."""

  _solution: solution_service_pb2_grpc.SolutionServiceStub
  _installed_assets: installed_assets_pb2_grpc.InstalledAssetsStub

  def __init__(
      self,
      solution: solution_service_pb2_grpc.SolutionServiceStub,
      installed_assets: installed_assets_pb2_grpc.InstalledAssetsStub,
  ):
    self._solution = solution
    self._installed_assets = installed_assets

  def keys(self) -> Iterable[str]:
    return self._list_all_processes(keys_only=True).keys()

  def items(self) -> Iterable[tuple[str, behavior_tree.BehaviorTree]]:
    for k, v in self._list_all_processes(keys_only=False).items():
      if v is not None:
        yield k, behavior_tree.BehaviorTree(bt=v.behavior_tree)

  def values(self) -> Iterable[behavior_tree.BehaviorTree]:
    for v in self._list_all_processes(keys_only=False).values():
      if v is not None:
        yield behavior_tree.BehaviorTree(bt=v.behavior_tree)

  def __iter__(self) -> Iterator[str]:
    return self._list_all_processes(keys_only=True).keys().__iter__()

  def __contains__(self, identifier: str) -> bool:
    return self._get_process(identifier) is not None

  def __getitem__(self, identifier: str) -> behavior_tree.BehaviorTree:
    process = self._get_process(identifier)
    if process is None:
      raise KeyError(f'Process "{identifier}" not found')
    return behavior_tree.BehaviorTree(bt=process.behavior_tree)

  def __setitem__(self, identifier: str, value: behavior_tree.BehaviorTree):
    value.name = identifier
    self._solution.UpdateBehaviorTree(
        solution_service_pb2.UpdateBehaviorTreeRequest(
            behavior_tree=value.proto,
            allow_missing=True,
        )
    )

  def __delitem__(self, identifier: str):
    try:
      self._solution.DeleteBehaviorTree(
          solution_service_pb2.DeleteBehaviorTreeRequest(name=identifier)
      )
    except Exception as e:
      raise KeyError(f"Failed to delete behavior tree '{identifier}'") from e

  def _get_process(
      self, identifier: str
  ) -> process_asset_pb2.ProcessAsset | None:
    if id_utils.is_id(identifier):
      process_asset = self._get_process_asset(identifier)
      if process_asset is not None:
        return process_asset
    # Fallback: Always try the legacy lookup. Even if it looks like an asset id
    # it can be a behavior tree name like "my_tree.bt.pb".
    tree = self._get_legacy_process(identifier)
    if tree is not None:
      return process_asset_pb2.ProcessAsset(behavior_tree=tree)
    return None

  def _get_process_asset(
      self, identifier: str
  ) -> process_asset_pb2.ProcessAsset:
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
      return response.deployment_data.process.process
    except grpc.RpcError as e:
      if hasattr(e, 'code') and e.code() == grpc.StatusCode.NOT_FOUND:
        return None
      raise e

  def _get_legacy_process(
      self, identifier: str
  ) -> behavior_tree_pb2.BehaviorTree | None:
    try:
      return self._solution.GetBehaviorTree(
          solution_service_pb2.GetBehaviorTreeRequest(name=identifier)
      )
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
  ) -> dict[str, process_asset_pb2.ProcessAsset | None]:
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

  def _list_all_process_assets(
      self,
      *,
      keys_only: bool,
  ) -> dict[str, process_asset_pb2.ProcessAsset | None]:
    view = (
        view_pb2.AssetViewType.ASSET_VIEW_TYPE_BASIC
        if keys_only
        else view_pb2.AssetViewType.ASSET_VIEW_TYPE_FULL
    )
    next_page_token = None

    result: dict[str, process_asset_pb2.ProcessAsset | None] = {}
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
        result[id_str] = (
            None
            if keys_only
            else installed_asset.deployment_data.process.process
        )
      if not response.next_page_token:
        break
      next_page_token = response.next_page_token

    return result

  def _list_all_legacy_processes(
      self,
      *,
      keys_only: bool,
  ) -> dict[str, process_asset_pb2.ProcessAsset | None]:
    view = (
        solution_service_pb2.BehaviorTreeView.BEHAVIOR_TREE_VIEW_BASIC
        if keys_only
        else solution_service_pb2.BehaviorTreeView.BEHAVIOR_TREE_VIEW_FULL
    )
    next_page_token = None

    result: dict[str, process_asset_pb2.ProcessAsset | None] = {}
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
            else process_asset_pb2.ProcessAsset(behavior_tree=bt)
        )
      if not response.next_page_token:
        break
      next_page_token = response.next_page_token

    return result
