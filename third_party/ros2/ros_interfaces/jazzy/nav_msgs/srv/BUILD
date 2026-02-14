load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/nav_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "get_map_proto",
    srcs = ["get_map.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/nav_msgs/msg:occupancy_grid_proto"],
)

py_proto_library(
    name = "get_map_py_pb2",
    deps = ["get_map_proto"],
)

cc_proto_library(
    name = "get_map_cc_proto",
    deps = ["get_map_proto"],
)

go_proto_library(
    name = "get_map_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/nav_msgs/srv/get_map_go_proto",
    protos = ["get_map_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/nav_msgs/msg:occupancy_grid_go_proto",
    ],
)

proto_library(
    name = "get_plan_proto",
    srcs = ["get_plan.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_stamped_proto",
        "//third_party/ros2/ros_interfaces/jazzy/nav_msgs/msg:path_proto",
    ],
)

py_proto_library(
    name = "get_plan_py_pb2",
    deps = ["get_plan_proto"],
)

cc_proto_library(
    name = "get_plan_cc_proto",
    deps = ["get_plan_proto"],
)

go_proto_library(
    name = "get_plan_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/nav_msgs/srv/get_plan_go_proto",
    protos = ["get_plan_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_stamped_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/nav_msgs/msg:path_go_proto",
    ],
)

proto_library(
    name = "load_map_proto",
    srcs = ["load_map.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/nav_msgs/msg:occupancy_grid_proto"],
)

py_proto_library(
    name = "load_map_py_pb2",
    deps = ["load_map_proto"],
)

cc_proto_library(
    name = "load_map_cc_proto",
    deps = ["load_map_proto"],
)

go_proto_library(
    name = "load_map_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/nav_msgs/srv/load_map_go_proto",
    protos = ["load_map_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/nav_msgs/msg:occupancy_grid_go_proto",
    ],
)

proto_library(
    name = "set_map_proto",
    srcs = ["set_map.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_with_covariance_stamped_proto",
        "//third_party/ros2/ros_interfaces/jazzy/nav_msgs/msg:occupancy_grid_proto",
    ],
)

py_proto_library(
    name = "set_map_py_pb2",
    deps = ["set_map_proto"],
)

cc_proto_library(
    name = "set_map_cc_proto",
    deps = ["set_map_proto"],
)

go_proto_library(
    name = "set_map_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/nav_msgs/srv/set_map_go_proto",
    protos = ["set_map_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_with_covariance_stamped_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/nav_msgs/msg:occupancy_grid_go_proto",
    ],
)
