load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/std_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "bool_proto",
    srcs = ["bool.proto"],
)

py_proto_library(
    name = "bool_py_pb2",
    deps = ["bool_proto"],
)

cc_proto_library(
    name = "bool_cc_proto",
    deps = ["bool_proto"],
)

go_proto_library(
    name = "bool_go_proto",
    deps = ["bool_proto"],
)

proto_library(
    name = "byte_proto",
    srcs = ["byte.proto"],
)

py_proto_library(
    name = "byte_py_pb2",
    deps = ["byte_proto"],
)

cc_proto_library(
    name = "byte_cc_proto",
    deps = ["byte_proto"],
)

go_proto_library(
    name = "byte_go_proto",
    deps = ["byte_proto"],
)

proto_library(
    name = "byte_multi_array_proto",
    srcs = ["byte_multi_array.proto"],
    deps = [":multi_array_layout_proto"],
)

py_proto_library(
    name = "byte_multi_array_py_pb2",
    deps = ["byte_multi_array_proto"],
)

cc_proto_library(
    name = "byte_multi_array_cc_proto",
    deps = ["byte_multi_array_proto"],
)

go_proto_library(
    name = "byte_multi_array_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:multi_array_layout_go_proto",
    ],
    deps = ["byte_multi_array_proto"],
)

proto_library(
    name = "char_proto",
    srcs = ["char.proto"],
)

py_proto_library(
    name = "char_py_pb2",
    deps = ["char_proto"],
)

cc_proto_library(
    name = "char_cc_proto",
    deps = ["char_proto"],
)

go_proto_library(
    name = "char_go_proto",
    deps = ["char_proto"],
)

proto_library(
    name = "color_rgba_proto",
    srcs = ["color_rgba.proto"],
)

py_proto_library(
    name = "color_rgba_py_pb2",
    deps = ["color_rgba_proto"],
)

cc_proto_library(
    name = "color_rgba_cc_proto",
    deps = ["color_rgba_proto"],
)

go_proto_library(
    name = "color_rgba_go_proto",
    deps = ["color_rgba_proto"],
)

proto_library(
    name = "empty_proto",
    srcs = ["empty.proto"],
)

py_proto_library(
    name = "empty_py_pb2",
    deps = ["empty_proto"],
)

cc_proto_library(
    name = "empty_cc_proto",
    deps = ["empty_proto"],
)

go_proto_library(
    name = "empty_go_proto",
    deps = ["empty_proto"],
)

proto_library(
    name = "float32_proto",
    srcs = ["float32.proto"],
)

py_proto_library(
    name = "float32_py_pb2",
    deps = ["float32_proto"],
)

cc_proto_library(
    name = "float32_cc_proto",
    deps = ["float32_proto"],
)

go_proto_library(
    name = "float32_go_proto",
    deps = ["float32_proto"],
)

proto_library(
    name = "float32_multi_array_proto",
    srcs = ["float32_multi_array.proto"],
    deps = [":multi_array_layout_proto"],
)

py_proto_library(
    name = "float32_multi_array_py_pb2",
    deps = ["float32_multi_array_proto"],
)

cc_proto_library(
    name = "float32_multi_array_cc_proto",
    deps = ["float32_multi_array_proto"],
)

go_proto_library(
    name = "float32_multi_array_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:multi_array_layout_go_proto",
    ],
    deps = ["float32_multi_array_proto"],
)

proto_library(
    name = "float64_proto",
    srcs = ["float64.proto"],
)

py_proto_library(
    name = "float64_py_pb2",
    deps = ["float64_proto"],
)

cc_proto_library(
    name = "float64_cc_proto",
    deps = ["float64_proto"],
)

go_proto_library(
    name = "float64_go_proto",
    deps = ["float64_proto"],
)

proto_library(
    name = "float64_multi_array_proto",
    srcs = ["float64_multi_array.proto"],
    deps = [":multi_array_layout_proto"],
)

py_proto_library(
    name = "float64_multi_array_py_pb2",
    deps = ["float64_multi_array_proto"],
)

cc_proto_library(
    name = "float64_multi_array_cc_proto",
    deps = ["float64_multi_array_proto"],
)

go_proto_library(
    name = "float64_multi_array_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:multi_array_layout_go_proto",
    ],
    deps = ["float64_multi_array_proto"],
)

proto_library(
    name = "header_proto",
    srcs = ["header.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_proto"],
)

py_proto_library(
    name = "header_py_pb2",
    deps = ["header_proto"],
)

cc_proto_library(
    name = "header_cc_proto",
    deps = ["header_proto"],
)

go_proto_library(
    name = "header_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_go_proto",
    ],
    deps = ["header_proto"],
)

proto_library(
    name = "int16_proto",
    srcs = ["int16.proto"],
)

py_proto_library(
    name = "int16_py_pb2",
    deps = ["int16_proto"],
)

cc_proto_library(
    name = "int16_cc_proto",
    deps = ["int16_proto"],
)

go_proto_library(
    name = "int16_go_proto",
    deps = ["int16_proto"],
)

proto_library(
    name = "int16_multi_array_proto",
    srcs = ["int16_multi_array.proto"],
    deps = [":multi_array_layout_proto"],
)

py_proto_library(
    name = "int16_multi_array_py_pb2",
    deps = ["int16_multi_array_proto"],
)

cc_proto_library(
    name = "int16_multi_array_cc_proto",
    deps = ["int16_multi_array_proto"],
)

go_proto_library(
    name = "int16_multi_array_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:multi_array_layout_go_proto",
    ],
    deps = ["int16_multi_array_proto"],
)

proto_library(
    name = "int32_proto",
    srcs = ["int32.proto"],
)

py_proto_library(
    name = "int32_py_pb2",
    deps = ["int32_proto"],
)

cc_proto_library(
    name = "int32_cc_proto",
    deps = ["int32_proto"],
)

go_proto_library(
    name = "int32_go_proto",
    deps = ["int32_proto"],
)

proto_library(
    name = "int32_multi_array_proto",
    srcs = ["int32_multi_array.proto"],
    deps = [":multi_array_layout_proto"],
)

py_proto_library(
    name = "int32_multi_array_py_pb2",
    deps = ["int32_multi_array_proto"],
)

cc_proto_library(
    name = "int32_multi_array_cc_proto",
    deps = ["int32_multi_array_proto"],
)

go_proto_library(
    name = "int32_multi_array_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:multi_array_layout_go_proto",
    ],
    deps = ["int32_multi_array_proto"],
)

proto_library(
    name = "int64_proto",
    srcs = ["int64.proto"],
)

py_proto_library(
    name = "int64_py_pb2",
    deps = ["int64_proto"],
)

cc_proto_library(
    name = "int64_cc_proto",
    deps = ["int64_proto"],
)

go_proto_library(
    name = "int64_go_proto",
    deps = ["int64_proto"],
)

proto_library(
    name = "int64_multi_array_proto",
    srcs = ["int64_multi_array.proto"],
    deps = [":multi_array_layout_proto"],
)

py_proto_library(
    name = "int64_multi_array_py_pb2",
    deps = ["int64_multi_array_proto"],
)

cc_proto_library(
    name = "int64_multi_array_cc_proto",
    deps = ["int64_multi_array_proto"],
)

go_proto_library(
    name = "int64_multi_array_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:multi_array_layout_go_proto",
    ],
    deps = ["int64_multi_array_proto"],
)

proto_library(
    name = "int8_proto",
    srcs = ["int8.proto"],
)

py_proto_library(
    name = "int8_py_pb2",
    deps = ["int8_proto"],
)

cc_proto_library(
    name = "int8_cc_proto",
    deps = ["int8_proto"],
)

go_proto_library(
    name = "int8_go_proto",
    deps = ["int8_proto"],
)

proto_library(
    name = "int8_multi_array_proto",
    srcs = ["int8_multi_array.proto"],
    deps = [":multi_array_layout_proto"],
)

py_proto_library(
    name = "int8_multi_array_py_pb2",
    deps = ["int8_multi_array_proto"],
)

cc_proto_library(
    name = "int8_multi_array_cc_proto",
    deps = ["int8_multi_array_proto"],
)

go_proto_library(
    name = "int8_multi_array_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:multi_array_layout_go_proto",
    ],
    deps = ["int8_multi_array_proto"],
)

proto_library(
    name = "multi_array_dimension_proto",
    srcs = ["multi_array_dimension.proto"],
)

py_proto_library(
    name = "multi_array_dimension_py_pb2",
    deps = ["multi_array_dimension_proto"],
)

cc_proto_library(
    name = "multi_array_dimension_cc_proto",
    deps = ["multi_array_dimension_proto"],
)

go_proto_library(
    name = "multi_array_dimension_go_proto",
    deps = ["multi_array_dimension_proto"],
)

proto_library(
    name = "multi_array_layout_proto",
    srcs = ["multi_array_layout.proto"],
    deps = [":multi_array_dimension_proto"],
)

py_proto_library(
    name = "multi_array_layout_py_pb2",
    deps = ["multi_array_layout_proto"],
)

cc_proto_library(
    name = "multi_array_layout_cc_proto",
    deps = ["multi_array_layout_proto"],
)

go_proto_library(
    name = "multi_array_layout_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:multi_array_dimension_go_proto",
    ],
    deps = ["multi_array_layout_proto"],
)

proto_library(
    name = "string_proto",
    srcs = ["string.proto"],
)

py_proto_library(
    name = "string_py_pb2",
    deps = ["string_proto"],
)

cc_proto_library(
    name = "string_cc_proto",
    deps = ["string_proto"],
)

go_proto_library(
    name = "string_go_proto",
    deps = ["string_proto"],
)

proto_library(
    name = "u_int16_proto",
    srcs = ["u_int16.proto"],
)

py_proto_library(
    name = "u_int16_py_pb2",
    deps = ["u_int16_proto"],
)

cc_proto_library(
    name = "u_int16_cc_proto",
    deps = ["u_int16_proto"],
)

go_proto_library(
    name = "u_int16_go_proto",
    deps = ["u_int16_proto"],
)

proto_library(
    name = "u_int16_multi_array_proto",
    srcs = ["u_int16_multi_array.proto"],
    deps = [":multi_array_layout_proto"],
)

py_proto_library(
    name = "u_int16_multi_array_py_pb2",
    deps = ["u_int16_multi_array_proto"],
)

cc_proto_library(
    name = "u_int16_multi_array_cc_proto",
    deps = ["u_int16_multi_array_proto"],
)

go_proto_library(
    name = "u_int16_multi_array_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:multi_array_layout_go_proto",
    ],
    deps = ["u_int16_multi_array_proto"],
)

proto_library(
    name = "u_int32_proto",
    srcs = ["u_int32.proto"],
)

py_proto_library(
    name = "u_int32_py_pb2",
    deps = ["u_int32_proto"],
)

cc_proto_library(
    name = "u_int32_cc_proto",
    deps = ["u_int32_proto"],
)

go_proto_library(
    name = "u_int32_go_proto",
    deps = ["u_int32_proto"],
)

proto_library(
    name = "u_int32_multi_array_proto",
    srcs = ["u_int32_multi_array.proto"],
    deps = [":multi_array_layout_proto"],
)

py_proto_library(
    name = "u_int32_multi_array_py_pb2",
    deps = ["u_int32_multi_array_proto"],
)

cc_proto_library(
    name = "u_int32_multi_array_cc_proto",
    deps = ["u_int32_multi_array_proto"],
)

go_proto_library(
    name = "u_int32_multi_array_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:multi_array_layout_go_proto",
    ],
    deps = ["u_int32_multi_array_proto"],
)

proto_library(
    name = "u_int64_proto",
    srcs = ["u_int64.proto"],
)

py_proto_library(
    name = "u_int64_py_pb2",
    deps = ["u_int64_proto"],
)

cc_proto_library(
    name = "u_int64_cc_proto",
    deps = ["u_int64_proto"],
)

go_proto_library(
    name = "u_int64_go_proto",
    deps = ["u_int64_proto"],
)

proto_library(
    name = "u_int64_multi_array_proto",
    srcs = ["u_int64_multi_array.proto"],
    deps = [":multi_array_layout_proto"],
)

py_proto_library(
    name = "u_int64_multi_array_py_pb2",
    deps = ["u_int64_multi_array_proto"],
)

cc_proto_library(
    name = "u_int64_multi_array_cc_proto",
    deps = ["u_int64_multi_array_proto"],
)

go_proto_library(
    name = "u_int64_multi_array_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:multi_array_layout_go_proto",
    ],
    deps = ["u_int64_multi_array_proto"],
)

proto_library(
    name = "u_int8_proto",
    srcs = ["u_int8.proto"],
)

py_proto_library(
    name = "u_int8_py_pb2",
    deps = ["u_int8_proto"],
)

cc_proto_library(
    name = "u_int8_cc_proto",
    deps = ["u_int8_proto"],
)

go_proto_library(
    name = "u_int8_go_proto",
    deps = ["u_int8_proto"],
)

proto_library(
    name = "u_int8_multi_array_proto",
    srcs = ["u_int8_multi_array.proto"],
    deps = [":multi_array_layout_proto"],
)

py_proto_library(
    name = "u_int8_multi_array_py_pb2",
    deps = ["u_int8_multi_array_proto"],
)

cc_proto_library(
    name = "u_int8_multi_array_cc_proto",
    deps = ["u_int8_multi_array_proto"],
)

go_proto_library(
    name = "u_int8_multi_array_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:multi_array_layout_go_proto",
    ],
    deps = ["u_int8_multi_array_proto"],
)
