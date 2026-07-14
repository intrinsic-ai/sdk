// Copyright 2023 Intrinsic Innovation LLC

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
