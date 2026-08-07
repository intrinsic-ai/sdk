// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_WORLD_CONVERSION_SDF_SENSOR_CONVERSION_H_
#define INTRINSIC_WORLD_CONVERSION_SDF_SENSOR_CONVERSION_H_

#include <memory>
#include <string>

#include "absl/status/statusor.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/world/component/sensor_component.h"
#include "sdf/Sensor.hh"

// Conversion utilities for the SDF <sensor> element
namespace intrinsic {
namespace sdf {

struct ParseSensorResult {
  std::string name;
  Pose3d parent_t_sensor;
  std::unique_ptr<SensorComponent> sensor_component;
};
absl::StatusOr<ParseSensorResult> ParseSensor(const ::sdf::Sensor& sdf_sensor);

}  // namespace sdf
}  // namespace intrinsic

#endif  // INTRINSIC_WORLD_CONVERSION_SDF_SENSOR_CONVERSION_H_
