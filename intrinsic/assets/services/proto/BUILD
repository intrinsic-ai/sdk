# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "service_manifest_proto",
    srcs = ["service_manifest.proto"],
    deps = [
        ":env_var_proto",
        ":service_http_config_proto",
        ":service_permissions_proto",
        ":service_volume_proto",
        "//intrinsic/assets/proto:asset_tag_proto",
        "//intrinsic/assets/proto:documentation_proto",
        "//intrinsic/assets/proto:id_proto",
        "//intrinsic/assets/proto:vendor_proto",
        "//intrinsic/assets/services/proto/v1:service_inspection_proto",
        "//intrinsic/kubernetes/workcell_spec/proto:image_proto",
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:descriptor_proto",
    ],
)

go_proto_library(
    name = "service_manifest_go_proto",
    protos = [":service_manifest_proto"],
    deps = [
        ":env_var_go_proto",
        ":service_http_config_go_proto",
        ":service_permissions_go_proto",
        ":service_volume_go_proto",
        "//intrinsic/assets/proto:asset_tag_go_proto",
        "//intrinsic/assets/proto:documentation_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/assets/services/proto/v1:service_inspection_go_proto",
        "//intrinsic/kubernetes/workcell_spec/proto:image_go_proto",
    ],
)

py_proto_library(
    name = "service_manifest_py_pb2",
    deps = [":service_manifest_proto"],
)

proto_library(
    name = "env_var_proto",
    srcs = ["env_var.proto"],
)

go_proto_library(
    name = "env_var_go_proto",
    protos = [":env_var_proto"],
)

py_proto_library(
    name = "env_var_py_pb2",
    deps = [":env_var_proto"],
)

proto_library(
    name = "service_http_config_proto",
    srcs = ["service_http_config.proto"],
)

go_proto_library(
    name = "service_http_config_go_proto",
    protos = [":service_http_config_proto"],
)

py_proto_library(
    name = "service_http_config_py_pb2",
    deps = [":service_http_config_proto"],
)

proto_library(
    name = "service_volume_proto",
    srcs = ["service_volume.proto"],
)

go_proto_library(
    name = "service_volume_go_proto",
    protos = [":service_volume_proto"],
)

py_proto_library(
    name = "service_volume_py_pb2",
    deps = [":service_volume_proto"],
)

proto_library(
    name = "service_permissions_proto",
    srcs = ["service_permissions.proto"],
)

go_proto_library(
    name = "service_permissions_go_proto",
    protos = [":service_permissions_proto"],
)

py_proto_library(
    name = "service_permissions_py_pb2",
    deps = [":service_permissions_proto"],
)
