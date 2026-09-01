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

package(
    default_visibility = ["//visibility:public"],
)

js_library(
    name = "signal_jspb_proto",
    visibility = ["//intrinsic:intrinsic_internal"],
    deps = ["//intrinsic/hardware/gpio/v1:signal_proto"],
)

proto_library(
    name = "signal_proto",
    srcs = ["signal.proto"],
)

cc_proto_library(
    name = "signal_cc_proto",
    deps = [":signal_proto"],
)

go_proto_library(
    name = "signal_go_proto",
    importpath = "intrinsic/hardware/gpio/v1/signal_go_proto",
    protos = [":signal_proto"],
)

py_proto_library(
    name = "signal_py_pb2",
    deps = [":signal_proto"],
)

proto_library(
    name = "gpio_service_proto",
    srcs = ["gpio_service.proto"],
    deps = [
        ":signal_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/rpc:status_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

cc_proto_library(
    name = "gpio_service_cc_proto",
    deps = [":gpio_service_proto"],
)

cc_grpc_library(
    name = "gpio_service_cc_grpc_proto",
    srcs = [":gpio_service_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":gpio_service_cc_proto"],
)

go_proto_library(
    name = "gpio_service_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-grpc-gateway:go_gen_grpc_gateway",
    ],
    importpath = "intrinsic/hardware/gpio/v1/gpio_service_go_proto",
    protos = [":gpio_service_proto"],
    deps = [
        ":signal_go_proto",
        "@org_golang_google_genproto_googleapis_api//annotations",
        "@org_golang_google_genproto_googleapis_rpc//status",
    ],
)

py_proto_library(
    name = "gpio_service_py_pb2",
    deps = [":gpio_service_proto"],
)

py_grpc_library(
    name = "gpio_service_py_pb2_grpc",
    srcs = [":gpio_service_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":gpio_service_py_pb2"],
)
