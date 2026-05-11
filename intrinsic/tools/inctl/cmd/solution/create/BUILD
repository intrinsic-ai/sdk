# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "create",
    srcs = ["create.go"],
    importpath = "intrinsic/tools/inctl/cmd/solution/create/create",
    visibility = ["//intrinsic/tools/inctl/cmd/solution:__pkg__"],
    deps = [
        "//intrinsic/assets:cmdutils",
        "//intrinsic/solution_versions/proto/v1:branch_go_proto",
        "//intrinsic/solution_versions/proto/v1:snapshot_go_proto",
        "//intrinsic/solution_versions/proto/v1:solution_version_service_go_proto",
        "//intrinsic/tools/inctl/auth",
        "//intrinsic/tools/inctl/util:color",
        "//intrinsic/tools/inctl/util:orgutil",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
    ],
)
