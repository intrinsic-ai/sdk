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

#ifndef INTRINSIC_MOTION_PLANNING_CONVERSIONS_H_
#define INTRINSIC_MOTION_PLANNING_CONVERSIONS_H_

#include <vector>

#include "google/protobuf/repeated_field.h"
#include "google/protobuf/repeated_ptr_field.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/proto/joint_space.pb.h"

namespace intrinsic {

std::vector<eigenmath::VectorXd> ToVectorXds(
    const google::protobuf::RepeatedPtrField<intrinsic_proto::icon::JointVec>&
        vectors);
// Clears any data currently in vectors.
void ToJointVecs(
    const std::vector<eigenmath::VectorXd>& path,
    google::protobuf::RepeatedPtrField<intrinsic_proto::icon::JointVec>*
        vectors);

}  // namespace intrinsic

#endif  // INTRINSIC_MOTION_PLANNING_CONVERSIONS_H_
