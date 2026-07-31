// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/geometry/internal/util/eigen_serialization.h"

#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/pose2.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/marshal/riegeli_coder.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/util/status/status_macros.h"
#include "riegeli/records/record_reader.h"
#include "riegeli/records/record_writer.h"

namespace intrinsic {

template <>
absl::Status RiegeliCoder<::intrinsic::eigenmath::Pose2d>::Encode(
    const ::intrinsic::eigenmath::Pose2d& value,
    riegeli::RecordWriterBase& writer) {
  INTR_RETURN_IF_ERROR(
      RiegeliCoder<double>::Encode(value.translation().x(), writer));
  INTR_RETURN_IF_ERROR(
      RiegeliCoder<double>::Encode(value.translation().y(), writer));
  INTR_RETURN_IF_ERROR(RiegeliCoder<double>::Encode(value.angle(), writer));
  return absl::OkStatus();
}

template <>
absl::StatusOr<::intrinsic::eigenmath::Pose2d>
RiegeliCoder<::intrinsic::eigenmath::Pose2d>::Decode(
    riegeli::RecordReaderBase& reader) {
  INTR_ASSIGN_OR_RETURN(double x, RiegeliCoder<double>::Decode(reader));
  INTR_ASSIGN_OR_RETURN(double y, RiegeliCoder<double>::Decode(reader));
  INTR_ASSIGN_OR_RETURN(double angle, RiegeliCoder<double>::Decode(reader));
  return ::intrinsic::eigenmath::Pose2d(::intrinsic::eigenmath::Vector2d(x, y),
                                        angle);
}

template <>
absl::Status RiegeliCoder<::intrinsic::Pose3d>::Encode(
    const ::intrinsic::Pose3d& value, riegeli::RecordWriterBase& writer) {
  const auto normalized_quaternion = value.quaternion().normalized();
  INTR_RETURN_IF_ERROR(
      RiegeliCoder<double>::Encode(value.translation().x(), writer));
  INTR_RETURN_IF_ERROR(
      RiegeliCoder<double>::Encode(value.translation().y(), writer));
  INTR_RETURN_IF_ERROR(
      RiegeliCoder<double>::Encode(value.translation().z(), writer));
  INTR_RETURN_IF_ERROR(
      RiegeliCoder<double>::Encode(normalized_quaternion.w(), writer));
  INTR_RETURN_IF_ERROR(
      RiegeliCoder<double>::Encode(normalized_quaternion.x(), writer));
  INTR_RETURN_IF_ERROR(
      RiegeliCoder<double>::Encode(normalized_quaternion.y(), writer));
  INTR_RETURN_IF_ERROR(
      RiegeliCoder<double>::Encode(normalized_quaternion.z(), writer));
  return absl::OkStatus();
}

template <>
absl::StatusOr<::intrinsic::Pose3d> RiegeliCoder<::intrinsic::Pose3d>::Decode(
    riegeli::RecordReaderBase& reader) {
  INTR_ASSIGN_OR_RETURN(double x, RiegeliCoder<double>::Decode(reader));
  INTR_ASSIGN_OR_RETURN(double y, RiegeliCoder<double>::Decode(reader));
  INTR_ASSIGN_OR_RETURN(double z, RiegeliCoder<double>::Decode(reader));
  INTR_ASSIGN_OR_RETURN(double q_w, RiegeliCoder<double>::Decode(reader));
  INTR_ASSIGN_OR_RETURN(double q_x, RiegeliCoder<double>::Decode(reader));
  INTR_ASSIGN_OR_RETURN(double q_y, RiegeliCoder<double>::Decode(reader));
  INTR_ASSIGN_OR_RETURN(double q_z, RiegeliCoder<double>::Decode(reader));
  return ::intrinsic::Pose3d(
      ::intrinsic::eigenmath::Quaterniond(q_w, q_x, q_y, q_z),
      ::intrinsic::eigenmath::Vector3d(x, y, z));
}

}  // namespace intrinsic
