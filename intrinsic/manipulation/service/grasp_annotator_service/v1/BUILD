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

# Grasp annotator service definition.

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")
load("//intrinsic/httpjson:intrinsic_http_image.bzl", "intrinsic_http_image")
load("//intrinsic/httpjson/openapi:protoc_gen_openapi.bzl", "protoc_gen_openapi")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

go_proto_library(
    name = "grasp_annotator_service_go_proto",
    compilers = [
        "//bazel:go_gen_grpc_gateway",
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/manipulation/service/grasp_annotator_service/v1/grasp_annotator_service_go_proto",
    protos = [":grasp_annotator_service_proto"],
    deps = [
        "//intrinsic/logging/proto:context_go_proto",
        "//intrinsic/manipulation/grasping:grasp_annotations_go_proto",
        "//intrinsic/manipulation/grasping:grasp_annotator_go_proto",
        "@org_golang_google_genproto_googleapis_api//annotations",
        "@org_golang_google_genproto_googleapis_rpc//status",
    ],
)

cc_proto_library(
    name = "grasp_annotator_service_cc_proto",
    deps = [":grasp_annotator_service_proto"],
)

proto_library(
    name = "grasp_annotator_service_proto",
    srcs = ["grasp_annotator_service.proto"],
    deps = [
        "//intrinsic/logging/proto:context_proto",
        "//intrinsic/manipulation/grasping:grasp_annotations_proto",
        "//intrinsic/manipulation/grasping:grasp_annotator_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/rpc:status_proto",
    ],
)

py_proto_library(
    name = "grasp_annotator_service_py_pb2",
    deps = [":grasp_annotator_service_proto"],
)

py_grpc_library(
    name = "grasp_annotator_service_py_pb2_grpc",
    srcs = [":grasp_annotator_service_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":grasp_annotator_service_py_pb2"],
)

protoc_gen_openapi(
    name = "grasp_annotator_service_openapi",
    protos = [
        ":grasp_annotator_service_proto",
    ],
)

intrinsic_http_image(
    name = "grasp_annotator_service_httpjson_bridge",
    openapi_yaml = ":grasp_annotator_service_openapi",
    services = {
        ":grasp_annotator_service_go_proto": "intrinsic_proto.grasping.service.grasp_annotator_service.v1.GraspAnnotator",
    },
)
