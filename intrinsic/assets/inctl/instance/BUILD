# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "instance",
    srcs = ["instance.go"],
    importpath = "intrinsic/assets/inctl/instance",
    visibility = ["//intrinsic/assets/inctl:__pkg__"],
    deps = [
        "//intrinsic/assets/inctl/instance/get",
        "//intrinsic/assets/inctl/instance/list",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)
