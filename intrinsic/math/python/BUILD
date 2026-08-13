# Copyright 2026 Intrinsic Innovation LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@rules_python//python:defs.bzl", "py_library", "py_test")

package(default_visibility = [
    "//visibility:public",
])

py_library(
    name = "proto_conversion",
    srcs = ["proto_conversion.py"],
    deps = [
        ":data_types",
        "//intrinsic/icon/proto:cart_space_py_pb2",
        "//intrinsic/math/proto:array_py_pb2",
        "//intrinsic/math/proto:matrix_py_pb2",
        "//intrinsic/math/proto:point_py_pb2",
        "//intrinsic/math/proto:pose_py_pb2",
        "//intrinsic/math/proto:quaternion_py_pb2",
        "//intrinsic/math/proto:twist_py_pb2",
        "//intrinsic/math/proto:vector3_py_pb2",
        "@ai_intrinsic_sdks_pip_deps//numpy",
    ],
)

py_test(
    name = "proto_conversion_test",
    srcs = ["proto_conversion_test.py"],
    deps = [
        ":data_types",
        ":proto_conversion",
        "//intrinsic/icon/proto:cart_space_py_pb2",
        "//intrinsic/math/proto:array_py_pb2",
        "//intrinsic/math/proto:matrix_py_pb2",
        "//intrinsic/math/proto:point_py_pb2",
        "//intrinsic/math/proto:pose_py_pb2",
        "//intrinsic/math/proto:quaternion_py_pb2",
        "//intrinsic/math/proto:vector3_py_pb2",
        "@ai_intrinsic_sdks_pip_deps//hypothesis",
        "@ai_intrinsic_sdks_pip_deps//numpy",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
    ],
)

py_library(
    name = "ros_proto_conversion",
    srcs = ["ros_proto_conversion.py"],
    deps = [
        ":data_types",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:point_py_pb2",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_py_pb2",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:quaternion_py_pb2",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:twist_py_pb2",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:vector3_py_pb2",
        "@ai_intrinsic_sdks_pip_deps//numpy",
    ],
)

py_test(
    name = "ros_proto_conversion_test",
    srcs = ["ros_proto_conversion_test.py"],
    deps = [
        ":data_types",
        ":ros_proto_conversion",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:point_py_pb2",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:pose_py_pb2",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:quaternion_py_pb2",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:twist_py_pb2",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:vector3_py_pb2",
        "@ai_intrinsic_sdks_pip_deps//numpy",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
    ],
)

py_library(
    name = "data_types",
    srcs = ["data_types.py"],
    deps = [
        ":math_types",
        ":pose3",
        ":quaternion",
        ":rotation3",
        ":vector_util",
        "//intrinsic/icon/proto:cart_space_py_pb2",
        "@ai_intrinsic_sdks_pip_deps//numpy",
    ],
)

py_test(
    name = "data_types_test",
    size = "small",
    srcs = ["data_types_test.py"],
    deps = [
        ":data_types",
        "//intrinsic/icon/proto:cart_space_py_pb2",
        "@ai_intrinsic_sdks_pip_deps//numpy",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

py_library(
    name = "math_test",
    testonly = 1,
    srcs = ["math_test.py"],
    deps = [
        ":math_types",
        ":pose3",
        ":quaternion",
        ":rotation3",
        ":vector_util",
        "@ai_intrinsic_sdks_pip_deps//numpy",
        "@com_google_absl_py//absl/logging",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

py_test(
    name = "math_test_test",
    size = "small",
    srcs = [
        "math_test_test.py",
    ],
    deps = [
        ":math_test",
        ":math_types",
        ":pose3",
        ":rotation3",
        ":vector_util",
        "@ai_intrinsic_sdks_pip_deps//numpy",
        "@com_google_absl_py//absl/logging",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
    ],
)

py_library(
    name = "math_types",
    srcs = ["math_types.py"],
    deps = ["@ai_intrinsic_sdks_pip_deps//numpy"],
)

py_test(
    name = "math_types_test",
    size = "small",
    srcs = [
        "math_types_test.py",
    ],
    deps = [
        ":math_test",
        ":math_types",
        "@ai_intrinsic_sdks_pip_deps//numpy",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
    ],
)

py_library(
    name = "pose3",
    srcs = [
        "pose3.py",
    ],
    deps = [
        ":math_types",
        ":quaternion",
        ":rotation3",
        ":vector_util",
        "@ai_intrinsic_sdks_pip_deps//numpy",
    ],
)

py_test(
    name = "pose3_test",
    size = "medium",
    srcs = [
        "pose3_test.py",
    ],
    deps = [
        ":math_test",
        ":pose3",
        ":quaternion",
        ":rotation3",
        ":vector_util",
        "@ai_intrinsic_sdks_pip_deps//numpy",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
    ],
)

py_library(
    name = "rotation3",
    srcs = [
        "rotation3.py",
    ],
    deps = [
        ":math_types",
        ":quaternion",
        ":vector_util",
        "@ai_intrinsic_sdks_pip_deps//numpy",
    ],
)

py_test(
    name = "rotation3_test",
    size = "small",
    timeout = "moderate",
    srcs = [
        "rotation3_test.py",
    ],
    deps = [
        ":math_test",
        ":math_types",
        ":quaternion",
        ":rotation3",
        ":vector_util",
        "@ai_intrinsic_sdks_pip_deps//numpy",
        "@com_google_absl_py//absl/logging",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
    ],
)

py_library(
    name = "quaternion",
    srcs = [
        "quaternion.py",
    ],
    deps = [
        ":math_types",
        ":vector_util",
        "@ai_intrinsic_sdks_pip_deps//numpy",
    ],
)

py_test(
    name = "quaternion_test",
    size = "small",
    srcs = [
        "quaternion_test.py",
    ],
    deps = [
        ":math_test",
        ":math_types",
        ":quaternion",
        ":vector_util",
        "@ai_intrinsic_sdks_pip_deps//numpy",
        "@com_google_absl_py//absl/logging",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_absl_py//absl/testing:parameterized",
    ],
)

py_library(
    name = "vector_util",
    srcs = [
        "vector_util.py",
    ],
    deps = [
        ":math_types",
        "@ai_intrinsic_sdks_pip_deps//numpy",
    ],
)
