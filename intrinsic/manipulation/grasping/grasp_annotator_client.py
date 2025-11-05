# Copyright 2023 Intrinsic Innovation LLC

"""Defines the GraspAnnotatorClient class."""

import threading

from absl import logging
from google.rpc import code_pb2
import grpc
from intrinsic.geometry.proto import triangle_mesh_pb2
from intrinsic.manipulation.grasping import grasp_annotations_pb2
from intrinsic.manipulation.grasping import grasp_annotator_pb2
from intrinsic.manipulation.service.grasp_annotator_service.v1 import grasp_annotator_service_pb2
from intrinsic.manipulation.service.grasp_annotator_service.v1 import grasp_annotator_service_pb2_grpc
from intrinsic.world.python import object_world_resources

DEFAULT_GRASP_ANNOTATOR_SERVICE_ADDRESS = (
    "istio-ingressgateway.app-ingress.svc.cluster.local:80"
)
DEFAULT_GRASP_ANNOTATOR_SERVICE_INSTANCE_NAME = "grasp_annotator_service"


class GraspAnnotatorClient:
  """Helper class for calling the rpcs in the GraspAnnotatorService."""

  def __init__(
      self,
      stub: grasp_annotator_service_pb2_grpc.GraspAnnotatorStub,
      instance_name: str = DEFAULT_GRASP_ANNOTATOR_SERVICE_INSTANCE_NAME,
  ):
    """Constructor.

    Args:
      stub: The GraspannotatorService stub.
      instance_name: The service instance name of the grasp annotator service.
        This is the name defined in `intrinsic_resource_instance`.
    """
    self._stub: grasp_annotator_service_pb2_grpc.GraspAnnotatorStub = stub
    self._connection_params = {
        "metadata": [(
            "x-resource-instance-name",
            instance_name,
        )]
    }

  @classmethod
  def connect(
      cls,
      address: str = DEFAULT_GRASP_ANNOTATOR_SERVICE_ADDRESS,
      instance_name: str = DEFAULT_GRASP_ANNOTATOR_SERVICE_INSTANCE_NAME,
  ) -> tuple[grpc.Channel, "GraspAnnotatorClient"]:
    """Connects to the grasp annotator service.

    Args:
      address: The address of the grasp annotator service.
      instance_name: The service instance name of the grasp annotator service.
        This is the name defined in `intrinsic_resource_instance`.

    Returns:
      gRpc channel, grasp annotator client
    """
    logging.info("Connecting to grasp_annotator_service at %s", address)
    channel = grpc.insecure_channel(address)
    return channel, GraspAnnotatorClient(
        stub=grasp_annotator_service_pb2_grpc.GraspAnnotatorStub(channel),
        instance_name=instance_name,
    )

  def annotate_grasps(
      self,
      triangle_mesh: triangle_mesh_pb2.TriangleMesh,
      gripper_specs: grasp_annotator_pb2.ParameterizedGripperSpecs,
      num_samples: int,
      annotation_metrics_weights: (
          grasp_annotator_pb2.MetricWeights | None
      ) = None,
      max_num_annotations: int | None = None,
      constraint: grasp_annotator_pb2.GraspAnnotationConstraint | None = None,
  ) -> grasp_annotations_pb2.GraspAnnotations:
    """Annotates grasps.

    Args:
      triangle_mesh: The mesh to annotate on. See
        `manipulation_utils.aggregate_object_meshes` or
        `manipulation_utils.aggregate_scene_object_meshes` for converting an
        object reference or SceneObject to a triangle mesh.
      gripper_specs: The gripper specifications.
      num_samples: The number of samples to query on the mesh.
      annotation_metrics_weights: The metrics weights to score annotation with.
      max_num_annotations: The maximum number of annotations to return.
      constraint: Constraints to filter grasp poses.

    Returns:
      The annotated grasps as a `GraspAnnotations` proto.
    """
    request = grasp_annotator_service_pb2.GraspAnnotatorRequest(
        mesh_data=grasp_annotator_pb2.MeshData(triangle_mesh=triangle_mesh),
        gripper_specs=gripper_specs,
        num_samples=num_samples,
    )
    if annotation_metrics_weights:
      request.annotation_metrics_weights.CopyFrom(annotation_metrics_weights)
    if max_num_annotations is not None:
      request.max_num_annotations = max_num_annotations
    if constraint:
      request.constraint.CopyFrom(constraint)
    response: grasp_annotator_service_pb2.GraspAnnotatorResponse = (
        self._stub.Annotate(request, **self._connection_params)
    )
    return response.annotations

  def generate_grasps(
      self,
      triangle_mesh: triangle_mesh_pb2.TriangleMesh,
      gripper_specs: grasp_annotator_pb2.ParameterizedGripperSpecs,
      num_samples: int,
  ) -> grasp_annotations_pb2.GraspAnnotations:
    """Generates raw grasp annotations.

    Args:
      triangle_mesh: The mesh to generate grasps on. See
        `manipulation_utils.aggregate_object_meshes` or
        `manipulation_utils.aggregate_scene_object_meshes` for converting an
        object reference or SceneObject to a triangle mesh.
      gripper_specs: The gripper specifications.
      num_samples: The number of samples to query on the mesh.

    Returns:
      The unprocessed grasp annotations as a `GraspAnnotations` proto.
    """
    request = grasp_annotator_service_pb2.GraspAnnotatorGenerateRequest(
        mesh_data=grasp_annotator_pb2.MeshData(triangle_mesh=triangle_mesh),
        gripper_specs=gripper_specs,
        num_samples=num_samples,
    )
    response: grasp_annotator_service_pb2.GraspAnnotatorGenerateResponse = (
        self._stub.Generate(request, **self._connection_params)
    )
    return response.unprocessed_annotations

  def filter_grasps(
      self,
      unfiltered_annotations: grasp_annotations_pb2.GraspAnnotations,
      triangle_mesh: triangle_mesh_pb2.TriangleMesh,
      gripper_specs: grasp_annotator_pb2.ParameterizedGripperSpecs,
      constraint: grasp_annotator_pb2.GraspAnnotationConstraint | None = None,
  ) -> grasp_annotations_pb2.GraspAnnotations:
    """Filters grasp annotations.

    Args:
      unfiltered_annotations: The unfiltered grasp annotations.
      triangle_mesh: The mesh to filter grasps on. See
        `manipulation_utils.aggregate_object_meshes` or
        `manipulation_utils.aggregate_scene_object_meshes` for converting an
        object reference or SceneObject to a triangle mesh.
      gripper_specs: The gripper specifications.
      constraint: Constraints to filter grasp poses.

    Returns:
      The filtered grasp annotations as a `GraspAnnotations` proto.
    """
    request = grasp_annotator_service_pb2.GraspAnnotatorFilterRequest(
        mesh_data=grasp_annotator_pb2.MeshData(triangle_mesh=triangle_mesh),
        gripper_specs=gripper_specs,
        unfiltered_annotations=unfiltered_annotations,
    )
    if constraint:
      request.constraint.CopyFrom(constraint)
    response: grasp_annotator_service_pb2.GraspAnnotatorFilterResponse = (
        self._stub.Filter(request, **self._connection_params)
    )
    return response.filtered_annotations

  def score_grasps(
      self,
      unscored_annotations: grasp_annotations_pb2.GraspAnnotations,
      triangle_mesh: triangle_mesh_pb2.TriangleMesh,
      gripper_specs: grasp_annotator_pb2.ParameterizedGripperSpecs,
      annotation_metrics_weights: (
          grasp_annotator_pb2.MetricWeights | None
      ) = None,
      max_num_annotations: int | None = None,
  ) -> grasp_annotations_pb2.GraspAnnotations:
    """Scores and sorts grasp annotations.

    Args:
      unscored_annotations: The unscored grasp annotations.
      triangle_mesh: The mesh to score grasps on. See
        `manipulation_utils.aggregate_object_meshes` or
        `manipulation_utils.aggregate_scene_object_meshes` for converting an
        object reference or SceneObject to a triangle mesh.
      gripper_specs: The gripper specifications.
      annotation_metrics_weights: The metrics weights to score annotation with.
      max_num_annotations: The maximum number of annotations to return.

    Returns:
      The scored and sorted grasp annotations as a `GraspAnnotations` proto.
    """
    request = grasp_annotator_service_pb2.GraspAnnotatorScoreRequest(
        mesh_data=grasp_annotator_pb2.MeshData(triangle_mesh=triangle_mesh),
        gripper_specs=gripper_specs,
        unscored_annotations=unscored_annotations,
    )
    if annotation_metrics_weights:
      request.annotation_metrics_weights.CopyFrom(annotation_metrics_weights)
    if max_num_annotations is not None:
      request.max_num_annotations = max_num_annotations
    response: grasp_annotator_service_pb2.GraspAnnotatorScoreResponse = (
        self._stub.Score(request, **self._connection_params)
    )
    return response.scored_annotations

  def visualize_grasps_async(
      self,
      annotations: grasp_annotations_pb2.GraspAnnotations,
      object_instance: object_world_resources.WorldObject,
      lifetime_sec: int,
      constraint: grasp_annotator_pb2.GraspAnnotationConstraint | None = None,
  ) -> None:
    """Starts visualizing grasp annotations in a background thread.

    This is a fire-and-forget call.

    Args:
      annotations: The annotations to visualize.
      object_instance: The object instance to visualize annotations on top of.
      lifetime_sec: The lifetime of the visualization in seconds.
      constraint: Constraints to filter grasp poses. If set, the constraints
        will be visualized as well.
    """
    logging.info("Starting async visualization request...")

    frame_id = object_instance.get_tf_name_for_root_entity()

    grasp_annotator_visualize_request = grasp_annotator_service_pb2.GraspAnnotatorVisualizeRequest(
        annotations=annotations,
        frame_id=frame_id,
        constraint=constraint,
        visualization_options=grasp_annotator_service_pb2.VisualizationOptions(
            lifetime_sec=lifetime_sec
        ),
    )

    def _consume_stream_in_background():
      """Consumes the stream to completion, logging any errors."""
      try:
        for response in self._stub.Visualize(
            grasp_annotator_visualize_request,
            **self._connection_params,
        ):
          if response.status.code != code_pb2.OK:
            logging.error(
                "Async visualization failed: %s", response.status.message
            )
            break
      except grpc.RpcError as e:
        logging.exception("gRPC error during async visualization: %s", e)

    thread = threading.Thread(
        target=_consume_stream_in_background,
        daemon=True,
    )
    thread.start()
    logging.info("Async visualization task started.")

  def visualize_grasps_blocking(
      self,
      annotations: grasp_annotations_pb2.GraspAnnotations,
      object_instance: object_world_resources.WorldObject,
      lifetime_sec: int,
      constraint: grasp_annotator_pb2.GraspAnnotationConstraint | None = None,
  ) -> None:
    """Visualizes grasp annotations and waits for completion.

    Args:
      annotations: The annotations to visualize.
      object_instance: The object instance to visualize annotations on top of.
      lifetime_sec: The lifetime of the visualization in seconds.
      constraint: Constraints to filter grasp poses. If set, the constraints
        will be visualized as well.

    Raises:
      RuntimeError: If visualization fails or there is a gRPC communication
        error.
    """
    logging.info("Sending blocking visualization request...")
    frame_id = object_instance.get_tf_name_for_root_entity()

    grasp_annotator_visualize_request = grasp_annotator_service_pb2.GraspAnnotatorVisualizeRequest(
        annotations=annotations,
        frame_id=frame_id,
        constraint=constraint,
        visualization_options=grasp_annotator_service_pb2.VisualizationOptions(
            lifetime_sec=lifetime_sec
        ),
    )
    try:
      final_status = None
      for response in self._stub.Visualize(
          grasp_annotator_visualize_request,
          **self._connection_params,
      ):
        final_status = response.status
        if final_status.code != code_pb2.OK:
          break

      logging.info("Visualization stream finished.")

      if not final_status:
        raise RuntimeError("No response received from visualization stream.")

      if final_status.code != code_pb2.OK:
        raise RuntimeError(f"Visualization failed: {final_status.message}")

    except grpc.RpcError as e:
      raise RuntimeError(f"gRPC error during visualization: {e}") from e
