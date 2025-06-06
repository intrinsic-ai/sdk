// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_UTILS_REALTIME_STATUS_MATCHERS_H_
#define INTRINSIC_ICON_UTILS_REALTIME_STATUS_MATCHERS_H_

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <ostream>
#include <string>
#include <type_traits>
#include <utility>

#include "absl/status/status.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_macro.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

// Status matchers for RealtimeStatus and RealtimeStatusOr.
//
// Example use:
//   EXPECT_THAT(status, RealtimeIsOk());
//   EXPECT_THAT(status, RealtimeStatusIs(absl::StatusCode::kAborted));
//   EXPECT_THAT(status, RealtimeStatusIs(absl::StatusCode::kNotFound,
//                                        HasSubstr("device")));
//   EXPECT_THAT(result_or, RealtimeIsOkAndHolds(ElementsAre(1, 2, 3)));
//
// Compared to Googletest's matchers absl_testing::IsOk,
// absl_testing::IsOkAndHold, absl_testing::StatusIs, these here are more
// limited.
// We cannot use the native Googletest matchers because they expect
// absl::Status with a std::string, while RealtimeStatus has a stack-allocated
// string.

namespace intrinsic::icon {
namespace internal {

inline const RealtimeStatus& GetRealtimeStatus(const RealtimeStatus& status) {
  return status;
}

template <typename T>
inline const RealtimeStatus& GetRealtimeStatus(
    const RealtimeStatusOr<T>& status_or) {
  return status_or.status();
}

template <typename T>
class RealtimeIsOkMatcherImpl : public ::testing::MatcherInterface<T> {
 public:
  void DescribeTo(std::ostream* os) const override { *os << "is OK"; }

  void DescribeNegationTo(std::ostream* os) const override {
    *os << "is not OK";
  }

  bool MatchAndExplain(
      T status, ::testing::MatchResultListener* listener) const override {
    return GetRealtimeStatus(status).ok();
  }
};

class RealtimeIsOkMatcher {
 public:
  template <typename T>
  operator ::testing::Matcher<T>() const {  // NOLINT
    return ::testing::Matcher<T>(new RealtimeIsOkMatcherImpl<T>());
  }
};

template <typename T>
class RealtimeStatusIsMatcherImpl : public ::testing::MatcherInterface<T> {
 public:
  RealtimeStatusIsMatcherImpl(
      ::testing::Matcher<absl::StatusCode> code_matcher,
      ::testing::Matcher<const std::string&> message_matcher)
      : code_matcher_(std::move(code_matcher)),
        message_matcher_(std::move(message_matcher)) {}

  void DescribeTo(std::ostream* os) const override {
    *os << "has a status code that ";
    code_matcher_.DescribeTo(os);
    *os << ", and has an error message that ";
    message_matcher_.DescribeTo(os);
  }

  void DescribeNegationTo(std::ostream* os) const override {
    *os << "has a status code that ";
    code_matcher_.DescribeNegationTo(os);
    *os << ", or has an error message that ";
    message_matcher_.DescribeNegationTo(os);
  }

  bool MatchAndExplain(
      T status,
      ::testing::MatchResultListener* result_listener) const override {
    ::testing::StringMatchResultListener inner_listener;
    if (!code_matcher_.MatchAndExplain(GetRealtimeStatus(status).code(),
                                       &inner_listener)) {
      *result_listener << (inner_listener.str().empty()
                               ? "whose status code is wrong"
                               : "which has a status code " +
                                     inner_listener.str());
      return false;
    }

    if (!message_matcher_.Matches(
            std::string(GetRealtimeStatus(status).message()))) {
      *result_listener << "whose error message is wrong";
      return false;
    }

    return true;
  }

 private:
  const ::testing::Matcher<absl::StatusCode> code_matcher_;
  const ::testing::Matcher<const std::string&> message_matcher_;
};

class RealtimeStatusIsMatcher {
 public:
  template <typename StatusCodeMatcher, typename StatusMessageMatcher>
  RealtimeStatusIsMatcher(StatusCodeMatcher&& code_matcher,
                          StatusMessageMatcher&& message_matcher)
      : code_matcher_(::testing::MatcherCast<absl::StatusCode>(
            std::forward<StatusCodeMatcher>(code_matcher))),
        message_matcher_(::testing::MatcherCast<const std::string&>(
            std::forward<StatusMessageMatcher>(message_matcher))) {}

  template <typename T>
  operator ::testing::Matcher<T>() const {  // NOLINT
    return ::testing::Matcher<T>(
        new RealtimeStatusIsMatcherImpl<T>(code_matcher_, message_matcher_));
  }

 private:
  const ::testing::Matcher<absl::StatusCode> code_matcher_;
  const ::testing::Matcher<const std::string&> message_matcher_;
};

template <typename T>
class RealtimeIsOkAndHoldsMatcherImpl
    : public ::testing::MatcherInterface<const RealtimeStatusOr<T>&> {
 public:
  using inner_type = const std::decay_t<T>&;

  template <typename InnerMatcher>
  explicit RealtimeIsOkAndHoldsMatcherImpl(InnerMatcher&& inner_matcher)
      : inner_matcher_(::testing::SafeMatcherCast<inner_type>(
            std::forward<InnerMatcher>(inner_matcher))) {}

  void DescribeTo(::std::ostream* os) const override {
    *os << "is OK and has a value that ";
    inner_matcher_.DescribeTo(os);
  }

  void DescribeNegationTo(::std::ostream* os) const override {
    *os << "is not OK or has a value that ";
    inner_matcher_.DescribeNegationTo(os);
  }

  bool MatchAndExplain(
      const RealtimeStatusOr<T>& actual_value,
      ::testing::MatchResultListener* result_listener) const override {
    if (!actual_value.ok()) {
      *result_listener << "which has status "
                       << actual_value.status().message();
      return false;
    }
    ::testing::StringMatchResultListener inner_listener;
    const bool matches =
        inner_matcher_.MatchAndExplain(actual_value.value(), &inner_listener);
    const std::string inner_explanation = inner_listener.str();
    if (!inner_explanation.empty()) {
      *result_listener << "which contains value "
                       << ::testing::PrintToString(actual_value.value()) << ", "
                       << inner_explanation;
    }
    return matches;
  }

 private:
  const ::testing::Matcher<inner_type> inner_matcher_;
};

template <typename InnerMatcher>
class RealtimeIsOkAndHoldsMatcher {
 public:
  explicit RealtimeIsOkAndHoldsMatcher(InnerMatcher inner_matcher)
      : inner_matcher_(std::move(inner_matcher)) {}

  template <typename T>
  operator ::testing::Matcher<const RealtimeStatusOr<T>&>() const {  // NOLINT
    return ::testing::Matcher<const RealtimeStatusOr<T>&>(
        new RealtimeIsOkAndHoldsMatcherImpl<T>(inner_matcher_));
  }

 private:
  const InnerMatcher inner_matcher_;
};

}  // namespace internal

constexpr int kNumExpectAssertMalloc = 2;

#if INTRINSIC_MALLOC_TEST
#define INTRINSIC_RT_ASSERT_OK(expression)                    \
  ASSERT_THAT(expression, ::intrinsic::icon::RealtimeIsOk()); \
  ::intrinsic::icon::MallocCounterSubtract(                   \
      ::intrinsic::icon::kNumExpectAssertMalloc);
#else
#define INTRINSIC_RT_ASSERT_OK(expression) \
  ASSERT_THAT(expression, ::intrinsic::icon::RealtimeIsOk());
#endif

#if INTRINSIC_MALLOC_TEST
#define INTRINSIC_RT_EXPECT_OK(expression)                    \
  EXPECT_THAT(expression, ::intrinsic::icon::RealtimeIsOk()); \
  ::intrinsic::icon::MallocCounterSubtract(                   \
      ::intrinsic::icon::kNumExpectAssertMalloc);
#else
#define INTRINSIC_RT_EXPECT_OK(expression) \
  EXPECT_THAT(expression, ::intrinsic::icon::RealtimeIsOk());
#endif

#define INTRINSIC_RT_ASSERT_OK_AND_ASSIGN_3_(statusor, lhs, expr)         \
  auto statusor = (expr);                                                 \
  INTRINSIC_RT_ASSERT_OK(statusor.status());                              \
  if (!statusor.ok()) {                                                   \
    return;                                                               \
  }                                                                       \
  INTRINSIC_RT_STATUS_MACROS_IMPL_UNPARENTHESIZE_IF_PARENTHESIZED_(lhs) = \
      std::move(statusor.value())

#define INTRINSIC_RT_ASSERT_OK_AND_ASSIGN(lhs, expr) \
  INTRINSIC_RT_ASSERT_OK_AND_ASSIGN_3_(              \
      INTRINSIC_STATUS_MACROS_CONCAT(_status_or_value, __LINE__), lhs, expr)

inline internal::RealtimeIsOkMatcher RealtimeIsOk() {
  return internal::RealtimeIsOkMatcher();
}

template <typename StatusCodeMatcher, typename StatusMessageMatcher>
internal::RealtimeStatusIsMatcher RealtimeStatusIs(
    StatusCodeMatcher&& code_matcher, StatusMessageMatcher&& message_matcher) {
  return internal::RealtimeStatusIsMatcher(
      std::forward<StatusCodeMatcher>(code_matcher),
      std::forward<StatusMessageMatcher>(message_matcher));
}

template <typename StatusCodeMatcher>
internal::RealtimeStatusIsMatcher RealtimeStatusIs(
    StatusCodeMatcher&& code_matcher) {
  return internal::RealtimeStatusIsMatcher(
      std::forward<StatusCodeMatcher>(code_matcher), ::testing::_);
}

template <typename InnerMatcher>
internal::RealtimeIsOkAndHoldsMatcher<std::decay_t<InnerMatcher>>
RealtimeIsOkAndHolds(InnerMatcher&& inner_matcher) {
  return internal::RealtimeIsOkAndHoldsMatcher<std::decay_t<InnerMatcher>>(
      std::forward<InnerMatcher>(inner_matcher));
}

inline void PrintTo(const RealtimeStatus& rtstatus, std::ostream* os) {
  *os << "RealtimeStatus(" << RealtimeStatusCodeToCharArray(rtstatus.code())
      << ", \"" << rtstatus.message() << "\")";
}

template <class T>
inline void PrintTo(const RealtimeStatusOr<T>& rtstatus_or, std::ostream* os) {
  if (rtstatus_or.ok()) {
    *os << "RealtimeStatusOr(value="
        << ::testing::PrintToString(rtstatus_or.value()) << ")";
  } else {
    *os << "RealtimeStatusOr("
        << RealtimeStatusCodeToCharArray(rtstatus_or.status().code()) << ", \""
        << rtstatus_or.status().message() << "\")";
  }
}

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_UTILS_REALTIME_STATUS_MATCHERS_H_
