# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "listchecks",
    srcs = ["listchecks.go"],
    importpath = "intrinsic/tools/inctl/cmd/doctor/cmds/list_checks/listchecks",
    deps = [
        "//intrinsic/tools/inctl/cmd/doctor/checks",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)
