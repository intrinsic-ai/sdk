# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

go_library(
    name = "config",
    srcs = ["config.go"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/resources/proto:runtime_context_go_proto",
        "//intrinsic/util/proto:protoio",
    ],
)

go_test(
    name = "config_test",
    srcs = ["config_test.go"],
    library = ":config",
    deps = [
        "//intrinsic/resources/proto:runtime_context_go_proto",
        "//intrinsic/util/proto:protoio",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "servicemanifest",
    srcs = ["servicemanifest.go"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:metadatautils",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/assets/services/proto:service_volume_go_proto",
        "//intrinsic/util/go:validate",
        "//intrinsic/util/proto:names",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)
