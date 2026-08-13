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

"""A skill that echos its parameters as its result."""

from intrinsic.skills.python import skill_interface as skl
from intrinsic.skills.python import skill_interface_utils
from intrinsic.skills.testing import echo_skill_pb2
from intrinsic.util import decorators


class EchoSkill(skl.Skill):
  """A skill that echos its parameters as its result."""

  @decorators.overrides(skl.Skill)
  def execute(
      self,
      request: skl.ExecuteRequest[echo_skill_pb2.EchoSkillParams],
      context: skl.ExecuteContext,
  ) -> echo_skill_pb2.EchoSkillReturn:
    return echo_skill_pb2.EchoSkillReturn(foo=request.params.foo)

  @decorators.overrides(skl.Skill)
  def preview(
      self,
      request: skl.PreviewRequest[echo_skill_pb2.EchoSkillParams],
      context: skl.PreviewContext,
  ) -> echo_skill_pb2.EchoSkillReturn:
    return skill_interface_utils.preview_via_execute(self, request, context)
