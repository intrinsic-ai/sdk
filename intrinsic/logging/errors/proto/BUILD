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

# This package is to provide a consistent way throughout Intrinsic components to
# store structured data about errors.

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@intrinsic_js_rules//:defs.bzl", "js_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "error_report_proto",
    srcs = ["error_report.proto"],
    deps = [
        "@com_google_googleapis//google/rpc:status_proto",
        "@com_google_protobuf//:any_proto",
    ],
)

cc_proto_library(
    name = "error_report_cc_proto",
    deps = [":error_report_proto"],
)

py_proto_library(
    name = "error_report_py_pb2",
    deps = [":error_report_proto"],
)

go_proto_library(
    name = "error_report_go_proto",
    importpath = "intrinsic/logging/errors/proto/error_report_go_proto",
    protos = [":error_report_proto"],
    deps = [
        "//intrinsic/executive/proto:clips_snapshot_go_proto",
        "@org_golang_google_genproto_googleapis_rpc//status",
    ],
)

js_library(
    name = "error_report_jspb_proto",
    deps = [":error_report_proto"],
)
