# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_binary")

go_binary(
    name = "inbuild",
    srcs = ["inbuild.go"],
    importpath = "intrinsic/tools/inbuild/inbuild",
    visibility = ["//intrinsic:intrinsic_internal"],
    deps = ["//intrinsic/tools/inbuild/cmd:root"],
)
