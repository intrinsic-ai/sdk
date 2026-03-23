# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

go_library(
    name = "anyresolver",
    srcs = ["anyresolver.go"],
    importpath = "intrinsic/httpjson/any/anyresolver",
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":greedyresolver",
        ":installedassetsresolver",
        ":protoregistryresolver",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
    ],
)

go_test(
    name = "anyresolver_test",
    srcs = ["anyresolver_test.go"],
    embed = [":anyresolver"],
    importpath = "intrinsic/httpjson/any/anyresolver_test",
    deps = [
        ":fakeserver",
        "@org_golang_google_protobuf//encoding/protojson:go_default_library",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
    ],
)

go_library(
    name = "greedyresolver",
    srcs = ["greedyresolver.go"],
    importpath = "intrinsic/httpjson/any/greedyresolver",
    deps = [
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
    ],
)

go_test(
    name = "greedyresolver_test",
    srcs = ["greedyresolver_test.go"],
    embed = [":greedyresolver"],
    importpath = "intrinsic/httpjson/any/greedyresolver_test",
    deps = [
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)

go_library(
    name = "protoregistryresolver",
    srcs = ["protoregistryresolver.go"],
    importpath = "intrinsic/httpjson/any/protoregistryresolver",
    deps = [
        "//intrinsic/proto_tools/proto:proto_registry_go_proto",
        "//intrinsic/resources/proto:runtime_context_go_proto",
        "//intrinsic/util/proto:protoio",
        "//intrinsic/util/proto:registryutil",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
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
    ],
)

go_test(
    name = "protoregistryresolver_test",
    srcs = ["protoregistryresolver_test.go"],
    embed = [":protoregistryresolver"],
    importpath = "intrinsic/httpjson/any/protoregistryresolver_test",
    deps = [
        ":fakeserver",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
    ],
)

go_library(
    name = "installedassetsresolver",
    srcs = ["installedassetsresolver.go"],
    importpath = "intrinsic/httpjson/any/installedassetsresolver",
    deps = [
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:installed_assets_go_proto",
        "//intrinsic/assets/proto:view_go_proto",
        "//intrinsic/resources/proto:runtime_context_go_proto",
        "//intrinsic/util/proto:protoio",
        "//intrinsic/util/proto:registryutil",
        "@com_github_masterminds_semver//:semver",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_grpc//credentials/insecure",
        "@org_golang_google_grpc//grpclog:go_default_library",
        "@org_golang_google_grpc//metadata",
        "@org_golang_google_protobuf//encoding/protojson:go_default_library",
        "@org_golang_google_protobuf//reflect/protodesc:go_default_library",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
    ],
)

go_test(
    name = "installedassetsresolver_test",
    srcs = ["installedassetsresolver_test.go"],
    embed = [":installedassetsresolver"],
    importpath = "intrinsic/httpjson/any/installedassetsresolver_test",
    deps = [
        ":fakeserver",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
    ],
)

go_library(
    name = "fakeserver",
    testonly = 1,
    srcs = ["fakeserver.go"],
    importpath = "intrinsic/httpjson/any/fakeserver",
    deps = [
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:installed_assets_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/proto_tools/proto:proto_registry_go_proto",
        "//intrinsic/util/proto/testing:prototestutil",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)
