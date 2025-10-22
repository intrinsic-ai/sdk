# Copyright 2023 Intrinsic Innovation LLC

"""Tests of process_providing.py."""

from unittest import mock

from absl.testing import absltest
import grpc
from intrinsic.assets.processes.proto import process_asset_pb2
from intrinsic.assets.proto import asset_type_pb2
from intrinsic.assets.proto import id_pb2
from intrinsic.assets.proto import installed_assets_pb2
from intrinsic.assets.proto import metadata_pb2
from intrinsic.assets.proto import view_pb2
from intrinsic.executive.proto import behavior_tree_pb2
from intrinsic.frontend.solution_service.proto import solution_service_pb2
from intrinsic.solutions import behavior_tree
from intrinsic.solutions.internal import process_providing


def _behavior_tree_with_name(name: str):
  bt = behavior_tree_pb2.BehaviorTree()
  bt.name = name
  bt.root.task.call_behavior.skill_id = "ai.intrinsic.skill-0"
  bt.root.name = "ai.intrinsic.skill-0"
  bt.root.id = 1
  bt.root.state = behavior_tree_pb2.BehaviorTree.State.ACCEPTED
  return bt


def _process_asset_with_id(
    identifier: str,
) -> installed_assets_pb2.InstalledAsset:
  id_parts = identifier.rpartition(".")
  metadata = metadata_pb2.Metadata(
      id_version=id_pb2.IdVersion(
          id=id_pb2.Id(package=id_parts[0], name=id_parts[2])
      )
  )
  return installed_assets_pb2.InstalledAsset(
      metadata=metadata,
      deployment_data=installed_assets_pb2.InstalledAsset.DeploymentData(
          process=installed_assets_pb2.InstalledAsset.ProcessDeploymentData(
              process=process_asset_pb2.ProcessAsset(
                  metadata=metadata,
                  behavior_tree=_behavior_tree_with_name(
                      identifier + " display name"
                  ),
              )
          )
      ),
  )


class ProcessProvidingTest(absltest.TestCase):

  def setUp(self):
    super().setUp()
    self._solution_service = mock.MagicMock()
    self._installed_assets = mock.MagicMock()
    self._processes = process_providing.Processes(
        self._solution_service, self._installed_assets
    )

    # Default behavior for the mocks
    self._solution_service.ListBehaviorTrees.return_value = (
        solution_service_pb2.ListBehaviorTreesResponse(
            behavior_trees=None,
            next_page_token=None,
        )
    )
    self._installed_assets.ListInstalledAssets.return_value = (
        installed_assets_pb2.ListInstalledAssetsResponse(
            installed_assets=None,
            next_page_token=None,
        )
    )

  # Tests keys(), items(), values(), __iter__() at the same time.
  def test_iterables_empty(self):
    self.assertEqual(list(self._processes.keys()), [])
    self.assertEqual(list(self._processes.items()), [])
    self.assertEqual(list(self._processes.values()), [])
    self.assertEqual(list(self._processes), [])

  # Tests keys(), items(), values(), __iter__() at the same time.
  def test_iterables_multiple_pages_legacy_processes(self):
    proto1 = _behavior_tree_with_name("tree1")
    proto2 = _behavior_tree_with_name("tree2")
    proto3 = _behavior_tree_with_name("tree3")
    # Two responses for each call to keys(), items(), values(), __iter__()
    self._solution_service.ListBehaviorTrees.side_effect = 4 * (
        solution_service_pb2.ListBehaviorTreesResponse(
            behavior_trees=[proto1, proto2], next_page_token="some_token"
        ),
        solution_service_pb2.ListBehaviorTreesResponse(
            behavior_trees=[proto3], next_page_token=None
        ),
    )

    self.assertEqual(list(self._processes.keys()), ["tree1", "tree2", "tree3"])
    self.assertEqual(
        [
            (name, bt.proto, bt.asset_metadata_proto)
            for name, bt in self._processes.items()
        ],
        [
            ("tree1", proto1, None),
            ("tree2", proto2, None),
            ("tree3", proto3, None),
        ],
    )
    self.assertEqual(
        [
            (bt.proto, bt.asset_metadata_proto)
            for bt in self._processes.values()
        ],
        [(proto1, None), (proto2, None), (proto3, None)],
    )
    self.assertEqual(list(self._processes), ["tree1", "tree2", "tree3"])

    def expected_call(*, page_token: str, full_view: bool):
      return mock.call(
          solution_service_pb2.ListBehaviorTreesRequest(
              page_size=process_providing._SOLUTION_SERVICE_MAX_PAGE_SIZE,
              view=(
                  solution_service_pb2.BehaviorTreeView.BEHAVIOR_TREE_VIEW_FULL
                  if full_view
                  else solution_service_pb2.BehaviorTreeView.BEHAVIOR_TREE_VIEW_BASIC
              ),
              page_token=page_token,
          )
      )

    self.assertSequenceEqual(
        self._solution_service.ListBehaviorTrees.mock_calls,
        (
            # Calls for keys()
            expected_call(page_token=None, full_view=False),
            expected_call(page_token="some_token", full_view=False),
            # Calls for items()
            expected_call(page_token=None, full_view=True),
            expected_call(page_token="some_token", full_view=True),
            # Calls for values()
            expected_call(page_token=None, full_view=True),
            expected_call(page_token="some_token", full_view=True),
            # Calls for __iter__()
            expected_call(page_token=None, full_view=False),
            expected_call(page_token="some_token", full_view=False),
        ),
    )

  # Tests keys(), items(), values(), __iter__() at the same time.
  def test_iterables_multiple_pages_process_assets(self):
    proto1 = _process_asset_with_id("ai.intrinsic.process1")
    proto2 = _process_asset_with_id("ai.intrinsic.process2")
    proto3 = _process_asset_with_id("ai.intrinsic.process3")
    # Two responses for each call to keys(), items(), values(), __iter__()
    self._installed_assets.ListInstalledAssets.side_effect = 4 * (
        installed_assets_pb2.ListInstalledAssetsResponse(
            installed_assets=[proto1, proto2], next_page_token="some_token"
        ),
        installed_assets_pb2.ListInstalledAssetsResponse(
            installed_assets=[proto3], next_page_token=None
        ),
    )

    self.assertEqual(
        list(self._processes.keys()),
        [
            "ai.intrinsic.process1",
            "ai.intrinsic.process2",
            "ai.intrinsic.process3",
        ],
    )
    self.assertEqual(
        [
            (identifier, bt.proto, bt.asset_metadata_proto)
            for identifier, bt in self._processes.items()
        ],
        [
            (
                "ai.intrinsic.process1",
                proto1.deployment_data.process.process.behavior_tree,
                proto1.deployment_data.process.process.metadata,
            ),
            (
                "ai.intrinsic.process2",
                proto2.deployment_data.process.process.behavior_tree,
                proto2.deployment_data.process.process.metadata,
            ),
            (
                "ai.intrinsic.process3",
                proto3.deployment_data.process.process.behavior_tree,
                proto3.deployment_data.process.process.metadata,
            ),
        ],
    )
    self.assertEqual(
        [
            (bt.proto, bt.asset_metadata_proto)
            for bt in self._processes.values()
        ],
        [
            (
                proto1.deployment_data.process.process.behavior_tree,
                proto1.deployment_data.process.process.metadata,
            ),
            (
                proto2.deployment_data.process.process.behavior_tree,
                proto2.deployment_data.process.process.metadata,
            ),
            (
                proto3.deployment_data.process.process.behavior_tree,
                proto3.deployment_data.process.process.metadata,
            ),
        ],
    )
    self.assertEqual(
        list(self._processes),
        [
            "ai.intrinsic.process1",
            "ai.intrinsic.process2",
            "ai.intrinsic.process3",
        ],
    )

    def expected_call(*, page_token: str, full_view: bool):
      return mock.call(
          installed_assets_pb2.ListInstalledAssetsRequest(
              strict_filter=installed_assets_pb2.ListInstalledAssetsRequest.Filter(
                  asset_types=[asset_type_pb2.AssetType.ASSET_TYPE_PROCESS]
              ),
              page_size=process_providing._INSTALLED_ASSETS_MAX_PAGE_SIZE,
              view=(
                  view_pb2.AssetViewType.ASSET_VIEW_TYPE_FULL
                  if full_view
                  else view_pb2.AssetViewType.ASSET_VIEW_TYPE_BASIC
              ),
              page_token=page_token,
          )
      )

    self.assertSequenceEqual(
        self._installed_assets.ListInstalledAssets.mock_calls,
        (
            # Calls for keys()
            expected_call(page_token=None, full_view=False),
            expected_call(page_token="some_token", full_view=False),
            # Calls for items()
            expected_call(page_token=None, full_view=True),
            expected_call(page_token="some_token", full_view=True),
            # Calls for values()
            expected_call(page_token=None, full_view=True),
            expected_call(page_token="some_token", full_view=True),
            # Calls for __iter__()
            expected_call(page_token=None, full_view=False),
            expected_call(page_token="some_token", full_view=False),
        ),
    )

  # Tests keys(), items(), values(), __iter__() at the same time.
  def test_iterables_legacy_processes_and_process_assets(self):
    asset_proto1 = _process_asset_with_id("ai.intrinsic.process")
    asset_proto2 = _process_asset_with_id("main.bt.pb")
    # Shadowed by the process asset with id "main.bt.pb".
    bt_proto3 = _behavior_tree_with_name("main.bt.pb")
    bt_proto4 = _behavior_tree_with_name("My tree")
    self._installed_assets.ListInstalledAssets.return_value = (
        installed_assets_pb2.ListInstalledAssetsResponse(
            installed_assets=[asset_proto1, asset_proto2],
            next_page_token=None,
        )
    )
    self._solution_service.ListBehaviorTrees.return_value = (
        solution_service_pb2.ListBehaviorTreesResponse(
            behavior_trees=[bt_proto3, bt_proto4],
            next_page_token=None,
        )
    )

    self.assertEqual(
        list(self._processes.keys()),
        ["ai.intrinsic.process", "main.bt.pb", "My tree"],
    )
    self.assertEqual(
        [
            (identifier, bt.proto, bt.asset_metadata_proto)
            for identifier, bt in self._processes.items()
        ],
        [
            (
                "ai.intrinsic.process",
                asset_proto1.deployment_data.process.process.behavior_tree,
                asset_proto1.deployment_data.process.process.metadata,
            ),
            (
                "main.bt.pb",
                asset_proto2.deployment_data.process.process.behavior_tree,
                asset_proto2.deployment_data.process.process.metadata,
            ),
            ("My tree", bt_proto4, None),
        ],
    )
    self.assertEqual(
        [
            (bt.proto, bt.asset_metadata_proto)
            for bt in self._processes.values()
        ],
        [
            (
                asset_proto1.deployment_data.process.process.behavior_tree,
                asset_proto1.deployment_data.process.process.metadata,
            ),
            (
                asset_proto2.deployment_data.process.process.behavior_tree,
                asset_proto2.deployment_data.process.process.metadata,
            ),
            (bt_proto4, None),
        ],
    )
    self.assertEqual(
        list(self._processes),
        ["ai.intrinsic.process", "main.bt.pb", "My tree"],
    )

  def test_contains(self):

    def mock_get_installed_asset(
        request: installed_assets_pb2.GetInstalledAssetRequest,
    ):
      if request.id == id_pb2.Id(package="ai.intrinsic", name="process"):
        return _process_asset_with_id("ai.intrinsic.process")
      else:
        error = grpc.RpcError(str(request.id) + " not found")
        error.code = lambda: grpc.StatusCode.NOT_FOUND
        raise error

    def mock_get_behavior_tree(
        request: solution_service_pb2.GetBehaviorTreeRequest,
    ):
      if request.name == "My tree":
        return _behavior_tree_with_name("My tree")
      # A tree with a name that conforms to the asset id format.
      elif request.name == "main.bt.pb":
        return _behavior_tree_with_name("main.bt.pb")
      else:
        error = grpc.RpcError(request.name + " not found")
        error.code = lambda: grpc.StatusCode.NOT_FOUND
        raise error

    self._solution_service.GetBehaviorTree.side_effect = mock_get_behavior_tree
    self._installed_assets.GetInstalledAsset.side_effect = (
        mock_get_installed_asset
    )

    self.assertIn("ai.intrinsic.process", self._processes)
    self.assertIn("main.bt.pb", self._processes)
    self.assertIn("My tree", self._processes)
    self.assertNotIn("non_existent_tree", self._processes)

  def test_getitem(self):
    asset_proto = _process_asset_with_id("ai.intrinsic.process")
    bt_proto1 = _behavior_tree_with_name("My tree")
    # A tree with a name that conforms to the asset id format.
    bt_proto2 = _behavior_tree_with_name("main.bt.pb")

    def mock_get_installed_asset(
        request: installed_assets_pb2.GetInstalledAssetRequest,
    ):
      if request.id == id_pb2.Id(package="ai.intrinsic", name="process"):
        return asset_proto
      else:
        error = grpc.RpcError(str(request.id) + " not found")
        error.code = lambda: grpc.StatusCode.NOT_FOUND
        raise error

    def mock_get_behavior_tree(
        request: solution_service_pb2.GetBehaviorTreeRequest,
    ):
      if request.name == bt_proto1.name:
        return bt_proto1
      elif request.name == bt_proto2.name:
        return bt_proto2
      else:
        error = grpc.RpcError(request.name + " not found")
        error.code = lambda: grpc.StatusCode.NOT_FOUND
        raise error

    self._solution_service.GetBehaviorTree.side_effect = mock_get_behavior_tree
    self._installed_assets.GetInstalledAsset.side_effect = (
        mock_get_installed_asset
    )

    self.assertEqual(
        self._processes["ai.intrinsic.process"].proto,
        asset_proto.deployment_data.process.process.behavior_tree,
    )
    self.assertEqual(
        self._processes["ai.intrinsic.process"].asset_metadata_proto,
        asset_proto.deployment_data.process.process.metadata,
    )
    self.assertEqual(self._processes["My tree"].proto, bt_proto1)
    self.assertIsNone(self._processes["My tree"].asset_metadata_proto)
    self.assertEqual(self._processes["main.bt.pb"].proto, bt_proto2)
    self.assertIsNone(self._processes["main.bt.pb"].asset_metadata_proto)
    with self.assertRaises(KeyError):
      self._processes["non_existent_tree"]  # pylint: disable=pointless-statement

  def test_setitem(self):
    bt = behavior_tree.BehaviorTree("tree1", root=behavior_tree.Fail("Failure"))

    self._processes["tree1"] = bt

    self._solution_service.UpdateBehaviorTree.assert_called_once_with(
        solution_service_pb2.UpdateBehaviorTreeRequest(
            behavior_tree=bt.proto,
            allow_missing=True,
        )
    )

  def test_delitem(self):
    del self._processes["tree1"]

    self._solution_service.DeleteBehaviorTree.assert_called_once_with(
        solution_service_pb2.DeleteBehaviorTreeRequest(name="tree1")
    )


if __name__ == "__main__":
  absltest.main()
