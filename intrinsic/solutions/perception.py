# Copyright 2023 Intrinsic Innovation LLC

"""Camera access within the workcell API."""

import datetime
import enum
import math
from typing import Optional, Tuple

import grpc
from intrinsic.perception.client.v1.python.camera import _camera_utils
from intrinsic.perception.client.v1.python.camera import camera_client
from intrinsic.perception.client.v1.python.camera import data_classes
from intrinsic.perception.proto.v1 import image_buffer_pb2
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
  _executive: execution.Executive
  _is_simulated: bool

  config: data_classes.CameraConfig

  def __init__(
      self,
      channel: grpc.Channel,
      resource_handle: resource_handle_pb2.ResourceHandle,
      executive: execution.Executive,
      is_simulated: bool,
  ):
    """Creates a Camera object.

    During construction the camera is not yet open. Opening the camera on the
    camera server will happen once it is needed as this function only requests
    the current camera config from the resource registry. The camera will be
    created during the first capture() call.

    Args:
      channel: The grpc channel to the respective camera server.
      resource_handle: Resource handle for the camera.
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

    self.config = data_classes.CameraConfig(camera_config)
    grpc_info = resource_handle.connection_info.grpc
    connection_params = connection.ConnectionParams(
        grpc_info.address, grpc_info.server_instance, grpc_info.header
    )
    self._client = camera_client.CameraClient(
        channel,
        connection_params,
        camera_config.identifier,
    )

    self._resource_handle = resource_handle
    self._executive = executive
    self._is_simulated = is_simulated

  @property
  def _resource_name(self) -> str:
    """Returns the resource id of the camera."""
    return self._resource_handle.name

  def capture(
      self,
      camera_config: Optional[data_classes.CameraConfig] = None,
      timeout: datetime.timedelta = datetime.timedelta(
          seconds=_MAX_FRAME_WAIT_TIME_SECONDS
      ),
      sensor_ids: Optional[list[int]] = None,
      skip_undistortion: bool = False,
  ) -> data_classes.CaptureResult:
    """Performs grpc request to capture sensor images from the camera.

    Args:
      camera_config: Optional. The camera config to use. If not specified, the
        config that was used to instantiate this class will be used. This can be
        the camera config which was retrieved from the resource handle.
      timeout: Timeout duration for Capture() service calls.
      sensor_ids: List of selected sensor identifiers for Capture() service
        calls.
      skip_undistortion: Whether to skip undistortion.

    Returns:
      The acquired list of sensor images.

    Raises:
      ValueError: The identifier in camera_config does not match the identifier
        this camera was instantiated with.
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
    if camera_config is None:
      camera_config = self.config
    elif camera_config.identifier != self.config.identifier:
      raise ValueError(
          'The identifier in camera_config does not match the identifier found'
          ' in the resource handle this camera was instantiated with.'
      )

    deadline = datetime.datetime.now() + timeout
    sensor_ids = sensor_ids or []
    capture_result = self._client.capture(
        camera_config=camera_config.proto,
        timeout=timeout,
        deadline=deadline,
        sensor_ids=sensor_ids,
        skip_undistortion=skip_undistortion,
    )
    return data_classes.CaptureResult(capture_result)

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


def _create_cameras(
    resource_registry: resource_registry_client.ResourceRegistryClient,
    grpc_channel: grpc.Channel,
    executive: execution.Executive,
    is_simulated: bool,
) -> dict[str, Camera]:
  """Creates cameras for each resource handle that is a camera.

  Please note that by calling this function the cameras are not opened on the
  camera server yet. They will be created during the first capture() call.

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
