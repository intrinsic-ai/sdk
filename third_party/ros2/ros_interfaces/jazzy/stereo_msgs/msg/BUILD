load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/stereo_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "disparity_image_proto",
    srcs = ["disparity_image.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:image_proto",
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:region_of_interest_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "disparity_image_py_pb2",
    deps = ["disparity_image_proto"],
)

cc_proto_library(
    name = "disparity_image_cc_proto",
    deps = ["disparity_image_proto"],
)

go_proto_library(
    name = "disparity_image_go_proto",
    protos = ["disparity_image_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:image_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:region_of_interest_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
)
