load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_proto_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "state_proto",
    srcs = ["state.proto"],
)

py_proto_library(
    name = "state_py_pb2",
    deps = ["state_proto"],
)

cc_proto_library(
    name = "state_cc_proto",
    deps = ["state_proto"],
)

go_proto_library(
    name = "state_go_proto",
    deps = ["state_proto"],
)

proto_library(
    name = "transition_proto",
    srcs = ["transition.proto"],
)

py_proto_library(
    name = "transition_py_pb2",
    deps = ["transition_proto"],
)

cc_proto_library(
    name = "transition_cc_proto",
    deps = ["transition_proto"],
)

go_proto_library(
    name = "transition_go_proto",
    deps = ["transition_proto"],
)

proto_library(
    name = "transition_description_proto",
    srcs = ["transition_description.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/msg:state_proto",
        "//third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/msg:transition_proto",
    ],
)

py_proto_library(
    name = "transition_description_py_pb2",
    deps = ["transition_description_proto"],
)

cc_proto_library(
    name = "transition_description_cc_proto",
    deps = ["transition_description_proto"],
)

go_proto_library(
    name = "transition_description_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/msg:state_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/msg:transition_go_proto",
    ],
    deps = ["transition_description_proto"],
)

proto_library(
    name = "transition_event_proto",
    srcs = ["transition_event.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/msg:state_proto",
        "//third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/msg:transition_proto",
    ],
)

py_proto_library(
    name = "transition_event_py_pb2",
    deps = ["transition_event_proto"],
)

cc_proto_library(
    name = "transition_event_cc_proto",
    deps = ["transition_event_proto"],
)

go_proto_library(
    name = "transition_event_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/msg:state_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/msg:transition_go_proto",
    ],
    deps = ["transition_event_proto"],
)
