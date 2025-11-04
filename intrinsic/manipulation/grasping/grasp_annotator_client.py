# Copyright 2023 Intrinsic Innovation LLC

"""Defines the GraspAnnotatorClient class."""

import threading

from absl import logging
from google.rpc import code_pb2
import grpc
from intrinsic.manipulation.grasping import grasp_annotations_pb2
from intrinsic.manipulation.service.grasp_annotator_service.v1 import grasp_annotator_service_pb2
from intrinsic.manipulation.service.grasp_annotator_service.v1 import grasp_annotator_service_pb2_grpc

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
      grasp_annotator_request: grasp_annotator_service_pb2.GraspAnnotatorRequest,
  ) -> grasp_annotations_pb2.GraspAnnotations:
    """Annotates grasps.

    Args:
      grasp_annotator_request: The parameters used to annotate grasps.

    Returns:
      The annotated grasps.
    """
    response: grasp_annotator_service_pb2.GraspAnnotatorResponse = (
        self._stub.Annotate(
            grasp_annotator_request,
            **self._connection_params,
        )
    )
    return response.annotations

  def generate_grasps(
      self,
      generate_request: grasp_annotator_service_pb2.GraspAnnotatorGenerateRequest,
  ) -> grasp_annotations_pb2.GraspAnnotations:
    """Generates raw grasp annotations."""
    response: grasp_annotator_service_pb2.GraspAnnotatorGenerateResponse = (
        self._stub.Generate(generate_request, **self._connection_params)
    )
    return response.unprocessed_annotations

  def filter_grasps(
      self,
      filter_request: grasp_annotator_service_pb2.GraspAnnotatorFilterRequest,
  ) -> grasp_annotations_pb2.GraspAnnotations:
    """Filters grasp annotations."""
    response: grasp_annotator_service_pb2.GraspAnnotatorFilterResponse = (
        self._stub.Filter(filter_request, **self._connection_params)
    )
    return response.filtered_annotations

  def score_grasps(
      self,
      score_request: grasp_annotator_service_pb2.GraspAnnotatorScoreRequest,
  ) -> grasp_annotations_pb2.GraspAnnotations:
    """Scores and sorts grasp annotations."""
    response: grasp_annotator_service_pb2.GraspAnnotatorScoreResponse = (
        self._stub.Score(score_request, **self._connection_params)
    )
    return response.scored_annotations

  def visualize_grasps_async(
      self,
      grasp_annotator_visualize_request: grasp_annotator_service_pb2.GraspAnnotatorVisualizeRequest,
  ) -> None:
    """Starts visualizing grasp annotations in a background thread.

    This is a fire-and-forget call.

    Args:
      grasp_annotator_visualize_request: The annotations to visualize.
    """
    logging.info("Starting async visualization request...")

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
      grasp_annotator_visualize_request: grasp_annotator_service_pb2.GraspAnnotatorVisualizeRequest,
  ) -> None:
    """Visualizes grasp annotations and waits for completion.

    Args:
      grasp_annotator_visualize_request: The annotations to visualize.

    Raises:
      RuntimeError: If visualization fails or there is a gRPC communication
        error.
    """
    logging.info("Sending blocking visualization request...")
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
