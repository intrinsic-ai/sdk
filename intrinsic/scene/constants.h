// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SCENE_CONSTANTS_H_
#define INTRINSIC_SCENE_CONSTANTS_H_

#include <array>

#include "absl/strings/string_view.h"
#include "intrinsic/scene/user_data_keys.h"

namespace intrinsic {
namespace scene_object {
constexpr absl::string_view kSceneObjectAssetUserDataKey =
    "FLOWSTATE_ASSET_USER_DATA";
constexpr int kSceneObjectNumReservedUserDataKeys = 4;
constexpr std::array<absl::string_view, kSceneObjectNumReservedUserDataKeys>
    kSceneObjectReservedUserDataKeys = {
        sdf::kGazeboPlugins,
        sdf::kGazeboCollisionSurface,
        sdf::kGazeboJointPhysics,
        sdf::kSdfLights,
    };
constexpr double kMaxApplicationLimitsMultiplier = 0.95;

}  // namespace scene_object
}  // namespace intrinsic

#endif  // INTRINSIC_SCENE_CONSTANTS_H_
