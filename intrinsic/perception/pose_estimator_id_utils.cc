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

#include "intrinsic/perception/pose_estimator_id_utils.h"

#include "intrinsic/perception/proto/v1/pose_estimator_id.pb.h"

namespace intrinsic::perception {
namespace {

constexpr char kDefaultPoseEstimatorIdPackage[] = "ai.intrinsic";

}  // namespace

intrinsic_proto::perception::v1::PoseEstimatorId WithDefaultPackageIfUnset(
    intrinsic_proto::perception::v1::PoseEstimatorId pose_estimator_id) {
  if (pose_estimator_id.package().empty()) {
    pose_estimator_id.set_package(kDefaultPoseEstimatorIdPackage);
  }
  return pose_estimator_id;
}

}  // namespace intrinsic::perception
