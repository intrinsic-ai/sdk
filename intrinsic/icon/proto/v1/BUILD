# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

proto_library(
    name = "service_proto",
    srcs = ["service.proto"],
    deps = [
        ":types_proto",
        "//intrinsic/icon/proto:joint_space_proto",
        "//intrinsic/icon/proto:logging_mode_proto",
        "//intrinsic/icon/proto:part_status_proto",
        "//intrinsic/icon/proto:safety_status_proto",
        "//intrinsic/icon/proto:streaming_output_proto",
        "//intrinsic/logging/proto:context_proto",
        "//intrinsic/world/proto:robot_payload_proto",
        "@com_google_googleapis//google/rpc:status_proto",
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:duration_proto",
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

cc_proto_library(
    name = "service_cc_proto",
    deps = [":service_proto"],
)

go_proto_library(
    name = "service_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/icon/proto/v1/service_go_proto",
    protos = [":service_proto"],
    deps = [
        "//intrinsic/icon/proto:joint_space_go_proto",
        "//intrinsic/icon/proto:logging_mode_go_proto",
        "//intrinsic/icon/proto:part_status_go_proto",
        "//intrinsic/icon/proto:safety_status_go_proto",
        "//intrinsic/icon/proto:streaming_output_go_proto",
        "//intrinsic/icon/proto/v1:types_go_proto",
        "//intrinsic/logging/proto:context_go_proto",
        "//intrinsic/world/proto:robot_payload_go_proto",
        "@org_golang_google_genproto_googleapis_rpc//status",
        "@org_golang_google_protobuf//types/known/anypb",
        "@org_golang_google_protobuf//types/known/durationpb",
        "@org_golang_google_protobuf//types/known/emptypb",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

py_proto_library(
    name = "service_py_pb2",
    deps = [":service_proto"],
)

cc_grpc_library(
    name = "service_cc_grpc_proto",
    srcs = [":service_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":service_cc_proto"],
)

py_grpc_library(
    name = "service_py_pb2_grpc",
    srcs = [":service_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":service_py_pb2"],
)

proto_library(
    name = "types_proto",
    srcs = ["types.proto"],
    deps = [
        ":condition_types_proto",
        "//intrinsic/icon/proto:generic_part_config_proto",
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:descriptor_proto",
    ],
)

cc_proto_library(
    name = "types_cc_proto",
    deps = [":types_proto"],
)

go_proto_library(
    name = "types_go_proto",
    importpath = "intrinsic/icon/proto/v1/types_go_proto",
    protos = [":types_proto"],
    deps = [
        ":condition_types_go_proto",
        "//intrinsic/icon/proto:generic_part_config_go_proto",
    ],
)

py_proto_library(
    name = "types_py_pb2",
    deps = [":types_proto"],
)

proto_library(
    name = "condition_types_proto",
    srcs = ["condition_types.proto"],
)

cc_proto_library(
    name = "condition_types_cc_proto",
    deps = [":condition_types_proto"],
)

go_proto_library(
    name = "condition_types_go_proto",
    importpath = "intrinsic/icon/proto/v1/condition_types_go_proto",
    protos = [":condition_types_proto"],
)

py_proto_library(
    name = "condition_types_py_pb2",
    deps = [":condition_types_proto"],
)

proto_library(
    name = "jogging_service_proto",
    srcs = ["jogging_service.proto"],
    deps = [
        "//intrinsic/icon/proto:cart_space_proto",
        "//intrinsic/kinematics/types:joint_limits_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
        "@com_google_protobuf//:duration_proto",
    ],
)

cc_proto_library(
    name = "jogging_service_cc_proto",
    deps = [":jogging_service_proto"],
)

cc_grpc_library(
    name = "jogging_service_cc_grpc_proto",
    srcs = [":jogging_service_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":jogging_service_cc_proto"],
)

go_proto_library(
    name = "jogging_service_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/icon/proto/v1/jogging_service_go_proto",
    protos = [":jogging_service_proto"],
    deps = [
        "//intrinsic/icon/proto:cart_space_go_proto",
        "//intrinsic/kinematics/types:joint_limits_go_proto",
        "//intrinsic/world/proto:object_world_refs_go_proto",
        "@org_golang_google_protobuf//types/known/durationpb",
    ],
)
