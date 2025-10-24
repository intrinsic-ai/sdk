# Copyright 2023 Intrinsic Innovation LLC

"""Tests of process_providing.py."""

from unittest import mock

from absl.testing import absltest
from google.longrunning import operations_pb2
from google.protobuf import any_pb2
from google.protobuf import descriptor_pb2
from google.rpc import code_pb2
from google.rpc import status_pb2
import grpc
from intrinsic.assets.processes.proto import process_asset_pb2
from intrinsic.assets.proto import asset_type_pb2
from intrinsic.assets.proto import id_pb2
from intrinsic.assets.proto import installed_assets_pb2
from intrinsic.assets.proto import metadata_pb2
from intrinsic.assets.proto import view_pb2
from intrinsic.executive.proto import behavior_tree_pb2
from intrinsic.frontend.solution_service.proto import solution_service_pb2
from intrinsic.skills.proto import skills_pb2
from intrinsic.solutions import behavior_tree
from intrinsic.solutions.internal import behavior_call
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
  return process_asset_pb2.ProcessAsset(
      metadata=metadata_pb2.Metadata(
          id_version=id_pb2.IdVersion(
              id=id_pb2.Id(package=id_parts[0], name=id_parts[2])
          )
      ),
      behavior_tree=_behavior_tree_with_name(identifier + " display name"),
  )


def _installed_asset_with_id(
    identifier: str,
) -> installed_assets_pb2.InstalledAsset:
  asset = _process_asset_with_id(identifier)
  return installed_assets_pb2.InstalledAsset(
      metadata=asset.metadata,
      deployment_data=installed_assets_pb2.InstalledAsset.DeploymentData(
          process=installed_assets_pb2.InstalledAsset.ProcessDeploymentData(
              process=asset
          )
      ),
  )


def _default_task() -> behavior_tree.Task:
  return behavior_tree.Task(
      behavior_call.Action(skill_id="ai.intrinsic.skill-0")
  )


class ProcessProvidingTest(absltest.TestCase):

  def setUp(self):
    super().setUp()
    self._solution_service = mock.MagicMock()
    self._installed_assets = mock.MagicMock()
    self._operations = mock.MagicMock()
    self._processes = process_providing.Processes(
        self._solution_service, self._installed_assets, self._operations
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
    proto1 = _installed_asset_with_id("ai.intrinsic.process1")
    proto2 = _installed_asset_with_id("ai.intrinsic.process2")
    proto3 = _installed_asset_with_id("ai.intrinsic.process3")
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
    asset_proto1 = _installed_asset_with_id("ai.intrinsic.process")
    asset_proto2 = _installed_asset_with_id("main.bt.pb")
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
        return _installed_asset_with_id("ai.intrinsic.process")
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
    asset_proto = _installed_asset_with_id("ai.intrinsic.process")
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
    bt = behavior_tree.BehaviorTree("tree1", root=_default_task())

    self._processes["tree1"] = bt

    self._solution_service.UpdateBehaviorTree.assert_called_once_with(
        solution_service_pb2.UpdateBehaviorTreeRequest(
            behavior_tree=bt.proto,
            allow_missing=True,
        )
    )

  def test_save_legacy_process(self):
    bt = behavior_tree.BehaviorTree("tree1", root=_default_task())

    self._processes.save(bt)

    self._solution_service.UpdateBehaviorTree.assert_called_once_with(
        solution_service_pb2.UpdateBehaviorTreeRequest(
            behavior_tree=bt.proto,
            allow_missing=True,
        )
    )

  def test_save_process_asset(self):
    # The behavior tree to be saved.
    bt_proto = behavior_tree_pb2.BehaviorTree(
        name="My tree",
        description=skills_pb2.Skill(
            id="ai.intrinsic.my_process",
            id_version="ai.intrinsic.my_process.0.0.1+old",
            skill_name="ai.intrinsic.my_process",
            package_name="ai.intrinsic",
            display_name="My tree",
        ),
        root=behavior_tree_pb2.BehaviorTree.Node(
            fail=behavior_tree_pb2.BehaviorTree.FailNode()
        ),
    )
    # The behavior tree as expected by the installed assets service when saving.
    bt_proto_saved = behavior_tree_pb2.BehaviorTree()
    bt_proto_saved.CopyFrom(bt_proto)
    bt_proto_saved.description.ClearField("id_version")

    # The metadata to be saved.
    metadata_proto = metadata_pb2.Metadata(
        id_version=id_pb2.IdVersion(
            id=id_pb2.Id(package="ai.intrinsic", name="my_process"),
            version="0.0.1+old",
        ),
        display_name="My tree",
        file_descriptor_set=descriptor_pb2.FileDescriptorSet(
            file=[descriptor_pb2.FileDescriptorProto(name="foo.proto")],
        ),
    )
    # The metadata as expected by the installed assets service when saving.
    metadata_proto_saved = metadata_pb2.Metadata()
    metadata_proto_saved.CopyFrom(metadata_proto)
    metadata_proto_saved.id_version.ClearField("version")
    metadata_proto_saved.ClearField("file_descriptor_set")
    # The metadata returned by the installed assets service after saving.
    metadata_proto_new = metadata_pb2.Metadata()
    metadata_proto_new.CopyFrom(metadata_proto)
    metadata_proto_new.id_version.version = "0.0.1+new"

    # The installed asset returned by the installed assets service in the
    # operation result after saving. This includes only the metadata, not the
    # behavior tree.
    installed_asset_any = any_pb2.Any()
    installed_asset_any.Pack(
        installed_assets_pb2.InstalledAsset(metadata=metadata_proto_new)
    )

    bt = behavior_tree.BehaviorTree.create_from_proto(
        process_asset_pb2.ProcessAsset(
            metadata=metadata_proto,
            behavior_tree=bt_proto,
        )
    )
    self._installed_assets.CreateInstalledAsset.return_value = (
        operations_pb2.Operation(name="the_operation", done=False)
    )
    self._operations.WaitOperation.side_effect = [
        operations_pb2.Operation(name="the_operation", done=False),
        operations_pb2.Operation(
            name="the_operation", done=True, response=installed_asset_any
        ),
    ]

    self._processes.save(bt)

    # Expect BehaviorTree to have been updated with the new version generated
    # by the installed assets service.
    self.assertEqual(bt.asset_metadata_proto.id_version.version, "0.0.1+new")
    self.assertEqual(
        bt.proto.description.id_version,
        "ai.intrinsic.my_process.0.0.1+new",
    )
    self._installed_assets.CreateInstalledAsset.assert_called_once_with(
        installed_assets_pb2.CreateInstalledAssetRequest(
            asset=installed_assets_pb2.CreateInstalledAssetRequest.Asset(
                process=process_asset_pb2.ProcessAsset(
                    metadata=metadata_proto_saved,
                    behavior_tree=bt_proto_saved,
                ),
            ),
        ),
    )
    self._operations.WaitOperation.assert_called_with(
        operations_pb2.WaitOperationRequest(
            name="the_operation",
            timeout=process_providing._WAIT_OPERATION_TIMEOUT,
        )
    )

  def test_save_process_asset_succeeds_on_already_exists(self):
    bt = behavior_tree.BehaviorTree.create_from_proto(
        _process_asset_with_id("ai.intrinsic.process")
    )
    self._installed_assets.CreateInstalledAsset.return_value = (
        operations_pb2.Operation(name="the_operation", done=False)
    )
    self._operations.WaitOperation.side_effect = [
        operations_pb2.Operation(
            name="the_operation",
            done=True,
            error=status_pb2.Status(
                code=code_pb2.ALREADY_EXISTS, message="process already exists"
            ),
        ),
    ]

    # Expect no exception.
    self._processes.save(bt)

  def test_save_process_asset_returns_operation_error(self):
    bt = behavior_tree.BehaviorTree.create_from_proto(
        _process_asset_with_id("ai.intrinsic.process")
    )
    self._installed_assets.CreateInstalledAsset.return_value = (
        operations_pb2.Operation(name="the_operation", done=False)
    )
    self._operations.WaitOperation.side_effect = [
        operations_pb2.Operation(
            name="the_operation",
            done=True,
            error=status_pb2.Status(
                code=code_pb2.INVALID_ARGUMENT, message="some error"
            ),
        ),
    ]

    with self.assertRaisesRegex(RuntimeError, "INVALID_ARGUMENT.*some error"):
      self._processes.save(bt)

  def test_delitem_legacy_process(self):
    del self._processes["My tree"]

    self._installed_assets.DeleteInstalledAsset.assert_not_called()
    self._solution_service.DeleteBehaviorTree.assert_called_once_with(
        solution_service_pb2.DeleteBehaviorTreeRequest(name="My tree")
    )

  def test_delitem_legacy_process_with_asset_id_like_name(self):
    self._installed_assets.DeleteInstalledAsset.return_value = (
        operations_pb2.Operation(name="the_operation", done=False)
    )
    self._operations.WaitOperation.side_effect = [
        operations_pb2.Operation(name="the_operation", done=False),
        operations_pb2.Operation(name="the_operation", done=False),
        operations_pb2.Operation(
            name="the_operation",
            done=True,
            error=status_pb2.Status(code=code_pb2.NOT_FOUND),
        ),
    ]

    del self._processes["main.bt.pb"]

    # Expect that first the Process asset "main.bt.pb" is tried to be deleted,
    # but it does not exist.
    self._installed_assets.DeleteInstalledAsset.assert_called_once_with(
        installed_assets_pb2.DeleteInstalledAssetRequest(
            asset=id_pb2.Id(package="main.bt", name="pb")
        )
    )
    self._operations.WaitOperation.assert_called_with(
        operations_pb2.WaitOperationRequest(
            name="the_operation",
            timeout=process_providing._WAIT_OPERATION_TIMEOUT,
        )
    )
    # Then a legacy deletion of the tree with name "main.bt.pb" should be
    # attempted.
    self._solution_service.DeleteBehaviorTree.assert_called_once_with(
        solution_service_pb2.DeleteBehaviorTreeRequest(name="main.bt.pb")
    )

  def test_delitem_process_asset(self):
    self._installed_assets.DeleteInstalledAsset.return_value = (
        operations_pb2.Operation(name="the_operation", done=False)
    )
    self._operations.WaitOperation.side_effect = [
        operations_pb2.Operation(name="the_operation", done=False),
        operations_pb2.Operation(name="the_operation", done=False),
        operations_pb2.Operation(name="the_operation", done=True),
    ]

    del self._processes["ai.intrinsic.process"]

    self._installed_assets.DeleteInstalledAsset.assert_called_once_with(
        installed_assets_pb2.DeleteInstalledAssetRequest(
            asset=id_pb2.Id(package="ai.intrinsic", name="process")
        )
    )
    self._operations.WaitOperation.assert_called_with(
        operations_pb2.WaitOperationRequest(
            name="the_operation",
            timeout=process_providing._WAIT_OPERATION_TIMEOUT,
        )
    )
    self._solution_service.DeleteBehaviorTree.assert_not_called()

  def test_delitem_legacy_process_raises_operation_error(self):
    self._installed_assets.DeleteInstalledAsset.return_value = (
        operations_pb2.Operation(name="the_operation", done=False)
    )
    self._operations.WaitOperation.side_effect = [
        operations_pb2.Operation(name="the_operation", done=False),
        operations_pb2.Operation(name="the_operation", done=False),
        operations_pb2.Operation(
            name="the_operation",
            done=True,
            error=status_pb2.Status(
                code=code_pb2.INVALID_ARGUMENT, message="some error"
            ),
        ),
    ]

    with self.assertRaises(KeyError) as e:
      del self._processes["ai.intrinsic.process"]

    cause = e.exception.__cause__
    self.assertIsInstance(cause, RuntimeError)
    self.assertRegex(str(cause), "INVALID_ARGUMENT.*some error")


if __name__ == "__main__":
  absltest.main()
