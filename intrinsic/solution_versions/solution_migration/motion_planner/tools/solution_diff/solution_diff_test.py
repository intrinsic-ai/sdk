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

"""Tests for solution_diff."""

import os
from unittest import mock

from absl import app
from absl.testing import absltest
from absl.testing import flagsaver
from google.protobuf import descriptor_pb2
from google.protobuf import text_encoding
from google.protobuf import text_format
import grpc

from intrinsic.assets.proto.v1 import solution_pb2
from intrinsic.solution_versions.proto.v1 import branch_pb2
from intrinsic.solution_versions.proto.v1 import snapshot_pb2
from intrinsic.solution_versions.proto.v1 import solution_version_service_pb2
from intrinsic.solution_versions.solution_migration.motion_planner.tools.solution_diff import solution_diff
from intrinsic.util.grpc import auth
from intrinsic.util.grpc import dialerutil

_FIELD = descriptor_pb2.FieldDescriptorProto


class _FakeRpcError(grpc.RpcError):
  """Fake `grpc.RpcError` with configurable status code and details."""

  def __init__(self, code: grpc.StatusCode, details: str):
    super().__init__(details)
    self._code = code
    self._details = details

  def code(self) -> grpc.StatusCode:
    return self._code

  def details(self) -> str:
    return self._details


def _file_descriptor_set(
    package: str, message: str, deps: tuple[str, ...] = ()
) -> descriptor_pb2.FileDescriptorSet:
  """Returns a set with one file that defines `<package>.<message>` with `string name = 1`."""
  file = descriptor_pb2.FileDescriptorProto(
      name=f"{package}/{message}.proto",
      package=package,
      syntax="proto3",
      dependency=list(deps),
      message_type=[
          descriptor_pb2.DescriptorProto(
              name=message,
              field=[
                  _FIELD(
                      name="name",
                      number=1,
                      type=_FIELD.TYPE_STRING,
                      label=_FIELD.LABEL_OPTIONAL,
                  )
              ],
          )
      ],
  )
  return descriptor_pb2.FileDescriptorSet(file=[file])


def _any_textproto(indent: str, url: str, value: bytes) -> str:
  escaped = text_encoding.CEscape(value, as_utf8=False)
  return f'{indent}type_url: "{url}"\n{indent}value: "{escaped}"\n'


# Wire encoding of `name: "motion_planner_service"` (field 1, length 22).
_VALUE = b"\n\x16motion_planner_service"
# Wire encoding of `dep { name: "motion_planner_service" }` (field 2, length 24).
_VALUE_WITH_DEP = b"\x12\x18" + _VALUE


def _params_with_dep() -> descriptor_pb2.FileDescriptorSet:
  """Returns p.Params with `string name = 1` and `Dep dep = 2`, where Dep has `string name = 1`."""
  fds = _file_descriptor_set("p", "Params")
  file = fds.file[0]
  file.message_type.add().CopyFrom(file.message_type[0])
  file.message_type[1].name = "Dep"
  file.message_type[0].field.add(
      name="dep",
      number=2,
      type=_FIELD.TYPE_MESSAGE,
      label=_FIELD.LABEL_OPTIONAL,
      type_name=".p.Dep",
  )
  return fds


class FindCatalogAssetsTest(absltest.TestCase):

  def test_finds_both_reference_forms(self):
    textproto = """
      catalog {
        id {
          package: "ai.intrinsic"
          name: "move_robot"
        }
        version: "0.20260831.0-RC01"
      }
      catalog {
        asset_type: ASSET_TYPE_SKILL
        id_version {
          id {
            package: "ai.intrinsic"
            name: "plan_grasp"
          }
          version: "0.20260924.0-RC01"
        }
      }
    """
    self.assertEqual(
        solution_diff.find_catalog_assets(textproto),
        {
            "ai.intrinsic.move_robot.0.20260831.0-RC01": (
                "ai.intrinsic.move_robot"
            ),
            "ai.intrinsic.plan_grasp.0.20260924.0-RC01": (
                "ai.intrinsic.plan_grasp"
            ),
        },
    )


class ResolveOrgInfoTest(absltest.TestCase):

  def test_parses_org_with_project(self):
    org_info = solution_diff.resolve_org_info("my_org@my_project")
    self.assertEqual(
        org_info, auth.OrgInfo(organization="my_org", project="my_project")
    )

  def test_raises_usage_error_on_invalid_org_string(self):
    with self.assertRaises(app.UsageError):
      solution_diff.resolve_org_info("a@b@c")

  def test_raises_usage_error_when_org_not_found(self):
    with mock.patch.object(
        auth,
        "read_org_info",
        side_effect=auth.OrgNotFoundError("not found", "unknown_org"),
    ):
      with self.assertRaises(app.UsageError):
        solution_diff.resolve_org_info("unknown_org")


class CreateChannelTest(absltest.TestCase):

  def test_creates_channel_with_org_info(self):
    fake_channel = mock.MagicMock()
    with mock.patch.object(
        dialerutil, "create_channel_from_org", return_value=fake_channel
    ) as mock_create:
      got = solution_diff.create_channel("my_org@my_project")

    self.assertIs(got, fake_channel)
    mock_create.assert_called_once_with(
        auth.OrgInfo(organization="my_org", project="my_project"),
        grpc_options=solution_diff._GRPC_OPTIONS,
    )

  def test_raises_usage_error_when_credentials_not_found(self):
    with mock.patch.object(
        dialerutil,
        "create_channel_from_org",
        side_effect=auth.CredentialsNotFoundError("not found", "my_project"),
    ):
      with self.assertRaises(app.UsageError):
        solution_diff.create_channel("my_org@my_project")

  def test_raises_usage_error_when_default_token_missing(self):
    with mock.patch.object(
        dialerutil,
        "create_channel_from_org",
        side_effect=KeyError("default"),
    ):
      with self.assertRaises(app.UsageError):
        solution_diff.create_channel("my_org@my_project")


class ResolveSnapshotIdsTest(absltest.TestCase):

  def test_returns_explicit_snapshot_ids_without_rpc(self):
    stub = mock.MagicMock()
    before_id, after_id = solution_diff.resolve_snapshot_ids(
        stub, "branch_1", "snap_before", "snap_after"
    )
    self.assertEqual((before_id, after_id), ("snap_before", "snap_after"))
    stub.GetSnapshot.assert_not_called()

  def test_resolves_default_tip_and_parent_snapshot_ids(self):
    stub = mock.MagicMock()
    stub.GetSnapshot.return_value = snapshot_pb2.Snapshot(
        id="snap_tip", parents=["snap_parent"]
    )
    before_id, after_id = solution_diff.resolve_snapshot_ids(stub, "branch_1")
    self.assertEqual((before_id, after_id), ("snap_parent", "snap_tip"))
    stub.GetSnapshot.assert_called_once_with(
        solution_version_service_pb2.GetSnapshotRequest(
            branch_id="branch_1", snapshot_id=""
        )
    )

  def test_resolves_tip_when_only_before_snapshot_id_is_set(self):
    stub = mock.MagicMock()
    stub.GetSnapshot.return_value = snapshot_pb2.Snapshot(
        id="snap_tip", parents=[]
    )
    before_id, after_id = solution_diff.resolve_snapshot_ids(
        stub, "branch_1", before_snapshot_id="snap_before"
    )
    self.assertEqual((before_id, after_id), ("snap_before", "snap_tip"))
    stub.GetSnapshot.assert_called_once_with(
        solution_version_service_pb2.GetSnapshotRequest(
            branch_id="branch_1", snapshot_id=""
        )
    )

  def test_resolves_parent_when_only_after_snapshot_id_is_set(self):
    stub = mock.MagicMock()
    stub.GetSnapshot.return_value = snapshot_pb2.Snapshot(
        id="snap_after", parents=["snap_parent"]
    )
    before_id, after_id = solution_diff.resolve_snapshot_ids(
        stub, "branch_1", after_snapshot_id="snap_after"
    )
    self.assertEqual((before_id, after_id), ("snap_parent", "snap_after"))
    stub.GetSnapshot.assert_called_once_with(
        solution_version_service_pb2.GetSnapshotRequest(
            branch_id="branch_1", snapshot_id="snap_after"
        )
    )

  def test_falls_back_to_upstream_branch_when_after_snapshot_not_on_deployment_branch(
      self,
  ):
    stub = mock.MagicMock()
    stub.GetSnapshot.side_effect = [
        _FakeRpcError(grpc.StatusCode.NOT_FOUND, "snapshot not found"),
        snapshot_pb2.Snapshot(id="snap_after", parents=["snap_parent"]),
    ]
    stub.GetBranch.return_value = branch_pb2.Branch(
        id="branch_1", upstream_branch_id="upstream_branch_1"
    )

    before_id, after_id = solution_diff.resolve_snapshot_ids(
        stub, "branch_1", after_snapshot_id="snap_after"
    )

    self.assertEqual((before_id, after_id), ("snap_parent", "snap_after"))
    self.assertEqual(
        stub.GetSnapshot.call_args_list,
        [
            mock.call(
                solution_version_service_pb2.GetSnapshotRequest(
                    branch_id="branch_1", snapshot_id="snap_after"
                )
            ),
            mock.call(
                solution_version_service_pb2.GetSnapshotRequest(
                    branch_id="upstream_branch_1", snapshot_id="snap_after"
                )
            ),
        ],
    )

  def test_raises_usage_error_when_branch_has_no_tip_snapshot(self):
    stub = mock.MagicMock()
    stub.GetSnapshot.return_value = snapshot_pb2.Snapshot(id="", parents=[])
    with self.assertRaises(app.UsageError):
      solution_diff.resolve_snapshot_ids(stub, "branch_1")

  def test_raises_usage_error_when_snapshot_has_no_parent(self):
    stub = mock.MagicMock()
    stub.GetSnapshot.return_value = snapshot_pb2.Snapshot(
        id="snap_tip", parents=[]
    )
    with self.assertRaises(app.UsageError):
      solution_diff.resolve_snapshot_ids(stub, "branch_1")


class GetSnapshotTextprotoTest(absltest.TestCase):

  def test_formats_solution_with_unexpanded_any(self):
    url = "type.intrinsic.ai/assets/ai.intrinsic.skill/p.Params"
    solution = text_format.Parse(
        f"""
        assets {{
          key: "main_process"
          value {{
            local {{
              process {{
                behavior_tree {{
                  root {{
                    task {{
                      call_behavior {{
                        skill_id: "ai.intrinsic.skill"
                        parameters {{
                          {_any_textproto("                          ", url, _VALUE).strip()}
                        }}
                      }}
                    }}
                  }}
                }}
              }}
            }}
          }}
        }}
        """,
        solution_pb2.Solution(),
    )
    stub = mock.MagicMock()
    stub.GetSolution.return_value = solution

    textproto = solution_diff.get_snapshot_textproto(stub, "branch_1", "snap_1")

    stub.GetSolution.assert_called_once()
    req = stub.GetSolution.call_args[0][0]
    self.assertEqual(req.source.branch_id, "branch_1")
    self.assertEqual(req.source.snapshot_id, "snap_1")
    self.assertEqual(
        req.view,
        solution_version_service_pb2.GetSolutionRequest.VIEW_FULL_NO_FDS,
    )
    decoder = solution_diff.Decoder(
        {"ai.intrinsic.skill": _file_descriptor_set("p", "Params")}
    )
    self.assertIn('name: "motion_planner_service"', decoder.decode(textproto))


class FetchDescriptorsTest(absltest.TestCase):

  def test_fetches_descriptors_and_retries_chunk_on_resource_exhausted(self):
    resp_cls = (
        solution_version_service_pb2.BatchGetAssetFileDescriptorSetsResponse
    )
    fds = _file_descriptor_set("p", "Params")

    def fake_rpc(request):
      if len(request.asset_id_versions) > 1:
        raise _FakeRpcError(
            grpc.StatusCode.RESOURCE_EXHAUSTED, "message larger than max"
        )
      id_version = request.asset_id_versions[0]
      if id_version == "ai.intrinsic.good.1.0.0":
        resp = resp_cls()
        item = resp.results.add(asset_id_version=id_version)
        item.file_descriptor_set.CopyFrom(fds)
        return resp
      resp = resp_cls()
      item = resp.results.add(asset_id_version=id_version)
      item.error_status.code = 5
      item.error_status.message = "not found"
      return resp

    stub = mock.MagicMock()
    stub.BatchGetAssetFileDescriptorSets.side_effect = fake_rpc

    got = solution_diff.fetch_descriptors(
        stub,
        {
            "ai.intrinsic.good.1.0.0": "ai.intrinsic.good",
            "ai.intrinsic.missing.1.0.0": "ai.intrinsic.missing",
        },
    )

    self.assertEqual(got, {"ai.intrinsic.good": fds})

  def test_propagates_non_size_rpc_error_without_retrying_items(self):
    stub = mock.MagicMock()
    stub.BatchGetAssetFileDescriptorSets.side_effect = _FakeRpcError(
        grpc.StatusCode.UNAUTHENTICATED, "unauthenticated"
    )

    with self.assertRaises(grpc.RpcError):
      solution_diff.fetch_descriptors(
          stub,
          {
              "ai.intrinsic.a.1.0.0": "ai.intrinsic.a",
              "ai.intrinsic.b.1.0.0": "ai.intrinsic.b",
          },
      )

    stub.BatchGetAssetFileDescriptorSets.assert_called_once()


class BuildPoolTest(absltest.TestCase):

  def test_adds_files_after_their_dependencies(self):
    base = _file_descriptor_set("base", "Base").file[0]
    child = _file_descriptor_set("child", "Child", deps=(base.name,)).file[0]
    pool, skipped = solution_diff.build_pool(
        descriptor_pb2.FileDescriptorSet(file=[child, base])
    )
    self.assertEmpty(skipped)
    self.assertIsNotNone(pool.FindMessageTypeByName("child.Child"))

  def test_reports_files_with_missing_dependencies(self):
    child = _file_descriptor_set(
        "child", "Child", deps=("missing.proto",)
    ).file[0]
    _, skipped = solution_diff.build_pool(
        descriptor_pb2.FileDescriptorSet(file=[child])
    )
    self.assertEqual(skipped, [child.name])


class DecoderTest(absltest.TestCase):

  def test_expands_any_with_the_named_assets_descriptors(self):
    decoder = solution_diff.Decoder(
        {"ai.intrinsic.skill": _file_descriptor_set("p", "Params")}
    )
    url = "type.intrinsic.ai/assets/ai.intrinsic.skill/p.Params"

    out = decoder.decode(_any_textproto("  ", url, _VALUE))

    self.assertEqual(
        out, f'  [{url}] {{\n    name: "motion_planner_service"\n  }}\n'
    )
    self.assertEmpty(decoder.unresolved)

  def test_keeps_value_when_named_asset_has_no_descriptor(self):
    # Another Asset defines p.Params, but the URL names a different Asset.
    decoder = solution_diff.Decoder(
        {"ai.intrinsic.other": _file_descriptor_set("p", "Params")}
    )
    url = "type.intrinsic.ai/assets/ai.intrinsic.skill/p.Params"
    textproto = _any_textproto("", url, _VALUE)

    self.assertEqual(decoder.decode(textproto), textproto)
    self.assertEqual(decoder.unresolved, {url})

  def test_googleapis_url_uses_any_asset_that_defines_the_message(self):
    decoder = solution_diff.Decoder(
        {"ai.intrinsic.skill": _file_descriptor_set("p", "Params")}
    )
    url = "type.googleapis.com/p.Params"

    out = decoder.decode(_any_textproto("", url, _VALUE))

    self.assertIn('name: "motion_planner_service"', out)

  def test_keeps_value_that_does_not_decode(self):
    decoder = solution_diff.Decoder(
        {"ai.intrinsic.skill": _file_descriptor_set("p", "Params")}
    )
    url = "type.intrinsic.ai/assets/ai.intrinsic.skill/p.Params"
    textproto = _any_textproto("", url, b"\n\x16short")  # Length past the end.

    self.assertEqual(decoder.decode(textproto), textproto)
    self.assertEqual(decoder.unresolved, {url})

  def test_googleapis_url_uses_the_schema_of_the_called_skill(self):
    # Both Assets define p.Params, but only the called Skill has the `dep` field.
    decoder = solution_diff.Decoder({
        "ai.intrinsic.other": _file_descriptor_set("p", "Params"),
        "ai.intrinsic.skill": _params_with_dep(),
    })
    url = "type.googleapis.com/p.Params"
    textproto = (
        'call_behavior {\n  skill_id: "ai.intrinsic.skill"\n  parameters {\n'
        + _any_textproto("    ", url, _VALUE_WITH_DEP)
        + "  }\n}\n"
    )

    out = decoder.decode(textproto)

    self.assertIn('      dep {\n        name: "motion_planner_service"\n', out)

  def test_skill_id_of_an_earlier_call_is_not_used(self):
    decoder = solution_diff.Decoder({
        "ai.intrinsic.other": _file_descriptor_set("p", "Params"),
        "ai.intrinsic.skill": _params_with_dep(),
    })
    url = "type.googleapis.com/p.Params"
    textproto = (
        "call_behavior {\n"
        '  skill_id: "ai.intrinsic.skill"\n'
        "}\n"
        "call_behavior {\n"
        "  parameters {\n"
        + _any_textproto("    ", url, _VALUE_WITH_DEP)
        + "  }\n}\n"
    )

    out = decoder.decode(textproto)

    self.assertNotIn("dep {", out)

  def test_prints_fields_that_the_schema_does_not_define(self):
    decoder = solution_diff.Decoder(
        {"ai.intrinsic.other": _file_descriptor_set("p", "Params")}
    )
    url = "type.googleapis.com/p.Params"

    out = decoder.decode(_any_textproto("", url, _VALUE_WITH_DEP))

    self.assertIn('2 {\n    1: "motion_planner_service"\n  }', out)


class SplitSkillCallHunksTest(absltest.TestCase):

  def test_keeps_small_skill_call_hunks_and_counts_large_ones(self):
    skill = "@@ -1,2 +1,2 @@\n call_behavior {\n-a\n+b\n"
    other = "@@ -9,2 +9,2 @@\n instances {\n-a\n+b\n"
    large = '@@ -20,300 +20,300 @@\n skill_id: "x"\n' + " line\n" * 300
    diff = "--- a\n+++ b\n" + skill + other + large

    keep, too_large = solution_diff.split_skill_call_hunks(diff)

    self.assertEqual(keep, [skill])
    self.assertEqual(too_large, 1)

  def test_keeps_hunk_whose_skill_id_is_outside_the_context(self):
    before = (
        ["call_behavior {\n", '  skill_id: "ai.intrinsic.plan_grasp"\n']
        + ["  parameters {\n"]
        + ["    x: 1\n"] * 40
        + ["  }\n", "}\n", "other {\n", "  y: 1\n", "}\n"]
    )
    inside = "@@ -40,2 +40,3 @@\n     x: 1\n+    dep {}\n     x: 1\n"
    outside = "@@ -47,2 +47,3 @@\n other {\n+  z: 1\n   y: 1\n"
    diff = "--- a\n+++ b\n" + inside + outside

    keep, _ = solution_diff.split_skill_call_hunks(diff, before)

    self.assertEqual(keep, [inside])


class MainTest(absltest.TestCase):

  def test_writes_decoded_snapshots_and_diffs(self):
    out_dir = self.create_tempdir().full_path
    url = "type.intrinsic.ai/assets/ai.intrinsic.skill/p.Params"
    before_text = (
        'call_behavior {\n  skill_id: "ai.intrinsic.skill"\n  parameters {\n'
        + _any_textproto("    ", url, b"")
        + "  }\n}\n"
    )
    after_text = (
        'call_behavior {\n  skill_id: "ai.intrinsic.skill"\n  parameters {\n'
        + _any_textproto("    ", url, _VALUE)
        + "  }\n}\n"
    )
    with (
        flagsaver.flagsaver(
            solution="branch_1",
            org="my_org@my_project",
            before_snapshot_id="snap_before",
            after_snapshot_id="snap_after",
            out=out_dir,
        ),
        mock.patch.object(solution_diff, "create_channel"),
        mock.patch.object(
            solution_diff,
            "get_snapshot_textproto",
            side_effect=[before_text, after_text],
        ),
        mock.patch.object(
            solution_diff,
            "fetch_descriptors",
            return_value={
                "ai.intrinsic.skill": _file_descriptor_set("p", "Params")
            },
        ),
    ):
      solution_diff.main(["solution_diff"])

    with open(os.path.join(out_dir, "after.textproto"), encoding="utf-8") as f:
      self.assertIn('name: "motion_planner_service"', f.read())
    with open(os.path.join(out_dir, "bt_changes.diff"), encoding="utf-8") as f:
      self.assertIn('+      name: "motion_planner_service"', f.read())

  def test_converts_rpc_error_to_usage_error(self):
    out_dir = self.create_tempdir().full_path
    with (
        flagsaver.flagsaver(
            solution="bad_branch",
            org="my_org@my_project",
            out=out_dir,
        ),
        mock.patch.object(solution_diff, "create_channel"),
        mock.patch.object(
            solution_diff,
            "resolve_snapshot_ids",
            side_effect=_FakeRpcError(
                grpc.StatusCode.NOT_FOUND, "branch not found"
            ),
        ),
    ):
      with self.assertRaisesRegex(
          app.UsageError, "NOT_FOUND: branch not found"
      ):
        solution_diff.main(["solution_diff"])

  def test_rejects_unexpected_positional_arguments(self):
    with self.assertRaises(app.UsageError):
      solution_diff.main(["solution_diff", "extra_arg"])


if __name__ == "__main__":
  absltest.main()
