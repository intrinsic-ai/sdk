# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_binary")

package(default_visibility = ["//visibility:public"])

go_binary(
    name = "inctl_external",
    srcs = ["inctl_external.go"],
    deps = [
        "//intrinsic/assets/inctl:assetcmd",
        "//intrinsic/assets/services/inctl:service",
        "//intrinsic/tools/inctl/cmd:markdown",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/cmd:skill",
        "//intrinsic/tools/inctl/cmd/auth",
        "//intrinsic/tools/inctl/cmd/bazel",
        "//intrinsic/tools/inctl/cmd/cluster",
        "//intrinsic/tools/inctl/cmd/customer",
        "//intrinsic/tools/inctl/cmd/device",
        "//intrinsic/tools/inctl/cmd/doctor",
        "//intrinsic/tools/inctl/cmd/logs",
        "//intrinsic/tools/inctl/cmd/notebook",
        "//intrinsic/tools/inctl/cmd/process",
        "//intrinsic/tools/inctl/cmd/recordings",
        "//intrinsic/tools/inctl/cmd/solution",
        "//intrinsic/tools/inctl/cmd/version",
    ],
)
