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
