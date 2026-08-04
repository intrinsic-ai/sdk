// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_WORLD_COMPONENT_SENSOR_COMPONENT_H_
#define INTRINSIC_WORLD_COMPONENT_SENSOR_COMPONENT_H_

#include <memory>
#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/world/proto/sensor_component.pb.h"

namespace intrinsic {

// An enum class holding the type of the sensor.
enum class SensorType {
  kInvalid = 0,
  kCamera,
  kDepthCamera,
  kForceTorque,
  kLidar,
};

// A component to hold information about sensors (measuring devices) in the
// World.
//
// Sensor types supported:
//   - Camera*
//   - Depth camera*
//   - Force torque sensor
//
// *Sensor is oriented according to the following convention: z forward, x
// right, y down (in the frame of the owning entity).
class SensorComponent {
 public:
  virtual ~SensorComponent() = default;

  // Returns a new SensorComponent instance.
  static std::unique_ptr<SensorComponent> Create();

  // Returns a new SensorComponent instance derived from the given proto.
  static absl::StatusOr<std::unique_ptr<SensorComponent>> FromProto(
      const intrinsic_proto::world::SensorComponent& proto);

  // Returns a new SensorComponent instance that is a copy of this one.
  virtual std::unique_ptr<SensorComponent> Clone() const = 0;

  // Returns a proto representation of this component.
  virtual absl::StatusOr<intrinsic_proto::world::SensorComponent> ToProto()
      const = 0;

  // Get whether the sensor is always updating.
  virtual bool GetAlwaysOn() const = 0;

  // Set whether the sensor is always updating.
  virtual void SetAlwaysOn(bool always_on) = 0;

  // Get the sensor update rate (in Hz).
  virtual double GetUpdateRate() const = 0;

  // Set the sensor update rate (in Hz).
  virtual void SetUpdateRate(double update_rate) = 0;

  // Get topic name.
  virtual const std::string& GetTopic() const = 0;

  // Set topic name.
  virtual void SetTopic(absl::string_view topic) = 0;

  // Get the type of the sensor.
  //
  // Note that the type might be SensorType::kInvalid.
  // To set the type, use one of the SetXxxSpec() member functions.
  virtual SensorType GetType() const = 0;

  // Get the camera spec.
  //
  // Returns NOT_FOUND error if the sensor type is not camera.
  virtual absl::StatusOr<intrinsic_proto::world::SensorComponent::Camera>
  GetCameraSpec() const = 0;

  // Get the depth camera spec.
  //
  // Returns NOT_FOUND error if the sensor type is not depth camera.
  virtual absl::StatusOr<intrinsic_proto::world::SensorComponent::DepthCamera>
  GetDepthCameraSpec() const = 0;

  // Get the force torque spec.
  //
  // Returns NOT_FOUND error if the sensor type is not force torque.
  virtual absl::StatusOr<intrinsic_proto::world::SensorComponent::ForceTorque>
  GetForceTorqueSpec() const = 0;

  // Get the lidar spec.
  //
  // Returns NOT_FOUND error if the sensor type is not lidar.
  virtual absl::StatusOr<intrinsic_proto::world::SensorComponent::Lidar>
  GetLidarSpec() const = 0;

  // Set the sensor type to camera and set the spec to `spec`.
  //
  // The previous spec is overwritten.
  virtual void SetCameraSpec(
      const intrinsic_proto::world::SensorComponent::Camera& spec) = 0;

  // Set the sensor type to depth camera and set the spec to `spec`.
  //
  // The previous spec is overwritten.
  virtual void SetDepthCameraSpec(
      const intrinsic_proto::world::SensorComponent::DepthCamera& spec) = 0;

  // Set the sensor type to force torque and set the spec to `spec`.
  //
  // The previous spec is overwritten.
  virtual void SetForceTorqueSpec(
      const intrinsic_proto::world::SensorComponent::ForceTorque& spec) = 0;

  // Set the sensor type to lidar and set the spec to `spec`.
  //
  // The previous spec is overwritten.
  virtual void SetLidarSpec(
      const intrinsic_proto::world::SensorComponent::Lidar& spec) = 0;

  // Updates the component based on the given proto. This will do a full
  // override, if something is missing from this proto it will override any
  // existing values with the defaults.
  virtual absl::Status UpdateFromProto(
      const intrinsic_proto::world::SensorComponent& proto) = 0;
};

}  // namespace intrinsic

#endif  // INTRINSIC_WORLD_COMPONENT_SENSOR_COMPONENT_H_
