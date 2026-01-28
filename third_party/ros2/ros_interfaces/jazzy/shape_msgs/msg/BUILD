load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/shape_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "mesh_proto",
    srcs = ["mesh.proto"],
    deps = [
        ":mesh_triangle_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:point_proto",
    ],
)

py_proto_library(
    name = "mesh_py_pb2",
    deps = ["mesh_proto"],
)

cc_proto_library(
    name = "mesh_cc_proto",
    deps = ["mesh_proto"],
)

go_proto_library(
    name = "mesh_go_proto",
    protos = ["mesh_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:point_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/shape_msgs/msg:mesh_triangle_go_proto",
    ],
)

proto_library(
    name = "mesh_triangle_proto",
    srcs = ["mesh_triangle.proto"],
)

py_proto_library(
    name = "mesh_triangle_py_pb2",
    deps = ["mesh_triangle_proto"],
)

cc_proto_library(
    name = "mesh_triangle_cc_proto",
    deps = ["mesh_triangle_proto"],
)

go_proto_library(
    name = "mesh_triangle_go_proto",
    protos = ["mesh_triangle_proto"],
)

proto_library(
    name = "plane_proto",
    srcs = ["plane.proto"],
)

py_proto_library(
    name = "plane_py_pb2",
    deps = ["plane_proto"],
)

cc_proto_library(
    name = "plane_cc_proto",
    deps = ["plane_proto"],
)

go_proto_library(
    name = "plane_go_proto",
    protos = ["plane_proto"],
)

proto_library(
    name = "solid_primitive_proto",
    srcs = ["solid_primitive.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:polygon_proto"],
)

py_proto_library(
    name = "solid_primitive_py_pb2",
    deps = ["solid_primitive_proto"],
)

cc_proto_library(
    name = "solid_primitive_cc_proto",
    deps = ["solid_primitive_proto"],
)

go_proto_library(
    name = "solid_primitive_go_proto",
    protos = ["solid_primitive_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:polygon_go_proto",
    ],
)
