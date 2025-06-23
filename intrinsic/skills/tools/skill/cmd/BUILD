# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "cmd",
    srcs = ["root.go"],
    deps = [
        "//intrinsic/tools/inctl/cmd:root",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_library(
    name = "solutionutil",
    srcs = ["solutionutil.go"],
    deps = [
        "//intrinsic/frontend/cloud/api/v1:clusterdiscovery_api_go_grpc_proto",
        "//intrinsic/tools/inctl/cmd/solution",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

go_library(
    name = "dialerutil",
    srcs = ["dialerutil.go"],
    deps = [
        "//intrinsic/assets:baseclientutils",
        "//intrinsic/kubernetes/acl:identity",
        "//intrinsic/tools/inctl/auth",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials:go_default_library",
        "@org_golang_google_grpc//credentials/insecure:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
    ],
)

go_library(
    name = "skillio",
    srcs = ["skillio.go"],
    deps = [
        "//intrinsic/assets:bundleio",
        "//intrinsic/assets:idutils",
    ],
)

go_library(
    name = "waitforskill",
    srcs = ["waitforskill.go"],
    deps = [
        "//intrinsic/skills/proto:skill_registry_go_grpc_proto",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_library(
    name = "listutil",
    srcs = ["listutil.go"],
    deps = ["//intrinsic/skills/proto:skills_go_proto"],
)
