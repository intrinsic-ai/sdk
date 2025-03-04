// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_PROTO_PARSE_TEXT_PROTO_H_
#define INTRINSIC_UTIL_PROTO_PARSE_TEXT_PROTO_H_

#include <string>
#include <string_view>
#include <type_traits>

#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "google/protobuf/message.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {

// Parses the given text proto into the given message. Returns an error if the
// parsing fails, e.g., if the text proto does not match the type of the given
// message.
absl::Status ParseTextProtoInto(std::string_view asciipb,
                                google::protobuf::Message* message);

// Parses the given text proto as a protocol message of type 'T' and returns the
// result in a StatusOr<T>.
template <typename T, typename = std::enable_if_t<
                          std::is_base_of_v<google::protobuf::Message, T>>>
absl::StatusOr<T> ParseTextProto(absl::string_view asciipb) {
  T message;
  INTR_RETURN_IF_ERROR(ParseTextProtoInto(asciipb, &message));
  return message;
}

namespace internal {

// Internal helper class for ParseTextProtoOrDie() / ParseTextOrDie().
// Do not use directly.
class ParseTextProtoHelper {
 public:
  explicit ParseTextProtoHelper(absl::string_view text) : text_(text) {}

  template <typename T, typename = std::enable_if_t<
                            std::is_base_of_v<google::protobuf::Message, T>>>
  operator T() const {  // NOLINT(google-explicit-constructor)
    auto parse_result = ParseTextProto<T>(text_);
    CHECK_OK(parse_result.status());
    return *parse_result;
  }

 private:
  std::string text_;
  friend ParseTextProtoHelper ParseTextProtoOrDie(absl::string_view);
};

}  // namespace internal

// Parses the given text proto as a protocol message whose type is automatically
// inferred from the return type. If the parsing fails, prints a failure message
// and terminates the program.
inline internal::ParseTextProtoHelper ParseTextProtoOrDie(
    absl::string_view asciipb) {
  return internal::ParseTextProtoHelper(asciipb);
}

// Parses the given text proto as a protocol message of type 'T'. If the parsing
// fails, prints a failure message and terminates the program.
template <typename T, typename = std::enable_if_t<
                          std::is_base_of_v<google::protobuf::Message, T>>>
T ParseTextOrDie(absl::string_view asciipb) {
  return ParseTextProtoOrDie(asciipb);
}

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_PROTO_PARSE_TEXT_PROTO_H_
