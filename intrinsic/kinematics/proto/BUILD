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
    name = "kinematics_proto",
    srcs = ["kinematics.proto"],
    visibility = ["//visibility:public"],
    deps = [
        "//intrinsic/math/proto:matrix_proto",
        "//intrinsic/math/proto:point_proto",
        "//intrinsic/math/proto:pose_proto",
    ],
)

proto_library(
    name = "skeleton_proto",
    srcs = ["skeleton.proto"],
    visibility = ["//visibility:public"],
    deps = [
        ":kinematics_proto",
        "//intrinsic/math/proto:pose_proto",
    ],
)

cc_proto_library(
    name = "kinematics_cc_proto",
    deps = [":kinematics_proto"],
)

go_proto_library(
    name = "kinematics_go_proto",
    importpath = "intrinsic/kinematics/proto/kinematics_go_proto",
    protos = [":kinematics_proto"],
    deps = [
        "//intrinsic/math/proto:matrix_go_proto",
        "//intrinsic/math/proto:point_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
    ],
)

py_proto_library(
    name = "kinematics_py_proto",
    deps = [":kinematics_proto"],
)

cc_proto_library(
    name = "skeleton_cc_proto",
    deps = [":skeleton_proto"],
)

go_proto_library(
    name = "skeleton_go_proto",
    importpath = "intrinsic/kinematics/proto/skeleton_go_proto",
    protos = [":skeleton_proto"],
    deps = [
        ":kinematics_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
    ],
)

py_proto_library(
    name = "skeleton_py_proto",
    deps = [":skeleton_proto"],
)
