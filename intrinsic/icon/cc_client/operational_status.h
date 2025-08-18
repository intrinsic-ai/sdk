// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_CC_CLIENT_OPERATIONAL_STATUS_H_
#define INTRINSIC_ICON_CC_CLIENT_OPERATIONAL_STATUS_H_

#include <iostream>
#include <ostream>
#include <string>

#include "absl/base/attributes.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/proto/v1/types.pb.h"

namespace intrinsic::icon {

// Types representing the status of the server or a group of hardware.
//
// Use `Client::Enable()`, `Client::Disable()`,
// `Client::ClearFaults()` and `Client::GetOperationalStatus()` to get
// and set the operational state. See client.h for details.

// The summarized state of a group of hardware (operational hardware modules or
// cell control hardware modules) or the real-time control service.
enum class OperationalState {
  // Indicates that this group of hardware (or the server) is not ready for
  // active control and that no sessions can be started that need to control
  // these parts.
  // This is possible when:
  // - The skill "disable_realtime_control" or `Client::Disable()` were called.
  // - The server or hardware is starting up.
  // - Faults are being cleared.
  // Read-only sessions are possible for all parts.
  // Part status is being published.
  // `icon_client.Enable()` can be called to enable full control for all parts.
  kDisabled,

  // Indicates that at least one part, possibly the entire real-time control
  // service, is faulted. `icon_client.ClearFaults()` is needed to re-enable
  // control.
  // Depending on the fault, real-time control may or may not be running the
  // safety actions.
  // An example for a single part fault is a robot hardware module reporting
  // an emergency stop but still being connected.
  // An example for a global fault that cannot be cleared is a mistake in the
  // hardware module names in the config.
  // An example for a global fault that can be cleared is a timeout in a
  // simulation reset.
  // If a part is not faulted, read-only sessions reading from them can
  // continue, and part status may still be published.
  kFaulted,

  // Indicates that the server is ready for a session to begin and all parts are
  // enabled.
  // Part status is being published.
  kEnabled
};

// The summarized state of a group of hardware (operational hardware modules or
// cell control hardware modules) or the real-time control service, along with a
// fault reason when the state is `kFaulted`.
class OperationalStatus final {
 public:
  // Constructs an OperationalStatus with state set to `kDisabled`
  OperationalStatus();

  // Creates an `OperationalStatus` object with state set to `kDisabled`.
  static OperationalStatus Disabled();
  // Creates an `OperationalStatus` object with state set to `kFaulted`.
  // `reason` is a human-readable description of what caused the fault.
  static OperationalStatus Faulted(absl::string_view reason);
  // Creates an `OperationalStatus` object with state set to `kEnabled`.
  static OperationalStatus Enabled();

  // Returns the operational state.
  OperationalState state() const { return state_; }

  // Returns a human-readable description of what caused the `kFaulted` state.
  // When not in the `kFaulted` state, returns an empty string.
  std::string fault_reason() const { return fault_reason_; }

  bool operator==(const OperationalStatus& other) const {
    return state_ == other.state_ && fault_reason_ == other.fault_reason_;
  }

 private:
  // Private constructor. Use static methods (`Disabled()`, etc.) to create an
  // `OperationalStatus` object.
  OperationalStatus(OperationalState state, absl::string_view fault_reason);

  OperationalState state_;
  std::string fault_reason_;
};

// These convenience functions return `true` if a given status matches the
// OperationalState of its associated function.
ABSL_MUST_USE_RESULT bool IsDisabled(const OperationalStatus& status);
ABSL_MUST_USE_RESULT bool IsFaulted(const OperationalStatus& status);
ABSL_MUST_USE_RESULT bool IsEnabled(const OperationalStatus& status);

// Converts an OperationalState to a string. For example,
// `ToString(OperationalState::kDisabled)` returns `"DISABLED"`.
std::string ToString(OperationalState state);

// Converts an OpertationalStatus to a string.
// If `IsDisabled(status)` returns `"DISABLED"`.
// If `IsFaulted(status)` returns `"FAULTED(reason)"` where `reason` is
// `status.fault_reason()`.
// If `IsEnabled(status)` returns `"ENABLED"`.
std::string ToString(const OperationalStatus& status);

// operator<<
//
// Prints a human-readable representation of `state` to `os`.
std::ostream& operator<<(std::ostream& os, OperationalState state);

// operator<<
//
// Prints a human-readable representation of `status` to `os`.
std::ostream& operator<<(std::ostream& os, const OperationalStatus& status);

// These functions convert the types declared in this header to/from proto.
intrinsic_proto::icon::v1::OperationalState ToProto(OperationalState state);
intrinsic_proto::icon::v1::OperationalStatus ToProto(
    const OperationalStatus& status);
absl::StatusOr<OperationalState> FromProto(
    const intrinsic_proto::icon::v1::OperationalState& proto);
absl::StatusOr<OperationalStatus> FromProto(
    const intrinsic_proto::icon::v1::OperationalStatus& proto);

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_CC_CLIENT_OPERATIONAL_STATUS_H_
