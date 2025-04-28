# Copyright 2023 Intrinsic Innovation LLC

"""Converters from intrinsic math protos to commonly used in-memory representations."""

from intrinsic.math.python import data_types
import numpy as np
from third_party.ros2.ros_interfaces.jazzy.geometry_msgs.msg import point_pb2
from third_party.ros2.ros_interfaces.jazzy.geometry_msgs.msg import pose_pb2
from third_party.ros2.ros_interfaces.jazzy.geometry_msgs.msg import quaternion_pb2
from third_party.ros2.ros_interfaces.jazzy.geometry_msgs.msg import vector3_pb2

_QUATERNION_UNITY_TOLERANCE = np.finfo(np.float64).eps * 32


def ndarray_from_point_proto(point_proto: point_pb2.Point) -> np.ndarray:
  """Convert a point_pb2.Point to a size 3 np.ndarray."""
  return np.array([point_proto.x, point_proto.y, point_proto.z])


def ndarray_to_point_proto(point: np.ndarray) -> point_pb2.Point:
  """Converts a size 3 np.ndarray to a point_pb2.Point.

  Args:
    point: An np.ndarray of size 3.

  Returns:
    A point_pb2.Point.

  Raises:
    ValueError if the input array has not a length of 3.
  """
  if point.shape != (3,):
    raise ValueError(
        'Received point of size {0} but expected a size of 3.'.format(
            point.size
        )
    )
  return point_pb2.Point(x=point[0], y=point[1], z=point[2])


def ndarray_to_vector3_proto(arr3: np.ndarray) -> vector3_pb2.Vector3:
  """Converts a size 3 np.ndarray to a vector3_pb2.Vector3.

  Args:
    arr3: An np.ndarray of size 3.

  Returns:
    A vector3_pb2.Vector3.

  Raises:
    ValueError if the input array has not a length of 3.
  """
  if arr3.shape != (3,):
    raise ValueError(
        'Received point of size {0} but expected a size of 3.'.format(arr3.size)
    )
  return vector3_pb2.Vector3(x=arr3[0], y=arr3[1], z=arr3[2])


def ndarray_from_vector3_proto(
    vector3_proto: vector3_pb2.Vector3,
) -> np.ndarray:
  """Converts a vector3_pb2.Vector3 to a size 3 np.ndarray."""
  return np.array([vector3_proto.x, vector3_proto.y, vector3_proto.z])


def quaternion_from_proto(
    quaternion_proto: quaternion_pb2.Quaternion,
) -> data_types.Quaternion:
  """Convert a quaternion proto to a quaternion."""
  return data_types.Quaternion(
      [
          quaternion_proto.x,
          quaternion_proto.y,
          quaternion_proto.z,
          quaternion_proto.w,
      ],
      normalize=False,
  )


def quaternion_to_proto(
    quat: data_types.Quaternion,
) -> quaternion_pb2.Quaternion:
  """Convert a quaternion to a quaternion proto."""
  return quaternion_pb2.Quaternion(x=quat.x, y=quat.y, z=quat.z, w=quat.w)


def pose_from_proto(
    pose_proto: pose_pb2.Pose, normalize_quaternion: bool = False
) -> data_types.Pose3:
  """Convert a pose proto to a pose."""
  point = ndarray_from_point_proto(pose_proto.position)
  rotation = data_types.Rotation3(
      quat=quaternion_from_proto(pose_proto.orientation),
      normalize=normalize_quaternion,
  )
  # We expect the quaternion in the input proto to be normalized. Pose3 does not
  # require this so we check this explicitly.
  rotation.quaternion.check_normalized(_QUATERNION_UNITY_TOLERANCE)
  return data_types.Pose3(translation=point, rotation=rotation)


def pose_to_proto(pose: data_types.Pose3) -> pose_pb2.Pose:
  """Convert a pose to a pose proto."""
  msg = pose_pb2.Pose()
  msg.position.CopyFrom(ndarray_to_point_proto(pose.translation))
  # Normalize quaternion if this is not already the case (don't re-normalize and
  # introduce numerical variations). Pose3 may contain a non-unit quaternion.
  quat = pose.quaternion
  if not quat.is_normalized():
    quat = quat.normalize()
  msg.orientation.CopyFrom(quaternion_to_proto(quat))
  return msg
