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

"""Readable diff of two Solution snapshots, with Skill parameters decoded.

Shows the changes between two snapshots with field names and values. For
example, use it to see what `inctl solution migrate --migration motion-planner`
changed and port the change into SBL Python scripts.

Usage:
  # From a Bazel workspace:
  bazel run \\
    //intrinsic/solution_versions/solution_migration/motion_planner/tools/solution_diff \\
    -- --solution=<BRANCH_ID> --org=<ORG> \\
    [--before_snapshot_id=<SNAP>] [--after_snapshot_id=<SNAP>] [--out=DIR]

  # From a Python environment with the Intrinsic solutions wheel installed:
  python3 -m \\
    intrinsic.solution_versions.solution_migration.motion_planner.tools.solution_diff.solution_diff \\
    --solution=<BRANCH_ID> --org=<ORG>
"""

import bisect
from collections.abc import Mapping
from collections.abc import Sequence
import difflib
import itertools
import os
import re
import shutil
import subprocess
from typing import Any

from absl import app
from absl import flags
from google.protobuf import descriptor_pb2
from google.protobuf import descriptor_pool
from google.protobuf import message
from google.protobuf import message_factory
from google.protobuf import text_encoding
from google.protobuf import text_format
import grpc

from intrinsic.solution_versions.proto.v1 import snapshot_source_pb2
from intrinsic.solution_versions.proto.v1 import solution_version_service_pb2
from intrinsic.solution_versions.proto.v1 import solution_version_service_pb2_grpc
from intrinsic.util.grpc import auth
from intrinsic.util.grpc import dialerutil

_SOLUTION = flags.DEFINE_string(
    "solution",
    None,
    "Deployment branch ID of the Solution (..._BRANCH).",
)
_BEFORE_SNAPSHOT_ID = flags.DEFINE_string(
    "before_snapshot_id",
    "",
    "Snapshot ID of the earlier version (..._SNAP). Defaults to the primary"
    " parent of --after_snapshot_id.",
)
_AFTER_SNAPSHOT_ID = flags.DEFINE_string(
    "after_snapshot_id",
    "",
    "Snapshot ID of the later version (..._SNAP). Defaults to the tip"
    " snapshot of --solution.",
)
_ORG = flags.DEFINE_string(
    "org",
    None,
    "Organization name (for example, my_org or my_org@my_project).",
)
_OUT = flags.DEFINE_string(
    "out",
    "solution_diff_out",
    "Output folder. A relative path is relative to the folder you run from.",
)

_GRPC_OPTIONS = [
    ("grpc.max_receive_message_length", -1),
    ("grpc.max_send_message_length", -1),
]
# The catalog returns at most about 4 MB per call; 8 Assets per call stays
# under that.
_FDS_CHUNK_SIZE = 8
# Hunks longer than this are moved Asset map entries or whole copied Processes.
_MAX_FOCUSED_HUNK_LINES = 200
_CONTEXT_LINES = 14

# Catalog Asset references in a Solution or ModifiedSolution textproto.
_ASSET_RE = re.compile(
    r"catalog \{\s*(?:asset_type: \w+\s*)?(?:id_version \{\s*)?id \{\s*"
    r'package: "([^"]+)"\s*name: "([^"]+)"\s*\}\s*version: "([^"]+)"'
)
# An Any in textproto form: a type_url line followed by a value line at the
# same indent.
_ANY_RE = re.compile(
    r'^( *)type_url: "([^"]+)"\n\1value: "((?:[^"\\]|\\.)*)"\n', re.M
)
_ASSET_URL_RE = re.compile(r"/assets/([^/]+)/")
_SKILL_ID_RE = re.compile(r'^( *)skill_id: "([^"]+)"$', re.M)
_HUNK_SPLIT_RE = re.compile(r"(?m)^(?=@@ )")
_HUNK_HEADER_RE = re.compile(r"@@ -(\d+)")


def find_catalog_assets(textproto: str) -> dict[str, str]:
  """Returns `{"<package>.<name>.<version>": "<package>.<name>"}` for Assets."""
  return {
      f"{package}.{name}.{version}": f"{package}.{name}"
      for package, name, version in _ASSET_RE.findall(textproto)
  }


def resolve_org_info(org: str) -> auth.OrgInfo:
  """Resolves organization info from `<org>@<project>` or saved `inctl` config.

  Args:
    org: Organization identifier (`<org>@<project>` or `<org>`).

  Returns:
    The resolved `auth.OrgInfo`.

  Raises:
    app.UsageError: If the organization cannot be parsed or found locally.
  """
  if "@" in org:
    try:
      return auth.parse_info_from_string(org)
    except ValueError as err:
      raise app.UsageError(str(err)) from err
  try:
    return auth.read_org_info(org)
  except auth.OrgNotFoundError as err:
    raise app.UsageError(
        f"Credentials for organization '{org}' not found. Run "
        f"'inctl auth login --org {org}' first."
    ) from err


def create_channel(org: str) -> grpc.Channel:
  """Opens an authenticated gRPC channel to the cloud organization."""
  org_info = resolve_org_info(org)
  try:
    return dialerutil.create_channel_from_org(
        org_info, grpc_options=_GRPC_OPTIONS
    )
  except (auth.CredentialsNotFoundError, KeyError) as err:
    raise app.UsageError(
        f"Credentials for project '{org_info.project}' not found. Run "
        f"'inctl auth login --org {org}' first."
    ) from err


def _rpc_error_message(err: grpc.RpcError) -> str:
  """Returns a concise human-readable message for a `grpc.RpcError`."""
  details = err.details() if hasattr(err, "details") else ""
  code = err.code() if hasattr(err, "code") else None
  if details and code is not None:
    return f"{code.name}: {details}"
  return details or str(err)


def _is_batch_too_large(err: grpc.RpcError) -> bool:
  """Returns True if `err` indicates a batch request exceeded message limits."""
  code = err.code() if hasattr(err, "code") else None
  details = err.details() if hasattr(err, "details") else ""
  return code == grpc.StatusCode.RESOURCE_EXHAUSTED or "ResourceExhausted" in (
      details or str(err)
  )


def resolve_snapshot_ids(
    stub: Any,
    branch_id: str,
    before_snapshot_id: str = "",
    after_snapshot_id: str = "",
) -> tuple[str, str]:
  """Returns `(before_snapshot_id, after_snapshot_id)`, resolving missing IDs.

  When `after_snapshot_id` is empty, uses the tip snapshot of `branch_id`. When
  `before_snapshot_id` is empty, uses the primary parent of the `after`
  snapshot. If `after_snapshot_id` is on the upstream version branch rather than
  the deployment branch's history, falls back to `GetBranch.upstream_branch_id`.

  Args:
    stub: A `SolutionVersionServiceStub` (or compatible test double).
    branch_id: Deployment branch ID of the Solution.
    before_snapshot_id: Optional explicit earlier snapshot ID.
    after_snapshot_id: Optional explicit later snapshot ID.

  Returns:
    A `(before_snapshot_id, after_snapshot_id)` tuple.

  Raises:
    app.UsageError: If a missing snapshot ID cannot be resolved from branch
      history.
  """
  if before_snapshot_id and after_snapshot_id:
    return before_snapshot_id, after_snapshot_id
  try:
    after_snapshot = stub.GetSnapshot(
        solution_version_service_pb2.GetSnapshotRequest(
            branch_id=branch_id,
            snapshot_id=after_snapshot_id,
        )
    )
  except grpc.RpcError as err:
    code = err.code() if hasattr(err, "code") else None
    if not after_snapshot_id or code != grpc.StatusCode.NOT_FOUND:
      raise
    branch = stub.GetBranch(
        solution_version_service_pb2.GetBranchRequest(branch_id=branch_id)
    )
    if not branch.upstream_branch_id:
      raise
    after_snapshot = stub.GetSnapshot(
        solution_version_service_pb2.GetSnapshotRequest(
            branch_id=branch.upstream_branch_id,
            snapshot_id=after_snapshot_id,
        )
    )
  resolved_after = after_snapshot_id or after_snapshot.id
  if not resolved_after:
    raise app.UsageError(
        f"Branch '{branch_id}' has no tip snapshot; specify"
        " --after_snapshot_id."
    )
  if before_snapshot_id:
    return before_snapshot_id, resolved_after
  if not after_snapshot.parents:
    raise app.UsageError(
        f"Snapshot '{resolved_after}' has no parent snapshot; specify"
        " --before_snapshot_id."
    )
  return after_snapshot.parents[0], resolved_after


def get_snapshot_textproto(
    stub: Any,
    branch_id: str,
    snapshot_id: str,
) -> str:
  """Fetches a Solution snapshot via gRPC and formats it as a textproto."""
  solution_proto = stub.GetSolution(
      solution_version_service_pb2.GetSolutionRequest(
          source=snapshot_source_pb2.SnapshotSource(
              branch_id=branch_id,
              snapshot_id=snapshot_id,
          ),
          view=(
              solution_version_service_pb2.GetSolutionRequest.VIEW_FULL_NO_FDS
          ),
      )
  )
  return text_format.MessageToString(
      solution_proto,
      descriptor_pool=descriptor_pool.DescriptorPool(),
      print_unknown_fields=True,
  )


def fetch_descriptors(
    stub: Any, assets: Mapping[str, str]
) -> dict[str, descriptor_pb2.FileDescriptorSet]:
  """Returns `{asset_id: FileDescriptorSet}` for `{id_version: asset_id}`."""
  req_cls = solution_version_service_pb2.BatchGetAssetFileDescriptorSetsRequest
  id_versions = sorted(assets)
  results = []
  errors: dict[str, str] = {}
  for i in range(0, len(id_versions), _FDS_CHUNK_SIZE):
    chunk = id_versions[i : i + _FDS_CHUNK_SIZE]
    try:
      results.extend(
          stub.BatchGetAssetFileDescriptorSets(
              req_cls(asset_id_versions=chunk)
          ).results
      )
      continue
    except grpc.RpcError as err:
      if len(chunk) <= 1 or not _is_batch_too_large(err):
        raise
    for one in chunk:  # Retry one by one when the batch response is too large.
      try:
        results.extend(
            stub.BatchGetAssetFileDescriptorSets(
                req_cls(asset_id_versions=[one])
            ).results
        )
      except grpc.RpcError as err:
        if not _is_batch_too_large(err):
          raise
        errors[one] = _rpc_error_message(err)

  out: dict[str, descriptor_pb2.FileDescriptorSet] = {}
  for item in results:
    id_version = item.asset_id_version
    if id_version not in assets:
      continue
    which = item.WhichOneof("result")
    if which == "file_descriptor_set":
      out[assets[id_version]] = item.file_descriptor_set
    elif which == "error_status":
      errors[id_version] = (
          item.error_status.message or f"code {item.error_status.code}"
      )
  print(f"{len(out)} descriptor sets fetched, {len(errors)} errors")
  for id_version, err in list(errors.items())[:5]:
    print(f"  {id_version}: {err}")
  return out


def build_pool(
    fds: descriptor_pb2.FileDescriptorSet,
) -> tuple[descriptor_pool.DescriptorPool, list[str]]:
  """Returns a pool with `fds` files and names of files that failed to add."""
  pool = descriptor_pool.DescriptorPool()
  pending = list(fds.file)
  while pending:  # Add files after their dependencies, in any input order.
    still = []
    for file_proto in pending:
      try:
        pool.Add(file_proto)
      except (TypeError, ValueError):
        still.append(file_proto)
    if len(still) == len(pending):
      break
    pending = still
  return pool, [file_proto.name for file_proto in pending]


class Decoder:
  """Expands Any values in a textproto using Asset descriptors."""

  def __init__(
      self,
      descriptor_sets: Mapping[str, descriptor_pb2.FileDescriptorSet],
  ):
    """Initializes the decoder with `{asset_id: FileDescriptorSet}`."""
    self._sets = descriptor_sets
    self._pools: dict[str, descriptor_pool.DescriptorPool | None] = {}
    self.unresolved: set[str] = set()
    self.skipped_files: set[str] = set()

  def _pool(self, asset_id: str) -> descriptor_pool.DescriptorPool | None:
    if asset_id not in self._pools:
      pool = None
      if asset_id in self._sets:
        pool, skipped = build_pool(self._sets[asset_id])
        self.skipped_files.update(skipped)
      self._pools[asset_id] = pool
    return self._pools[asset_id]

  def _message_type(
      self, url: str, skill_id: str | None = None
  ) -> tuple[
      type[message.Message] | None, descriptor_pool.DescriptorPool | None
  ]:
    """Returns `(message class, pool)`, or `(None, None)` if not found.

    A URL that names an Asset uses only that Asset's descriptors, so a message
    name that several Assets define is not decoded with the wrong schema. Other
    URLs (`type.googleapis.com`) try the Skill that owns the parameters first,
    then any other Asset that defines the message.

    Args:
      url: Type URL of the `google.protobuf.Any` message.
      skill_id: Optional `skill_id` of the enclosing Skill call.

    Returns:
      A `(message_class, descriptor_pool)` tuple, or `(None, None)` when no
      matching descriptor is found.
    """
    name = url.rsplit("/", 1)[-1]
    url_match = _ASSET_URL_RE.search(url)
    if url_match:
      candidates = [url_match.group(1)]
    else:
      # Other Assets can define an older message with the same name. For
      # example, grasp_object has a MoveRobotParams without the
      # motion_planner_service field that move_robot has.
      first = [skill_id] if skill_id and skill_id in self._sets else []
      candidates = first + [a for a in self._sets if a != skill_id]
    for asset_id in candidates:
      pool = self._pool(asset_id)
      if pool is None:
        continue
      try:
        desc = pool.FindMessageTypeByName(name)
      except KeyError:
        continue
      return message_factory.GetMessageClass(desc), pool
    return None, None

  def expand(self, match: re.Match[str], skill_id: str | None = None) -> str:
    """Replaces one `type_url`/`value` pair with the decoded message."""
    indent, url, escaped = match.group(1), match.group(2), match.group(3)
    cls, pool = self._message_type(url, skill_id)
    if cls is None:
      self.unresolved.add(url)
      return match.group(0)
    try:
      msg = cls.FromString(text_encoding.CUnescape(escaped))
    except message.DecodeError:
      self.unresolved.add(url)
      return match.group(0)
    # Print unknown fields, so a schema that lacks a field still shows its
    # number and value.
    body = text_format.MessageToString(
        msg,
        indent=len(indent) + 2,
        descriptor_pool=pool,
        print_unknown_fields=True,
    )
    return f"{indent}[{url}] {{\n{body}{indent}}}\n"

  def decode(self, textproto: str) -> str:
    """Decodes all resolvable `google.protobuf.Any` blocks in `textproto`."""
    lines = textproto.splitlines(keepends=True)
    line_ends = list(itertools.accumulate(len(line) for line in lines))

    def expand_in_call(any_match: re.Match[str]) -> str:
      index = bisect.bisect_right(line_ends, any_match.start())
      return self.expand(
          any_match, enclosing_skill_id(lines, index, len(any_match.group(1)))
      )

    return _ANY_RE.sub(expand_in_call, textproto)


def unified_diff(before_path: str, after_path: str, context: int) -> str:
  """Returns `diff -u` of the files `before_path` and `after_path`."""
  if shutil.which("diff"):
    return subprocess.run(
        ["diff", f"-U{context}", before_path, after_path],
        capture_output=True,
        encoding="utf-8",
        check=False,
    ).stdout
  with (
      open(before_path, encoding="utf-8") as before_file,
      open(after_path, encoding="utf-8") as after_file,
  ):
    return "".join(
        difflib.unified_diff(
            before_file.readlines(),
            after_file.readlines(),
            before_path,
            after_path,
            n=context,
        )
    )


def enclosing_skill_id(
    lines: Sequence[str], index: int, indent: int
) -> str | None:
  """Returns the `skill_id` of the Skill call enclosing `lines[index]`, or None.

  Walks up through the parent blocks. The `skill_id` is a sibling of one of the
  parents, so it has the same indent as that parent.

  Args:
    lines: Lines of the textproto document.
    index: Zero-based line index of the inner block.
    indent: Indentation (in spaces) of the inner block.

  Returns:
    The enclosing `skill_id` string, or `None` if not inside a Skill call.
  """
  level = indent
  for i in range(index - 1, -1, -1):
    line = lines[i]
    if not line.strip():
      continue
    line_indent = len(line) - len(line.lstrip(" "))
    level = min(level, line_indent)
    if line_indent == level:
      skill_match = _SKILL_ID_RE.match(line.rstrip("\n"))
      if skill_match:
        return skill_match.group(2)
    if level == 0:
      return None
  return None


def _first_change(hunk: str) -> tuple[int, int] | None:
  """Returns `(index in before file, indent)` of the first change in `hunk`."""
  header_match = _HUNK_HEADER_RE.match(hunk)
  if not header_match:
    return None
  index = int(header_match.group(1)) - 1
  for line in hunk.splitlines()[1:]:
    if line[:1] in "-+":
      return index, len(line) - 1 - len(line[1:].lstrip(" "))
    index += 1
  return None


def split_skill_call_hunks(
    diff: str, before: Sequence[str] | None = None
) -> tuple[list[str], int]:
  """Returns hunks touching a Skill call and count of hunks too large to keep.

  With the lines of the `before` file, a hunk also counts if its change is
  inside a Skill call whose `skill_id` is outside the context lines (for
  example, in long Skill parameters).

  Args:
    diff: Unified diff string.
    before: Optional lines of the `before` textproto file.

  Returns:
    A `(kept_hunks, num_too_large)` tuple.
  """
  hunks = _HUNK_SPLIT_RE.split(diff)[1:]

  def touches_call(hunk: str) -> bool:
    if "call_behavior" in hunk or "skill_id:" in hunk:
      return True
    change = _first_change(hunk) if before else None
    return (
        change is not None and enclosing_skill_id(before, *change) is not None
    )

  calls = [hunk for hunk in hunks if touches_call(hunk)]
  keep = [hunk for hunk in calls if hunk.count("\n") <= _MAX_FOCUSED_HUNK_LINES]
  return keep, len(calls) - len(keep)


def main(argv: Sequence[str]) -> None:
  if len(argv) > 1:
    raise app.UsageError(f"Unexpected arguments: {' '.join(argv[1:])}")
  out = os.path.expanduser(_OUT.value)
  # `bazel run` starts in the runfiles tree; resolve a relative --out against
  # the caller's folder.
  if not os.path.isabs(out) and os.environ.get("BUILD_WORKING_DIRECTORY"):
    out = os.path.join(os.environ["BUILD_WORKING_DIRECTORY"], out)
  os.makedirs(out, exist_ok=True)

  try:
    with create_channel(_ORG.value) as channel:
      stub = solution_version_service_pb2_grpc.SolutionVersionServiceStub(
          channel
      )
      before_id, after_id = resolve_snapshot_ids(
          stub,
          _SOLUTION.value,
          _BEFORE_SNAPSHOT_ID.value,
          _AFTER_SNAPSHOT_ID.value,
      )
      print(f"== Download snapshots ({before_id} -> {after_id})", flush=True)
      raw = {
          "before": get_snapshot_textproto(stub, _SOLUTION.value, before_id),
          "after": get_snapshot_textproto(stub, _SOLUTION.value, after_id),
      }

      assets: dict[str, str] = {}
      for text in raw.values():
        assets.update(find_catalog_assets(text))
      print(
          f"== Fetch parameter descriptors for {len(assets)} Catalog Assets",
          flush=True,
      )
      decoder = Decoder(fetch_descriptors(stub, assets))
  except grpc.RpcError as err:
    raise app.UsageError(_rpc_error_message(err)) from err

  print("== Decode Skill parameters", flush=True)
  paths = {}
  for label, text in raw.items():
    paths[label] = os.path.join(out, f"{label}.textproto")
    with open(paths[label], "w", encoding="utf-8") as out_file:
      out_file.write(decoder.decode(text))
  print(
      f"{len(decoder.unresolved)} parameter types stay encoded: no descriptor"
      " (for example sideloaded Assets) or the value did not decode"
  )
  if decoder.skipped_files:
    print(f"{len(decoder.skipped_files)} descriptor files could not be loaded")

  print("== Diff", flush=True)
  with open(os.path.join(out, "full.diff"), "w", encoding="utf-8") as out_file:
    out_file.write(unified_diff(paths["before"], paths["after"], 3))
  with open(paths["before"], encoding="utf-8") as before_file:
    before_lines = before_file.readlines()
  keep, too_large = split_skill_call_hunks(
      unified_diff(paths["before"], paths["after"], _CONTEXT_LINES),
      before_lines,
  )
  focused = os.path.join(out, "bt_changes.diff")
  with open(focused, "w", encoding="utf-8") as out_file:
    out_file.write("".join(keep))
  print(f"{len(keep)} Behavior Tree hunks -> {focused}")
  if too_large:
    print(
        f"Large Skill call hunks (over {_MAX_FOCUSED_HUNK_LINES} lines) only in"
        f" {os.path.join(out, 'full.diff')}: {too_large}"
    )


if __name__ == "__main__":
  flags.mark_flags_as_required(["solution", "org"])
  app.run(main)
