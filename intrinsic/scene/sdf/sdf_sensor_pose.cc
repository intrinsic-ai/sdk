// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/scene/sdf/sdf_sensor_pose.h"

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/math/pose3.h"
#include "sdf/Sensor.hh"

namespace intrinsic {
namespace sdf {

namespace {

bool IsOrientedSensorType(::sdf::SensorType sensor_type) {
  switch (sensor_type) {
    case ::sdf::SensorType::CAMERA:
    case ::sdf::SensorType::DEPTH_CAMERA:
    case ::sdf::SensorType::GPU_LIDAR:
    case ::sdf::SensorType::LIDAR:
    case ::sdf::SensorType::RGBD_CAMERA:
    case ::sdf::SensorType::THERMAL_CAMERA:
      return true;
    default:
      return false;
  }
}

Pose3d SdfSensorTIntrinsicSenzor(::sdf::SensorType sensor_type) {
  if (IsOrientedSensorType(sensor_type)) {
    // Convert between the following to conventions:
    // SDF: z up, x forward, y left
    // intrinsic::World: z forward, x right, y down
    // The following is the inverse of "<pose>0 0 0 1.570796326794896619
    // -1.570796326794896619 0</pose>".
    return Pose3d(eigenmath::Quaterniond(0.5, -0.5, 0.5, -0.5));
  } else {
    return Pose3d::Identity();
  }
}

}  // namespace

Pose3d SensorPoseFromSdf(const Pose3d& parent_t_sdf_sensor,
                         ::sdf::SensorType sensor_type) {
  return parent_t_sdf_sensor * SdfSensorTIntrinsicSenzor(sensor_type);
}

Pose3d SensorPoseToSdf(const Pose3d& parent_t_intrinsic_sensor,
                       ::sdf::SensorType sensor_type) {
  return parent_t_intrinsic_sensor *
         SdfSensorTIntrinsicSenzor(sensor_type).inverse();
}

}  // namespace sdf
}  // namespace intrinsic
