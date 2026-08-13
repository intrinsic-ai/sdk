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

#ifndef INTRINSIC_ICON_TOOLS_GENERATE_DOCUMENTATION_H_
#define INTRINSIC_ICON_TOOLS_GENERATE_DOCUMENTATION_H_

#include <string>
#include <vector>

#include "absl/status/statusor.h"
#include "absl/types/span.h"
#include "intrinsic/icon/proto/v1/types.pb.h"

namespace intrinsic {
namespace icon {

std::string GenerateActionNames(
    absl::Span<const intrinsic_proto::icon::v1::ActionSignature> signatures);

absl::StatusOr<std::string> GenerateSingleActionDocumentation(
    const intrinsic_proto::icon::v1::ActionSignature& signature);

absl::StatusOr<std::string> GenerateDocumentation(
    absl::Span<const intrinsic_proto::icon::v1::ActionSignature> signatures);

absl::StatusOr<std::string> GenerateDocumentation(
    absl::Span<const intrinsic_proto::icon::v1::ActionSignature> signatures,
    absl::Span<const std::vector<std::string>> compatible_parts,
    bool with_toc_header = true, bool with_devsite_header = false);

}  // namespace icon
}  // namespace intrinsic

#endif  // INTRINSIC_ICON_TOOLS_GENERATE_DOCUMENTATION_H_
