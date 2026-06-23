// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ASSETS_DEPENDENCIES_PLATFORM_H_
#define INTRINSIC_ASSETS_DEPENDENCIES_PLATFORM_H_

#include "absl/strings/string_view.h"

namespace intrinsic::assets::dependencies::platform {

inline constexpr absl::string_view kRuntimeAssetID = "ai.intrinsic.runtime";
inline constexpr absl::string_view kRuntimeInstanceName = "intrinsic_runtime";

}  // namespace intrinsic::assets::dependencies::platform

#endif  // INTRINSIC_ASSETS_DEPENDENCIES_PLATFORM_H_
