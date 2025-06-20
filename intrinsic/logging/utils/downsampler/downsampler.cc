// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/logging/utils/downsampler/downsampler.h"

#include <optional>
#include <utility>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "absl/time/time.h"
#include "intrinsic/logging/proto/downsampler.pb.h"
#include "intrinsic/logging/proto/log_item.pb.h"
#include "intrinsic/util/proto_time.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::data_logger {

using ::intrinsic_proto::data_logger::LogItem;

using DownsamplerOptionsProto =
    ::intrinsic_proto::data_logger::DownsamplerOptions;
using DownsamplerEventSourceStateProto =
    ::intrinsic_proto::data_logger::DownsamplerEventSourceState;
using DownsamplerStateProto = ::intrinsic_proto::data_logger::DownsamplerState;

// Type equality.

bool DownsamplerOptions::operator==(const DownsamplerOptions& other) const {
  return sampling_interval_time == other.sampling_interval_time &&
         sampling_interval_count == other.sampling_interval_count;
}

bool DownsamplerEventSourceState::operator==(
    const DownsamplerEventSourceState& other) const {
  return last_use_time == other.last_use_time &&
         count_since_last_use == other.count_since_last_use;
}

bool DownsamplerState::operator==(const DownsamplerState& other) const {
  return event_source_states == other.event_source_states;
}

// Downsampler.

Downsampler::Downsampler(DownsamplerOptions options)
    : options_(std::move(options)) {}

absl::StatusOr<bool> Downsampler::ShouldDownsample(const LogItem& item) {
  auto it = state_.event_source_states.find(item.metadata().event_source());

  // Never seen before.
  if (it == state_.event_source_states.end()) {
    return false;
  }

  auto& seen_event_source = it->second;
  seen_event_source.count_since_last_use++;

  // Time-based.
  if (options_.sampling_interval_time.has_value()) {
    INTR_ASSIGN_OR_RETURN(absl::Time acquisition_time,
                          ToAbslTime(item.metadata().acquisition_time()));
    if (acquisition_time - seen_event_source.last_use_time <
        options_.sampling_interval_time) {
      return true;
    }
  }

  // Count-based.
  if (options_.sampling_interval_count.has_value()) {
    if (seen_event_source.count_since_last_use <
        options_.sampling_interval_count) {
      return true;
    }
  }

  return false;
}

absl::StatusOr<DownsamplerEventSourceState> Downsampler::GetEventSourceState(
    absl::string_view event_source) const {
  auto it = state_.event_source_states.find(event_source);
  if (it == state_.event_source_states.end()) {
    return absl::NotFoundError(
        absl::StrCat("event source not found: ", event_source));
  }
  return DownsamplerEventSourceState{
      .last_use_time = it->second.last_use_time,
      .count_since_last_use = it->second.count_since_last_use};
}

absl::Status Downsampler::SetEventSourceState(
    absl::string_view event_source, const DownsamplerEventSourceState& state) {
  state_.event_source_states.insert_or_assign(event_source, state);
  return absl::OkStatus();
}

absl::StatusOr<DownsamplerState> Downsampler::GetState() const {
  DownsamplerState state;
  for (const auto& [event_source, tracker] : state_.event_source_states) {
    state.event_source_states[event_source] = DownsamplerEventSourceState{
        .last_use_time = tracker.last_use_time,
        .count_since_last_use = tracker.count_since_last_use};
  }
  return state;
}

absl::Status Downsampler::SetState(const DownsamplerState& state) {
  state_.event_source_states.clear();
  for (const auto& [event_source, event_state] : state.event_source_states) {
    state_.event_source_states[event_source] = event_state;
  }
  return absl::OkStatus();
}

absl::Status Downsampler::RegisterIngest(const LogItem& item) {
  INTR_ASSIGN_OR_RETURN(absl::Time acquisition_time,
                        ToAbslTime(item.metadata().acquisition_time()));
  state_.event_source_states.insert_or_assign(
      item.metadata().event_source(),
      DownsamplerEventSourceState{.last_use_time = acquisition_time,
                                  .count_since_last_use = 0});
  return absl::OkStatus();
}

void Downsampler::Reset() { state_.event_source_states.clear(); }

}  // namespace intrinsic::data_logger
