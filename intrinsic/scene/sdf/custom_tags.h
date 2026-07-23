// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SCENE_SDF_CUSTOM_TAGS_H_
#define INTRINSIC_SCENE_SDF_CUSTOM_TAGS_H_

#include "absl/strings/string_view.h"

namespace intrinsic {
namespace sdf {
// Custom element under sdf <model> elements.
// This element is expected to contain data matching intrinsic::CartesianLimits.
inline constexpr absl::string_view kCartesianLimitsCustomElement =
    "intrinsic:cartesian_limits";
// Custom elements under <intrinsic:cartesian_limits> elements.
inline constexpr absl::string_view kMinTranslationalPositionCustomElement =
    "intrinsic:min_translational_position";
inline constexpr absl::string_view kMaxTranslationalPositionCustomElement =
    "intrinsic:max_translational_position";
inline constexpr absl::string_view kMinTranslationalVelocityCustomElement =
    "intrinsic:min_translational_velocity";
inline constexpr absl::string_view kMaxTranslationalVelocityCustomElement =
    "intrinsic:max_translational_velocity";
inline constexpr absl::string_view kMinTranslationalAccelerationCustomElement =
    "intrinsic:min_translational_acceleration";
inline constexpr absl::string_view kMaxTranslationalAccelerationCustomElement =
    "intrinsic:max_translational_acceleration";
inline constexpr absl::string_view kMinTranslationalJerkCustomElement =
    "intrinsic:min_translational_jerk";
inline constexpr absl::string_view kMaxTranslationalJerkCustomElement =
    "intrinsic:max_translational_jerk";
inline constexpr absl::string_view kMaxRotationalVelocityCustomElement =
    "intrinsic:max_rotational_velocity";
inline constexpr absl::string_view kMaxRotationalAccelerationCustomElement =
    "intrinsic:max_rotational_acceleration";
inline constexpr absl::string_view kMaxRotationalJerkCustomElement =
    "intrinsic:max_rotational_jerk";
inline constexpr absl::string_view kControlFrequencyHz =
    "intrinsic:control_frequency_hz";

// The name of the IK solver used by ICON for the kinematic chain
// defined in the model.
inline constexpr absl::string_view kIkSolverCustomElement =
    "intrinsic:ik_solver";

// Custom element under sdf <model> elements containing user data.
inline constexpr absl::string_view kUserDataCustomElement =
    "intrinsic:user_data";
// The name of the attributed used to specify the name of the final link of the
// ik solver.
inline constexpr absl::string_view kIkSolverLinkNameAttribute = "tip_link_name";

// Custom element under sdf joint <Limit> elements.
inline constexpr absl::string_view kAccelerationCustomElement =
    "intrinsic:acceleration";
inline constexpr absl::string_view kJerkCustomElement = "intrinsic:jerk";

// Custom attribute under sdf <frame> elements.
// This attribute is used to specify if a custom object world Frame should be
// created for a <frame> element.
inline constexpr absl::string_view kCreateEntityCustomAttribute =
    "intrinsic:create_entity";

// Custom attribute under sdf <frame> elements.
// This attribute is used to specify if a custom object world Frame should be
// marked as an attachment frame.
// Implies `intrinsic:create_entity` attribute.
inline constexpr absl::string_view kCreateAttachmentEntityCustomAttribute =
    "intrinsic:create_attachment_entity";

}  // namespace sdf
}  // namespace intrinsic

#endif  // INTRINSIC_SCENE_SDF_CUSTOM_TAGS_H_
