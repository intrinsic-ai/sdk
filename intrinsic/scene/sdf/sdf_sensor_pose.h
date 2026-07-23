// Copyright 2023 Intrinsic Innovation LLC

// This file provides a utilities for converting between World and SDF pose
// conventions.

#ifndef INTRINSIC_SCENE_SDF_SDF_SENSOR_POSE_H_
#define INTRINSIC_SCENE_SDF_SDF_SENSOR_POSE_H_

#include "intrinsic/math/pose3.h"
#include "sdf/Sensor.hh"

namespace intrinsic {
namespace sdf {

// Rotates the given frame of an SDF sensor to comply with the convention used
// by intrinsic::World. Only has an effect for sensor types that follow an
// orientation convention. E.g., affects sensors with type "camera" and
// "lidar"/"ray" frames, but does not affect sensors with type "force_torque".
//
// 'parent_t_sensor_sdf' is the pose of the sensor in the frame of its parent as
// parsed from an SDF. For valid values of 'sensor_type', see
// http://sdformat.org/spec?ver=1.7&elem=sensor#sensor_type.
// Returns a corrected pose "parent_t_intrinsic_sensor".
Pose3d SensorPoseFromSdf(const Pose3d& parent_t_sdf_sensor,
                         ::sdf::SensorType sensor_type);

// Inverse of SensorPoseFromSdf(). Transforms the given
// 'parent_t_intrinsic_sensor' to "parent_t_sdf_sensor" depending on the given
// 'sensor_type'.
Pose3d SensorPoseToSdf(const Pose3d& parent_t_intrinsic_sensor,
                       ::sdf::SensorType sensor_type);

}  // namespace sdf
}  // namespace intrinsic

#endif  // INTRINSIC_SCENE_SDF_SDF_SENSOR_POSE_H_
