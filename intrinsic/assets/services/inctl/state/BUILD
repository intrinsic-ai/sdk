# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "get",
    srcs = ["get.go"],
    deps = [
        ":stateutils",
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets/services/proto/v1:system_service_state_go_grpc_proto",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_library(
    name = "list",
    srcs = ["list.go"],
    deps = [
        ":stateutils",
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets/services/proto/v1:system_service_state_go_grpc_proto",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_library(
    name = "stateutils",
    srcs = ["stateutils.go"],
    deps = [
        "//intrinsic/assets/services/proto/v1:system_service_state_go_grpc_proto",
        "//intrinsic/tools/inctl/util:printer",
        "@org_golang_google_protobuf//encoding/protojson:go_default_library",
        "@org_golang_x_text//cases:go_default_library",
        "@org_golang_x_text//language:go_default_library",
    ],
)

go_library(
    name = "state",
    srcs = ["state.go"],
    deps = [
        ":disable",
        ":enable",
        ":get",
        ":list",
        ":restart",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_library(
    name = "disable",
    srcs = ["disable.go"],
    deps = [
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets/services/proto/v1:system_service_state_go_grpc_proto",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_library(
    name = "enable",
    srcs = ["enable.go"],
    deps = [
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets/services/proto/v1:system_service_state_go_grpc_proto",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_library(
    name = "restart",
    srcs = ["restart.go"],
    deps = [
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets/services/proto/v1:system_service_state_go_grpc_proto",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)
