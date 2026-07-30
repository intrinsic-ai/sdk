# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

go_library(
    name = "serialization",
    srcs = ["serialization.go"],
    importpath = "intrinsic/httpjson/serialization",
    visibility = ["//visibility:public"],
    deps = [
        "//intrinsic/httpjson/any",
        "@com_github_grpc_ecosystem_grpc_gateway_v2//runtime",
        "@org_golang_google_protobuf//encoding/protojson:go_default_library",
    ],
)

go_test(
    name = "serialization_test",
    srcs = ["serialization_test.go"],
    embed = [":serialization"],
    deps = [
        "@com_github_grpc_ecosystem_grpc_gateway_v2//runtime",
        "@org_golang_google_protobuf//proto:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb:go_default_library",
    ],
)
