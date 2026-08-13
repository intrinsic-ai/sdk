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

"""Helper module to construct a TrajectoryTrackingAction."""

from intrinsic.icon.actions import trajectory_tracking_action_pb2
from intrinsic.icon.proto import joint_space_pb2
from intrinsic.icon.python import actions

ACTION_TYPE_NAME: str = 'intrinsic.trajectory_tracking'


# If the initial robot position deviates more from the first trajectory position
# than this threshold, the Action will return FailedPrecondition.
MAX_INITIAL_JOINT_POSITION_DEVIATION: float = 0.1  # [rad]

# If the initial robot position deviates more from the first trajectory position
# than this threshold, the Action will return FailedPrecondition.
MAX_INITIAL_JOINT_VELOCITY_DEVIATION: float = 0.15  # [rad/sec]


def create_trajectory_tracking_action(
    action_id: int,
    part: str,
    trajectory: joint_space_pb2.JointTrajectoryPVA,
) -> actions.Action:
  """Creates a TrajectoryTrackingAction.

  Follows a given JointTrajectoryPVA and interpolates between the discretized
  states. The previously commanded robot position/velocity setpoint must be
  within `MAX_INITIAL_JOINT_POSITION_DEVIATION` and
  `MAX_INITIAL_JOINT_VELOCITY_DEVIATION` of the first joint state of the
  command. A residual controller is applied to compensate for initial state
  deviations smaller than above thresholds.

  Args:
    action_id: The ID of the action.
    part: The part has to provide a JointPosition interface.
    trajectory: Joint trajectory the robot should follow. the number of joints
      has to match the robot.

  Returns:
    The TrajectoryTrackingAction.
  """
  return actions.Action(
      action_id,
      ACTION_TYPE_NAME,
      part,
      trajectory_tracking_action_pb2.TrajectoryTrackingActionFixedParams(
          trajectory=trajectory
      ),
  )
