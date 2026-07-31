// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_INTERNAL_UTIL_EIGEN_SERIALIZATION_H_
#define INTRINSIC_GEOMETRY_INTERNAL_UTIL_EIGEN_SERIALIZATION_H_

#include <cstddef>

#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "intrinsic/eigenmath/pose2.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/marshal/riegeli_coder.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/util/status/status_macros.h"
#include "riegeli/records/record_reader.h"
#include "riegeli/records/record_writer.h"

namespace intrinsic {

template <typename _Scalar, int _Rows, int _Cols, int _Options, int _MaxRows,
          int _MaxCols>
struct RiegeliCoder<
    Eigen::Matrix<_Scalar, _Rows, _Cols, _Options, _MaxRows, _MaxCols>> {
  static absl::Status Encode(
      const Eigen::Matrix<_Scalar, _Rows, _Cols, _Options, _MaxRows, _MaxCols>&
          value,
      riegeli::RecordWriterBase& writer) {
    INTR_RETURN_IF_ERROR(RiegeliCoder<size_t>::Encode(value.rows(), writer))
        << "Failed to write matrix rows";
    INTR_RETURN_IF_ERROR(RiegeliCoder<size_t>::Encode(value.cols(), writer))
        << "Failed to write matrix cols";
    for (int i = 0; i < value.rows(); ++i) {
      for (int j = 0; j < value.cols(); ++j) {
        INTR_RETURN_IF_ERROR(RiegeliCoder<_Scalar>::Encode(value(i, j), writer))
            << absl::StrCat("Failed to write matrix element [", i, ",", j, "]");
      }
    }
    return absl::OkStatus();
  }
  static absl::StatusOr<
      Eigen::Matrix<_Scalar, _Rows, _Cols, _Options, _MaxRows, _MaxCols>>
  Decode(riegeli::RecordReaderBase& reader) {
    auto output =
        Eigen::Matrix<_Scalar, _Rows, _Cols, _Options, _MaxRows, _MaxCols>();
    INTR_ASSIGN_OR_RETURN(size_t rows, RiegeliCoder<size_t>::Decode(reader));
    INTR_ASSIGN_OR_RETURN(size_t cols, RiegeliCoder<size_t>::Decode(reader));
    output.resize(rows, cols);
    for (int i = 0; i < rows; ++i) {
      for (int j = 0; j < cols; ++j) {
        INTR_ASSIGN_OR_RETURN(output(i, j),
                              RiegeliCoder<_Scalar>::Decode(reader));
      }
    }
    return std::move(output);
  }
};

template <>
absl::Status RiegeliCoder<::intrinsic::eigenmath::Pose2d>::Encode(
    const ::intrinsic::eigenmath::Pose2d& value,
    riegeli::RecordWriterBase& writer);

template <>
absl::StatusOr<::intrinsic::eigenmath::Pose2d>
RiegeliCoder<::intrinsic::eigenmath::Pose2d>::Decode(
    riegeli::RecordReaderBase& reader);

template <>
absl::Status RiegeliCoder<::intrinsic::Pose3d>::Encode(
    const ::intrinsic::Pose3d& value, riegeli::RecordWriterBase& writer);

template <>
absl::StatusOr<::intrinsic::Pose3d> RiegeliCoder<::intrinsic::Pose3d>::Decode(
    riegeli::RecordReaderBase& reader);

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_UTIL_EIGEN_SERIALIZATION_H_
