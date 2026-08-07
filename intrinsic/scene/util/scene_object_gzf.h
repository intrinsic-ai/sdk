// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SCENE_UTIL_SCENE_OBJECT_GZF_H_
#define INTRINSIC_SCENE_UTIL_SCENE_OBJECT_GZF_H_

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/scene/proto/v1/scene_object.pb.h"
#include "intrinsic/world/gzfile/gzfile.h"

namespace intrinsic {
namespace scene_object {

// Adds the given SceneObject proto to the given GZFile.
absl::Status AddSceneObjectToGzf(
    const intrinsic_proto::scene_object::v1::SceneObject& scene_object,
    GZFile& gzfile);

// Extracts a SceneObject proto from the given GZFile that has previously been
// added with AddSceneObjectToGzf(...) above. Returns an error if no SceneObject
// is stored in the GZFile.
absl::StatusOr<intrinsic_proto::scene_object::v1::SceneObject>
GetSceneObjectFromGzf(const GZFile& gzfile);

}  // namespace scene_object
}  // namespace intrinsic

#endif  // INTRINSIC_SCENE_UTIL_SCENE_OBJECT_GZF_H_
