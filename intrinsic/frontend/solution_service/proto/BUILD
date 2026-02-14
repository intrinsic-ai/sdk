# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")
load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "proto_source_code_info_transitive_descriptor_set")

package(default_visibility = ["//intrinsic:public_api_users"])

proto_library(
    name = "solution_service_proto",
    srcs = ["solution_service.proto"],
    deps = [
        ":status_proto",
        "//intrinsic/executive/proto:behavior_tree_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_protobuf//:empty_proto",
    ],
)

proto_source_code_info_transitive_descriptor_set(
    name = "solution_service_proto_descriptor_set",
    deps = [":solution_service_proto"],
)

cc_proto_library(
    name = "solution_service_cc_proto",
    deps = [":solution_service_proto"],
)

cc_grpc_library(
    name = "solution_service_cc_grpc_proto",
    srcs = [":solution_service_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":solution_service_cc_proto"],
)

py_proto_library(
    name = "solution_service_py_pb2",
    deps = [":solution_service_proto"],
)

py_grpc_library(
    name = "solution_service_py_pb2_grpc",
    srcs = [":solution_service_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":solution_service_py_pb2"],
)

go_proto_library(
    name = "solution_service_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
        "//bazel:go_gen_grpc_gateway",
    ],
    importpath = "intrinsic/frontend/solution_service/proto/solution_service_go_proto",
    protos = [":solution_service_proto"],
    deps = [
        ":status_go_proto",
        "//intrinsic/executive/proto:behavior_tree_go_proto",
        "@org_golang_google_genproto_googleapis_api//annotations",
        "@org_golang_google_protobuf//types/known/emptypb",
    ],
)

proto_library(
    name = "status_proto",
    srcs = ["status.proto"],
)

cc_proto_library(
    name = "status_cc_proto",
    deps = [":status_proto"],
)

py_proto_library(
    name = "status_py_pb2",
    deps = [":status_proto"],
)

go_proto_library(
    name = "status_go_proto",
    importpath = "intrinsic/frontend/solution_service/proto/status_go_proto",
    protos = [":status_proto"],
)
