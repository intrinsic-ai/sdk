# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "start",
    srcs = ["start.go"],
    visibility = ["//intrinsic/tools/inctl/cmd/solution:__pkg__"],
    deps = [
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/config:operationmode",
        "//intrinsic/config/proto:operation_mode_go_proto",
        "//intrinsic/kubernetes/workcell_spec/proto:deploy_go_grpc_proto",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
