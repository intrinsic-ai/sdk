// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_UTILS_LOG_H_
#define INTRINSIC_ICON_UTILS_LOG_H_

#include "absl/base/log_severity.h"                  // IWYU pragma: export
#include "intrinsic/icon/release/source_location.h"  // IWYU pragma: export
#include "intrinsic/icon/utils/log_internal.h"       // IWYU pragma: export
#include "intrinsic/icon/utils/log_sink.h"           // IWYU pragma: export

// A real-time logging interface.
//
// This is a small, real-time safe variant of Google C++ Logging
// http://abseil.io/docs/cpp/guides/logging. It does not allocate and
// truncates messages longer than LogSinkInterface::kLogMessageMaxSize.
// Allowed log levels are INFO, WARNING and ERROR. Instead of FATAL, you could
// log ERROR and then call CHECK from absl/log/check.h. It supports all types
// that absl/strings/str_cat.h can convert, including absl::string_view, which
// works well with icon::FixedString.
//
// Examples:
//   #include "intrinsic/icon/utils/log.h"
//   INTRINSIC_RT_LOG(INFO) << "first joint position: " << joint_position[0];
//   INTRINSIC_RT_LOG(ERROR) << "part: " << part.name();  // string_view
//
//   // Logs at most once every 2 seconds.
//   INTRINSIC_RT_LOG_THROTTLED(WARNING) << "limit exceeded";
//
// FATAL LOGGING
// -------------
// None of these macros is fatal.  For a non-recoverable error use
// LOG(FATAL) or or CHECK_* which will log a message and terminate the program
// (see absl/log/check.h and absl/log/log.h).
//
// THROTTLED LOGGING
// -----------------
// INTRINSIC_RT_LOG_THROTTLED collects repetitions of a message at the
// same call site over a short period of time (2 seconds).
// It also prints a count how many messages were ignored.
// This logging function is useful to avoid log spam for
// high-frequency calls (for example, every millisecond).

namespace intrinsic {

// Not RT safe.
// Must be called before using any of the logging macros below, unless running
// in a intrinsic::Thread.
// Otherwise, INTRINSIC_RT_LOG* is not real-time safe.
void RtLogInitForThisThread();

}  // namespace intrinsic

// NOLINTBEGIN(readability/braces)
#define INTRINSIC_RT_LOG(SEVERITY)                          \
  if (true)                                                 \
  ::intrinsic::icon::internal::LogClient() +=               \
      ::intrinsic::icon::internal::LogEntryBuilder::Create( \
          ::intrinsic::icon::LogPriority::SEVERITY, INTRINSIC_LOC)
// NOLINTEND(readability/braces)

// NOLINTBEGIN(readability/braces)
#define INTRINSIC_RT_LOG_THROTTLED(SEVERITY)                              \
  if (static ::intrinsic::icon::internal::LogThrottler throttler; true)   \
    if (auto result =                                                     \
            throttler.Tick(::intrinsic::icon::GlobalLogContext::GetTime); \
        result.has_value())                                               \
  ::intrinsic::icon::internal::LogClient() +=                             \
      ::intrinsic::icon::internal::LogEntryBuilder::Create(               \
          ::intrinsic::icon::LogPriority::SEVERITY, result.value(),       \
          INTRINSIC_LOC)
// NOLINTEND(readability/braces)

// Logs the message the first N times this macro instance is called.
// NOLINTBEGIN(readability/braces)
#define INTRINSIC_RT_LOG_FIRST_N(SEVERITY, N)                         \
  if (static size_t COUNTER_##__LINE__ = 0; COUNTER_##__LINE__++ < N) \
  ::intrinsic::icon::internal::LogClient() +=                         \
      ::intrinsic::icon::internal::LogEntryBuilder::Create(           \
          ::intrinsic::icon::LogPriority::SEVERITY, INTRINSIC_LOC)
// NOLINTEND(readability/braces)

// Logs the first time it is called.
#define INTRINSIC_RT_LOG_FIRST(SEVERITY) INTRINSIC_RT_LOG_FIRST_N(SEVERITY, 1)

// Logs if the condition is met.
#define INTRINSIC_RT_LOG_IF(SEVERITY, CONDITION) \
  if (CONDITION) INTRINSIC_RT_LOG(SEVERITY)

// Documentation for developers of logging:
// Filtering is implemented similar to absl/log/internal/conditions.h
// Also, the if clause will error if prefixes (like intrinsic::) are used,
// which we don't want call sites to rely on.

#endif  // INTRINSIC_ICON_UTILS_LOG_H_
