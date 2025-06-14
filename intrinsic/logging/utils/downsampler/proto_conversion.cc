// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/logging/utils/downsampler/proto_conversion.h"

#include <cstdint>
#include <optional>

#include "absl/status/statusor.h"
#include "absl/time/time.h"
#include "intrinsic/logging/proto/downsampler.pb.h"
#include "intrinsic/logging/utils/downsampler/downsampler.h"
#include "intrinsic/util/proto_time.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic_proto::data_logger {

absl::StatusOr<intrinsic::data_logger::DownsamplerOptions> FromProto(
    const intrinsic_proto::data_logger::DownsamplerOptions& proto) {
  std::optional<absl::Duration> sampling_interval_time = std::nullopt;
  if (proto.has_sampling_interval_time()) {
    INTR_ASSIGN_OR_RETURN(
        sampling_interval_time,
        intrinsic::ToAbslDuration(proto.sampling_interval_time()),
        _ << "invalid sampling_interval_time");
  }

  std::optional<int32_t> sampling_interval_count = std::nullopt;
  if (proto.has_sampling_interval_count()) {
    sampling_interval_count = proto.sampling_interval_count();
  }

  return intrinsic::data_logger::DownsamplerOptions{
      .sampling_interval_time = sampling_interval_time,
      .sampling_interval_count = sampling_interval_count,
  };
}

absl::StatusOr<intrinsic::data_logger::DownsamplerEventSourceState> FromProto(
    const intrinsic_proto::data_logger::DownsamplerEventSourceState& proto) {
  INTR_ASSIGN_OR_RETURN(absl::Time last_use_time,
                        intrinsic::ToAbslTime(proto.last_use_time()),
                        _ << "invalid last_use_time");
  return intrinsic::data_logger::DownsamplerEventSourceState{
      .last_use_time = last_use_time,
      .count_since_last_use = proto.count_since_last_use()};
}

absl::StatusOr<intrinsic::data_logger::DownsamplerState> FromProto(
    const intrinsic_proto::data_logger::DownsamplerState& proto) {
  intrinsic::data_logger::DownsamplerState state;
  for (const auto& [event_source, event_state_proto] :
       proto.event_source_states()) {
    INTR_ASSIGN_OR_RETURN(state.event_source_states[event_source],
                          FromProto(event_state_proto));
  }
  return state;
}

}  // namespace intrinsic_proto::data_logger

namespace intrinsic::data_logger {

absl::StatusOr<intrinsic_proto::data_logger::DownsamplerOptions> ToProto(
    const intrinsic::data_logger::DownsamplerOptions& options) {
  intrinsic_proto::data_logger::DownsamplerOptions proto;
  if (options.sampling_interval_time.has_value()) {
    INTR_ASSIGN_OR_RETURN(*proto.mutable_sampling_interval_time(),
                          FromAbslDuration(*options.sampling_interval_time),
                          _ << "invalid sampling_interval_time");
  }
  if (options.sampling_interval_count.has_value()) {
    proto.set_sampling_interval_count(*options.sampling_interval_count);
  }
  return proto;
}

absl::StatusOr<intrinsic_proto::data_logger::DownsamplerEventSourceState>
ToProto(const intrinsic::data_logger::DownsamplerEventSourceState& state) {
  intrinsic_proto::data_logger::DownsamplerEventSourceState proto;
  INTR_ASSIGN_OR_RETURN(*proto.mutable_last_use_time(),
                        FromAbslTime(state.last_use_time),
                        _ << "invalid last_use_time");
  proto.set_count_since_last_use(state.count_since_last_use);
  return proto;
}

absl::StatusOr<intrinsic_proto::data_logger::DownsamplerState> ToProto(
    const intrinsic::data_logger::DownsamplerState& state) {
  intrinsic_proto::data_logger::DownsamplerState proto;
  for (const auto& [event_source, event_state] : state.event_source_states) {
    INTR_ASSIGN_OR_RETURN((*proto.mutable_event_source_states())[event_source],
                          ToProto(event_state));
  }
  return proto;
}

}  // namespace intrinsic::data_logger
