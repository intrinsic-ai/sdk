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

#ifndef INTRINSIC_ICON_PROTO_CONCATENATE_TRAJECTORY_PROTOS_H_
#define INTRINSIC_ICON_PROTO_CONCATENATE_TRAJECTORY_PROTOS_H_

#include <vector>

#include "absl/status/statusor.h"
#include "intrinsic/icon/proto/joint_space.pb.h"

namespace intrinsic {

// Splits a trajectory `proto` into a vector of sequential sub-trajectories with
// `max_subtrajectory_length`. Time stamps remain untouched in the different
// trajectory segments. Returns kFailedPrecondition in case of invalid
// `max_subtrajectory_length` or in case of an empty `proto`.
absl::StatusOr<std::vector<intrinsic_proto::icon::JointTrajectoryPVA>>
SplitTrajectoryProto(const intrinsic_proto::icon::JointTrajectoryPVA& proto,
                     int max_subtrajectory_length);

// Concatenates `trajectory_segments` to a single trajectory in a first-in-first
// out fashion. It is assumed that time stamps throughout the
// `trajectory_segments` are monotonically increasing, and that the first time
// stamp of a segment is greater than the last time stamp of the preceding
// segment. Returns kFailedPrecondition if `trajectories` is empty.
absl::StatusOr<intrinsic_proto::icon::JointTrajectoryPVA>
ConcatenateTrajectoryProtos(
    const std::vector<intrinsic_proto::icon::JointTrajectoryPVA>&
        trajectory_segments);

}  // namespace intrinsic

#endif  // INTRINSIC_ICON_PROTO_CONCATENATE_TRAJECTORY_PROTOS_H_
