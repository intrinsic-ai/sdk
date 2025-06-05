load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_proto_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/nav_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "goals_proto",
    srcs = ["goals.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_stamped_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "goals_py_pb2",
    deps = ["goals_proto"],
)

cc_proto_library(
    name = "goals_cc_proto",
    deps = ["goals_proto"],
)

go_proto_library(
    name = "goals_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_stamped_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["goals_proto"],
)

proto_library(
    name = "grid_cells_proto",
    srcs = ["grid_cells.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:point_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "grid_cells_py_pb2",
    deps = ["grid_cells_proto"],
)

cc_proto_library(
    name = "grid_cells_cc_proto",
    deps = ["grid_cells_proto"],
)

go_proto_library(
    name = "grid_cells_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:point_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["grid_cells_proto"],
)

proto_library(
    name = "map_meta_data_proto",
    srcs = ["map_meta_data.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_proto",
    ],
)

py_proto_library(
    name = "map_meta_data_py_pb2",
    deps = ["map_meta_data_proto"],
)

cc_proto_library(
    name = "map_meta_data_cc_proto",
    deps = ["map_meta_data_proto"],
)

go_proto_library(
    name = "map_meta_data_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_go_proto",
    ],
    deps = ["map_meta_data_proto"],
)

proto_library(
    name = "occupancy_grid_proto",
    srcs = ["occupancy_grid.proto"],
    deps = [
        ":map_meta_data_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "occupancy_grid_py_pb2",
    deps = ["occupancy_grid_proto"],
)

cc_proto_library(
    name = "occupancy_grid_cc_proto",
    deps = ["occupancy_grid_proto"],
)

go_proto_library(
    name = "occupancy_grid_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/nav_msgs/msg:map_meta_data_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["occupancy_grid_proto"],
)

proto_library(
    name = "odometry_proto",
    srcs = ["odometry.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_with_covariance_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:twist_with_covariance_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "odometry_py_pb2",
    deps = ["odometry_proto"],
)

cc_proto_library(
    name = "odometry_cc_proto",
    deps = ["odometry_proto"],
)

go_proto_library(
    name = "odometry_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_with_covariance_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:twist_with_covariance_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["odometry_proto"],
)

proto_library(
    name = "path_proto",
    srcs = ["path.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_stamped_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "path_py_pb2",
    deps = ["path_proto"],
)

cc_proto_library(
    name = "path_cc_proto",
    deps = ["path_proto"],
)

go_proto_library(
    name = "path_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_stamped_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["path_proto"],
)
