# Copyright 2023 Intrinsic Innovation LLC

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "promptutil",
    srcs = ["prompt.go"],
    importpath = "intrinsic/tools/inctl/util/promptutil",
    visibility = ["//visibility:public"],
    deps = [
        "//intrinsic/tools/inctl/util:color",
        "@com_github_spf13_cobra//:cobra",
    ],
)

go_test(
    name = "promptutil_test",
    srcs = ["prompt_test.go"],
    embed = [":promptutil"],
    deps = [
        "@com_github_spf13_cobra//:cobra",
        "@com_github_stretchr_testify//assert",
    ],
)
