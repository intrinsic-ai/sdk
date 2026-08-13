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

#include "intrinsic/world/component/sensor_component.h"

#include <memory>
#include <string>
#include <utility>
#include <variant>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/util/status/status_builder.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/world/proto/sensor_component.pb.h"

namespace intrinsic {

namespace {

class SensorComponentImpl : public SensorComponent {
 public:
  explicit SensorComponentImpl() = default;
  explicit SensorComponentImpl(bool always_on, double update_rate,
                               absl::string_view topic)
      : always_on_(always_on),
        update_rate_(update_rate),
        topic_(std::string(topic)) {}

  absl::StatusOr<intrinsic_proto::world::SensorComponent> ToProto()
      const override;
  std::unique_ptr<SensorComponent> Clone() const override;
  bool GetAlwaysOn() const override { return always_on_; }
  void SetAlwaysOn(bool always_on) override { always_on_ = always_on; };
  double GetUpdateRate() const override { return update_rate_; }
  void SetUpdateRate(double update_rate) override {
    update_rate_ = update_rate;
  }
  const std::string& GetTopic() const override { return topic_; }
  void SetTopic(absl::string_view topic) override { topic_ = topic; }
  SensorType GetType() const override { return type_; }
  absl::StatusOr<intrinsic_proto::world::SensorComponent::Camera>
  GetCameraSpec() const override {
    if (type_ != SensorType::kInvalid &&
        std::holds_alternative<intrinsic_proto::world::SensorComponent::Camera>(
            spec_)) {
      return std::get<intrinsic_proto::world::SensorComponent::Camera>(spec_);
    }
    return intrinsic::NotFoundErrorBuilder()
           << "The type of the sensor is not camera.";
  }
  absl::StatusOr<intrinsic_proto::world::SensorComponent::DepthCamera>
  GetDepthCameraSpec() const override {
    if (type_ != SensorType::kInvalid &&
        std::holds_alternative<
            intrinsic_proto::world::SensorComponent::DepthCamera>(spec_)) {
      return std::get<intrinsic_proto::world::SensorComponent::DepthCamera>(
          spec_);
    }
    return intrinsic::NotFoundErrorBuilder()
           << "The type of the sensor is not depth camera.";
  }
  absl::StatusOr<intrinsic_proto::world::SensorComponent::ForceTorque>
  GetForceTorqueSpec() const override {
    if (type_ != SensorType::kInvalid &&
        std::holds_alternative<
            intrinsic_proto::world::SensorComponent::ForceTorque>(spec_)) {
      return std::get<intrinsic_proto::world::SensorComponent::ForceTorque>(
          spec_);
    }
    return intrinsic::NotFoundErrorBuilder()
           << "The type of the sensor is not force torque.";
  }
  absl::StatusOr<intrinsic_proto::world::SensorComponent::Lidar> GetLidarSpec()
      const override {
    if (type_ != SensorType::kInvalid &&
        std::holds_alternative<intrinsic_proto::world::SensorComponent::Lidar>(
            spec_)) {
      return std::get<intrinsic_proto::world::SensorComponent::Lidar>(spec_);
    }
    return intrinsic::NotFoundErrorBuilder()
           << "The type of the sensor is not lidar.";
  }
  void SetCameraSpec(
      const intrinsic_proto::world::SensorComponent::Camera& spec) override {
    spec_ = spec;
    type_ = SensorType::kCamera;
  }
  void SetDepthCameraSpec(
      const intrinsic_proto::world::SensorComponent::DepthCamera& spec)
      override {
    spec_ = spec;
    type_ = SensorType::kDepthCamera;
  }
  void SetForceTorqueSpec(
      const intrinsic_proto::world::SensorComponent::ForceTorque& spec)
      override {
    spec_ = spec;
    type_ = SensorType::kForceTorque;
  }
  void SetLidarSpec(
      const intrinsic_proto::world::SensorComponent::Lidar& spec) override {
    spec_ = spec;
    type_ = SensorType::kLidar;
  }
  absl::Status UpdateFromProto(
      const intrinsic_proto::world::SensorComponent& proto) override;

 private:
  bool always_on_ = false;
  double update_rate_ = 0;
  std::string topic_;
  SensorType type_ = SensorType::kInvalid;
  std::variant<intrinsic_proto::world::SensorComponent::Camera,
               intrinsic_proto::world::SensorComponent::DepthCamera,
               intrinsic_proto::world::SensorComponent::ForceTorque,
               intrinsic_proto::world::SensorComponent::Lidar>
      spec_;
};

absl::StatusOr<intrinsic_proto::world::SensorComponent>
SensorComponentImpl::ToProto() const {
  intrinsic_proto::world::SensorComponent result;
  result.set_always_on(always_on_);
  result.set_update_rate(update_rate_);
  result.set_topic(topic_);
  switch (type_) {
    case SensorType::kCamera:
      *result.mutable_camera() =
          std::get<intrinsic_proto::world::SensorComponent::Camera>(spec_);
      break;
    case SensorType::kDepthCamera:
      *result.mutable_depth_camera() =
          std::get<intrinsic_proto::world::SensorComponent::DepthCamera>(spec_);
      break;
    case SensorType::kForceTorque:
      *result.mutable_force_torque() =
          std::get<intrinsic_proto::world::SensorComponent::ForceTorque>(spec_);
      break;
    case SensorType::kLidar:
      *result.mutable_lidar() =
          std::get<intrinsic_proto::world::SensorComponent::Lidar>(spec_);
      break;
    case SensorType::kInvalid:
      break;
  }

  return std::move(result);
}

std::unique_ptr<SensorComponent> SensorComponentImpl::Clone() const {
  auto result = std::make_unique<SensorComponentImpl>(
      always_on_, update_rate_,
      topic_);
  result->type_ = type_;
  result->spec_ = spec_;
  return result;
}

absl::Status SensorComponentImpl::UpdateFromProto(
    const intrinsic_proto::world::SensorComponent& proto) {
  INTR_ASSIGN_OR_RETURN(auto other_ptr, FromProto(proto));
  auto* other = dynamic_cast<SensorComponentImpl*>(other_ptr.get());

  always_on_ = other->always_on_;
  update_rate_ = other->update_rate_;
  type_ = other->type_;
  spec_ = other->spec_;
  topic_ = other->topic_;

  return absl::OkStatus();
}

}  // namespace

std::unique_ptr<SensorComponent> SensorComponent::Create() {
  return std::make_unique<SensorComponentImpl>();
}

absl::StatusOr<std::unique_ptr<SensorComponent>> SensorComponent::FromProto(
    const intrinsic_proto::world::SensorComponent& proto) {
  std::string device_name = "";
  std::string device_type = "";

  auto result = std::make_unique<SensorComponentImpl>(
      proto.always_on(), proto.update_rate(),
      proto.topic());

  switch (proto.type_oneof_case()) {
    case intrinsic_proto::world::SensorComponent::kCamera:
      result->SetCameraSpec(proto.camera());
      break;
    case intrinsic_proto::world::SensorComponent::kDepthCamera:
      result->SetDepthCameraSpec(proto.depth_camera());
      break;
    case intrinsic_proto::world::SensorComponent::kForceTorque:
      result->SetForceTorqueSpec(proto.force_torque());
      break;
    case intrinsic_proto::world::SensorComponent::kLidar:
      result->SetLidarSpec(proto.lidar());
      break;
    case intrinsic_proto::world::SensorComponent::TYPE_ONEOF_NOT_SET:
      break;
  }
  return result;
}

}  // namespace intrinsic
