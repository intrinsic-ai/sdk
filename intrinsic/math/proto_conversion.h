// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_MATH_PROTO_CONVERSION_H_
#define INTRINSIC_MATH_PROTO_CONVERSION_H_

#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/math/proto/accel.pb.h"
#include "intrinsic/math/proto/affine.pb.h"
#include "intrinsic/math/proto/matrix.pb.h"
#include "intrinsic/math/proto/point.pb.h"
#include "intrinsic/math/proto/pose.pb.h"
#include "intrinsic/math/proto/quaternion.pb.h"
#include "intrinsic/math/proto/twist.pb.h"
#include "intrinsic/math/proto/vector3.pb.h"
#include "intrinsic/math/twist.h"

// All conversions from INTRINSIC protos to their respective C++ types should be
// declared in namespace intrinsic_proto. This makes it possible to make
// unqualified calls to FromProto() throughout our code base.
namespace intrinsic_proto {

constexpr int kMaxMatrixProtoDimension = 2048;
absl::StatusOr<intrinsic::eigenmath::MatrixXd> FromProto(
    const Matrixd& proto_matrix);

absl::StatusOr<intrinsic::eigenmath::AffineTransform3d> FromProto(
    const Affine3d& affine3d);
// The is the maximum allowable dimension of a matrix stored in a proto.
intrinsic::Twist FromProto(const intrinsic_proto::Twist& twist);
intrinsic::Acceleration FromProto(const intrinsic_proto::Accel& accel);

}  // namespace intrinsic_proto

// All conversions from Intrinsic protos to their respective C++ types should be
// declared in namespace intrinsic_proto. This makes it possible to make
// unqualified calls to FromProto() throughout our code base.
namespace intrinsic_proto {

intrinsic::eigenmath::Vector3d FromProto(const Point& point);
intrinsic::eigenmath::Quaterniond FromProto(const Quaternion& quaternion);
absl::StatusOr<intrinsic::Pose> FromProto(const Pose& pose);
// Only checks if the quaternion of the pose is roughly normalized, in which
// case it normalizes the input quaternion before generating the pose. If it is
// as normalized as expected in `FromProto`, then no normalization is performed.
absl::StatusOr<intrinsic::Pose> FromProtoNormalized(const Pose& pose);

intrinsic::eigenmath::Vector3d FromProto(const Vector3& vec3);

}  // namespace intrinsic_proto

// To enable unqualified calls to ToProto() throughout our code base, we declare
// functions which convert C++ types to protos in namespace intrinsic.
namespace intrinsic {

intrinsic_proto::Pose ToProto(const Pose& pose);
intrinsic_proto::Point ToProto(const eigenmath::Vector3d& point);
intrinsic_proto::Quaternion ToProto(const eigenmath::Quaterniond& quaternion);

// Instantiate specific common matrix sizes to remove ambiguity between matrix
// and vector.
intrinsic_proto::Affine3d ToProto(
    const eigenmath::AffineTransform3d& affine_transform);
absl::StatusOr<intrinsic_proto::Matrixd> ToProto(
    const eigenmath::MatrixXd& matrix);
intrinsic_proto::Matrixd ToProto(const eigenmath::Matrix3d& matrix);
intrinsic_proto::Matrixd ToProto(const eigenmath::Matrix4d& matrix);
intrinsic_proto::Matrixd ToProto(const eigenmath::Matrix4dAligned& matrix);
intrinsic_proto::Matrixd ToProto(const eigenmath::Matrix6d& matrix);

intrinsic_proto::Twist ToProto(const Twist& twist);
intrinsic_proto::Accel ToProto(const intrinsic::Acceleration& acceleration);

}  // namespace intrinsic

#endif  // INTRINSIC_MATH_PROTO_CONVERSION_H_
