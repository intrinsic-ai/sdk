// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_CONNECT_CC_GRPC_CHANNEL_H_
#define INTRINSIC_CONNECT_CC_GRPC_CHANNEL_H_

#include <memory>
#include <optional>
#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "grpcpp/channel.h"
#include "grpcpp/security/credentials.h"
#include "grpcpp/support/channel_arguments.h"
#include "intrinsic/util/time/clock.h"

namespace intrinsic::connect {

// Default timeout for the initial GRPC connection made by client libraries.
inline constexpr absl::Duration kGrpcClientConnectDefaultTimeout =
    absl::Seconds(5);

// Wait for a newly created channel to be connected
absl::Status WaitForChannelConnected(absl::string_view address,
                                     std::shared_ptr<::grpc::Channel> channel,
                                     absl::Time deadline = absl::Now());

// Get recommended default gRPC channel arguments.
::grpc::ChannelArguments DefaultGrpcChannelArgs();

// Get gRPC channel arguments with unlimited send/receive message size.
// This also includes all settings from DefaultGrpcChannelArgs(). This can be
// used for services that send large messages, e.g., the geometry service.
::grpc::ChannelArguments UnlimitedMessageSizeGrpcChannelArgs();

// Apply default configuration of our project and create a new channel.
//
// DEPRECATED: Use GrpcChannel instead.
absl::StatusOr<std::shared_ptr<::grpc::Channel>> CreateClientChannel(
    absl::string_view address, absl::Time deadline,
    const ::grpc::ChannelArguments& channel_args =
        connect::DefaultGrpcChannelArgs(),
    bool use_default_application_credentials = false,
    std::optional<std::string> server_instance_name = std::nullopt);

// Builder for creating a gRPC channel.
class GrpcChannel {
 public:
  // Constructs a GrpcChannel.
  explicit GrpcChannel(absl::string_view address);

  // Constructs a GrpcChannel with a custom clock.
  //
  // Typically, this is only used when the caller needs to control the clock,
  // e.g., in tests.
  //
  // The pointer to clock will be stored for the life of the GrpcChannel, and
  // must outlive it.
  GrpcChannel(absl::string_view address, ClockInterface* clock);

  // Not copyable or movable.
  GrpcChannel(const GrpcChannel&) = delete;
  GrpcChannel& operator=(const GrpcChannel&) = delete;
  GrpcChannel(GrpcChannel&&) = delete;
  GrpcChannel& operator=(GrpcChannel&&) = delete;

  ~GrpcChannel() = default;

  // Sets the deadline for connecting to the channel.
  GrpcChannel& WithDeadline(absl::Time deadline);

  // Sets the timeout for connecting to the channel.
  //
  // The default timeout is 60 seconds.
  //
  // This is a convenience method for WithDeadline(ClockInterface::Now() +
  // timeout).
  GrpcChannel& WithTimeout(absl::Duration timeout);

  // Specifies the channel credentials to use.
  //
  // If not set, grpc::GoogleDefaultCredentials() will be used.
  GrpcChannel& WithChannelCredentials(
      std::shared_ptr<grpc::ChannelCredentials> credentials);

  // Performs a health check on the channel.
  //
  // Note: Do not call this unless providing custom channel credentials.
  // This does not work with grpc::GoogleDefaultCredentials(). Health checks
  // will fail with "UNKNOWN: Received http2 header with status: 302".
  //
  // A server instance name can be provided to the health check.
  //
  // This is disabled by default.
  GrpcChannel& WithCheckChannelHealth(
      std::optional<absl::string_view> server_instance_name = std::nullopt);

  // Sets the channel arguments to use unlimited send/receive message size.
  //
  // This also includes all settings from DefaultGrpcChannelArgs(). This can be
  // used for services that send large messages.
  GrpcChannel& WithUnlimitedMessageSizeChannelArgs();

  // Provides custom channel arguments.
  //
  // Note: No default channel arguments are applied. If you want to apply
  // modifications, or additions, to the default channel arguments, you'll need
  // to call DefaultGrpcChannelArgs(), or UnlimitedMessageSizeGrpcChannelArgs(),
  // and then apply your modifications to the returned grpc::ChannelArguments
  // to use as the arg to WithCustomChannelArgs().
  GrpcChannel& WithCustomChannelArgs(
      const grpc::ChannelArguments& channel_args);

  // Constructs the grpc::Channel and connects to the server.
  //
  // This can only be called once.
  absl::StatusOr<std::shared_ptr<grpc::Channel>> Connect();

 private:
  struct CheckChannelHealthOptions {
    std::optional<absl::string_view> server_instance_name;
  };

  std::string address_;
  ClockInterface* clock_ = RealClock::GetInstance();
  std::shared_ptr<grpc::ChannelCredentials> credentials_;
  absl::Time deadline_ = clock_->Now() + absl::Seconds(60);
  grpc::ChannelArguments channel_args_ = DefaultGrpcChannelArgs();

  // If set, check the channel health.
  std::optional<CheckChannelHealthOptions> check_channel_health_ = std::nullopt;

  bool consumed_ = false;
};

}  // namespace intrinsic::connect

#endif  // INTRINSIC_CONNECT_CC_GRPC_CHANNEL_H_
