// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/interfaces/force_torque_utils.h"

#include "flatbuffers/detached_buffer.h"
#include "flatbuffers/flatbuffer_builder.h"
#include "intrinsic/icon/flatbuffers/transform_types.fbs.h"
#include "intrinsic/icon/hal/interfaces/force_sensor.fbs.h"
#include "intrinsic/icon/hal/interfaces/force_torque.fbs.h"
#include "intrinsic/icon/utils/fixed_string.h"

namespace intrinsic_fbs {

flatbuffers::DetachedBuffer CreateFbsForceTorqueStatus() {
  flatbuffers::FlatBufferBuilder builder;
  builder.ForceDefaults(true);
  ForceTorqueStatusBuilder status_builder(builder);

  intrinsic_fbs::Wrench wrench;
  status_builder.add_wrench(&wrench);

  status_builder.add_enabled(false);
  status_builder.add_retare_suggested(false);
  status_builder.add_retare_completed(false);
  status_builder.add_status_code(
      intrinsic_fbs::ForceSensorStatusCode::GenericError);

  auto status = status_builder.Finish();
  builder.Finish(status);
  return builder.Release();
}

flatbuffers::DetachedBuffer CreateFbsForceTorqueCommand() {
  flatbuffers::FlatBufferBuilder builder;
  builder.ForceDefaults(true);
  ForceTorqueCommandBuilder command_builder(builder);

  command_builder.add_retare(false);
  command_builder.add_num_taring_cycles(1);

  auto command = command_builder.Finish();
  builder.Finish(command);
  return builder.Release();
}

intrinsic::icon::FixedString<kMaxFaultLength> ToFixedString(
    intrinsic_fbs::ForceSensorStatusCode status_code) {
  switch (status_code) {
    case intrinsic_fbs::ForceSensorStatusCode::Ok:
      return intrinsic::icon::FixedString<kMaxFaultLength>(
          "Force-torque sensor status is 'Ok'. Force-sensor operating "
          "normally.");
    case intrinsic_fbs::ForceSensorStatusCode::ReportedConstantReadings:
      return intrinsic::icon::FixedString<kMaxFaultLength>(
          "Reported constant readings over extended time period.");
    case intrinsic_fbs::ForceSensorStatusCode::GenericError:
      return intrinsic::icon::FixedString<kMaxFaultLength>(
          "Force-torque sensor status is 'GenericError'. Detailed error was "
          "not specified.");
    case intrinsic_fbs::ForceSensorStatusCode::
        InternalErrorRequiringPowerCycling:
      return intrinsic::icon::FixedString<kMaxFaultLength>(
          "Force-torque sensor status is 'InternalErrorRequiringPowerCycling'. "
          "Power cycle sensor to clear failure.");
    case intrinsic_fbs::ForceSensorStatusCode::MechanicalOverloadDetected:
      return intrinsic::icon::FixedString<kMaxFaultLength>(
          "Force-torque sensor status is 'MechanicalOverloadDetected'.");
    case intrinsic_fbs::ForceSensorStatusCode::SupplyOutOfRange:
      return intrinsic::icon::FixedString<kMaxFaultLength>(
          "Force-torque sensor status is 'SupplyOutOfRange'. Either voltage or "
          "current exceeded allowed limits.");
    case intrinsic_fbs::ForceSensorStatusCode::ReadingsSaturated:
      return intrinsic::icon::FixedString<kMaxFaultLength>(
          "Force-torque sensor status is 'ReadingsSaturated'.");
    case intrinsic_fbs::ForceSensorStatusCode::CalibrationChecksumError:
      return intrinsic::icon::FixedString<kMaxFaultLength>(
          "Force-torque sensor status is 'CalibrationChecksumError'.");
    case intrinsic_fbs::ForceSensorStatusCode::SimulatedError:
      return intrinsic::icon::FixedString<kMaxFaultLength>(
          "Force-torque sensor status is 'SimulatedError'.");
    default:
      return intrinsic::icon::FixedString<kMaxFaultLength>(
          "Unknown error type.");
  }
}

}  // namespace intrinsic_fbs
