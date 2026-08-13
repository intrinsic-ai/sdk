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

"""Helper module to construct a TareForceTorqueSensorAction."""

from typing import Optional

from intrinsic.icon.actions import tare_force_torque_sensor_pb2
from intrinsic.icon.python import actions

ACTION_TYPE_NAME = "intrinsic.tare_force_torque_sensor"


def create_tare_force_torque_sensor_action(
    action_id: int,
    force_torque_sensor_part_name: str,
    num_taring_cycles: Optional[int] = None,
) -> actions.Action:
  """Creates a TareForceTorqueSensorAction.

  Resets the force-torque sensor bias under the assumption that no forces or
  torques are applied to the sensor, other than those resulting from the
  attached tool. The saved bias will subsequently be subtracted from all sensor
  readings, which means that the sensor readings are 'calibrated' around the
  setpoint for which the bias was stored. Modelled post-sensor inertia (or
  configured payload) is not affected by the taring step. Run this Action before
  calling the CartesianAdmittanceAction or other compliance control Actions that
  involve a force-torque sensor.

  Only execute this action if the robot is not moving.

  Args:
    action_id: The ID of the action.
    force_torque_sensor_part_name: The name of the part providing the
      ForceTorqueSensorInterface.
    num_taring_cycles:  The number of real-time control cycles over which the
      force torque sensor offset is averaged. Defaults to 1.

  Returns:
    The TareForceTorqueSensorAction.
  """
  params = tare_force_torque_sensor_pb2.TareForceTorqueSensorParams()
  if num_taring_cycles is not None:
    params.num_taring_cycles = num_taring_cycles

  return actions.Action(
      action_id,
      ACTION_TYPE_NAME,
      force_torque_sensor_part_name,
      params,
  )
