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

"""Helpers for reading, trimming, and formatting grasp annotation diagnostics.

`GraspAnnotationDiagnostics` (see `grasp_annotator.proto`) reports the
candidate funnel of one grasp annotation run as a list of `StageCount`s, each
with up to `max_num_rejected_annotations_per_stage` rejected candidates
ordered nearest-to-threshold first. These helpers are shared by the producer
(the grasp annotator service) and its consumers (such as `annotate_grasp` and
custom skills) so that they agree on trimming, empty-stage lookup, and
human-readable funnel and stage diagnostics.
"""

import dataclasses
import logging

from intrinsic.manipulation.grasping import grasp_annotator_pb2

_AnnotationCheck = grasp_annotator_pb2.AnnotationCheck

_EXCLUDED_BOXES_KEY = "exclude_poses_in_box"
_INCLUDED_BOXES_KEY = "include_poses_in_box"

# Human-readable labels used by `format_candidate_funnel` to render a stage
# chain such as "100 rays → 96 hit rays → 0 narrower than the opening". Each
# entry maps an `AnnotationCheck` stage to `(entering_label, surviving_label)`:
# `entering_label` qualifies `stage.input_count` at the start of the funnel (or
# whenever the candidate count/unit changes between stages), and
# `surviving_label` qualifies `stage.output_count` for candidates that pass the
# stage.
_STAGE_LABELS: dict[int, tuple[str, str]] = {
    _AnnotationCheck.RAY_MISSED_SURFACE: ("rays", "hit rays"),
    _AnnotationCheck.RAY_SINGLE_HIT: ("hit rays", "multi-hit rays"),
    _AnnotationCheck.INSUFFICIENT_ANTIPODAL_NORMALS: (
        "point pairs",
        "antipodal pairs",
    ),
    _AnnotationCheck.MIN_CLOSING_DISTANCE_VIOLATED: (
        "pairs",
        "wider than the minimum closing distance",
    ),
    _AnnotationCheck.OPENING_EXCEEDED: ("pairs", "narrower than the opening"),
    _AnnotationCheck.NO_CONTACT_SURFACE_IN_FINGER_REACH: (
        "pairs",
        "with a contact surface in finger reach",
    ),
    _AnnotationCheck.FINGER_COLLISION_WITH_PART: (
        "candidate poses",
        "collision-free",
    ),
    _AnnotationCheck.SUCTION_CUP_OVER_HOLE: (
        "sampled poses",
        "without holes under the cups",
    ),
    _AnnotationCheck.SURFACE_ROUGHNESS_EXCEEDED: (
        "poses",
        "within the roughness threshold",
    ),
    _AnnotationCheck.INCOMPATIBLE_OBJECT_SHAPE: (
        "object mesh",
        "cylindrical fit",
    ),
    _AnnotationCheck.INCOMPATIBLE_RADIAL_DIMENSION: (
        "fitted cylinder",
        "within the gripper radii",
    ),
    _AnnotationCheck.BOUNDING_BOX_CONSTRAINT_EXCLUDED: (
        "poses",
        "passing the bounding box constraint",
    ),
    _AnnotationCheck.AXIS_CONSTRAINT_MISALIGNED: (
        "poses",
        "aligned with the axis constraint",
    ),
    _AnnotationCheck.INSUFFICIENT_CONTACT_AREA: (
        "poses",
        "with enough finger contact area",
    ),
}


@dataclasses.dataclass(frozen=True)
class StageDiagnosticText:
  """Human-readable diagnosis and recovery action for an empty funnel stage.

  Attributes:
    diagnosis: Sentence clause (without trailing period) describing what failed
      at the stage, including the nearest-miss measurement and any candidate
      debug note when available.
    action: Capitalized recovery instruction sentence ending with a period.
  """

  diagnosis: str
  action: str


@dataclasses.dataclass(frozen=True)
class _StageTemplate:
  diagnosis: str
  action: str
  nearest_clause: str | None = None


# Keyed by (AnnotationCheck, context_key), where context_key is a gripper-specs
# oneof field name, bounding-box mode, or None for the generic fallback.
_STAGE_TEMPLATES: dict[tuple[int, str | None], _StageTemplate] = {
    (_AnnotationCheck.RAY_MISSED_SURFACE, None): _StageTemplate(
        diagnosis="No ray hit the object mesh",
        action=(
            "Check that the object has non-degenerate, closed (watertight)"
            " visual geometry (for example, no zero-area faces or planar mesh)"
            " and increase `num_samples`."
        ),
    ),
    (_AnnotationCheck.RAY_SINGLE_HIT, None): _StageTemplate(
        diagnosis=(
            "No ray passed through the object twice, so no opposing surface"
            " pairs exist"
        ),
        action=(
            "Check that the object has closed (watertight) visual geometry"
            " with non-zero thickness and increase `num_samples`."
        ),
    ),
    (_AnnotationCheck.INSUFFICIENT_ANTIPODAL_NORMALS, None): _StageTemplate(
        diagnosis=(
            "None of the {input_count} surface point pairs is antipodal: the"
            " angle between the surface normals must exceed"
            " `antipodal_threshold_in_radians` = {threshold:g} rad"
        ),
        nearest_clause="closest {nearest:g} rad",
        action=(
            "Lower `antipodal_threshold_in_radians` for rough or uneven meshes."
        ),
    ),
    (_AnnotationCheck.MIN_CLOSING_DISTANCE_VIOLATED, None): _StageTemplate(
        diagnosis=(
            "All {input_count} opposing surface pairs are narrower than"
            " `min_closing_distance_in_meters` = {threshold:g} m"
        ),
        nearest_clause="widest {nearest:g} m",
        action="Lower `min_closing_distance_in_meters`.",
    ),
    (_AnnotationCheck.OPENING_EXCEEDED, None): _StageTemplate(
        diagnosis=(
            "All {input_count} opposing surface pairs are wider than"
            " `opening_in_meters` = {threshold:g} m"
        ),
        nearest_clause="narrowest {nearest:g} m",
        action="Use a gripper with a wider `opening_in_meters`.",
    ),
    (_AnnotationCheck.NO_CONTACT_SURFACE_IN_FINGER_REACH, None): _StageTemplate(
        diagnosis=(
            "No surface lies within `finger_length_in_meters` = {threshold:g}"
            " m of the grasp line"
        ),
        nearest_clause="closest {nearest:g} m",
        action="Increase `finger_length_in_meters`.",
    ),
    (_AnnotationCheck.FINGER_COLLISION_WITH_PART, None): _StageTemplate(
        diagnosis=(
            "The gripper collides with the object at all {input_count}"
            " candidate poses"
        ),
        action=(
            "Check `finger_length_in_meters` (too short causes the gripper"
            " base to hit the approach surface; too long can hit surrounding"
            " features), `finger_width_in_meters`, `opening_in_meters`, or"
            " `collision_mesh`."
        ),
    ),
    (
        _AnnotationCheck.FINGER_COLLISION_WITH_PART,
        "v_notch_pinch_gripper_specs",
    ): _StageTemplate(
        diagnosis=(
            "The V-notch gripper collides with the object at all"
            " {input_count} candidate poses sampled on the fitted cylinder"
        ),
        action=(
            "Increase `notch_offset_in_meters` if the palm hits the object, or"
            " adjust `collision_mesh`."
        ),
    ),
    (_AnnotationCheck.SUCTION_CUP_OVER_HOLE, None): _StageTemplate(
        diagnosis="Every sampled patch has a hole under at least one cup",
        nearest_clause="best surface hit ratio {nearest:g} < {threshold:g}",
        action="Use a smaller cup `radius_in_meters` or fewer cups.",
    ),
    (_AnnotationCheck.SURFACE_ROUGHNESS_EXCEEDED, None): _StageTemplate(
        diagnosis=(
            "No patch is flat enough within effective"
            " `surface_roughness_threshold_in_meters` = {threshold:g} m"
        ),
        nearest_clause=(
            "smallest max ray distance {nearest:g} m, including the 1 mm"
            " ray-origin offset"
        ),
        action=(
            "Increase `surface_roughness_threshold_in_meters` or use a smaller"
            " cup `radius_in_meters`."
        ),
    ),
    (_AnnotationCheck.INCOMPATIBLE_OBJECT_SHAPE, None): _StageTemplate(
        diagnosis=(
            "The object mesh is not sufficiently cylindrical for this gripper;"
            " minimum hull-to-cylinder volume ratio is {threshold:g}"
        ),
        nearest_clause="fitted ratio {nearest:g}",
        action=(
            "Use a parallel-jaw pinch or suction gripper for non-cylindrical"
            " objects, or check the object's visual mesh."
        ),
    ),
    (_AnnotationCheck.INCOMPATIBLE_RADIAL_DIMENSION, None): _StageTemplate(
        diagnosis=(
            "The object's fitted radial dimension violates gripper radius"
            " bound {threshold:g} m"
        ),
        nearest_clause="fitted radius {nearest:g} m",
        action=(
            "Adjust `closing_radius_in_meters` or `opening_radius_in_meters`."
        ),
    ),
    (
        _AnnotationCheck.INCOMPATIBLE_RADIAL_DIMENSION,
        "v_notch_pinch_gripper_specs",
    ): _StageTemplate(
        diagnosis=(
            "The object's fitted radial dimension violates V-notch graspable"
            " radius bound {threshold:g} m"
        ),
        nearest_clause="fitted radius {nearest:g} m",
        action=(
            "Adjust `opening_in_meters`, `min_closing_distance_in_meters`,"
            " `notch_depth_in_meters`, or `notch_angle_in_radians`."
        ),
    ),
    (_AnnotationCheck.BOUNDING_BOX_CONSTRAINT_EXCLUDED, None): _StageTemplate(
        diagnosis=(
            "All {input_count} grasps were rejected by"
            " `bounding_box_constraint`"
        ),
        action=(
            "Adjust the boxes in `bounding_box_constraint`, or remove the"
            " constraint."
        ),
    ),
    (
        _AnnotationCheck.BOUNDING_BOX_CONSTRAINT_EXCLUDED,
        _INCLUDED_BOXES_KEY,
    ): _StageTemplate(
        diagnosis=(
            "All {input_count} grasps fell outside `bounding_box_constraint`"
        ),
        action=(
            "Enlarge or reposition the boxes in `bounding_box_constraint`, or"
            " remove the constraint."
        ),
    ),
    (
        _AnnotationCheck.BOUNDING_BOX_CONSTRAINT_EXCLUDED,
        _EXCLUDED_BOXES_KEY,
    ): _StageTemplate(
        diagnosis=(
            "All {input_count} grasps fell inside the excluded"
            " `bounding_box_constraint` boxes"
        ),
        action=(
            "Shrink or reposition the boxes in `bounding_box_constraint`, or"
            " remove the constraint."
        ),
    ),
    (_AnnotationCheck.AXIS_CONSTRAINT_MISALIGNED, None): _StageTemplate(
        diagnosis=(
            "No grasp approaches within `angle_tolerance_in_radians` ="
            " {threshold:g} rad of `axis_constraint`"
        ),
        nearest_clause="closest {nearest:g} rad",
        action=(
            "Widen `angle_tolerance_in_radians` or add an axis to"
            " `axis_constraint`."
        ),
    ),
    (_AnnotationCheck.INSUFFICIENT_CONTACT_AREA, None): _StageTemplate(
        diagnosis=(
            "Finger contact area is below the effective minimum of"
            " {threshold:g} m^2 per finger (`min_finger_contact_area` times the"
            " finger area)"
        ),
        nearest_clause="largest {nearest:g} m^2",
        action=(
            "Lower `min_finger_contact_area` or use fingers no larger than the"
            " contact face."
        ),
    ),
}


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


def _stage_surviving_label(
    stage_name: int,
    constraint: grasp_annotator_pb2.GraspAnnotationConstraint | None,
) -> str:
  if (
      stage_name == _AnnotationCheck.BOUNDING_BOX_CONSTRAINT_EXCLUDED
      and constraint is not None
      and constraint.HasField("bounding_box_constraint")
  ):
    if constraint.bounding_box_constraint.exclude_poses_in_box:
      return "outside the excluded bounding boxes"
    return "inside the bounding boxes"
  return _STAGE_LABELS.get(stage_name, ("candidates", "remaining"))[1]


def format_candidate_funnel(
    diagnostics: grasp_annotator_pb2.GraspAnnotationDiagnostics,
    constraint: grasp_annotator_pb2.GraspAnnotationConstraint | None = None,
) -> str:
  """Formats the candidate funnel up to and including the first empty stage.

  Args:
    diagnostics: The funnel diagnostics to format.
    constraint: Optional active constraint used to disambiguate inclusion vs.
      exclusion bounding boxes.

  Returns:
    A human-readable funnel string, or an empty string if `stage_breakdown` is
    empty.
  """
  if not diagnostics.stage_breakdown:
    return ""
  first = diagnostics.stage_breakdown[0]
  first_in_label = _STAGE_LABELS.get(
      first.stage_name, ("candidates", "remaining")
  )[0]
  parts = [f"{first.input_count} {first_in_label}"]
  prev_output = first.input_count
  for stage in diagnostics.stage_breakdown:
    in_label = _STAGE_LABELS.get(stage.stage_name, ("candidates", "remaining"))[
        0
    ]
    out_label = _stage_surviving_label(stage.stage_name, constraint)
    if stage.HasField("input_count") and stage.input_count != prev_output:
      parts.append(f"{stage.input_count} {in_label}")
    parts.append(f"{stage.output_count} {out_label}")
    prev_output = stage.output_count
    if stage.HasField("output_count") and stage.output_count == 0:
      break
  return " \u2192 ".join(parts)


def _nearest_miss(stage: grasp_annotator_pb2.StageCount) -> float | None:
  if stage.rejected_annotations and stage.rejected_annotations[0].HasField(
      "measured_value_against_threshold"
  ):
    return stage.rejected_annotations[0].measured_value_against_threshold
  return None


def _context_key(
    stage_name: int,
    gripper_specs: grasp_annotator_pb2.ParameterizedGripperSpecs | None,
    constraint: grasp_annotator_pb2.GraspAnnotationConstraint | None,
) -> str | None:
  if (
      stage_name == _AnnotationCheck.BOUNDING_BOX_CONSTRAINT_EXCLUDED
      and constraint is not None
      and constraint.HasField("bounding_box_constraint")
  ):
    return (
        _EXCLUDED_BOXES_KEY
        if constraint.bounding_box_constraint.exclude_poses_in_box
        else _INCLUDED_BOXES_KEY
    )
  if gripper_specs is not None:
    return gripper_specs.WhichOneof("gripper_specs")
  return None


def format_stage_diagnostic(
    stage: grasp_annotator_pb2.StageCount,
    gripper_specs: grasp_annotator_pb2.ParameterizedGripperSpecs | None = None,
    constraint: grasp_annotator_pb2.GraspAnnotationConstraint | None = None,
) -> StageDiagnosticText | None:
  """Formats the diagnosis and recovery action for an empty `StageCount`.

  Args:
    stage: The stage that emptied the funnel.
    gripper_specs: Optional gripper specification used to tailor gripper-specific
      diagnostics (such as V-notch radius or collision levers).
    constraint: Optional constraint used to distinguish inclusion vs. exclusion
      bounding-box checks.

  Returns:
    Structured `StageDiagnosticText`, or None if `stage.stage_name` is unknown.
  """
  reason = stage.stage_name
  if (
      reason == _AnnotationCheck.INSUFFICIENT_ANTIPODAL_NORMALS
      and stage.HasField("input_count")
      and stage.input_count == 0
  ):
    reason = _AnnotationCheck.RAY_SINGLE_HIT

  key = _context_key(reason, gripper_specs, constraint)
  template = _STAGE_TEMPLATES.get((reason, key)) or _STAGE_TEMPLATES.get(
      (reason, None)
  )
  if template is None:
    return None

  nearest = _nearest_miss(stage)
  values = {
      "input_count": stage.input_count,
      "threshold": stage.threshold,
      "nearest": nearest,
  }
  try:
    diagnosis = template.diagnosis.format(**values)
    if nearest is not None and template.nearest_clause is not None:
      diagnosis = f"{diagnosis} ({template.nearest_clause.format(**values)})"
    action = template.action.format(**values)
  except (KeyError, ValueError, TypeError) as err:
    logging.warning(
        "Failed to format grasp annotation diagnostic for stage %s: %s",
        stage.stage_name,
        err,
    )
    return None

  if stage.rejected_annotations and stage.rejected_annotations[0].debug_message:
    diagnosis = f"{diagnosis} ({stage.rejected_annotations[0].debug_message})"
  return StageDiagnosticText(diagnosis=diagnosis, action=action)
