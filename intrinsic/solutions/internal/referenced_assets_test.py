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

"""Tests of referenced_assets.py."""

from unittest import mock

from absl.testing import absltest
import grpc

from intrinsic.assets import id_utils
from intrinsic.assets.install import installed_assets_client
from intrinsic.assets.processes.proto import process_asset_pb2
from intrinsic.assets.proto import asset_type_pb2
from intrinsic.assets.proto import installed_assets_pb2
from intrinsic.assets.proto import view_pb2
from intrinsic.executive.proto import behavior_tree_pb2
from intrinsic.solutions.internal import referenced_assets


def _create_task_node(skill_id: str) -> behavior_tree_pb2.BehaviorTree.Node:
  node = behavior_tree_pb2.BehaviorTree.Node()
  node.task.call_behavior.skill_id = skill_id
  return node


def _create_installed_asset(
    package: str, name: str, version: str, asset_type: asset_type_pb2.AssetType
) -> installed_assets_pb2.InstalledAsset:
  installed = installed_assets_pb2.InstalledAsset()
  installed.asset.catalog.id_version.id.package = package
  installed.asset.catalog.id_version.id.name = name
  installed.asset.catalog.id_version.version = version
  installed.asset.catalog.asset_type = asset_type
  return installed


class ReferencedAssetsTest(absltest.TestCase):

  def setUp(self):
    super().setUp()
    self._installed_assets = mock.MagicMock(
        spec=installed_assets_client.InstalledAssetsClient
    )

  def test_update_referenced_assets_with_skill_references(self):
    pa = process_asset_pb2.ProcessAsset()
    seq = behavior_tree_pb2.BehaviorTree.SequenceNode()
    seq.children.append(_create_task_node("ai.intrinsic.skill_a"))
    seq.children.append(_create_task_node("ai.intrinsic.skill_b"))
    seq.children.append(_create_task_node("ai.intrinsic.skill_a"))
    pa.behavior_tree.root.sequence.CopyFrom(seq)

    installed_assets = [
        _create_installed_asset(
            "ai.intrinsic",
            "skill_a",
            "1.0.0",
            asset_type_pb2.AssetType.ASSET_TYPE_SKILL,
        ),
        _create_installed_asset(
            "ai.intrinsic",
            "skill_b",
            "1.0.0",
            asset_type_pb2.AssetType.ASSET_TYPE_SKILL,
        ),
    ]
    self._installed_assets.batch_get_installed_assets.return_value = (
        installed_assets
    )

    referenced_assets.update_referenced_assets(pa, self._installed_assets)
    self.assertEqual(
        set(pa.assets.keys()), {"ai.intrinsic.skill_a", "ai.intrinsic.skill_b"}
    )
    self._installed_assets.batch_get_installed_assets.assert_called_once()
    called_args, called_kwargs = (
        self._installed_assets.batch_get_installed_assets.call_args
    )
    self.assertEqual(
        set(called_args[0]), {"ai.intrinsic.skill_a", "ai.intrinsic.skill_b"}
    )
    self.assertEqual(
        called_kwargs.get("view"), view_pb2.AssetViewType.ASSET_VIEW_TYPE_BASIC
    )

  def test_update_referenced_assets_skips_called_tree_state(self):
    pa = process_asset_pb2.ProcessAsset()
    task_node = _create_task_node("ai.intrinsic.skill_outer")
    called_tree = behavior_tree_pb2.BehaviorTree()
    called_tree.root.CopyFrom(
        _create_task_node("ai.intrinsic.skill_inner_leaked")
    )
    task_node.task.called_tree_state.CopyFrom(called_tree)
    pa.behavior_tree.root.CopyFrom(task_node)

    self._installed_assets.batch_get_installed_assets.return_value = [
        _create_installed_asset(
            "ai.intrinsic",
            "skill_outer",
            "1.0.0",
            asset_type_pb2.AssetType.ASSET_TYPE_SKILL,
        )
    ]

    referenced_assets.update_referenced_assets(pa, self._installed_assets)
    self.assertEqual(set(pa.assets.keys()), {"ai.intrinsic.skill_outer"})
    called_args, _ = self._installed_assets.batch_get_installed_assets.call_args
    self.assertEqual(called_args[0], ["ai.intrinsic.skill_outer"])

  def test_update_referenced_assets_empty_tree_clears_assets(self):
    pa = process_asset_pb2.ProcessAsset()
    pa.assets["old.stale.asset"].catalog.id_version.id.package = "old"
    referenced_assets.update_referenced_assets(pa, self._installed_assets)
    self.assertEqual(len(pa.assets), 0)
    self._installed_assets.batch_get_installed_assets.assert_not_called()

  def test_update_referenced_assets_populates_installed_asset(self):
    pa = process_asset_pb2.ProcessAsset()
    pa.behavior_tree.root.CopyFrom(_create_task_node("ai.intrinsic.skill_1"))

    self._installed_assets.batch_get_installed_assets.return_value = [
        _create_installed_asset(
            "ai.intrinsic",
            "skill_1",
            "1.2.3",
            asset_type_pb2.AssetType.ASSET_TYPE_SKILL,
        )
    ]
    referenced_assets.update_referenced_assets(pa, self._installed_assets)
    self.assertIn("ai.intrinsic.skill_1", pa.assets)
    processed = pa.assets["ai.intrinsic.skill_1"]
    self.assertEqual(
        processed.catalog.asset_type, asset_type_pb2.AssetType.ASSET_TYPE_SKILL
    )
    self.assertEqual(processed.catalog.id_version.id.package, "ai.intrinsic")
    self.assertEqual(processed.catalog.id_version.id.name, "skill_1")
    self.assertEqual(processed.catalog.id_version.version, "1.2.3")

  def test_update_referenced_assets_preserves_fallback_for_uninstalled(
      self,
  ):
    pa = process_asset_pb2.ProcessAsset()
    pa.behavior_tree.root.CopyFrom(
        _create_task_node("ai.intrinsic.skill_missing")
    )
    fallback = process_asset_pb2.ProcessAsset.ProcessedAsset()
    fallback.catalog.asset_type = asset_type_pb2.AssetType.ASSET_TYPE_SKILL
    fallback.catalog.id_version.id.package = "ai.intrinsic"
    fallback.catalog.id_version.id.name = "skill_missing"
    fallback.catalog.id_version.version = "0.9.0"
    pa.assets["ai.intrinsic.skill_missing"].CopyFrom(fallback)

    error = grpc.RpcError("skill not found")
    error.code = lambda: grpc.StatusCode.NOT_FOUND
    self._installed_assets.batch_get_installed_assets.side_effect = error

    referenced_assets.update_referenced_assets(pa, self._installed_assets)
    self.assertIn("ai.intrinsic.skill_missing", pa.assets)
    self.assertEqual(
        pa.assets["ai.intrinsic.skill_missing"].catalog.id_version.version,
        "0.9.0",
    )

  def test_update_referenced_assets_removes_stale_entries(self):
    pa = process_asset_pb2.ProcessAsset()
    pa.behavior_tree.root.CopyFrom(
        _create_task_node("ai.intrinsic.skill_active")
    )
    pa.assets["ai.intrinsic.skill_stale"].catalog.id_version.version = "1.0.0"

    self._installed_assets.batch_get_installed_assets.return_value = [
        _create_installed_asset(
            "ai.intrinsic",
            "skill_active",
            "2.0.0",
            asset_type_pb2.AssetType.ASSET_TYPE_SKILL,
        )
    ]
    referenced_assets.update_referenced_assets(pa, self._installed_assets)
    self.assertNotIn("ai.intrinsic.skill_stale", pa.assets)
    self.assertIn("ai.intrinsic.skill_active", pa.assets)

  def test_update_referenced_assets_batch_error_non_not_found_raises(self):
    pa = process_asset_pb2.ProcessAsset()
    pa.behavior_tree.root.CopyFrom(_create_task_node("ai.intrinsic.skill_0"))

    error = grpc.RpcError("permission denied")
    error.code = lambda: grpc.StatusCode.PERMISSION_DENIED
    self._installed_assets.batch_get_installed_assets.side_effect = error

    with self.assertRaises(grpc.RpcError) as e:
      referenced_assets.update_referenced_assets(pa, self._installed_assets)
    self.assertEqual(e.exception.code(), grpc.StatusCode.PERMISSION_DENIED)


if __name__ == "__main__":
  absltest.main()
