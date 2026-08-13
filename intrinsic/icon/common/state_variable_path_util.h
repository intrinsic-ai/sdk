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

#ifndef INTRINSIC_ICON_COMMON_STATE_VARIABLE_PATH_UTIL_H_
#define INTRINSIC_ICON_COMMON_STATE_VARIABLE_PATH_UTIL_H_

#include <cstddef>
#include <optional>
#include <string>

#include "absl/types/span.h"

namespace intrinsic::icon {

// This struct represents one node of a state variable path.
// A node consists of a name and optionally an array index.
struct StateVariablePathNode {
  std::string ToString() const;

  template <typename Sink>
  friend void AbslStringify(Sink& sink, const StateVariablePathNode& node) {
    sink.Append(node.ToString());
  }

  // Name of the node. Must not be empty.
  std::string name;
  // Optional index of an array that is represented by this node.
  std::optional<size_t> index;
};

inline bool operator==(const StateVariablePathNode& lhs,
                       const StateVariablePathNode& rhs) {
  return lhs.name == rhs.name && lhs.index == rhs.index;
}

// Builds a state variable path from a span of nodes by prepending
// `kStateVariablePathPrefix` and joining the nodes with
// `kStateVariablePathSeparator`.
std::string BuildStateVariablePath(
    absl::Span<const StateVariablePathNode> nodes);

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_COMMON_STATE_VARIABLE_PATH_UTIL_H_
