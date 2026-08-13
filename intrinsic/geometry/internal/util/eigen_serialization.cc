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
