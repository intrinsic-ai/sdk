// Copyright 2023 Intrinsic Innovation LLC

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
