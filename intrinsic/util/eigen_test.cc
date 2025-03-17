// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/eigen.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <array>
#include <cstddef>
#include <limits>
#include <string>
#include <vector>

#include "google/protobuf/repeated_field.h"
#include "intrinsic/eigenmath/rotation_utils.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

namespace intrinsic {

TEST(UtilEigen, VectorXdUtils) {
  eigenmath::VectorXd value(6);
  for (size_t i = 0; i < value.size(); ++i) {
    value[i] = i + 1;
  }

  google::protobuf::RepeatedField<double> rpt_field;
  VectorXdToRepeatedDouble(value, &rpt_field);

  eigenmath::VectorXd decoded_value = RepeatedDoubleToVectorXd(rpt_field);

  EXPECT_EQ(value, decoded_value);

  std::vector<double> vector_value = VectorXdToVector(value);
  EXPECT_EQ(vector_value.size(), value.size());
  EXPECT_EQ(value, VectorToVectorXd(vector_value));

  std::array<double, 6> array_value = VectorXdToArray<6>(value);
  EXPECT_EQ(array_value.size(), value.size());
  EXPECT_EQ(value, ArrayToVectorXd(array_value));
}

TEST(UtilEigen, MatrixXdUtils) {
  std::vector<eigenmath::VectorXd> vector_of_vectorxd = {
      (eigenmath::VectorXd(4) << 1.2, 4.5, 7.8, 10.1).finished(),
      (eigenmath::VectorXd(4) << 2.3, 5.6, 8.9, 11.2).finished(),
      (eigenmath::VectorXd(4) << 3.4, 6.7, 9.0, 12.3).finished()};
  eigenmath::MatrixXd matrixxd = VectorOfVectorXdToMatrixXd(vector_of_vectorxd);
  eigenmath::MatrixXd expected_matrixxd(4, 3);
  expected_matrixxd << 1.2, 2.3, 3.4, 4.5, 5.6, 6.7, 7.8, 8.9, 9.0, 10.1, 11.2,
      12.3;
  EXPECT_EQ(matrixxd, expected_matrixxd);
}

TEST(UtilEigen, VectorUtils) {
  std::vector<double> value(6);
  for (size_t i = 0; i < value.size(); ++i) {
    value[i] = i + 1;
  }

  google::protobuf::RepeatedField<double> rpt_field;
  VectorDoubleToRepeatedDouble(value, &rpt_field);

  std::vector<double> decoded_value = RepeatedDoubleToVectorDouble(rpt_field);

  EXPECT_EQ(value, decoded_value);
}

}  // namespace intrinsic
