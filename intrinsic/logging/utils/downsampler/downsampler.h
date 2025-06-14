// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_LOGGING_UTILS_DOWNSAMPLER_DOWNSAMPLER_H_
#define INTRINSIC_LOGGING_UTILS_DOWNSAMPLER_DOWNSAMPLER_H_

#include <cstdint>
#include <optional>
#include <string>

#include "absl/container/flat_hash_map.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/time/time.h"
#include "intrinsic/logging/proto/downsampler.pb.h"
#include "intrinsic/logging/proto/log_item.pb.h"

namespace intrinsic::data_logger {

// Options for downsampling.
//
// Only applies a sampling interval if it is set.
//
// If multiple sampling intervals are set, the candidate will only be sampled if
// all intervals are met.
struct DownsamplerOptions {
  std::optional<absl::Duration> sampling_interval_time = std::nullopt;
  std::optional<int32_t> sampling_interval_count = std::nullopt;

  bool operator==(const DownsamplerOptions& other) const;
};

// Tracks the state of an event source in the downsampler.
struct DownsamplerEventSourceState {
  absl::Time last_use_time = absl::InfinitePast();
  int32_t count_since_last_use = 0;

  bool operator==(const DownsamplerEventSourceState& other) const;
};

// Tracks the state of the downsampler.
struct DownsamplerState {
  absl::flat_hash_map<std::string, DownsamplerEventSourceState>
      event_source_states;

  bool operator==(const DownsamplerState& other) const;
};

// Determines whether a LogItem should be downsampled.
// Downsamples per event source.
//
// This class is conditionally thread-compatible, please ensure external
// synchronization if used in a multi-threaded environment, and ensure that
// calls to ShouldDownsample and RegisterUse are made in a strictly time-ordered
// manner with no duplicate calls per LogItem.
class Downsampler {
 public:
  typedef DownsamplerOptions Options;

  explicit Downsampler(DownsamplerOptions options);

  // Returns true if the LogItem should be downsampled.
  //
  // Each call will increment the counter for the appropriate event source for
  // count-based sampling. Do not call on the same LogItem multiple times unless
  // the intention is to treat it as distinct duplicate copies.
  absl::StatusOr<bool> ShouldDownsample(
      const intrinsic_proto::data_logger::LogItem& item);

  // Gets the state of the downsampler for the given event source.
  // If the event source is not registered, returns an error.
  absl::StatusOr<DownsamplerEventSourceState> GetEventSourceState(
      absl::string_view event_source) const;

  // Sets the state of the downsampler for the given event_source.
  // If the event source is not registered, it will be registered.
  //
  // This is useful for restoring the downsampler state (e.g., using an external
  // pagination cursor).
  absl::Status SetEventSourceState(absl::string_view event_source,
                                   const DownsamplerEventSourceState& state);

  // Gets the entire state of the downsampler.
  absl::StatusOr<DownsamplerState> GetState() const;

  // Sets the entire state of the downsampler.
  // This completely replaces the existing state.
  absl::Status SetState(const DownsamplerState& state);

  // Registers the ingestion by the caller of a non-downsampled LogItem.
  //
  // This should be called if a LogItem is ingested after the downsampler
  // reports that it should not be downsampled, which factors into the
  // downsampling logic for future items under consideration.
  //
  // If a LogItem is not ingested, then it is treated as if it was downsampled.
  //
  // Registration does not happen automatically because even if a Downsampler
  // reports that a LogItem should be used, the caller might not have used it.
  absl::Status RegisterIngest(
      const intrinsic_proto::data_logger::LogItem& item);

  // Resets the downsampler to its initial state.
  void Reset();

 private:
  DownsamplerOptions options_;
  DownsamplerState state_;
};

}  // namespace intrinsic::data_logger

#endif  // INTRINSIC_LOGGING_UTILS_DOWNSAMPLER_DOWNSAMPLER_H_
