# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

proto_library(
    name = "asset_graph_proto",
    srcs = ["asset_graph.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

cc_proto_library(
    name = "asset_graph_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":asset_graph_proto"],
)

py_proto_library(
    name = "asset_graph_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":asset_graph_proto"],
)

go_proto_library(
    name = "asset_graph_go_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":asset_graph_proto"],
)

proto_library(
    name = "instance_proto",
    srcs = ["instance.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/scene/proto/v1:scene_object_config_proto",
        "@com_google_protobuf//:any_proto",
    ],
)

cc_proto_library(
    name = "instance_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":instance_proto"],
)

go_proto_library(
    name = "instance_go_proto",
    go_deps = [
        "//intrinsic/scene/proto/v1:scene_object_config_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":instance_proto"],
)

py_proto_library(
    name = "instance_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":instance_proto"],
)

proto_library(
    name = "reference_proto",
    srcs = ["reference.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/proto:asset_type_proto",
        "//intrinsic/assets/proto:id_proto",
    ],
)

cc_proto_library(
    name = "reference_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":reference_proto"],
)

go_proto_library(
    name = "reference_go_proto",
    go_deps = [
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":reference_proto"],
)

py_proto_library(
    name = "reference_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":reference_proto"],
)

proto_library(
    name = "search_proto",
    srcs = ["search.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

py_proto_library(
    name = "search_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":search_proto"],
)

go_proto_library(
    name = "search_go_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":search_proto"],
)

cc_proto_library(
    name = "search_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":search_proto"],
)

# intrinsic:assets_dependencies:strip_begin

proto_library(
    name = "resolved_dependency_proto",
    srcs = ["resolved_dependency.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/proto:id_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/api:field_behavior_proto",
    ],
)

go_proto_library(
    name = "resolved_dependency_go_proto",
    go_deps = [
        "//intrinsic/assets/proto:id_go_proto",
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":resolved_dependency_proto"],
)

cc_proto_library(
    name = "resolved_dependency_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":resolved_dependency_proto"],
)

py_proto_library(
    name = "resolved_dependency_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":resolved_dependency_proto"],
)

proto_library(
    name = "dependency_proto",
    srcs = ["dependency.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

go_proto_library(
    name = "dependency_go_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":dependency_proto"],
)

cc_proto_library(
    name = "dependency_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":dependency_proto"],
)

py_proto_library(
    name = "dependency_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":dependency_proto"],
)

# intrinsic:assets_dependencies:strip_end
