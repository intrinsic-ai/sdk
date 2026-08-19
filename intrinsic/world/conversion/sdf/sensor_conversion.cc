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

#include "intrinsic/world/conversion/sdf/sensor_conversion.h"

#include <cmath>
#include <cstdint>
#include <cstdlib>
#include <string>

#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/strings/substitute.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/scene/sdf/sdf_sensor_pose.h"
#include "intrinsic/scene/sdf/sdf_util.h"
#include "intrinsic/scene/sdf/xml_utils.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/world/component/sensor_component.h"
#include "intrinsic/world/hashing/hashing.h"
#include "intrinsic/world/proto/sensor_component.pb.h"
#include "sdf/Camera.hh"
#include "sdf/Element.hh"
#include "sdf/ForceTorque.hh"
#include "sdf/Lidar.hh"
#include "sdf/Noise.hh"
#include "sdf/Sensor.hh"

namespace intrinsic {
namespace sdf {

namespace {

using ForceTorqueDevicePluginSpec =
    intrinsic_proto::world::SensorComponent::ForceTorqueDevicePluginSpec;
using CameraPluginSpec =
    intrinsic_proto::world::SensorComponent::CameraPluginSpec;
using RangeFinderDevicePluginSpec =
    intrinsic_proto::world::SensorComponent::RangeFinderDevicePluginSpec;
using Noise = intrinsic_proto::world::SensorComponent::Noise;

const auto* kSdfFormatToIntrinsicFormat =
    new WorldHashMap< ::sdf::PixelFormatType,
                      intrinsic_proto::world::SensorComponent::Image::Format>({
        {::sdf::PixelFormatType::L_INT8,
         intrinsic_proto::world::SensorComponent::Image::FORMAT_L8},
        {::sdf::PixelFormatType::RGB_INT8,
         intrinsic_proto::world::SensorComponent::Image::FORMAT_R8G8B8},
        {::sdf::PixelFormatType::BGR_INT8,
         intrinsic_proto::world::SensorComponent::Image::FORMAT_B8G8R8},
        {::sdf::PixelFormatType::BAYER_RGGB8,
         intrinsic_proto::world::SensorComponent::Image::FORMAT_BAYER_RGGB8},
        {::sdf::PixelFormatType::BAYER_BGGR8,
         intrinsic_proto::world::SensorComponent::Image::FORMAT_BAYER_BGGR8},
        {::sdf::PixelFormatType::BAYER_GBRG8,
         intrinsic_proto::world::SensorComponent::Image::FORMAT_BAYER_GBRG8},
        {::sdf::PixelFormatType::BAYER_GRBG8,
         intrinsic_proto::world::SensorComponent::Image::FORMAT_BAYER_GRBG8},
    });

const auto* kSdfForceTorqueFrameToIntrinsicForceTorqueFrame = new WorldHashMap<
    ::sdf::ForceTorqueFrame,
    intrinsic_proto::world::SensorComponent::ForceTorque::Frame>({
    {::sdf::ForceTorqueFrame::PARENT,
     intrinsic_proto::world::SensorComponent::ForceTorque::FRAME_PARENT},
    {::sdf::ForceTorqueFrame::CHILD,
     intrinsic_proto::world::SensorComponent::ForceTorque::FRAME_CHILD},
    {::sdf::ForceTorqueFrame::SENSOR,
     intrinsic_proto::world::SensorComponent::ForceTorque::FRAME_SENSOR},
});

const auto*
    kSdfForceTorqueMeasureDirectionToIntrinsicForceTorqueMeasureDirection =
        new WorldHashMap< ::sdf::ForceTorqueMeasureDirection,
                          intrinsic_proto::world::SensorComponent::ForceTorque::
                              MeasureDirection>({
            {::sdf::ForceTorqueMeasureDirection::CHILD_TO_PARENT,
             intrinsic_proto::world::SensorComponent::ForceTorque::
                 MEASURE_DIRECTION_CHILD_TO_PARENT},
            {::sdf::ForceTorqueMeasureDirection::PARENT_TO_CHILD,
             intrinsic_proto::world::SensorComponent::ForceTorque::
                 MEASURE_DIRECTION_PARENT_TO_CHILD},
        });

// The camera visibility mask is for simulation use only.
// The default camera visibility mask here (0x0000FFFF) is different from the
// default value in SDF (0xFFFFFFFF). The default SDF visibility mask allows
// cameras to see everything in simulation. However, that is not always what
// we want. For example, the plenoptic unit has an IR camera that can see
// projected IR pattern which is not visible to other regular RGB cameras.
// In order to achieve this effect in simulation, we create a custom visibility
// mask for the IR camera (with one of the 16 bits enabled) so that it matches
// the custom visibility flags of certain targets (IR beacons / projections)
// in the scene that have matching visibility flags, i.e. visibility flags with
// one fo the first 16 bits enabled only and rest are zeros.
constexpr uint32_t kDefaultCameraVisibilityMask = 0x0000FFFF;

Noise IntrinsicNoiseFromSdfNoise(const ::sdf::Noise& noise_sdf) {
  Noise noise_proto;
  switch (noise_sdf.Type()) {
    case ::sdf::NoiseType::GAUSSIAN:
      noise_proto.set_type(Noise::TYPE_GAUSSIAN);
      break;
    case ::sdf::NoiseType::GAUSSIAN_QUANTIZED:
      noise_proto.set_type(Noise::TYPE_GAUSSIAN_QUANTIZED);
      break;
    default:
      noise_proto.set_type(Noise::TYPE_UNSPECIFIED);
      break;
  }
  noise_proto.set_mean(noise_sdf.Mean());
  noise_proto.set_stddev(noise_sdf.StdDev());
  noise_proto.set_bias_mean(noise_sdf.BiasMean());
  noise_proto.set_bias_stddev(noise_sdf.BiasStdDev());
  noise_proto.set_dynamic_bias_stddev(noise_sdf.DynamicBiasStdDev());
  noise_proto.set_dynamic_bias_correlation_time(
      noise_sdf.DynamicBiasCorrelationTime());
  noise_proto.set_precision(noise_sdf.Precision());
  return noise_proto;
}

absl::StatusOr<intrinsic_proto::world::SensorComponent::CommonCameraProperties>
IntrinsicCommonCameraPropertiesFromSdfCamera(const ::sdf::Camera& sdf_camera) {
  intrinsic_proto::world::SensorComponent::CommonCameraProperties properties;
  properties.set_horizontal_fov(sdf_camera.HorizontalFov().Radian());

  // Process <image>
  if (sdf_camera.ImageWidth() == 0) {
    return absl::InvalidArgumentError(
        "Expected positive camera image width, provided 0.");
  }
  if (sdf_camera.ImageHeight() == 0) {
    return absl::InvalidArgumentError(
        "Expected positive camera image height, provided 0.");
  }
  auto& intr_image = *properties.mutable_image();
  intr_image.set_width(sdf_camera.ImageWidth());
  intr_image.set_height(sdf_camera.ImageHeight());
  if (!kSdfFormatToIntrinsicFormat->contains(sdf_camera.PixelFormat())) {
    return absl::InvalidArgumentError(
        absl::Substitute("Unsupported sdf camera pixel format: $0",
                         sdf_camera.PixelFormatStr()));
  }
  intr_image.set_format(
      kSdfFormatToIntrinsicFormat->at(sdf_camera.PixelFormat()));

  // Process <clip>
  auto& intr_clip = *properties.mutable_clip();
  intr_clip.set_near(sdf_camera.NearClip());
  intr_clip.set_far(sdf_camera.FarClip());

  // Process <noise>
  const auto& sdf_noise = sdf_camera.ImageNoise();
  if (sdf_noise.Type() == ::sdf::NoiseType::GAUSSIAN) {
    auto& intr_noise = *properties.mutable_noise();
    intr_noise.set_type(
        intrinsic_proto::world::SensorComponent::Noise::TYPE_GAUSSIAN);
    intr_noise.set_mean(sdf_noise.Mean());
    intr_noise.set_stddev(sdf_noise.StdDev());
  }

  // Process optional <lens>.
  if (sdf_camera.HasLensIntrinsics()) {
    auto& intr_intrinsics = *properties.mutable_intrinsics();
    intr_intrinsics.set_fx(sdf_camera.LensIntrinsicsFx());
    intr_intrinsics.set_fy(sdf_camera.LensIntrinsicsFy());
    intr_intrinsics.set_cx(sdf_camera.LensIntrinsicsCx());
    intr_intrinsics.set_cy(sdf_camera.LensIntrinsicsCy());
    const double width = intr_image.width();
    const double fx = intr_intrinsics.fx();
    const double fov = 2 * std::atan2(width, 2 * fx);
    if (std::abs(fov - properties.horizontal_fov()) > 1e-3) {
      LOG(WARNING) << "The computed horizontal FOV is " << fov
                   << " while the specified horizontal FOV is "
                   << properties.horizontal_fov();
      LOG(WARNING) << "The computed value will be used.";
      properties.set_horizontal_fov(fov);
    }
  }

  // Process <distortion>
  auto& intr_distortion = *properties.mutable_distortion();
  intr_distortion.set_k1(sdf_camera.DistortionK1());
  intr_distortion.set_k2(sdf_camera.DistortionK2());
  intr_distortion.set_k3(sdf_camera.DistortionK3());
  intr_distortion.set_p1(sdf_camera.DistortionP1());
  intr_distortion.set_p2(sdf_camera.DistortionP2());

  // Process <visibility_mask>
  uint32_t visibility_mask = sdf_camera.VisibilityMask();
  visibility_mask = visibility_mask == 0xFFFFFFFF ? kDefaultCameraVisibilityMask
                                                  : visibility_mask;
  properties.set_visibility_mask(visibility_mask);

  return properties;
}

absl::StatusOr<intrinsic_proto::world::SensorComponent::ForceTorque>
IntrinsicForceTorqueFromSdfForceTorque(
    const ::sdf::ForceTorque& sdf_force_torque) {
  intrinsic_proto::world::SensorComponent::ForceTorque intr_ft;
  intr_ft.set_frame(
      intrinsic_proto::world::SensorComponent::ForceTorque::FRAME_CHILD);
  intr_ft.set_measure_direction(
      intrinsic_proto::world::SensorComponent::ForceTorque::
          MEASURE_DIRECTION_CHILD_TO_PARENT);

  auto sdf_frame = sdf_force_torque.Frame();
  if (!kSdfForceTorqueFrameToIntrinsicForceTorqueFrame->contains(sdf_frame)) {
    return absl::InvalidArgumentError(
        absl::Substitute("Unsupported frame in sdf force torque element: $0",
                         sdf_force_torque.Element()->ToString("")));
  }
  intr_ft.set_frame(
      kSdfForceTorqueFrameToIntrinsicForceTorqueFrame->at(sdf_frame));

  auto sdf_measure_direction = sdf_force_torque.MeasureDirection();
  if (!kSdfForceTorqueMeasureDirectionToIntrinsicForceTorqueMeasureDirection
           ->contains(sdf_measure_direction)) {
    return absl::InvalidArgumentError(absl::Substitute(
        "Unsupported measure direction in sdf force torque element: $0",
        sdf_force_torque.Element()->ToString("")));
  }
  intr_ft.set_measure_direction(
      kSdfForceTorqueMeasureDirectionToIntrinsicForceTorqueMeasureDirection->at(
          sdf_measure_direction));
  // If present in the SDF, populate the force_noise and torque_noise fields in
  // the proto.
  if (sdf_force_torque.ForceXNoise().Element() != nullptr) {
    *intr_ft.mutable_force_noise()->mutable_x() =
        IntrinsicNoiseFromSdfNoise(sdf_force_torque.ForceXNoise());
  }
  if (sdf_force_torque.ForceYNoise().Element() != nullptr) {
    *intr_ft.mutable_force_noise()->mutable_y() =
        IntrinsicNoiseFromSdfNoise(sdf_force_torque.ForceYNoise());
  }
  if (sdf_force_torque.ForceZNoise().Element() != nullptr) {
    *intr_ft.mutable_force_noise()->mutable_z() =
        IntrinsicNoiseFromSdfNoise(sdf_force_torque.ForceZNoise());
  }
  // Same for torque
  if (sdf_force_torque.TorqueXNoise().Element() != nullptr) {
    *intr_ft.mutable_torque_noise()->mutable_x() =
        IntrinsicNoiseFromSdfNoise(sdf_force_torque.TorqueXNoise());
  }
  if (sdf_force_torque.TorqueYNoise().Element() != nullptr) {
    *intr_ft.mutable_torque_noise()->mutable_y() =
        IntrinsicNoiseFromSdfNoise(sdf_force_torque.TorqueYNoise());
  }
  if (sdf_force_torque.TorqueZNoise().Element() != nullptr) {
    *intr_ft.mutable_torque_noise()->mutable_z() =
        IntrinsicNoiseFromSdfNoise(sdf_force_torque.TorqueZNoise());
  }

  return intr_ft;
}

absl::StatusOr<intrinsic_proto::world::SensorComponent::Lidar>
IntrinsicLidarFromSdfLidar(const ::sdf::Lidar& sdf_lidar) {
  intrinsic_proto::world::SensorComponent::Lidar intr_lidar;
  auto& horizontal_scan = *intr_lidar.mutable_horizontal();
  horizontal_scan.set_samples(sdf_lidar.HorizontalScanSamples());
  horizontal_scan.set_resolution(sdf_lidar.HorizontalScanResolution());
  horizontal_scan.set_min_angle(sdf_lidar.HorizontalScanMinAngle().Radian());
  horizontal_scan.set_max_angle(sdf_lidar.HorizontalScanMaxAngle().Radian());

  auto& vertical_scan = *intr_lidar.mutable_vertical();
  vertical_scan.set_samples(sdf_lidar.VerticalScanSamples());
  vertical_scan.set_resolution(sdf_lidar.VerticalScanResolution());
  vertical_scan.set_min_angle(sdf_lidar.VerticalScanMinAngle().Radian());
  vertical_scan.set_max_angle(sdf_lidar.VerticalScanMaxAngle().Radian());

  auto& intr_range = *intr_lidar.mutable_range();
  intr_range.set_min_distance(sdf_lidar.RangeMin());
  intr_range.set_max_distance(sdf_lidar.RangeMax());
  intr_range.set_resolution((sdf_lidar.RangeResolution()));

  auto& intr_noise = *intr_lidar.mutable_noise();
  intr_noise.set_type(
      intrinsic_proto::world::SensorComponent::Noise::TYPE_GAUSSIAN);
  intr_noise.set_mean(sdf_lidar.LidarNoise().Mean());
  intr_noise.set_stddev(sdf_lidar.LidarNoise().StdDev());

  return intr_lidar;
}

}  // namespace

absl::StatusOr<ParseSensorResult> ParseSensor(const ::sdf::Sensor& sdf_sensor) {
  ParseSensorResult result;

  result.name = sdf_sensor.Name();

  // Process <pose>.
  INTR_ASSIGN_OR_RETURN(const Pose3d parent_t_this,
                        ParseSemanticPose(sdf_sensor.SemanticPose()));

  ::sdf::SensorType type = sdf_sensor.Type();
  result.parent_t_sensor = SensorPoseFromSdf(parent_t_this, type);

  result.sensor_component = SensorComponent::Create();

  // Process optional <always_on> and <update_rate>.
  if (sdf_sensor.Element()->HasElement("always_on")) {
    INTR_ASSIGN_OR_RETURN(
        bool always_on, ParseChildAs<bool>(sdf_sensor.Element(), "always_on"));
    result.sensor_component->SetAlwaysOn(always_on);
  }

  if (sdf_sensor.Element()->HasElement("update_rate")) {
    INTR_ASSIGN_OR_RETURN(
        double update_rate,
        ParseChildAs<double>(sdf_sensor.Element(), "update_rate"));
    result.sensor_component->SetUpdateRate(update_rate);
  }

  if (sdf_sensor.Element()->HasElement("topic")) {
    INTR_ASSIGN_OR_RETURN(
        std::string topic,
        ParseChildAs<std::string>(sdf_sensor.Element(), "topic"));
    result.sensor_component->SetTopic(topic);
  }

  // Process "type" attribute.
  switch (type) {
    case ::sdf::SensorType::CAMERA: {
      if (sdf_sensor.CameraSensor() == nullptr) {
        return absl::InvalidArgumentError(
            absl::Substitute("SDF <sensor> has type camera but has no valid "
                             "Camerasdf_sensor. Full element:\n $0",
                             sdf_sensor.Element()->ToString("")));
      }
      intrinsic_proto::world::SensorComponent::Camera camera_proto;
      INTR_ASSIGN_OR_RETURN(*camera_proto.mutable_properties(),
                            IntrinsicCommonCameraPropertiesFromSdfCamera(
                                *sdf_sensor.CameraSensor()));
      result.sensor_component->SetCameraSpec(camera_proto);
    } break;
    case ::sdf::SensorType::DEPTH_CAMERA:
    case ::sdf::SensorType::RGBD_CAMERA: {
      if (sdf_sensor.CameraSensor() == nullptr) {
        return absl::InvalidArgumentError(
            absl::Substitute("SDF <sensor> has type depth camera or rgbd "
                             "camera but has no valid "
                             "Camerasdf_sensor. Full element:\n $0",
                             sdf_sensor.Element()->ToString("")));
      }

      intrinsic_proto::world::SensorComponent::DepthCamera depth_camera_proto;
      INTR_ASSIGN_OR_RETURN(*depth_camera_proto.mutable_properties(),
                            IntrinsicCommonCameraPropertiesFromSdfCamera(
                                *sdf_sensor.CameraSensor()));
      result.sensor_component->SetDepthCameraSpec(depth_camera_proto);
      if (sdf_sensor.CameraSensor()->HasDepthCamera()) {
        LOG(WARNING) << "<depth_camera> under <camera> is not supported.";
      }
    } break;
    case ::sdf::SensorType::FORCE_TORQUE: {
      if (sdf_sensor.ForceTorqueSensor() == nullptr) {
        return absl::InvalidArgumentError(absl::Substitute(
            "SDF <sensor> has type force torque but has no valid "
            "ForceTorquesdf_sensor. Full element:\n $0",
            sdf_sensor.Element()->ToString("")));
      }
      INTR_ASSIGN_OR_RETURN(auto force_torque_proto,
                            IntrinsicForceTorqueFromSdfForceTorque(
                                *sdf_sensor.ForceTorqueSensor()));
      result.sensor_component->SetForceTorqueSpec(force_torque_proto);
    } break;
    case ::sdf::SensorType::LIDAR: {
      if (sdf_sensor.LidarSensor() == nullptr) {
        return absl::InvalidArgumentError(
            absl::Substitute("SDF <sensor> has type lidar but has no valid "
                             "Lidarsdf_sensor. Full "
                             "element:\n $0",
                             sdf_sensor.Element()->ToString("")));
      }
      INTR_ASSIGN_OR_RETURN(auto lidar_proto, IntrinsicLidarFromSdfLidar(
                                                  *sdf_sensor.LidarSensor()));
      result.sensor_component->SetLidarSpec(lidar_proto);
    } break;
    case ::sdf::SensorType::NONE:
    case ::sdf::SensorType::ALTIMETER:
    case ::sdf::SensorType::CONTACT:
    case ::sdf::SensorType::GPS:
    case ::sdf::SensorType::GPU_LIDAR:
    case ::sdf::SensorType::IMU:
    case ::sdf::SensorType::LOGICAL_CAMERA:
    case ::sdf::SensorType::MAGNETOMETER:
    case ::sdf::SensorType::MULTICAMERA:
    case ::sdf::SensorType::RFID:
    case ::sdf::SensorType::RFIDTAG:
    case ::sdf::SensorType::SONAR:
    case ::sdf::SensorType::WIRELESS_RECEIVER:
    case ::sdf::SensorType::WIRELESS_TRANSMITTER:
    case ::sdf::SensorType::AIR_PRESSURE:
    case ::sdf::SensorType::THERMAL_CAMERA:
    case ::sdf::SensorType::NAVSAT:
    case ::sdf::SensorType::SEGMENTATION_CAMERA:
    case ::sdf::SensorType::BOUNDINGBOX_CAMERA:
    case ::sdf::SensorType::CUSTOM:
    case ::sdf::SensorType::WIDE_ANGLE_CAMERA:
    case ::sdf::SensorType::AIR_SPEED: {
      LOG(WARNING) << "Sensor type '" << sdf_sensor.TypeStr()
                   << "' is not supported by World.";
    }
  }

  return result;
}

}  // namespace sdf
}  // namespace intrinsic
