load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/pcl_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "model_coefficients_proto",
    srcs = ["model_coefficients.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto"],
)

py_proto_library(
    name = "model_coefficients_py_pb2",
    deps = ["model_coefficients_proto"],
)

cc_proto_library(
    name = "model_coefficients_cc_proto",
    deps = ["model_coefficients_proto"],
)

go_proto_library(
    name = "model_coefficients_go_proto",
    protos = ["model_coefficients_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
)

proto_library(
    name = "point_indices_proto",
    srcs = ["point_indices.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto"],
)

py_proto_library(
    name = "point_indices_py_pb2",
    deps = ["point_indices_proto"],
)

cc_proto_library(
    name = "point_indices_cc_proto",
    deps = ["point_indices_proto"],
)

go_proto_library(
    name = "point_indices_go_proto",
    protos = ["point_indices_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
)

proto_library(
    name = "polygon_mesh_proto",
    srcs = ["polygon_mesh.proto"],
    deps = [
        ":vertices_proto",
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:point_cloud2_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "polygon_mesh_py_pb2",
    deps = ["polygon_mesh_proto"],
)

cc_proto_library(
    name = "polygon_mesh_cc_proto",
    deps = ["polygon_mesh_proto"],
)

go_proto_library(
    name = "polygon_mesh_go_proto",
    protos = ["polygon_mesh_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/pcl_msgs/msg:vertices_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:point_cloud2_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
)

proto_library(
    name = "vertices_proto",
    srcs = ["vertices.proto"],
)

py_proto_library(
    name = "vertices_py_pb2",
    deps = ["vertices_proto"],
)

cc_proto_library(
    name = "vertices_cc_proto",
    deps = ["vertices_proto"],
)

go_proto_library(
    name = "vertices_go_proto",
    protos = ["vertices_proto"],
)
