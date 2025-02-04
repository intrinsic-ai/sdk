# Copyright 2023 Intrinsic Innovation LLC

"""Defines the GraspPlannerClient class."""

from typing import Sequence, Tuple
from absl import logging
from google.protobuf import empty_pb2
import grpc
from intrinsic.manipulation.grasping import grasp_execution_planner_params_pb2
from intrinsic.manipulation.grasping import grasp_pb2
from intrinsic.manipulation.grasping import grasp_planner_params_pb2
from intrinsic.manipulation.grasping import grasp_planning_interfaces
from intrinsic.manipulation.grasping import grasp_ranker_params_pb2
from intrinsic.manipulation.service import grasp_planner_service_pb2
from intrinsic.manipulation.service import grasp_planner_service_pb2_grpc
from intrinsic.world.proto import object_world_refs_pb2
import numpy as np

DEFAULT_GRASP_PLANNER_SERVICE_ADDRESS = (
    "istio-ingressgateway.app-ingress.svc.cluster.local:80"
)
DEFAULT_GRASP_PLANNER_SERVICE_INSTANCE_NAME = "grasp_planner_service"


class GraspPlannerClient:
  """Helper class for calling the rpcs in the GraspPlannerService."""

  def __init__(
      self,
      stub: grasp_planner_service_pb2_grpc.GraspPlannerServiceStub,
      instance_name: str = DEFAULT_GRASP_PLANNER_SERVICE_INSTANCE_NAME,
  ):
    """Constructor.

    Args:
      stub: The GraspPlannerService stub.
      instance_name: The service instance name of the grasp planner service.
        This is the name defined in `intrinsic_resource_instance`.
    """
    self._stub: grasp_planner_service_pb2_grpc.GraspPlannerServiceStub = stub
    self._connection_params = {
        "metadata": [(
            "x-resource-instance-name",
            instance_name,
        )]
    }

  @classmethod
  def connect(
      cls,
      address: str = DEFAULT_GRASP_PLANNER_SERVICE_ADDRESS,
      instance_name: str = DEFAULT_GRASP_PLANNER_SERVICE_INSTANCE_NAME,
  ) -> Tuple[grpc.Channel, "GraspPlannerClient"]:
    """Connects to the grasp planner service.

    Args:
      address: The address of the grasp planner service.
      instance_name: The service instance name of the grasp planner service.
        This is the name defined in `intrinsic_resource_instance`.

    Returns:
      gRpc channel, grasp planner client
    """
    # Create connection to the service
    logging.info("Connecting to grasp_planner_service at %s", address)

    channel = grpc.insecure_channel(address)
    return channel, GraspPlannerClient(
        stub=grasp_planner_service_pb2_grpc.GraspPlannerServiceStub(channel),
        instance_name=instance_name,
    )

  def register_grasp_planner(
      self,
      planner_id: str,
      grasp_planner_params: grasp_planner_params_pb2.GraspPlannerParams,
      grasp_execution_planner_params: grasp_execution_planner_params_pb2.GraspExecutionPlannerParams,
  ) -> None:
    """Registers a GraspPlanner.

    Args:
      planner_id: The id of grasp planner to register.
      grasp_planner_params: The parameters used to construct a grasp planner.
      grasp_execution_planner_params: The parameters used to construct a grasp
        execution planner.

    Raises:
      RuntimeError: If registration fails.
    """
    response = self._stub.RegisterGraspPlanner(
        grasp_planner_service_pb2.RegisterGraspPlannerRequest(
            planner_id=planner_id,
            grasp_planner_params=grasp_planner_params,
            grasp_execution_planner_params=grasp_execution_planner_params,
        ),
        **self._connection_params,
    )
    if response.success:
      logging.info("Successfully registered a grasp planner %s.", planner_id)
    else:
      raise RuntimeError(
          f"Failed to register a grasp planner {planner_id}:"
          f" {response.debug_message}.",
      )

  def remove_all_grasp_planners(self) -> None:
    """Removes all registered grasp planners."""
    self._stub.RemoveAllGraspPlanners(
        empty_pb2.Empty(), **self._connection_params
    )

  def plan_grasps(
      self,
      planner_id: str,
      plan_grasps_params: grasp_planner_service_pb2.GPSPlanGraspsParams,
  ) -> grasp_pb2.GraspPlan:
    """Plan grasps.

    Args:
      planner_id: The id of grasp planner to plan grasps.
      plan_grasps_params: The parameters used to plan grasps.

    Returns:
      The planned grasps.
    """
    logging.info("Calling PlanGrasps for planner %s.", planner_id)
    plan = self._stub.PlanGrasps(
        grasp_planner_service_pb2.PlanGraspsRequest(
            planner_id=planner_id,
            plan_grasps_params=plan_grasps_params,
        ),
        **self._connection_params,
    )
    logging.info("Planning successful: %r.", bool(plan.grasps))
    return plan

  def notify_grasp_results(
      self,
      planner_id: str,
      executed_grasps: Sequence[grasp_pb2.AttemptedGrasp],
  ) -> None:
    """Feeds execution results to the grasp planner manager.

    Args:
      planner_id: The id of grasp planner manager to notify the grasps.
      executed_grasps: The executed grasp results.
    """
    logging.info("Calling NotifyGraspResults for planner %s.", planner_id)
    self._stub.NotifyGraspResults(
        grasp_planner_service_pb2.NotifyGraspResultsRequest(
            planner_id=planner_id,
            executed_grasps=executed_grasps,
        ),
        **self._connection_params,
    )
    return

  def plan_grasp_execution(
      self,
      planner_id: str,
      plan_grasp_execution_params: grasp_planner_service_pb2.GPSPlanGraspExecutionParams,
  ) -> grasp_planning_interfaces.GraspExecutionPlanningResult:
    """Plan a grasp execution.

    Args:
      planner_id: The id of grasp planner to plan the grasp execution.
      plan_grasp_execution_params: The parameters used to plan the grasp
        execution.

    Returns:
      A grasp execution planning result.
    """
    logging.info(
        "Calling PlanGraspExecution for grasp %s.",
        plan_grasp_execution_params.grasp.grasp_id,
    )
    planning_result_proto = self._stub.PlanGraspExecution(
        grasp_planner_service_pb2.PlanGraspExecutionRequest(
            planner_id=planner_id,
            plan_grasp_execution_params=plan_grasp_execution_params,
        ),
        **self._connection_params,
    )
    planning_result = (
        grasp_planning_interfaces.GraspExecutionPlanningResult.from_proto(
            planning_result_proto
        )
    )
    logging.info("Planning successful: %r.", bool(planning_result.success))
    return planning_result

  def rank_grasps(
      self,
      grasps: Sequence[grasp_pb2.Grasp],
      robot_name: str,
      tool_frame: object_world_refs_pb2.TransformNodeReference,
      rankers_params: grasp_ranker_params_pb2.GraspRankersParams,
      world_id: str,
      score_threshold: float,
      recent_grasps: Sequence[grasp_pb2.AttemptedGrasp] | None = None,
      max_num_grasps: int | None = None,
      grasp_ranking_soft_timeout_in_sec: float = np.inf,
  ) -> grasp_pb2.GraspPlan:
    """Rank grasps.

    Args:
      grasps: The grasps to rank.
      robot_name: The name of the robot.
      tool_frame: The reference to the tool frame.
      rankers_params: The parameters used to construct the grasp rankers.
      world_id: The id of the world.
      score_threshold: The threshold to filter out low scored grasps.
      recent_grasps: The recently executed grasps. If None, use an empty list.
      max_num_grasps: The maximum number of grasps to return.
      grasp_ranking_soft_timeout_in_sec: The soft timeout to return early if at
        least one valid grasp is found.

    Returns:
      The ranked grasps.
    """
    logging.info("Ranking %d grasps.", len(grasps))

    return self._stub.RankGrasps(
        grasp_planner_service_pb2.RankGraspsRequest(
            grasps=grasps,
            robot_name=robot_name,
            tool_frame=tool_frame,
            rankers_params=rankers_params,
            world_id=world_id,
            score_threshold=score_threshold,
            max_num_grasps=max_num_grasps,
            grasp_ranking_soft_timeout_in_sec=grasp_ranking_soft_timeout_in_sec,
            recent_grasps=recent_grasps,
        ),
        **self._connection_params,
    )
