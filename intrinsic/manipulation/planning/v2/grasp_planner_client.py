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

"""Defines the v2 GraspPlannerClient class."""

from absl import logging
import grpc

from intrinsic.assets import interface_utils
from intrinsic.assets.dependencies import utils
from intrinsic.assets.proto.v1 import resolved_dependency_pb2
from intrinsic.manipulation.service.grasp_planner_service.v2 import grasp_planner_service_pb2
from intrinsic.manipulation.service.grasp_planner_service.v2 import grasp_planner_service_pb2_grpc

DEFAULT_GRASP_PLANNER_SERVICE_ADDRESS: str = (
    "istio-ingressgateway.app-ingress.svc.cluster.local:80"
)
DEFAULT_GRASP_PLANNER_SERVICE_INSTANCE_NAME: str = "grasp_planner_service"


def grasp_planner_service_interface_uri() -> str:
  """Returns the interface URI of the GraspPlannerService v2 dependency."""
  return (
      f"{interface_utils.GRPC_URI_PREFIX}"
      f"{grasp_planner_service_pb2.DESCRIPTOR.services_by_name['GraspPlannerService'].full_name}"
  )


class GraspPlannerClient:
  """Helper class for calling the RPCs in GraspPlannerService v2."""

  def __init__(
      self,
      stub: grasp_planner_service_pb2_grpc.GraspPlannerServiceStub,
      instance_name: str | None = DEFAULT_GRASP_PLANNER_SERVICE_INSTANCE_NAME,
  ):
    """Constructor.

    Args:
      stub: The GraspPlannerService stub.
      instance_name: The service instance name of the grasp planner service.
        This is the name defined in `intrinsic_resource_instance`. Defaults to
        `DEFAULT_GRASP_PLANNER_SERVICE_INSTANCE_NAME` if omitted. If set to
        None or empty string, connection params will not be set in RPC calls.
    """
    self._stub: grasp_planner_service_pb2_grpc.GraspPlannerServiceStub = stub
    self._connection_params = {}
    if instance_name:
      self._connection_params["metadata"] = [(
          "x-resource-instance-name",
          instance_name,
      )]

  @classmethod
  def connect(
      cls,
      address: str = DEFAULT_GRASP_PLANNER_SERVICE_ADDRESS,
      instance_name: str = DEFAULT_GRASP_PLANNER_SERVICE_INSTANCE_NAME,
  ) -> tuple[grpc.Channel, "GraspPlannerClient"]:
    """Connects to the grasp planner service.

    Args:
      address: The address of the grasp planner service.
      instance_name: The service instance name of the grasp planner service.
        This is the name defined in `intrinsic_resource_instance`.

    Returns:
      gRPC channel, grasp planner client
    """
    logging.info("Connecting to grasp_planner_service v2 at %s", address)
    channel = grpc.insecure_channel(address)
    return channel, GraspPlannerClient(
        stub=grasp_planner_service_pb2_grpc.GraspPlannerServiceStub(channel),
        instance_name=instance_name,
    )

  @classmethod
  def connect_dependency(
      cls,
      dep: resolved_dependency_pb2.ResolvedDependency,
  ) -> tuple[grpc.Channel, "GraspPlannerClient"]:
    """Connects to the grasp planner service via a ResolvedDependency.

    Args:
      dep: The resolved dependency for the grasp planner service.

    Returns:
      gRPC channel, grasp planner client
    """
    channel = utils.connect(dep, grasp_planner_service_interface_uri())
    return channel, GraspPlannerClient(
        stub=grasp_planner_service_pb2_grpc.GraspPlannerServiceStub(channel),
        instance_name=None,
    )

  def plan_grasps(
      self,
      request: grasp_planner_service_pb2.PlanGraspsRequest,
  ) -> grasp_planner_service_pb2.PlanGraspsResponse:
    """Plans grasps using the GraspPlannerService v2.

    Args:
      request: The PlanGraspsRequest specifying targets, grasp annotations,
        robot name, tool frame, and ranking configuration.

    Returns:
      PlanGraspsResponse containing the planned grasps and/or failure diagnostics.
    """
    return self._stub.PlanGrasps(request, **self._connection_params)
