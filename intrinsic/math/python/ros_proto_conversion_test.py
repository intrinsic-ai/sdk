# Copyright 2023 Intrinsic Innovation LLC

"""Tests for math proto conversion utils."""

from absl.testing import absltest
from absl.testing import parameterized
import numpy as np

from intrinsic.math.python import data_types
from intrinsic.math.python import ros_proto_conversion
from third_party.ros2.ros_interfaces.jazzy.geometry_msgs.msg import point_pb2
from third_party.ros2.ros_interfaces.jazzy.geometry_msgs.msg import pose_pb2
from third_party.ros2.ros_interfaces.jazzy.geometry_msgs.msg import quaternion_pb2
from third_party.ros2.ros_interfaces.jazzy.geometry_msgs.msg import vector3_pb2

_rng = np.random.RandomState(seed=0)


class ROSProtoConversionTest(parameterized.TestCase):
  """Tests conversion of ROS math protos to/from corresponding python classes."""

  def test_quaternion_from_proto(self):
    quat_proto_expected = quaternion_pb2.Quaternion(x=0.1, y=0.2, z=0.3, w=0.4)
    quat = ros_proto_conversion.quaternion_from_proto(quat_proto_expected)
    quat_proto = ros_proto_conversion.quaternion_to_proto(quat)
    self.assertEqual(quat_proto, quat_proto_expected)

  def test_quaternion_to_proto(self):
    quat_expected = data_types.Quaternion.random_unit()
    quat_proto = ros_proto_conversion.quaternion_to_proto(quat_expected)
    quat = ros_proto_conversion.quaternion_from_proto(quat_proto)
    self.assertEqual(quat, quat_expected)

  def test_pose_from_proto_fails_for_non_unit_quaterions(self):
    pose_proto_expected = pose_pb2.Pose(
        position=point_pb2.Point(x=0.1, y=0.2, z=0.3),
        orientation=quaternion_pb2.Quaternion(x=0.4, y=0.5, z=0.6, w=0.7),
    )
    # Can't construct a pose with non-normalized quaternion!
    with self.assertRaises(ValueError):
      ros_proto_conversion.pose_from_proto(pose_proto_expected)

  def test_pose_from_proto_quaternion_normalizing(self):
    pose_proto_expected = pose_pb2.Pose(
        position=point_pb2.Point(x=0.1, y=0.2, z=0.3),
        orientation=quaternion_pb2.Quaternion(x=0.4, y=0.5, z=0.6, w=0.7),
    )

    quaternion = np.array([0.4, 0.5, 0.6, 0.7])
    pose_expected = data_types.Pose3(
        translation=[0.1, 0.2, 0.3],
        rotation=data_types.Rotation3(
            quat=data_types.Quaternion(quaternion / np.linalg.norm(quaternion))
        ),
    )

    normalized_pose = ros_proto_conversion.pose_from_proto(
        pose_proto_expected, normalize_quaternion=True
    )

    self.assertEqual(normalized_pose, pose_expected)

  def test_pose_from_proto(self):
    pose_proto_expected = pose_pb2.Pose(
        position=point_pb2.Point(x=0.1, y=0.2, z=0.3),
        orientation=quaternion_pb2.Quaternion(x=0.0, y=0.0, z=0.0, w=1.0),
    )
    pose = ros_proto_conversion.pose_from_proto(pose_proto_expected)
    pose_proto = ros_proto_conversion.pose_to_proto(pose)
    self.assertEqual(pose_proto, pose_proto_expected)

  def test_pose_to_proto_normalized(self):
    pose_expected = data_types.Pose3(
        translation=np.random.randn(3), rotation=data_types.Rotation3.random()
    )
    pose_proto = ros_proto_conversion.pose_to_proto(pose_expected)
    pose = ros_proto_conversion.pose_from_proto(pose_proto)
    self.assertEqual(pose, pose_expected)

    pose = data_types.Pose3(
        translation=[1, 2, 3],
        rotation=data_types.Rotation3(
            quat=data_types.Quaternion([0.5, -0.5, 0.5, -0.5])
        ),
    )
    pose_proto_expected = pose_pb2.Pose(
        position=point_pb2.Point(x=1, y=2, z=3),
        orientation=quaternion_pb2.Quaternion(x=0.5, y=-0.5, z=0.5, w=-0.5),
    )
    pose_proto = ros_proto_conversion.pose_to_proto(pose)
    self.assertEqual(pose_proto, pose_proto_expected)

  def test_pose_to_proto_not_normalized(self):
    pose = data_types.Pose3(
        translation=[1, 2, 3],
        rotation=data_types.Rotation3(
            quat=data_types.Quaternion([0, 0, 0, 1.1])
        ),
    )
    pose_proto_expected = pose_pb2.Pose(
        position=point_pb2.Point(x=1, y=2, z=3),
        orientation=quaternion_pb2.Quaternion(x=0, y=0, z=0, w=1),
    )
    pose_proto = ros_proto_conversion.pose_to_proto(pose)
    self.assertEqual(pose_proto, pose_proto_expected)

  def test_pose_roundtrip(self):
    pose_proto = pose_pb2.Pose(
        position=point_pb2.Point(x=-1.32635246, y=-0.20890486, z=-0.16996824),
        orientation=quaternion_pb2.Quaternion(
            x=-0.5366317784871225,
            y=-0.4971772185983531,
            z=0.24928426929721537,
            w=-0.634585298210971,
        ),
    )
    result_proto = ros_proto_conversion.pose_to_proto(
        ros_proto_conversion.pose_from_proto(pose_proto)
    )
    # We expect bit-wise equality
    self.assertEqual(result_proto, pose_proto)


if __name__ == '__main__':
  absltest.main()
