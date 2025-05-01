# Copyright 2023 Intrinsic Innovation LLC

"""Defines the MultiAxisAutoAnnotatorClient class."""

from absl import logging
import grpc
from intrinsic.manipulation.grasping import schmalz_multi_axis_grasp_pb2
from intrinsic.manipulation.service import multi_axis_auto_annotator_service_pb2
from intrinsic.manipulation.service import multi_axis_auto_annotator_service_pb2_grpc


DEFAULT_SERVICE_ADDRESS = (
    "istio-ingressgateway.app-ingress.svc.cluster.local:80"
)
DEFAULT_SERVICE_INSTANCE_NAME = "multi_axis_auto_annotator_service"


class MultiAxisAutoAnnotatorClient:
  """Helper class for calling the rpcs in the MultiAxisAutoAnnotatorService."""

  def __init__(
      self,
      stub: multi_axis_auto_annotator_service_pb2_grpc.MultiAxisAutoAnnotatorServiceStub,
      instance_name: str = DEFAULT_SERVICE_INSTANCE_NAME,
  ):
    """Constructor.

    Args:
      stub: The MultiAxisAutoAnnotatorServiceStub.
      instance_name: The service instance name of the auto annotator service.
    """
    self._stub: (
        multi_axis_auto_annotator_service_pb2_grpc.MultiAxisAutoAnnotatorServiceStub
    ) = stub
    self._connection_params = {
        "metadata": [(
            "x-resource-instance-name",
            instance_name,
        )]
    }

  @classmethod
  def connect(
      cls,
      address: str = DEFAULT_SERVICE_ADDRESS,
      instance_name: str = DEFAULT_SERVICE_INSTANCE_NAME,
  ) -> tuple[grpc.Channel, "MultiAxisAutoAnnotatorClient"]:
    """Connects to the auto annotator service.

    Args:
      address: The address of the auto annotator service.
      instance_name: The service instance name of the auto annotator service.

    Returns:
      gRpc channel, grasp annotator client
    """
    logging.info(
        "Connecting to multi_axis_auto_annotator_service at %s", address
    )
    channel = grpc.insecure_channel(address)
    return channel, MultiAxisAutoAnnotatorClient(
        stub=multi_axis_auto_annotator_service_pb2_grpc.MultiAxisAutoAnnotatorServiceStub(
            channel
        ),
        instance_name=instance_name,
    )

  def get_annotations(
      self,
      get_annotations_request: multi_axis_auto_annotator_service_pb2.GetAnnotationsRequest,
  ) -> list[schmalz_multi_axis_grasp_pb2.SchmalzMultiAxisGraspAnnotation]:
    """Get annotations for a given triangle mesh.

    Args:
      get_annotations_request: The parameters used to get annotations.

    Returns:
      The annotations from the auto annotator service.
    """
    response = self._stub.GetAnnotations(
        get_annotations_request,
        **self._connection_params,
    )
    return response.annotations

  def get_annotations_and_commands(
      self,
      get_annotations_and_commands_request: multi_axis_auto_annotator_service_pb2.GetAnnotationsAndCommandsRequest,
  ) -> tuple[
      list[schmalz_multi_axis_grasp_pb2.SchmalzMultiAxisGraspAnnotation],
      list[schmalz_multi_axis_grasp_pb2.SchmalzMultiAxisGraspCommand],
  ]:
    """Get annotations and grasp commands for a given triangle mesh.

    Args:
      get_annotations_and_commands_request: The parameters used to get
        annotations and grasp commands.

    Returns:
      The annotations and grasp commands from the auto annotator service.
    """
    response = self._stub.GetAnnotationsAndCommands(
        get_annotations_and_commands_request,
        **self._connection_params,
    )
    return response.annotations, response.commands

  def show_annotations(
      self,
      show_annotations_request: multi_axis_auto_annotator_service_pb2.ShowAnnotationsRequest,
  ) -> bool:
    """Show annotations for a given part.

    This visualizes the annotations in the world with
    time_between_annotations_in_sec time between spawning annotations. The
    annotations are shown in the order they are provided. The plates and cups
    are deleted after all annotations are shown.

    Args:
      show_annotations_request: The parameters used to show annotations.

    Returns:
      True if the annotations are shown successfully, False otherwise.
    """
    response = self._stub.ShowAnnotations(
        show_annotations_request,
        **self._connection_params,
    )
    return response.success
