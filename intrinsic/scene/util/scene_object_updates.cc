// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/scene/util/scene_object_updates.h"

#include <cstddef>
#include <utility>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "intrinsic/scene/proto/v1/scene_object.pb.h"
#include "intrinsic/scene/proto/v1/scene_object_updates.pb.h"
#include "intrinsic/scene/util/scene_object_updates_internal.h"
#include "intrinsic/scene/validate/scene_object_validation.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {
namespace scene_object {
namespace {

using ::intrinsic_proto::scene_object::v1::SceneObject;
using ::intrinsic_proto::scene_object::v1::SceneObjectInstanceUpdates;
using ::intrinsic_proto::scene_object::v1::SceneObjectUpdates;

// We use Extras... so that we can have 0 or 1 params specified. For instance
// editing there is no need to specify the update type because it is implied to
// be instance.
template <typename ProtoT, typename... Extras>
absl::StatusOr<SceneObjectUpdateResult> ProcessSceneObjectUpdatesDefault(
    SceneObject&& object, const ProtoT& updates, Extras... extras) {
  for (const auto& update : updates.updates()) {
    INTR_ASSIGN_OR_RETURN(object, internal::ProcessSceneObjectUpdates(
                                      std::move(object), update, extras...));
  }

  return SceneObjectUpdateResult{
      .result = std::move(object),
      .update_errors = {},
  };
}

template <typename ProtoT, typename... Extras>
absl::StatusOr<SceneObjectUpdateResult> ProcessSceneObjectUpdatesSkipFailed(
    SceneObject&& object, const ProtoT& updates, Extras... extras) {
  SceneObject fallback;
  std::vector<SceneObjectUpdateResult::UpdateError> update_errors;

  for (size_t i = 0; i < updates.updates_size(); ++i) {
    const auto& update = updates.updates(i);
    // Grab a copy of the object before the update so we can revert to it in
    // case of an error.
    fallback = object;

    // Attempt the update on the object.
    auto object_result = internal::ProcessSceneObjectUpdates(std::move(object),
                                                             update, extras...);

    if (object_result.ok()) {
      // If the update worked, we take the result and store it inside object.
      object = std::move(object_result).value();
    } else {
      // Add the status to the set of received statuses from the updates.
      update_errors.push_back({
          .index = i,
          .status = object_result.status(),
      });

      // If the update failed, we take the fallback object and revert to it,
      // adding the error to our list of non fatal errors.
      object = std::move(fallback);
    }
  }

  return SceneObjectUpdateResult{
      .result = std::move(object),
      .update_errors = std::move(update_errors),
  };
}

// Process the switch between the different UpdatePolicy types.
template <typename ProtoT, typename... Extras>
absl::StatusOr<SceneObjectUpdateResult> ProcessSceneObjectUpdatesHelper(
    SceneObject&& object, const ProtoT& updates, UpdatePolicy update_policy,
    bool validate_original_scene_object, Extras... extras) {
  if (validate_original_scene_object) {
    // Validate the input scene object before doing anything else.
    INTR_RETURN_IF_ERROR(scene_object::ValidateSceneObject(object));
  }

  switch (update_policy) {
    case UpdatePolicy::kDefault:
      return ProcessSceneObjectUpdatesDefault(std::move(object), updates,
                                              extras...);
    case UpdatePolicy::kSkipFailed:
      return ProcessSceneObjectUpdatesSkipFailed(std::move(object), updates,
                                                 extras...);
    default:
      return absl::InvalidArgumentError(
          absl::StrCat("Unsupported update policy: ", update_policy));
  }
}

}  // namespace

absl::StatusOr<SceneObjectUpdateResult> ProcessSceneObjectUpdates(
    SceneObject object, const SceneObjectUpdates& updates,
    SceneObjectUpdateOptions update_options) {
  return ProcessSceneObjectUpdatesHelper(
      std::move(object), updates, update_options.update_policy,
      update_options.validate_original_scene_object,
      update_options.update_type);
}

absl::StatusOr<SceneObjectUpdateResult> ProcessSceneObjectUpdates(
    SceneObject object, const SceneObjectInstanceUpdates& updates,
    UpdatePolicy update_policy) {
  return ProcessSceneObjectUpdatesHelper(
      std::move(object), updates, update_policy,
      /*validate_original_scene_object=*/true);
}

}  // namespace scene_object
}  // namespace intrinsic
