// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/lib/ati/ati_force_torque_bus_component.h"

#include <array>
#include <cstdint>
#include <functional>
#include <memory>
#include <utility>

#include "absl/log/log.h"
#include "absl/memory/memory.h"
#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/flatbuffers/transform_types.fbs.h"
#include "intrinsic/icon/hal/default_hardware_interfaces.h"  // IWYU pragma: keep
#include "intrinsic/icon/hal/hardware_interface_handle.h"
#include "intrinsic/icon/hal/hardware_interface_registry.h"
#include "intrinsic/icon/hal/interfaces/force_sensor.fbs.h"
#include "intrinsic/icon/hal/interfaces/force_torque.fbs.h"
#include "intrinsic/icon/hal/lib/ati/ati_constants.h"
#include "intrinsic/icon/hal/lib/ati/v1/ati_force_torque_bus_component_config.pb.h"
#include "intrinsic/icon/hal/lib/fieldbus/async_device_request.h"
#include "intrinsic/icon/hal/lib/fieldbus/device_init_context.h"
#include "intrinsic/icon/hal/lib/fieldbus/process_variable.h"
#include "intrinsic/icon/hal/lib/fieldbus/variable_registry.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/icon/utils/sensor_utils.h"
#include "intrinsic/util/status/annotate.h"
#include "intrinsic/util/status/status_macros.h"

using WrenchDataT = int32_t;
using StatusDataT = uint32_t;
using ControlDataT = uint32_t;

using ::intrinsic::fieldbus::ProcessVariable;
using ::intrinsic::fieldbus::RequestStatus;
using ::intrinsic::fieldbus::RequestType;
using ::intrinsic::fieldbus::VariableRegistry;

namespace intrinsic::ati {

namespace {

// The following translation is lossy in that the code can represent multiple
// statuses (since any combination of bits can be set). In practice, this means
// that only the first condition that matches here is returned!
//
// Note this mapping is also specific to the axia80 which you can find in here:
// https://www.ati-ia.com/app_content/Documents/9610-05-EtherCAT%20Axia80.pdf
static intrinsic_fbs::ForceSensorStatusCode StatusCodeToCanonicalError(
    uint32_t code) {
  if (code == AtiAxiaObjects::kOk) {
    return intrinsic_fbs::ForceSensorStatusCode::Ok;
  }
  if (code & AtiAxiaObjects::kBusy) {
    return intrinsic_fbs::ForceSensorStatusCode::Busy;
  }
  if (code & AtiAxiaObjects::kSupplyVoltageOutOfRange) {
    return intrinsic_fbs::ForceSensorStatusCode::SupplyOutOfRange;
  }
  if (code & AtiAxiaObjects::kBrokenGage) {
    return intrinsic_fbs::ForceSensorStatusCode::BrokenGage;
  }
  if (code & AtiAxiaObjects::kGageOutOfRange) {
    return intrinsic_fbs::ForceSensorStatusCode::MechanicalOverloadDetected;
  }
  if (code & AtiAxiaObjects::kSimulatedError) {
    return intrinsic_fbs::ForceSensorStatusCode::SimulatedError;
  }
  if (code & AtiAxiaObjects::kCalibrationChecksumError) {
    return intrinsic_fbs::ForceSensorStatusCode::CalibrationChecksumError;
  }
  if (code & AtiAxiaObjects::kSaturation) {
    return intrinsic_fbs::ForceSensorStatusCode::ReadingsSaturated;
  }
  return intrinsic_fbs::ForceSensorStatusCode::GenericError;
}

static ControlDataT Control1Word(uint32_t filter_index) {
  static ControlDataT calibration_index = 0;
  static ControlDataT internal_sampling_rate_index = 3;  // 3900hz

  ControlDataT control_word =
      (calibration_index << static_cast<ControlDataT>(
           AtiAxiaObjects::Control1Bits::kCalibrationSelectionScaler));
  control_word |= (filter_index << static_cast<ControlDataT>(
                       AtiAxiaObjects::Control1Bits::kFilterSelectionScaler));
  control_word |=
      // This implementation writes zeroes here for the oem type.
      (internal_sampling_rate_index << static_cast<ControlDataT>(
           AtiAxiaObjects::Control1Bits::kInternalSampleRateSelectionScaler));
  return control_word;
}

// Returns zero values while taring the sensor bias.
void Tare(double fx, double fy, double fz, double tx, double ty, double tz,
          std::array<::intrinsic::icon::DofSensorBias, 6>& ft_sensor_bias,
          intrinsic_fbs::Wrench* wrench) {
  auto add_bias_sample = [](::intrinsic::icon::DofSensorBias& bias,
                            double sensor_reading) -> void {
    if (bias.SampleCount() == 0) {
      bias.ResetBias();
    }
    bias.AddSample(sensor_reading);
  };
  add_bias_sample(ft_sensor_bias[0], fx);
  add_bias_sample(ft_sensor_bias[1], fy);
  add_bias_sample(ft_sensor_bias[2], fz);
  add_bias_sample(ft_sensor_bias[3], tx);
  add_bias_sample(ft_sensor_bias[4], ty);
  add_bias_sample(ft_sensor_bias[5], tz);

  wrench->mutate_x(0.0);
  wrench->mutate_y(0.0);
  wrench->mutate_z(0.0);
  wrench->mutate_rx(0.0);
  wrench->mutate_ry(0.0);
  wrench->mutate_rz(0.0);
}

// Updates the wrench with the latest sensor readings.
void ReadSensor(double fx, double fy, double fz, double tx, double ty,
                double tz,
                std::array<::intrinsic::icon::DofSensorBias, 6>& ft_sensor_bias,
                intrinsic_fbs::Wrench* wrench) {
  wrench->mutate_x(fx - ft_sensor_bias[0].Bias());
  wrench->mutate_y(fy - ft_sensor_bias[1].Bias());
  wrench->mutate_z(fz - ft_sensor_bias[2].Bias());
  wrench->mutate_rx(tx - ft_sensor_bias[3].Bias());
  wrench->mutate_ry(ty - ft_sensor_bias[4].Bias());
  wrench->mutate_rz(tz - ft_sensor_bias[5].Bias());
}

}  // namespace

absl::StatusOr<std::unique_ptr<AtiForceTorqueBusComponent>>
AtiForceTorqueBusComponent::Create(
    fieldbus::DeviceInitContext& init_context,
    const intrinsic_proto::ati::v1::ForceTorqueBusComponentConfig&
        device_config) {
  const VariableRegistry& variable_registry =
      init_context.GetVariableRegistry();
  intrinsic::icon::HardwareInterfaceRegistry& interface_registry =
      init_context.GetInterfaceRegistry();
  INTR_ASSIGN_OR_RETURN(
      auto fx, variable_registry.GetInputVariable(device_config.input_f_x()));
  INTR_RETURN_IF_ERROR(fx.IsCompatibleType<WrenchDataT>());
  INTR_ASSIGN_OR_RETURN(
      auto fy, variable_registry.GetInputVariable(device_config.input_f_y()));
  INTR_RETURN_IF_ERROR(fy.IsCompatibleType<WrenchDataT>());
  INTR_ASSIGN_OR_RETURN(
      auto fz, variable_registry.GetInputVariable(device_config.input_f_z()));
  INTR_RETURN_IF_ERROR(fz.IsCompatibleType<WrenchDataT>());
  INTR_ASSIGN_OR_RETURN(
      auto tx, variable_registry.GetInputVariable(device_config.input_t_x()));
  INTR_RETURN_IF_ERROR(tx.IsCompatibleType<WrenchDataT>());
  INTR_ASSIGN_OR_RETURN(
      auto ty, variable_registry.GetInputVariable(device_config.input_t_y()));
  INTR_RETURN_IF_ERROR(ty.IsCompatibleType<WrenchDataT>());
  INTR_ASSIGN_OR_RETURN(
      auto tz, variable_registry.GetInputVariable(device_config.input_t_z()));
  INTR_RETURN_IF_ERROR(tz.IsCompatibleType<WrenchDataT>());
  INTR_ASSIGN_OR_RETURN(auto status_code, variable_registry.GetInputVariable(
                                              device_config.status_code()));
  INTR_RETURN_IF_ERROR(status_code.IsCompatibleType<StatusDataT>());
  if (device_config.control_codes().empty()) {
    return absl::NotFoundError("ATI f/t sensor requires one control word.");
  }
  INTR_ASSIGN_OR_RETURN(auto control_1, variable_registry.GetOutputVariable(
                                            device_config.control_codes()[0]));
  INTR_RETURN_IF_ERROR(control_1.IsCompatibleType<ControlDataT>());

  absl::string_view command_interface = kCommandInterfaceName;
  if (!device_config.hardware_module_command_interface().empty()) {
    command_interface = device_config.hardware_module_command_interface();
  }
  absl::string_view status_interface = kStatusInterfaceName;
  if (!device_config.hardware_module_status_interface().empty()) {
    status_interface = device_config.hardware_module_status_interface();
  }
  INTR_ASSIGN_OR_RETURN(
      auto force_torque_status,
      interface_registry
          .AdvertiseMutableInterface<intrinsic_fbs::ForceTorqueStatus>(
              status_interface));
  INTR_ASSIGN_OR_RETURN(
      auto force_torque_command,
      interface_registry.AdvertiseInterface<intrinsic_fbs::ForceTorqueCommand>(
          command_interface));

  // Default is 0 which on the ATI is no filtering. Watch out for this for other
  // force torque sensor models.
  static constexpr ControlDataT kDefaultFilterIndex = 0;
  uint32_t filter_index = kDefaultFilterIndex;
  if (device_config.has_filter_index()) {
    filter_index = device_config.filter_index();
  }

  if (device_config.has_force_count() && device_config.has_torque_count()) {
    // Configured to read force_count and torque_count from ProcessVariables.
    INTR_ASSIGN_OR_RETURN(auto force_count, variable_registry.GetInputVariable(
                                                device_config.force_count()));
    INTR_RETURN_IF_ERROR(force_count.IsCompatibleType<int32_t>());
    INTR_ASSIGN_OR_RETURN(auto torque_count, variable_registry.GetInputVariable(
                                                 device_config.torque_count()));
    INTR_RETURN_IF_ERROR(torque_count.IsCompatibleType<int32_t>());
    return absl::WrapUnique(new AtiForceTorqueBusComponent(
        std::move(force_torque_status), std::move(force_torque_command), fx, fy,
        fz, tx, ty, tz, status_code, control_1, force_count, torque_count,
        filter_index));
  } else {
    uint32_t calibration_address = AtiEcatOemObjects::kCalibrationAddress;
    uint32_t calibration_counts_per_force =
        AtiEcatOemObjects::kCalibrationCountsPerForceSubIndex;
    uint32_t calibration_counts_per_torque =
        AtiEcatOemObjects::kCalibrationCountsPerTorqueSubIndex;
    if (!device_config.ati_oem_device()) {
      calibration_address = AtiAxiaObjects::kCalibrationAddress;
      calibration_counts_per_force =
          AtiAxiaObjects::kCalibrationCountsPerForceSubIndex;
      calibration_counts_per_torque =
          AtiAxiaObjects::kCalibrationCountsPerTorqueSubIndex;
    }
    // Configured to read force_count and torque_count via service variable.
    INTR_ASSIGN_OR_RETURN(
        auto service_variable_force_counts,
        variable_registry.GetServiceVariable(calibration_address,
                                             calibration_counts_per_force,
                                             device_config.bus_position()),
        AnnotateError(
            _.LogError(),
            "Check setting ForceTorqueBusComponentConfig.ati_oem_device. "));
    INTR_ASSIGN_OR_RETURN(
        int32_t counts_per_force, service_variable_force_counts.Read<int32_t>(),
        AnnotateError(
            _.LogError(),
            "Check setting ForceTorqueBusComponentConfig.ati_oem_device. "));
    INTR_ASSIGN_OR_RETURN(
        auto service_variable_torque_counts,
        variable_registry.GetServiceVariable(calibration_address,
                                             calibration_counts_per_torque,
                                             device_config.bus_position()));
    INTR_ASSIGN_OR_RETURN(int32_t counts_per_torque,
                          service_variable_torque_counts.Read<int32_t>());

    return absl::WrapUnique(new AtiForceTorqueBusComponent(
        std::move(force_torque_status), std::move(force_torque_command), fx, fy,
        fz, tx, ty, tz, status_code, control_1, counts_per_force,
        counts_per_torque, filter_index));
  }
}

AtiForceTorqueBusComponent::AtiForceTorqueBusComponent(
    intrinsic::icon::MutableHardwareInterfaceHandle<
        intrinsic_fbs::ForceTorqueStatus>
        force_torque_status,
    intrinsic::icon::HardwareInterfaceHandle<intrinsic_fbs::ForceTorqueCommand>
        force_torque_command,
    ProcessVariable fx, ProcessVariable fy, ProcessVariable fz,
    ProcessVariable tx, ProcessVariable ty, ProcessVariable tz,
    ProcessVariable status_code, ProcessVariable control_1,
    int32_t counts_per_force, int32_t counts_per_torque, uint32_t filter_index)
    : force_torque_status_(std::move(force_torque_status)),
      force_torque_command_(std::move(force_torque_command)),
      fx_(std::move(fx)),
      fy_(std::move(fy)),
      fz_(std::move(fz)),
      tx_(std::move(tx)),
      ty_(std::move(ty)),
      tz_(std::move(tz)),
      status_code_(std::move(status_code)),
      control_1_(std::move(control_1)),
      cyclic_force_count_variable_(nullptr, ProcessVariable::kUnknown, 0),
      cyclic_torque_count_variable_(nullptr, ProcessVariable::kUnknown, 0),
      read_cyclic_force_and_torque_count_(false),
      counts_per_force_(static_cast<double>(counts_per_force)),
      counts_per_torque_(static_cast<double>(counts_per_torque)),
      filter_index_(filter_index) {
  force_torque_status_->mutable_wrench()->mutate_x(0.0);
  force_torque_status_->mutable_wrench()->mutate_y(0.0);
  force_torque_status_->mutable_wrench()->mutate_z(0.0);
  force_torque_status_->mutable_wrench()->mutate_rx(0.0);
  force_torque_status_->mutable_wrench()->mutate_ry(0.0);
  force_torque_status_->mutable_wrench()->mutate_rz(0.0);
  force_torque_status_->mutate_status_code(
      intrinsic_fbs::ForceSensorStatusCode::Ok);
  force_torque_status_->mutate_raw_status_code(0);
  control_1_data_ = Control1Word(filter_index_);
  EndTaring();

  // Initialize the force torque sensor.
  control_1_.WriteUnchecked(control_1_data_);
}

AtiForceTorqueBusComponent::AtiForceTorqueBusComponent(
    intrinsic::icon::MutableHardwareInterfaceHandle<
        intrinsic_fbs::ForceTorqueStatus>
        force_torque_status,
    intrinsic::icon::HardwareInterfaceHandle<intrinsic_fbs::ForceTorqueCommand>
        force_torque_command,
    ProcessVariable fx, ProcessVariable fy, ProcessVariable fz,
    ProcessVariable tx, ProcessVariable ty, ProcessVariable tz,
    ProcessVariable status_code, ProcessVariable control_1,
    ProcessVariable counts_per_force, ProcessVariable counts_per_torque,
    uint32_t filter_index)
    : AtiForceTorqueBusComponent(
          std::move(force_torque_status), std::move(force_torque_command), fx,
          fy, fz, tx, ty, tz, status_code, control_1, 1, 1, filter_index) {
  cyclic_force_count_variable_ = std::move(counts_per_force),
  cyclic_torque_count_variable_ = std::move(counts_per_torque),
  read_cyclic_force_and_torque_count_ = true;
}

intrinsic::icon::RealtimeStatusOr<RequestStatus>
AtiForceTorqueBusComponent::CyclicWrite(RequestType request_type) {
  auto control_1_data = control_1_data_;
  if (request_type == RequestType::kEnableMotion ||
      request_type == RequestType::kClearFaults) {
    // Reset the sensor's bias.
    // See `Set bias against current load.` in chapter 5.2.9 of
    // https://www.ati-ia.com/app_content/documents/9620-05-EtherCAT.pdf
    control_1_data |=
        1 << static_cast<ControlDataT>(
            AtiAxiaObjects::Control1Bits::kSetBiasAgainstCurrentLoadScaler);
  }
  control_1_.WriteUnchecked(control_1_data);

  // Taring is already running. Ignore any incoming requests.
  if (IsTaring()) {
    return RequestStatus::kDone;
  }
  // Taring requested.
  if (force_torque_command_->retare()) {
    StartTaring(force_torque_command_->num_taring_cycles());
  } else {
    EndTaring();
  }

  return RequestStatus::kDone;
}

bool AtiForceTorqueBusComponent::IsTaring() const {
  return taring_in_progress_ &&
         ft_sensor_bias_[0].SampleCount() < taring_cycles_;
}

intrinsic::icon::RealtimeStatusOr<RequestStatus>
AtiForceTorqueBusComponent::CyclicRead(RequestType) {
  if (read_cyclic_force_and_torque_count_) {
    counts_per_force_ = cyclic_force_count_variable_.ReadUnchecked<int32_t>();
    counts_per_torque_ = cyclic_torque_count_variable_.ReadUnchecked<int32_t>();
  }
  double fx = fx_.ReadUnchecked<WrenchDataT>() / counts_per_force_;
  double fy = fy_.ReadUnchecked<WrenchDataT>() / counts_per_force_;
  double fz = fz_.ReadUnchecked<WrenchDataT>() / counts_per_force_;
  double tx = tx_.ReadUnchecked<WrenchDataT>() / counts_per_torque_;
  double ty = ty_.ReadUnchecked<WrenchDataT>() / counts_per_torque_;
  double tz = tz_.ReadUnchecked<WrenchDataT>() / counts_per_torque_;

  // Either update with sensor reading or tare.
  update_force_torque_status_(fx, fy, fz, tx, ty, tz, ft_sensor_bias_,
                              force_torque_status_->mutable_wrench());

  force_torque_status_->mutate_raw_status_code(
      status_code_.ReadUnchecked<StatusDataT>());
  force_torque_status_->mutate_status_code(
      StatusCodeToCanonicalError(force_torque_status_->raw_status_code()));
  force_torque_status_->mutate_retare_completed(!IsTaring());

  return RequestStatus::kDone;
}

void AtiForceTorqueBusComponent::StartTaring(int taring_cycles) {
  for (auto& bias : ft_sensor_bias_) {
    bias.ResetBias();
    bias.ResetSampleCount();
  }
  update_force_torque_status_ = std::bind(
      &Tare, std::placeholders::_1, std::placeholders::_2,
      std::placeholders::_3, std::placeholders::_4, std::placeholders::_5,
      std::placeholders::_6, std::placeholders::_7, std::placeholders::_8);

  taring_cycles_ = taring_cycles;
  taring_in_progress_ = true;
}

void AtiForceTorqueBusComponent::EndTaring() {
  update_force_torque_status_ = std::bind(
      &ReadSensor, std::placeholders::_1, std::placeholders::_2,
      std::placeholders::_3, std::placeholders::_4, std::placeholders::_5,
      std::placeholders::_6, std::placeholders::_7, std::placeholders::_8);
  taring_in_progress_ = false;
  taring_cycles_ = 0;
}

}  // namespace intrinsic::ati
