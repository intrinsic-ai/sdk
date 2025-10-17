# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "solution",
    srcs = [
        "solution.go",
    ],
    visibility = [
        "//intrinsic/skills/tools:__subpackages__",
        "//intrinsic/tools/inctl:__subpackages__",
    ],
    deps = [
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/cmd/solution/get",
        "//intrinsic/tools/inctl/cmd/solution/list",
        "//intrinsic/tools/inctl/cmd/solution/start",
        "//intrinsic/tools/inctl/cmd/solution/stop",
        "//intrinsic/tools/inctl/util:orgutil",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
    ],
)
