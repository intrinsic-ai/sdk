# Copyright 2023 Intrinsic Innovation LLC

# The device manager allows to configure specific devices and facilitates the onboarding flow for new installations.

load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "version",
    srcs = ["version.go"],
    visibility = [
        "//intrinsic/tools/inctl/cmd/cluster:__subpackages__",
        "//visibility:public",
    ],
)

go_test(
    name = "version_test",
    size = "small",
    srcs = ["version_test.go"],
    library = ":version",
)
