# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_library", "go_test")
load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "proto_source_code_info_transitive_descriptor_set")

proto_library(
    name = "test_messages_proto",
    srcs = ["test_messages.proto"],
)

proto_source_code_info_transitive_descriptor_set(
    name = "test_messages_descriptors",
    deps = [":test_messages_proto"],
)

go_library(
    name = "any",
    srcs = [
        "anyresolver.go",
        "filedescriptorsetresolver.go",
        "greedyresolver.go",
        "installedassetsresolver.go",
        "protoregistryresolver.go",
        "resolver.go",
    ],
    importpath = "intrinsic/httpjson/any",
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:installed_assets_go_proto",
        "//intrinsic/assets/proto:view_go_proto",
        "//intrinsic/proto_tools/proto:proto_registry_go_proto",
        "//intrinsic/resources/proto:runtime_context_go_proto",
        "//intrinsic/util/proto:protoio",
        "//intrinsic/util/proto:registryutil",
        "@com_github_masterminds_semver//:semver",
        "@io_bazel_rules_go//go/runfiles:go_default_library",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//credentials/insecure",
        "@org_golang_google_grpc//grpclog:go_default_library",
        "@org_golang_google_grpc//metadata",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//encoding/protojson:go_default_library",
        "@org_golang_google_protobuf//reflect/protodesc:go_default_library",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
        "@org_golang_google_protobuf//types/descriptorpb:go_default_library",
    ],
)

go_test(
    name = "any_test",
    srcs = [
        "anyresolver_test.go",
        "fakeserver.go",
        "filedescriptorsetresolver_test.go",
        "greedyresolver_test.go",
        "installedassetsresolver_test.go",
        "protoregistryresolver_test.go",
    ],
    data = [
        ":test_messages_descriptors",
    ],
    embed = [":any"],
    importpath = "intrinsic/httpjson/any_test",
    deps = [
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:installed_assets_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/proto_tools/proto:proto_registry_go_proto",
        "//intrinsic/util/proto/testing:prototestutil",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//encoding/protojson:go_default_library",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
        "@org_golang_google_protobuf//types/descriptorpb:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)
