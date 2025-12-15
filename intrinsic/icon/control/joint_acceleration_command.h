// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_CONTROL_JOINT_ACCELERATION_COMMAND_H_
#define INTRINSIC_ICON_CONTROL_JOINT_ACCELERATION_COMMAND_H_

#include <cstddef>
#include <optional>

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/utils/realtime_status_or.h"

namespace intrinsic::icon {

// Represents a set of command parameters for joint acceleration control.
class JointAccelerationCommand {
 public:
  // Default constructor to make this play nice with containers and StatusOr.
  JointAccelerationCommand() : acceleration_(0) {}

  explicit JointAccelerationCommand(const eigenmath::VectorNd& acceleration)
      : acceleration_(acceleration) {}

  // Builds a JointAccelerationCommand object.
  //
  // Returns InvalidArgument if any of the two vectors' sizes don't match up.
  static RealtimeStatusOr<JointAccelerationCommand> Create(
      const eigenmath::VectorNd& acceleration,
      const std::optional<eigenmath::VectorNd>& torque = std::nullopt);

  const eigenmath::VectorNd& acceleration() const;

  const std::optional<eigenmath::VectorNd>& torque() const;

  size_t Size() const;

 private:
  JointAccelerationCommand(const eigenmath::VectorNd& acceleration,
                           const std::optional<eigenmath::VectorNd>& torque);

  eigenmath::VectorNd acceleration_;
  std::optional<eigenmath::VectorNd> torque_ = std::nullopt;
};

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_CONTROL_JOINT_ACCELERATION_COMMAND_H_
