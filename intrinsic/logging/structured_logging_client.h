// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_LOGGING_STRUCTURED_LOGGING_CLIENT_H_
#define INTRINSIC_LOGGING_STRUCTURED_LOGGING_CLIENT_H_

#include <functional>
#include <map>
#include <memory>
#include <optional>
#include <string>
#include <vector>

#include "absl/container/flat_hash_map.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "absl/types/span.h"
#include "grpcpp/channel.h"
#include "intrinsic/logging/proto/log_item.pb.h"
#include "intrinsic/logging/proto/logger_service.grpc.pb.h"
#include "intrinsic/logging/proto/logger_service.pb.h"

namespace intrinsic {

// A client class to interact with the structured logging service.
// The class is thread-safe. If multiple gRPC services are available at the same
// address it is recommended to share a channel between them.
class StructuredLoggingClient {
 public:
  using LogItem = ::intrinsic_proto::data_logger::LogItem;
  using LogOptions = ::intrinsic_proto::data_logger::LogOptions;
  using LoggerStub = ::intrinsic_proto::data_logger::DataLogger::StubInterface;

  struct ListResult {
    std::vector<LogItem> log_items;
    std::string next_page_token;
  };

  using GetResult = ListResult;

  // Creates a structured logging client by connecting to the specified address.
  // If the connection cannot be established when the deadline is met, the
  // function returns an error.
  static absl::StatusOr<StructuredLoggingClient> Create(
      absl::string_view address, absl::Time deadline);

  // Constructs a client from an existing gRPC channel.
  explicit StructuredLoggingClient(
      const std::shared_ptr<grpc::Channel>& channel);

  // Direct stub injection, typically used to inject mocks for testing.
  explicit StructuredLoggingClient(std::unique_ptr<LoggerStub> stub);

  StructuredLoggingClient(const StructuredLoggingClient&) = delete;
  StructuredLoggingClient& operator=(const StructuredLoggingClient&) = delete;

  // Move construction and assignment.
  StructuredLoggingClient(StructuredLoggingClient&&);
  StructuredLoggingClient& operator=(StructuredLoggingClient&&);

  ~StructuredLoggingClient();

  // Logs an r-value item.
  absl::Status Log(LogItem&& item) const;

  // Logs an item. The item will be internally copied to the logging request.
  absl::Status Log(const LogItem& item) const;

  // Performs asynchronous logging of an r-value item. A default callback is
  // installed, which prints a warning message in case of a logging failure.
  void LogAsync(LogItem&& item) const;

  // Performs asynchronous logging of an r-value item and calls the user
  // specified callback when done.
  void LogAsync(LogItem&& item,
                std::function<void(absl::Status)> callback) const;

  // Performs asynchronous logging of an item. The item will be internally
  // copied to the logging request. A default callback is
  // installed, which prints a warning message in case of a logging failure.
  void LogAsync(const LogItem& item) const;

  // Performs asynchronous logging of an item and calls the user specified
  // callback when done. The item will be internally copied to the logging
  // request.
  void LogAsync(const LogItem& item,
                std::function<void(absl::Status)> callback) const;

  // Returns a list of `event_source` that can be requested using list requests.
  absl::StatusOr<std::vector<std::string>> ListLogSources() const;

  // Returns a list of log items for the specified event source. If no data is
  // available, an empty vector is returned and the function does not generate
  // an error.
  absl::StatusOr<GetResult> GetLogItems(absl::string_view event_source) const;

  // Returns a list of log items for the specified event source. If no data is
  // available, an empty vector is returned and the function does not generate
  // an error.
  absl::StatusOr<GetResult> GetLogItems(absl::string_view event_source,
                                        absl::Time start_time,
                                        absl::Time end_time) const;

  // Returns a list of log items for the specified event source. If no data is
  // available, an empty vector is returned and the function does not generate
  // an error.
  //
  // If start_time and end_time are not specified, the default is to read all
  // data from the start of time until now.
  //
  // The function supports pagination. On each request 'page_size' items are
  // returned if that many are available. In addition a 'page_token' is returned
  // which can be used on the next request to request the next batch of items.
  //
  // Filtering is supported in the same way as documented on the logging
  // service.
  absl::StatusOr<GetResult> GetLogItems(
      absl::string_view event_source, int page_size,
      absl::string_view page_token = "",
      absl::Time start_time = absl::UniversalEpoch(),
      absl::Time end_time = absl::Now(),
      absl::flat_hash_map<std::string, std::string> filter_labels = {}) const;

  // Returns the most recent LogItem that has been logged for the given event
  // source. If no LogItem with a matching event_source has been logged since
  // --file_ttl, then NOT_FOUND will be returned instead.
  absl::StatusOr<LogItem> GetMostRecentItem(
      absl::string_view event_source) const;

  // Set the logging configuration for an event_source
  absl::Status SetLogOptions(
      const std::map<std::string, intrinsic_proto::data_logger::LogOptions>&
          options) const;

  // Get the logging configuration for an event_source
  absl::StatusOr<LogOptions> GetLogOptions(
      absl::string_view event_source) const;

  // Writes all log files of the specified 'event_sources' to GCS.
  // Might be throttled per-event-source if called too frequently.
  //
  // Returns absl::ResourceExhaustedError if any sync for any event source was
  // throttled.
  absl::StatusOr<std::vector<std::string>> SyncAndRotateLogs(
      absl::Span<const absl::string_view> event_sources) const;

  // Writes all log files to GCS.
  // Might be throttled per-event-source if called too frequently.
  //
  // Returns absl::ResourceExhaustedError if any sync for any event source was
  // throttled.
  absl::StatusOr<std::vector<std::string>> SyncAndRotateLogs() const;

  // Creates a local recording from structured logging data.
  //
  // The data will be copied to its own table for safekeeping, to ensure
  // that it does not get deleted from the rolling-buffer TimescaleDB tables if
  // the upload process is very slow due to limited internet connectivity.
  absl::StatusOr<intrinsic_proto::data_logger::BagMetadata>
  CreateLocalRecording(
      absl::Time start_time, absl::Time end_time, absl::string_view description,
      absl::Span<const absl::string_view> event_sources_to_record) const;

  // List recordings stored locally.
  absl::StatusOr<std::vector<intrinsic_proto::data_logger::BagMetadata>>
  ListLocalRecordings(std::optional<absl::Time> start_time,
                      std::optional<absl::Time> end_time,
                      bool only_summary_metadata,
                      absl::Span<const absl::string_view> bag_ids) const;

 private:
  // Use of pimpl / firewall idiom to hide gRPC details.
  struct StructuredLoggingClientImpl;
  std::unique_ptr<StructuredLoggingClientImpl> impl_;
};

}  // namespace intrinsic

#endif  // INTRINSIC_LOGGING_STRUCTURED_LOGGING_CLIENT_H_
