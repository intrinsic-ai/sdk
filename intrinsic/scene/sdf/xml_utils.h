// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SCENE_SDF_XML_UTILS_H_
#define INTRINSIC_SCENE_SDF_XML_UTILS_H_

#include <string>

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "sdf/Element.hh"

namespace intrinsic::sdf {

// Returns a normalized XML string by parsing and pretty-printing it with
// collapsed whitespaces. This is useful for normalizing XML strings
// for comparison in tests.
absl::StatusOr<std::string> NormalizeXml(absl::string_view xml);

// Returns a compact XML string from `element` without extra whitespaces.
// This is useful for generating XML strings that are as small as possible.
absl::StatusOr<std::string> GetCompactXml(
    const ::sdf::ElementConstPtr& element);

// Returns a prettified XML string from `element` with indentation and
// newlines. This is useful for generating human-readable XML strings.
absl::StatusOr<std::string> GetPrettifiedXml(
    const ::sdf::ElementConstPtr& element);

// Escapes special XML characters in a string.
std::string EscapeXml(absl::string_view s);

}  // namespace intrinsic::sdf

#endif  // INTRINSIC_SCENE_SDF_XML_UTILS_H_
