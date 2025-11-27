# Copyright 2023 Intrinsic Innovation LLC

"""Convenience class for Camera use within skills."""

from __future__ import annotations

from collections.abc import Mapping
import datetime
from typing import cast
from typing import Optional
from typing import Union
import warnings

from absl import logging
from google.protobuf import empty_pb2
import grpc
from intrinsic.math.python import pose3
from intrinsic.perception.proto.v1 import settings_pb2
from intrinsic.perception.python.camera import _camera_utils
from intrinsic.perception.python.camera import camera_client
from intrinsic.perception.python.camera import data_classes
from intrinsic.resources.client import resource_registry_client
from intrinsic.resources.proto import resource_handle_pb2
from intrinsic.skills.proto import equipment_pb2
from intrinsic.skills.python import skill_interface
from intrinsic.util.grpc import connection
from intrinsic.world.python import object_world_client
from intrinsic.world.python import object_world_resources
import numpy as np


def make_camera_resource_selector() -> equipment_pb2.ResourceSelector:
  """Creates the default resource selector for a camera equipment slot.

  Used in a skill's `required_equipment` implementation.

  Returns:
    A resource selector that is valid for cameras.
  """
  return equipment_pb2.ResourceSelector(
      capability_names=[
          _camera_utils.CAMERA_RESOURCE_CAPABILITY,
      ]
  )


class Camera:
  """Convenience class for Camera use within skills.

  This class provides a more pythonic interface than the `CameraClient` which
  wraps the gRPC calls for interacting with cameras.

  Typical usage example:

  - Add a camera slot to the skill, e.g.:
    ```
    @classmethod
    @overrides(skl.Skill)
    def required_equipment(cls) -> Mapping[str,
    equipment_pb2.ResourceSelector]:
      # create a camera equipment slot for the skill
      return {
          camera_slot: cameras.make_camera_resource_selector()
      }
    ```
  - Create and use a camera in the skill:
    ```
    def execute(
        self, request: skl.ExecuteRequest, context: skl.ExecuteContext
    ) -> ...:
    ...

    # access the camera equipment slot added in `required_equipment`
    camera = cameras.Camera.create(context, "camera_slot")

    # get the camera's intrinsic matrix as a numpy array
    intrinsic_matrix = camera.intrinsic_matrix()

    # capture from the camera's primary sensor
    sensor_image = camera.capture()
    # access image buffer as a numpy array
    img = sensor_image.array

    # or capture from all of the camera's currently configured sensors
    capture_result = camera.multi_sensor_capture()
    for sensor_name, sensor_image in capture_result.sensor_images.items():
      pass  # access each sensor's image buffer using sensor_image.array
    ```
    ...
  """

  _client: camera_client.CameraClient
  _resource_registry: Optional[resource_registry_client.ResourceRegistryClient]
  _world_client: Optional[object_world_client.ObjectWorldClient]
  _resource_handle: resource_handle_pb2.ResourceHandle
  _world_object: Optional[object_world_resources.WorldObject]
  _sensor_id_to_name: Mapping[int, str]

  config: data_classes.CameraConfig
  factory_config: Optional[data_classes.CameraConfig]
  factory_sensor_info: Mapping[str, data_classes.SensorInformation]

  @classmethod
  def create(
      cls,
      context: skill_interface.ExecuteContext,
      slot: str,
  ) -> Camera:
    """Creates a Camera object from the skill's execution context.

    Args:
      context: The skill's current skill_interface.ExecuteContext.
      slot: The camera slot created in skill's required_equipment
        implementation.

    Returns:
      A connected Camera object with sensor information cached.
    """
    resource_handle = context.resource_handles[slot]
    world_client = context.object_world
    return cls.create_from_resource_handle(
        resource_handle=resource_handle,
        world_client=world_client,
    )

  @classmethod
  def create_from_resource_registry(
      cls,
      resource_registry: resource_registry_client.ResourceRegistryClient,
      resource_name: str,
      world_client: Optional[object_world_client.ObjectWorldClient] = None,
      channel: Optional[grpc.Channel] = None,
      channel_creds: Optional[grpc.ChannelCredentials] = None,
  ) -> Camera:
    """Creates a Camera object from the given resource registry and resource name.

    Args:
      resource_registry: The resource registry client.
      resource_name: The resource name of the camera.
      world_client: Optional. The current world client, for camera pose
        information.
      channel: Optional. The gRPC channel to the camera service.
      channel_creds: Optional. The gRPC channel credentials to use for the
        connection.

    Returns:
      A connected Camera object with sensor information cached. If no object or
      world information is available, an identity pose will be used for
      world_t_camera and all the world update methods will be a no-op.
    """
    resource_handle = resource_registry.get_resource_instance(
        resource_name
    ).resource_handle
    return cls.create_from_resource_handle(
        resource_handle=resource_handle,
        world_client=world_client,
        resource_registry=resource_registry,
        channel=channel,
        channel_creds=channel_creds,
    )

  @classmethod
  def create_from_resource_handle(
      cls,
      resource_handle: resource_handle_pb2.ResourceHandle,
      world_client: Optional[object_world_client.ObjectWorldClient] = None,
      resource_registry: Optional[
          resource_registry_client.ResourceRegistryClient
      ] = None,
      channel: Optional[grpc.Channel] = None,
      channel_creds: Optional[grpc.ChannelCredentials] = None,
  ) -> Camera:
    """Creates a Camera object from the given resource handle.

    Args:
      resource_handle: The resource handle with which to connect to the camera.
      world_client: Optional. The current world client, for camera pose
        information.
      resource_registry: Optional. The resource registry client.
      channel: Optional. The gRPC channel to the camera service.
      channel_creds: Optional. The gRPC channel credentials to use for the
        connection.

    Returns:
      A connected Camera object with sensor information cached. If no object or
      world information is available, an identity pose will be used for
      world_t_camera and all the world update methods will be a no-op.
    """
    if channel is None:
      channel = _camera_utils.initialize_camera_grpc_channel(
          resource_handle,
          channel_creds,
      )
    return cls(
        channel=channel,
        resource_handle=resource_handle,
        resource_registry=resource_registry,
        world_client=world_client,
    )

  def __init__(
      self,
      channel: grpc.Channel,
      resource_handle: resource_handle_pb2.ResourceHandle,
      resource_registry: Optional[
          resource_registry_client.ResourceRegistryClient
      ] = None,
      world_client: Optional[object_world_client.ObjectWorldClient] = None,
  ):
    """Creates a Camera object from the given camera equipment and world.

    Args:
      channel: The gRPC channel to the camera service.
      resource_handle: The resource handle with which to connect to the camera.
      resource_registry: Optional. The resource registry client.
      world_client: Optional. The current world client, for camera pose
        information.

    Raises:
      RuntimeError: The camera's config could not be parsed from the
        resource handle.
    """
    self._resource_registry = resource_registry
    self._world_client = world_client
    self._resource_handle = resource_handle
    self._world_object = (
        self._world_client.get_object(resource_handle)
        if self._world_client
        else None
    )
    self._sensor_id_to_name = {}

    # parse config
    camera_config = _camera_utils.unpack_camera_config(self._resource_handle)
    if not camera_config:
      raise RuntimeError(
          "Could not parse camera config from resource handle: %s."
          % self.display_name
      )
    self.config = data_classes.CameraConfig(camera_config)
    self.factory_config = None
    self.factory_sensor_info = {}

    # create camera client
    grpc_info = resource_handle.connection_info.grpc
    connection_params = connection.ConnectionParams(
        grpc_info.address, grpc_info.server_instance, grpc_info.header
    )
    self._client = camera_client.CameraClient(
        channel, connection_params, camera_config
    )

    # attempt to describe cameras to get factory configurations
    try:
      describe_camera_proto = self._client.describe_camera()

      self.factory_config = data_classes.CameraConfig(
          describe_camera_proto.camera_config
      )
      self.factory_sensor_info = {
          sensor_info.display_name: data_classes.SensorInformation(sensor_info)
          for sensor_info in describe_camera_proto.sensors
      }

      # map sensor_ids to human readable sensor names from camera description
      # for capture result
      self._sensor_id_to_name = {
          sensor_info.sensor_id: sensor_name
          for sensor_name, sensor_info in self.factory_sensor_info.items()
      }
    except grpc.RpcError as e:
      logging.warning("Could not load factory configuration: %s", e)

    if not self.created:
      self._client.create_camera(camera_config)

  def _reinitialize(
      self,
      error: grpc.RpcError,
      deadline: Optional[datetime.datetime] = None,
  ) -> None:
    """Create camera handle from resources."""
    if self._resource_registry is not None:
      self._resource_handle = self._resource_registry.get_resource_instance(
          self.display_name
      ).resource_handle
      camera_config = _camera_utils.unpack_camera_config(
          self._resource_handle.resource_data
      )
      if camera_config is None:
        raise ValueError(
            "CameraConfig not found in resource handle %s" % self.display_name
        ) from error
      self.config = data_classes.CameraConfig(camera_config)
    self._client.create_camera(self.config.proto, deadline=deadline)

  @property
  def created(self) -> bool:
    """Returns whether the camera client is created."""
    return self._client.created

  @property
  def identifier(self) -> Optional[str]:
    """Camera identifier."""
    return self.config.identifier

  @property
  def equipment_name(self) -> str:
    """Deprecated: Use display_name instead.

    Camera equipment name.
    """
    warnings.warn(
        "equipment_name() is deprecated. Use display_name() instead.",
        DeprecationWarning,
        stacklevel=2,
    )
    return self._resource_handle.name

  @property
  def display_name(self) -> str:
    """Camera display name."""
    return self._resource_handle.name

  @property
  def resource_handle(self) -> resource_handle_pb2.ResourceHandle:
    """Camera resource handle."""
    return self._resource_handle

  @property
  def dimensions(self) -> Optional[tuple[int, int]]:
    """Deprecated: Use the sensor_dimensions property instead.

    Camera intrinsic dimensions (width, height).
    """
    warnings.warn(
        "dimensions() is deprecated. Use sensor_dimensions() instead.",
        DeprecationWarning,
        stacklevel=2,
    )
    return self.config.dimensions

  @property
  def sensor_names(self) -> list[str]:
    """List of sensor names."""
    return list(self.factory_sensor_info.keys())

  @property
  def sensor_ids(self) -> list[int]:
    """List of sensor ids."""
    return list(self.sensor_id_to_name.keys())

  @property
  def sensor_id_to_name(self) -> Mapping[int, str]:
    """Mapping of sensor ids to sensor names."""
    return self._sensor_id_to_name

  @property
  def sensor_dimensions(self) -> Mapping[str, tuple[int, int]]:
    """Mapping of sensor name to the sensor's intrinsic dimensions (width, height)."""
    return {
        sensor_name: sensor_info.dimensions
        for sensor_name, sensor_info in self.factory_sensor_info.items()
    }

  def intrinsic_matrix(
      self,
      sensor_name: Optional[str] = None,
  ) -> Optional[np.ndarray]:
    """Get the camera intrinsic matrix or that of a specific sensor (for multi-sensor cameras), falling back to factory settings or the camera intrinsic matrix if intrinsic params are missing from the requested sensor config.

    Args:
      sensor_name: The desired sensor name, or None for the primary sensor
        intrinsic matrix or to fall back to the camera intrinsic matrix
        (deprecated).

    Returns:
      The sensor's intrinsic matrix or the deprecated camera intrinsic matrix or
      None if no intrinsic matrix could be found.
    """
    if sensor_name is None:
      sensor_name = next(iter(self.sensor_names)) if self.sensor_names else None

    if sensor_name is None or sensor_name not in self.factory_sensor_info:
      warnings.warn(
          "Using deprecated camera.intrinsic_matrix() because the primary"
          " sensor could not be found.",
          DeprecationWarning,
          stacklevel=2,
      )
      return self.config.intrinsic_matrix

    sensor_info = self.factory_sensor_info[sensor_name]
    sensor_id = sensor_info.sensor_id
    sensor_config = (
        self.config.sensor_configs[sensor_id]
        if sensor_id in self.config.sensor_configs
        else None
    )

    # If the camera config provides sensor config overrides, return those.
    # Otherwise, fall back to the factory settings if available or finally the
    # deprecated camera config intrinsic matrix.
    if sensor_config is not None and sensor_config.camera_params is not None:
      return sensor_config.camera_params.intrinsic_matrix
    if sensor_info is not None:
      factory_camera_params = sensor_info.factory_camera_params
      if factory_camera_params is not None:
        return factory_camera_params.intrinsic_matrix
    return self.config.intrinsic_matrix

  def distortion_params(
      self,
      sensor_name: Optional[str] = None,
  ) -> Optional[np.ndarray]:
    """Get the camera distortion params or that of a specific sensor (for multi-sensor cameras), falling back to factory settings if distortion params are missing from the sensor config.

    Args:
      sensor_name: The desired sensor name, or None for the primary sensor
        distortion params or to fall back to the camera distortion params
        (deprecated).

    Returns:
      The distortion params (k1, k2, p1, p2, [k3, [k4, k5, k6, [s1, s2, s3, s4,
        [tx, ty]]]]) or None if it couldn't be found.
    """
    if sensor_name is None:
      sensor_name = next(iter(self.sensor_names)) if self.sensor_names else None

    if sensor_name is None or sensor_name not in self.factory_sensor_info:
      warnings.warn(
          "Using deprecated camera.distortion_params() because the primary"
          " sensor could not be found.",
          DeprecationWarning,
          stacklevel=2,
      )
      return self.config.distortion_params

    sensor_info = self.factory_sensor_info[sensor_name]
    sensor_id = sensor_info.sensor_id
    sensor_config = (
        self.config.sensor_configs[sensor_id]
        if sensor_id in self.config.sensor_configs
        else None
    )

    # If the camera config provides sensor config overrides, return those.
    # Otherwise, fall back to the factory settings if available or finally the
    # deprecated camera config distortion params.
    if sensor_config is not None and sensor_config.camera_params is not None:
      return sensor_config.camera_params.distortion_params
    if sensor_info is not None:
      factory_camera_params = sensor_info.factory_camera_params
      if factory_camera_params is not None:
        return factory_camera_params.distortion_params
    return self.config.distortion_params

  @property
  def world_object(self) -> Optional[object_world_resources.WorldObject]:
    """Camera world object."""
    return self._world_object

  @property
  def world_t_camera(self) -> pose3.Pose3:
    """Camera world pose."""
    if self._world_client is None:
      logging.warning("World client is None, returning identity pose.")
      return pose3.Pose3()
    return self._world_client.get_transform(
        node_a=self._world_client.root,
        node_b=self._world_object,
    )

  def camera_t_sensor(self, sensor_name: str) -> Optional[pose3.Pose3]:
    """Get the sensor camera_t_sensor pose, falling back to factory settings if pose is missing from the sensor config.

    Args:
      sensor_name: The desired sensor's name.

    Returns:
      The pose3.Pose3 of the sensor relative to the pose of the camera itself or
        None if it couldn't be found.
    """
    if sensor_name not in self.factory_sensor_info:
      return None
    sensor_info = self.factory_sensor_info[sensor_name]
    sensor_id = sensor_info.sensor_id

    sensor_config = (
        self.config.sensor_configs[sensor_id]
        if sensor_id in self.config.sensor_configs
        else None
    )

    if sensor_config is not None and sensor_config.camera_t_sensor is not None:
      return sensor_config.camera_t_sensor
    elif sensor_info is not None and sensor_info.camera_t_sensor is not None:
      return sensor_info.camera_t_sensor
    else:
      return None

  def world_t_sensor(self, sensor_name: str) -> Optional[pose3.Pose3]:
    """Get the sensor world_t_sensor pose, falling back to factory settings for camera_t_sensor if pose is missing from the sensor config.

    Args:
      sensor_name: The desired sensor's name.

    Returns:
      The pose3.Pose3 of the sensor relative to the pose of the world or None if
        it couldn't be found.
    """
    camera_t_sensor = self.camera_t_sensor(sensor_name)
    if camera_t_sensor is None:
      return None
    return self.world_t_camera.multiply(camera_t_sensor)

  def update_world_t_camera(self, world_t_camera: pose3.Pose3) -> None:
    """Update camera world pose relative to world root.

    Args:
      world_t_camera: The new world_t_camera pose.
    """
    if self._world_client is None:
      return
    self._world_client.update_transform(
        node_a=self._world_client.root,
        node_b=self._world_object,
        a_t_b=world_t_camera,
        node_to_update=self._world_object,
    )

  def update_camera_t_other(
      self,
      other: object_world_resources.TransformNode,
      camera_t_other: pose3.Pose3,
  ) -> None:
    """Update camera world pose relative to another object.

    Args:
      other: The other object.
      camera_t_other: The relative transform.
    """
    if self._world_client is None:
      return
    self._world_client.update_transform(
        node_a=self._world_object,
        node_b=other,
        a_t_b=camera_t_other,
        node_to_update=self._world_object,
    )

  def update_other_t_camera(
      self,
      other: object_world_resources.TransformNode,
      other_t_camera: pose3.Pose3,
  ) -> None:
    """Update camera world pose relative to another object.

    Args:
      other: The other object.
      other_t_camera: The relative transform.
    """
    if self._world_client is None:
      return
    self._world_client.update_transform(
        node_a=other,
        node_b=self._world_object,
        a_t_b=other_t_camera,
        node_to_update=self._world_object,
    )

  def _capture(
      self,
      timeout: Optional[datetime.timedelta] = None,
      deadline: Optional[datetime.datetime] = None,
      sensor_ids: Optional[list[int]] = None,
      skip_undistortion: bool = False,
  ) -> data_classes.CaptureResult:
    """Capture from the camera and return a CaptureResult."""
    deadline = deadline or (
        datetime.datetime.now() + timeout if timeout is not None else None
    )
    try:
      capture_result_proto = self._client.capture(
          timeout=timeout,
          deadline=deadline,
          sensor_ids=sensor_ids,
          skip_undistortion=skip_undistortion,
      )
      return data_classes.CaptureResult(
          capture_result_proto, self._sensor_id_to_name, self.world_t_camera
      )
    except grpc.RpcError as e:
      if cast(grpc.Call, e).code() != grpc.StatusCode.NOT_FOUND:
        raise
      # If the camera was not found, recreate the camera. This can happen when
      # switching between sim/real or when a service restarts.
      self._reinitialize(e, deadline)
      return self._capture(timeout, deadline, sensor_ids, skip_undistortion)

  def capture(
      self,
      sensor_name: Optional[str] = None,
      timeout: Optional[datetime.timedelta] = None,
      skip_undistortion: bool = False,
  ) -> data_classes.SensorImage:
    """Capture from the camera and return a SensorImage from the selected sensor or the primary sensor if None.

    Args:
      sensor_name: An optional sensor name to capture from, if it's available.
        If it is None, the camera's primary sensor will be selected.
      timeout: An optional timeout which is used for retrieving a sensor image
        from the underlying driver implementation. If this timeout is
        implemented by the underlying camera driver, it will not spend more than
        the specified time when waiting for the new sensor image, after which it
        will throw a deadline exceeded error. The timeout should be greater than
        the combined exposure and processing time. Processing times can be
        roughly estimated as a value between 10 - 50 ms. The timeout just serves
        as an upper limit to prevent blocking calls within the camera driver. In
        case of intermittent network errors users can try to increase the
        timeout. The default timeout (if None) of 500 ms works well in common
        setups.
      skip_undistortion: Whether to skip undistortion.

    Returns:
      A SensorImage from the selected sensor.

    Raises:
      ValueError: The matching sensor could not be found or the capture result
        could not be parsed.
      grpc.RpcError: A gRPC error occurred.
    """
    try:
      if sensor_name is not None:
        if not self.factory_sensor_info:
          raise ValueError(
              "No factory sensor info found, cannot find sensor id for"
              f" {sensor_name}"
          )
        if sensor_name not in self.factory_sensor_info:
          raise ValueError(f"Invalid sensor name: {sensor_name}")
        sensor_ids = [self.factory_sensor_info[sensor_name].sensor_id]
      else:
        sensor_ids = None

      capture_result = self._capture(
          timeout=timeout,
          sensor_ids=sensor_ids,
          skip_undistortion=skip_undistortion,
      )
      if sensor_name is not None:
        return capture_result.sensor_images[sensor_name]
      else:
        return next(iter(capture_result.sensor_images.values()))
    except grpc.RpcError as e:
      logging.warning("Could not capture from camera.")
      raise e

  def multi_sensor_capture(
      self,
      sensor_names: Optional[list[str]] = None,
      timeout: Optional[datetime.timedelta] = None,
      skip_undistortion: bool = False,
  ) -> data_classes.CaptureResult:
    """Capture from the camera and return a CaptureResult.

    Args:
      sensor_names: An optional list of sensor names that will be transmitted in
        the response, if data was collected for them. This acts as a mask to
        limit the number of transmitted `SensorImage`s. If it is None or empty,
        all `SensorImage`s will be transferred.
      timeout: An optional timeout which is used for retrieving sensor images
        from the underlying driver implementation. If this timeout is
        implemented by the underlying camera driver, it will not spend more than
        the specified time when waiting for new sensor images, after which it
        will throw a deadline exceeded error. The timeout should be greater than
        the combined exposure and processing time. Processing times can be
        roughly estimated as a value between 10 - 50 ms. The timeout just serves
        as an upper limit to prevent blocking calls within the camera driver. In
        case of intermittent network errors users can try to increase the
        timeout. The default timeout (if None) of 500 ms works well in common
        setups.
      skip_undistortion: Whether to skip undistortion.

    Returns:
      A CaptureResult which contains the selected sensor images.

    Raises:
      ValueError: The matching sensors could not be found or the capture result
        could not be parsed.
      grpc.RpcError: A gRPC error occurred.
    """
    try:
      if sensor_names is not None:
        if not self.factory_sensor_info:
          raise ValueError(
              "No factory sensor info found, cannot find sensor ids for"
              f" {sensor_names}"
          )
        sensor_ids: list[int] = []
        for sensor_name in sensor_names:
          if sensor_name not in self.factory_sensor_info:
            raise ValueError(f"Invalid sensor name: {sensor_name}")
          sensor_id = self.factory_sensor_info[sensor_name].sensor_id
          sensor_ids.append(sensor_id)
      else:
        sensor_ids = None

      return self._capture(
          timeout=timeout,
          sensor_ids=sensor_ids,
          skip_undistortion=skip_undistortion,
      )
    except grpc.RpcError as e:
      logging.warning("Could not capture from camera.")
      raise e

  def read_camera_setting_properties(
      self,
      name: str,
  ) -> Union[
      settings_pb2.FloatSettingProperties,
      settings_pb2.IntegerSettingProperties,
      settings_pb2.EnumSettingProperties,
  ]:
    """Read the properties of a camera setting by name.

    These settings vary for different types of cameras, but generally conform to
    the GenICam Standard Features Naming
    Convention (SFNC):
    https://www.emva.org/wp-content/uploads/GenICam_SFNC_v2_7.pdf.

    Args:
      name: The setting name.

    Returns:
      The setting properties, which can be used to validate that a particular
        setting is supported.

    Raises:
      ValueError: Setting properties type could not be parsed.
      grpc.RpcError: A gRPC error occurred.
    """
    try:
      camera_setting_properties_proto = (
          self._client.read_camera_setting_properties(name=name)
      )

      setting_properties = camera_setting_properties_proto.WhichOneof(
          "setting_properties"
      )
      if setting_properties == "float_properties":
        return camera_setting_properties_proto.float_properties
      elif setting_properties == "integer_properties":
        return camera_setting_properties_proto.integer_properties
      elif setting_properties == "enum_properties":
        return camera_setting_properties_proto.enum_properties
      else:
        raise ValueError(
            f"Could not parse setting_properties: {setting_properties}."
        )
    except grpc.RpcError as e:
      logging.warning("Could not read camera setting properties.")
      raise e

  def read_camera_setting(
      self,
      name: str,
  ) -> Union[int, float, bool, str]:
    """Read a camera setting by name.

    These settings vary for different types of cameras, but generally conform to
    the GenICam Standard Features Naming
    Convention (SFNC):
    https://www.emva.org/wp-content/uploads/GenICam_SFNC_v2_7.pdf.

    Args:
      name: The setting name.

    Returns:
      The current camera setting.

    Raises:
      ValueError: Setting type could not be parsed.
      grpc.RpcError: A gRPC error occurred.
    """
    try:
      camera_setting_proto = self._client.read_camera_setting(name=name)

      value = camera_setting_proto.WhichOneof("value")
      if value == "integer_value":
        return camera_setting_proto.integer_value
      elif value == "float_value":
        return camera_setting_proto.float_value
      elif value == "bool_value":
        return camera_setting_proto.bool_value
      elif value == "string_value":
        return camera_setting_proto.string_value
      elif value == "enumeration_value":
        return camera_setting_proto.enumeration_value
      elif value == "command_value":
        return "command"
      else:
        raise ValueError(f"Could not parse value: {value}.")
    except grpc.RpcError as e:
      logging.warning("Could not read camera setting.")
      raise e

  def update_camera_setting(
      self,
      name: str,
      value: Union[int, float, bool, str],
  ) -> None:
    """Update a camera setting.

    These settings vary for different types of cameras, but generally conform to
    the GenICam Standard Features Naming
    Convention (SFNC):
    https://www.emva.org/wp-content/uploads/GenICam_SFNC_v2_7.pdf.

    Args:
      name: The setting name.
      value: The desired setting value.

    Raises:
      ValueError: Setting type could not be parsed or value doesn't match type.
      grpc.RpcError: A gRPC error occurred.
    """
    try:
      # Cannot get sufficient type information from just
      # `Union[int, float, bool, str]`, so read the setting first and then
      # update its value.
      setting = self._client.read_camera_setting(name=name)
      value_type = setting.WhichOneof("value")
      if value_type == "integer_value":
        if not isinstance(value, int):
          raise ValueError(f"Expected int value for {name} but got '{value}'")
        setting.integer_value = value
      elif value_type == "float_value":
        # allow int values to be casted to float, but not vice versa
        if isinstance(value, int):
          value = float(value)
        if not isinstance(value, float):
          raise ValueError(f"Expected float value for {name} but got '{value}'")
        setting.float_value = value
      elif value_type == "bool_value":
        if not isinstance(value, bool):
          raise ValueError(f"Expected bool value for {name} but got '{value}'")
        setting.bool_value = value
      elif value_type == "string_value":
        if not isinstance(value, str):
          raise ValueError(
              f"Expected string value for {name} but got '{value}'"
          )
        setting.string_value = value
      elif value_type == "enumeration_value":
        if not isinstance(value, str):
          raise ValueError(
              f"Expected enumeration value string for {name} but got '{value}'"
          )
        setting.enumeration_value = value
      elif value_type == "command_value":
        # no need to check value contents
        setting.command_value = empty_pb2.Empty()
      else:
        raise ValueError(f"Could not parse value: {value_type}.")

      self._client.update_camera_setting(setting=setting)
    except grpc.RpcError as e:
      logging.warning("Could not update camera setting.")
      raise e
