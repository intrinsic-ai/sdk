# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "data",
    srcs = ["data.go"],
    visibility = ["//intrinsic/tools/inctl:__pkg__"],
    deps = [
        "//intrinsic/assets/inctl:install",
        "//intrinsic/assets/inctl:release",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:cobrautil",
    ],
)
