// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_HARDWARE_GPIO_GPIO_CLIENT_H_
#define INTRINSIC_HARDWARE_GPIO_GPIO_CLIENT_H_

#include <atomic>
#include <memory>
#include <string>

#include "absl/container/flat_hash_set.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/mutex.h"
#include "intrinsic/hardware/gpio/v1/gpio_service.grpc.pb.h"
#include "intrinsic/hardware/gpio/v1/gpio_service.pb.h"
#include "intrinsic/util/grpc/channel_interface.h"
#include "intrinsic/util/grpc/connection_params.h"

namespace intrinsic::gpio {

// Client class that talks to GPIO grpc service to read from and write to GPIO
// ports. On the first write call, a stream session is opened to claim exclusive
// write access to all the signals that could be written to during the session.
// The stream session is kept alive for the entire lifetime of the class
// instance to guarantee this exclusive write access. If the stream session
// fails to open on the first write, it is retried on successive write calls.
// While write calls are thread safe, only one call is allowed to proceed at any
// given time and others are blocked. Read calls are thread safe and are
// independent of the stream session.

class GPIOClient {
 public:
  // `gpio_service_name` should be name of the already configured GPIO service,
  // potentially behind an ingress.
  // `signals_to_claim` should be the union of all the signals that could be
  // written to using this client
  GPIOClient(
      std::unique_ptr<intrinsic_proto::gpio::v1::GPIOService::StubInterface>
          stub,
      absl::string_view gpio_service_name,
      const absl::flat_hash_set<std::string>& signals_to_claim);

  // Creates an instance of GPIOClient that does not attempt to create the grpc
  // client channel immediately. Instead, this is delayed till an rpc call needs
  // to be made to the GPIO service. This is useful in situations where the GPIO
  // service may not be running when an instance of GPIOClient is constructed.
  // `gpio_service_address` is the server address used in creating the client
  // channel.
  // `signals_to_claim` should be the union of all the signals that could be
  // written to using this client.
  GPIOClient(const ConnectionParams& connection_params,
             const absl::flat_hash_set<std::string>& signals_to_claim);

  // Close the stream session (if valid) in the destructor
  ~GPIOClient();

  // The client can't be copied or assigned
  GPIOClient(const GPIOClient&) = delete;
  GPIOClient& operator=(const GPIOClient&) = delete;

  // Writes the desired values for the specified signals. Write calls are thread
  // safe and only one call is allowed to communicate with the GPIO service at
  // any given time.
  // On the first write, a stream session is opened to claim exclusive write
  // access to the signals. If this fails, it is retried on consecutive writes.
  // The stream session is kept open for the lifetime of this client instance.
  //
  // `retry_on_session_error`: by default, if a previously valid session stream
  // becomes invalid (e.g. due to GPIO service being restarted), an attempt is
  // made to open a new session followed by the actual write operation.
  absl::Status Write(
      const intrinsic_proto::gpio::v1::SignalValueSet& desired_values,
      bool retry_on_session_error = true) ABSL_LOCKS_EXCLUDED(write_mutex_);

  // Returns the values read for the given signal names.
  absl::StatusOr<intrinsic_proto::gpio::v1::ReadSignalsResponse> Read(
      const intrinsic_proto::gpio::v1::ReadSignalsRequest& request);

  // Blocks till the specified condition for the signal values is met. `Write`
  // is allowed while this method blocks.
  absl::StatusOr<intrinsic_proto::gpio::v1::WaitForValueResponse> WaitForValue(
      const intrinsic_proto::gpio::v1::WaitForValueRequest& request,
      absl::Duration timeout);

  // Reads the signal values and matches them against the desired value(s).
  // Returns True when all the values match.
  absl::StatusOr<bool> ReadAndMatch(
      const intrinsic_proto::gpio::v1::SignalValueSet& match_values);

  // Returns all the signals known to the GPIO server
  absl::StatusOr<intrinsic_proto::gpio::v1::GetSignalDescriptionsResponse>
  GetSignalDescriptions();

 private:
  // Creates grpc client channel to GPIO service. This method should be called
  // before initiating an rpc call. This method is thread-safe and can be called
  // repeatedly.
  absl::Status CreateClientChannel() ABSL_LOCKS_EXCLUDED(create_channel_mutex_);

  // Internal implementation of write that opens the write session stream (if
  // necessary) and does a single write.
  absl::Status WriteInternal(
      const intrinsic_proto::gpio::v1::SignalValueSet& desired_values)
      ABSL_EXCLUSIVE_LOCKS_REQUIRED(write_mutex_);

  std::unique_ptr<intrinsic_proto::gpio::v1::GPIOService::StubInterface> stub_;

  // Connection parameters for the GPIO service to connect to.
  const ConnectionParams connection_params_;
  ClientContextFactory client_context_factory_ =
      intrinsic::DefaultClientContextFactory;

  // Request containing all the signals that can be written to during the
  // lifetime of the client instance.
  intrinsic_proto::gpio::v1::OpenWriteSessionRequest claim_signals_request_;

  // Client context for `stream_session_`. This should have the same lifetime as
  // `session_stream_`.
  std::unique_ptr<grpc::ClientContext> stream_session_ctx_
      ABSL_GUARDED_BY(write_mutex_);

  // Write session stream to claim exclusive write access to signals specified
  // in `claim_signals_request_` for the lifetime of this client instance.
  std::unique_ptr<grpc::ClientReaderWriterInterface<
      intrinsic_proto::gpio::v1::OpenWriteSessionRequest,
      intrinsic_proto::gpio::v1::OpenWriteSessionResponse>>
      session_stream_ ABSL_GUARDED_BY(write_mutex_);

  // Tears down the session and reports errors.
  absl::Status CleanUpSessionAfterError()
      ABSL_EXCLUSIVE_LOCKS_REQUIRED(write_mutex_);

  // Mutex to ensure only one write can communicate with the GPIO server at any
  // given time
  absl::Mutex write_mutex_;

  // Flag indicating if write access to signals was successfully claimed
  std::atomic<bool> signals_claimed_ = false;

  // Mutex to guard client channel creation
  absl::Mutex create_channel_mutex_;
};

};  // namespace intrinsic::gpio

#endif  // INTRINSIC_HARDWARE_GPIO_GPIO_CLIENT_H_
