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

#ifndef INTRINSIC_SCENE_SDF_SIM_SPEC_FROM_SDF_H__
#define INTRINSIC_SCENE_SDF_SIM_SPEC_FROM_SDF_H__

#include <optional>
#include <vector>

#include "absl/status/statusor.h"
#include "intrinsic/scene/proto/v1/scene_object.pb.h"
#include "intrinsic/scene/proto/v1/simulation_spec.pb.h"
#include "sdf/Model.hh"

namespace intrinsic {
namespace scene_object {

enum class UnsupportedPluginsProcessing {
  kDefault = 0,
  kFail = 0,
  kSkip = 1,
  kInline = 2,
};

// Returns the simulation spec from the given SDF model after other entities
// like links, joints, etc. have already been extracted.
//
// entities: already extracted entities from the SDF model like links, joints,
// etc. These are required to validate that references to links, joints, etc.
// in the plugin are valid.
// unsupported_plugins: Describes how to handle unsupported plugins.
absl::StatusOr<std::optional<intrinsic_proto::scene_object::v1::SimulationSpec>>
ExtractSimulationSpecFromSdf(
    const ::sdf::Model& model,
    const std::vector<intrinsic_proto::scene_object::v1::Entity>& entities,
    UnsupportedPluginsProcessing unsupported_plugins);

}  // namespace scene_object

}  // namespace intrinsic

#endif  // INTRINSIC_SCENE_SDF_SIM_SPEC_FROM_SDF_H__
