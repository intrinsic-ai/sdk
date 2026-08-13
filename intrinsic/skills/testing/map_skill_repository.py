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

"""Implementation of SkillRepository to store skills in a map."""

import threading

from intrinsic.skills.internal import runtime_data as rd
from intrinsic.skills.internal import single_skill_factory as skill_factory
from intrinsic.skills.internal import skill_repository as repo
from intrinsic.skills.python import skill_interface as skl
from intrinsic.util.decorators import overrides


class MapSkillRepository(repo.SkillRepository):
  """Implementation of SkillRepository to store skills in a map."""

  def __init__(
      self,
      skill_factories: dict[str, skill_factory.SingleSkillFactory],
  ):
    """Initializes the MapSkillRepository.

    Args:
      skill_factories: Maps a skill alias to a single skill factory.
    """
    self._lock = threading.Lock()
    self._skill_factories: dict[str, skill_factory.SingleSkillFactory] = (
        skill_factories
    )

  @overrides(repo.SkillRepository)
  def get_skill(self, skill_alias: str) -> skl.Skill:
    """Returns a new skill instance.

    Args:
      skill_alias: The skill alias.

    Raises:
      InvalidSkillAliasError: If the skill does not exist in the repository.
    """
    try:
      factory = self._skill_factories[skill_alias]
    except KeyError as err:
      raise repo.InvalidSkillAliasError(
          f'Skill with alias [{skill_alias}]" not found in the repository'
      ) from err

    return factory.get_skill(skill_alias)

  @overrides(repo.SkillRepository)
  def get_skill_execute(self, skill_alias: str) -> skl.SkillExecuteInterface:
    """Returns a new skill execute instance.

    Args:
      skill_alias: The skill's alias.

    Raises:
      SkillNotFoundError if the skill is not found in the repository.
    """
    try:
      factory = self._skill_factories[skill_alias]
    except KeyError as err:
      raise repo.InvalidSkillAliasError(
          f'Skill with alias [{skill_alias}]" not found in the repository'
      ) from err

    return factory.get_skill_execute(skill_alias)

  @overrides(repo.SkillRepository)
  def get_skill_project(self, skill_alias: str) -> skl.SkillProjectInterface:
    """Returns a new skill project instance.

    Args:
      skill_alias: The skill's alias.

    Raises:
      SkillNotFoundError if the skill is not found in the repository.
    """
    try:
      factory = self._skill_factories[skill_alias]
    except KeyError as err:
      raise repo.InvalidSkillAliasError(
          f'Skill with alias [{skill_alias}]" not found in the repository'
      ) from err

    return factory.get_skill_project(skill_alias)

  @overrides(repo.SkillRepository)
  def get_skill_runtime_data(self, skill_alias: str) -> rd.SkillRuntimeData:
    """Returns runtime data for the skill.

    Args:
      skill_alias: The skill's alias.

    Raises:
      SkillNotFoundError if the skill is not found in the repository.
    """
    try:
      factory = self._skill_factories[skill_alias]
    except KeyError as err:
      raise repo.InvalidSkillAliasError(
          f'Skill with alias [{skill_alias}]" not found in the repository'
      ) from err

    return factory.get_skill_runtime_data(skill_alias)

  @overrides(repo.SkillRepository)
  def get_skill_aliases(self) -> list[str]:
    """Returns the list of aliases of all registered skills."""
    with self._lock:
      return self._skill_factories.keys()
