# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "cmd",
    srcs = ["root.go"],
    deps = [
        "//intrinsic/skills/tools/skill/cmd/create",
        "//intrinsic/skills/tools/skill/cmd/install",
        "//intrinsic/skills/tools/skill/cmd/list",
        "//intrinsic/skills/tools/skill/cmd/logs",
        "//intrinsic/skills/tools/skill/cmd/release",
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

# NOTE: Do not use this library in new code. Use auth.NewCloudConnection instead.
go_library(
    name = "dialerutil",
    srcs = ["dialerutil.go"],
    visibility = [
        "//intrinsic/skills/tools/resource/cmd:__pkg__",
        "//intrinsic/skills/tools/skill/cmd:__pkg__",
        "//intrinsic/skills/tools/skill/cmd/logs:__pkg__",
        "//intrinsic/storage/artifacts/client:__pkg__",
        "//intrinsic/tools/inctl/cmd:__pkg__",
        "//intrinsic/tools/inctl/cmd/cluster:__pkg__",
        "//intrinsic/tools/inctl/cmd/device:__pkg__",
        "//intrinsic/tools/inctl/cmd/logs:__pkg__",
        "//intrinsic/tools/inctl/cmd/modified_solution:__pkg__",
        "//intrinsic/tools/inctl/cmd/process:__pkg__",
        "//intrinsic/tools/inctl/cmd/solution:__pkg__",
    ],
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
