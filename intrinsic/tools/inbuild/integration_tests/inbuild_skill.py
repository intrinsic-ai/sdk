# Copyright 2023 Intrinsic Innovation LLC

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
