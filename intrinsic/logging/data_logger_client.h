// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#ifndef INTRINSIC_LOGGING_DATA_LOGGER_CLIENT_H_
#define INTRINSIC_LOGGING_DATA_LOGGER_CLIENT_H_

#include <cstdint>
#include <functional>
#include <memory>

#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "absl/time/time.h"
#include "intrinsic/connect/cc/grpc/channel.h"
#include "intrinsic/logging/proto/log_item.pb.h"
#include "intrinsic/logging/proto/logger_service.grpc.pb.h"

namespace intrinsic::data_logger {

// Initializes a connection to the DataLogger gRPC server. Call once during
// program startup.
absl::Status StartUpIntrinsicLoggerViaGrpc(
    absl::string_view target_address,
    absl::Duration timeout =
        intrinsic::connect::kGrpcClientConnectDefaultTimeout);

// Initializes the DataLogger client from a stub. Call once during program
// startup. Intended for testing.
absl::Status StartUpIntrinsicLoggerViaStub(
    std::unique_ptr<intrinsic_proto::data_logger::DataLogger::StubInterface>
        stub);

// Generates a random integer with sufficient entropy to be considered globally
// unique.
uint64_t GenerateUid();

// Asynchronously logs a message.
void LogAsync(const intrinsic_proto::data_logger::LogItem& item);
void LogAsync(const intrinsic_proto::data_logger::LogItem& item,
              std::function<void(absl::Status)> callback);
void LogAsync(intrinsic_proto::data_logger::LogItem&& item);
void LogAsync(intrinsic_proto::data_logger::LogItem&& item,
              std::function<void(absl::Status)> callback);

// Sends a request to log `item` to the logging service, waits for a response
// and returns the Status. Use this instead of the async call when:
//   1) You don't mind blocking for a few ms.
//   2) You need to know whether the item was successfully logged to disk.
absl::Status LogAndAwaitResponse(
    const intrinsic_proto::data_logger::LogItem& item);
absl::Status LogAndAwaitResponse(intrinsic_proto::data_logger::LogItem&& item);

}  // namespace intrinsic::data_logger

#endif  // INTRINSIC_LOGGING_DATA_LOGGER_CLIENT_H_
