# Copyright 2023 Intrinsic Innovation LLC

"""Base camera class wrapping gRPC connection and calls."""

from __future__ import annotations

import datetime
from typing import Any, Optional, cast

import grpc
from intrinsic.perception.proto.v1 import camera_config_pb2
from intrinsic.perception.proto.v1 import camera_identifier_pb2
from intrinsic.perception.proto.v1 import camera_service_pb2
from intrinsic.perception.proto.v1 import camera_service_pb2_grpc
from intrinsic.perception.proto.v1 import camera_settings_pb2
from intrinsic.perception.proto.v1 import capture_result_pb2
from intrinsic.util.grpc import connection
from intrinsic.util.grpc import error_handling
from intrinsic.util.grpc import interceptor


class CameraClient:
  """Base camera class wrapping gRPC connection and calls.

  Skill users should use the `Camera` class, which provides a more pythonic
  interface.
  """

  camera_identifier: camera_identifier_pb2.CameraIdentifier
  _camera_stub: camera_service_pb2_grpc.CameraServiceStub

  def __init__(
      self,
      camera_channel: grpc.Channel,
      connection_params: connection.ConnectionParams,
      camera_identifier: camera_identifier_pb2.CameraIdentifier,
  ):
    """Creates a CameraClient object."""
    self.camera_identifier = camera_identifier

    # Create stub.
    intercepted_camera_channel = grpc.intercept_channel(
        camera_channel,
        interceptor.HeaderAdderInterceptor(connection_params.headers),
    )
    self._camera_stub = camera_service_pb2_grpc.CameraServiceStub(
        intercepted_camera_channel
    )

  @error_handling.retry_on_grpc_unavailable
  def describe_camera(
      self,
  ) -> camera_service_pb2.DescribeCameraResponse:
    """Enumerates connected sensors.

    Returns:
      A camera_server_pb2.DescribeCameraResponse with the camera's sensor
      information.

    Raises:
      grpc.RpcError: A gRPC error occurred.
    """
    request = camera_service_pb2.DescribeCameraRequest(
        camera_identifier=self.camera_identifier
    )
    response = self._camera_stub.DescribeCamera(request)
    return response

  def capture(
      self,
      camera_config: Optional[camera_config_pb2.CameraConfig] = None,
      timeout: Optional[datetime.timedelta] = None,
      deadline: Optional[datetime.datetime] = None,
      sensor_ids: Optional[list[int]] = None,
      skip_undistortion: bool = False,
  ) -> capture_result_pb2.CaptureResult:
    """Captures image data from the requested sensors of the specified camera.

    Args:
      camera_config: Optional. The camera config to use. If not specified, a
        default camera config will be created from the camera identifier that
        was used to instantiate this camera client.
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
      ValueError: When the identifier in camera_config does not match the
        identifier this client was instantiated with.
    """
    if camera_config is None:
      camera_config = camera_config_pb2.CameraConfig(
          identifier=self.camera_identifier
      )
    elif camera_config.identifier != self.camera_identifier:
      raise ValueError(
          "The identifier in camera_config does not match the identifier this"
          " camera client was instantiated with."
      )

    deadline = deadline or (
        datetime.datetime.now() + timeout if timeout is not None else None
    )
    sensor_ids = sensor_ids or []
    response = self._capture(
        camera_config, deadline, sensor_ids, skip_undistortion
    )
    return response

  @error_handling.retry_on_grpc_unavailable
  def _capture(
      self,
      camera_config: camera_config_pb2.CameraConfig,
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
    request = camera_service_pb2.CaptureRequest(camera_config=camera_config)
    if timeout is not None:
      request.timeout.FromTimedelta(timeout)
    request.sensor_ids[:] = sensor_ids
    if skip_undistortion:
      sensor_ids_skip_undistortion = []
      if not sensor_ids:
        describe_camera_response = self.describe_camera()
        for sensor in describe_camera_response.sensors:
          sensor_ids_skip_undistortion.append(sensor.id)
      else:
        sensor_ids_skip_undistortion = sensor_ids
      for sensor_id in sensor_ids_skip_undistortion:
        request.post_processing_by_sensor_id[sensor_id].skip_undistortion = True
    if timeout is not None:
      response, _ = self._camera_stub.Capture.with_call(
          request,
          timeout=timeout.seconds,
      )
    else:
      response = self._camera_stub.Capture(request)
    return response.capture_result

  @error_handling.retry_on_grpc_unavailable
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
    request = camera_service_pb2.ReadCameraSettingPropertiesRequest(
        camera_identifier=self.camera_identifier,
        name=name,
    )
    response = self._camera_stub.ReadCameraSettingProperties(request)
    return response.properties

  @error_handling.retry_on_grpc_unavailable
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
    request = camera_service_pb2.ReadCameraSettingRequest(
        camera_identifier=self.camera_identifier,
        name=name,
    )
    response = self._camera_stub.ReadCameraSetting(request)
    return response.setting

  @error_handling.retry_on_grpc_unavailable
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
    request = camera_service_pb2.UpdateCameraSettingRequest(
        camera_identifier=self.camera_identifier,
        setting=setting,
    )
    self._camera_stub.UpdateCameraSetting(request)
