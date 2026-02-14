# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//intrinsic/tools/inbuild:__subpackages__"])

go_library(
    name = "root",
    srcs = ["root.go"],
    importpath = "intrinsic/tools/inbuild/cmd/root",
    deps = [
        "//intrinsic/production:intrinsic",
        "//intrinsic/tools/inbuild/cmd/service",
        "//intrinsic/tools/inbuild/cmd/skill",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_test(
    name = "root_test",
    srcs = ["root_test.go"],
    embed = [":root"],
    importpath = "intrinsic/tools/inbuild/cmd/root_test",
)
