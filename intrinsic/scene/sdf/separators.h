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

#ifndef INTRINSIC_SCENE_SDF_SEPARATORS_H
#define INTRINSIC_SCENE_SDF_SEPARATORS_H

#include "absl/strings/string_view.h"

namespace intrinsic {

// The separator that gazebo uses to scope names of entities in the SDF.
inline constexpr absl::string_view kSdfNameSeparator = "::";

// The default separator that we use internally to scope gazebo names in order
// to avoid ambiguity when link names contain kSdfNameSeparator.
inline constexpr absl::string_view kWorldSdfDefaultSeparator = "/";

// The separator we use to append intrinsic world entity ids to model names in
// order to preserve the property that model names are unique in the SDF.
inline constexpr absl::string_view kSdfNamePartsSeparator = "__";

}  // namespace intrinsic

#endif  // INTRINSIC_SCENE_SDF_SEPARATORS_H
