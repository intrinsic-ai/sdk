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

#include "intrinsic/scene/sdf/sim_spec_from_sdf.h"

#include <optional>
#include <string>
#include <utility>
#include <vector>

#include "absl/container/flat_hash_set.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/strings/substitute.h"
#include "intrinsic/scene/proto/v1/entity.pb.h"
#include "intrinsic/scene/proto/v1/simulation_spec.pb.h"
#include "intrinsic/scene/sdf/sdf_util.h"
#include "intrinsic/scene/sdf/xml_utils.h"
#include "intrinsic/util/status/status_macros.h"
#include "sdf/Element.hh"
#include "sdf/Model.hh"

namespace intrinsic {
namespace scene_object {

namespace {

using ::intrinsic::sdf::GetAttributeAsString;
using ::intrinsic::sdf::GetChildrenByTag;
using ::intrinsic::sdf::GetCompactXml;
using ::intrinsic_proto::scene_object::v1::Entity;
using ::intrinsic_proto::scene_object::v1::SimulationSpec;
using ::sdf::ElementConstPtr;

}  // namespace

absl::StatusOr<std::optional<SimulationSpec>> ExtractSimulationSpecFromSdf(
    const ::sdf::Model& model, const std::vector<Entity>& entities,
    UnsupportedPluginsProcessing unsupported_plugins) {
  std::optional<SimulationSpec> sim_spec;
  if (model.Static()) {
    sim_spec = SimulationSpec();
    sim_spec->set_is_static(true);
  }

  // Processes relevant "plugins" from the SDF.
  for (const ElementConstPtr& plugin :
       GetChildrenByTag(model.Element(), "plugin")) {
    INTR_ASSIGN_OR_RETURN(std::string xml, GetCompactXml(plugin));
    INTR_ASSIGN_OR_RETURN(std::string filename,
                          GetAttributeAsString(plugin, "filename"));
    {
      switch (unsupported_plugins) {
        case UnsupportedPluginsProcessing::kFail:
          return absl::UnimplementedError(
              absl::Substitute("Unsupported plugin: '$0'.", filename));
        case UnsupportedPluginsProcessing::kSkip:
          LOG(WARNING) << absl::Substitute("Skipping unsupported plugin: '$0'.",
                                           filename);
          break;
        case UnsupportedPluginsProcessing::kInline:
          if (!sim_spec.has_value()) {
            sim_spec = SimulationSpec();
          }
          LOG(WARNING) << absl::Substitute("Inlining unsupported plugin: '$0'.",
                                           filename);
          sim_spec->add_extra_inlined_plugins(xml);
          break;
      }
    }
  }

  return sim_spec;
}

}  // namespace scene_object
}  // namespace intrinsic
