# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "doctor",
    srcs = ["doctor.go"],
    importpath = "intrinsic/tools/inctl/cmd/doctor/doctor",
    deps = [
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/cmd/doctor/api",
        "//intrinsic/tools/inctl/cmd/doctor/cmds/check",
        "//intrinsic/tools/inctl/cmd/doctor/cmds/list_checks:listchecks",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)
