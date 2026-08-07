// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SCENE_UTIL_SCENE_OBJECT_UPDATES_H_
#define INTRINSIC_SCENE_UTIL_SCENE_OBJECT_UPDATES_H_

#include <cstddef>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/scene/proto/v1/scene_object.pb.h"
#include "intrinsic/scene/proto/v1/scene_object_updates.pb.h"

namespace intrinsic {
namespace scene_object {

// Describes how to treat errors when applying updates to the SceneObject.
enum class UpdatePolicy {
  // Returns fatal errors on any update that cannot be applied.
  kDefault = 0,
  // Skips any updates that fail to apply, and continues to apply the rest of
  // the updates as best as it can. Returning all the statuses that were
  // received along the way.
  kSkipFailed,
};

enum class UpdateType {
  // All updates are permitted using the default setting.
  kDefault = 0,
  // All updates are permitted using kObjectType, it is equivalent to default.
  kObjectType = 1,
  // Structural changes are not permitted, but configuration changes are ok.
  kObjectInstance = 2,
};

// Options specified for the ProcessSceneObjectUpdates calls. See individual
// options for documentation on what they do.
struct SceneObjectUpdateOptions {
  UpdatePolicy update_policy = UpdatePolicy::kDefault;
  UpdateType update_type = UpdateType::kDefault;
  bool validate_original_scene_object = true;

  static SceneObjectUpdateOptions Default() {
    return SceneObjectUpdateOptions{
        .update_policy = UpdatePolicy::kDefault,
        .update_type = UpdateType::kDefault,
        .validate_original_scene_object = true,
    };
  }
};

// The result struct from calling ProcessSceneObjectUpdates, contains the
// resulting proto as well as the optionally filled update errors for any
// updates the had an error when being applied. The errors are only filled if
// the update policy is kSkipFailed.
struct SceneObjectUpdateResult {
  struct UpdateError {
    // The index of the update in the SceneObjectUpdates proto.
    size_t index;
    // The status of the failed update.
    absl::Status status;
  };

  // The resulting SceneObject with the updates that succeeded.
  intrinsic_proto::scene_object::v1::SceneObject result;
  // The errors status of all the updates that failed.
  std::vector<UpdateError> update_errors;
};

// Applies the given SceneObjectUpdates to the given SceneObject. It may return
// errors as part of the status or potentially as part of the result struct
// depending on the UpdatePolicy specified. Some updates may also be denied
// based on the update type specified.
absl::StatusOr<SceneObjectUpdateResult> ProcessSceneObjectUpdates(
    intrinsic_proto::scene_object::v1::SceneObject object,
    const intrinsic_proto::scene_object::v1::SceneObjectUpdates& updates,
    SceneObjectUpdateOptions update_options =
        SceneObjectUpdateOptions::Default());

// Applies the given SceneObjectInstanceUpdates to the given SceneObject. It may
// return errors as part of the status or potentially as part of the result
// struct depending on the UpdatePolicy specified. Some updates may also be
// denied based on the update type specified.
absl::StatusOr<SceneObjectUpdateResult> ProcessSceneObjectUpdates(
    intrinsic_proto::scene_object::v1::SceneObject object,
    const intrinsic_proto::scene_object::v1::SceneObjectInstanceUpdates&
        updates,
    UpdatePolicy update_policy = UpdatePolicy::kDefault);

}  // namespace scene_object
}  // namespace intrinsic

#endif  // INTRINSIC_SCENE_UTIL_SCENE_OBJECT_UPDATES_H_
