// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
