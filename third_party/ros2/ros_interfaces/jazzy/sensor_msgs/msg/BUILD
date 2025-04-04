load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_proto_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/sensor_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "battery_state_proto",
    srcs = ["battery_state.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto"],
)

py_proto_library(
    name = "battery_state_py_pb2",
    deps = ["battery_state_proto"],
)

cc_proto_library(
    name = "battery_state_cc_proto",
    deps = ["battery_state_proto"],
)

go_proto_library(
    name = "battery_state_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["battery_state_proto"],
)

proto_library(
    name = "camera_info_proto",
    srcs = ["camera_info.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:region_of_interest_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "camera_info_py_pb2",
    deps = ["camera_info_proto"],
)

cc_proto_library(
    name = "camera_info_cc_proto",
    deps = ["camera_info_proto"],
)

go_proto_library(
    name = "camera_info_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:region_of_interest_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["camera_info_proto"],
)

proto_library(
    name = "channel_float32_proto",
    srcs = ["channel_float32.proto"],
)

py_proto_library(
    name = "channel_float32_py_pb2",
    deps = ["channel_float32_proto"],
)

cc_proto_library(
    name = "channel_float32_cc_proto",
    deps = ["channel_float32_proto"],
)

go_proto_library(
    name = "channel_float32_go_proto",
    deps = ["channel_float32_proto"],
)

proto_library(
    name = "compressed_image_proto",
    srcs = ["compressed_image.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto"],
)

py_proto_library(
    name = "compressed_image_py_pb2",
    deps = ["compressed_image_proto"],
)

cc_proto_library(
    name = "compressed_image_cc_proto",
    deps = ["compressed_image_proto"],
)

go_proto_library(
    name = "compressed_image_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["compressed_image_proto"],
)

proto_library(
    name = "fluid_pressure_proto",
    srcs = ["fluid_pressure.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto"],
)

py_proto_library(
    name = "fluid_pressure_py_pb2",
    deps = ["fluid_pressure_proto"],
)

cc_proto_library(
    name = "fluid_pressure_cc_proto",
    deps = ["fluid_pressure_proto"],
)

go_proto_library(
    name = "fluid_pressure_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["fluid_pressure_proto"],
)

proto_library(
    name = "illuminance_proto",
    srcs = ["illuminance.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto"],
)

py_proto_library(
    name = "illuminance_py_pb2",
    deps = ["illuminance_proto"],
)

cc_proto_library(
    name = "illuminance_cc_proto",
    deps = ["illuminance_proto"],
)

go_proto_library(
    name = "illuminance_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["illuminance_proto"],
)

proto_library(
    name = "image_proto",
    srcs = ["image.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto"],
)

py_proto_library(
    name = "image_py_pb2",
    deps = ["image_proto"],
)

cc_proto_library(
    name = "image_cc_proto",
    deps = ["image_proto"],
)

go_proto_library(
    name = "image_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["image_proto"],
)

proto_library(
    name = "imu_proto",
    srcs = ["imu.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:quaternion_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:vector3_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "imu_py_pb2",
    deps = ["imu_proto"],
)

cc_proto_library(
    name = "imu_cc_proto",
    deps = ["imu_proto"],
)

go_proto_library(
    name = "imu_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:quaternion_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:vector3_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["imu_proto"],
)

proto_library(
    name = "joint_state_proto",
    srcs = ["joint_state.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto"],
)

py_proto_library(
    name = "joint_state_py_pb2",
    deps = ["joint_state_proto"],
)

cc_proto_library(
    name = "joint_state_cc_proto",
    deps = ["joint_state_proto"],
)

go_proto_library(
    name = "joint_state_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["joint_state_proto"],
)

proto_library(
    name = "joy_proto",
    srcs = ["joy.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto"],
)

py_proto_library(
    name = "joy_py_pb2",
    deps = ["joy_proto"],
)

cc_proto_library(
    name = "joy_cc_proto",
    deps = ["joy_proto"],
)

go_proto_library(
    name = "joy_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["joy_proto"],
)

proto_library(
    name = "joy_feedback_proto",
    srcs = ["joy_feedback.proto"],
)

py_proto_library(
    name = "joy_feedback_py_pb2",
    deps = ["joy_feedback_proto"],
)

cc_proto_library(
    name = "joy_feedback_cc_proto",
    deps = ["joy_feedback_proto"],
)

go_proto_library(
    name = "joy_feedback_go_proto",
    deps = ["joy_feedback_proto"],
)

proto_library(
    name = "joy_feedback_array_proto",
    srcs = ["joy_feedback_array.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:joy_feedback_proto"],
)

py_proto_library(
    name = "joy_feedback_array_py_pb2",
    deps = ["joy_feedback_array_proto"],
)

cc_proto_library(
    name = "joy_feedback_array_cc_proto",
    deps = ["joy_feedback_array_proto"],
)

go_proto_library(
    name = "joy_feedback_array_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:joy_feedback_go_proto",
    ],
    deps = ["joy_feedback_array_proto"],
)

proto_library(
    name = "laser_echo_proto",
    srcs = ["laser_echo.proto"],
)

py_proto_library(
    name = "laser_echo_py_pb2",
    deps = ["laser_echo_proto"],
)

cc_proto_library(
    name = "laser_echo_cc_proto",
    deps = ["laser_echo_proto"],
)

go_proto_library(
    name = "laser_echo_go_proto",
    deps = ["laser_echo_proto"],
)

proto_library(
    name = "laser_scan_proto",
    srcs = ["laser_scan.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto"],
)

py_proto_library(
    name = "laser_scan_py_pb2",
    deps = ["laser_scan_proto"],
)

cc_proto_library(
    name = "laser_scan_cc_proto",
    deps = ["laser_scan_proto"],
)

go_proto_library(
    name = "laser_scan_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["laser_scan_proto"],
)

proto_library(
    name = "magnetic_field_proto",
    srcs = ["magnetic_field.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:vector3_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "magnetic_field_py_pb2",
    deps = ["magnetic_field_proto"],
)

cc_proto_library(
    name = "magnetic_field_cc_proto",
    deps = ["magnetic_field_proto"],
)

go_proto_library(
    name = "magnetic_field_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:vector3_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["magnetic_field_proto"],
)

proto_library(
    name = "multi_dof_joint_state_proto",
    srcs = ["multi_dof_joint_state.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:transform_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:twist_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:wrench_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "multi_dof_joint_state_py_pb2",
    deps = ["multi_dof_joint_state_proto"],
)

cc_proto_library(
    name = "multi_dof_joint_state_cc_proto",
    deps = ["multi_dof_joint_state_proto"],
)

go_proto_library(
    name = "multi_dof_joint_state_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:transform_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:twist_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:wrench_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["multi_dof_joint_state_proto"],
)

proto_library(
    name = "multi_echo_laser_scan_proto",
    srcs = ["multi_echo_laser_scan.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:laser_echo_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "multi_echo_laser_scan_py_pb2",
    deps = ["multi_echo_laser_scan_proto"],
)

cc_proto_library(
    name = "multi_echo_laser_scan_cc_proto",
    deps = ["multi_echo_laser_scan_proto"],
)

go_proto_library(
    name = "multi_echo_laser_scan_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:laser_echo_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["multi_echo_laser_scan_proto"],
)

proto_library(
    name = "nav_sat_fix_proto",
    srcs = ["nav_sat_fix.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:nav_sat_status_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "nav_sat_fix_py_pb2",
    deps = ["nav_sat_fix_proto"],
)

cc_proto_library(
    name = "nav_sat_fix_cc_proto",
    deps = ["nav_sat_fix_proto"],
)

go_proto_library(
    name = "nav_sat_fix_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:nav_sat_status_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["nav_sat_fix_proto"],
)

proto_library(
    name = "nav_sat_status_proto",
    srcs = ["nav_sat_status.proto"],
)

py_proto_library(
    name = "nav_sat_status_py_pb2",
    deps = ["nav_sat_status_proto"],
)

cc_proto_library(
    name = "nav_sat_status_cc_proto",
    deps = ["nav_sat_status_proto"],
)

go_proto_library(
    name = "nav_sat_status_go_proto",
    deps = ["nav_sat_status_proto"],
)

proto_library(
    name = "point_cloud_proto",
    srcs = ["point_cloud.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:point32_proto",
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:channel_float32_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "point_cloud_py_pb2",
    deps = ["point_cloud_proto"],
)

cc_proto_library(
    name = "point_cloud_cc_proto",
    deps = ["point_cloud_proto"],
)

go_proto_library(
    name = "point_cloud_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:point32_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:channel_float32_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["point_cloud_proto"],
)

proto_library(
    name = "point_cloud2_proto",
    srcs = ["point_cloud2.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:point_field_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "point_cloud2_py_pb2",
    deps = ["point_cloud2_proto"],
)

cc_proto_library(
    name = "point_cloud2_cc_proto",
    deps = ["point_cloud2_proto"],
)

go_proto_library(
    name = "point_cloud2_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:point_field_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["point_cloud2_proto"],
)

proto_library(
    name = "point_field_proto",
    srcs = ["point_field.proto"],
)

py_proto_library(
    name = "point_field_py_pb2",
    deps = ["point_field_proto"],
)

cc_proto_library(
    name = "point_field_cc_proto",
    deps = ["point_field_proto"],
)

go_proto_library(
    name = "point_field_go_proto",
    deps = ["point_field_proto"],
)

proto_library(
    name = "range_proto",
    srcs = ["range.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto"],
)

py_proto_library(
    name = "range_py_pb2",
    deps = ["range_proto"],
)

cc_proto_library(
    name = "range_cc_proto",
    deps = ["range_proto"],
)

go_proto_library(
    name = "range_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["range_proto"],
)

proto_library(
    name = "region_of_interest_proto",
    srcs = ["region_of_interest.proto"],
)

py_proto_library(
    name = "region_of_interest_py_pb2",
    deps = ["region_of_interest_proto"],
)

cc_proto_library(
    name = "region_of_interest_cc_proto",
    deps = ["region_of_interest_proto"],
)

go_proto_library(
    name = "region_of_interest_go_proto",
    deps = ["region_of_interest_proto"],
)

proto_library(
    name = "relative_humidity_proto",
    srcs = ["relative_humidity.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto"],
)

py_proto_library(
    name = "relative_humidity_py_pb2",
    deps = ["relative_humidity_proto"],
)

cc_proto_library(
    name = "relative_humidity_cc_proto",
    deps = ["relative_humidity_proto"],
)

go_proto_library(
    name = "relative_humidity_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["relative_humidity_proto"],
)

proto_library(
    name = "temperature_proto",
    srcs = ["temperature.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto"],
)

py_proto_library(
    name = "temperature_py_pb2",
    deps = ["temperature_proto"],
)

cc_proto_library(
    name = "temperature_cc_proto",
    deps = ["temperature_proto"],
)

go_proto_library(
    name = "temperature_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["temperature_proto"],
)

proto_library(
    name = "time_reference_proto",
    srcs = ["time_reference.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "time_reference_py_pb2",
    deps = ["time_reference_proto"],
)

cc_proto_library(
    name = "time_reference_cc_proto",
    deps = ["time_reference_proto"],
)

go_proto_library(
    name = "time_reference_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["time_reference_proto"],
)
