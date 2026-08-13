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

"""Helper module to construct actions.

This module only should redirect Create functions to other modules to make the
usage in the external API easier.
"""

from intrinsic.icon.actions import adio_utils

# isort: off
# isort: on

from intrinsic.icon.actions import point_to_point_move_utils
from intrinsic.icon.actions import stop_utils
from intrinsic.icon.actions import tare_force_torque_sensor_utils

# isort: off
# isort: on

from intrinsic.icon.actions import trajectory_tracking_action_utils
from intrinsic.icon.actions import wait_for_settling_utils

create_trajectory_tracking_action = (
    trajectory_tracking_action_utils.create_trajectory_tracking_action
)

create_tare_force_torque_sensor_action = (
    tare_force_torque_sensor_utils.create_tare_force_torque_sensor_action
)

create_point_to_point_move_action = (
    point_to_point_move_utils.create_point_to_point_move_action
)

create_stop_action = stop_utils.create_stop_action

create_wait_for_settling_action = (
    wait_for_settling_utils.create_wait_for_settling_action
)

create_digital_output_action = adio_utils.create_digital_output_action
