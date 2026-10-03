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
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@intrinsic_js_rules//:defs.bzl", "js_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "asset_proto",
    srcs = ["asset.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":processed_asset_proto",
        "//intrinsic/assets/proto/v1:reference_proto",
    ],
)

cc_proto_library(
    name = "asset_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":asset_proto"],
)

py_proto_library(
    name = "asset_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":asset_proto"],
)

go_proto_library(
    name = "asset_go_proto",
    importpath = "intrinsic/assets/proto/v1/asset_go_proto",
    protos = [":asset_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":processed_asset_go_proto",
        "//intrinsic/assets/proto/v1:reference_go_proto",
    ],
)

js_library(
    name = "asset_jspb_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":asset_proto"],
)

js_library(
    name = "asset_graph_jspb_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["//intrinsic/assets/proto/v1:asset_graph_proto"],
)

js_library(
    name = "asset_instances_jspb_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["//intrinsic/assets/proto/v1:asset_instances_proto"],
)

js_library(
    name = "grpc_connection_jspb_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["//intrinsic/assets/proto/v1:grpc_connection_proto"],
)

js_library(
    name = "instance_config_jspb_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["//intrinsic/assets/proto/v1:instance_config_proto"],
)

proto_library(
    name = "processed_asset_proto",
    srcs = ["processed_asset.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/data/proto/v1:data_asset_proto",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_proto",
        "//intrinsic/assets/processes/proto:process_asset_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_proto",
        "//intrinsic/assets/services/proto:service_manifest_proto",
        "//intrinsic/skills/proto:processed_skill_manifest_proto",
    ],
)

cc_proto_library(
    name = "processed_asset_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":processed_asset_proto"],
)

py_proto_library(
    name = "processed_asset_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":processed_asset_proto"],
)

go_proto_library(
    name = "processed_asset_go_proto",
    importpath = "intrinsic/assets/proto/v1/processed_asset_go_proto",
    protos = [":processed_asset_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_go_proto",
        "//intrinsic/assets/processes/proto:process_asset_go_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/skills/proto:processed_skill_manifest_go_proto",
    ],
)

js_library(
    name = "processed_asset_jspb_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":processed_asset_proto"],
)

js_library(
    name = "reference_jspb_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["//intrinsic/assets/proto/v1:reference_proto"],
)

js_library(
    name = "search_jspb_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["//intrinsic/assets/proto/v1:search_proto"],
)

js_library(
    name = "resolved_dependency_jspb_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["//intrinsic/assets/proto/v1:resolved_dependency_proto"],
)

proto_library(
    name = "solution_proto",
    srcs = ["solution.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":asset_proto",
        "//intrinsic/assets/proto/v1:asset_instances_proto",
        "//intrinsic/world/proto:object_world_updates_proto",
    ],
)

cc_proto_library(
    name = "solution_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":solution_proto"],
)

py_proto_library(
    name = "solution_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":solution_proto"],
)

go_proto_library(
    name = "solution_go_proto",
    importpath = "intrinsic/assets/proto/v1/solution_go_proto",
    protos = [":solution_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":asset_go_proto",
        "//intrinsic/assets/proto/v1:asset_instances_go_proto",
        "//intrinsic/world/proto:object_world_updates_go_proto",
    ],
)

js_library(
    name = "solution_jspb_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":solution_proto"],
)

proto_library(
    name = "solution_deployment_proto",
    srcs = ["solution_deployment.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":solution_proto",
        "//intrinsic/config/proto:operation_mode_proto",
        "//intrinsic/util/status:extended_status_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/longrunning:operations_proto",
    ],
)

go_proto_library(
    name = "solution_deployment_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/assets/proto/v1/solution_deployment_go_proto",
    protos = [":solution_deployment_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":solution_go_proto",
        "//intrinsic/config/proto:operation_mode_go_proto",
        "//intrinsic/util/status:extended_status_go_proto",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)

cc_proto_library(
    name = "solution_deployment_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":solution_deployment_proto"],
)

py_proto_library(
    name = "solution_deployment_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":solution_deployment_proto"],
)

js_library(
    name = "solution_deployment_jspb_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":solution_deployment_proto"],
)

js_library(
    name = "dependency_jspb_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["//intrinsic/assets/proto/v1:dependency_proto"],
)

js_library(
    name = "asset_configuration_jspb_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["//intrinsic/assets/proto/v1:asset_configuration_proto"],
)

js_library(
    name = "recommended_configuration_config_jspb_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["//intrinsic/assets/proto/v1:recommended_configuration_config_proto"],
)

proto_library(
    name = "asset_artifacts_proto",
    srcs = ["asset_artifacts.proto"],
    deps = [
        "//intrinsic/assets/data/proto/v1:referenced_data_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/longrunning:operations_proto",
    ],
)

go_proto_library(
    name = "asset_artifacts_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-grpc-gateway:go_gen_grpc_gateway",
    ],
    importpath = "intrinsic/assets/proto/v1/asset_artifacts_go_proto",
    protos = [":asset_artifacts_proto"],
    deps = [
        "//intrinsic/assets/data/proto/v1:referenced_data_go_proto",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)

proto_library(
    name = "asset_configuration_proto",
    srcs = ["asset_configuration.proto"],
    deps = [
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_protobuf//:any_proto",
    ],
)

cc_proto_library(
    name = "asset_configuration_cc_proto",
    deps = [":asset_configuration_proto"],
)

go_proto_library(
    name = "asset_configuration_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-grpc-gateway:go_gen_grpc_gateway",
    ],
    importpath = "intrinsic/assets/proto/v1/asset_configuration_go_proto",
    protos = [":asset_configuration_proto"],
    deps = ["@org_golang_google_genproto_googleapis_api//annotations"],
)

py_proto_library(
    name = "asset_configuration_py_pb2",
    deps = [":asset_configuration_proto"],
)

py_grpc_library(
    name = "asset_configuration_pb2_grpc",
    srcs = [":asset_configuration_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":asset_configuration_py_pb2"],
)

proto_library(
    name = "asset_graph_proto",
    srcs = ["asset_graph.proto"],
)

cc_proto_library(
    name = "asset_graph_cc_proto",
    deps = [":asset_graph_proto"],
)

go_proto_library(
    name = "asset_graph_go_proto",
    importpath = "intrinsic/assets/proto/v1/asset_graph_go_proto",
    protos = [":asset_graph_proto"],
)

py_proto_library(
    name = "asset_graph_py_pb2",
    deps = [":asset_graph_proto"],
)

proto_library(
    name = "dependency_proto",
    srcs = ["dependency.proto"],
)

cc_proto_library(
    name = "dependency_cc_proto",
    deps = [":dependency_proto"],
)

go_proto_library(
    name = "dependency_go_proto",
    importpath = "intrinsic/assets/proto/v1/dependency_go_proto",
    protos = [":dependency_proto"],
)

py_proto_library(
    name = "dependency_py_pb2",
    deps = [":dependency_proto"],
)

proto_library(
    name = "grpc_connection_proto",
    srcs = ["grpc_connection.proto"],
)

cc_proto_library(
    name = "grpc_connection_cc_proto",
    deps = [":grpc_connection_proto"],
)

go_proto_library(
    name = "grpc_connection_go_proto",
    importpath = "intrinsic/assets/proto/v1/grpc_connection_go_proto",
    protos = [":grpc_connection_proto"],
)

py_proto_library(
    name = "grpc_connection_py_pb2",
    deps = [":grpc_connection_proto"],
)

proto_library(
    name = "instance_config_proto",
    srcs = ["instance_config.proto"],
    deps = [
        "//intrinsic/scene/proto/v1:scene_object_config_proto",
        "@com_google_protobuf//:any_proto",
    ],
)

cc_proto_library(
    name = "instance_config_cc_proto",
    deps = [":instance_config_proto"],
)

go_proto_library(
    name = "instance_config_go_proto",
    importpath = "intrinsic/assets/proto/v1/instance_config_go_proto",
    protos = [":instance_config_proto"],
    deps = [
        "//intrinsic/scene/proto/v1:scene_object_config_go_proto",
    ],
)

py_proto_library(
    name = "instance_config_py_pb2",
    deps = [":instance_config_proto"],
)

proto_library(
    name = "asset_instances_proto",
    srcs = ["asset_instances.proto"],
    deps = [
        ":dependency_proto",
        ":grpc_connection_proto",
        ":instance_config_proto",
        "//intrinsic/assets/proto:asset_tag_proto",
        "//intrinsic/assets/proto:asset_type_proto",
        "//intrinsic/assets/proto:id_proto",
        "//intrinsic/assets/proto:metadata_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/api:field_behavior_proto",
    ],
)

cc_proto_library(
    name = "asset_instances_cc_proto",
    deps = [":asset_instances_proto"],
)

cc_grpc_library(
    name = "asset_instances_cc_grpc",
    srcs = [":asset_instances_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":asset_instances_cc_proto"],
)

go_proto_library(
    name = "asset_instances_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/assets/proto/v1/asset_instances_go_proto",
    protos = [":asset_instances_proto"],
    deps = [
        ":dependency_go_proto",
        ":grpc_connection_go_proto",
        ":instance_config_go_proto",
        "//intrinsic/assets/proto:asset_tag_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)

py_proto_library(
    name = "asset_instances_py_pb2",
    deps = [":asset_instances_proto"],
)

py_grpc_library(
    name = "asset_instances_py_pb2_grpc",
    srcs = [":asset_instances_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":asset_instances_py_pb2"],
)

proto_library(
    name = "recommended_configuration_config_proto",
    srcs = ["recommended_configuration_config.proto"],
)

cc_proto_library(
    name = "recommended_configuration_config_cc_proto",
    deps = [":recommended_configuration_config_proto"],
)

go_proto_library(
    name = "recommended_configuration_config_go_proto",
    importpath = "intrinsic/assets/proto/v1/recommended_configuration_config_go_proto",
    protos = [":recommended_configuration_config_proto"],
)

py_proto_library(
    name = "recommended_configuration_config_py_pb2",
    deps = [":recommended_configuration_config_proto"],
)

proto_library(
    name = "reference_proto",
    srcs = ["reference.proto"],
    deps = [
        "//intrinsic/assets/proto:asset_type_proto",
        "//intrinsic/assets/proto:id_proto",
    ],
)

cc_proto_library(
    name = "reference_cc_proto",
    deps = [":reference_proto"],
)

go_proto_library(
    name = "reference_go_proto",
    importpath = "intrinsic/assets/proto/v1/reference_go_proto",
    protos = [":reference_proto"],
    deps = [
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
    ],
)

py_proto_library(
    name = "reference_py_pb2",
    deps = [":reference_proto"],
)

proto_library(
    name = "resolved_dependency_proto",
    srcs = ["resolved_dependency.proto"],
    deps = [
        ":grpc_connection_proto",
        "//intrinsic/assets/proto:id_proto",
        "@com_google_googleapis//google/api:field_behavior_proto",
    ],
)

cc_proto_library(
    name = "resolved_dependency_cc_proto",
    deps = [":resolved_dependency_proto"],
)

go_proto_library(
    name = "resolved_dependency_go_proto",
    importpath = "intrinsic/assets/proto/v1/resolved_dependency_go_proto",
    protos = [":resolved_dependency_proto"],
    deps = [
        ":grpc_connection_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)

py_proto_library(
    name = "resolved_dependency_py_pb2",
    deps = [":resolved_dependency_proto"],
)

proto_library(
    name = "search_proto",
    srcs = ["search.proto"],
)

cc_proto_library(
    name = "search_cc_proto",
    deps = [":search_proto"],
)

go_proto_library(
    name = "search_go_proto",
    importpath = "intrinsic/assets/proto/v1/search_go_proto",
    protos = [":search_proto"],
)

py_proto_library(
    name = "search_py_pb2",
    deps = [":search_proto"],
)
