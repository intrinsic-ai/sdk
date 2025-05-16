# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_binary")

go_binary(
    name = "inbuild",
    srcs = ["inbuild.go"],
    visibility = ["//intrinsic:__subpackages__"],
    deps = ["//intrinsic/tools/inbuild/cmd:root"],
)
