// Copyright 2023 Intrinsic Innovation LLC

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
