# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "origin",
    srcs = ["origin.go"],
    importpath = "intrinsic/frontend/go/origin",
    visibility = [
        "//intrinsic/cloud/devicemanager/device:__pkg__",
        "//intrinsic/frontend:__subpackages__",
        "//intrinsic/kubernetes/acl:__pkg__",
    ],
)
