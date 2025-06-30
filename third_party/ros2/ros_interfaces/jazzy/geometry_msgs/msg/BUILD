load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/geometry_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "accel_proto",
    srcs = ["accel.proto"],
    deps = [":vector3_proto"],
)

py_proto_library(
    name = "accel_py_pb2",
    deps = ["accel_proto"],
)

cc_proto_library(
    name = "accel_cc_proto",
    deps = ["accel_proto"],
)

go_proto_library(
    name = "accel_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:vector3_go_proto",
    ],
    deps = ["accel_proto"],
)

proto_library(
    name = "accel_stamped_proto",
    srcs = ["accel_stamped.proto"],
    deps = [
        ":accel_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "accel_stamped_py_pb2",
    deps = ["accel_stamped_proto"],
)

cc_proto_library(
    name = "accel_stamped_cc_proto",
    deps = ["accel_stamped_proto"],
)

go_proto_library(
    name = "accel_stamped_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:accel_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["accel_stamped_proto"],
)

proto_library(
    name = "accel_with_covariance_proto",
    srcs = ["accel_with_covariance.proto"],
    deps = [":accel_proto"],
)

py_proto_library(
    name = "accel_with_covariance_py_pb2",
    deps = ["accel_with_covariance_proto"],
)

cc_proto_library(
    name = "accel_with_covariance_cc_proto",
    deps = ["accel_with_covariance_proto"],
)

go_proto_library(
    name = "accel_with_covariance_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:accel_go_proto",
    ],
    deps = ["accel_with_covariance_proto"],
)

proto_library(
    name = "accel_with_covariance_stamped_proto",
    srcs = ["accel_with_covariance_stamped.proto"],
    deps = [
        ":accel_with_covariance_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "accel_with_covariance_stamped_py_pb2",
    deps = ["accel_with_covariance_stamped_proto"],
)

cc_proto_library(
    name = "accel_with_covariance_stamped_cc_proto",
    deps = ["accel_with_covariance_stamped_proto"],
)

go_proto_library(
    name = "accel_with_covariance_stamped_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:accel_with_covariance_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["accel_with_covariance_stamped_proto"],
)

proto_library(
    name = "inertia_proto",
    srcs = ["inertia.proto"],
    deps = [":vector3_proto"],
)

py_proto_library(
    name = "inertia_py_pb2",
    deps = ["inertia_proto"],
)

cc_proto_library(
    name = "inertia_cc_proto",
    deps = ["inertia_proto"],
)

go_proto_library(
    name = "inertia_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:vector3_go_proto",
    ],
    deps = ["inertia_proto"],
)

proto_library(
    name = "inertia_stamped_proto",
    srcs = ["inertia_stamped.proto"],
    deps = [
        ":inertia_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "inertia_stamped_py_pb2",
    deps = ["inertia_stamped_proto"],
)

cc_proto_library(
    name = "inertia_stamped_cc_proto",
    deps = ["inertia_stamped_proto"],
)

go_proto_library(
    name = "inertia_stamped_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:inertia_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["inertia_stamped_proto"],
)

proto_library(
    name = "point_proto",
    srcs = ["point.proto"],
)

py_proto_library(
    name = "point_py_pb2",
    deps = ["point_proto"],
)

cc_proto_library(
    name = "point_cc_proto",
    deps = ["point_proto"],
)

go_proto_library(
    name = "point_go_proto",
    deps = ["point_proto"],
)

proto_library(
    name = "point32_proto",
    srcs = ["point32.proto"],
)

py_proto_library(
    name = "point32_py_pb2",
    deps = ["point32_proto"],
)

cc_proto_library(
    name = "point32_cc_proto",
    deps = ["point32_proto"],
)

go_proto_library(
    name = "point32_go_proto",
    deps = ["point32_proto"],
)

proto_library(
    name = "point_stamped_proto",
    srcs = ["point_stamped.proto"],
    deps = [
        ":point_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "point_stamped_py_pb2",
    deps = ["point_stamped_proto"],
)

cc_proto_library(
    name = "point_stamped_cc_proto",
    deps = ["point_stamped_proto"],
)

go_proto_library(
    name = "point_stamped_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:point_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["point_stamped_proto"],
)

proto_library(
    name = "polygon_proto",
    srcs = ["polygon.proto"],
    deps = [":point32_proto"],
)

py_proto_library(
    name = "polygon_py_pb2",
    deps = ["polygon_proto"],
)

cc_proto_library(
    name = "polygon_cc_proto",
    deps = ["polygon_proto"],
)

go_proto_library(
    name = "polygon_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:point32_go_proto",
    ],
    deps = ["polygon_proto"],
)

proto_library(
    name = "polygon_instance_proto",
    srcs = ["polygon_instance.proto"],
    deps = [":polygon_proto"],
)

py_proto_library(
    name = "polygon_instance_py_pb2",
    deps = ["polygon_instance_proto"],
)

cc_proto_library(
    name = "polygon_instance_cc_proto",
    deps = ["polygon_instance_proto"],
)

go_proto_library(
    name = "polygon_instance_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:polygon_go_proto",
    ],
    deps = ["polygon_instance_proto"],
)

proto_library(
    name = "polygon_instance_stamped_proto",
    srcs = ["polygon_instance_stamped.proto"],
    deps = [
        ":polygon_instance_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "polygon_instance_stamped_py_pb2",
    deps = ["polygon_instance_stamped_proto"],
)

cc_proto_library(
    name = "polygon_instance_stamped_cc_proto",
    deps = ["polygon_instance_stamped_proto"],
)

go_proto_library(
    name = "polygon_instance_stamped_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:polygon_instance_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["polygon_instance_stamped_proto"],
)

proto_library(
    name = "polygon_stamped_proto",
    srcs = ["polygon_stamped.proto"],
    deps = [
        ":polygon_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "polygon_stamped_py_pb2",
    deps = ["polygon_stamped_proto"],
)

cc_proto_library(
    name = "polygon_stamped_cc_proto",
    deps = ["polygon_stamped_proto"],
)

go_proto_library(
    name = "polygon_stamped_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:polygon_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["polygon_stamped_proto"],
)

proto_library(
    name = "pose_proto",
    srcs = ["pose.proto"],
    deps = [
        ":point_proto",
        ":quaternion_proto",
    ],
)

py_proto_library(
    name = "pose_py_pb2",
    deps = ["pose_proto"],
)

cc_proto_library(
    name = "pose_cc_proto",
    deps = ["pose_proto"],
)

go_proto_library(
    name = "pose_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:point_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:quaternion_go_proto",
    ],
    deps = ["pose_proto"],
)

proto_library(
    name = "pose2_d_proto",
    srcs = ["pose2_d.proto"],
)

py_proto_library(
    name = "pose2_d_py_pb2",
    deps = ["pose2_d_proto"],
)

cc_proto_library(
    name = "pose2_d_cc_proto",
    deps = ["pose2_d_proto"],
)

go_proto_library(
    name = "pose2_d_go_proto",
    deps = ["pose2_d_proto"],
)

proto_library(
    name = "pose_array_proto",
    srcs = ["pose_array.proto"],
    deps = [
        ":pose_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "pose_array_py_pb2",
    deps = ["pose_array_proto"],
)

cc_proto_library(
    name = "pose_array_cc_proto",
    deps = ["pose_array_proto"],
)

go_proto_library(
    name = "pose_array_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["pose_array_proto"],
)

proto_library(
    name = "pose_stamped_proto",
    srcs = ["pose_stamped.proto"],
    deps = [
        ":pose_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "pose_stamped_py_pb2",
    deps = ["pose_stamped_proto"],
)

cc_proto_library(
    name = "pose_stamped_cc_proto",
    deps = ["pose_stamped_proto"],
)

go_proto_library(
    name = "pose_stamped_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["pose_stamped_proto"],
)

proto_library(
    name = "pose_with_covariance_proto",
    srcs = ["pose_with_covariance.proto"],
    deps = [":pose_proto"],
)

py_proto_library(
    name = "pose_with_covariance_py_pb2",
    deps = ["pose_with_covariance_proto"],
)

cc_proto_library(
    name = "pose_with_covariance_cc_proto",
    deps = ["pose_with_covariance_proto"],
)

go_proto_library(
    name = "pose_with_covariance_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_go_proto",
    ],
    deps = ["pose_with_covariance_proto"],
)

proto_library(
    name = "pose_with_covariance_stamped_proto",
    srcs = ["pose_with_covariance_stamped.proto"],
    deps = [
        ":pose_with_covariance_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "pose_with_covariance_stamped_py_pb2",
    deps = ["pose_with_covariance_stamped_proto"],
)

cc_proto_library(
    name = "pose_with_covariance_stamped_cc_proto",
    deps = ["pose_with_covariance_stamped_proto"],
)

go_proto_library(
    name = "pose_with_covariance_stamped_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_with_covariance_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["pose_with_covariance_stamped_proto"],
)

proto_library(
    name = "quaternion_proto",
    srcs = ["quaternion.proto"],
)

py_proto_library(
    name = "quaternion_py_pb2",
    deps = ["quaternion_proto"],
)

cc_proto_library(
    name = "quaternion_cc_proto",
    deps = ["quaternion_proto"],
)

go_proto_library(
    name = "quaternion_go_proto",
    deps = ["quaternion_proto"],
)

proto_library(
    name = "quaternion_stamped_proto",
    srcs = ["quaternion_stamped.proto"],
    deps = [
        ":quaternion_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "quaternion_stamped_py_pb2",
    deps = ["quaternion_stamped_proto"],
)

cc_proto_library(
    name = "quaternion_stamped_cc_proto",
    deps = ["quaternion_stamped_proto"],
)

go_proto_library(
    name = "quaternion_stamped_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:quaternion_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["quaternion_stamped_proto"],
)

proto_library(
    name = "transform_proto",
    srcs = ["transform.proto"],
    deps = [
        ":quaternion_proto",
        ":vector3_proto",
    ],
)

py_proto_library(
    name = "transform_py_pb2",
    deps = ["transform_proto"],
)

cc_proto_library(
    name = "transform_cc_proto",
    deps = ["transform_proto"],
)

go_proto_library(
    name = "transform_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:quaternion_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:vector3_go_proto",
    ],
    deps = ["transform_proto"],
)

proto_library(
    name = "transform_stamped_proto",
    srcs = ["transform_stamped.proto"],
    deps = [
        ":transform_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "transform_stamped_py_pb2",
    deps = ["transform_stamped_proto"],
)

cc_proto_library(
    name = "transform_stamped_cc_proto",
    deps = ["transform_stamped_proto"],
)

go_proto_library(
    name = "transform_stamped_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:transform_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["transform_stamped_proto"],
)

proto_library(
    name = "twist_proto",
    srcs = ["twist.proto"],
    deps = [":vector3_proto"],
)

py_proto_library(
    name = "twist_py_pb2",
    deps = ["twist_proto"],
)

cc_proto_library(
    name = "twist_cc_proto",
    deps = ["twist_proto"],
)

go_proto_library(
    name = "twist_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:vector3_go_proto",
    ],
    deps = ["twist_proto"],
)

proto_library(
    name = "twist_stamped_proto",
    srcs = ["twist_stamped.proto"],
    deps = [
        ":twist_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "twist_stamped_py_pb2",
    deps = ["twist_stamped_proto"],
)

cc_proto_library(
    name = "twist_stamped_cc_proto",
    deps = ["twist_stamped_proto"],
)

go_proto_library(
    name = "twist_stamped_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:twist_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["twist_stamped_proto"],
)

proto_library(
    name = "twist_with_covariance_proto",
    srcs = ["twist_with_covariance.proto"],
    deps = [":twist_proto"],
)

py_proto_library(
    name = "twist_with_covariance_py_pb2",
    deps = ["twist_with_covariance_proto"],
)

cc_proto_library(
    name = "twist_with_covariance_cc_proto",
    deps = ["twist_with_covariance_proto"],
)

go_proto_library(
    name = "twist_with_covariance_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:twist_go_proto",
    ],
    deps = ["twist_with_covariance_proto"],
)

proto_library(
    name = "twist_with_covariance_stamped_proto",
    srcs = ["twist_with_covariance_stamped.proto"],
    deps = [
        ":twist_with_covariance_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "twist_with_covariance_stamped_py_pb2",
    deps = ["twist_with_covariance_stamped_proto"],
)

cc_proto_library(
    name = "twist_with_covariance_stamped_cc_proto",
    deps = ["twist_with_covariance_stamped_proto"],
)

go_proto_library(
    name = "twist_with_covariance_stamped_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:twist_with_covariance_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["twist_with_covariance_stamped_proto"],
)

proto_library(
    name = "vector3_proto",
    srcs = ["vector3.proto"],
)

py_proto_library(
    name = "vector3_py_pb2",
    deps = ["vector3_proto"],
)

cc_proto_library(
    name = "vector3_cc_proto",
    deps = ["vector3_proto"],
)

go_proto_library(
    name = "vector3_go_proto",
    deps = ["vector3_proto"],
)

proto_library(
    name = "vector3_stamped_proto",
    srcs = ["vector3_stamped.proto"],
    deps = [
        ":vector3_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "vector3_stamped_py_pb2",
    deps = ["vector3_stamped_proto"],
)

cc_proto_library(
    name = "vector3_stamped_cc_proto",
    deps = ["vector3_stamped_proto"],
)

go_proto_library(
    name = "vector3_stamped_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:vector3_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["vector3_stamped_proto"],
)

proto_library(
    name = "velocity_stamped_proto",
    srcs = ["velocity_stamped.proto"],
    deps = [
        ":twist_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "velocity_stamped_py_pb2",
    deps = ["velocity_stamped_proto"],
)

cc_proto_library(
    name = "velocity_stamped_cc_proto",
    deps = ["velocity_stamped_proto"],
)

go_proto_library(
    name = "velocity_stamped_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:twist_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["velocity_stamped_proto"],
)

proto_library(
    name = "wrench_proto",
    srcs = ["wrench.proto"],
    deps = [":vector3_proto"],
)

py_proto_library(
    name = "wrench_py_pb2",
    deps = ["wrench_proto"],
)

cc_proto_library(
    name = "wrench_cc_proto",
    deps = ["wrench_proto"],
)

go_proto_library(
    name = "wrench_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:vector3_go_proto",
    ],
    deps = ["wrench_proto"],
)

proto_library(
    name = "wrench_stamped_proto",
    srcs = ["wrench_stamped.proto"],
    deps = [
        ":wrench_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "wrench_stamped_py_pb2",
    deps = ["wrench_stamped_proto"],
)

cc_proto_library(
    name = "wrench_stamped_cc_proto",
    deps = ["wrench_stamped_proto"],
)

go_proto_library(
    name = "wrench_stamped_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:wrench_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["wrench_stamped_proto"],
)
