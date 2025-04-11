# Copyright 2023 Intrinsic Innovation LLC

"""Camera access within the workcell API."""

import datetime
import enum
import math
from typing import Optional, Tuple, cast

import grpc
from intrinsic.perception.proto import image_buffer_pb2
from intrinsic.perception.python.camera import _camera_utils
from intrinsic.perception.python.camera import camera_client
from intrinsic.perception.python.camera import data_classes
from intrinsic.resources.client import resource_registry_client
from intrinsic.resources.proto import resource_handle_pb2
from intrinsic.solutions import deployments
from intrinsic.solutions import execution
from intrinsic.solutions import utils
from intrinsic.util.grpc import connection
import matplotlib.pyplot as plt


# On the guitar cluster grabbing a frame can take more than 7s; if another
# camera is rendering as well this time can go up 16s. To be on the safe side
# the timeout is set to a high value.
# Since frame grabbing tends to timeout on overloaded guitar clusters we
# increase this value even further.
_MAX_FRAME_WAIT_TIME_SECONDS = 120
_PLOT_WIDTH_INCHES = 40
_PLOT_HEIGHT_INCHES = 20


@utils.protoenum(
    proto_enum_type=image_buffer_pb2.Encoding,
    unspecified_proto_enum_map_to_none=image_buffer_pb2.Encoding.ENCODING_UNSPECIFIED,
    strip_prefix='ENCODING_',
)
class ImageEncoding(enum.Enum):
  """Represents the encoding of an image."""


class Camera:
  """Convenience wrapper for Camera."""

  _client: camera_client.CameraClient
  _resource_handle: resource_handle_pb2.ResourceHandle
  _resource_id: str
  _resource_registry: resource_registry_client.ResourceRegistryClient
  _executive: execution.Executive
  _is_simulated: bool

  def __init__(
      self,
      channel: grpc.Channel,
      resource_handle: resource_handle_pb2.ResourceHandle,
      resource_registry: resource_registry_client.ResourceRegistryClient,
      executive: execution.Executive,
      is_simulated: bool,
  ):
    """Creates a Camera object.

    During construction the camera is not yet open. Opening the camera on the
    camera server will happen once it's needed: first, the current camera
    config will be requested from the resource registry, then, the camera
    handle will be created on the camera server.

    Args:
      channel: The grpc channel to the respective camera server.
      resource_handle: Resource handle for the camera.
      resource_registry: Resource registry to fetch camera resources from.
      executive: The executive for checking the state.
      is_simulated: Whether or not the world is being simulated.

    Raises:
      RuntimeError: The camera's config could not be parsed from the
        resource handle.
    """
    camera_config = _camera_utils.unpack_camera_config(resource_handle)
    if not camera_config:
      raise RuntimeError(
          'Could not parse camera config from resource handle: %s.'
          % resource_handle.name
      )

    grpc_info = resource_handle.connection_info.grpc
    connection_params = connection.ConnectionParams(
        grpc_info.address, grpc_info.server_instance, grpc_info.header
    )
    self._client = camera_client.CameraClient(
        channel,
        connection_params,
        camera_config,
    )

    self._resource_handle = resource_handle
    self._resource_registry = resource_registry
    self._executive = executive
    self._is_simulated = is_simulated

  @property
  def _resource_name(self) -> str:
    """Returns the resource id of the camera."""
    return self._resource_handle.name

  def capture(
      self,
      timeout: datetime.timedelta = datetime.timedelta(
          seconds=_MAX_FRAME_WAIT_TIME_SECONDS
      ),
      sensor_ids: Optional[list[int]] = None,
      skip_undistortion: bool = False,
  ) -> data_classes.CaptureResult:
    """Performs grpc request to capture sensor images from the camera.

    If the camera handle is no longer valid (eg, if the server returns a
    NOT_FOUND status), the camera will be reopened (once) on the camera server;
    if re-opening the camera fails an exception is raised.

    Args:
      timeout: Timeout duration for Capture() service calls.
      sensor_ids: List of selected sensor identifiers for Capture() service
        calls.
      skip_undistortion: Whether to skip undistortion.

    Returns:
      The acquired list of sensor images.

    Raises:
      grpc.RpcError from the camera or resource service.
    """

    if self._is_simulated:
      try:
        _ = self._executive.operation
      except execution.OperationNotFoundError:
        print(
            'Note: The image could be showing an outdated simulation state. Run'
            ' `simulation.reset()` to resolve this.'
        )

    deadline = datetime.datetime.now() + timeout
    if not self._client.created:
      self._reinitialize_from_resources(deadline)

    sensor_ids = sensor_ids or []
    try:
      result = self._client.capture(
          timeout=timeout,
          deadline=deadline,
          sensor_ids=sensor_ids,
          skip_undistortion=skip_undistortion,
      )
    except grpc.RpcError as e:
      if cast(grpc.Call, e).code() != grpc.StatusCode.NOT_FOUND:
        raise
      # If the camera was not found, recreate the camera. This can happen when
      # switching between sim/real or when a service restarts.
      self._reinitialize_from_resources(deadline)
      result = self._client.capture(
          timeout=timeout,
          deadline=deadline,
          sensor_ids=sensor_ids,
          skip_undistortion=skip_undistortion,
      )
    return data_classes.CaptureResult(result)

  def show_capture(
      self,
      figsize: Tuple[float, float] = (_PLOT_WIDTH_INCHES, _PLOT_HEIGHT_INCHES),
  ) -> None:
    """Acquires and plots all sensor images from a capture call in a grid plot.

    Args:
      figsize: Size of grid plot. It is defined as a (width, height) tuple with
        the dimensions in inches.
    """
    capture_result = self.capture()
    fig = plt.figure(figsize=figsize)
    nrows = math.ceil(len(capture_result.sensor_images) / 2)
    ncols = 2

    for i, sensor_image in enumerate(capture_result.sensor_images.values()):
      # The first half sensor images are shown on the left side of the plot grid
      # and the second half on the right side.
      if i < nrows:
        fig.add_subplot(nrows, ncols, 2 * i + 1)
      else:
        fig.add_subplot(nrows, ncols, 2 * (i % nrows) + 2)

      if sensor_image.shape[-1] == 1:
        plt.imshow(sensor_image.array, cmap='gray')
      else:
        plt.imshow(sensor_image.array)
      plt.axis('off')
      plt.title(f'Sensor {sensor_image.sensor_id}')

  def _reinitialize_from_resources(self, deadline: datetime.datetime) -> None:
    """Create camera handle from resources."""
    resource_handle = self._resource_registry.get_resource_instance(
        name=self._resource_name,
    ).resource_handle
    camera_config = _camera_utils.unpack_camera_config(resource_handle)
    if camera_config is None:
      raise ValueError(
          'CameraConfig not found in resource handle %s' % self._resource_name
      )
    self._client.create_camera(camera_config, deadline=deadline)
    self._resource_handle = resource_handle


def _create_cameras(
    resource_registry: resource_registry_client.ResourceRegistryClient,
    grpc_channel: grpc.Channel,
    executive: execution.Executive,
    is_simulated: bool,
) -> dict[str, Camera]:
  """Creates cameras for each resource handle that is a camera.

  Please note that the cameras are not opened directly on the camera service.
  The CreateCamera request is delayed until first use of get_frame.

  Args:
    resource_registry: Resource registry to fetch camera resources from.
    grpc_channel: Channel to the camera service.
    executive: The executive for checking the state.
    is_simulated: Whether or not the world is being simulated.

  Returns:
    A dict with camera handles keyed by camera name.

  Raises:
      status.StatusNotOk: If the grpc request failed (propagates grpc error).
  """
  cameras = {}
  for resource_handle in resource_registry.list_all_resource_handles():
    if _camera_utils.unpack_camera_config(resource_handle) is None:
      continue

    cameras[resource_handle.name] = Camera(
        channel=grpc_channel,
        resource_handle=resource_handle,
        resource_registry=resource_registry,
        executive=executive,
        is_simulated=is_simulated,
    )
  return cameras


class Cameras:
  """Convenience wrapper for camera access."""

  _cameras: dict[str, Camera]

  def __init__(
      self,
      resource_registry: resource_registry_client.ResourceRegistryClient,
      grpc_channel: grpc.Channel,
      executive: execution.Executive,
      is_simulated: bool,
  ):
    """Initializes camera handles for all camera resources.

    Note that grpc calls are performed in this constructor.

    Args:
      resource_registry: Resource registry to fetch camera resources from.
      grpc_channel: Channel to the camera grpc service.
      executive: The executive for checking the state.
      is_simulated: Whether or not the world is being simulated.

    Raises:
      status.StatusNotOk: If the grpc request failed (propagates grpc error).
    """
    self._cameras = _create_cameras(
        resource_registry, grpc_channel, executive, is_simulated
    )

  @classmethod
  def for_solution(cls, solution: deployments.Solution) -> 'Cameras':
    """Creates a Cameras instance for the given Solution.

    Args:
      solution: The deployed solution.

    Returns:
      The new Cameras instance.
    """
    resource_registry = resource_registry_client.ResourceRegistryClient.connect(
        solution.grpc_channel
    )

    return cls(
        resource_registry=resource_registry,
        grpc_channel=solution.grpc_channel,
        executive=solution.executive,
        is_simulated=solution.is_simulated,
    )

  def __getitem__(self, camera_name: str) -> Camera:
    """Returns camera wrapper for the specified identifier.

    Args:
      camera_name: Unique identifier of the camera.

    Returns:
      A camera wrapper object that contains a handle to the camera.

    Raises:
      KeyError: if there is no camera with available with the given name.
    """
    return self._cameras[camera_name]

  def __getattr__(self, camera_name: str) -> Camera:
    """Returns camera wrapper for the specified identifier.

    Args:
      camera_name: Unique identifier of the camera.

    Returns:
      A camera wrapper object that contains a handle to the camera.

    Raises:
      AttributeError: if there is no camera with available with the given name.
    """
    if camera_name not in self._cameras:
      raise AttributeError(f'Camera {camera_name} is unknown.')
    return self._cameras[camera_name]

  def __len__(self) -> int:
    """Returns the number of cameras."""
    return len(self._cameras)

  def __str__(self) -> str:
    """Concatenates all camera keys into a string."""
    return '\n'.join(self._cameras.keys())

  def __dir__(self) -> list[str]:
    """Lists all cameras by key (sorted)."""
    return sorted(self._cameras.keys())
