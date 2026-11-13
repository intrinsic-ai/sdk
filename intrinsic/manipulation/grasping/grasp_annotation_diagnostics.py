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

"""Helpers for reading and shaping grasp annotation rejection diagnostics.

`GraspAnnotationDiagnostics` (see `grasp_annotator.proto`) reports the
candidate funnel of one grasp annotation run as a list of `StageCount`s, each
with up to `max_num_rejected_annotations_per_stage` rejected candidates
ordered nearest-to-threshold first. These helpers are shared by the producer
(the grasp annotator service) and its consumers so that they agree on the
validation and trimming behavior.
"""

from intrinsic.manipulation.grasping import grasp_annotator_pb2


def trim_rejected_annotations(
    diagnostics: grasp_annotator_pb2.GraspAnnotationDiagnostics,
    max_rejections_per_stage: int,
) -> grasp_annotator_pb2.GraspAnnotationDiagnostics:
  """Returns a copy of `diagnostics` with capped rejections per stage.

  Rejected candidates are stored nearest-to-threshold first, so keeping the
  first `max_rejections_per_stage` entries keeps the nearest misses. Counts,
  thresholds and timings are copied unchanged. `diagnostics` is not modified.

  Args:
    diagnostics: The diagnostics to trim.
    max_rejections_per_stage: Maximum number of `rejected_annotations` to keep
      per stage.

  Returns:
    The trimmed copy.

  Raises:
    ValueError: If `max_rejections_per_stage` is negative.
  """
  if max_rejections_per_stage < 0:
    raise ValueError(
        "max_rejections_per_stage must be non-negative, got"
        f" {max_rejections_per_stage}."
    )
  trimmed = grasp_annotator_pb2.GraspAnnotationDiagnostics()
  trimmed.CopyFrom(diagnostics)
  for stage in trimmed.stage_breakdown:
    del stage.rejected_annotations[max_rejections_per_stage:]
  return trimmed


def get_first_empty_stage(
    diagnostics: grasp_annotator_pb2.GraspAnnotationDiagnostics,
) -> grasp_annotator_pb2.StageCount | None:
  """Returns the first stage that left no candidates, or None.

  Args:
    diagnostics: The diagnostics to inspect.

  Returns:
    The first `StageCount` whose `output_count` is set and equal to 0, i.e. the
    stage that rejected the last remaining candidates, or None if every stage
    kept at least one candidate.
  """
  for stage in diagnostics.stage_breakdown:
    if stage.HasField("output_count") and stage.output_count == 0:
      return stage
  return None
