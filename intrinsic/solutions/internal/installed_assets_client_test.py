# Copyright 2023 Intrinsic Innovation LLC

"""Tests of InstalledAssetsClient."""

from unittest import mock

from absl.testing import absltest
from absl.testing import parameterized
from google.longrunning import operations_pb2
from google.protobuf import any_pb2
from google.rpc import code_pb2
from google.rpc import status_pb2
import grpc

from intrinsic.assets.processes.proto import process_asset_pb2
from intrinsic.assets.proto import asset_tag_pb2
from intrinsic.assets.proto import asset_type_pb2
from intrinsic.assets.proto import id_pb2
from intrinsic.assets.proto import installed_assets_pb2
from intrinsic.assets.proto import metadata_pb2
from intrinsic.assets.proto import view_pb2
from intrinsic.executive.proto import behavior_tree_pb2
from intrinsic.solutions.internal import installed_assets_client
from intrinsic.solutions.testing import compare


class _MockGrpcError(grpc.RpcError):
  _code: grpc.StatusCode

  def __init__(self, code: grpc.StatusCode):
    self._code = code

  def code(self) -> grpc.StatusCode:
    return self._code


def _process_asset_with_id(
    identifier: str,
) -> installed_assets_pb2.InstalledAsset:
  id_parts = identifier.rpartition(".")
  return process_asset_pb2.ProcessAsset(
      metadata=metadata_pb2.Metadata(
          id_version=id_pb2.IdVersion(
              id=id_pb2.Id(package=id_parts[0], name=id_parts[2])
          )
      ),
      behavior_tree=behavior_tree_pb2.BehaviorTree(
          name="My process",
          root=behavior_tree_pb2.BehaviorTree.Node(
              sequence=behavior_tree_pb2.BehaviorTree.SequenceNode()
          ),
      ),
  )


def _installed_asset_with_id(
    identifier: str,
) -> installed_assets_pb2.InstalledAsset:
  # Create a Process asset as an example
  asset = _process_asset_with_id(identifier)
  return installed_assets_pb2.InstalledAsset(
      metadata=asset.metadata,
      deployment_data=installed_assets_pb2.InstalledAsset.DeploymentData(
          process=installed_assets_pb2.InstalledAsset.ProcessDeploymentData(
              process=asset
          )
      ),
  )


class InstalledAssetsClientTest(parameterized.TestCase):

  _installed_assets: mock.MagicMock
  _operations: mock.MagicMock
  _client: installed_assets_client.InstalledAssetsClient

  def setUp(self):
    super().setUp()
    self._installed_assets = mock.MagicMock()
    self._operations = mock.MagicMock()
    self._client = installed_assets_client.InstalledAssetsClient(
        self._installed_assets, self._operations
    )

  @parameterized.named_parameters(
      (
          "id_proto",
          id_pb2.Id(package="ai.intrinsic", name="process"),
      ),
      (
          "id_str",
          "ai.intrinsic.process",
      ),
  )
  def test_get_installed_asset(self, id):  # pylint: disable=redefined-builtin
    proto = _installed_asset_with_id("ai.intrinsic.process")
    self._installed_assets.GetInstalledAsset.return_value = proto

    result = self._client.get_installed_asset(id)

    compare.assertProto2Equal(self, result, proto)
    self._installed_assets.GetInstalledAsset.assert_called_once_with(
        installed_assets_pb2.GetInstalledAssetRequest(
            id=id_pb2.Id(package="ai.intrinsic", name="process"),
            view=view_pb2.AssetViewType.ASSET_VIEW_TYPE_FULL,
        )
    )

  def test_get_installed_asset_with_view(self):
    proto = _installed_asset_with_id("ai.intrinsic.process")
    self._installed_assets.GetInstalledAsset.return_value = proto

    result = self._client.get_installed_asset(
        id_pb2.Id(package="ai.intrinsic", name="process"),
        view=view_pb2.AssetViewType.ASSET_VIEW_TYPE_DETAIL,
    )

    compare.assertProto2Equal(self, result, proto)
    self._installed_assets.GetInstalledAsset.assert_called_once_with(
        installed_assets_pb2.GetInstalledAssetRequest(
            id=id_pb2.Id(package="ai.intrinsic", name="process"),
            view=view_pb2.AssetViewType.ASSET_VIEW_TYPE_DETAIL,
        )
    )

  def test_get_installed_asset_returns_error(self):
    self._installed_assets.GetInstalledAsset.side_effect = _MockGrpcError(
        code=grpc.StatusCode.PERMISSION_DENIED
    )

    with self.assertRaises(grpc.RpcError) as e:
      self._client.get_installed_asset("ai.intrinsic.process")

    self.assertEqual(e.exception.code(), grpc.StatusCode.PERMISSION_DENIED)

  def test_list_installed_assets(self):
    proto1 = _installed_asset_with_id("ai.intrinsic.process1")
    proto2 = _installed_asset_with_id("ai.intrinsic.process2")
    proto3 = _installed_asset_with_id("ai.intrinsic.process3")
    self._installed_assets.ListInstalledAssets.side_effect = (
        installed_assets_pb2.ListInstalledAssetsResponse(
            installed_assets=[proto1, proto2],
            next_page_token="some_page_token",
        ),
        installed_assets_pb2.ListInstalledAssetsResponse(
            installed_assets=[proto3],
            next_page_token=None,
        ),
    )
    result = self._client.list_all_installed_assets()

    self.assertEqual(result, [proto1, proto2, proto3])
    self._installed_assets.ListInstalledAssets.assert_has_calls([
        mock.call(
            installed_assets_pb2.ListInstalledAssetsRequest(
                strict_filter=installed_assets_pb2.ListInstalledAssetsRequest.Filter(),
                page_size=installed_assets_client._MAX_PAGE_SIZE,
                page_token=None,
                view=view_pb2.AssetViewType.ASSET_VIEW_TYPE_BASIC,
            )
        ),
        mock.call(
            installed_assets_pb2.ListInstalledAssetsRequest(
                strict_filter=installed_assets_pb2.ListInstalledAssetsRequest.Filter(),
                page_size=installed_assets_client._MAX_PAGE_SIZE,
                page_token="some_page_token",
                view=view_pb2.AssetViewType.ASSET_VIEW_TYPE_BASIC,
            )
        ),
    ])

  def test_list_installed_assets_supports_asset_types_filter(self):
    proto = _installed_asset_with_id("ai.intrinsic.process")
    self._installed_assets.ListInstalledAssets.return_value = (
        installed_assets_pb2.ListInstalledAssetsResponse(
            installed_assets=[proto],
            next_page_token=None,
        )
    )

    self._client.list_all_installed_assets(
        asset_types=[
            asset_type_pb2.AssetType.ASSET_TYPE_PROCESS,
            asset_type_pb2.AssetType.ASSET_TYPE_SKILL,
        ]
    )

    self._installed_assets.ListInstalledAssets.assert_called_once_with(
        installed_assets_pb2.ListInstalledAssetsRequest(
            strict_filter=installed_assets_pb2.ListInstalledAssetsRequest.Filter(
                asset_types=[
                    asset_type_pb2.AssetType.ASSET_TYPE_PROCESS,
                    asset_type_pb2.AssetType.ASSET_TYPE_SKILL,
                ]
            ),
            page_size=installed_assets_client._MAX_PAGE_SIZE,
            page_token=None,
            view=view_pb2.AssetViewType.ASSET_VIEW_TYPE_BASIC,
        )
    )

  def test_list_installed_assets_supports_asset_tag_filter(self):
    proto = _installed_asset_with_id("ai.intrinsic.process")
    self._installed_assets.ListInstalledAssets.return_value = (
        installed_assets_pb2.ListInstalledAssetsResponse(
            installed_assets=[proto],
            next_page_token=None,
        )
    )

    self._client.list_all_installed_assets(
        asset_tag=asset_tag_pb2.AssetTag.ASSET_TAG_SUBPROCESS
    )

    self._installed_assets.ListInstalledAssets.assert_called_once_with(
        installed_assets_pb2.ListInstalledAssetsRequest(
            strict_filter=installed_assets_pb2.ListInstalledAssetsRequest.Filter(
                asset_tag=asset_tag_pb2.AssetTag.ASSET_TAG_SUBPROCESS
            ),
            page_size=installed_assets_client._MAX_PAGE_SIZE,
            page_token=None,
            view=view_pb2.AssetViewType.ASSET_VIEW_TYPE_BASIC,
        )
    )

  def test_list_installed_assets_supports_view(self):
    proto = _installed_asset_with_id("ai.intrinsic.process")
    self._installed_assets.ListInstalledAssets.return_value = (
        installed_assets_pb2.ListInstalledAssetsResponse(
            installed_assets=[proto],
            next_page_token=None,
        )
    )

    self._client.list_all_installed_assets(
        view=view_pb2.AssetViewType.ASSET_VIEW_TYPE_DETAIL
    )

    self._installed_assets.ListInstalledAssets.assert_called_once_with(
        installed_assets_pb2.ListInstalledAssetsRequest(
            strict_filter=installed_assets_pb2.ListInstalledAssetsRequest.Filter(),
            page_size=installed_assets_client._MAX_PAGE_SIZE,
            page_token=None,
            view=view_pb2.AssetViewType.ASSET_VIEW_TYPE_DETAIL,
        )
    )

  def test_list_installed_assets_returns_error(self):
    self._installed_assets.ListInstalledAssets.side_effect = _MockGrpcError(
        code=grpc.StatusCode.PERMISSION_DENIED
    )

    with self.assertRaises(grpc.RpcError) as e:
      self._client.list_all_installed_assets()
    self.assertEqual(e.exception.code(), grpc.StatusCode.PERMISSION_DENIED)

  def test_create_installed_asset(self):
    proto = _installed_asset_with_id("ai.intrinsic.process")
    any_proto = any_pb2.Any()
    any_proto.Pack(proto)
    self._installed_assets.CreateInstalledAsset.return_value = (
        operations_pb2.Operation(name="the_operation", done=False)
    )
    self._operations.WaitOperation.side_effect = [
        operations_pb2.Operation(name="the_operation", done=False),
        operations_pb2.Operation(name="the_operation", done=False),
        operations_pb2.Operation(
            name="the_operation", done=True, response=any_proto
        ),
    ]

    result = self._client.create_installed_asset(
        installed_assets_pb2.CreateInstalledAssetRequest.Asset(
            process=proto.deployment_data.process.process
        )
    )

    compare.assertProto2Equal(self, result, proto)
    self._installed_assets.CreateInstalledAsset.assert_called_once_with(
        installed_assets_pb2.CreateInstalledAssetRequest(
            asset=installed_assets_pb2.CreateInstalledAssetRequest.Asset(
                process=proto.deployment_data.process.process
            ),
            policy=installed_assets_pb2.UpdatePolicy.UPDATE_POLICY_UNSPECIFIED,
        )
    )
    self._operations.WaitOperation.assert_called_with(
        operations_pb2.WaitOperationRequest(
            name="the_operation",
            timeout=installed_assets_client._WAIT_OPERATION_TIMEOUT,
        )
    )

  def test_create_installed_asset_supports_update_policy(self):
    proto = _installed_asset_with_id("ai.intrinsic.process")
    any_proto = any_pb2.Any()
    any_proto.Pack(proto)
    self._installed_assets.CreateInstalledAsset.return_value = (
        operations_pb2.Operation(name="the_operation", done=False)
    )
    self._operations.WaitOperation.return_value = operations_pb2.Operation(
        name="the_operation", done=True, response=any_proto
    )

    self._client.create_installed_asset(
        installed_assets_pb2.CreateInstalledAssetRequest.Asset(
            process=proto.deployment_data.process.process
        ),
        update_policy=installed_assets_pb2.UpdatePolicy.UPDATE_POLICY_UPDATE_UNUSED,
    )

    self._installed_assets.CreateInstalledAsset.assert_called_once_with(
        installed_assets_pb2.CreateInstalledAssetRequest(
            asset=installed_assets_pb2.CreateInstalledAssetRequest.Asset(
                process=proto.deployment_data.process.process
            ),
            policy=installed_assets_pb2.UpdatePolicy.UPDATE_POLICY_UPDATE_UNUSED,
        )
    )

  def test_create_installed_asset_returns_error(self):
    proto = _installed_asset_with_id("ai.intrinsic.process")
    self._installed_assets.CreateInstalledAsset.side_effect = _MockGrpcError(
        code=grpc.StatusCode.PERMISSION_DENIED
    )

    with self.assertRaises(grpc.RpcError) as e:
      self._client.create_installed_asset(
          installed_assets_pb2.CreateInstalledAssetRequest.Asset(
              process=proto.deployment_data.process.process
          ),
      )
    self.assertEqual(e.exception.code(), grpc.StatusCode.PERMISSION_DENIED)

  def test_create_installed_asset_returns_operation_error(self):
    proto = _installed_asset_with_id("ai.intrinsic.process")
    self._installed_assets.CreateInstalledAsset.return_value = (
        operations_pb2.Operation(name="the_operation", done=False)
    )
    self._operations.WaitOperation.return_value = operations_pb2.Operation(
        name="the_operation",
        done=True,
        error=status_pb2.Status(code=code_pb2.INTERNAL, message="some error"),
    )

    with self.assertRaises(installed_assets_client.OperationError) as e:
      self._client.create_installed_asset(
          installed_assets_pb2.CreateInstalledAssetRequest.Asset(
              process=proto.deployment_data.process.process
          ),
      )
    self.assertEqual(e.exception.code(), grpc.StatusCode.INTERNAL)
    self.assertEqual(e.exception.message(), "some error")

  @parameterized.named_parameters(
      (
          "id_proto",
          id_pb2.Id(package="ai.intrinsic", name="process"),
      ),
      (
          "id_str",
          "ai.intrinsic.process",
      ),
  )
  def test_delete_installed_asset(self, id):  # pylint: disable=redefined-builtin
    self._installed_assets.DeleteInstalledAsset.return_value = (
        operations_pb2.Operation(name="the_operation", done=False)
    )
    self._operations.WaitOperation.side_effect = [
        operations_pb2.Operation(name="the_operation", done=False),
        operations_pb2.Operation(name="the_operation", done=False),
        operations_pb2.Operation(name="the_operation", done=True),
    ]

    self._client.delete_installed_asset(id)

    self._installed_assets.DeleteInstalledAsset.assert_called_once_with(
        installed_assets_pb2.DeleteInstalledAssetRequest(
            asset=id_pb2.Id(package="ai.intrinsic", name="process")
        )
    )
    self._operations.WaitOperation.assert_called_with(
        operations_pb2.WaitOperationRequest(
            name="the_operation",
            timeout=installed_assets_client._WAIT_OPERATION_TIMEOUT,
        )
    )

  def test_delete_installed_asset_supports_delete_policy(self):
    self._installed_assets.DeleteInstalledAsset.return_value = (
        operations_pb2.Operation(name="the_operation", done=False)
    )
    self._operations.WaitOperation.return_value = operations_pb2.Operation(
        name="the_operation", done=True
    )

    self._client.delete_installed_asset(
        id_pb2.Id(package="ai.intrinsic", name="process"),
        delete_policy=installed_assets_pb2.DeletePolicy.POLICY_REJECT_USED,
    )

    self._installed_assets.DeleteInstalledAsset.assert_called_once_with(
        installed_assets_pb2.DeleteInstalledAssetRequest(
            asset=id_pb2.Id(package="ai.intrinsic", name="process"),
            policy=installed_assets_pb2.DeletePolicy.POLICY_REJECT_USED,
        )
    )

  def test_delete_installed_asset_returns_error(self):
    self._installed_assets.DeleteInstalledAsset.side_effect = _MockGrpcError(
        code=grpc.StatusCode.PERMISSION_DENIED
    )

    with self.assertRaises(grpc.RpcError) as e:
      self._client.delete_installed_asset("ai.intrinsic.process")
    self.assertEqual(e.exception.code(), grpc.StatusCode.PERMISSION_DENIED)

  def test_delete_installed_asset_returns_operation_error(self):
    self._installed_assets.DeleteInstalledAsset.return_value = (
        operations_pb2.Operation(name="the_operation", done=False)
    )
    self._operations.WaitOperation.return_value = operations_pb2.Operation(
        name="the_operation",
        done=True,
        error=status_pb2.Status(code=code_pb2.INTERNAL, message="some error"),
    )

    with self.assertRaises(installed_assets_client.OperationError) as e:
      self._client.delete_installed_asset("ai.intrinsic.process")
    self.assertEqual(e.exception.code(), grpc.StatusCode.INTERNAL)
    self.assertEqual(e.exception.message(), "some error")


if __name__ == "__main__":
  absltest.main()
