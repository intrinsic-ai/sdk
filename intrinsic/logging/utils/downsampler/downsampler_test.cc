// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/logging/utils/downsampler/downsampler.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "absl/time/time.h"
#include "intrinsic/logging/proto/downsampler.pb.h"
#include "intrinsic/logging/proto/log_item.pb.h"
#include "intrinsic/logging/proto/logger_service.pb.h"
#include "intrinsic/util/proto_time.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

namespace intrinsic::data_logger {
namespace {

using DownsamplerOptionsProto =
    ::intrinsic_proto::data_logger::DownsamplerOptions;
using DownsamplerEventSourceStateProto =
    ::intrinsic_proto::data_logger::DownsamplerEventSourceState;
using DownsamplerStateProto = ::intrinsic_proto::data_logger::DownsamplerState;

using ::intrinsic_proto::data_logger::LogItem;

using ::absl_testing::IsOkAndHolds;
using ::absl_testing::StatusIs;
using ::testing::Eq;
using ::testing::HasSubstr;
using ::testing::Pair;
using ::testing::UnorderedElementsAre;

// Test helpers.

LogItem CreateLogItem(absl::string_view event_source,
                      absl::Time acquisition_time) {
  LogItem item;
  item.mutable_metadata()->set_event_source(event_source);
  *item.mutable_metadata()->mutable_acquisition_time() =
      *FromAbslTime(acquisition_time);
  return item;
}

// Type equality.

TEST(TypeEqualityTest, DownsamplerOptions) {
  DownsamplerOptions options_1 = {.sampling_interval_time = absl::Seconds(1),
                                  .sampling_interval_count = 10};
  DownsamplerOptions options_2 = {.sampling_interval_time = absl::Seconds(1),
                                  .sampling_interval_count = 10};
  EXPECT_EQ(options_1, options_2);
}

TEST(TypeEqualityTest, DownsamplerEventSourceState) {
  DownsamplerEventSourceState state_1 = {
      .last_use_time = absl::FromUnixSeconds(100), .count_since_last_use = 1};
  DownsamplerEventSourceState state_2 = {
      .last_use_time = absl::FromUnixSeconds(100), .count_since_last_use = 1};
  EXPECT_EQ(state_1, state_2);
}

TEST(TypeEqualityTest, DownsamplerState) {
  DownsamplerState state_1 = {
      .event_source_states = {{"test_source",
                               {.last_use_time = absl::FromUnixSeconds(100),
                                .count_since_last_use = 1}}}};
  DownsamplerState state_2 = {
      .event_source_states = {{"test_source",
                               {.last_use_time = absl::FromUnixSeconds(100),
                                .count_since_last_use = 1}}}};
  EXPECT_EQ(state_1, state_2);
}

// Downsampler.

TEST(DownsamplerTest, FirstItemShouldNotDownsample) {
  Downsampler downsampler({.sampling_interval_count = 2});
  LogItem item = CreateLogItem("test_source", absl::FromUnixSeconds(0));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));
  EXPECT_OK(downsampler.RegisterIngest(item));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(true));
}

TEST(DownsamplerTest, NoOptionsShouldNotDownsample) {
  Downsampler downsampler({});
  LogItem item = CreateLogItem("test_source", absl::FromUnixSeconds(0));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));
  EXPECT_OK(downsampler.RegisterIngest(item));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));
}

TEST(DownsamplerTest, TimeIntervalDownsamplingZero) {
  Downsampler downsampler({.sampling_interval_time = absl::ZeroDuration()});
  LogItem item = CreateLogItem("test_source", absl::FromUnixSeconds(0));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));
  EXPECT_OK(downsampler.RegisterIngest(item));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));
}

TEST(DownsamplerTest, TimeIntervalDownsamplingWorks) {
  Downsampler downsampler({.sampling_interval_time = absl::Seconds(2)});
  LogItem item = CreateLogItem("test_source", absl::FromUnixSeconds(0));

  // First message should go through.
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));

  // It should continue to be passed until a use is registered.
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));
  EXPECT_OK(downsampler.RegisterIngest(item));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(true));

  // Within interval.
  item = CreateLogItem("test_source", absl::FromUnixSeconds(1));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(true));

  item = CreateLogItem("test_source", absl::FromUnixSeconds(2));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));
  EXPECT_OK(downsampler.RegisterIngest(item));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(true));

  // Outside interval.
  item = CreateLogItem("test_source", absl::FromUnixSeconds(10));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));

  // Again, it should continue to be passed until a use is registered.
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));
  EXPECT_OK(downsampler.RegisterIngest(item));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(true));
}

TEST(DownsamplerTest, CountIntervalDownsamplingZero) {
  Downsampler downsampler({.sampling_interval_count = 0});
  LogItem item = CreateLogItem("test_source", absl::FromUnixSeconds(0));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));
  EXPECT_OK(downsampler.RegisterIngest(item));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));
}

TEST(DownsamplerTest, CountIntervalDownsamplingWorks) {
  Downsampler downsampler({.sampling_interval_count = 3});
  LogItem item = CreateLogItem("test_source", absl::FromUnixSeconds(0));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));  // 1
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));  // 2
  EXPECT_OK(downsampler.RegisterIngest(item));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(true));   // 1
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(true));   // 2
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));  // 3
}

TEST(DownsamplerTest, CombinedDownsampling) {
  Downsampler downsampler({.sampling_interval_time = absl::Seconds(2),
                           .sampling_interval_count = 5});
  LogItem item = CreateLogItem("test_source", absl::FromUnixSeconds(0));

  // First message should go through.
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));

  // It should continue to be passed until a use is registered.
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));
  EXPECT_OK(downsampler.RegisterIngest(item));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(true));  // 1

  // Within time interval, but out of count interval.
  item = CreateLogItem("test_source", absl::FromUnixSeconds(1));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(true));  // 2
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(true));  // 3
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(true));  // 4
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(true));  // 5
  EXPECT_OK(downsampler.RegisterIngest(item));

  // Out of time interval, but not count interval.
  item = CreateLogItem("test_source", absl::FromUnixSeconds(5));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(true));   // 1
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(true));   // 2
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(true));   // 3
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(true));   // 4
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));  // 5

  // Again, it should continue to be passed until a use is registered.
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));
  EXPECT_OK(downsampler.RegisterIngest(item));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(true));
}

TEST(DownsamplerTest, TracksDifferentEventSourcesSeparately) {
  Downsampler downsampler({.sampling_interval_count = 3});

  LogItem a = CreateLogItem("test_source_a", absl::FromUnixSeconds(0));
  LogItem b = CreateLogItem("test_source_b", absl::FromUnixSeconds(0));

  EXPECT_THAT(downsampler.ShouldDownsample(a), IsOkAndHolds(false));  // A: 1
  EXPECT_OK(downsampler.RegisterIngest(a));                           // A: 0
  EXPECT_THAT(downsampler.ShouldDownsample(a), IsOkAndHolds(true));   // A: 1

  EXPECT_THAT(downsampler.ShouldDownsample(b), IsOkAndHolds(false));  // B: 1
  EXPECT_OK(downsampler.RegisterIngest(b));                           // B: 0
  EXPECT_THAT(downsampler.ShouldDownsample(b), IsOkAndHolds(true));   // B: 1

  EXPECT_THAT(downsampler.ShouldDownsample(a), IsOkAndHolds(true));   // A: 2
  EXPECT_THAT(downsampler.ShouldDownsample(a), IsOkAndHolds(false));  // A: 3

  EXPECT_THAT(downsampler.ShouldDownsample(b), IsOkAndHolds(true));   // B: 2
  EXPECT_THAT(downsampler.ShouldDownsample(b), IsOkAndHolds(false));  // B: 3

  EXPECT_OK(downsampler.RegisterIngest(a));                           // A: 0
  EXPECT_THAT(downsampler.ShouldDownsample(a), IsOkAndHolds(true));   // A: 1
  EXPECT_THAT(downsampler.ShouldDownsample(b), IsOkAndHolds(false));  // B: 4
}

TEST(DownsamplerTest, GetAndSetEventSourceStateWorks) {
  Downsampler downsampler({.sampling_interval_count = 3});

  LogItem item = CreateLogItem("test_source", absl::FromUnixSeconds(0));

  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(false));  // 1
  EXPECT_OK(downsampler.RegisterIngest(item));
  EXPECT_THAT(downsampler.ShouldDownsample(item), IsOkAndHolds(true));  // 1

  // Get ok.
  ASSERT_OK_AND_ASSIGN(DownsamplerEventSourceState state,
                       downsampler.GetEventSourceState("test_source"));
  EXPECT_EQ(state.count_since_last_use, 1);
  EXPECT_EQ(state.last_use_time, absl::FromUnixSeconds(0));

  // Get non-existent.
  EXPECT_THAT(
      downsampler.GetEventSourceState("non_existent_source"),
      StatusIs(absl::StatusCode::kNotFound,
               HasSubstr("event source not found: non_existent_source")));

  // Tracks use.
  item = CreateLogItem("test_source", absl::FromUnixSeconds(1));
  EXPECT_OK(downsampler.RegisterIngest(item));
  ASSERT_OK_AND_ASSIGN(state, downsampler.GetEventSourceState("test_source"));
  EXPECT_EQ(state.count_since_last_use, 0);
  EXPECT_EQ(state.last_use_time, absl::FromUnixSeconds(1));

  // Restores across instances.
  Downsampler downsampler_2({.sampling_interval_count = 3});
  EXPECT_OK(downsampler_2.SetEventSourceState("test_source", state));
  EXPECT_THAT(downsampler_2.ShouldDownsample(item), IsOkAndHolds(true));   // 1
  EXPECT_THAT(downsampler_2.ShouldDownsample(item), IsOkAndHolds(true));   // 2
  EXPECT_THAT(downsampler_2.ShouldDownsample(item), IsOkAndHolds(false));  // 3

  // Also works if downsampler noops.
  Downsampler downsampler_3({});
  EXPECT_THAT(downsampler_3.ShouldDownsample(item), IsOkAndHolds(false));  // 1
  EXPECT_THAT(downsampler_3.GetEventSourceState("test_source"),
              StatusIs(absl::StatusCode::kNotFound));
  EXPECT_OK(downsampler_3.SetEventSourceState("test_source", state));
  EXPECT_THAT(downsampler_3.ShouldDownsample(item), IsOkAndHolds(false));  // 1
  ASSERT_OK_AND_ASSIGN(state, downsampler_3.GetEventSourceState("test_source"));
  EXPECT_EQ(state.count_since_last_use, 1);
  EXPECT_EQ(state.last_use_time, absl::FromUnixSeconds(1));
}

TEST(DownsamplerTest, GetAndSetStateWorks) {
  Downsampler downsampler({.sampling_interval_count = 2});
  ASSERT_OK_AND_ASSIGN(DownsamplerState state, downsampler.GetState());
  EXPECT_TRUE(state.event_source_states.empty());

  LogItem a = CreateLogItem("source_a", absl::FromUnixSeconds(10));
  LogItem b = CreateLogItem("source_b", absl::FromUnixSeconds(20));

  // Doesn't track unseen event sources.
  EXPECT_THAT(downsampler.ShouldDownsample(a), IsOkAndHolds(false));
  ASSERT_OK_AND_ASSIGN(state, downsampler.GetState());
  EXPECT_TRUE(state.event_source_states.empty());

  // Tracks registrations.
  EXPECT_OK(downsampler.RegisterIngest(a));
  ASSERT_OK_AND_ASSIGN(state, downsampler.GetState());
  EXPECT_THAT(state.event_source_states,
              UnorderedElementsAre(Pair(
                  "source_a", Eq(DownsamplerEventSourceState{
                                  .last_use_time = absl::FromUnixSeconds(10),
                                  .count_since_last_use = 0}))));

  EXPECT_THAT(downsampler.ShouldDownsample(b), IsOkAndHolds(false));
  EXPECT_OK(downsampler.RegisterIngest(b));
  EXPECT_THAT(downsampler.ShouldDownsample(b), IsOkAndHolds(true));
  ASSERT_OK_AND_ASSIGN(state, downsampler.GetState());
  EXPECT_THAT(
      state.event_source_states,
      UnorderedElementsAre(
          Pair("source_a", Eq(DownsamplerEventSourceState{
                               .last_use_time = absl::FromUnixSeconds(10),
                               .count_since_last_use = 0})),
          Pair("source_b", Eq(DownsamplerEventSourceState{
                               .last_use_time = absl::FromUnixSeconds(20),
                               .count_since_last_use = 1}))));

  // Restores across instances.
  Downsampler downsampler_2({.sampling_interval_count = 2});
  EXPECT_OK(downsampler_2.SetState(state));  // Sets B: 1

  LogItem a_2 = CreateLogItem("source_a", absl::FromUnixSeconds(11));
  LogItem b_2 = CreateLogItem("source_b", absl::FromUnixSeconds(21));

  EXPECT_THAT(downsampler_2.ShouldDownsample(a_2),
              IsOkAndHolds(true));  // A: 1
  EXPECT_THAT(downsampler_2.ShouldDownsample(a_2),
              IsOkAndHolds(false));  // A: 2
  EXPECT_THAT(downsampler_2.ShouldDownsample(b_2),
              IsOkAndHolds(false));  // B: 2

  // Also works if downsampler noops.
  Downsampler downsampler_3({});
  EXPECT_OK(downsampler_3.SetState(state));
  EXPECT_THAT(downsampler_3.ShouldDownsample(a_2),
              IsOkAndHolds(false));  // A: 1
  EXPECT_THAT(downsampler_3.ShouldDownsample(a_2),
              IsOkAndHolds(false));  // A: 2
}

TEST(DownsamplerTest, Reset) {
  Downsampler downsampler({
      .sampling_interval_time = absl::Seconds(2),
  });

  LogItem a = CreateLogItem("test_source_a", absl::FromUnixSeconds(0));
  LogItem b = CreateLogItem("test_source_b", absl::FromUnixSeconds(0));

  EXPECT_THAT(downsampler.ShouldDownsample(a), IsOkAndHolds(false));
  EXPECT_OK(downsampler.RegisterIngest(a));
  EXPECT_THAT(downsampler.ShouldDownsample(a), IsOkAndHolds(true));

  EXPECT_THAT(downsampler.ShouldDownsample(b), IsOkAndHolds(false));
  EXPECT_OK(downsampler.RegisterIngest(b));
  EXPECT_THAT(downsampler.ShouldDownsample(b), IsOkAndHolds(true));

  downsampler.Reset();

  EXPECT_THAT(downsampler.ShouldDownsample(a), IsOkAndHolds(false));
  EXPECT_THAT(downsampler.ShouldDownsample(b), IsOkAndHolds(false));
}

}  // namespace
}  // namespace intrinsic::data_logger
