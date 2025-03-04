// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/proto/parsed_type_url.h"

#include <ostream>
#include <string_view>

#include "absl/base/attributes.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_format.h"
#include "absl/strings/string_view.h"
#include "intrinsic/util/proto/status_specs.h"
#include "intrinsic/util/proto/type_url.h"
#include "intrinsic/util/status/status_builder.h"

namespace intrinsic {

absl::StatusOr<ParsedUrl> ParseTypeUrl(
    std::string_view type_url ABSL_ATTRIBUTE_LIFETIME_BOUND) {
  ParsedUrl parsed_url = {.type_url = type_url};

  if (!type_url.starts_with(kIntrinsicTypeUrlPrefix)) {
    return StatusBuilder(absl::StatusCode::kInvalidArgument)
        .SetExtendedStatus(
            util::proto::kExtendedStatusComponent, util::proto::kInvalidUrlCode,
            {.title = util::proto::kInvalidUrlTitle,
             .user_message =
                 absl::StrFormat("Type URL '%s' does not start with '%s'",
                                 type_url, kIntrinsicTypeUrlPrefix),
             .user_instructions = util::proto::kInvalidUrlInstructions});
  }

  parsed_url.prefix = type_url.substr(0, kIntrinsicTypeUrlPrefix.length());

  std::string_view::size_type second_slash_pos =
      type_url.find(kTypeUrlSeparator, kIntrinsicTypeUrlPrefix.length());
  if (second_slash_pos == std::string_view::npos) {
    return StatusBuilder(absl::StatusCode::kInvalidArgument)
        .SetExtendedStatus(
            util::proto::kExtendedStatusComponent, util::proto::kInvalidUrlCode,
            {.title = util::proto::kInvalidUrlTitle,
             .user_message = absl::StrFormat(
                 "Type URL '%s' is missing separator after prefix", type_url),
             .user_instructions = util::proto::kInvalidUrlInstructions});
  }

  parsed_url.area =
      type_url.substr(kIntrinsicTypeUrlPrefix.length(),
                      second_slash_pos - kIntrinsicTypeUrlPrefix.length());

  std::string_view::size_type last_slash_pos =
      type_url.find_last_of(kTypeUrlSeparator);
  if (last_slash_pos <= second_slash_pos) {
    return StatusBuilder(absl::StatusCode::kInvalidArgument)
        .SetExtendedStatus(
            util::proto::kExtendedStatusComponent, util::proto::kInvalidUrlCode,
            {.title = util::proto::kInvalidUrlTitle,
             .user_message = absl::StrFormat(
                 "Type URL '%s' is missing are or message type", type_url),
             .user_instructions = util::proto::kInvalidUrlInstructions});
  }

  parsed_url.path = type_url.substr(second_slash_pos + 1,
                                    last_slash_pos - second_slash_pos - 1);
  parsed_url.message_type = type_url.substr(last_slash_pos + 1);

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
