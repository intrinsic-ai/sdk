// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_CONTROL_C_API_C_TYPES_H_
#define INTRINSIC_ICON_CONTROL_C_API_C_TYPES_H_

#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

constexpr size_t kIntrinsicIconMaxNumberOfJoints = 25;

// C API type to convey joint position commands with optional feedforwards.
// The arrays are fixed-size so they can live on the stack, but only the given
// number of elements is valid to read.
//
// It is the responsibility of the creator to ensure that the arrays (if
// present) all have `size` *valid* elements.
struct IntrinsicIconJointPositionCommand {
  size_t size;
  // Array of position setpoints.
  double position_setpoints[kIntrinsicIconMaxNumberOfJoints];
  // Array of velocity feedforwards. Set velocity_feedforwards_size to zero if
  // there are no velocity feedforwards.
  double velocity_feedforwards[kIntrinsicIconMaxNumberOfJoints];
  bool has_velocity_feedforwards;
  // Array of acceleration feedforwards. Set acceleration_feedforwards_size to
  // zero if there are no acceleration feedforwards.
  double acceleration_feedforwards[kIntrinsicIconMaxNumberOfJoints];
  bool has_acceleration_feedforwards;
};

// C API type to convey joint limits.
// The arrays are fixed-size so they can live on the stack, but only the given
// number of elements is valid to read.
//
// It is the responsibility of the creator to ensure that the arrays (if
// present) all have `size` *valid* elements.
struct IntrinsicIconJointLimits {
  size_t size;
  // Array of minimum position values.
  double min_position[kIntrinsicIconMaxNumberOfJoints];
  // Array of maximum position values.
  double max_position[kIntrinsicIconMaxNumberOfJoints];
  // Array of maximum velocity values.
  double max_velocity[kIntrinsicIconMaxNumberOfJoints];
  // Array of maximum acceleration values.
  double max_acceleration[kIntrinsicIconMaxNumberOfJoints];
  // Array of maximum jerk values.
  double max_jerk[kIntrinsicIconMaxNumberOfJoints];
  // Array of maximum torque values.
  double max_torque[kIntrinsicIconMaxNumberOfJoints];
};

// C API type to convey positional joint state.
struct IntrinsicIconJointStateP {
  size_t size;
  double positions[kIntrinsicIconMaxNumberOfJoints];
};

// C API type to convey velocity joint state.
struct IntrinsicIconJointStateV {
  size_t size;
  double velocities[kIntrinsicIconMaxNumberOfJoints];
};

// C API type to convey acceleration joint state.
struct IntrinsicIconJointStateA {
  size_t size;
  double accelerations[kIntrinsicIconMaxNumberOfJoints];
};

struct IntrinsicIconQuaternion {
  double w;
  double x;
  double y;
  double z;
};

struct IntrinsicIconPoint {
  double x;
  double y;
  double z;
};

struct IntrinsicIconPose3d {
  IntrinsicIconQuaternion rotation;
  IntrinsicIconPoint translation;
};

struct IntrinsicIconWrench {
  double x;
  double y;
  double z;
  double rx;
  double ry;
  double rz;
};

// Note that this is *not* zero-terminated, so dereferencing `data + size` is an
// error.
struct IntrinsicIconString {
  char* data;
  size_t size;
};

typedef void (*IntrinsicIconStringDestroy)(IntrinsicIconString* str);

// Same as above, not zero-terminated. In addition, this is a pure view whose
// storage is not valid outside of its original scope (i.e. functions that take
// IntrinsicIconStringView must not hold references after they finish).
struct IntrinsicIconStringView {
  const char* data;
  const size_t size;
};

// C API type for 6 x N matrices of double values,
// with N <= kIntrinsicIconMaxNumberOfJoints.
struct IntrinsicIconMatrix6Nd {
  // The number of columns in the matrix. Must not be greater than
  // kIntrinsicIconMaxNumberOfJoints.
  size_t num_cols;
  // Matrix values, in column-major order (to match Eigen's default order).
  // Indices >= 6 * num_cols are invalid and contain unspecified data!
  double data[6 * kIntrinsicIconMaxNumberOfJoints];
};

struct IntrinsicIconSignalValue {
  bool current_value;
  bool previous_value;
};

#ifdef __cplusplus
}
#endif

#endif  // INTRINSIC_ICON_CONTROL_C_API_C_TYPES_H_
