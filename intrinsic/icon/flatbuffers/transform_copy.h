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


#ifndef INTRINSIC_ICON_FLATBUFFERS_TRANSFORM_COPY_H_
#define INTRINSIC_ICON_FLATBUFFERS_TRANSFORM_COPY_H_

#include "Eigen/Dense"
#include "flatbuffers/vector.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/flatbuffers/transform_types.fbs.h"
#include "intrinsic/icon/flatbuffers/transform_types.h"

namespace intrinsic_fbs {

// Copies the raw data memory from Eigen::Vector3d to memory in
// Point.
void CopyTo(const Eigen::Vector3d& vector, Point* point);

// Copies the raw data memory from Eigen::Quaterniond to memory in
// Rotation.
void CopyTo(const Eigen::Quaterniond& quaternion, Rotation* rotation);

// Copies the raw data memory from Eigen::VectorXd to memory in
// VectorNd.
void CopyTo(const Eigen::VectorXd& vector, VectorNd* vector_fbs);

// Copies the raw data memory from intrinsic::eigenmath::Vector6d to memory in
// Twist.
void CopyTo(const intrinsic::eigenmath::Vector6d& vector, Twist* twist);

// Copies the raw data memory from intrinsic::eigenmath::Vector6d to memory in
// Acceleration.
void CopyTo(const intrinsic::eigenmath::Vector6d& vector,
            Acceleration* acceleration);

// Copies the raw data memory from intrinsic::eigenmath::Vector6d to memory in
// Jerk.
void CopyTo(const intrinsic::eigenmath::Vector6d& vector, Jerk* jerk);

// Copies the raw data memory from intrinsic::eigenmath::Vector6d to memory in
// Wrench.
void CopyTo(const intrinsic::eigenmath::Vector6d& vector, Wrench* wrench);

// Copies the raw data memory from intrinsic::eigenmath::VectorNd to memory in
// VectorNd.
void CopyTo(const intrinsic::eigenmath::VectorNd& vector, VectorNd* vector_fbs);

// Assigns a value to a flatbuffer vector.
template <typename T, typename U>
void CopyTo(flatbuffers::Vector<T>* vector, const U& value) {
  for (int idx = 0; idx < vector->Length(); idx++) {
    vector->Mutate(idx, value);
  }
}

}  // namespace intrinsic_fbs

#endif  // INTRINSIC_ICON_FLATBUFFERS_TRANSFORM_COPY_H_
