// Copyright 2023 Intrinsic Innovation LLC

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
#include "intrinsic/scene/conversion/scene_object_model_utils.h"
#include "intrinsic/scene/proto/v1/entity.pb.h"
#include "intrinsic/scene/proto/v1/simulation_spec.pb.h"
#include "intrinsic/scene/sdf/sdf_util.h"
#include "intrinsic/scene/sdf/xml_utils.h"
#include "intrinsic/simulation/world/sim_plugins.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/world/proto/robot_component.pb.h"
#include "intrinsic/world/proto/kinematics_component.pb.h"
#include "sdf/Element.hh"
#include "sdf/Model.hh"

namespace intrinsic {
namespace scene_object {

namespace {

using DeviceSpec =
    ::intrinsic_proto::scene_object::v1::RobotSimPluginSpec::DeviceSpec;
using ::intrinsic_proto::scene_object::v1::Entity;
using ::intrinsic_proto::scene_object::v1::SimulationSpec;
using IconSimSpec =
    ::intrinsic_proto::scene_object::v1::RobotSimPluginSpec::IconSimSpec;

using IconSimDevice = ::intrinsic_proto::world::RobotComponent::IconSimDevice;
using IconSimPluginSpec =
    ::intrinsic_proto::world::RobotComponent::IconSimPluginSpec;
using MultiCameraPluginSpec =
    ::intrinsic_proto::scene_object::v1::MultiCameraPluginSpec;
using RobotSimPluginSpec =
    ::intrinsic_proto::scene_object::v1::RobotSimPluginSpec;

using ::intrinsic::sdf::GetAttributeAsString;
using ::intrinsic::sdf::GetChildrenByTag;
using ::intrinsic::sdf::GetCompactXml;
using ::intrinsic::simulation::PluginSdfTrait;
using ::sdf::ElementConstPtr;

DeviceSpec ConvertFromIconSimDevice(const IconSimDevice& device) {
  DeviceSpec device_spec;
  device_spec.set_name(device.name());
  if (!device.type().empty()) {
    device_spec.set_type(device.type());
  }
  if (!device.joint().empty()) {
    device_spec.set_joint_entity(device.joint());
  }
  if (device.has_initial()) {
    if (device_spec.has_joint_entity()) {
      device_spec.set_initial_state(device.initial());
    } else {
      LOG(WARNING) << absl::Substitute(
          "Not setting initial state $0 as no joint was provided.",
          device.initial());
    }
  }

  return device_spec;
}

}  // namespace

absl::StatusOr<std::optional<SimulationSpec>> ExtractSimulationSpecFromSdf(
    const ::sdf::Model& model, const std::vector<Entity>& entities,
    UnsupportedPluginsProcessing unsupported_plugins) {
  const auto kinematic_joints = GetNonFixedJointNames(entities);

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
    if (simulation::FilenameMatchesPluginSpec<IconSimPluginSpec>(filename)) {
      if (kinematic_joints.empty()) {
        LOG(WARNING) << absl::Substitute(
            "Skip '$0': cannot configure this plugin since the model does not "
            "have any kinematic joints.",
            filename);
        continue;
      }
      INTR_ASSIGN_OR_RETURN(const IconSimPluginSpec icon_sim_spec,
                            PluginSdfTrait<IconSimPluginSpec>::Parse(xml));
      if (!sim_spec.has_value()) {
        sim_spec = SimulationSpec();
      }

      // Parsing the icon_sim_spec is enough to set the field.
      sim_spec->mutable_robot()->mutable_icon_sim_spec();
    } else if (simulation::FilenameMatchesPluginSpec<IconSimDevice>(filename)) {
      std::vector<IconSimDevice> icon_sim_devices;
      INTR_RETURN_IF_ERROR(
          PluginSdfTrait<IconSimDevice>::Parse(xml, &icon_sim_devices));

      if (!icon_sim_devices.empty()) {
        if (!sim_spec.has_value()) {
          sim_spec = SimulationSpec();
        }
      }

      // Validates that joints specified in the sim plugin are valid.
      for (const auto& device : icon_sim_devices) {
        *sim_spec->mutable_robot()->add_device_specs() =
            ConvertFromIconSimDevice(device);

        if (device.name().empty()) {
          return absl::InvalidArgumentError(
              "Failed to configure device simulation plugin since device name "
              "was not provided");
        }
        if (!device.has_joint()) {
          continue;
        }

        if (device.joint().empty()) {
          return absl::InvalidArgumentError(absl::Substitute(
              "Failed to configure device simulation plugin as device '$0' is "
              "configured to control an empty joint",
              device.name()));
        }

        if (!kinematic_joints.contains(device.joint())) {
          return absl::InvalidArgumentError(absl::Substitute(
              "Failed to configure device simulation plugin as device '$0' is "
              "configured to control a non-existent joint '$1'",
              device.name(), device.joint()));
        }
      }
    } else if (simulation::FilenameMatchesPluginSpec<MultiCameraPluginSpec>(
                   filename)) {
      INTR_ASSIGN_OR_RETURN(auto multi_camera_plugin_spec,
                            PluginSdfTrait<MultiCameraPluginSpec>::Parse(xml));

      if (multi_camera_plugin_spec.sensors_size() == 0) {
        return absl::InvalidArgumentError(
            "Failed to configure multi camera plugin as no sensors were "
            "specified");
      }

      // Collect the sensor names from the entities.
      absl::flat_hash_set<std::string> sensor_names;
      for (const auto& entity : entities) {
        if (entity.has_sensor()) {
          sensor_names.insert(entity.name());
        }
      }

      // Validate that all sensors referenced by the plugin exist.
      for (int i = 0; i < multi_camera_plugin_spec.sensors_size(); ++i) {
        const std::string& sensor_name_in_spec =
            multi_camera_plugin_spec.sensors(i).name();
        if (!sensor_names.contains(sensor_name_in_spec)) {
          return absl::NotFoundError(absl::Substitute(
              "MultiCameraPlugin references sensor '$0', which does not exist.",
              sensor_name_in_spec));
        }
      }

      if (!sim_spec.has_value()) {
        sim_spec = SimulationSpec();
      }
      *sim_spec->mutable_multi_camera_plugin() =
          std::move(multi_camera_plugin_spec);
    } else {
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
