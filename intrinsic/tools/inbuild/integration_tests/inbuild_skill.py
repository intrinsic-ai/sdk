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

"""A Python skill that is built with inbuild."""

import logging

from intrinsic.skills.python import skill_interface
from intrinsic.tools.inbuild.integration_tests import inbuild_skill_pb2
from intrinsic.util.decorators import overrides

logger = logging.getLogger(__name__)


class InbuildSkill(skill_interface.Skill):

  @overrides(skill_interface.Skill)
  def execute(
      self,
      request: skill_interface.ExecuteRequest[
          inbuild_skill_pb2.InbuildSkillParams
      ],
      context: skill_interface.ExecuteContext,
  ) -> None:
    logger.info("Hello from InbuildSkill.Execute: %s", request.params.foo)
    return None
