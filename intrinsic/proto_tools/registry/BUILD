# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@rules_python//python:defs.bzl", "py_library", "py_test")
load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//intrinsic/apps/common/resources:__subpackages__"])

py_library(
    name = "proto_registry_client",
    srcs = ["proto_registry_client.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/proto_tools/proto:proto_registry_py_pb2",
        "//intrinsic/proto_tools/proto:proto_registry_py_pb2_grpc",
        requirement("grpcio"),
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_test(
    name = "proto_registry_client_test",
    srcs = ["proto_registry_client_test.py"],
    data = [
        "//intrinsic/proto_tools/registry/test_data:descriptor_test_message_proto_descriptor_set",
    ],
    deps = [
        ":proto_registry_client",
        "//intrinsic/proto_tools/proto:proto_registry_py_pb2",
        "//intrinsic/proto_tools/proto:proto_registry_py_pb2_grpc",
        "//intrinsic/util/path_resolver:path_resolver_py",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_protobuf//:protobuf_python",
    ],
)

go_library(
    name = "protoregistryclient",
    srcs = ["proto_registry_client.go"],
    importpath = "intrinsic/proto_tools/registry/protoregistryclient",
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/proto_tools/proto:proto_registry_go_proto",
        "//intrinsic/util/proto:registryutil",
        "//intrinsic/util/proto:typeurl",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
    ],
)

go_test(
    name = "protoregistryclient_test",
    srcs = ["proto_registry_client_test.go"],
    embed = [":protoregistryclient"],
    importpath = "intrinsic/proto_tools/registry/protoregistryclient_test",
    deps = [
        "//intrinsic/proto_tools/proto:proto_registry_go_proto",
        "//intrinsic/testing:grpctest",
        "//intrinsic/util/proto:registryutil",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//credentials/insecure:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
    ],
)
