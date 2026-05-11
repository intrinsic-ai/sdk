# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "pubsub",
    srcs = [
        "forward.go",
        "hub_create.go",
        "hub_delete.go",
        "pubsub.go",
    ],
    importpath = "intrinsic/tools/inctl/cmd/pubsub/pubsub",
    deps = [
        "//intrinsic/platform/pubsub/common:command_go_proto",
        "//intrinsic/platform/pubsub/connect/cloud/proto/v1alpha1:pubsub_connect_go_proto",
        "//intrinsic/tools/inctl/auth",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:orgutil",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_grpc//metadata:go_default_library",
    ],
)

go_test(
    name = "pubsub_test",
    srcs = [
        "forward_test.go",
        "hub_create_test.go",
        "hub_delete_test.go",
        "pubsub_test.go",
    ],
    embed = [":pubsub"],
    importpath = "intrinsic/tools/inctl/cmd/pubsub/pubsub_test",
    deps = [
        "//intrinsic/platform/pubsub/common:command_go_proto",
        "//intrinsic/platform/pubsub/connect/cloud/proto/v1alpha1:pubsub_connect_go_proto",
        "@com_github_spf13_viper//:go_default_library",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)
