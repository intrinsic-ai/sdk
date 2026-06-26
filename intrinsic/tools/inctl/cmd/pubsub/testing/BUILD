# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "testing",
    srcs = [
        "fake_asset_deployment_server.go",
        "fake_asset_instances_server.go",
        "fake_installed_assets_server.go",
        "fake_operations_server.go",
    ],
    importpath = "intrinsic/tools/inctl/cmd/pubsub/testing",
    deps = [
        "//intrinsic/assets/proto:asset_deployment_go_proto",
        "//intrinsic/assets/proto:installed_assets_go_proto",
        "//intrinsic/assets/proto/v1:asset_instances_go_proto",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
