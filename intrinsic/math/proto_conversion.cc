// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/math/proto_conversion.h"

#include <cmath>

#include "absl/algorithm/container.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_format.h"
#include "absl/strings/substitute.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/math/almost_equals.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/math/proto/affine.pb.h"
#include "intrinsic/math/proto/twist.pb.h"
#include "intrinsic/math/proto/vector3.pb.h"
#include "intrinsic/math/twist.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic_proto {

absl::StatusOr<intrinsic::eigenmath::MatrixXd> FromProto(
    const Matrixd& proto_matrix) {
  if (proto_matrix.rows() > intrinsic_proto::kMaxMatrixProtoDimension ||
      proto_matrix.rows() < 1) {
    return absl::InvalidArgumentError(absl::Substitute(
        "Invalid number of rows in matrix proto: $0. Must be "
        "in the range [1, $1]",
        proto_matrix.rows(), intrinsic_proto::kMaxMatrixProtoDimension));
  }

  if (proto_matrix.cols() > intrinsic_proto::kMaxMatrixProtoDimension ||
      proto_matrix.cols() < 1) {
    return absl::InvalidArgumentError(absl::Substitute(
        "Invalid number of columns in matrix proto: $0. Must "
        "be in the range [1, $1]",
        proto_matrix.cols(), intrinsic_proto::kMaxMatrixProtoDimension));
  }

  intrinsic::eigenmath::MatrixXd eigen_matrix(proto_matrix.rows(),
                                              proto_matrix.cols());
  if (proto_matrix.values().size() != eigen_matrix.size()) {
    return absl::InvalidArgumentError(
        absl::Substitute("The number of elements in the matrix doesn't match "
                         "the size (cols x rows) definition: $0 vs $1",
                         proto_matrix.values().size(), eigen_matrix.size()));
  }
  absl::c_copy(proto_matrix.values(), eigen_matrix.reshaped().begin());
  return eigen_matrix;
}

absl::StatusOr<intrinsic::eigenmath::AffineTransform3d> FromProto(
    const Affine3d& affine3d) {
  intrinsic::eigenmath::AffineTransform3d eigen_affine;
  INTR_ASSIGN_OR_RETURN(eigen_affine.linear(), FromProto(affine3d.linear()));
  eigen_affine.translation() = FromProto(affine3d.translation());
  return eigen_affine;
}

intrinsic::Twist FromProto(const intrinsic_proto::Twist& twist) {
  return {twist.linear().x(),  twist.linear().y(),  twist.linear().z(),
          twist.angular().x(), twist.angular().y(), twist.angular().z()};
}

intrinsic::Acceleration FromProto(const intrinsic_proto::Accel& accel) {
  return {accel.linear().x(),  accel.linear().y(),  accel.linear().z(),
          accel.angular().x(), accel.angular().y(), accel.angular().z()};
}

}  // namespace intrinsic_proto

namespace intrinsic_proto {

intrinsic::eigenmath::Vector3d FromProto(const Point& point) {
  return {point.x(), point.y(), point.z()};
}

intrinsic::eigenmath::Quaterniond FromProto(const Quaternion& quaternion) {
  // Eigen's Quaternion ctor takes parameters in the order {w, x, y, z}.
  return {quaternion.w(), quaternion.x(), quaternion.y(), quaternion.z()};
}

absl::StatusOr<intrinsic::Pose> FromProto(const Pose& pose) {
  intrinsic::eigenmath::Quaterniond quaternion = FromProto(pose.orientation());
  // We need to perform a soft-check in here, since otherwise, we might raise an
  // error status due to numeric errors introduced by the squared norm
  // computation.
  // Using exact float comparison, it is not necessarily true that
  //   Quaterniond::UnitRandom().norm() == 1.0
  if (const double squared_norm = quaternion.squaredNorm();
      !intrinsic::AlmostEquals(squared_norm, 1.0)) {
    const intrinsic::eigenmath::Quaterniond normalized_quat =
        quaternion.normalized();
    return absl::InvalidArgumentError(absl::StrFormat(
        "Failed to create Pose from proto which contains a "
        "non-unit quaternion with norm(quat) == %.17f . The normalized "
        "quaternion would be %.17f, %.17f, %.17f, %.17f",
        std::sqrt(squared_norm), normalized_quat.x(), normalized_quat.y(),
        normalized_quat.z(), normalized_quat.w()));
  }

  intrinsic::eigenmath::Vector3d position = FromProto(pose.position());
  if (position.hasNaN()) {
    return absl::InvalidArgumentError(
        absl::StrFormat("Failed to create Pose from proto which contains a "
                        "nan position values: {%.17f, %.17f, %.17f}",
                        position.x(), position.y(), position.z()));
  }

  return intrinsic::Pose(quaternion, position,
                         intrinsic::eigenmath::kDoNotNormalize);
}

absl::StatusOr<intrinsic::Pose> FromProtoNormalized(const Pose& pose) {
  intrinsic::eigenmath::Quaterniond quaternion = FromProto(pose.orientation());
  constexpr double kNormalizationError = 1e-3;
  const double squared_norm = quaternion.squaredNorm();
  // If we're already normalized to a reasonable degree, then simply don't
  // renormalize at all. This preserves the property that if we call
  // FromProtoNormalized(ToProto(pose)) that we get the same result back.
  if (intrinsic::AlmostEquals(squared_norm, 1.0)) {
    return intrinsic::Pose(quaternion, FromProto(pose.position()),
                           intrinsic::eigenmath::kDoNotNormalize);
  }

  const intrinsic::eigenmath::Quaterniond normalized_quat =
      quaternion.normalized();
  // We need to perform a soft-check in here, since otherwise, we might raise an
  // error status due to numeric errors introduced by the squared norm
  // computation or due to rounding. To enforce higher precision checks and not
  // allow normalization of the quaternion, please use FromProto directly.
  if (!intrinsic::AlmostEquals(squared_norm, 1.0, kNormalizationError)) {
    return absl::InvalidArgumentError(absl::StrFormat(
        "Failed to create Pose from proto which contains a "
        "non-unit quaternion with norm(quat) == %.6f . The normalized "
        "quaternion would be %.4f, %.4f, %.4f, %.4f. Provided quaternion %v",
        std::sqrt(squared_norm), normalized_quat.x(), normalized_quat.y(),
        normalized_quat.z(), normalized_quat.w(), pose));
  }
  return intrinsic::Pose(normalized_quat, FromProto(pose.position()),
                         intrinsic::eigenmath::kDoNotNormalize);
}

intrinsic::eigenmath::Vector3d FromProto(const Vector3& vec3) {
  return {vec3.x(), vec3.y(), vec3.z()};
}

}  // namespace intrinsic_proto

namespace intrinsic {

intrinsic_proto::Pose ToProto(const Pose& pose) {
  intrinsic_proto::Pose proto_pose;
  *proto_pose.mutable_position() = ToProto(pose.translation());
  // Compensate for numerical inaccuracies which may accumulate from many
  // composed rotations to maintain serialization idempotency.
  eigenmath::Quaterniond quat = pose.quaternion();
  if (!intrinsic::AlmostEquals(quat.squaredNorm(), 1)) {
    quat.normalize();
  }
  *proto_pose.mutable_orientation() = ToProto(quat);
  return proto_pose;
}

intrinsic_proto::Point ToProto(const eigenmath::Vector3d& point) {
  intrinsic_proto::Point proto_point;
  proto_point.set_x(point.x());
  proto_point.set_y(point.y());
  proto_point.set_z(point.z());
  return proto_point;
}

intrinsic_proto::Quaternion ToProto(const eigenmath::Quaterniond& quaternion) {
  intrinsic_proto::Quaternion proto_quaternion;
  proto_quaternion.set_x(quaternion.x());
  proto_quaternion.set_y(quaternion.y());
  proto_quaternion.set_z(quaternion.z());
  proto_quaternion.set_w(quaternion.w());
  return proto_quaternion;
}

intrinsic_proto::Affine3d ToProto(
    const eigenmath::AffineTransform3d& affine_transform) {
  intrinsic_proto::Affine3d proto_affine3d;
  const eigenmath::Matrix3d& eigen_linear = affine_transform.linear();
  const eigenmath::Vector3d& eigen_position = affine_transform.translation();
  *proto_affine3d.mutable_linear() = ToProto(eigen_linear);
  *proto_affine3d.mutable_translation() = ToProto(eigen_position);
  return proto_affine3d;
}
namespace {
template <typename MatrixType>
absl::StatusOr<intrinsic_proto::Matrixd> ToProtoImpl(
    const MatrixType& eigen_matrix) {
  if (eigen_matrix.rows() < 1 ||
      eigen_matrix.rows() > intrinsic_proto::kMaxMatrixProtoDimension) {
    return absl::InvalidArgumentError(absl::Substitute(
        "Invalid number of rows in matrix for serialization: $0. Must be "
        "in the range [1, $1]",
        eigen_matrix.rows(), intrinsic_proto::kMaxMatrixProtoDimension));
  }

  if (eigen_matrix.cols() < 1 ||
      eigen_matrix.cols() > intrinsic_proto::kMaxMatrixProtoDimension) {
    return absl::InvalidArgumentError(absl::Substitute(
        "Invalid number of columns in matrix for serialization: $0. Must be "
        "in the range [1, $1]",
        eigen_matrix.cols(), intrinsic_proto::kMaxMatrixProtoDimension));
  }

  const auto reshaped_matrix = eigen_matrix.reshaped();  // NOLINT
  intrinsic_proto::Matrixd proto_matrix;
  proto_matrix.set_rows(eigen_matrix.rows());
  proto_matrix.set_cols(eigen_matrix.cols());
  proto_matrix.mutable_values()->Add(reshaped_matrix.begin(),
                                     reshaped_matrix.end());
  return proto_matrix;
}
}  // namespace

intrinsic_proto::Matrixd ToProto(const eigenmath::Matrix3d& matrix) {
  return ToProtoImpl(matrix).value();
}
absl::StatusOr<intrinsic_proto::Matrixd> ToProto(
    const eigenmath::MatrixXd& matrix) {
  return ToProtoImpl(matrix);
}

intrinsic_proto::Matrixd ToProto(const eigenmath::Matrix4d& matrix) {
  return ToProtoImpl(matrix).value();
}
intrinsic_proto::Matrixd ToProto(const eigenmath::Matrix4dAligned& matrix) {
  return ToProtoImpl(matrix).value();
}
intrinsic_proto::Matrixd ToProto(const eigenmath::Matrix6d& matrix) {
  return ToProtoImpl(matrix).value();
}

intrinsic_proto::Twist ToProto(const Twist& twist) {
  intrinsic_proto::Twist proto_twist;
  proto_twist.mutable_linear()->set_x(twist[0]);
  proto_twist.mutable_linear()->set_y(twist[1]);
  proto_twist.mutable_linear()->set_z(twist[2]);
  proto_twist.mutable_angular()->set_x(twist[3]);
  proto_twist.mutable_angular()->set_y(twist[4]);
  proto_twist.mutable_angular()->set_z(twist[5]);
  return proto_twist;
}
intrinsic_proto::Accel ToProto(const intrinsic::Acceleration& accel) {
  intrinsic_proto::Accel proto_accel;
  proto_accel.mutable_linear()->set_x(accel[0]);
  proto_accel.mutable_linear()->set_y(accel[1]);
  proto_accel.mutable_linear()->set_z(accel[2]);
  proto_accel.mutable_angular()->set_x(accel[3]);
  proto_accel.mutable_angular()->set_y(accel[4]);
  proto_accel.mutable_angular()->set_z(accel[5]);
  return proto_accel;
}

}  // namespace intrinsic
