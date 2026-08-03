# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

go_library(
    name = "origin",
    srcs = ["origin.go"],
    importpath = "intrinsic/frontend/go/origin",
    visibility = [
        "//intrinsic/cloud/devicemanager/device:__pkg__",
        "//intrinsic/frontend:__subpackages__",
        "//intrinsic/kubernetes/acl:__pkg__",
    ],
    deps = [
        "@com_github_golang_glog//:go_default_library",
    ],
)

go_test(
    name = "origin_test",
    srcs = ["origin_test.go"],
    embed = [":origin"],
    importpath = "intrinsic/frontend/go/origin_test",
)
