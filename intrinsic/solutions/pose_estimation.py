# Copyright 2023 Intrinsic Innovation LLC

"""Pose estimator access within the workcell API."""

import dataclasses
import datetime
from typing import Dict, Iterator, List, Optional

from intrinsic.assets import interface_utils
from intrinsic.assets.proto import asset_type_pb2
from intrinsic.assets.proto import installed_assets_pb2
from intrinsic.assets.proto import installed_assets_pb2_grpc
from intrinsic.assets.proto import view_pb2
from intrinsic.perception.proto.v1 import perception_model_pb2
from intrinsic.solutions import ipython
from intrinsic.util.grpc import error_handling


_CSS_FAILURE_STYLE = (
    'color: #ab0000; font-family: monospace; font-weight: bold; '
    'padding-left: var(--jp-code-padding);'
)
_LAST_RESULT_TIMEOUT_SECONDS = 5
_POSE_ESTIMATOR_RESOURCE_FAMILY_ID = 'perception_model'
_DEFAULT_PACKAGE_NAME_RESOURCE = 'ai.intrinsic'
_DEFAULT_PAGE_SIZE = 200


@error_handling.retry_on_grpc_unavailable
def _list_pose_estimators(
    stub: installed_assets_pb2_grpc.InstalledAssetsStub,
) -> List[installed_assets_pb2.InstalledAsset]:
  """Lists installed data assets."""
  list_installed_assets_request = (
      installed_assets_pb2.ListInstalledAssetsRequest(
          strict_filter=installed_assets_pb2.ListInstalledAssetsRequest.Filter(
              asset_types=[
                  asset_type_pb2.AssetType.ASSET_TYPE_DATA,
              ]
          ),
          view=view_pb2.AssetViewType.ASSET_VIEW_TYPE_DETAIL,
          page_size=_DEFAULT_PAGE_SIZE,
      )
  )
  list_response = stub.ListInstalledAssets(list_installed_assets_request)
  pose_estimators: List[installed_assets_pb2.InstalledAsset] = []
  perception_model_provides = (
      interface_utils.DATA_URI_PREFIX
      + perception_model_pb2.PerceptionModel.DESCRIPTOR.full_name
  )
  for installed_asset in list_response.installed_assets:
    for provides in installed_asset.metadata.provides:
      if provides.uri == perception_model_provides:
        pose_estimators.append(installed_asset)
        break  # break out of the for loop over provides
  return pose_estimators


@dataclasses.dataclass(frozen=True)
class PoseEstimatorId:
  """Wrapper for a PoseEstimatorId proto.

  Attributes:
    id: name of the pose estimator.
    package: Package of the pose estimator.
  """

  id: str
  package: str


class PoseEstimators:
  """Convenience wrapper for pose estimator access."""

  _installed_assets_stub: installed_assets_pb2_grpc.InstalledAssetsStub

  def __init__(
      self,
      installed_assets_stub: installed_assets_pb2_grpc.InstalledAssetsStub,
  ):
    # pyformat: disable
    """Initializes all available pose estimator configs.

    Args:
      installed_assets_stub: Stub for the installed assets service.
    """
    # pyformat: enable
    self._installed_assets_stub = installed_assets_stub

  @error_handling.retry_on_grpc_unavailable
  def _get_pose_estimators(self) -> Dict[str, PoseEstimatorId]:
    """Query pose estimators.

    Returns:
      A dict of pose estimator ids keyed by resource id.

    Raises:
      status.StatusNotOk: If the grpc request failed (propagates grpc error).
    """
    installed_data_assets = _list_pose_estimators(self._installed_assets_stub)
    pose_estimators: Dict[str, PoseEstimatorId] = {
        installed_data_asset.metadata.id_version.id.name: PoseEstimatorId(
            id=installed_data_asset.metadata.id_version.id.name,
            package=installed_data_asset.metadata.id_version.id.package,
        )
        for installed_data_asset in installed_data_assets
    }
    return pose_estimators

  def __getattr__(self, pose_estimator_id: str) -> PoseEstimatorId:
    """Returns the id of the pose estimator.

    Args:
      pose_estimator_id: Resource id of the pose estimator.

    Returns:
      Pose estimator id.

    Raises:
      AttributeError: if there is no pose estimator resource id with the given
      name.
    """
    pose_estimators = self._get_pose_estimators()
    if pose_estimator_id not in pose_estimators:
      raise AttributeError(f'Pose estimator {pose_estimator_id} is unknown')
    return pose_estimators[pose_estimator_id]

  def __len__(self) -> int:
    """Returns the number of pose estimators."""
    return len(self._get_pose_estimators())

  def __str__(self) -> str:
    """Concatenates all pose estimator config paths into a string."""
    return '\n'.join(sorted(self._get_pose_estimators().keys()))

  def __dir__(self) -> List[str]:
    """Lists all pose estimators by key (sorted)."""
    return sorted(list(self._get_pose_estimators().keys()))

  def __getitem__(self, pose_estimator_id: str) -> PoseEstimatorId:
    """Returns the id of the pose estimator.

    Args:
      pose_estimator_id: Resource id of the pose estimator.

    Returns:
      Pose estimator id.

    Raises:
      AttributeError: if there is no pose estimator resource id with the given
      name.
    """
    return self._get_pose_estimators()[pose_estimator_id]

  def __iter__(self) -> Iterator[PoseEstimatorId]:
    """Returns an iterator over all pose estimators.

    Returns:
      Pose estimator ids.
    """
    return iter(self._get_pose_estimators().values())
