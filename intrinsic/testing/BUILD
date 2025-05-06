# Copyright 2023 Intrinsic Innovation LLC

# Package containing integration tests for Intrinsic.

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "grpctest",
    srcs = ["grpctest.go"],
    visibility = [
        "//intrinsic:__subpackages__",
    ],
    deps = [
        "@com_github_bazelbuild_remote_apis_sdks//go/pkg/portpicker",
    ],
)
