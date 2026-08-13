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

#ifndef INTRINSIC_KINEMATICS_TYPES_JOINT_STATE_H_
#define INTRINSIC_KINEMATICS_TYPES_JOINT_STATE_H_

#include "intrinsic/kinematics/types/state_rn.h"

namespace intrinsic {

// A collection of convenience typedefs for joint states as Euclidean states.
using JointStateP = StateRnP;
using JointStateV = StateRnV;
using JointStateA = StateRnA;
using JointStateJ = StateRnJ;
using JointStateT = StateRnT;
using JointStatePV = StateRnPV;
using JointStatePA = StateRnPA;
using JointStatePVA = StateRnPVA;
using JointStatePVT = StateRnPVT;
using JointStatePVAJ = StateRnPVAJ;
using JointStatePVAT = StateRnPVAT;
using JointStatePVAJT = StateRnPVAJT;
using JointStateVAJ = StateRnVAJ;

template <int N = eigenmath::MAX_EIGEN_VECTOR_SIZE>
using JointStatePVAWithMaxSize = StateRnPVAWithMaxSize<N>;

}  // namespace intrinsic

#endif  // INTRINSIC_KINEMATICS_TYPES_JOINT_STATE_H_
