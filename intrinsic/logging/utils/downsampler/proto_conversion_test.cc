// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/logging/utils/downsampler/proto_conversion.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <optional>

#include "absl/status/status.h"
#include "absl/time/time.h"
#include "intrinsic/logging/proto/downsampler.pb.h"
#include "intrinsic/logging/utils/downsampler/downsampler.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

namespace intrinsic {
namespace {

using ::intrinsic::data_logger::DownsamplerEventSourceState;
using ::intrinsic::data_logger::DownsamplerOptions;
using ::intrinsic::data_logger::DownsamplerState;

using DownsamplerOptionsProto =
    ::intrinsic_proto::data_logger::DownsamplerOptions;
using DownsamplerEventSourceStateProto =
    ::intrinsic_proto::data_logger::DownsamplerEventSourceState;
using DownsamplerStateProto = ::intrinsic_proto::data_logger::DownsamplerState;

using ::absl_testing::IsOkAndHolds;
using ::absl_testing::StatusIs;
using ::testing::HasSubstr;

// DownsamplerOptions.

TEST(DownsamplerOptions, RoundtripWorks) {
  DownsamplerOptions options{
      .sampling_interval_time = absl::Seconds(1),
      .sampling_interval_count = 10,
  };
  ASSERT_OK_AND_ASSIGN(DownsamplerOptionsProto proto, ToProto(options));
  EXPECT_THAT(FromProto(proto), IsOkAndHolds(options));

  DownsamplerOptions nullopt_options{
      .sampling_interval_time = std::nullopt,
      .sampling_interval_count = std::nullopt,
  };
  ASSERT_OK_AND_ASSIGN(DownsamplerOptionsProto nullopt_options_proto,
                       ToProto(nullopt_options));
  EXPECT_THAT(FromProto(nullopt_options_proto), IsOkAndHolds(nullopt_options));
}

TEST(DownsamplerOptions, ErrorsOnInvalidSamplingIntervalTime) {
  DownsamplerOptions invalid_options{.sampling_interval_time =
                                         absl::InfiniteDuration()};
  EXPECT_THAT(ToProto(invalid_options),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("invalid sampling_interval_time")));

  DownsamplerOptionsProto invalid_options_proto;
  invalid_options_proto.mutable_sampling_interval_time()->set_seconds(
      315576000001);
  EXPECT_THAT(FromProto(invalid_options_proto),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("invalid sampling_interval_time")));
}

// DownsamplerEventSourceState.

TEST(DownsamplerEventSourceState, RoundtripWorks) {
  DownsamplerEventSourceState state{
      .last_use_time = absl::FromUnixSeconds(123),
      .count_since_last_use = 42,
  };
  ASSERT_OK_AND_ASSIGN(DownsamplerEventSourceStateProto state_proto,
                       ToProto(state));
  EXPECT_THAT(FromProto(state_proto), IsOkAndHolds(state));
}

TEST(DownsamplerEventSourceState, ErrorsOnInvalidLastUseTime) {
  DownsamplerEventSourceState invalid_state{.last_use_time =
                                                absl::InfiniteFuture()};
  EXPECT_THAT(ToProto(invalid_state),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("invalid last_use_time")));

  DownsamplerEventSourceStateProto invalid_state_proto;
  invalid_state_proto.mutable_last_use_time()->set_seconds(253402300800);
  EXPECT_THAT(FromProto(invalid_state_proto),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("invalid last_use_time")));
}

// DownsamplerState.

TEST(DownsamplerState, RoundtripWorks) {
  DownsamplerState state{
      .event_source_states = {{"event_source_1",
                               {
                                   .last_use_time = absl::FromUnixSeconds(123),
                                   .count_since_last_use = 42,
                               }},
                              {"event_source_2",
                               {
                                   .last_use_time = absl::FromUnixSeconds(456),
                                   .count_since_last_use = 84,
                               }}},
  };
  ASSERT_OK_AND_ASSIGN(DownsamplerStateProto state_proto, ToProto(state));
  EXPECT_THAT(FromProto(state_proto), IsOkAndHolds(state));

  DownsamplerState empty_state{};
  ASSERT_OK_AND_ASSIGN(DownsamplerStateProto empty_state_proto,
                       ToProto(empty_state));
  EXPECT_THAT(FromProto(empty_state_proto), IsOkAndHolds(empty_state));
}

TEST(DownsamplerState, ErrorsOnInvalidLastUseTime) {
  DownsamplerState invalid_state{
      .event_source_states = {{"event_source_1",
                               {
                                   .last_use_time = absl::InfiniteFuture(),
                               }}},
  };
  EXPECT_THAT(ToProto(invalid_state),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("invalid last_use_time")));

  DownsamplerStateProto invalid_state_proto;
  (*invalid_state_proto.mutable_event_source_states())["event_source_1"]
      .mutable_last_use_time()
      ->set_seconds(253402300800);
  EXPECT_THAT(FromProto(invalid_state_proto),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("invalid last_use_time")));
}

}  // namespace
}  // namespace intrinsic
