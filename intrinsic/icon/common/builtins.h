// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#ifndef INTRINSIC_ICON_COMMON_BUILTINS_H_
#define INTRINSIC_ICON_COMMON_BUILTINS_H_

// This header defines the names of built-in action types and state variables.

namespace intrinsic {
namespace icon {

// Name and description of the builtin Stop action.
inline constexpr char kStopAction[] = "intrinsic.stop";
inline constexpr char kStopActionDescription[] =
    "Stops all motion of the part while respecting joint limits. Computes a "
    "time-optimal stopping trajectory with zero terminal velocity subject to "
    "maximum hardware acceleration and jerk limits. Use this for all "
    "applications that do not involve force control. The stop motion is "
    "synchronized in joint space. Therefore, the resulting stop motion does "
    "not continue tracking previously commanded motions in a path-accurate "
    "way. The JointStopAction also holds a settling state estimator which "
    "monitors residual oscillations after the stop trajectory has been played "
    "back, see state variable documentation of the JointStopAction.";
inline constexpr char kStopPartSlot[] = "position_part";
inline constexpr char kStopPartSlotDescription[] =
    "The Action smoothly decelerates the joints of this Part to reach a stop.";

// Name of the builtin state variable that tells if an action has completed.
// All actions expose this, but the semantics are action-specific and actions
// are allowed to never be done. Value is a boolean.
inline constexpr char kIsDone[] = "intrinsic.is_done";
inline constexpr char kIsDoneDescription[] =
    "Builtin state variable that tells whether an action has completed";
// Name of the builtin state variable that tells if an action has brought the
// robot to a stop. Stop actions must expose this. Other actions may optionally
// expose this. Value is a boolean.
inline constexpr char kIsStopped[] = "intrinsic.is_stopped";
inline constexpr char kIsStoppedDescription[] =
    "Builtin state variable that tells if an action has brought the robot to a "
    "stop";
// Name of the builtin state variable that reports the amount of time since an
// action became active, in seconds. Actions optionally expose this.
// Value is a double.
inline constexpr char kActionElapsedTime[] = "intrinsic.action_elapsed_time";
inline constexpr char kActionElapsedTimeDescription[] =
    "Builtin state variable that reports the amount of time since an action "
    "became active, in seconds";

}  // namespace icon
}  // namespace intrinsic

#endif  // INTRINSIC_ICON_COMMON_BUILTINS_H_
