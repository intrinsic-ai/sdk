// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/lib/service_variable_value/service_variable_device.h"

#include <cstdint>
#include <limits>
#include <list>
#include <memory>
#include <string>
#include <utility>
#include <vector>

#include "absl/functional/any_invocable.h"
#include "absl/log/log.h"
#include "absl/memory/memory.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/str_join.h"
#include "absl/strings/string_view.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "google/protobuf/message.h"
#include "intrinsic/icon/hal/default_hardware_interfaces.h"  // IWYU pragma: keep
#include "intrinsic/icon/hal/hardware_interface_handle.h"
#include "intrinsic/icon/hal/hardware_interface_registry.h"
#include "intrinsic/icon/hal/interfaces/io_controller.fbs.h"
#include "intrinsic/icon/hal/lib/fieldbus/async_device_request.h"
#include "intrinsic/icon/hal/lib/fieldbus/device_init_context.h"
#include "intrinsic/icon/hal/lib/fieldbus/service_variable.h"
#include "intrinsic/icon/hal/lib/fieldbus/v1/value_parsing.pb.h"
#include "intrinsic/icon/hal/lib/fieldbus/variable_registry.h"
#include "intrinsic/icon/hal/lib/service_variable_value/v1/service_variable_config.pb.h"
#include "intrinsic/icon/hal/lib/service_variable_value/v1/service_variable_export.pb.h"
#include "intrinsic/icon/utils/async_buffer.h"
#include "intrinsic/icon/utils/log.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_macro.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/logging/proto/log_item.pb.h"
#include "intrinsic/platform/pubsub/publisher.h"
#include "intrinsic/platform/pubsub/pubsub.h"
#include "intrinsic/util/proto_time.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/util/thread/stop_token.h"
#include "intrinsic/util/thread/thread.h"
#include "intrinsic/util/thread/thread_utils.h"

namespace intrinsic::fieldbus {

static constexpr absl::Duration kDefaultSdoReadInterval = absl::Seconds(5);
constexpr absl::string_view kEcatTopicPrefix = "/ethercat/";
constexpr absl::string_view kSdoExportTopicSuffix = "/sdo_export";

namespace {

using ::intrinsic_proto::fieldbus::v1::ServiceVariableExportType;

// This class is used for thread-safe communication between the SDO reader
// thread and CyclicRead.
class SdoToInterfaceWriter {
 public:
  // Constructs an SdoToInterfaceWriter with the given hardware interface
  // handle.
  explicit SdoToInterfaceWriter(
      intrinsic::icon::MutableHardwareInterfaceHandle<intrinsic_fbs::AIOStatus>
          i_handle)
      : async_buffer_(0.0), interface_handle_(std::move(i_handle)) {}

  // Writes a `value` to the async buffer. This is thread-safe.
  void Write(double value) {
    double* free_buffer = async_buffer_.GetFreeBuffer();
    *free_buffer = value;
    async_buffer_.CommitFreeBuffer();
  }

  // Updates the `interface_handle_` with the latest value from the
  // `async_buffer_`. Returns OkStatus on success.
  intrinsic::icon::RealtimeStatus UpdateInterface() {
    double* active_buffer = nullptr;
    bool new_data = async_buffer_.GetActiveBuffer(&active_buffer);
    if (new_data && active_buffer != nullptr) {
      interface_handle_->mutable_signals()->GetMutableObject(0)->mutate_value(
          *active_buffer);
    }
    return intrinsic::icon::OkStatus();
  }

 private:
  // Thread-safe buffer to pass double values between threads.
  intrinsic::AsyncBuffer<double> async_buffer_;
  // Handle to the mutable hardware interface for AIOStatus.
  intrinsic::icon::MutableHardwareInterfaceHandle<intrinsic_fbs::AIOStatus>
      interface_handle_;
};

template <typename SourceType, typename TargetType>
absl::Status RangeCheck(SourceType value, TargetType min_value,
                        TargetType max_value) {
  if (value < min_value) {
    return absl::OutOfRangeError(absl::StrCat("Value is ", value,
                                              ", but must not be less than ",
                                              static_cast<double>(min_value)));
  }
  if (value > max_value) {
    return absl::OutOfRangeError(absl::StrCat("Value is ", value,
                                              ", but must not be greater than ",
                                              static_cast<double>(max_value)));
  }
  return absl::OkStatus();
}

// Returns the string representation of the ServiceVariableValueType.
std::string ToString(
    const intrinsic_proto::fieldbus::v1::ServiceVariableValueType& value) {
  switch (value.value_case()) {
    case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::kBoolValue: {
      return absl::StrCat(value.bool_value() ? "true" : "false");
    }
    case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::kUint8Value: {
      return absl::StrCat(value.uint8_value());
    }
    case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::kInt8Value: {
      return absl::StrCat(value.int8_value());
    }
    case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::
        kUint16Value: {
      return absl::StrCat(value.uint16_value());
    }
    case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::kInt16Value: {
      return absl::StrCat(value.int16_value());
    }
    case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::
        kUint32Value: {
      return absl::StrCat(value.uint32_value());
    }
    case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::kInt32Value: {
      return absl::StrCat(value.int32_value());
    }
    case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::
        kUint64Value: {
      return absl::StrCat(value.uint64_value());
    }
    case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::kInt64Value: {
      return absl::StrCat(value.int64_value());
    }
    case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::kFloatValue: {
      return absl::StrCat(value.float_value());
    }
    case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::
        kDoubleValue: {
      return absl::StrCat(value.double_value());
    }
    case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::
        VALUE_NOT_SET: {
      return "VALUE_NOT_SET";
    }
    default: {
      // google::protobuf::ShortFormat unfortunately prints go/debugstr.
      return absl::StrCat("unknown value type: ",
                          google::protobuf::ShortFormat(value));
    }
  }
}

absl::StatusOr<intrinsic::Thread> SetupSdoPublisher(
    absl::string_view module_name,
    const intrinsic::fieldbus::VariableRegistry& variable_registry,
    const intrinsic_proto::fieldbus::v1::ServiceVariableConfig& config,
    absl::Duration cycle_time,
    std::vector<std::unique_ptr<SdoToInterfaceWriter>>
        service_variable_to_interface_writers) {
  const ServiceVariableExportType export_type = config.sdo_export_type();
  if (export_type == ServiceVariableExportType::SDO_EXPORT_TYPE_NONE) {
    LOG(INFO) << "No service variable export configured.";
    return intrinsic::Thread();
  }

  if (config.sdo_reads().empty()) {
    LOG(INFO) << "No service variable reads configured.";
    return intrinsic::Thread();
  }

  absl::Duration export_interval = kDefaultSdoReadInterval;
  if (config.has_sdo_export_interval()) {
    INTR_ASSIGN_OR_RETURN(
        export_interval,
        intrinsic::ToAbslDurationRejectNegative(config.sdo_export_interval()));

    if (export_interval == absl::ZeroDuration()) {
      return absl::InvalidArgumentError(
          "service variable export interval can not be zero.");
    }
  }

  // To avoid flooding the bus with too many reads, we assume
  // ServiceVariableDevice::kMinCyclesPerSdoRead cycles per service variable
  // read. Thus, the export interval cannot be smaller than:
  // ServiceVariableDevice::kMinCyclesPerSdoRead * cycle_time *
  // #service_variable_reads. This is an estimate, the bus might be able to
  // handle more reads, but we don't want to investigate the limits of the bus
  // here, as we believe this feature is rarely pushed to the this limit.
  auto min_export_interval = ServiceVariableDevice::kMinCyclesPerSdoRead *
                             cycle_time * config.sdo_reads_size();
  if (export_interval < min_export_interval) {
    return absl::InvalidArgumentError(absl::StrCat(
        "sdo_export_interval is too short: must be at least ",
        absl::ToDoubleMilliseconds(min_export_interval), " ms, which is ",
        ServiceVariableDevice::kMinCyclesPerSdoRead,
        " * cycle_time * #service_variable_reads."));
  }

  std::vector<absl::AnyInvocable<
      absl::StatusOr<intrinsic_proto::fieldbus::v1::SdoStatus>()>>
      service_variable_readers;
  service_variable_readers.reserve(config.sdo_reads_size());

  for (int i = 0; i < config.sdo_reads_size(); ++i) {
    const intrinsic_proto::fieldbus::v1::ServiceVariableRead&
        service_variable_read = config.sdo_reads(i);
    std::unique_ptr<SdoToInterfaceWriter> writer =
        std::move(service_variable_to_interface_writers[i]);
    const intrinsic_proto::fieldbus::v1::ServiceVariable&
        service_variable_config = service_variable_read.sdo_variable();
    INTR_ASSIGN_OR_RETURN(auto variable,
                          variable_registry.GetServiceVariable(
                              service_variable_config.index(),
                              service_variable_config.sub_index(),
                              service_variable_config.bus_position()));

    auto parse_sdo = [variable, service_variable_read,
                      writer = std::move(writer)]() mutable
        -> absl::StatusOr<intrinsic_proto::fieldbus::v1::SdoStatus> {
      intrinsic_proto::fieldbus::v1::SdoStatus service_variable_status;
      *service_variable_status.mutable_sdo_variable() =
          service_variable_read.sdo_variable();
      const auto type = service_variable_read.type();
      service_variable_status.set_type(type);
      service_variable_status.set_alias(service_variable_read.alias());
      absl::Time service_variable_read_time = absl::Now();

      switch (type) {
        case intrinsic_proto::fieldbus::v1::BOOL_SDO_TYPE: {
          INTR_ASSIGN_OR_RETURN(const auto value, variable.Read<bool>());
          service_variable_read_time = absl::Now();
          if (writer != nullptr) {
            writer->Write(static_cast<double>(value));
          }
          service_variable_status.mutable_value()->set_bool_value(value);
          const auto iterpretations =
              ServiceVariableDevice::InterpretationsOfValue(
                  service_variable_read, value);
          service_variable_status.mutable_value_interpretations()->Assign(
              iterpretations.begin(), iterpretations.end());
          break;
        }
        case intrinsic_proto::fieldbus::v1::UINT8_SDO_TYPE: {
          INTR_ASSIGN_OR_RETURN(auto value, variable.Read<uint8_t>());
          service_variable_read_time = absl::Now();
          if (writer != nullptr) {
            writer->Write(static_cast<double>(value));
          }
          service_variable_status.mutable_value()->set_uint8_value(value);
          const auto iterpretations =
              ServiceVariableDevice::InterpretationsOfValue(
                  service_variable_read, value);
          service_variable_status.mutable_value_interpretations()->Assign(
              iterpretations.begin(), iterpretations.end());
          break;
        }
        case intrinsic_proto::fieldbus::v1::INT8_SDO_TYPE: {
          INTR_ASSIGN_OR_RETURN(auto value, variable.Read<int8_t>());
          service_variable_read_time = absl::Now();
          if (writer != nullptr) {
            writer->Write(static_cast<double>(value));
          }
          service_variable_status.mutable_value()->set_int8_value(value);
          const auto iterpretations =
              ServiceVariableDevice::InterpretationsOfValue(
                  service_variable_read, value);
          service_variable_status.mutable_value_interpretations()->Assign(
              iterpretations.begin(), iterpretations.end());
          break;
        }
        case intrinsic_proto::fieldbus::v1::UINT16_SDO_TYPE: {
          INTR_ASSIGN_OR_RETURN(auto value, variable.Read<uint16_t>());
          service_variable_read_time = absl::Now();
          if (writer != nullptr) {
            writer->Write(static_cast<double>(value));
          }
          service_variable_status.mutable_value()->set_uint16_value(value);
          const auto iterpretations =
              ServiceVariableDevice::InterpretationsOfValue(
                  service_variable_read, value);
          service_variable_status.mutable_value_interpretations()->Assign(
              iterpretations.begin(), iterpretations.end());
          break;
        }
        case intrinsic_proto::fieldbus::v1::INT16_SDO_TYPE: {
          INTR_ASSIGN_OR_RETURN(auto value, variable.Read<int16_t>());
          service_variable_read_time = absl::Now();
          if (writer != nullptr) {
            writer->Write(static_cast<double>(value));
          }
          service_variable_status.mutable_value()->set_int16_value(value);
          const auto iterpretations =
              ServiceVariableDevice::InterpretationsOfValue(
                  service_variable_read, value);
          service_variable_status.mutable_value_interpretations()->Assign(
              iterpretations.begin(), iterpretations.end());
          break;
        }
        case intrinsic_proto::fieldbus::v1::UINT32_SDO_TYPE: {
          INTR_ASSIGN_OR_RETURN(auto value, variable.Read<uint32_t>());
          service_variable_read_time = absl::Now();
          if (writer != nullptr) {
            writer->Write(static_cast<double>(value));
          }
          service_variable_status.mutable_value()->set_uint32_value(value);
          const auto iterpretations =
              ServiceVariableDevice::InterpretationsOfValue(
                  service_variable_read, value);
          service_variable_status.mutable_value_interpretations()->Assign(
              iterpretations.begin(), iterpretations.end());
          break;
        }
        case intrinsic_proto::fieldbus::v1::INT32_SDO_TYPE: {
          INTR_ASSIGN_OR_RETURN(auto value, variable.Read<int32_t>());
          service_variable_read_time = absl::Now();
          if (writer != nullptr) {
            writer->Write(static_cast<double>(value));
          }
          service_variable_status.mutable_value()->set_int32_value(value);
          const auto iterpretations =
              ServiceVariableDevice::InterpretationsOfValue(
                  service_variable_read, value);
          service_variable_status.mutable_value_interpretations()->Assign(
              iterpretations.begin(), iterpretations.end());
          break;
        }
        case intrinsic_proto::fieldbus::v1::UINT64_SDO_TYPE: {
          INTR_ASSIGN_OR_RETURN(auto value, variable.Read<uint64_t>());
          service_variable_read_time = absl::Now();
          if (writer != nullptr) {
            writer->Write(static_cast<double>(value));
          }
          service_variable_status.mutable_value()->set_uint64_value(value);
          const auto iterpretations =
              ServiceVariableDevice::InterpretationsOfValue(
                  service_variable_read, value);
          service_variable_status.mutable_value_interpretations()->Assign(
              iterpretations.begin(), iterpretations.end());
          break;
        }
        case intrinsic_proto::fieldbus::v1::INT64_SDO_TYPE: {
          INTR_ASSIGN_OR_RETURN(auto value, variable.Read<int64_t>());
          service_variable_read_time = absl::Now();
          if (writer != nullptr) {
            writer->Write(static_cast<double>(value));
          }
          service_variable_status.mutable_value()->set_int64_value(value);
          const auto iterpretations =
              ServiceVariableDevice::InterpretationsOfValue(
                  service_variable_read, value);
          service_variable_status.mutable_value_interpretations()->Assign(
              iterpretations.begin(), iterpretations.end());
          break;
        }
        case intrinsic_proto::fieldbus::v1::FLOAT_SDO_TYPE: {
          INTR_ASSIGN_OR_RETURN(auto value, variable.Read<float>());
          service_variable_read_time = absl::Now();
          if (writer != nullptr) {
            writer->Write(static_cast<double>(value));
          }
          service_variable_status.mutable_value()->set_float_value(value);
          break;
        }
        case intrinsic_proto::fieldbus::v1::DOUBLE_SDO_TYPE: {
          INTR_ASSIGN_OR_RETURN(auto value, variable.Read<double>());
          service_variable_read_time = absl::Now();
          if (writer != nullptr) {
            writer->Write(static_cast<double>(value));
          }
          service_variable_status.mutable_value()->set_double_value(value);
          break;
        }
        default:
          return absl::InvalidArgumentError(absl::StrCat(
              "Unsupported value type for service variable with "
              "alias: '",
              service_variable_status.alias(), "', index: ",
              absl::Hex(service_variable_status.sdo_variable().index()),
              ", subindex: ",
              absl::Hex(service_variable_status.sdo_variable().sub_index()),
              " and pos: ",
              service_variable_status.sdo_variable().bus_position()));
      }

      *service_variable_status.mutable_read_time() =
          intrinsic::FromAbslTimeClampToValidRange(absl::Now());
      return service_variable_status;
    };
    service_variable_readers.push_back(std::move(parse_sdo));
  }
  const std::string topic =
      absl::StrCat(kEcatTopicPrefix, module_name, kSdoExportTopicSuffix);
  const bool should_publish =
      export_type == ServiceVariableExportType::SDO_EXPORT_TYPE_DEFAULT ||
      export_type ==
          ServiceVariableExportType::SDO_EXPORT_TYPE_PUBSUB_AND_LOG ||
      export_type ==
          ServiceVariableExportType::SDO_EXPORT_TYPE_ANALOG_INPUT_AND_PUBSUB ||
      export_type == ServiceVariableExportType::
                         SDO_EXPORT_TYPE_ANALOG_INPUT_AND_PUBSUB_AND_LOG;
  const bool should_log =
      export_type == ServiceVariableExportType::SDO_EXPORT_TYPE_LOG ||
      export_type ==
          ServiceVariableExportType::SDO_EXPORT_TYPE_PUBSUB_AND_LOG ||
      export_type == ServiceVariableExportType::
                         SDO_EXPORT_TYPE_ANALOG_INPUT_AND_PUBSUB_AND_LOG;
  intrinsic::Thread
  thread([export_interval, topic, export_type, should_publish, should_log,
          service_variable_readers = std::move(service_variable_readers)](
             intrinsic::StopToken stop_token) mutable -> void {
    intrinsic::PubSub pubsub;
    absl::StatusOr<intrinsic::Publisher> publisher;
    if (should_publish) {
      publisher = pubsub.CreatePublisher(topic, intrinsic::TopicConfig());
      if (!publisher.ok()) {
        if (export_type ==
                ServiceVariableExportType::SDO_EXPORT_TYPE_PUBSUB_AND_LOG ||
            export_type ==
                ServiceVariableExportType::
                    SDO_EXPORT_TYPE_ANALOG_INPUT_AND_PUBSUB_AND_LOG) {
          LOG(WARNING) << "Failed to create publisher with: '"
                       << publisher.status()
                       << "', will still log service variables.";
        } else {
          LOG(ERROR) << "Failed to create publisher: " << publisher.status();
          return;
        }
      }
    }

    LOG(INFO) << "Started service variable reader with "
              << service_variable_readers.size() << " variables on topic '"
              << topic << "' and export_interval of '" << export_interval
              << "'.";

    std::vector<intrinsic_proto::fieldbus::v1::SdoStatus>
        service_variable_status_items;
    service_variable_status_items.reserve(service_variable_readers.size());
    while (!stop_token.stop_requested()) {
      service_variable_status_items.clear();
      for (auto& service_variable_reader : service_variable_readers) {
        auto status = service_variable_reader();
        if (!status.ok()) {
          LOG_EVERY_N_SEC(ERROR, 10) << status;
        } else {
          service_variable_status_items.push_back(status.value());
        }
      }

      if (should_log) {
        for (const auto& service_variable_status_item :
             service_variable_status_items) {
          std::string value_interpretation = ".";

          if (!service_variable_status_item.value_interpretations().empty()) {
            value_interpretation = absl::StrCat(
                " with interpretations: '",
                absl::StrJoin(
                    service_variable_status_item.value_interpretations(), ", "),
                "'.");
          }

          LOG(INFO)
              << "'service variable at bus position '"
              << service_variable_status_item.sdo_variable().bus_position()
              << "' at 0x"
              << absl::Hex(service_variable_status_item.sdo_variable().index())
              << "."
              << absl::Hex(
                     service_variable_status_item.sdo_variable().sub_index())
              << " with alias '" << service_variable_status_item.alias()
              << " has value '"
              << ToString(service_variable_status_item.value()) << "'"
              << value_interpretation;
        }
      }

      if (publisher.ok() && should_publish) {
        intrinsic_proto::fieldbus::v1::ServiceVariableExport
            service_variable_export;
        for (const auto& service_variable_item :
             service_variable_status_items) {
          *service_variable_export.add_sdo_status() = service_variable_item;
        }
        if (const auto pub_status = publisher->Publish(service_variable_export);
            !pub_status.ok()) {
          LOG_EVERY_N_SEC(ERROR, 10) << "Failed to publish: " << pub_status;
        }
      }
      // Prevent busy loop.
      absl::SleepFor(export_interval);
    }
  });
  if (auto status = intrinsic::SetThreadName("service_variable_reader",
                                             thread.native_handle());
      !status.ok()) {
    LOG(WARNING) << "Failed to set thread name: " << status
                 << ". Continuing without setting thread name.";
  }
  return thread;
}

}  // namespace

absl::StatusOr<std::unique_ptr<ServiceVariableDevice>>
ServiceVariableDevice::Create(
    fieldbus::DeviceInitContext& init_context,
    const intrinsic_proto::fieldbus::v1::ServiceVariableConfig& config,
    absl::string_view context_name, absl::Duration cycle_time) {
  const fieldbus::VariableRegistry& variable_registry =
      init_context.GetVariableRegistry();

  std::list<fieldbus::ServiceVariable> service_variables;
  for (const auto& service_variable_write : config.sdo_writes()) {
    const auto& service_variable = service_variable_write.sdo_variable();
    INTR_ASSIGN_OR_RETURN(auto variable, variable_registry.GetServiceVariable(
                                             service_variable.index(),
                                             service_variable.sub_index(),
                                             service_variable.bus_position()));

    auto message = absl::StrCat(
        "service variable (0x", absl::Hex(service_variable.index()), ".",
        absl::Hex(service_variable.sub_index()),
        ", bus_pos: ", service_variable.bus_position(), ") value: ");
    const auto& value = service_variable_write.value();

    switch (value.value_case()) {
      case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::
          kBoolValue: {
        INTR_ASSIGN_OR_RETURN(const auto value_change,
                              ReadWriteRead(variable, value.bool_value()));

        INTRINSIC_RT_LOG(INFO)
            << message << value_change.first << " -> " << value_change.second;
        break;
      }
      case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::
          kUint8Value: {
        // Need to check the range of the value since uint8_value is a uint32_t.
        INTR_RETURN_IF_ERROR(RangeCheck(value.uint8_value(),
                                        std::numeric_limits<uint8_t>::min(),
                                        std::numeric_limits<uint8_t>::max()));
        INTR_ASSIGN_OR_RETURN(
            auto value_change,
            ReadWriteRead(variable, static_cast<uint8_t>(value.uint8_value())));
        INTRINSIC_RT_LOG(INFO)
            << message << value_change.first << " -> " << value_change.second;
        break;
      }
      case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::
          kInt8Value: {
        // Need to check the range of the value since int8_value is a int32_t.
        INTR_RETURN_IF_ERROR(RangeCheck(value.int8_value(),
                                        std::numeric_limits<int8_t>::min(),
                                        std::numeric_limits<int8_t>::max()));
        INTR_ASSIGN_OR_RETURN(
            auto value_change,
            ReadWriteRead(variable, static_cast<int8_t>(value.int8_value())));
        INTRINSIC_RT_LOG(INFO)
            << message << value_change.first << " -> " << value_change.second;
        break;
      }
      case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::
          kUint16Value: {
        // Need to check the range of the value since uint16_value is a
        // uint32_t.
        INTR_RETURN_IF_ERROR(RangeCheck(value.uint16_value(),
                                        std::numeric_limits<uint16_t>::min(),
                                        std::numeric_limits<uint16_t>::max()));
        INTR_ASSIGN_OR_RETURN(
            auto value_change,
            ReadWriteRead(variable,
                          static_cast<uint16_t>(value.uint16_value())));
        INTRINSIC_RT_LOG(INFO)
            << message << value_change.first << " -> " << value_change.second;
        break;
      }
      case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::
          kInt16Value: {
        // Need to check the range of the value since int8_value is a int32_t.
        INTR_RETURN_IF_ERROR(RangeCheck(value.int16_value(),
                                        std::numeric_limits<int16_t>::min(),
                                        std::numeric_limits<int16_t>::max()));
        INTR_ASSIGN_OR_RETURN(
            auto value_change,
            ReadWriteRead(variable, static_cast<int16_t>(value.int16_value())));
        INTRINSIC_RT_LOG(INFO)
            << message << value_change.first << " -> " << value_change.second;
        break;
      }
      case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::
          kUint32Value: {
        INTR_ASSIGN_OR_RETURN(auto value_change,
                              ReadWriteRead(variable, value.uint32_value()));
        INTRINSIC_RT_LOG(INFO)
            << message << value_change.first << " -> " << value_change.second;
        break;
      }
      case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::
          kInt32Value: {
        INTR_ASSIGN_OR_RETURN(auto value_change,
                              ReadWriteRead(variable, value.int32_value()));
        INTRINSIC_RT_LOG(INFO)
            << message << value_change.first << " -> " << value_change.second;
        break;
      }
      case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::
          kUint64Value: {
        INTR_ASSIGN_OR_RETURN(auto value_change,
                              ReadWriteRead(variable, value.uint64_value()));
        INTRINSIC_RT_LOG(INFO)
            << message << value_change.first << " -> " << value_change.second;
        break;
      }
      case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::
          kInt64Value: {
        INTR_ASSIGN_OR_RETURN(auto value_change,
                              ReadWriteRead(variable, value.int64_value()));
        INTRINSIC_RT_LOG(INFO)
            << message << value_change.first << " -> " << value_change.second;
        break;
      }
      case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::
          kFloatValue: {
        INTR_ASSIGN_OR_RETURN(auto value_change,
                              ReadWriteRead(variable, value.float_value()));
        INTRINSIC_RT_LOG(INFO)
            << message << value_change.first << " -> " << value_change.second;
        break;
      }
      case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::
          kDoubleValue: {
        INTR_ASSIGN_OR_RETURN(auto value_change,
                              ReadWriteRead(variable, value.double_value()));
        INTRINSIC_RT_LOG(INFO)
            << message << value_change.first << " -> " << value_change.second;
        break;
      }
      case intrinsic_proto::fieldbus::v1::ServiceVariableValueType::
          VALUE_NOT_SET: {
        return absl::InvalidArgumentError(absl::StrCat(
            "No value configured for service variable with index: ",
            absl::Hex(service_variable.index()),
            ", subindex: ", absl::Hex(service_variable.sub_index()),
            " and pos: ", service_variable.bus_position()));
        break;
      }
      default: {
        return absl::InvalidArgumentError(absl::StrCat(
            "Unsupported value type for service variable with index: ",
            absl::Hex(service_variable.index()),
            ", subindex: ", absl::Hex(service_variable.sub_index()),
            " and pos: ", service_variable.bus_position()));
      }
    }
    service_variables.push_back(variable);
  }

  std::vector<std::unique_ptr<SdoToInterfaceWriter>>
      service_variable_interface_writers;
  service_variable_interface_writers.reserve(config.sdo_reads_size());
  std::vector<absl::AnyInvocable<intrinsic::icon::RealtimeStatus()>>
      interface_updaters;
  interface_updaters.reserve(config.sdo_reads_size());

  const ServiceVariableExportType export_type = config.sdo_export_type();
  const bool should_create_analog_interface =
      export_type == ServiceVariableExportType::SDO_EXPORT_TYPE_ANALOG_INPUT ||
      export_type ==
          ServiceVariableExportType::SDO_EXPORT_TYPE_ANALOG_INPUT_AND_PUBSUB ||
      export_type == ServiceVariableExportType::
                         SDO_EXPORT_TYPE_ANALOG_INPUT_AND_PUBSUB_AND_LOG;

  for (const intrinsic_proto::fieldbus::v1::ServiceVariableRead&
           service_variable_read : config.sdo_reads()) {
    if (should_create_analog_interface) {
      if (service_variable_read.alias().empty()) {
        return absl::InvalidArgumentError(absl::StrCat(
            "service variable at bus position '",
            service_variable_read.sdo_variable().bus_position(), "' at 0x",
            absl::Hex(service_variable_read.sdo_variable().index()), ".",
            absl::Hex(service_variable_read.sdo_variable().sub_index()),
            " is configured to be exported as an analog input, but has no "
            "alias configured. service variable read variables must have a "
            "unique alias "
            "name when any of the ANALOG_IO export types are enabled."));
      }
      INTR_ASSIGN_OR_RETURN(
          auto handle,
          init_context.GetInterfaceRegistry()
              .AdvertiseMutableInterface<intrinsic_fbs::AIOStatus>(
                  service_variable_read.alias(),
                  std::vector<std::string>{service_variable_read.alias()}));
      auto writer = std::make_unique<SdoToInterfaceWriter>(std::move(handle));
      interface_updaters.push_back([writer_ptr = writer.get()]() {
        return writer_ptr->UpdateInterface();
      });
      service_variable_interface_writers.push_back(std::move(writer));
    } else {
      service_variable_interface_writers.push_back(nullptr);
    }
  }

  intrinsic::Thread service_variable_read_thread;
  if (config.sdo_reads_size() > 0) {
    INTR_ASSIGN_OR_RETURN(
        service_variable_read_thread,
        SetupSdoPublisher(context_name, variable_registry, config, cycle_time,
                          std::move(service_variable_interface_writers)));
  }

  return absl::WrapUnique(new ServiceVariableDevice(
      std::move(service_variable_read_thread),
      std::move(interface_updaters)));
}

ServiceVariableDevice::ServiceVariableDevice(
    intrinsic::Thread service_variable_read_thread,
    std::vector<absl::AnyInvocable<intrinsic::icon::RealtimeStatus()>>
        interface_updaters)
    : sdo_read_thread_(std::move(service_variable_read_thread)),
      interface_updaters_(std::move(interface_updaters))
{}

ServiceVariableDevice::~ServiceVariableDevice() = default;

intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus>
ServiceVariableDevice::CyclicRead(fieldbus::RequestType) {
  for (auto& updater : interface_updaters_) {
    INTRINSIC_RT_RETURN_IF_ERROR(updater());
  }
  return fieldbus::RequestStatus::kDone;
}

intrinsic::icon::RealtimeStatusOr<fieldbus::RequestStatus>
ServiceVariableDevice::CyclicWrite(fieldbus::RequestType) {
  return fieldbus::RequestStatus::kDone;
}

}  // namespace intrinsic::fieldbus
