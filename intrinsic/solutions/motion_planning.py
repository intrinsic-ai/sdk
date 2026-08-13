# Copyright 2026 Intrinsic Innovation LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""A wrapper of MotionPlannerClientExternal that is more convenient to use."""

from intrinsic.motion_planning import motion_planner_client
from intrinsic.motion_planning.proto.v1 import motion_planner_service_pb2_grpc
from intrinsic.solutions import deployments


class MotionPlannerClient(motion_planner_client.MotionPlannerClientBase):
  """A wrapper of MotionPlannerClientBase that is more convenient to use."""

  @classmethod
  def for_solution(
      cls, solution: deployments.Solution
  ) -> "MotionPlannerClient":
    """Connects to the motion planner gRPC service for a given solution."""
    return cls(
        world_id=solution.world.world_id,
        stub=motion_planner_service_pb2_grpc.MotionPlannerServiceStub(
            solution.grpc_channel
        ),
    )
