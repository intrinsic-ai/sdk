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

"""A skill that always returns OK and has an empty equipment set."""

from intrinsic.skills.proto import footprint_pb2
from intrinsic.skills.python import skill_interface as skl
from intrinsic.skills.python import skill_interface_utils
from intrinsic.skills.testing import no_op_skill_pb2
from intrinsic.util.decorators import overrides


class NoOpSkill(skl.Skill):
  """Skill that always returns OK and has an empty equipment set."""

  @overrides(skl.Skill)
  def get_footprint(
      self,
      params: skl.GetFootprintRequest[no_op_skill_pb2.NoOpSkillParams],
      context: skl.GetFootprintContext,
  ) -> footprint_pb2.Footprint:
    """See base class."""
    return footprint_pb2.Footprint()

  @overrides(skl.Skill)
  def execute(
      self,
      request: skl.ExecuteRequest[no_op_skill_pb2.NoOpSkillParams],
      context: skl.ExecuteContext,
  ) -> None:
    """See base class."""
    return None

  @overrides(skl.Skill)
  def preview(
      self,
      request: skl.PreviewRequest[no_op_skill_pb2.NoOpSkillParams],
      context: skl.PreviewContext,
  ) -> None:
    return skill_interface_utils.preview_via_execute(self, request, context)
