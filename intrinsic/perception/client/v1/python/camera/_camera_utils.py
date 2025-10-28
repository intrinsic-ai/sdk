# Copyright 2023 Intrinsic Innovation LLC

"""Camera misc helper methods."""

from __future__ import annotations

from collections.abc import Mapping
from typing import Optional

from absl import logging
import grpc
from intrinsic.perception.proto import camera_config_pb2 as camera_config_v0_pb2
from intrinsic.perception.proto import camera_identifier_pb2 as camera_identifier_v0_pb2
from intrinsic.perception.proto import camera_params_pb2 as camera_params_v0_pb2
from intrinsic.perception.proto import camera_settings_pb2 as camera_settings_v0_pb2
from intrinsic.perception.proto import dimensions_pb2 as dimensions_v0_pb2
from intrinsic.perception.proto import distortion_params_pb2 as distortion_params_v0_pb2
from intrinsic.perception.proto import intrinsic_params_pb2 as intrinsic_params_v0_pb2
from intrinsic.perception.proto import sensor_config_pb2 as sensor_config_v0_pb2
from intrinsic.perception.proto.v1 import camera_config_pb2
from intrinsic.perception.proto.v1 import camera_identifier_pb2
from intrinsic.perception.proto.v1 import camera_params_pb2
from intrinsic.perception.proto.v1 import camera_settings_pb2
from intrinsic.perception.proto.v1 import dimensions_pb2
from intrinsic.perception.proto.v1 import distortion_params_pb2
from intrinsic.perception.proto.v1 import intrinsic_params_pb2
from intrinsic.perception.proto.v1 import sensor_config_pb2
from intrinsic.resources.proto import resource_handle_pb2
from intrinsic.skills.python import proto_utils
import numpy as np

CAMERA_RESOURCE_CAPABILITY = "CameraConfig"


def extract_identifier(config: camera_config_pb2.CameraConfig) -> Optional[str]:
  """Extract the camera identifier from the camera config."""
  # extract device_id from oneof
  identifier = config.identifier
  camera_driver = identifier.WhichOneof("drivers")
  if camera_driver == "genicam":
    return identifier.genicam.device_id
  else:
    return None


def extract_dimensions(
    dimensions: dimensions_pb2.Dimensions,
) -> tuple[int, int]:
  """Extract dimensions into a tuple."""
  return (dimensions.cols, dimensions.rows)


def extract_intrinsic_dimensions(
    ip: intrinsic_params_pb2.IntrinsicParams,
) -> tuple[int, int]:
  """Extract dimensions from intrinsic params into a tuple."""
  return extract_dimensions(ip.dimensions)


def extract_intrinsic_matrix(
    ip: intrinsic_params_pb2.IntrinsicParams,
) -> np.ndarray:
  """Extract intrinsic matrix from intrinsic params as a numpy array."""
  return np.array([
      [ip.focal_length_x, 0, ip.principal_point_x],
      [0, ip.focal_length_y, ip.principal_point_y],
      [0, 0, 1],
  ])


def extract_distortion_params(
    dp: distortion_params_pb2.DistortionParams,
) -> np.ndarray:
  """Extract distortion parameters from distortion params as a numpy array."""
  params = [dp.k1, dp.k2, dp.p1, dp.p2]
  if any(dp.HasField(f) for f in ["tx", "ty"]):
    params.extend(
        [dp.k3, dp.k4, dp.k5, dp.k6, dp.s4, dp.s3, dp.s2, dp.s1, dp.tx, dp.ty]
    )
  elif any(dp.HasField(f) for f in ["s4", "s3", "s2", "s1"]):
    params.extend([dp.k3, dp.k4, dp.k5, dp.k6, dp.s4, dp.s3, dp.s2, dp.s1])
  elif any(dp.HasField(f) for f in ["k6", "k5", "k4"]):
    params.extend([dp.k3, dp.k4, dp.k5, dp.k6])
  elif dp.HasField("k3"):
    params.append(dp.k3)
  return np.array(params)


def unpack_camera_config(
    resource_handle: resource_handle_pb2.ResourceHandle,
) -> Optional[camera_config_pb2.CameraConfig]:
  """Returns the camera config from a camera resource handle or None if equipment is not a camera."""
  data: Mapping[str, resource_handle_pb2.ResourceHandle.ResourceData] = (
      resource_handle.resource_data
  )
  config = data.get(CAMERA_RESOURCE_CAPABILITY, None)

  if config is None:
    return None

  try:
    camera_config = camera_config_pb2.CameraConfig()
    proto_utils.unpack_any(config.contents, camera_config)
  except TypeError:
    try:
      camera_config_v0 = camera_config_v0_pb2.CameraConfig()
      proto_utils.unpack_any(config.contents, camera_config_v0)
      camera_config = to_v1_camera_config(camera_config_v0)
    except TypeError as e:
      logging.exception("Failed to unpack camera config: %s", e)
      return None

  return camera_config


def initialize_camera_grpc_channel(
    resource_handle: resource_handle_pb2.ResourceHandle,
    channel_creds: Optional[grpc.ChannelCredentials] = None,
) -> grpc.Channel:
  """Initializes a gRPC channel to the camera service."""
  # use unlimited message size for receiving images (e.g. -1)
  options = [("grpc.max_receive_message_length", -1)]
  grpc_info = resource_handle.connection_info.grpc
  if channel_creds is not None:
    channel = grpc.secure_channel(
        grpc_info.address, channel_creds, options=options
    )
  else:
    channel = grpc.insecure_channel(grpc_info.address, options=options)
  return channel


def from_v1_camera_identifier(
    camera_identifier: camera_identifier_pb2.CameraIdentifier,
) -> camera_identifier_v0_pb2.CameraIdentifier:
  """Converts camera identifier from v1 to v0."""
  camera_identifier_v0 = camera_identifier_v0_pb2.CameraIdentifier()
  if camera_identifier.HasField("genicam"):
    camera_identifier_v0.genicam.device_id = camera_identifier.genicam.device_id
  else:
    raise ValueError(
        "Unsupported camera identifier type: %s"
        % camera_identifier.WhichOneof("drivers")
    )
  return camera_identifier_v0


def to_v1_camera_identifier(
    camera_identifier: camera_identifier_v0_pb2.CameraIdentifier,
) -> camera_identifier_pb2.CameraIdentifier:
  """Converts camera identifier from v0 to v1."""
  camera_identifier_v1 = camera_identifier_pb2.CameraIdentifier()
  if camera_identifier.HasField("genicam"):
    camera_identifier_v1.genicam.device_id = camera_identifier.genicam.device_id
  else:
    raise ValueError(
        "Unsupported camera identifier type: %s"
        % camera_identifier.WhichOneof("drivers")
    )
  return camera_identifier_v1


def from_v1_camera_setting(
    camera_setting: camera_settings_pb2.CameraSetting,
) -> camera_settings_v0_pb2.CameraSetting:
  """Converts CameraSettingProperties from v1 to v0."""
  camera_setting_v0 = camera_settings_v0_pb2.CameraSetting(
      name=camera_setting.name
  )
  if camera_setting.HasField("integer_value"):
    camera_setting_v0.integer_value = camera_setting.integer_value
  elif camera_setting.HasField("float_value"):
    camera_setting_v0.float_value = camera_setting.float_value
  elif camera_setting.HasField("bool_value"):
    camera_setting_v0.bool_value = camera_setting.bool_value
  elif camera_setting.HasField("string_value"):
    camera_setting_v0.string_value = camera_setting.string_value
  elif camera_setting.HasField("enumeration_value"):
    camera_setting_v0.enumeration_value = camera_setting.enumeration_value
  elif camera_setting.HasField("command_value"):
    camera_setting_v0.command_value.CopyFrom(camera_setting.command_value)
  return camera_setting_v0


def to_v1_camera_setting(
    camera_setting: camera_settings_v0_pb2.CameraSetting,
) -> camera_settings_pb2.CameraSetting:
  """Converts CameraSettingProperties from v0 to v1."""
  camera_setting_v1 = camera_settings_pb2.CameraSetting(
      name=camera_setting.name
  )
  if camera_setting.HasField("integer_value"):
    camera_setting_v1.integer_value = camera_setting.integer_value
  elif camera_setting.HasField("float_value"):
    camera_setting_v1.float_value = camera_setting.float_value
  elif camera_setting.HasField("bool_value"):
    camera_setting_v1.bool_value = camera_setting.bool_value
  elif camera_setting.HasField("string_value"):
    camera_setting_v1.string_value = camera_setting.string_value
  elif camera_setting.HasField("enumeration_value"):
    camera_setting_v1.enumeration_value = camera_setting.enumeration_value
  elif camera_setting.HasField("command_value"):
    camera_setting_v1.command_value.CopyFrom(camera_setting.command_value)
  return camera_setting_v1


def from_v1_sensor_config(
    sensor_config: sensor_config_pb2.SensorConfig,
) -> sensor_config_v0_pb2.SensorConfig:
  """Converts SensorConfig from v1 to v0."""
  sensor_config_v0 = sensor_config_v0_pb2.SensorConfig(id=sensor_config.id)
  if sensor_config.HasField("camera_t_sensor"):
    sensor_config_v0.camera_t_sensor.CopyFrom(sensor_config.camera_t_sensor)
  if sensor_config.HasField("camera_params"):
    sensor_config_v0.camera_params.CopyFrom(
        from_v1_camera_params(sensor_config.camera_params)
    )
  return sensor_config_v0


def to_v1_sensor_config(
    sensor_config: sensor_config_v0_pb2.SensorConfig,
) -> sensor_config_pb2.SensorConfig:
  """Converts SensorConfig from v0 to v1."""
  sensor_config_v1 = sensor_config_pb2.SensorConfig(id=sensor_config.id)
  if sensor_config.HasField("camera_t_sensor"):
    sensor_config_v1.camera_t_sensor.CopyFrom(sensor_config.camera_t_sensor)
  if sensor_config.HasField("camera_params"):
    sensor_config_v1.camera_params.CopyFrom(
        to_v1_camera_params(sensor_config.camera_params)
    )
  return sensor_config_v1


def from_v1_dimensions(
    dimensions: dimensions_pb2.Dimensions,
) -> dimensions_v0_pb2.Dimensions:
  """Converts Dimensions from v1 to v0."""
  dimensions_v0 = dimensions_v0_pb2.Dimensions(
      cols=dimensions.cols, rows=dimensions.rows
  )
  return dimensions_v0


def to_v1_dimensions(
    dimensions: dimensions_v0_pb2.Dimensions,
) -> dimensions_pb2.Dimensions:
  """Converts Dimensions from v0 to v1."""
  dimensions_v1 = dimensions_pb2.Dimensions(
      cols=dimensions.cols, rows=dimensions.rows
  )
  return dimensions_v1


def from_v1_intrinsic_params(
    intrinsic_params: intrinsic_params_pb2.IntrinsicParams,
) -> intrinsic_params_v0_pb2.IntrinsicParams:
  """Converts IntrinsicParams from v1 to v0."""
  intrinsic_params_v0 = intrinsic_params_v0_pb2.IntrinsicParams(
      focal_length_x=intrinsic_params.focal_length_x,
      focal_length_y=intrinsic_params.focal_length_y,
      principal_point_x=intrinsic_params.principal_point_x,
      principal_point_y=intrinsic_params.principal_point_y,
  )
  intrinsic_params_v0.dimensions.CopyFrom(
      from_v1_dimensions(intrinsic_params.dimensions)
  )
  return intrinsic_params_v0


def to_v1_intrinsic_params(
    intrinsic_params: intrinsic_params_v0_pb2.IntrinsicParams,
) -> intrinsic_params_pb2.IntrinsicParams:
  """Converts IntrinsicParams from v0 to v1."""
  intrinsic_params_v1 = intrinsic_params_pb2.IntrinsicParams(
      focal_length_x=intrinsic_params.focal_length_x,
      focal_length_y=intrinsic_params.focal_length_y,
      principal_point_x=intrinsic_params.principal_point_x,
      principal_point_y=intrinsic_params.principal_point_y,
  )
  intrinsic_params_v1.dimensions.CopyFrom(
      to_v1_dimensions(intrinsic_params.dimensions)
  )
  return intrinsic_params_v1


def from_v1_distortion_params(
    distortion_params: distortion_params_pb2.DistortionParams,
) -> distortion_params_v0_pb2.DistortionParams:
  """Converts DistortionParams from v1 to v0."""
  distortion_params_v0 = distortion_params_v0_pb2.DistortionParams(
      k1=distortion_params.k1,
      k2=distortion_params.k2,
      p1=distortion_params.p1,
      p2=distortion_params.p2,
  )
  if distortion_params.HasField("k3"):
    distortion_params_v0.k3 = distortion_params.k3
  if distortion_params.HasField("k4"):
    distortion_params_v0.k4 = distortion_params.k4
  if distortion_params.HasField("k5"):
    distortion_params_v0.k5 = distortion_params.k5
  if distortion_params.HasField("k6"):
    distortion_params_v0.k6 = distortion_params.k6
  if distortion_params.HasField("s1"):
    distortion_params_v0.s1 = distortion_params.s1
  if distortion_params.HasField("s2"):
    distortion_params_v0.s2 = distortion_params.s2
  if distortion_params.HasField("s3"):
    distortion_params_v0.s3 = distortion_params.s3
  if distortion_params.HasField("s4"):
    distortion_params_v0.s4 = distortion_params.s4
  if distortion_params.HasField("tx"):
    distortion_params_v0.tx = distortion_params.tx
  if distortion_params.HasField("ty"):
    distortion_params_v0.ty = distortion_params.ty
  return distortion_params_v0


def to_v1_distortion_params(
    distortion_params: distortion_params_v0_pb2.DistortionParams,
) -> distortion_params_pb2.DistortionParams:
  """Converts DistortionParams from v0 to v1."""
  distortion_params_v1 = distortion_params_pb2.DistortionParams(
      k1=distortion_params.k1,
      k2=distortion_params.k2,
      p1=distortion_params.p1,
      p2=distortion_params.p2,
  )
  if distortion_params.HasField("k3"):
    distortion_params_v1.k3 = distortion_params.k3
  if distortion_params.HasField("k4"):
    distortion_params_v1.k4 = distortion_params.k4
  if distortion_params.HasField("k5"):
    distortion_params_v1.k5 = distortion_params.k5
  if distortion_params.HasField("k6"):
    distortion_params_v1.k6 = distortion_params.k6
  if distortion_params.HasField("s1"):
    distortion_params_v1.s1 = distortion_params.s1
  if distortion_params.HasField("s2"):
    distortion_params_v1.s2 = distortion_params.s2
  if distortion_params.HasField("s3"):
    distortion_params_v1.s3 = distortion_params.s3
  if distortion_params.HasField("s4"):
    distortion_params_v1.s4 = distortion_params.s4
  if distortion_params.HasField("tx"):
    distortion_params_v1.tx = distortion_params.tx
  if distortion_params.HasField("ty"):
    distortion_params_v1.ty = distortion_params.ty
  return distortion_params_v1


def from_v1_camera_params(
    camera_params: camera_params_pb2.CameraParams,
) -> camera_params_v0_pb2.CameraParams:
  """Converts CameraParams from v1 to v0."""
  camera_params_v0 = camera_params_v0_pb2.CameraParams()
  camera_params_v0.intrinsic_params.CopyFrom(
      from_v1_intrinsic_params(camera_params.intrinsic_params)
  )
  if camera_params.HasField("distortion_params"):
    camera_params_v0.distortion_params.CopyFrom(
        from_v1_distortion_params(camera_params.distortion_params)
    )
  return camera_params_v0


def to_v1_camera_params(
    camera_params: camera_params_v0_pb2.CameraParams,
) -> camera_params_pb2.CameraParams:
  """Converts CameraParams from v0 to v1."""
  camera_params_v1 = camera_params_pb2.CameraParams()
  camera_params_v1.intrinsic_params.CopyFrom(
      to_v1_intrinsic_params(camera_params.intrinsic_params)
  )
  if camera_params.HasField("distortion_params"):
    camera_params_v1.distortion_params.CopyFrom(
        to_v1_distortion_params(camera_params.distortion_params)
    )
  return camera_params_v1


def from_v1_camera_config(
    camera_config: camera_config_pb2.CameraConfig,
) -> camera_config_v0_pb2.CameraConfig:
  """Converts camera config from v1 to v0."""
  camera_config_v0 = camera_config_v0_pb2.CameraConfig()
  camera_config_v0.identifier.CopyFrom(
      from_v1_camera_identifier(camera_config.identifier)
  )
  for camera_setting in camera_config.camera_settings:
    camera_config_v0.camera_settings.append(
        from_v1_camera_setting(camera_setting)
    )
  for sensor_config in camera_config.sensor_configs:
    camera_config_v0.sensor_configs.append(from_v1_sensor_config(sensor_config))
  return camera_config_v0


def to_v1_camera_config(
    camera_config: camera_config_v0_pb2.CameraConfig,
) -> camera_config_pb2.CameraConfig:
  """Converts camera config from v0 to v1."""
  camera_config_v1 = camera_config_pb2.CameraConfig()
  camera_config_v1.identifier.CopyFrom(
      to_v1_camera_identifier(camera_config.identifier)
  )
  for camera_setting in camera_config.camera_settings:
    camera_config_v1.camera_settings.append(
        to_v1_camera_setting(camera_setting)
    )
  for sensor_config in camera_config.sensor_configs:
    camera_config_v1.sensor_configs.append(to_v1_sensor_config(sensor_config))
  return camera_config_v1
