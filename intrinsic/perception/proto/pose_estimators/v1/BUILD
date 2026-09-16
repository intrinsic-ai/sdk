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

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "charuco_marker_pose_estimation_config_proto",
    srcs = ["charuco_marker_pose_estimation_config.proto"],
)

cc_proto_library(
    name = "charuco_marker_pose_estimation_config_cc_proto",
    deps = [":charuco_marker_pose_estimation_config_proto"],
)

go_proto_library(
    name = "charuco_marker_pose_estimation_config_go_proto",
    importpath = "intrinsic/perception/proto/pose_estimators/v1/charuco_marker_pose_estimation_config_go_proto",
    protos = [":charuco_marker_pose_estimation_config_proto"],
)

py_proto_library(
    name = "charuco_marker_pose_estimation_config_py_pb2",
    deps = [":charuco_marker_pose_estimation_config_proto"],
)

proto_library(
    name = "one_shot_base_model_config_proto",
    srcs = ["one_shot_base_model_config.proto"],
)

cc_proto_library(
    name = "one_shot_base_model_config_cc_proto",
    deps = [":one_shot_base_model_config_proto"],
)

go_proto_library(
    name = "one_shot_base_model_config_go_proto",
    importpath = "intrinsic/perception/proto/pose_estimators/v1/one_shot_base_model_config_go_proto",
    protos = [":one_shot_base_model_config_proto"],
)

py_proto_library(
    name = "one_shot_base_model_config_py_pb2",
    deps = [":one_shot_base_model_config_proto"],
)
