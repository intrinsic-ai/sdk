// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_CONTROL_ALGORITHMS_JOINT_POSITION_PID_VELOCITY_CONTROLLER_H_
#define INTRINSIC_ICON_CONTROL_ALGORITHMS_JOINT_POSITION_PID_VELOCITY_CONTROLLER_H_

#include <memory>
#include <optional>

#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/proto/joint_position_pid_velocity_controller_config.pb.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/math/signals/butter_filter2.h"

namespace intrinsic {
namespace icon {

// PI controller to convert cyclic position setpoints to velocity setpoints.
class JointPositionPIDVelocityController {
 public:
  struct Params {
    eigenmath::VectorNd k_p;
    eigenmath::VectorNd k_i;
    eigenmath::VectorNd k_d;
    eigenmath::VectorNd k_ff;
    eigenmath::VectorNd max_integral_control;
    eigenmath::VectorNd max_velocity_command;
    double cycle_time_sec;
    std::optional<double> position_filter_cuttoff_frequency_hz = std::nullopt;
    std::optional<double> velocity_filter_cuttoff_frequency_hz = std::nullopt;
  };
  struct Filters {
    std::unique_ptr<ButterFilter2<eigenmath::VectorNd>>
        butterworth_position_filter;
    std::unique_ptr<ButterFilter2<eigenmath::VectorNd>>
        butterworth_velocity_filter;
    bool position_butterworth_initialized = false;
    bool velocity_butterworth_initialized = false;
  };

  struct State {
    explicit State(int num_joints) {
      integral_control = eigenmath::VectorNd::Zero(num_joints);
      filtered_position = eigenmath::VectorNd::Zero(num_joints);
      previous_velocity_command = eigenmath::VectorNd::Zero(num_joints);
    }
    eigenmath::VectorNd filtered_position;
    eigenmath::VectorNd integral_control;
    eigenmath::VectorNd previous_velocity_command;
  };

  static absl::StatusOr<std::unique_ptr<JointPositionPIDVelocityController>>
  Create(
      intrinsic_proto::icon::JointPositionPidVelocityControllerConfig config);

  // Returns the velocity setpoint when it was calculated successfully and
  // velocity_setpoint_ was updated. Returns error otherwise.
  RealtimeStatusOr<eigenmath::VectorNd> CalculateSetpoints(
      const eigenmath::VectorNd& position_desired,
      const eigenmath::VectorNd& velocity_feedforward,
      const eigenmath::VectorNd& position_state,
      const eigenmath::VectorNd& velocity_state);

  // Resets any internal state (integrator, setpoint, etc.).
  void Reset();

 protected:
  explicit JointPositionPIDVelocityController(
      Params params,
      std::unique_ptr<ButterFilter2<eigenmath::VectorNd>>
          butterworth_position_filter,
      std::unique_ptr<ButterFilter2<eigenmath::VectorNd>>
          butterworth_velocity_filter);
  Params params_;
  State state_;
  Filters filters_;
};
}  // namespace icon
}  // namespace intrinsic
#endif  // INTRINSIC_ICON_CONTROL_ALGORITHMS_JOINT_POSITION_PID_VELOCITY_CONTROLLER_H_
