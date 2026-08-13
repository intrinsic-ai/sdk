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

#include "intrinsic/icon/hal/interfaces/imu_utils.h"

#include "flatbuffers/detached_buffer.h"
#include "flatbuffers/flatbuffer_builder.h"
#include "intrinsic/icon/flatbuffers/transform_types.fbs.h"
#include "intrinsic/icon/hal/interfaces/imu.fbs.h"

namespace intrinsic_fbs {

flatbuffers::DetachedBuffer BuildImuStatus() {
  flatbuffers::FlatBufferBuilder builder;
  builder.ForceDefaults(true);

  ImuStatusBuilder status_builder(builder);

  intrinsic_fbs::Point linear_acceleration;
  status_builder.add_linear_acceleration(&linear_acceleration);

  intrinsic_fbs::Point angular_velocity;
  status_builder.add_angular_velocity(&angular_velocity);

  intrinsic_fbs::Rotation orientation;
  orientation.mutate_qw(1.0);
  status_builder.add_orientation(&orientation);

  auto status = status_builder.Finish();
  builder.Finish(status);

  return builder.Release();
}

}  // namespace intrinsic_fbs
