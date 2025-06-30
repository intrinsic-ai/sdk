load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/visualization_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "image_marker_proto",
    srcs = ["image_marker.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:point_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:color_rgba_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "image_marker_py_pb2",
    deps = ["image_marker_proto"],
)

cc_proto_library(
    name = "image_marker_cc_proto",
    deps = ["image_marker_proto"],
)

go_proto_library(
    name = "image_marker_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:point_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:color_rgba_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["image_marker_proto"],
)

proto_library(
    name = "interactive_marker_proto",
    srcs = ["interactive_marker.proto"],
    deps = [
        ":interactive_marker_control_proto",
        ":menu_entry_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "interactive_marker_py_pb2",
    deps = ["interactive_marker_proto"],
)

cc_proto_library(
    name = "interactive_marker_cc_proto",
    deps = ["interactive_marker_proto"],
)

go_proto_library(
    name = "interactive_marker_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/visualization_msgs/msg:interactive_marker_control_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/visualization_msgs/msg:menu_entry_go_proto",
    ],
    deps = ["interactive_marker_proto"],
)

proto_library(
    name = "interactive_marker_control_proto",
    srcs = ["interactive_marker_control.proto"],
    deps = [
        ":marker_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:quaternion_proto",
    ],
)

py_proto_library(
    name = "interactive_marker_control_py_pb2",
    deps = ["interactive_marker_control_proto"],
)

cc_proto_library(
    name = "interactive_marker_control_cc_proto",
    deps = ["interactive_marker_control_proto"],
)

go_proto_library(
    name = "interactive_marker_control_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:quaternion_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/visualization_msgs/msg:marker_go_proto",
    ],
    deps = ["interactive_marker_control_proto"],
)

proto_library(
    name = "interactive_marker_feedback_proto",
    srcs = ["interactive_marker_feedback.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:point_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "interactive_marker_feedback_py_pb2",
    deps = ["interactive_marker_feedback_proto"],
)

cc_proto_library(
    name = "interactive_marker_feedback_cc_proto",
    deps = ["interactive_marker_feedback_proto"],
)

go_proto_library(
    name = "interactive_marker_feedback_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:point_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["interactive_marker_feedback_proto"],
)

proto_library(
    name = "interactive_marker_init_proto",
    srcs = ["interactive_marker_init.proto"],
    deps = [":interactive_marker_proto"],
)

py_proto_library(
    name = "interactive_marker_init_py_pb2",
    deps = ["interactive_marker_init_proto"],
)

cc_proto_library(
    name = "interactive_marker_init_cc_proto",
    deps = ["interactive_marker_init_proto"],
)

go_proto_library(
    name = "interactive_marker_init_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/visualization_msgs/msg:interactive_marker_go_proto",
    ],
    deps = ["interactive_marker_init_proto"],
)

proto_library(
    name = "interactive_marker_pose_proto",
    srcs = ["interactive_marker_pose.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "interactive_marker_pose_py_pb2",
    deps = ["interactive_marker_pose_proto"],
)

cc_proto_library(
    name = "interactive_marker_pose_cc_proto",
    deps = ["interactive_marker_pose_proto"],
)

go_proto_library(
    name = "interactive_marker_pose_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["interactive_marker_pose_proto"],
)

proto_library(
    name = "interactive_marker_update_proto",
    srcs = ["interactive_marker_update.proto"],
    deps = [
        ":interactive_marker_pose_proto",
        ":interactive_marker_proto",
    ],
)

py_proto_library(
    name = "interactive_marker_update_py_pb2",
    deps = ["interactive_marker_update_proto"],
)

cc_proto_library(
    name = "interactive_marker_update_cc_proto",
    deps = ["interactive_marker_update_proto"],
)

go_proto_library(
    name = "interactive_marker_update_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/visualization_msgs/msg:interactive_marker_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/visualization_msgs/msg:interactive_marker_pose_go_proto",
    ],
    deps = ["interactive_marker_update_proto"],
)

proto_library(
    name = "marker_proto",
    srcs = ["marker.proto"],
    deps = [
        ":mesh_file_proto",
        ":uv_coordinate_proto",
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:point_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:vector3_proto",
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:compressed_image_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:color_rgba_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "marker_py_pb2",
    deps = ["marker_proto"],
)

cc_proto_library(
    name = "marker_cc_proto",
    deps = ["marker_proto"],
)

go_proto_library(
    name = "marker_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:point_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:vector3_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:compressed_image_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:color_rgba_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/visualization_msgs/msg:mesh_file_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/visualization_msgs/msg:uv_coordinate_go_proto",
    ],
    deps = ["marker_proto"],
)

proto_library(
    name = "marker_array_proto",
    srcs = ["marker_array.proto"],
    deps = [":marker_proto"],
)

py_proto_library(
    name = "marker_array_py_pb2",
    deps = ["marker_array_proto"],
)

cc_proto_library(
    name = "marker_array_cc_proto",
    deps = ["marker_array_proto"],
)

go_proto_library(
    name = "marker_array_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/visualization_msgs/msg:marker_go_proto",
    ],
    deps = ["marker_array_proto"],
)

proto_library(
    name = "menu_entry_proto",
    srcs = ["menu_entry.proto"],
)

py_proto_library(
    name = "menu_entry_py_pb2",
    deps = ["menu_entry_proto"],
)

cc_proto_library(
    name = "menu_entry_cc_proto",
    deps = ["menu_entry_proto"],
)

go_proto_library(
    name = "menu_entry_go_proto",
    deps = ["menu_entry_proto"],
)

proto_library(
    name = "mesh_file_proto",
    srcs = ["mesh_file.proto"],
)

py_proto_library(
    name = "mesh_file_py_pb2",
    deps = ["mesh_file_proto"],
)

cc_proto_library(
    name = "mesh_file_cc_proto",
    deps = ["mesh_file_proto"],
)

go_proto_library(
    name = "mesh_file_go_proto",
    deps = ["mesh_file_proto"],
)

proto_library(
    name = "uv_coordinate_proto",
    srcs = ["uv_coordinate.proto"],
)

py_proto_library(
    name = "uv_coordinate_py_pb2",
    deps = ["uv_coordinate_proto"],
)

cc_proto_library(
    name = "uv_coordinate_cc_proto",
    deps = ["uv_coordinate_proto"],
)

go_proto_library(
    name = "uv_coordinate_go_proto",
    deps = ["uv_coordinate_proto"],
)
