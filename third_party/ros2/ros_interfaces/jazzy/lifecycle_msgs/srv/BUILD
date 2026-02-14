load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "change_state_proto",
    srcs = ["change_state.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/msg:transition_proto"],
)

py_proto_library(
    name = "change_state_py_pb2",
    deps = ["change_state_proto"],
)

cc_proto_library(
    name = "change_state_cc_proto",
    deps = ["change_state_proto"],
)

go_proto_library(
    name = "change_state_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/srv/change_state_go_proto",
    protos = ["change_state_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/msg:transition_go_proto",
    ],
)

proto_library(
    name = "get_available_states_proto",
    srcs = ["get_available_states.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/msg:state_proto"],
)

py_proto_library(
    name = "get_available_states_py_pb2",
    deps = ["get_available_states_proto"],
)

cc_proto_library(
    name = "get_available_states_cc_proto",
    deps = ["get_available_states_proto"],
)

go_proto_library(
    name = "get_available_states_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/srv/get_available_states_go_proto",
    protos = ["get_available_states_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/msg:state_go_proto",
    ],
)

proto_library(
    name = "get_available_transitions_proto",
    srcs = ["get_available_transitions.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/msg:transition_description_proto"],
)

py_proto_library(
    name = "get_available_transitions_py_pb2",
    deps = ["get_available_transitions_proto"],
)

cc_proto_library(
    name = "get_available_transitions_cc_proto",
    deps = ["get_available_transitions_proto"],
)

go_proto_library(
    name = "get_available_transitions_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/srv/get_available_transitions_go_proto",
    protos = ["get_available_transitions_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/msg:transition_description_go_proto",
    ],
)

proto_library(
    name = "get_state_proto",
    srcs = ["get_state.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/msg:state_proto"],
)

py_proto_library(
    name = "get_state_py_pb2",
    deps = ["get_state_proto"],
)

cc_proto_library(
    name = "get_state_cc_proto",
    deps = ["get_state_proto"],
)

go_proto_library(
    name = "get_state_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/srv/get_state_go_proto",
    protos = ["get_state_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/lifecycle_msgs/msg:state_go_proto",
    ],
)
