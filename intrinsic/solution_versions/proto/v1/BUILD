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

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@intrinsic_js_rules//:defs.bzl", "js_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

js_library(
    name = "solution_version_service_jspb_proto",
    visibility = [
        "//intrinsic/frontend:__subpackages__",
        "//intrinsic/portal/frontend:__pkg__",
    ],
    deps = [":solution_version_service_proto"],
)

js_library(
    name = "snapshot_jspb_proto",
    visibility = [
        "//intrinsic/frontend:__subpackages__",
        "//intrinsic/portal/frontend:__pkg__",
    ],
    deps = ["//intrinsic/solution_versions/proto/v1:snapshot_proto"],
)

js_library(
    name = "snapshot_source_jspb_proto",
    visibility = [
        "//intrinsic/frontend:__subpackages__",
        "//intrinsic/portal/frontend:__pkg__",
    ],
    deps = ["//intrinsic/solution_versions/proto/v1:snapshot_source_proto"],
)

js_library(
    name = "branch_jspb_proto",
    visibility = [
        "//intrinsic/frontend:__subpackages__",
    ],
    deps = ["//intrinsic/solution_versions/proto/v1:branch_proto"],
)

proto_library(
    name = "solution_version_service_proto",
    srcs = ["solution_version_service.proto"],
    visibility = [
    ],
    deps = [
        "//intrinsic/solution_versions/proto/v1:branch_proto",
        "//intrinsic/solution_versions/proto/v1:conflict_proto",
        "//intrinsic/solution_versions/proto/v1:conflict_resolution_proto",
        "//intrinsic/solution_versions/proto/v1:snapshot_proto",
        "//intrinsic/solution_versions/proto/v1:snapshot_source_proto",
        "//intrinsic/solution_versions/proto/v1:warning_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/api:field_behavior_proto",
        "@com_google_googleapis//google/api:resource_proto",
        "@com_google_googleapis//google/longrunning:operations_proto",
        "@com_google_googleapis//google/rpc:status_proto",
        "@com_google_protobuf//:descriptor_proto",
        "@com_google_protobuf//:field_mask_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

py_proto_library(
    name = "solution_version_service_py_pb2",
    visibility = ["//visibility:public"],
    deps = [":solution_version_service_proto"],
)

py_grpc_library(
    name = "solution_version_service_pb2_grpc",
    srcs = [":solution_version_service_proto"],
    grpc_library = requirement("grpcio"),
    visibility = ["//visibility:public"],
    deps = [":solution_version_service_py_pb2"],
)

go_proto_library(
    name = "solution_version_service_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/solution_versions/proto/v1/solution_version_service_go_proto",
    protos = [":solution_version_service_proto"],
    visibility = [
        "//intrinsic/assets/deploy:__pkg__",
        "//intrinsic/assets/install:__pkg__",
        "//intrinsic/frontend/cloud:__subpackages__",
        "//intrinsic/frontend/go/api:__subpackages__",
        "//intrinsic/kubernetes/data_store/http_endpoints:__pkg__",
        "//intrinsic/kubernetes/workcell_spec:__pkg__",
        "//intrinsic/performance/inctl:__subpackages__",
        "//intrinsic/solution_copy:__subpackages__",
        "//intrinsic/solution_versions:__subpackages__",
        "//intrinsic/tools/inctl/cmd:__subpackages__",
    ],
    deps = [
        "//intrinsic/solution_versions/proto/v1:branch_go_proto",
        "//intrinsic/solution_versions/proto/v1:conflict_go_proto",
        "//intrinsic/solution_versions/proto/v1:conflict_resolution_go_proto",
        "//intrinsic/solution_versions/proto/v1:snapshot_go_proto",
        "//intrinsic/solution_versions/proto/v1:snapshot_source_go_proto",
        "//intrinsic/solution_versions/proto/v1:warning_go_proto",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@io_bazel_rules_go//proto/wkt:field_mask_go_proto",
        "@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
        "@org_golang_google_genproto//googleapis/longrunning",
        "@org_golang_google_genproto_googleapis_api//annotations",
        "@org_golang_google_genproto_googleapis_rpc//status",
    ],
)

js_library(
    name = "diff_mask_jspb_proto",
    visibility = [
        "//intrinsic/frontend:__subpackages__",
        "//intrinsic/portal/frontend:__pkg__",
    ],
    deps = ["//intrinsic/solution_versions/proto/v1:diff_mask_proto"],
)

js_library(
    name = "conflict_jspb_proto",
    visibility = [
        "//intrinsic/frontend:__subpackages__",
        "//intrinsic/portal/frontend:__pkg__",
    ],
    deps = ["//intrinsic/solution_versions/proto/v1:conflict_proto"],
)

js_library(
    name = "conflict_resolution_jspb_proto",
    visibility = [
        "//intrinsic/frontend:__subpackages__",
        "//intrinsic/portal/frontend:__pkg__",
    ],
    deps = ["//intrinsic/solution_versions/proto/v1:conflict_resolution_proto"],
)

js_library(
    name = "warning_jspb_proto",
    visibility = [
        "//intrinsic/frontend:__subpackages__",
    ],
    deps = ["//intrinsic/solution_versions/proto/v1:warning_proto"],
)

proto_library(
    name = "solution_diff_proto",
    srcs = ["solution_diff.proto"],
    visibility = [
    ],
    deps = [
        "//intrinsic/solution_versions/proto/v1:diff_mask_proto",
    ],
)

go_proto_library(
    name = "solution_diff_go_proto",
    importpath = "intrinsic/solution_versions/proto/v1/solution_diff_go_proto",
    protos = [":solution_diff_proto"],
    visibility = [
        "//intrinsic/frontend/cloud/portal_bridge:__subpackages__",
        "//intrinsic/solution_versions:__subpackages__",
        "//intrinsic/tools/inctl/cmd/solution_version:__pkg__",
    ],
    deps = [
        "//intrinsic/solution_versions/proto/v1:diff_mask_go_proto",
    ],
)

py_proto_library(
    name = "solution_diff_py_pb2",
    visibility = ["//visibility:public"],
    deps = [":solution_diff_proto"],
)

js_library(
    name = "solution_diff_jspb_proto",
    visibility = [
        "//intrinsic/frontend:__subpackages__",
        "//intrinsic/portal/frontend:__pkg__",
    ],
    deps = [":solution_diff_proto"],
)

proto_library(
    name = "branch_proto",
    srcs = ["branch.proto"],
    deps = [
        "//intrinsic/config/proto:operation_mode_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

go_proto_library(
    name = "branch_go_proto",
    importpath = "intrinsic/solution_versions/proto/v1/branch_go_proto",
    protos = [":branch_proto"],
    deps = ["//intrinsic/config/proto:operation_mode_go_proto"],
)

py_proto_library(
    name = "branch_py_pb2",
    deps = [":branch_proto"],
)

proto_library(
    name = "conflict_proto",
    srcs = ["conflict.proto"],
)

go_proto_library(
    name = "conflict_go_proto",
    importpath = "intrinsic/solution_versions/proto/v1/conflict_go_proto",
    protos = [":conflict_proto"],
)

py_proto_library(
    name = "conflict_py_pb2",
    deps = [":conflict_proto"],
)

proto_library(
    name = "conflict_resolution_proto",
    srcs = ["conflict_resolution.proto"],
)

go_proto_library(
    name = "conflict_resolution_go_proto",
    importpath = "intrinsic/solution_versions/proto/v1/conflict_resolution_go_proto",
    protos = [":conflict_resolution_proto"],
)

py_proto_library(
    name = "conflict_resolution_py_pb2",
    deps = [":conflict_resolution_proto"],
)

proto_library(
    name = "diff_mask_proto",
    srcs = ["diff_mask.proto"],
)

go_proto_library(
    name = "diff_mask_go_proto",
    importpath = "intrinsic/solution_versions/proto/v1/diff_mask_go_proto",
    protos = [":diff_mask_proto"],
)

py_proto_library(
    name = "diff_mask_py_pb2",
    deps = [":diff_mask_proto"],
)

proto_library(
    name = "snapshot_proto",
    srcs = ["snapshot.proto"],
    deps = ["@com_google_protobuf//:timestamp_proto"],
)

go_proto_library(
    name = "snapshot_go_proto",
    importpath = "intrinsic/solution_versions/proto/v1/snapshot_go_proto",
    protos = [":snapshot_proto"],
)

py_proto_library(
    name = "snapshot_py_pb2",
    deps = [":snapshot_proto"],
)

proto_library(
    name = "snapshot_source_proto",
    srcs = ["snapshot_source.proto"],
)

go_proto_library(
    name = "snapshot_source_go_proto",
    importpath = "intrinsic/solution_versions/proto/v1/snapshot_source_go_proto",
    protos = [":snapshot_source_proto"],
)

py_proto_library(
    name = "snapshot_source_py_pb2",
    deps = [":snapshot_source_proto"],
)

proto_library(
    name = "warning_proto",
    srcs = ["warning.proto"],
    deps = ["//intrinsic/assets/proto:id_proto"],
)

go_proto_library(
    name = "warning_go_proto",
    importpath = "intrinsic/solution_versions/proto/v1/warning_go_proto",
    protos = [":warning_proto"],
    deps = [
        "//intrinsic/assets/proto:id_go_proto",
    ],
)

py_proto_library(
    name = "warning_py_pb2",
    deps = [":warning_proto"],
)
