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

#ifndef INTRINSIC_PERCEPTION_PUBLIC_POSE_ESTIMATOR_ID_UTILS_H_
#define INTRINSIC_PERCEPTION_PUBLIC_POSE_ESTIMATOR_ID_UTILS_H_

#include "intrinsic/perception/proto/v1/pose_estimator_id.pb.h"

namespace intrinsic::perception {

// Returns a copy of the given `pose_estimator_id` with `package` set to
// "ai.intrinsic" if it is empty/unset.
intrinsic_proto::perception::v1::PoseEstimatorId WithDefaultPackageIfUnset(
    intrinsic_proto::perception::v1::PoseEstimatorId pose_estimator_id);

}  // namespace intrinsic::perception

#endif  // INTRINSIC_PERCEPTION_PUBLIC_POSE_ESTIMATOR_ID_UTILS_H_
