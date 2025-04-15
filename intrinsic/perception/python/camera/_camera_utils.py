# Copyright 2023 Intrinsic Innovation LLC

"""Camera misc helper methods."""

from __future__ import annotations

from typing import Mapping, Optional, Tuple

from absl import logging
import grpc
from intrinsic.perception.client.v1.python.camera import _camera_utils as _camera_utils_v1
from intrinsic.perception.proto import camera_config_pb2
from intrinsic.perception.proto import dimensions_pb2
from intrinsic.perception.proto import distortion_params_pb2
from intrinsic.perception.proto import intrinsic_params_pb2
from intrinsic.perception.proto.v1 import camera_config_pb2 as camera_config_v1_pb2
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
  elif camera_driver == "photoneo":
    return identifier.photoneo.device_id
  elif camera_driver == "realsense":
    return identifier.realsense.device_id
  elif camera_driver == "plenoptic_unit":
    return identifier.plenoptic_unit.device_id
  elif camera_driver == "fake_genicam":
    return "fake_genicam"
  else:
    return None


def extract_dimensions(
    dimensions: dimensions_pb2.Dimensions,
) -> Tuple[int, int]:
  """Extract dimensions into a tuple."""
  return (dimensions.cols, dimensions.rows)


def extract_intrinsic_dimensions(
    ip: intrinsic_params_pb2.IntrinsicParams,
) -> Tuple[int, int]:
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
  if dp.k4 or dp.k5 or dp.k6:
    return np.array([dp.k1, dp.k2, dp.p1, dp.p2, dp.k3, dp.k4, dp.k5, dp.k6])
  else:
    return np.array([dp.k1, dp.k2, dp.p1, dp.p2, dp.k3])


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
      camera_config_v1 = camera_config_v1_pb2.CameraConfig()
      proto_utils.unpack_any(config.contents, camera_config_v1)
      camera_config = _camera_utils_v1.from_v1_camera_config(camera_config_v1)
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
