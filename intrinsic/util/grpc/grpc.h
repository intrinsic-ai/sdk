// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_GRPC_GRPC_H_
#define INTRINSIC_UTIL_GRPC_GRPC_H_

#include <cstdint>
#include <memory>
#include <optional>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/notification.h"
#include "absl/time/time.h"
#include "grpcpp/grpcpp.h"
#include "intrinsic/connect/cc/grpc/channel.h"  // IWYU pragma: export

namespace intrinsic {
using ::intrinsic::connect::DefaultGrpcChannelArgs;
using ::intrinsic::connect::kGrpcClientConnectDefaultTimeout;
using ::intrinsic::connect::UnlimitedMessageSizeGrpcChannelArgs;
using ::intrinsic::connect::WaitForChannelConnected;

// Default client-side timeout for invoking services (when configuring the
// client context with ConfigureClientContext().
constexpr absl::Duration kGrpcClientServiceCallDefaultTimeout =
    absl::Seconds(60);

// Options for CreateServer calls.
struct CreateServerOptions {
  // Set the maximum receive message size.
  std::optional<int> max_receive_message_size;
};

/**
 * Create a grpc server using the listen port on the default interface
 * and the set of services provided
 */
absl::StatusOr<std::unique_ptr<::grpc::Server>> CreateServer(
    uint16_t listen_port, const std::vector<::grpc::Service*>& services,
    const CreateServerOptions& options = CreateServerOptions());

/**
 * Create a grpc server using a specific address to listen to.
 */
absl::StatusOr<std::unique_ptr<::grpc::Server>> CreateServer(
    absl::string_view address, const std::vector<::grpc::Service*>& services,
    const CreateServerOptions& options = CreateServerOptions());

/**
 * Apply the default configuration of our project to the given ClientContext.
 *
 * Configurations set on the context:
 * - enable initial waiting for a connection to the service
 * - set a fixed maximum deadline (see kGrpcClientServiceCallDefaultTimeout),
 *   override if you know that this can be shorter, or if it needs to be
 *   longer. In particular mind deadlines for wait calls (e.g., long-running
 *   operations) and streaming calls which often need to be increased or be
 *   set to absl::InfiniteFuture().
 */
void ConfigureClientContext(::grpc::ClientContext* client_context);

// Parameters to configure the shutdown behavior of a gRPC server.
struct ShutdownParams {
  // Duration to wait for the grpc's health service state (if relevant) to
  // propagate to the load balancers.
  absl::Duration health_grace_duration;
  // Timeout passed into grpc::Server::Shutdown on a sigterm.
  absl::Duration shutdown_timeout;

  // Returns params that aggressively shutdowns the server.
  static ShutdownParams Aggressive();
};

// Registers a custom signal handler for SIGTERM, serves the server and blocks
// till it is shutdown. The custom handler is left registered when the function
// returns.
//
// `handlers_registered` notification is triggered once the signal handler is
// registered. This is mainly useful in unit tests to know when it is okay to
// raise a SIGTERM signal.
//
// Returns an error if registering the signal handler fails.
//
// Typical usage:
// ```
//   int main() {
//      auto server = intrinsic::CreateServer(...);
//      absl::Notification registered;
//      QCHECK_OK(RegisterSignalHandlerAndWait(server.get(),
//                ShutdownParams{...}, registered));
//      return EXIT_SUCCESS;
//   }
// ```
absl::Status RegisterSignalHandlerAndWait(
    ::grpc::Server* server, const ShutdownParams& params,
    absl::Notification& handlers_registered);

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_GRPC_GRPC_H_
