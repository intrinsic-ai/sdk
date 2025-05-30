// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_EIGENMATH_SO3_H_
#define INTRINSIC_EIGENMATH_SO3_H_

#include <cmath>
#include <cstdlib>
#include <iomanip>
#include <ios>
#include <ostream>
#include <sstream>
#include <string>
#include <type_traits>

#include "Eigen/Core"
#include "absl/log/check.h"
#include "absl/status/status.h"
#include "intrinsic/eigenmath/rotation_utils.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/testing/realtime_annotations.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/util/status/status_builder.h"

namespace intrinsic {
namespace eigenmath {

// A representation of 3D rotations using unit quaternions
template <typename Scalar, int Options = kDefaultOptions>
class SO3 {
 public:
  template <int OtherOptions>
  using Quaternion = Quaternion<Scalar, OtherOptions>;

  // Initializes to the identity rotation
  EIGEN_DEVICE_FUNC SO3()
      : quaternion_(Scalar(1), Scalar(0), Scalar(0), Scalar(0)) {}

  // Initializes with a rotation matrix
  //
  // If do_orthogonalize is set to true (default), the given matrix will be
  // orthogonalized through a singular value decomposition.
  explicit SO3(const Matrix3<Scalar>& matrix, bool do_orthogonalize = true) {
    if (do_orthogonalize) {
      quaternion_ = OrthogonalizeRotationMatrix(matrix);
    } else {
      quaternion_ = matrix;
    }
  }

  // Creates a SO3 from `matrix`. This method is real-time safe.
  // If do_orthogonalize is true (default), the matrix will be orthogonalized
  // through a singular value decomposition.
  // An error is returned if either
  // - do_orthogonalize is true, but the derived rotation after
  //   orthogonalization is not valid, or
  // - do_orthogonalize is false and the given matrix is not orthogonal.
  template <int OtherOptions = kDefaultOptions>
  static icon::RealtimeStatusOr<SO3> FromMatrixRealtimeSafe(
      const Matrix3<Scalar>& matrix,
      bool do_orthogonalize = true) INTRINSIC_CHECK_REALTIME_SAFE {
    Quaternion<OtherOptions> quaternion;
    if (do_orthogonalize) {
      quaternion = OrthogonalizeRotationMatrix(matrix);
    } else {
      quaternion = matrix;
    }
    if (!IsNormalizedQuaternion(quaternion)) {
      return icon::InvalidArgumentError(
          "Cannot create rotation from given matrix.");
    }
    return SO3(quaternion, kUnsafeCtor);
  }

  // Creates a SO3 from `matrix`. This method is not real-time safe but provides
  // a detailed error message in case of failure.
  // If do_orthogonalize is true (default), the matrix will be orthogonalized
  // through a singular value decomposition.
  // An error is returned if either
  // - do_orthogonalize is true, but the derived rotation after
  //   orthogonalization is not valid, or
  // - do_orthogonalize is false and the given matrix is not orthogonal.
  template <int OtherOptions = kDefaultOptions>
  static absl::StatusOr<SO3> FromMatrix(const Matrix3<Scalar>& matrix,
                                        bool do_orthogonalize = true) {
    Quaternion<OtherOptions> quaternion;
    if (do_orthogonalize) {
      quaternion = OrthogonalizeRotationMatrix(matrix);
    } else {
      quaternion = matrix;
    }
    if (!IsNormalizedQuaternion(quaternion)) {
      const Eigen::IOFormat inline_format(
          Eigen::StreamPrecision, Eigen::DontAlignCols,
          /*_coeffSeparator=*/", ", /*_rowSeparator=*/", ",
          /*_rowPrefix=*/"[", /*_rowSuffix=*/"]",
          /*_matPrefix=*/"[", /*_matSuffix=*/"].");
      return intrinsic::InvalidArgumentErrorBuilder()
             << "Cannot create rotation from matrix: "
             << matrix.format(inline_format);
    }
    return SO3(quaternion, kUnsafeCtor);
  }

  // Initializes with RPY angles
  //
  // RPY angles are usually really useful as a user representation of
  // orientation.
  SO3(Scalar roll, Scalar pitch, Scalar yaw) {
    RotationFromRPY(roll, pitch, yaw, &quaternion_);
  }

  // Initializes using quaternion
  //
  // Either do_normalize must be true, or quaternion must be normalized.
  //
  // By default the quaternion is normalized.
  template <int OtherOptions = kDefaultOptions>
  explicit SO3(const Quaternion<OtherOptions>& quaternion,
               bool do_normalize = true)
      : quaternion_(quaternion) {
    if (do_normalize) {
      quaternion_.normalize();
    }
    const bool is_normalized = IsNormalized();
    if (!is_normalized) {
      // Allocating here is OK, as the program will be terminated in any case.
      CHECK(is_normalized) << ExplainUnNormalizedQuaternion(quaternion_);
    }
  }

  // Creates a SO3 from `quaternion`. This method is real-time safe.
  // If do_normalize is true (default), the quaternion will be normalized.
  // An error is returned if either
  // - do_normalize is true, but the quaternion cannot be normalized, or
  // - do_normalize is false and the quaternion is not normalized.
  template <int OtherOptions = kDefaultOptions>
  static icon::RealtimeStatusOr<SO3> FromQuaternionRealtimeSafe(
      const Quaternion<OtherOptions>& quaternion,
      bool do_normalize = true) INTRINSIC_CHECK_REALTIME_SAFE {
    Quaternion<OtherOptions> quaternion_normalized = quaternion;
    if (do_normalize) {
      quaternion_normalized.normalize();
    }
    bool is_normalized = IsNormalizedQuaternion(quaternion_normalized);
    if (!is_normalized) {
      return icon::InvalidArgumentError(
          "Cannot create rotation from quaternion. " +
          ExplainUnNormalizedQuaternion(quaternion));
    }
    return SO3(quaternion_normalized, kUnsafeCtor);
  }

  // Creates a SO3 from `quaternion`. This method is not real-time safe but
  // provides a detailed error message in case of failure.
  // If do_normalize is true (default), the quaternion will be normalized.
  // An error is returned if either
  // - do_normalize is true, but the quaternion cannot be normalized, or
  // - do_normalize is false and the quaternion is not normalized.
  template <int OtherOptions = kDefaultOptions>
  static absl::StatusOr<SO3> FromQuaternion(
      const Quaternion<OtherOptions>& quaternion, bool do_normalize = true) {
    Quaternion<OtherOptions> quaternion_normalized = quaternion;
    if (do_normalize) {
      quaternion_normalized.normalize();
    }
    bool is_normalized = IsNormalizedQuaternion(quaternion_normalized);
    if (!is_normalized) {
      return absl::InvalidArgumentError(
          "Cannot create rotation from quaternion. " +
          ExplainUnNormalizedQuaternion(quaternion));
    }
    return SO3(quaternion_normalized, kUnsafeCtor);
  }

  // Conversion operator for other SO3 types with different Eigen::Options.
  template <int OtherOptions>
  EIGEN_DEVICE_FUNC SO3(const SO3<Scalar, OtherOptions>& other)  // NOLINT
      : quaternion_(other.quaternion()) {}

  // Assignment operator for other SO3 types with different Eigen::Options.
  template <int OtherOptions>
  EIGEN_DEVICE_FUNC SO3& operator=(const SO3<Scalar, OtherOptions>& other) {
    quaternion_ = other.quaternion();
    return *this;
  }

  // Assignment operator for other SO3 types with different Eigen::Options.
  template <int OtherOptions>
  SO3& operator=(const Quaternion<OtherOptions>& quaternion) {
    quaternion_ = quaternion;
    quaternion_.normalize();
    const bool is_normalized = IsNormalized();
    if (!is_normalized) {
      // Allocating here is OK, as the program will be terminated in any case.
      CHECK(is_normalized) << ExplainUnNormalizedQuaternion(quaternion_);
    }
    return *this;
  }

  // The quaternion.
  EIGEN_DEVICE_FUNC const Quaternion<Options>& quaternion() const {
    return quaternion_;
  }

  // Writeable quaternion.
  EIGEN_DEVICE_FUNC Quaternion<Options>& quaternion() { return quaternion_; }

  // Computes and returns the magnitude of the rotation in radians.
  EIGEN_DEVICE_FUNC Scalar norm() const {
    return Scalar(2) * std::atan2(quaternion_.vec().norm(),
                                  Eigen::numext::abs(quaternion_.w()));
  }

  // Returns the corresponding 3D rotation matrix.
  EIGEN_DEVICE_FUNC Matrix3<Scalar> matrix() const {
    return quaternion_.toRotationMatrix();
  }

  // The inverse rotation.
  EIGEN_DEVICE_FUNC SO3<Scalar> inverse() const {
    return SO3<Scalar>(quaternion_.conjugate(), kUnsafeCtor);
  }

  // Read-only pointer to underlying data.
  EIGEN_DEVICE_FUNC const Scalar* data() const {
    return quaternion_.coeffs().data();
  }

  // Cast SO3 instance to other scalar type.
  template <typename OtherScalar>
  EIGEN_DEVICE_FUNC SO3<OtherScalar> cast() const {
    if constexpr (std::is_same<OtherScalar, Scalar>::value) {
      return *this;
    } else {
      // force normalize call inside constructor
      // (crucial when going from lower precision to higher precision)
      return SO3<OtherScalar>(quaternion_.template cast<OtherScalar>());
    }
  }

  // Checks if identical to another pose under a given tolerance.
  template <int OtherOptions>
  EIGEN_DEVICE_FUNC bool isApprox(const SO3<Scalar, OtherOptions>& other,
                                  Scalar tolerance) const {
    // Let θ be the angle of rotation required to get from this orientation to
    // "other". Below quantity is equal to (1−cosθ)/2, and gives a rough
    // estimate of the distance between "this" and "other". In particular, it
    // gives 0 whenever the quaternions represent the same orientation, and it
    // gives 1 whenever the two orientations are 180∘ apart.
    return 1.0 - std::pow(quaternion_.normalized().dot(
                              other.quaternion().normalized()),
                          2) <
           tolerance;
  }

  // Checks if identical to another pose under default tolerance.
  template <int OtherOptions>
  EIGEN_DEVICE_FUNC bool isApprox(
      const SO3<Scalar, OtherOptions>& other) const {
    return isApprox(other, Eigen::NumTraits<Scalar>::dummy_precision());
  }

  // Multiplies rotation in-place.
  template <typename OtherScalar, int OtherOptions>
  EIGEN_DEVICE_FUNC SO3& operator*=(const SO3<OtherScalar, OtherOptions>& rhs) {
    quaternion_ *= rhs.quaternion();
    // Assuming the rotation's length is 1+epsilon due to floating-point
    // rounding errors, the code below reduces error to order epsilon^3 / 32.
    const Scalar nsq = quaternion_.squaredNorm();
    if (nsq != Scalar(1)) {
      quaternion_.coeffs() *= (Scalar(3) + nsq) / (Scalar(1) + Scalar(3) * nsq);
    }
    return *this;
  }

  // Composes two rotations.
  template <typename OtherScalar, int OtherOptions>
  EIGEN_DEVICE_FUNC auto operator*(
      const SO3<OtherScalar, OtherOptions>& rhs) const {
    using ResultScalar = std::common_type_t<Scalar, OtherScalar>;
    SO3<ResultScalar> result(*this);
    return result *= rhs;
  }

  // Ensures that the dot product of this quaternion and the provided
  // `reference` quaternion is positive, flipping the sign of this quaternion
  // if needed.
  //
  // The sign of each element (w, x, y, z) is flipped individually, thus
  // the resulting rotation is equivalent.
  template <typename OtherScalar, int OtherOptions>
  EIGEN_DEVICE_FUNC void MakeDotProductPositive(
      const SO3<OtherScalar, OtherOptions>& reference) {
    bool do_flip = quaternion().dot(reference.quaternion()) < 0.0;
    if (do_flip) {
      quaternion_ = Quaternion<Options>((-quaternion_.coeffs()).eval().data());
    }
  }

  template <int OtherOptions = kDefaultOptions>
  static bool IsNormalizedQuaternion(
      const Quaternion<OtherOptions>& quaternion) {
    using std::abs;  // for ADL
    return abs(quaternion.squaredNorm() - Scalar(1)) <
           Eigen::NumTraits<Scalar>::dummy_precision();
  }

  // Check whether the representation is normalized.
  bool IsNormalized() const { return IsNormalizedQuaternion(quaternion_); }

 private:
  Quaternion<Options> quaternion_;

  enum UnsafeCtorSignal { kUnsafeCtor };

  // Used to more efficiently construct an SO3 object when the invariants are
  // already guaranteed to hold by construction.
  template <int OtherOptions = kDefaultOptions>
  EIGEN_DEVICE_FUNC SO3(const Quaternion<OtherOptions>& quaternion,
                        UnsafeCtorSignal /*unused*/)
      : quaternion_(quaternion) {}

  // Generate a string explaining that the quaternion is not normalized. Used
  // to create detailed error messages in constructor/ factory functions above.
  template <int OtherOptions = kDefaultOptions>
  static std::string ExplainUnNormalizedQuaternion(
      const Quaternion<OtherOptions>& quaternion) {
    std::stringstream ss;
    ss << std::scientific << std::setprecision(18)
       << "Quaternion must be normalized (quaternion= " << quaternion
       << ", quaternion.squaredNorm()= " << quaternion.squaredNorm() << ")";
    return ss.str();
  }
};

// Outputs a SO3 to an ostream.
template <typename T, int Options>
std::ostream& operator<<(std::ostream& os,
                         const SO3<T, Options>& a_rotation_b) {
  os << "quaternion: " << a_rotation_b.quaternion().coeffs().transpose();
  return os;
}

// Rotates a 3D vector.
template <typename Scalar, int Options, typename Derived>
EIGEN_DEVICE_FUNC Vector<Scalar, 3> operator*(
    const SO3<Scalar, Options>& a_R_b,
    const Eigen::MatrixBase<Derived>& point_b) {
  static_assert(Derived::RowsAtCompileTime == 3,
                "Point must be three-dimensional");
  static_assert(Derived::ColsAtCompileTime == 1,
                "Point must be three-dimensional");
  return a_R_b.quaternion() * point_b;
}

// Overload for converting from SO3 to roll-pitch-yaw.
template <typename Scalar, int Options>
void SO3ToRPY(const SO3<Scalar, Options>& s, Scalar* roll, Scalar* pitch,
              Scalar* yaw) {
  RotationToRPY(s.quaternion(), roll, pitch, yaw);
}

using SO3dAligned = SO3<double, Eigen::AutoAlign>;
using SO3d = SO3<double, Eigen::DontAlign>;
using SO3fAligned = SO3<float, Eigen::AutoAlign>;
using SO3f = SO3<float, Eigen::DontAlign>;

}  // namespace eigenmath
}  // namespace intrinsic

#endif  // INTRINSIC_EIGENMATH_SO3_H_
