// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/proto/parsed_type_url.h"

#include <ostream>
#include <string>
#include <string_view>
#include <utility>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_format.h"
#include "absl/strings/str_split.h"
#include "absl/strings/string_view.h"
#include "intrinsic/util/proto/status_specs.h"
#include "intrinsic/util/proto/type_url.h"
#include "intrinsic/util/status/status_builder.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {

namespace {

// Parses a type_url or type_url_prefix up until and including the area part.
// On success returns the partially filled ParsedUrl and the remainder of the
// type_url starting after the '/' from the 'area' part.
absl::StatusOr<std::pair<ParsedUrl, std::string_view>> ParseTypeUrlToArea(
    std::string_view type_url) {
  ParsedUrl parsed_url = {.type_url = std::string(type_url)};

  if (!type_url.starts_with(kIntrinsicTypeUrlPrefix)) {
    std::string message =
        absl::StrFormat("Type URL '%s' does not start with '%s'", type_url,
                        kIntrinsicTypeUrlPrefix);
    return (StatusBuilder(absl::StatusCode::kInvalidArgument) << message)
        .AttachExtendedStatus(
            util::proto::kExtendedStatusComponent, util::proto::kInvalidUrlCode,
            {.title = util::proto::kInvalidUrlTitle,
             .user_message = message,
             .user_instructions = util::proto::kInvalidUrlInstructions});
  }

  parsed_url.prefix = type_url.substr(0, kIntrinsicTypeUrlPrefix.length());
  std::string_view type_url_parsed = type_url;
  type_url_parsed.remove_prefix(kIntrinsicTypeUrlPrefix.length());

  std::pair<std::string_view, std::string_view> area_and_remainder =
      absl::StrSplit(type_url_parsed, absl::MaxSplits(kTypeUrlSeparator, 1));

  // If no split has happened, the entire input ends up in the
  // first element of the pair, and the second is empty.
  if (area_and_remainder.first.length() == type_url_parsed.length() ||
      area_and_remainder.first.empty()) {
    std::string message = absl::StrFormat(
        "Type URL '%s' is missing area after Intrinsic prefix", type_url);
    return (StatusBuilder(absl::StatusCode::kInvalidArgument) << message)
        .AttachExtendedStatus(
            util::proto::kExtendedStatusComponent, util::proto::kInvalidUrlCode,
            {.title = util::proto::kInvalidUrlTitle,
             .user_message = message,
             .user_instructions = util::proto::kInvalidUrlInstructions});
  }

  parsed_url.area = area_and_remainder.first;

  return std::make_pair(std::move(parsed_url), area_and_remainder.second);
}

}  // namespace

absl::StatusOr<ParsedUrl> ParseTypeUrl(std::string_view type_url) {
  INTR_ASSIGN_OR_RETURN(auto parsed_url_and_remainder,
                        ParseTypeUrlToArea(type_url));
  ParsedUrl parsed_url = std::move(parsed_url_and_remainder.first);
  std::string_view remainder_after_area = parsed_url_and_remainder.second;

  std::string_view::size_type last_slash_pos =
      remainder_after_area.rfind(kTypeUrlSeparator);
  if (last_slash_pos == std::string_view::npos) {
    std::string message = absl::StrFormat(
        "Type URL '%s' is missing separator after area", type_url);
    return (StatusBuilder(absl::StatusCode::kInvalidArgument) << message)
        .AttachExtendedStatus(
            util::proto::kExtendedStatusComponent, util::proto::kInvalidUrlCode,
            {.title = util::proto::kInvalidUrlTitle,
             .user_message = message,
             .user_instructions = util::proto::kInvalidUrlInstructions});
  }

  parsed_url.path = remainder_after_area.substr(0, last_slash_pos);
  parsed_url.message_type = remainder_after_area.substr(last_slash_pos + 1);

  if (parsed_url.path.empty() || parsed_url.message_type.empty()) {
    std::string message = absl::StrFormat(
        "Type URL '%s' is missing path or message type", type_url);
    return (StatusBuilder(absl::StatusCode::kInvalidArgument) << message)
        .AttachExtendedStatus(
            util::proto::kExtendedStatusComponent, util::proto::kInvalidUrlCode,
            {.title = util::proto::kInvalidUrlTitle,
             .user_message = message,
             .user_instructions = util::proto::kInvalidUrlInstructions});
  }

  return parsed_url;
}

absl::StatusOr<ParsedUrl> ParseTypeUrlPrefix(std::string_view type_url_prefix) {
  INTR_ASSIGN_OR_RETURN(auto parsed_url_and_remainder,
                        ParseTypeUrlToArea(type_url_prefix));
  ParsedUrl parsed_url = std::move(parsed_url_and_remainder.first);
  std::string_view remainder_after_area = parsed_url_and_remainder.second;

  if (remainder_after_area.ends_with(kTypeUrlSeparator)) {
    remainder_after_area.remove_suffix(kTypeUrlSeparator.length());
  }
  if (remainder_after_area.empty()) {
    return absl::InvalidArgumentError(absl::StrFormat(
        "Type URL prefix '%s' is missing path", type_url_prefix));
  }

  parsed_url.path = remainder_after_area;

  return parsed_url;
}

std::ostream& operator<<(std::ostream& os, const ParsedUrl& parsed_url) {
  os << "ParsedUrl{type_url: " << parsed_url.type_url
     << ", prefix: " << parsed_url.prefix << ", area: " << parsed_url.area
     << ", path: " << parsed_url.path
     << ", message_type: " << parsed_url.message_type << "}";
  return os;
}

}  // namespace intrinsic
