// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_CC_CLIENT_CLIENT_H_
#define INTRINSIC_ICON_CC_CLIENT_CLIENT_H_

#include <memory>
#include <optional>
#include <string>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/time/time.h"
#include "absl/types/span.h"
#include "intrinsic/icon/cc_client/operational_status.h"
#include "intrinsic/icon/cc_client/robot_config.h"
#include "intrinsic/icon/common/part_properties.h"
#include "intrinsic/icon/common/slot_part_map.h"
#include "intrinsic/icon/control/logging_mode.h"
#include "intrinsic/icon/proto/part_status.pb.h"
#include "intrinsic/icon/proto/v1/service.grpc.pb.h"
#include "intrinsic/icon/proto/v1/service.pb.h"
#include "intrinsic/icon/proto/v1/types.pb.h"
#include "intrinsic/util/grpc/channel_interface.h"
#include "intrinsic/world/robot_payload/robot_payload.h"

// This header defines the ICON Application Layer C++ client library, which is a
// thin wrapper around the ICON Application Layer GRPC Service.

namespace intrinsic {
namespace icon {

// Default timeout for client GRPC requests.
constexpr absl::Duration kClientDefaultTimeout = absl::Seconds(20);

// A client for the ICON Application Layer GRPC Service.
//
// This object is moveable but not copyable.
class Client {
 public:
  enum HardwareGroup {
    kAllHardware,
    kOperationalHardwareOnly,
  };

  // Constructs a Client that uses the provided `icon_channel`.
  //
  // The factory returned by `icon_channel.GetClientContextFactory()` is invoked
  // before each gRPC request to obtain a ::grpc::ClientContext.  This provides
  // an opportunity to set client metadata, or other ClientContext settings, for
  // all ICON API requests made by the Client.
  explicit Client(std::shared_ptr<ChannelInterface> icon_channel);

  // Constructs a Client client that wraps `stub`.
  //
  // The resulting client uses `client_context_factory()` to obtain
  // ::grpc::ClientContext objects before each gRPC request.
  explicit Client(
      std::unique_ptr<intrinsic_proto::icon::v1::IconApi::StubInterface> stub,
      ClientContextFactory client_context_factory =
          DefaultClientContextFactory);

  // Makes a request to the server to get an Action Sigature by action type
  // name.
  //
  // Returns `NotFoundError` if the action type name is not found if the action
  // type name is not found.
  // Propagates gRPC communication errors.
  absl::StatusOr<intrinsic_proto::icon::v1::ActionSignature>
  GetActionSignatureByName(absl::string_view action_type_name) const;

  // Makes a request to the server to get part-specific config properties.
  //
  // These are fixed properties for the lifetime of the server (for example, the
  // number of DOFs for a robot arm.)
  //
  // Example:
  //
  //  INTR_ASSIGN_OR_RETURN(RobotConfig robot_config,
  //                   icon_client.GetConfig());
  //  INTR_ASSIGN_OR_RETURN(intrinsic_proto::icon::GenericPartConfig
  //  part_config,
  //                   robot_config.GetGenericPartConfig("robot_arm"));
  absl::StatusOr<RobotConfig> GetConfig() const;

  // Requests a restart of the entire server after all sessions are closed.
  //
  // If sessions are open, the restart will be delayed.
  // Devices are disabled similar to when an application is stopped.
  absl::Status RestartServer() const;

  // Makes a request to the server to get a snapshot of the server-side status,
  // including part-specific status.
  //
  // Example:
  //
  //  INTR_ASSIGN_OR_RETURN(intrinsic_proto::icon::v1::GetStatusResponse
  //  robot_status,
  //                   icon_client.GetRobotStatus());
  //  intrinsic_proto::icon::PartStatus my_part_status =
  //  robot_status.part_status.at("my_part");
  absl::StatusOr<intrinsic_proto::icon::v1::GetStatusResponse> GetStatus()
      const;

  // Makes a request to the server to get a snapshot of the server-side status,
  // including part-specific status, then looks up the Part status for
  // `part_name`.
  //
  // Each call makes a new requests, so *do not use this* if you are
  //
  // a) Interested in the server-wide data contained in
  //    intrinsic_proto::icon::v1::GetStatusResponse
  // b) Inspecting data for multiple parts
  //
  // Returns NotFoundError if a response was received, but does not contain a
  // status for `part_name`.
  absl::StatusOr<intrinsic_proto::icon::PartStatus> GetSinglePartStatus(
      absl::string_view part_name) const;

  // Makes a request to the server to determine if action type
  // `action_type_name` is compatible with part `part_name`.
  absl::StatusOr<bool> IsActionCompatible(
      absl::string_view part_name, absl::string_view action_type_name) const;

  // Makes a request to the server to determine if action type
  // `action_type_name` is compatible with the part assignments specified in
  // SlotPartMap.
  absl::StatusOr<bool> IsActionCompatible(
      const SlotPartMap& slot_part_map,
      absl::string_view action_type_name) const;

  // Makes a request to the server to list all available Action Signatures. The
  // results are sorted by Action type name.
  absl::StatusOr<std::vector<intrinsic_proto::icon::v1::ActionSignature>>
  ListActionSignatures() const;

  // Makes a request to the server to list all parts that are compatible with
  // all listed action types. If `action_type_names` is empty, returns all
  // parts.
  absl::StatusOr<std::vector<std::string>> ListCompatibleParts(
      absl::Span<const std::string> action_type_names) const;

  // Makes a request to the server to list all available part names.
  absl::StatusOr<std::vector<std::string>> ListParts() const;

  // Enables all parts on the server, which performs all steps necessary to get
  // the parts ready to receive commands.
  // Since the server auto-enables at startup, this is only needed after a
  // call to Disable().
  //
  // NOTE: Enabling a server is something the user does directly. DO NOT call
  // this from library code automatically to make things more convenient. Human
  // users must be able to rely on the robot to stay still unless they enable
  // it.
  //
  // If the operational state of the server is already kEnabled, then this does
  // nothing and returns absl::OkStatus(). Returns an error if the server is
  // faulted.
  absl::Status Enable() const;

  // Disables all parts on the server. Ends all currently-active sessions.
  //
  // NOTE: Disabling a server is something the user does directly. DO NOT call
  // this from library code automatically to make things more convenient. Human
  // users must be able to rely on the robot to stay enabled unless they
  // explicitly disable it (or the robot encounters a fault).
  //
  // With `kOperationalHardwareOnly`, parts that only use hardware
  // modules that are configured with
  // `IconMainConfig.hardware_config.cell_control_hardware` will be skipped,
  // keeping them enabled if they are enabled.
  // One use case is to integrate cell-level control where
  // operational robot hardware can be paused such that automatic
  // mode is not needed, while still reading/writing input/output on a fieldbus
  // hardware module for cell-level control.
  // By default, the `Disable` request disables all hardware.
  //
  // If no state changes are needed (because the selected group of hardware is
  // already disabled), then this does nothing and returns absl::OkStatus().
  // Returns an error if the server is faulted.
  absl::Status Disable(HardwareGroup group = kAllHardware) const;

  // Clears all faults and returns the server to an enabled state. Returns OK if
  // faults were successfully cleared.
  //
  // NOTE: ICON automatically enables after clearing faults!
  //
  // NOTE: Clearing faults is something the user does directly. DO NOT call this
  // from library code automatically to make things more convenient. ICON will
  // automatically re-enable when faults are cleared! Human users must be able
  // to rely on the robot to stay still unless they explicitly clear the
  // fault(s).
  //
  // Some classes of faults (internal server errors, or issues that have a
  // physical root cause) may require additional server- or hardware-specific
  // mitigation before ClearFaults() can successfully clear the fault.
  // If a server restart is needed, returns `ResourceExhausted`; in this case,
  // the client should retry until receiving `OkStatus`.
  absl::Status ClearFaults() const;

  // Returns the summarized state of the server.
  // This is the status of all hardware and the server.
  // It can differ from `GetCellControlHardwareStatus`, which is the state of
  // a subset of hardware.
  absl::StatusOr<OperationalStatus> GetOperationalStatus() const;

  // Returns the status of cell control hardware, which is marked with
  // `IconMainConfig.hardware_config.cell_control_hardware`.
  // Cell control hardware is a group of hardware modules that does not inherit
  // faults from operational hardware, and only gets disabled when a
  // any cell control hardware module faults (or when `Disable` is called).
  absl::StatusOr<OperationalStatus> GetCellControlHardwareStatus() const;

  absl::Status SetSpeedOverride(double new_speed_override);
  absl::StatusOr<double> GetSpeedOverride() const;

  absl::Status SetLoggingMode(LoggingMode logging_mode) const;

  absl::StatusOr<LoggingMode> GetLoggingMode() const;

  absl::Status SetPartProperties(const PartPropertyMap& property_map) const;
  absl::StatusOr<TimestampedPartProperties> GetPartProperties() const;

 private:
  // Hold onto the channel, if any, so that callers do not need to worry about
  // its lifetime.
  std::shared_ptr<ChannelInterface> channel_;

  // The stub for communicating with the backend.
  std::unique_ptr<intrinsic_proto::icon::v1::IconApi::StubInterface> stub_;

  // The timeout to apply to all requests made through this Client.
  absl::Duration timeout_;

  // Factory function that produces ::grpc::ClientContext objects before each
  // gRPC request.
  ClientContextFactory client_context_factory_;
};

}  // namespace icon
}  // namespace intrinsic

#endif  // INTRINSIC_ICON_CC_CLIENT_CLIENT_H_
