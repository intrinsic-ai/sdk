// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_LIB_ATI_ATI_FORCE_TORQUE_BUS_COMPONENT_H_
#define INTRINSIC_ICON_HAL_LIB_ATI_ATI_FORCE_TORQUE_BUS_COMPONENT_H_

#include <array>
#include <cstdint>
#include <functional>
#include <memory>

#include "absl/status/statusor.h"
#include "intrinsic/icon/flatbuffers/transform_types.fbs.h"
#include "intrinsic/icon/hal/default_hardware_interfaces.h"  // IWYU pragma: keep
#include "intrinsic/icon/hal/hardware_interface_handle.h"
#include "intrinsic/icon/hal/interfaces/force_torque.fbs.h"
#include "intrinsic/icon/hal/lib/ati/v1/ati_force_torque_bus_component_config.pb.h"
#include "intrinsic/icon/hal/lib/fieldbus/async_device_request.h"
#include "intrinsic/icon/hal/lib/fieldbus/bus_component.h"
#include "intrinsic/icon/hal/lib/fieldbus/bus_component_factory.h"
#include "intrinsic/icon/hal/lib/fieldbus/device_init_context.h"
#include "intrinsic/icon/hal/lib/fieldbus/process_variable.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/icon/utils/sensor_utils.h"

namespace intrinsic::ati {

// Fieldbus component implementation for Ati Force Torque sensors.
class AtiForceTorqueBusComponent final : public fieldbus::BusComponent {
 public:
  static constexpr char kStatusInterfaceName[] = "force_torque_status";
  static constexpr char kCommandInterfaceName[] = "force_torque_command";

  AtiForceTorqueBusComponent(const AtiForceTorqueBusComponent& other) noexcept =
      delete;
  AtiForceTorqueBusComponent& operator=(
      const AtiForceTorqueBusComponent& other) = delete;
  AtiForceTorqueBusComponent(AtiForceTorqueBusComponent&& other) noexcept =
      default;
  AtiForceTorqueBusComponent& operator=(AtiForceTorqueBusComponent&& other) =
      default;

  // Creates and initializes a new Ati bus device.
  // The devices registers the force torque wrench hardware interface and sets
  // up the connection to the bus to query the necessary process image fields.
  static absl::StatusOr<std::unique_ptr<AtiForceTorqueBusComponent>> Create(
      fieldbus::DeviceInitContext& init_context,
      const intrinsic_proto::ati::v1::ForceTorqueBusComponentConfig&
          device_config);

  // Reads the current force torque values and updates the wrench hardware
  // interface.
  intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus> CyclicRead(
      fieldbus::RequestType request_type) override;

  // Writes control words to the force torque sensor.
  intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus> CyclicWrite(
      fieldbus::RequestType request_type) override;

  // Indicates whether taring is currently happening.
  bool IsTaring() const;

 private:
  // Constructor.
  // Reads `counts_per_force` and `counts_per_torque` via SDO.
  AtiForceTorqueBusComponent(
      intrinsic::icon::MutableHardwareInterfaceHandle<
          intrinsic_fbs::ForceTorqueStatus>
          force_torque_status,
      intrinsic::icon::HardwareInterfaceHandle<
          intrinsic_fbs::ForceTorqueCommand>
          force_torque_command,
      fieldbus::ProcessVariable fx, fieldbus::ProcessVariable fy,
      fieldbus::ProcessVariable fz, fieldbus::ProcessVariable tx,
      fieldbus::ProcessVariable ty, fieldbus::ProcessVariable tz,
      fieldbus::ProcessVariable status_code,
      fieldbus::ProcessVariable control_1, int32_t counts_per_force,
      int32_t counts_per_torque, uint32_t filter_index);
  // Constructor.
  // Reads `counts_per_force` and `counts_per_torque` via ProcessVariables.
  AtiForceTorqueBusComponent(
      intrinsic::icon::MutableHardwareInterfaceHandle<
          intrinsic_fbs::ForceTorqueStatus>
          force_torque_status,
      intrinsic::icon::HardwareInterfaceHandle<
          intrinsic_fbs::ForceTorqueCommand>
          force_torque_command,
      fieldbus::ProcessVariable fx, fieldbus::ProcessVariable fy,
      fieldbus::ProcessVariable fz, fieldbus::ProcessVariable tx,
      fieldbus::ProcessVariable ty, fieldbus::ProcessVariable tz,
      fieldbus::ProcessVariable status_code,
      fieldbus::ProcessVariable control_1,
      fieldbus::ProcessVariable counts_per_force,
      fieldbus::ProcessVariable counts_per_torque, uint32_t filter_index);

  // Starts taring and sets the force torque readings to 0 until `EndTaring` is
  // called.
  void StartTaring(int taring_cycles);

  // Ends the taring process and updates the force torque readings to the actual
  // sensor readings.
  void EndTaring();

  intrinsic::icon::MutableHardwareInterfaceHandle<
      intrinsic_fbs::ForceTorqueStatus>
      force_torque_status_;
  intrinsic::icon::HardwareInterfaceHandle<intrinsic_fbs::ForceTorqueCommand>
      force_torque_command_;
  fieldbus::ProcessVariable fx_;
  fieldbus::ProcessVariable fy_;
  fieldbus::ProcessVariable fz_;
  fieldbus::ProcessVariable tx_;
  fieldbus::ProcessVariable ty_;
  fieldbus::ProcessVariable tz_;
  fieldbus::ProcessVariable status_code_;
  fieldbus::ProcessVariable
      control_1_;  // maps to control 1, setting default control 2
  fieldbus::ProcessVariable cyclic_force_count_variable_;
  fieldbus::ProcessVariable cyclic_torque_count_variable_;
  bool read_cyclic_force_and_torque_count_;
  double counts_per_force_;
  double counts_per_torque_;
  uint32_t control_1_data_;

  // Index  the filter to be used. Note the value is directly incorporated into
  // control_1_data_.
  uint32_t filter_index_;

  std::function<void(double, double, double, double, double, double,
                     std::array<::intrinsic::icon::DofSensorBias, 6>&,
                     intrinsic_fbs::Wrench*)>
      update_force_torque_status_;

  // Sensor bias which gets set during taring operation and is subsequently
  // subtracted from wrench measurements to correct for unmodelled payload, etc.
  std::array<::intrinsic::icon::DofSensorBias, 6> ft_sensor_bias_;
  bool taring_in_progress_ = false;
  int taring_cycles_ = 0;
};

}  // namespace intrinsic::ati

// Registers the AtiForceTorqueBusComponent and its config type with the
// bus component factory. Allows constructing the device from its config
// type via `CreateBusComponentFromConfig`.
REGISTER_BUS_COMPONENT(intrinsic::ati::AtiForceTorqueBusComponent,
                       intrinsic_proto::ati::v1::ForceTorqueBusComponentConfig);

#endif  // INTRINSIC_ICON_HAL_LIB_ATI_ATI_FORCE_TORQUE_BUS_COMPONENT_H_
