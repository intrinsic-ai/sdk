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

#include "intrinsic/motion_planning/conversions.h"

#include <vector>

#include "google/protobuf/repeated_ptr_field.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/proto/joint_space.pb.h"
#include "intrinsic/util/eigen.h"

namespace intrinsic {

std::vector<eigenmath::VectorXd> ToVectorXds(
    const google::protobuf::RepeatedPtrField<intrinsic_proto::icon::JointVec>&
        vectors) {
  std::vector<eigenmath::VectorXd> out;
  out.reserve(vectors.size());
  for (const auto& v : vectors) {
    out.push_back(RepeatedDoubleToVectorXd(v.joints()));
  }
  return out;
}

void ToJointVecs(
    const std::vector<eigenmath::VectorXd>& path,
    google::protobuf::RepeatedPtrField<intrinsic_proto::icon::JointVec>*
        vectors) {
  vectors->Clear();
  vectors->Reserve(path.size());
  for (const auto& point : path) {
    auto* new_point_proto = vectors->Add();
    VectorXdToRepeatedDouble(point, new_point_proto->mutable_joints());
  }
}

}  // namespace intrinsic
