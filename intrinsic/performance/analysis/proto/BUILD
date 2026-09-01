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

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "event_proto",
    srcs = ["event.proto"],
    deps = [
        "@com_google_protobuf//:struct_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

proto_library(
    name = "performance_metrics_proto",
    srcs = ["performance_metrics.proto"],
    deps = ["@com_google_protobuf//:struct_proto"],
)

cc_proto_library(
    name = "event_cc_proto",
    deps = [":event_proto"],
)

py_proto_library(
    name = "event_py_pb2",
    deps = [":event_proto"],
)

cc_proto_library(
    name = "performance_metrics_cc_proto",
    deps = [":performance_metrics_proto"],
)

go_proto_library(
    name = "performance_metrics_go_proto",
    importpath = "intrinsic/performance/analysis/proto/performance_metrics_go_proto",
    protos = [":performance_metrics_proto"],
)

py_proto_library(
    name = "performance_metrics_py_pb2",
    deps = [":performance_metrics_proto"],
)
