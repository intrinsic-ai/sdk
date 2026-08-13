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

#include "intrinsic/icon/common/state_variable_path_util.h"

#include <iterator>
#include <string>
#include <vector>

#include "absl/algorithm/container.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/str_join.h"
#include "absl/types/span.h"
#include "intrinsic/icon/common/state_variable_path_constants.h"

namespace intrinsic::icon {

std::string BuildStateVariablePath(
    const absl::Span<const StateVariablePathNode> nodes) {
  std::vector<std::string> node_strings;
  absl::c_transform(
      nodes, std::back_inserter(node_strings),
      [](const StateVariablePathNode& node) { return node.ToString(); });
  return absl::StrCat(kStateVariablePathPrefix,
                      absl::StrJoin(node_strings, kStateVariablePathSeparator));
}

std::string StateVariablePathNode::ToString() const {
  std::string result = name;
  if (index.has_value()) {
    absl::StrAppend(&result, "[", index.value(), "]");
  }
  return result;
}
}  // namespace intrinsic::icon
