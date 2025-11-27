# Copyright 2023 Intrinsic Innovation LLC

"""Convenience class for Camera use within skills."""

from __future__ import annotations

from collections.abc import Mapping
import datetime
from typing import Optional
from typing import Union

from absl import logging
from google.protobuf import empty_pb2
import grpc
import numpy as np

from intrinsic.math.python import pose3
from intrinsic.perception.client.v1.python.camera import _camera_utils
from intrinsic.perception.client.v1.python.camera import camera_client
from intrinsic.perception.client.v1.python.camera import data_classes
from intrinsic.perception.proto.v1 import settings_pb2
from intrinsic.resources.client import resource_registry_client
from intrinsic.resources.proto import resource_handle_pb2
from intrinsic.skills.proto import equipment_pb2
from intrinsic.skills.python import skill_interface
from intrinsic.util.grpc import connection
from intrinsic.world.python import object_world_client
from intrinsic.world.python import object_world_resources


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
    camera = cameras.Camera.create(context, 'camera_slot')

    # get the camera's intrinsic matrix (for a particular sensor) as a numpy
    # array
    intrinsic_matrix = camera.intrinsic_matrix('sensor_name')

    # capture from all of the camera's currently configured sensors
    capture_result = camera.capture()
    for sensor_name, sensor_image in capture_result.sensor_images.items():
      pass  # access each sensor's image buffer using sensor_image.array
    ```
    ...
  """

  _client: camera_client.CameraClient
  _world_client: Optional[object_world_client.ObjectWorldClient]
  _resource_handle: resource_handle_pb2.ResourceHandle
  _world_object: Optional[object_world_resources.WorldObject]
  _sensor_id_to_name: Mapping[int, str]

  config: data_classes.CameraConfig
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
        channel=channel,
        channel_creds=channel_creds,
    )

  @classmethod
  def create_from_resource_handle(
      cls,
      resource_handle: resource_handle_pb2.ResourceHandle,
      world_client: Optional[object_world_client.ObjectWorldClient] = None,
      channel: Optional[grpc.Channel] = None,
      channel_creds: Optional[grpc.ChannelCredentials] = None,
  ) -> Camera:
    """Creates a Camera object from the given resource handle.

    Args:
      resource_handle: The resource handle with which to connect to the camera.
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
    if channel is None:
      channel = _camera_utils.initialize_camera_grpc_channel(
          resource_handle,
          channel_creds,
      )
    return cls(
        channel=channel,
        resource_handle=resource_handle,
        world_client=world_client,
    )

  def __init__(
      self,
      channel: grpc.Channel,
      resource_handle: resource_handle_pb2.ResourceHandle,
      world_client: Optional[object_world_client.ObjectWorldClient] = None,
  ):
    """Creates a Camera object from the given camera equipment and world.

    Args:
      channel: The gRPC channel to the camera service.
      resource_handle: The resource handle with which to connect to the camera.
      world_client: Optional. The current world client, for camera pose
        information.

    Raises:
      RuntimeError: The camera's config could not be parsed from the
        resource handle.
    """
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
    self.factory_sensor_info = {}

    # create camera client
    grpc_info = resource_handle.connection_info.grpc
    connection_params = connection.ConnectionParams(
        grpc_info.address, grpc_info.server_instance, grpc_info.header
    )
    self._client = camera_client.CameraClient(
        channel, connection_params, camera_config.identifier
    )

    # attempt to describe cameras to get factory configurations
    try:
      describe_camera_proto = self._client.describe_camera()
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

  @property
  def identifier(self) -> Optional[str]:
    """Camera identifier."""
    return self.config.identifier

  @property
  def display_name(self) -> str:
    """Camera display name."""
    return self._resource_handle.name

  @property
  def resource_handle(self) -> resource_handle_pb2.ResourceHandle:
    """Camera resource handle."""
    return self._resource_handle

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

  def intrinsic_matrix(self, sensor_name: str) -> Optional[np.ndarray]:
    """Get the intrinsic matrix of a specific sensor (for multi-sensor cameras), falling back to factory settings if the intrinsic params are missing from the sensor config.

    Args:
      sensor_name: The desired sensor name.

    Returns:
      The sensor's intrinsic matrix or None if it could not be found.
    """
    sensor_info = self.factory_sensor_info.get(sensor_name)
    if sensor_info is None:
      return None

    sensor_id = sensor_info.sensor_id
    sensor_config = (
        self.config.sensor_configs[sensor_id]
        if sensor_id in self.config.sensor_configs
        else None
    )

    if sensor_config is not None and sensor_config.intrinsic_matrix is not None:
      return sensor_config.intrinsic_matrix
    elif (
        sensor_info is not None
        and sensor_info.factory_intrinsic_matrix is not None
    ):
      return sensor_info.factory_intrinsic_matrix
    else:
      return None

  def distortion_params(self, sensor_name: str) -> Optional[np.ndarray]:
    """Get the distortion params of a specific sensor (for multi-sensor cameras), falling back to factory settings if distortion params are missing from the sensor config.

    Args:
      sensor_name: The desired sensor name.

    Returns:
      The distortion params (k1, k2, p1, p2, [k3, [k4, k5, k6, [s1, s2, s3, s4,
        [tx, ty]]]]) or None if it couldn't be found.
    """
    sensor_info = self.factory_sensor_info.get(sensor_name)
    if sensor_info is None:
      return None

    sensor_id = sensor_info.sensor_id
    sensor_config = (
        self.config.sensor_configs[sensor_id]
        if sensor_id in self.config.sensor_configs
        else None
    )

    if (
        sensor_config is not None
        and sensor_config.distortion_params is not None
    ):
      return sensor_config.distortion_params
    elif (
        sensor_info is not None
        and sensor_info.factory_distortion_params is not None
    ):
      return sensor_info.factory_distortion_params
    else:
      return None

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
    sensor_info = self.factory_sensor_info.get(sensor_name)
    if sensor_info is None:
      return None

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
      camera_config: Optional[data_classes.CameraConfig] = None,
      timeout: Optional[datetime.timedelta] = None,
      deadline: Optional[datetime.datetime] = None,
      sensor_ids: Optional[list[int]] = None,
      skip_undistortion: bool = False,
  ) -> data_classes.CaptureResult:
    """Capture from the camera and return a CaptureResult."""
    deadline = deadline or (
        datetime.datetime.now() + timeout if timeout is not None else None
    )
    camera_config_proto = (
        camera_config.proto if camera_config is not None else None
    )
    capture_result_proto = self._client.capture(
        camera_config=camera_config_proto,
        timeout=timeout,
        deadline=deadline,
        sensor_ids=sensor_ids,
        skip_undistortion=skip_undistortion,
    )
    return data_classes.CaptureResult(
        capture_result_proto, self._sensor_id_to_name, self.world_t_camera
    )

  def capture(
      self,
      camera_config: Optional[data_classes.CameraConfig] = None,
      sensor_names: Optional[list[str]] = None,
      timeout: Optional[datetime.timedelta] = None,
      skip_undistortion: bool = False,
  ) -> data_classes.CaptureResult:
    """Capture from the camera and return a CaptureResult.

    Args:
      camera_config: Optional. The camera config to use. If not specified, the
        config that was used to instantiate this class will be used. This can be
        the camera config which was retrieved from the skill context or the
        resource handle.
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
      ValueError: The identifier in camera_config does not match the identifier
        this camera was instantiated with.
      grpc.RpcError: A gRPC error occurred.
    """
    if camera_config is None:
      camera_config = self.config
    elif camera_config.identifier != self.config.identifier:
      raise ValueError(
          "The identifier in camera_config does not match the identifier found"
          " in the resource handle this camera was instantiated with."
      )
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
          camera_config=camera_config,
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
