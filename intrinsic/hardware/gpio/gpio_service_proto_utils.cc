// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/hardware/gpio/gpio_service_proto_utils.h"

#include <algorithm>
#include <cmath>
#include <string>
#include <vector>

#include "absl/container/flat_hash_map.h"
#include "intrinsic/hardware/gpio/v1/signal.pb.h"

namespace intrinsic_proto::gpio::v1 {

bool operator==(const SignalValue& lhs, const SignalValue& rhs) {
  if (lhs.value_case() != rhs.value_case()) {
    return false;
  }

  switch (lhs.value_case()) {
    case SignalValue::ValueCase::kBoolValue:
      return lhs.bool_value() == rhs.bool_value();
    case SignalValue::ValueCase::kUnsignedIntValue:
      return lhs.unsigned_int_value() == rhs.unsigned_int_value();
    case SignalValue::ValueCase::kIntValue:
      return lhs.int_value() == rhs.int_value();
    case SignalValue::ValueCase::kFloatValue: {
      return lhs.float_value() == rhs.float_value();
    }
    case SignalValue::ValueCase::kDoubleValue: {
      return lhs.double_value() == rhs.double_value();
    }
    case SignalValue::ValueCase::kInt8Value:
      return lhs.int8_value().value() == rhs.int8_value().value();
    case SignalValue::ValueCase::kUnsignedInt8Value:
      return lhs.unsigned_int8_value().value() ==
             rhs.unsigned_int8_value().value();
    case SignalValue::ValueCase::VALUE_NOT_SET:
      return false;
  }

  return false;
}

bool operator==(const SignalValueSet& lhs, const SignalValueSet& rhs) {
  if (lhs.values().size() != rhs.values().size()) {
    return false;
  }

  auto MakeSignalToValueMap = [](const SignalValueSet& value_set)
      -> absl::flat_hash_map<std::string, SignalValue> {
    absl::flat_hash_map<std::string, SignalValue> signal_map;
    for (const auto& [name, value] : value_set.values()) {
      signal_map[name] = value;
    }
    return signal_map;
  };

  return MakeSignalToValueMap(lhs) == MakeSignalToValueMap(rhs);
}

bool operator==(const OpenWriteSessionRequest& lhs,
                const OpenWriteSessionRequest& rhs) {
  auto SortedInitialSessionSignals =
      [](const OpenWriteSessionRequest& req) -> std::vector<std::string> {
    // Technically, signal names can be non-unique. While that should not make
    // any difference in practice, let's return a sorted vector (instead of
    // a set) and not make any assumptions.
    std::vector<std::string> names;
    names.reserve(req.initial_session_data().signal_names().size());
    for (const auto& name : req.initial_session_data().signal_names()) {
      names.push_back(name);
    }
    std::sort(names.begin(), names.end());
    return names;
  };

  return (SortedInitialSessionSignals(lhs) ==
          SortedInitialSessionSignals(rhs)) &&
         (lhs.write_signals().signal_values() ==
          rhs.write_signals().signal_values());
}

bool operator==(const ReadSignalsRequest& lhs, const ReadSignalsRequest& rhs) {
  auto sorted_signal_names =
      [](const ReadSignalsRequest& req) -> std::vector<std::string> {
    std::vector<std::string> names;
    for (const auto& name : req.signal_names()) {
      names.push_back(name);
    }
    std::sort(names.begin(), names.end());
    return names;
  };

  return sorted_signal_names(lhs) == sorted_signal_names(rhs);
}

}  // namespace intrinsic_proto::gpio::v1

namespace intrinsic::gpio {

using ::intrinsic_proto::gpio::v1::SignalValue;

bool SignalValuesAreApproxEqual(const SignalValue& a, const SignalValue& b) {
  if (a.value_case() != b.value_case()) {
    return false;
  }

  switch (a.value_case()) {
    case intrinsic_proto::gpio::v1::SignalValue::ValueCase::kBoolValue:
      return a.bool_value() == b.bool_value();
    case intrinsic_proto::gpio::v1::SignalValue::ValueCase::kUnsignedIntValue:
      return a.unsigned_int_value() == b.unsigned_int_value();
    case intrinsic_proto::gpio::v1::SignalValue::ValueCase::kIntValue:
      return a.int_value() == b.int_value();
    case intrinsic_proto::gpio::v1::SignalValue::ValueCase::kFloatValue: {
      const auto kTolerance = 1e-3f;
      return std::fabs(a.float_value() - b.float_value()) < kTolerance;
    }
    case intrinsic_proto::gpio::v1::SignalValue::ValueCase::kDoubleValue: {
      const auto kTolerance = 1e-3;
      return std::fabs(a.double_value() - b.double_value()) < kTolerance;
    }
    case intrinsic_proto::gpio::v1::SignalValue::ValueCase::kInt8Value:
      return a.int8_value().value() == b.int8_value().value();
    case intrinsic_proto::gpio::v1::SignalValue::ValueCase::kUnsignedInt8Value:
      return a.unsigned_int8_value().value() == b.unsigned_int8_value().value();
    case intrinsic_proto::gpio::v1::SignalValue::ValueCase::VALUE_NOT_SET:
      return false;
  }

  return false;
}

intrinsic_proto::gpio::v1::SignalValue SignalFalseValue() {
  intrinsic_proto::gpio::v1::SignalValue value;
  value.set_bool_value(false);
  return value;
}

intrinsic_proto::gpio::v1::SignalValue SignalTrueValue() {
  intrinsic_proto::gpio::v1::SignalValue value;
  value.set_bool_value(true);
  return value;
}

intrinsic_proto::gpio::v1::SignalType SignalTypeFromValue(
    const intrinsic_proto::gpio::v1::SignalValue& value) {
  using ::intrinsic_proto::gpio::v1::SignalType;
  using ::intrinsic_proto::gpio::v1::SignalValue;
  switch (value.value_case()) {
    case SignalValue::kBoolValue:
      return SignalType::SIGNAL_TYPE_BOOL;
    case SignalValue::kUnsignedIntValue:
      return SignalType::SIGNAL_TYPE_UNSIGNED_INT;
    case SignalValue::kIntValue:
      return SignalType::SIGNAL_TYPE_INT;
    case SignalValue::kFloatValue:
      return SignalType::SIGNAL_TYPE_FLOAT;
    case SignalValue::kDoubleValue:
      return SignalType::SIGNAL_TYPE_DOUBLE;
    case SignalValue::kInt8Value:
      return SignalType::SIGNAL_TYPE_INT8;
    case SignalValue::kUnsignedInt8Value:
      return SignalType::SIGNAL_TYPE_UNSIGNED_INT8;
    case SignalValue::VALUE_NOT_SET:
      return SignalType::SIGNAL_TYPE_UNKNOWN;
  }
}

}  // namespace intrinsic::gpio
