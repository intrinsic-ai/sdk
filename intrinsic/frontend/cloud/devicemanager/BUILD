# Copyright 2023 Intrinsic Innovation LLC

# The device manager allows to configure specific devices and facilitates the onboarding flow for new installations.

load("//bazel:go_macros.bzl", "go_library", "go_test")

go_library(
    name = "info",
    srcs = ["info.go"],
    visibility = [
        "//intrinsic/kubernetes/inversion:__pkg__",
        "//intrinsic/tools/inctl/cmd/cluster:__pkg__",
    ],
)

go_library(
    name = "version",
    srcs = ["version.go"],
    visibility = ["//intrinsic/tools/inctl/cmd/cluster:__subpackages__"],
)

go_test(
    name = "version_test",
    size = "small",
    srcs = ["version_test.go"],
    library = ":version",
)
