// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_HARDWARE_GPIO_GPIO_SERVICE_PROTO_UTILS_H_
#define INTRINSIC_HARDWARE_GPIO_GPIO_SERVICE_PROTO_UTILS_H_

#include "intrinsic/hardware/gpio/v1/gpio_service.pb.h"
#include "intrinsic/hardware/gpio/v1/signal.pb.h"

namespace intrinsic_proto::gpio::v1 {
// Inject equality operator in the proto's namespace so that argument lookup
// just works.

// Returns true if the signal values are exactly the same. Uninitialized values
// are not considered equal.
bool operator==(const SignalValue& lhs, const SignalValue& rhs);

inline bool operator!=(const SignalValue& lhs, const SignalValue& rhs) {
  return !(lhs == rhs);
}

// Returns true if the signal value sets have the same signal name to value
// unordered map. Uninitialized values are not considered equal.
bool operator==(const SignalValueSet& lhs, const SignalValueSet& rhs);

inline bool operator!=(const SignalValueSet& lhs, const SignalValueSet& rhs) {
  return !(lhs == rhs);
}

// Two instances of `OpenWriteSessionRequest` are considered equal if they
// contain the same set of names (i.e. the order in which the signal names
// appear does not matter).
bool operator==(const intrinsic_proto::gpio::v1::OpenWriteSessionRequest& lhs,
                const intrinsic_proto::gpio::v1::OpenWriteSessionRequest& rhs);

inline bool operator!=(
    const intrinsic_proto::gpio::v1::OpenWriteSessionRequest& lhs,
    const intrinsic_proto::gpio::v1::OpenWriteSessionRequest& rhs) {
  return !(lhs == rhs);
}

// Two instances of `ReadSignalsRequest` are considered equal if they contain
// the same set of names (i.e. the order in which the signal names appear does
// not matter).
bool operator==(const intrinsic_proto::gpio::v1::ReadSignalsRequest& lhs,
                const intrinsic_proto::gpio::v1::ReadSignalsRequest& rhs);

inline bool operator!=(
    const intrinsic_proto::gpio::v1::ReadSignalsRequest& lhs,
    const intrinsic_proto::gpio::v1::ReadSignalsRequest& rhs) {
  return !(lhs == rhs);
}

}  // namespace intrinsic_proto::gpio::v1

namespace intrinsic::gpio {

// Returns true if the signal values are exactly the same for non-floating point
// types and approximately the same for floating point types.
// Notes about equality operation:
// - Uninitialized values are not considered equal.
// - A hard-coded tolerance is used to compare floating point types.
bool SignalValuesAreApproxEqual(
    const intrinsic_proto::gpio::v1::SignalValue& a,
    const intrinsic_proto::gpio::v1::SignalValue& b);

// Returns a false boolean signal value.
intrinsic_proto::gpio::v1::SignalValue SignalFalseValue();

// Returns a true boolean signal value.
intrinsic_proto::gpio::v1::SignalValue SignalTrueValue();

intrinsic_proto::gpio::v1::SignalType SignalTypeFromValue(
    const intrinsic_proto::gpio::v1::SignalValue& value);
}  // namespace intrinsic::gpio

#endif  // INTRINSIC_HARDWARE_GPIO_GPIO_SERVICE_PROTO_UTILS_H_
