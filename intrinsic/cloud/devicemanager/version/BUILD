# Copyright 2023 Intrinsic Innovation LLC

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "version",
    srcs = ["version.go"],
    importpath = "intrinsic/cloud/devicemanager/version",
    visibility = [
        "//intrinsic/cloud/devicemanager:__subpackages__",
        "//intrinsic/frontend/cloud/devicemanager:__subpackages__",
        "//intrinsic/kubernetes/vmpool/service/internal/clients:__subpackages__",
        "//intrinsic/kubernetes/vmpool/updater/internal:__pkg__",
        "//intrinsic/tools/inctl/cmd/cluster:__subpackages__",
    ],
)

go_test(
    name = "version_test",
    size = "small",
    srcs = ["version_test.go"],
    embed = [":version"],
)
