// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_PROTO_PARSED_TYPE_URL_H_
#define INTRINSIC_UTIL_PROTO_PARSED_TYPE_URL_H_

#include <ostream>
#include <string_view>

#include "absl/base/attributes.h"
#include "absl/status/statusor.h"

namespace intrinsic {

// Structure of separate parsed elements of a Type URL.
// Intrinsic Type URLs have the form:
// type.intrinsic.ai/<area>/<path>
// The elements are:
// prefix: always type.intrinsic.ai
// area: the designated resolver responsible, e.g., skill
// path: the resolver/area-specific path, e.g., <id>/<version> for a skill.
// message type: a specific full name of a proto
//
// Note that resolvers must identify a particular file descriptor set only by
// the area and path, not by message type. Caching will be performed based on
// the type URL without the message type.
//
// Example:
// type.intrinsic.ai/skills/my_skill/1.0.0/com.example.MyParameterProto
// |---- prefix ----|-area-|-----path-----|----- message type --------|
//
struct ParsedUrl {
  std::string_view type_url;
  std::string_view prefix;
  std::string_view area;
  std::string_view path;
  std::string_view message_type;
};

std::ostream& operator<<(std::ostream& os, const ParsedUrl& parsed_url);

absl::StatusOr<ParsedUrl> ParseTypeUrl(
    std::string_view type_url ABSL_ATTRIBUTE_LIFETIME_BOUND);

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_PROTO_PARSED_TYPE_URL_H_
