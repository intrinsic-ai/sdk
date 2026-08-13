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

"""Helper module to construct a Stop action."""

import dataclasses

from intrinsic.icon.python import actions

ACTION_TYPE_NAME = "intrinsic.stop"


@dataclasses.dataclass(frozen=True)
class StateVariables:
  IS_SETTLED = "is_settled"


def create_stop_action(
    action_id: int, joint_position_part_name: str
) -> actions.Action:
  """Creates a Stop action.

  Args:
    action_id: The ID of the action.
    joint_position_part_name:  The name of the part providing the JointPosition
      interface.

  Returns:
    The Stop action.
  """
  return actions.Action(
      action_id, ACTION_TYPE_NAME, joint_position_part_name, None
  )
