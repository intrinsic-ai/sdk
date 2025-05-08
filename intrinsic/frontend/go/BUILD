# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = [
    "//intrinsic/frontend:__subpackages__",
])

go_library(
    name = "origin",
    srcs = ["origin.go"],
    visibility = [
        "//intrinsic/frontend:__subpackages__",
        "//intrinsic/kubernetes/acl:__pkg__",
    ],
)
