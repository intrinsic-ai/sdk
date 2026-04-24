// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/lib/robot_arm/robot_arm_bus_component.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <algorithm>
#include <array>
#include <cstddef>
#include <cstdint>
#include <memory>
#include <optional>
#include <string>

#include "absl/status/status.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/hal/hardware_interface_handle.h"
#include "intrinsic/icon/hal/hardware_interface_registry.h"
#include "intrinsic/icon/hal/icon_state_register.h"
#include "intrinsic/icon/hal/interfaces/icon_state.fbs.h"
#include "intrinsic/icon/hal/interfaces/joint_command.fbs.h"
#include "intrinsic/icon/hal/interfaces/joint_state.fbs.h"
#include "intrinsic/icon/hal/lib/fieldbus/async_device_request.h"
#include "intrinsic/icon/hal/lib/fieldbus/device_init_context.h"
#include "intrinsic/icon/hal/lib/fieldbus/fake_variable_registry.h"
#include "intrinsic/icon/hal/lib/robot_arm/v1/robot_arm_bus_component_config.pb.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/shared_memory_manager.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/testing/unique_segment_name.h"
#include "intrinsic/icon/utils/clock.h"
#include "intrinsic/icon/utils/current_cycle.h"
#include "intrinsic/icon/utils/realtime_status_matchers.h"
#include "intrinsic/util/proto/parse_text_proto.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

namespace intrinsic::robot_arm {
namespace {

using ::absl_testing::StatusIs;
using ::intrinsic::ParseTextProtoOrDie;
using ::intrinsic::fieldbus::RequestStatus;
using ::intrinsic::fieldbus::RequestType;
using ::intrinsic::icon::HardwareInterfaceRegistry;
using ::intrinsic::icon::RealtimeIsOkAndHolds;
using ::intrinsic_proto::icon::v1::RobotArmBusComponentConfig;
using ::testing::DoubleNear;
using ::testing::Eq;
using ::testing::HasSubstr;
using ::testing::Pointwise;

template <typename T, std::size_t kNDof>
class LoopbackFakeVariableRegistry
    : public intrinsic::fieldbus::FakeVariableRegistry {
 public:
  explicit LoopbackFakeVariableRegistry(
      const std::array<T, kNDof>& joint_states,
      const std::array<T, kNDof>& position_commands,
      const std::array<T, kNDof>& ff_velocity_commands)
      : joint_states_(joint_states),
        position_commands_(position_commands),
        ff_velocity_commands_(ff_velocity_commands) {
    for (std::size_t i = 0; i < kNDof; ++i) {
      // We map position and velocity states to the same data.
      AddInputProcessVariable(absl::StrCat("position_state_", i),
                              &joint_states_[i]);
      AddInputProcessVariable(absl::StrCat("velocity_state_", i),
                              &joint_states_[i]);
      AddInputProcessVariable(absl::StrCat("acceleration_state_", i),
                              &joint_states_[i]);
      AddOutputProcessVariable(absl::StrCat("position_command_", i),
                               &position_commands_[i]);
      AddOutputProcessVariable(absl::StrCat("ff_velocity_command_", i),
                               &ff_velocity_commands_[i]);
    }
  }

  T GetJointState(std::size_t i) const { return joint_states_[i]; }

  std::array<T, kNDof> GetJointStates() const { return joint_states_; }

  T GetPositionCommand(std::size_t i) const { return position_commands_[i]; }

  std::array<T, kNDof> GetPositionCommands() const {
    return position_commands_;
  }

  std::array<T, kNDof> GetFeedforwardVelocityCommands() const {
    return ff_velocity_commands_;
  }

 private:
  std::array<T, kNDof> joint_states_;
  std::array<T, kNDof> position_commands_;
  std::array<T, kNDof> ff_velocity_commands_;
};

constexpr absl::string_view kTypical3DofDeviceConfig = R"pb(
  device_prefix: "foo"
  position_state_variables:
  [ { variable_name: "position_state_0", scale: 2.0 }
    , { variable_name: "position_state_1", scale: 2.0 }
    , { variable_name: "position_state_2", scale: 2.0 }]
  velocity_state_variables:
  [ { variable_name: "velocity_state_0", scale: 2.0 }
    , { variable_name: "velocity_state_1", scale: 2.0 }
    , { variable_name: "velocity_state_2", scale: 2.0 }]
  acceleration_state_variables:
  [ { variable_name: "acceleration_state_0", scale: 2.0 }
    , { variable_name: "acceleration_state_1", scale: 2.0 }
    , { variable_name: "acceleration_state_2", scale: 2.0 }]
  position_command_variables:
  [ { variable_name: "position_command_0", scale: 0.5 }
    , { variable_name: "position_command_1", scale: 0.5 }
    , { variable_name: "position_command_2", scale: 0.5 }]
  feedforward_velocity_command_variables:
  [ { variable_name: "ff_velocity_command_0", scale: 0.5 }
    , { variable_name: "ff_velocity_command_1", scale: 0.5 }
    , { variable_name: "ff_velocity_command_2", scale: 0.5 }]
)pb";

constexpr absl::string_view kNegative3DofDeviceConfig = R"pb(
  device_prefix: "foo"
  position_state_variables:
  [ { variable_name: "position_state_0", scale: -2.0 }
    , { variable_name: "position_state_1", scale: -2.0 }
    , { variable_name: "position_state_2", scale: -2.0 }]
  velocity_state_variables:
  [ { variable_name: "velocity_state_0", scale: -2.0 }
    , { variable_name: "velocity_state_1", scale: -2.0 }
    , { variable_name: "velocity_state_2", scale: -2.0 }]
  acceleration_state_variables:
  [ { variable_name: "acceleration_state_0", scale: -2.0 }
    , { variable_name: "acceleration_state_1", scale: -2.0 }
    , { variable_name: "acceleration_state_2", scale: -2.0 }]
  position_command_variables:
  [ { variable_name: "position_command_0", scale: -0.5 }
    , { variable_name: "position_command_1", scale: -0.5 }
    , { variable_name: "position_command_2", scale: -0.5 }]
  feedforward_velocity_command_variables:
  [ { variable_name: "ff_velocity_command_0", scale: -0.5 }
    , { variable_name: "ff_velocity_command_1", scale: -0.5 }
    , { variable_name: "ff_velocity_command_2", scale: -0.5 }]
)pb";

class RobotArmBusComponentTestFixture : public ::testing::Test {
 public:
  void SetUp() override {
    ASSERT_OK_AND_ASSIGN(shm_manager_,
                         intrinsic::icon::SharedMemoryManager::Create(
                             kMemoryNamespace, "some_module"));
    interface_registry_ =
        intrinsic::icon::HardwareInterfaceRegistry(*shm_manager_);

    ASSERT_OK_AND_ASSIGN(
        icon_state_, interface_registry_
                         ->AdvertiseMutableInterface<intrinsic_fbs::IconState>(
                             intrinsic::icon::kIconStateInterfaceName));
  }

 protected:
  const std::string kMemoryNamespace = intrinsic::icon::UniqueMemoryNamespace();
  std::unique_ptr<intrinsic::icon::SharedMemoryManager> shm_manager_;
  std::optional<HardwareInterfaceRegistry> interface_registry_;
  intrinsic::icon::MutableHardwareInterfaceHandle<intrinsic_fbs::IconState>
      icon_state_;
};

TEST_F(RobotArmBusComponentTestFixture, CreateWhenTypical) {
  LoopbackFakeVariableRegistry<double, 3> variable_registry =
      LoopbackFakeVariableRegistry<double, 3>(
          /*joint_states=*/{2, 4, 8}, /*position_commands=*/{0, 0, 0},
          /*ff_velocity_commands=*/{0, 0, 0});

  fieldbus::DeviceInitContext init_context(*interface_registry_,
                                           variable_registry);
  ASSERT_OK_AND_ASSIGN(
      std::unique_ptr<RobotArmBusComponent> robot_arm_bus_component,
      RobotArmBusComponent::Create(
          init_context, ParseTextProtoOrDie(kTypical3DofDeviceConfig)));
}

TEST_F(RobotArmBusComponentTestFixture, CreateFailsWithUnsupportedBusDataType) {
  LoopbackFakeVariableRegistry<uint8_t, 3> variable_registry =
      LoopbackFakeVariableRegistry<uint8_t, 3>(
          /*joint_states=*/{2, 4, 8}, /*position_commands=*/{0, 0, 0},
          /*ff_velocity_commands=*/{0, 0, 0});

  fieldbus::DeviceInitContext init_context(*interface_registry_,
                                           variable_registry);
  EXPECT_THAT(RobotArmBusComponent::Create(
                  init_context, ParseTextProtoOrDie(kTypical3DofDeviceConfig)),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("Unsupported bus variable data type")));
}

TEST_F(RobotArmBusComponentTestFixture, CreateFailsWithMissingPositionStates) {
  LoopbackFakeVariableRegistry<double, 3> variable_registry =
      LoopbackFakeVariableRegistry<double, 3>(
          /*joint_states=*/{2, 4, 8}, /*position_commands=*/{0, 0, 0},
          /*ff_velocity_commands=*/{0, 0, 0});

  RobotArmBusComponentConfig device_config =
      ParseTextProtoOrDie(kTypical3DofDeviceConfig);
  // Remove the state variables, so that the number of joints mismatch.
  device_config.mutable_position_state_variables()->Clear();
  fieldbus::DeviceInitContext init_context(*interface_registry_,
                                           variable_registry);
  EXPECT_THAT(RobotArmBusComponent::Create(init_context, device_config),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("Number of joints mismatch")));
}

TEST_F(RobotArmBusComponentTestFixture,
       CreateSucceedsWithMissingVelocityStates) {
  LoopbackFakeVariableRegistry<double, 3> variable_registry =
      LoopbackFakeVariableRegistry<double, 3>(
          /*joint_states=*/{2, 4, 8}, /*position_commands=*/{0, 0, 0},
          /*ff_velocity_commands=*/{0, 0, 0});

  RobotArmBusComponentConfig device_config =
      ParseTextProtoOrDie(kTypical3DofDeviceConfig);
  // Remove the state variables, so that the number of joints mismatch.
  device_config.mutable_velocity_state_variables()->Clear();
  fieldbus::DeviceInitContext init_context(*interface_registry_,
                                           variable_registry);
  EXPECT_OK(RobotArmBusComponent::Create(init_context, device_config));
}

TEST_F(RobotArmBusComponentTestFixture, CreateFailsWithIllsizedVelocityStates) {
  LoopbackFakeVariableRegistry<double, 3> variable_registry =
      LoopbackFakeVariableRegistry<double, 3>(
          /*joint_states=*/{2, 4, 8}, /*position_commands=*/{0, 0, 0},
          /*ff_velocity_commands=*/{0, 0, 0});

  RobotArmBusComponentConfig device_config =
      ParseTextProtoOrDie(kTypical3DofDeviceConfig);
  // Remove the state variables, so that the number of joints mismatch.
  device_config.mutable_velocity_state_variables()->DeleteSubrange(0, 1);
  fieldbus::DeviceInitContext init_context(*interface_registry_,
                                           variable_registry);
  EXPECT_THAT(RobotArmBusComponent::Create(init_context, device_config),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("Number of joints mismatch")));
}

TEST_F(RobotArmBusComponentTestFixture, CreateFailsWithMissingPositionCommand) {
  LoopbackFakeVariableRegistry<double, 3> variable_registry =
      LoopbackFakeVariableRegistry<double, 3>(
          /*joint_states=*/{2, 4, 8}, /*position_commands=*/{0, 0, 0},
          /*ff_velocity_commands=*/{0, 0, 0});

  RobotArmBusComponentConfig device_config =
      ParseTextProtoOrDie(kTypical3DofDeviceConfig);
  // Remove the command variables, so that the number of joints mismatch.
  device_config.mutable_position_command_variables()->Clear();
  fieldbus::DeviceInitContext init_context(*interface_registry_,
                                           variable_registry);
  EXPECT_THAT(RobotArmBusComponent::Create(init_context, device_config),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("Number of joints mismatch")));
}

TEST_F(RobotArmBusComponentTestFixture,
       CreateSucceedsWithMissingFfVelocityStates) {
  LoopbackFakeVariableRegistry<double, 3> variable_registry =
      LoopbackFakeVariableRegistry<double, 3>(
          /*joint_states=*/{2, 4, 8}, /*position_commands=*/{0, 0, 0},
          /*ff_velocity_commands=*/{0, 0, 0});

  RobotArmBusComponentConfig device_config =
      ParseTextProtoOrDie(kTypical3DofDeviceConfig);
  // Remove the command variables, so that there's a config without feedforward
  // velocity commands.
  device_config.mutable_feedforward_velocity_command_variables()->Clear();
  fieldbus::DeviceInitContext init_context(*interface_registry_,
                                           variable_registry);
  EXPECT_OK(RobotArmBusComponent::Create(init_context, device_config));
}

TEST_F(RobotArmBusComponentTestFixture, ReadProvidesScaledJointStates) {
  std::array<double, 3> joint_states = {2, 4, 8};
  LoopbackFakeVariableRegistry<double, 3> variable_registry =
      LoopbackFakeVariableRegistry<double, 3>(
          /*joint_states=*/joint_states, /*position_commands=*/{0, 0, 0},
          /*ff_velocity_commands=*/{0, 0, 0});

  fieldbus::DeviceInitContext init_context(*interface_registry_,
                                           variable_registry);
  ASSERT_OK_AND_ASSIGN(
      std::unique_ptr<RobotArmBusComponent> robot_arm_bus_component,
      RobotArmBusComponent::Create(
          init_context, ParseTextProtoOrDie(kTypical3DofDeviceConfig)));

  EXPECT_THAT(
      robot_arm_bus_component->CyclicRead(RequestType::kNormalOperation),
      RealtimeIsOkAndHolds(RequestStatus::kDone));

  std::array<double, 3> expected_joint_values;
  // Scale `position_values` according to the scale factor in
  // `kTypical3DofDeviceConfig`.
  std::transform(std::begin(joint_states), std::end(joint_states),
                 std::begin(expected_joint_values),
                 [](auto value) { return 2.0 * value; });

  ASSERT_OK_AND_ASSIGN(
      auto position_state,
      interface_registry_
          ->GetInterfaceHandle<intrinsic_fbs::JointPositionState>(
              "foo_joint_position_state"));

  EXPECT_THAT(*position_state->position(),
              Pointwise(DoubleNear(0.01), expected_joint_values));

  // Read also writes the sensed position states as position commands
  // (b/297333030).
  EXPECT_THAT(variable_registry.GetPositionCommands(),
              Pointwise(DoubleNear(0.01), joint_states));

  ASSERT_OK_AND_ASSIGN(
      auto velocity_state,
      interface_registry_
          ->GetInterfaceHandle<intrinsic_fbs::JointVelocityState>(
              "foo_joint_velocity_state"));
  EXPECT_THAT(*velocity_state->velocity(),
              Pointwise(DoubleNear(0.01), expected_joint_values));

  // Acceleration is same same as velocity, since they're mapped to the same
  // fake variables.
  ASSERT_OK_AND_ASSIGN(
      auto acceleration_state,
      interface_registry_
          ->GetInterfaceHandle<intrinsic_fbs::JointAccelerationState>(
              "foo_joint_acceleration_state"));
  EXPECT_THAT(*acceleration_state->acceleration(),
              Pointwise(DoubleNear(0.01), expected_joint_values));
}

TEST_F(RobotArmBusComponentTestFixture,
       ReadProvidesScaledJointStatesWithNegativeScale) {
  std::array<double, 3> joint_states = {2, 4, 8};
  LoopbackFakeVariableRegistry<double, 3> variable_registry =
      LoopbackFakeVariableRegistry<double, 3>(
          /*joint_states=*/joint_states, /*position_commands=*/{0, 0, 0},
          /*ff_velocity_commands=*/{0, 0, 0});

  fieldbus::DeviceInitContext init_context(*interface_registry_,
                                           variable_registry);
  ASSERT_OK_AND_ASSIGN(
      std::unique_ptr<RobotArmBusComponent> robot_arm_bus_component,
      RobotArmBusComponent::Create(
          init_context, ParseTextProtoOrDie(kNegative3DofDeviceConfig)));

  EXPECT_THAT(
      robot_arm_bus_component->CyclicRead(RequestType::kNormalOperation),
      RealtimeIsOkAndHolds(RequestStatus::kDone));

  std::array<double, 3> expected_joint_values;
  // Scale `position_values` according to the scale factor in
  // `kNegative3DofDeviceConfig`.
  std::transform(std::begin(joint_states), std::end(joint_states),
                 std::begin(expected_joint_values),
                 [](auto value) { return -2.0 * value; });

  ASSERT_OK_AND_ASSIGN(
      auto position_state,
      interface_registry_
          ->GetInterfaceHandle<intrinsic_fbs::JointPositionState>(
              "foo_joint_position_state"));

  EXPECT_THAT(*position_state->position(),
              Pointwise(DoubleNear(0.01), expected_joint_values));

  // Read also writes the sensed position states as position commands
  // (b/297333030).
  EXPECT_THAT(variable_registry.GetPositionCommands(),
              Pointwise(DoubleNear(0.01), joint_states));

  ASSERT_OK_AND_ASSIGN(
      auto velocity_state,
      interface_registry_
          ->GetInterfaceHandle<intrinsic_fbs::JointVelocityState>(
              "foo_joint_velocity_state"));
  EXPECT_THAT(*velocity_state->velocity(),
              Pointwise(DoubleNear(0.01), expected_joint_values));

  // Acceleration is same same as velocity, since they're mapped to the same
  // fake variables.
  ASSERT_OK_AND_ASSIGN(
      auto acceleration_state,
      interface_registry_
          ->GetInterfaceHandle<intrinsic_fbs::JointAccelerationState>(
              "foo_joint_acceleration_state"));
  EXPECT_THAT(*acceleration_state->acceleration(),
              Pointwise(DoubleNear(0.01), expected_joint_values));
}

TEST_F(
    RobotArmBusComponentTestFixture,
    ReadWithoutVelocityVariablesProvidesJointPositionAndAccelerationStateOnly) {
  std::array<double, 3> joint_states = {2, 4, 8};
  LoopbackFakeVariableRegistry<double, 3> variable_registry =
      LoopbackFakeVariableRegistry<double, 3>(
          /*joint_states=*/joint_states, /*position_commands=*/{0, 0, 0},
          /*ff_velocity_commands=*/{0, 0, 0});

  RobotArmBusComponentConfig device_config =
      ParseTextProtoOrDie(kTypical3DofDeviceConfig);
  // Remove the velocity state variables, so that the velocity interface is not
  // advertised.
  device_config.mutable_velocity_state_variables()->Clear();
  fieldbus::DeviceInitContext init_context(*interface_registry_,
                                           variable_registry);
  ASSERT_OK_AND_ASSIGN(
      std::unique_ptr<RobotArmBusComponent> robot_arm_bus_component,
      RobotArmBusComponent::Create(init_context, device_config));

  EXPECT_THAT(
      robot_arm_bus_component->CyclicRead(RequestType::kNormalOperation),
      RealtimeIsOkAndHolds(RequestStatus::kDone));

  std::array<double, 3> expected_joint_values;
  // Scale `position_values` according to the scale factor in
  // `kTypical3DofDeviceConfig`.
  std::transform(std::begin(joint_states), std::end(joint_states),
                 std::begin(expected_joint_values),
                 [](auto value) { return 2.0 * value; });

  ASSERT_OK_AND_ASSIGN(
      auto position_state,
      interface_registry_
          ->GetInterfaceHandle<intrinsic_fbs::JointPositionState>(
              "foo_joint_position_state"));

  EXPECT_THAT(*position_state->position(),
              Pointwise(DoubleNear(0.01), expected_joint_values));

  // Read also writes the sensed position states as position commands
  // (b/297333030).
  EXPECT_THAT(variable_registry.GetPositionCommands(),
              Pointwise(DoubleNear(0.01), joint_states));

  EXPECT_THAT(interface_registry_
                  ->GetInterfaceHandle<intrinsic_fbs::JointVelocityState>(
                      "foo_joint_velocity_state"),
              StatusIs(absl::StatusCode::kNotFound));

  // Acceleration is same same as position, since they're mapped to the same
  // fake variables.
  ASSERT_OK_AND_ASSIGN(
      auto acceleration_state,
      interface_registry_
          ->GetInterfaceHandle<intrinsic_fbs::JointAccelerationState>(
              "foo_joint_acceleration_state"));
  EXPECT_THAT(*acceleration_state->acceleration(),
              Pointwise(DoubleNear(0.01), expected_joint_values));
}

TEST_F(
    RobotArmBusComponentTestFixture,
    ReadWithoutAccelerationVariablesProvidesJointPositionAndVelocityStateOnly) {
  std::array<double, 3> joint_states = {2, 4, 8};
  LoopbackFakeVariableRegistry<double, 3> variable_registry =
      LoopbackFakeVariableRegistry<double, 3>(
          /*joint_states=*/joint_states, /*position_commands=*/{0, 0, 0},
          /*ff_velocity_commands=*/{0, 0, 0});

  RobotArmBusComponentConfig device_config =
      ParseTextProtoOrDie(kTypical3DofDeviceConfig);
  // Remove the acceleration state variables, so that the acceleration interface
  // is not advertised.
  device_config.mutable_acceleration_state_variables()->Clear();
  fieldbus::DeviceInitContext init_context(*interface_registry_,
                                           variable_registry);
  ASSERT_OK_AND_ASSIGN(
      std::unique_ptr<RobotArmBusComponent> robot_arm_bus_component,
      RobotArmBusComponent::Create(init_context, device_config));

  EXPECT_THAT(
      robot_arm_bus_component->CyclicRead(RequestType::kNormalOperation),
      RealtimeIsOkAndHolds(RequestStatus::kDone));

  std::array<double, 3> expected_joint_values;
  // Scale `position_values` according to the scale factor in
  // `kTypical3DofDeviceConfig`.
  std::transform(std::begin(joint_states), std::end(joint_states),
                 std::begin(expected_joint_values),
                 [](auto value) { return 2.0 * value; });

  ASSERT_OK_AND_ASSIGN(
      auto position_state,
      interface_registry_
          ->GetInterfaceHandle<intrinsic_fbs::JointPositionState>(
              "foo_joint_position_state"));

  EXPECT_THAT(*position_state->position(),
              Pointwise(DoubleNear(0.01), expected_joint_values));

  // Read also writes the sensed position states as position commands
  // (b/297333030).
  EXPECT_THAT(variable_registry.GetPositionCommands(),
              Pointwise(DoubleNear(0.01), joint_states));

  // Velocity is same same as position, since they're mapped to the same
  // fake variables.
  ASSERT_OK_AND_ASSIGN(
      auto velocity_state,
      interface_registry_
          ->GetInterfaceHandle<intrinsic_fbs::JointVelocityState>(
              "foo_joint_velocity_state"));
  EXPECT_THAT(*velocity_state->velocity(),
              Pointwise(DoubleNear(0.01), expected_joint_values));

  EXPECT_THAT(interface_registry_
                  ->GetInterfaceHandle<intrinsic_fbs::JointAccelerationState>(
                      "foo_joint_acceleration_state"),
              StatusIs(absl::StatusCode::kNotFound));
}

TEST_F(
    RobotArmBusComponentTestFixture,
    ReadWithoutVelocityAndAccelerationVariablesProvidesJointPositionStateOnly) {
  std::array<double, 3> joint_states = {2, 4, 8};
  LoopbackFakeVariableRegistry<double, 3> variable_registry =
      LoopbackFakeVariableRegistry<double, 3>(
          /*joint_states=*/joint_states, /*position_commands=*/{0, 0, 0},
          /*ff_velocity_commands=*/{0, 0, 0});

  RobotArmBusComponentConfig device_config =
      ParseTextProtoOrDie(kTypical3DofDeviceConfig);
  // Remove the velocity and acceleration state variables, so that both
  // interfaces are not advertised.
  device_config.mutable_velocity_state_variables()->Clear();
  device_config.mutable_acceleration_state_variables()->Clear();
  fieldbus::DeviceInitContext init_context(*interface_registry_,
                                           variable_registry);
  ASSERT_OK_AND_ASSIGN(
      std::unique_ptr<RobotArmBusComponent> robot_arm_bus_component,
      RobotArmBusComponent::Create(init_context, device_config));

  EXPECT_THAT(
      robot_arm_bus_component->CyclicRead(RequestType::kNormalOperation),
      RealtimeIsOkAndHolds(RequestStatus::kDone));

  std::array<double, 3> expected_joint_values;
  // Scale `position_values` according to the scale factor in
  // `kTypical3DofDeviceConfig`.
  std::transform(std::begin(joint_states), std::end(joint_states),
                 std::begin(expected_joint_values),
                 [](auto value) { return 2.0 * value; });

  ASSERT_OK_AND_ASSIGN(
      auto position_state,
      interface_registry_
          ->GetInterfaceHandle<intrinsic_fbs::JointPositionState>(
              "foo_joint_position_state"));

  EXPECT_THAT(*position_state->position(),
              Pointwise(DoubleNear(0.01), expected_joint_values));

  // Read also writes the sensed position states as position commands
  // (b/297333030).
  EXPECT_THAT(variable_registry.GetPositionCommands(),
              Pointwise(DoubleNear(0.01), joint_states));

  EXPECT_THAT(interface_registry_
                  ->GetInterfaceHandle<intrinsic_fbs::JointVelocityState>(
                      "foo_joint_velocity_state"),
              StatusIs(absl::StatusCode::kNotFound));

  EXPECT_THAT(interface_registry_
                  ->GetInterfaceHandle<intrinsic_fbs::JointAccelerationState>(
                      "foo_joint_acceleration_state"),
              StatusIs(absl::StatusCode::kNotFound));
}

TEST_F(RobotArmBusComponentTestFixture,
       ReadProvidesScaledJointStatesWithIntegralType) {
  std::array<int32_t, 3> joint_states = {2, 4, 8};
  LoopbackFakeVariableRegistry<int32_t, 3> variable_registry =
      LoopbackFakeVariableRegistry<int32_t, 3>(
          /*joint_states=*/joint_states, /*position_commands=*/{0, 0, 0},
          /*ff_velocity_commands=*/{0, 0, 0});

  fieldbus::DeviceInitContext init_context(*interface_registry_,
                                           variable_registry);
  ASSERT_OK_AND_ASSIGN(
      std::unique_ptr<RobotArmBusComponent> robot_arm_bus_component,
      RobotArmBusComponent::Create(
          init_context, ParseTextProtoOrDie(kTypical3DofDeviceConfig)));

  EXPECT_THAT(
      robot_arm_bus_component->CyclicRead(RequestType::kNormalOperation),
      RealtimeIsOkAndHolds(RequestStatus::kDone));

  std::array<int32_t, 3> expected_joint_values;
  // Scale `position_values` according to the scale factor in
  // `kTypical3DofDeviceConfig`.
  std::transform(std::begin(joint_states), std::end(joint_states),
                 std::begin(expected_joint_values),
                 [](auto value) { return 2.0 * value; });

  ASSERT_OK_AND_ASSIGN(
      auto position_state,
      interface_registry_
          ->GetInterfaceHandle<intrinsic_fbs::JointPositionState>(
              "foo_joint_position_state"));

  EXPECT_THAT(*position_state->position(),
              Pointwise(Eq(), expected_joint_values));

  // Read also writes the sensed position states as position commands
  // (b/297333030).
  EXPECT_THAT(variable_registry.GetPositionCommands(),
              Pointwise(Eq(), joint_states));

  ASSERT_OK_AND_ASSIGN(
      auto velocity_state,
      interface_registry_
          ->GetInterfaceHandle<intrinsic_fbs::JointVelocityState>(
              "foo_joint_velocity_state"));
  EXPECT_THAT(*velocity_state->velocity(),
              Pointwise(Eq(), expected_joint_values));

  ASSERT_OK_AND_ASSIGN(
      auto acceleration_state,
      interface_registry_
          ->GetInterfaceHandle<intrinsic_fbs::JointAccelerationState>(
              "foo_joint_acceleration_state"));
  EXPECT_THAT(*acceleration_state->acceleration(),
              Pointwise(Eq(), expected_joint_values));
}

TEST_F(RobotArmBusComponentTestFixture,
       WriteProvidesScaledJointPositionAndFfVelocityCommands) {
  LoopbackFakeVariableRegistry<double, 3> variable_registry =
      LoopbackFakeVariableRegistry<double, 3>(
          /*joint_states=*/{0, 0, 0},
          /*position_commands=*/{0, 0, 0},
          /*ff_velocity_commands=*/{0, 0, 0});

  fieldbus::DeviceInitContext init_context(*interface_registry_,
                                           variable_registry);
  ASSERT_OK_AND_ASSIGN(
      std::unique_ptr<RobotArmBusComponent> robot_arm_bus_component,
      RobotArmBusComponent::Create(
          init_context, ParseTextProtoOrDie(kTypical3DofDeviceConfig)));

  ASSERT_OK_AND_ASSIGN(
      auto position_command,
      interface_registry_
          ->GetMutableInterfaceHandle<intrinsic_fbs::JointPositionCommand>(
              "foo_joint_position_command"));

  position_command->mutable_position()->Mutate(0, 2);
  position_command->mutable_position()->Mutate(1, 4);
  position_command->mutable_position()->Mutate(2, 8);
  position_command->mutable_velocity_feedforward()->Mutate(0, 16);
  position_command->mutable_velocity_feedforward()->Mutate(1, 32);
  position_command->mutable_velocity_feedforward()->Mutate(2, 64);

  // Not simulating ICON tick.
  EXPECT_THAT(
      robot_arm_bus_component->CyclicWrite(RequestType::kNormalOperation),
      RealtimeIsOkAndHolds(RequestStatus::kDone));

  // Expect the sensed position {0, 0, 0} if the position command update
  // timestamp is not updated.
  std::array<double, 3> expected_position_values = {0, 0, 0};
  EXPECT_THAT(variable_registry.GetPositionCommands(),
              Pointwise(DoubleNear(0.01), expected_position_values));
  // Expect the velocity feedforward {0, 0, 0} if the position command update
  // timestamp is not updated.
  std::array<double, 3> expected_feedforward_velocity_values = {0, 0, 0};
  EXPECT_THAT(
      variable_registry.GetFeedforwardVelocityCommands(),
      Pointwise(DoubleNear(0.01), expected_feedforward_velocity_values));

  // Try again with updated timestamp.
  {
    // Simulates an ICON tick.
    intrinsic::icon::Cycle::SetCurrentCycle(42);
    icon_state_->mutate_current_cycle(42);
    icon_state_.UpdatedAt(intrinsic::Clock::now());
    position_command.UpdatedAt(intrinsic::Clock::Now());
  }
  EXPECT_THAT(
      robot_arm_bus_component->CyclicWrite(RequestType::kNormalOperation),
      RealtimeIsOkAndHolds(RequestStatus::kDone));

  // Now we expect the commanded position values.
  expected_position_values = {1, 2, 4};
  EXPECT_THAT(variable_registry.GetPositionCommands(),
              Pointwise(DoubleNear(0.01), expected_position_values));
  // Now we expect the commanded feedforward velocity values.
  expected_feedforward_velocity_values = {8, 16, 32};
  EXPECT_THAT(
      variable_registry.GetFeedforwardVelocityCommands(),
      Pointwise(DoubleNear(0.01), expected_feedforward_velocity_values));
}

TEST_F(RobotArmBusComponentTestFixture,
       WriteProvidesScaledJointPositionAndFfVelocityCommandsWithNegativeScale) {
  LoopbackFakeVariableRegistry<double, 3> variable_registry =
      LoopbackFakeVariableRegistry<double, 3>(
          /*joint_states=*/{0, 0, 0},
          /*position_commands=*/{0, 0, 0},
          /*ff_velocity_commands=*/{0, 0, 0});

  fieldbus::DeviceInitContext init_context(*interface_registry_,
                                           variable_registry);
  ASSERT_OK_AND_ASSIGN(
      std::unique_ptr<RobotArmBusComponent> robot_arm_bus_component,
      RobotArmBusComponent::Create(
          init_context, ParseTextProtoOrDie(kNegative3DofDeviceConfig)));

  ASSERT_OK_AND_ASSIGN(
      auto position_command,
      interface_registry_
          ->GetMutableInterfaceHandle<intrinsic_fbs::JointPositionCommand>(
              "foo_joint_position_command"));

  position_command->mutable_position()->Mutate(0, 2);
  position_command->mutable_position()->Mutate(1, 4);
  position_command->mutable_position()->Mutate(2, 8);
  position_command->mutable_velocity_feedforward()->Mutate(0, 16);
  position_command->mutable_velocity_feedforward()->Mutate(1, 32);
  position_command->mutable_velocity_feedforward()->Mutate(2, 64);

  // Not simulating ICON tick.
  EXPECT_THAT(
      robot_arm_bus_component->CyclicWrite(RequestType::kNormalOperation),
      RealtimeIsOkAndHolds(RequestStatus::kDone));

  // Expect the sensed position {0, 0, 0} if the position command update
  // timestamp is not updated.
  std::array<double, 3> expected_position_values = {0, 0, 0};
  EXPECT_THAT(variable_registry.GetPositionCommands(),
              Pointwise(DoubleNear(0.01), expected_position_values));
  // Expect the velocity feedforward {0, 0, 0} if the position command update
  // timestamp is not updated.
  std::array<double, 3> expected_feedforward_velocity_values = {0, 0, 0};
  EXPECT_THAT(
      variable_registry.GetFeedforwardVelocityCommands(),
      Pointwise(DoubleNear(0.01), expected_feedforward_velocity_values));

  // Try again with updated timestamp.
  {
    // Simulates an ICON tick.
    intrinsic::icon::Cycle::SetCurrentCycle(42);
    icon_state_->mutate_current_cycle(42);
    icon_state_.UpdatedAt(intrinsic::Clock::now());
    position_command.UpdatedAt(intrinsic::Clock::Now());
  }
  EXPECT_THAT(
      robot_arm_bus_component->CyclicWrite(RequestType::kNormalOperation),
      RealtimeIsOkAndHolds(RequestStatus::kDone));

  // Now we expect the commanded position values.
  // Scale factor in `kNegative3DofDeviceConfig` is -0.5.
  expected_position_values = {-1, -2, -4};
  EXPECT_THAT(variable_registry.GetPositionCommands(),
              Pointwise(DoubleNear(0.01), expected_position_values));
  // Now we expect the commanded feedforward velocity values.
  expected_feedforward_velocity_values = {-8, -16, -32};
  EXPECT_THAT(
      variable_registry.GetFeedforwardVelocityCommands(),
      Pointwise(DoubleNear(0.01), expected_feedforward_velocity_values));
}

TEST_F(RobotArmBusComponentTestFixture,
       WriteProvidesScaledJointPositionCommandsWithIntegralType) {
  LoopbackFakeVariableRegistry<int32_t, 3> variable_registry =
      LoopbackFakeVariableRegistry<int32_t, 3>(
          /*joint_states=*/{0, 0, 0},
          /*position_commands=*/{0, 0, 0},
          /*ff_velocity_commands=*/{0, 0, 0});

  fieldbus::DeviceInitContext init_context(*interface_registry_,
                                           variable_registry);
  ASSERT_OK_AND_ASSIGN(
      std::unique_ptr<RobotArmBusComponent> robot_arm_bus_component,
      RobotArmBusComponent::Create(
          init_context, ParseTextProtoOrDie(kTypical3DofDeviceConfig)));

  ASSERT_OK_AND_ASSIGN(
      auto position_command,
      interface_registry_
          ->GetMutableInterfaceHandle<intrinsic_fbs::JointPositionCommand>(
              "foo_joint_position_command"));

  position_command->mutable_position()->Mutate(0, 2);
  position_command->mutable_position()->Mutate(1, 4);
  position_command->mutable_position()->Mutate(2, 8);
  position_command->mutable_velocity_feedforward()->Mutate(0, 16);
  position_command->mutable_velocity_feedforward()->Mutate(1, 32);
  position_command->mutable_velocity_feedforward()->Mutate(2, 64);

  // Not simulating ICON tick.
  EXPECT_THAT(
      robot_arm_bus_component->CyclicWrite(RequestType::kNormalOperation),
      RealtimeIsOkAndHolds(RequestStatus::kDone));

  // Expect the sensed position {0, 0, 0} if the position command update
  // timestamp is not updated.
  std::array<int32_t, 3> expected_position_values = {0, 0, 0};
  EXPECT_THAT(variable_registry.GetPositionCommands(),
              Pointwise(Eq(), expected_position_values));
  // Expect the velocity feedforward {0, 0, 0} if the position command update
  // timestamp is not updated.
  std::array<int32_t, 3> expected_feedforward_velocity_values = {0, 0, 0};
  EXPECT_THAT(variable_registry.GetFeedforwardVelocityCommands(),
              Pointwise(Eq(), expected_feedforward_velocity_values));

  // Try again with updated timestamp.
  {
    // Simulates an ICON tick.
    intrinsic::icon::Cycle::SetCurrentCycle(42);
    icon_state_->mutate_current_cycle(42);
    icon_state_.UpdatedAt(intrinsic::Clock::now());
    position_command.UpdatedAt(intrinsic::Clock::Now());
  }
  EXPECT_THAT(
      robot_arm_bus_component->CyclicWrite(RequestType::kNormalOperation),
      RealtimeIsOkAndHolds(RequestStatus::kDone));

  // Now we expect the commanded position values.
  expected_position_values = {1, 2, 4};
  EXPECT_THAT(variable_registry.GetPositionCommands(),
              Pointwise(Eq(), expected_position_values));
  // Now we expect the commanded feedforward velocity values.
  expected_feedforward_velocity_values = {8, 16, 32};
  EXPECT_THAT(variable_registry.GetFeedforwardVelocityCommands(),
              Pointwise(Eq(), expected_feedforward_velocity_values));
}

}  // namespace
}  // namespace intrinsic::robot_arm
