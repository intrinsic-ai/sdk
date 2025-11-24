// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_PROTO_PARSED_TYPE_URL_H_
#define INTRINSIC_UTIL_PROTO_PARSED_TYPE_URL_H_

#include <ostream>
#include <string>
#include <string_view>

#include "absl/status/statusor.h"

namespace intrinsic {

// A type URL of an Any proto consists of a type URL prefix followed by a '/'
// and the full message type. Intrinsic Type URLs have a specific format for the
// type URL prefix.
// The structure of separate parsed elements of a Type URL of
// Intrinsic Type URLs has the form:
// type.intrinsic.ai/<area>/<path>
// The elements are:
// custom prefix: always type.intrinsic.ai
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
// |---------- type URL prefix -----------|----- message type --------|
// |- custom prefix-|-area-|-----path-----|----- message type --------|
//
struct ParsedUrl {
  std::string type_url;
  std::string prefix;
  std::string area;
  std::string path;
  std::string message_type;
};

std::ostream& operator<<(std::ostream& os, const ParsedUrl& parsed_url);

// Parses a complete type URL into its parts.
absl::StatusOr<ParsedUrl> ParseTypeUrl(std::string_view type_url);

// Parses a type URL prefix, i.e., a type URL without the message type.
// This can end in a '/' or not.
// The message_type field in the returned ParsedUrl will be empty.
absl::StatusOr<ParsedUrl> ParseTypeUrlPrefix(std::string_view type_url_prefix);

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_PROTO_PARSED_TYPE_URL_H_
