# Copyright 2023 Intrinsic Innovation LLC

"""Base camera class wrapping gRPC connection and calls."""

from __future__ import annotations

import datetime
from typing import Optional

import grpc
from intrinsic.perception.proto import camera_config_pb2
from intrinsic.perception.proto import camera_settings_pb2
from intrinsic.perception.proto import capture_result_pb2
from intrinsic.perception.service.proto import camera_server_pb2
from intrinsic.perception.service.proto import camera_server_pb2_grpc
from intrinsic.util.grpc import connection
from intrinsic.util.grpc import error_handling
from intrinsic.util.grpc import interceptor


class CameraClient:
  """Base camera class wrapping gRPC connection and calls.

  Skill users should use the `Camera` class, which provides a more pythonic
  interface.
  """

  camera_config: camera_config_pb2.CameraConfig
  _camera_stub: camera_server_pb2_grpc.CameraServerStub
  _camera_handle: Optional[str]

  def __init__(
      self,
      camera_channel: grpc.Channel,
      connection_params: connection.ConnectionParams,
      camera_config: camera_config_pb2.CameraConfig,
  ):
    """Creates a CameraClient object."""
    self.camera_config = camera_config
    self._camera_handle = None

    # Create stub.
    intercepted_camera_channel = grpc.intercept_channel(
        camera_channel,
        interceptor.HeaderAdderInterceptor(connection_params.headers),
    )
    self._camera_stub = camera_server_pb2_grpc.CameraServerStub(
        intercepted_camera_channel
    )

  @property
  def created(self) -> bool:
    """Returns whether the camera client is created."""
    return self._camera_handle is not None

  def reset(
      self,
      camera_config: camera_config_pb2.CameraConfig,
  ):
    """Resets the camera client to be created on the next API call."""
    self.camera_config = camera_config
    self._camera_handle = None

  def create_camera(
      self,
      camera_config: camera_config_pb2.CameraConfig,
      timeout: Optional[datetime.timedelta] = None,
      deadline: Optional[datetime.datetime] = None,
  ) -> str:
    """Creates a camera instance."""
    deadline = deadline or (
        datetime.datetime.now() + timeout if timeout is not None else None
    )
    self.reset(camera_config)
    return self._create_camera(deadline)

  def _create_camera(
      self,
      deadline: Optional[datetime.datetime] = None,
  ) -> str:
    """Initializes the camera handle."""
    handle = self._create_camera_with_retry(self.camera_config, deadline)
    if not handle:
      raise grpc.RpcError(
          grpc.StatusCode.FAILED_PRECONDITION, "Could not create camera handle."
      )
    self._camera_handle = handle
    return handle

  @error_handling.retry_on_grpc_unavailable
  def _create_camera_with_retry(
      self,
      camera_config: camera_config_pb2.CameraConfig,
      deadline: Optional[datetime.datetime] = None,
  ) -> Optional[str]:
    """Creates a camera instance."""
    if deadline is not None:
      timeout = deadline - datetime.datetime.now()
      if timeout <= datetime.timedelta(seconds=0):
        raise grpc.RpcError(grpc.StatusCode.DEADLINE_EXCEEDED)
    request = camera_server_pb2.CreateCameraRequest(camera_config=camera_config)
    response = self._camera_stub.CreateCamera(request)
    return response.camera_handle or None

  def describe_camera(
      self,
  ) -> camera_server_pb2.DescribeCameraResponse:
    """Describes the camera and its sensors. Enumerates connected sensors.

    Returns:
      A camera_server_pb2.DescribeCameraResponse with the camera's config and
      sensor information.

    Raises:
      grpc.RpcError: A gRPC error occurred.
    """
    if self._camera_handle is None:
      self._create_camera()

    request = camera_server_pb2.DescribeCameraRequest(
        camera_handle=self._camera_handle
    )
    response = self._camera_stub.DescribeCamera(request)
    return response

  def capture(
      self,
      timeout: Optional[datetime.timedelta] = None,
      deadline: Optional[datetime.datetime] = None,
      sensor_ids: Optional[list[int]] = None,
      skip_undistortion: bool = False,
  ) -> capture_result_pb2.CaptureResult:
    """Captures image data from the requested sensors of the specified camera.

    Args:
      timeout: Optional. The timeout which is used for retrieving frames from
        the underlying driver implementation. If this timeout is implemented by
        the underlying camera driver, it will not spend more than the specified
        time when waiting for new frames. The timeout should be greater than the
        combined exposure and processing time. Processing times can be roughly
        estimated as a value between 10 - 50 ms. The timeout just serves as an
        upper limit to prevent blocking calls within the camera driver. In case
        of intermittent network errors users can try to increase the timeout.
        The default timeout (if unspecified) of 500 ms works well in common
        setups.
      deadline: Optional. The deadline corresponding to the timeout. This takes
        priority over the timeout.
      sensor_ids: Optional. Request data only for the following sensor ids (i.e.
        transmit mask). Empty returns all sensor images.
      skip_undistortion: Whether to skip undistortion.

    Returns:
      A capture_result_pb2.CaptureResult with the requested sensor images.

    Raises:
      grpc.RpcError: A gRPC error occurred.
    """
    deadline = deadline or (
        datetime.datetime.now() + timeout if timeout is not None else None
    )
    sensor_ids = sensor_ids or []

    if self._camera_handle is None:
      self._create_camera(deadline=deadline)
    return self._capture(deadline, sensor_ids, skip_undistortion)

  @error_handling.retry_on_grpc_unavailable
  def _capture(
      self,
      deadline: Optional[datetime.datetime],
      sensor_ids: list[int],
      skip_undistortion: bool,
  ) -> capture_result_pb2.CaptureResult:
    """Captures image data from the requested sensors of the specified camera."""
    timeout = None
    if deadline is not None:
      timeout = deadline - datetime.datetime.now()
      if timeout <= datetime.timedelta(seconds=0):
        raise grpc.RpcError(grpc.StatusCode.DEADLINE_EXCEEDED)
    request = camera_server_pb2.CaptureRequest(
        camera_handle=self._camera_handle
    )
    if timeout is not None:
      request.timeout.FromTimedelta(timeout)
    request.sensor_ids[:] = sensor_ids
    request.post_processing.skip_undistortion = skip_undistortion
    if timeout is not None:
      response, _ = self._camera_stub.Capture.with_call(
          request,
          timeout=timeout.seconds,
      )
    else:
      response = self._camera_stub.Capture(request)
    return response.capture_result

  def read_camera_setting_properties(
      self,
      name: str,
  ) -> camera_settings_pb2.CameraSettingProperties:
    """Read the properties of the setting.

    The function returns an error if the setting is not supported. If specific
    properties of a setting are not supported, they are not added to the result.
    The function only returns existing properties and triggers no errors for
    non-existing properties as these are optional to be implemented by the
    camera vendors.

    Args:
      name: The setting name. The setting name must be defined by the Standard
        Feature Naming Conventions (SFNC) which is part of the GenICam standard.

    Returns:
      A camera_settings_pb2.CameraSettingProperties with the requested setting
      properties.

    Raises:
      grpc.RpcError: A gRPC error occurred.
    """
    if self._camera_handle is None:
      self._create_camera()

    request = camera_server_pb2.ReadCameraSettingPropertiesRequest(
        camera_handle=self._camera_handle,
        name=name,
    )
    response = self._camera_stub.ReadCameraSettingProperties(request)
    return response.properties

  def read_camera_setting(
      self,
      name: str,
  ) -> camera_settings_pb2.CameraSetting:
    """Reads and returns the current value of a specific setting from a camera.

    The function returns an error if the setting is not supported.

    Args:
      name: The setting name. The setting name must be defined by the Standard
        Feature Naming Conventions (SFNC) which is part of the GenICam standard.

    Returns:
      A camera_settings_pb2.CameraSetting with the requested setting value.

    Raises:
      grpc.RpcError: A gRPC error occurred.
    """
    if self._camera_handle is None:
      self._create_camera()

    request = camera_server_pb2.ReadCameraSettingRequest(
        camera_handle=self._camera_handle,
        name=name,
    )
    response = self._camera_stub.ReadCameraSetting(request)
    return response.setting

  def update_camera_setting(
      self,
      setting: camera_settings_pb2.CameraSetting,
  ) -> None:
    """Update the value of a specific camera setting.

    The function returns an error if the setting is not supported.
    Note: When updating camera parameters, beware that the
    modifications will apply to all instances. I.e. it will also affect all
    other clients who are using the same camera.

    Args:
      setting: A camera_settings_pb2.CameraSetting with a value to update to.

    Raises:
      grpc.RpcError: A gRPC error occurred.
    """
    if self._camera_handle is None:
      self._create_camera()

    request = camera_server_pb2.UpdateCameraSettingRequest(
        camera_handle=self._camera_handle,
        setting=setting,
    )
    self._camera_stub.UpdateCameraSetting(request)
