# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "reporters",
    srcs = ["reporters.go"],
    importpath = "intrinsic/tools/inctl/cmd/doctor/reporters/reporters",
    deps = [
        "//intrinsic/tools/inctl/cmd/doctor/api",
        "//intrinsic/tools/inctl/cmd/doctor/reporters/cluster",
        "//intrinsic/tools/inctl/cmd/doctor/reporters/env_vars:envvars",
        "//intrinsic/tools/inctl/cmd/doctor/reporters/inctl_details:inctldetails",
        "//intrinsic/tools/inctl/cmd/doctor/reporters/org",
        "//intrinsic/tools/inctl/cmd/doctor/reporters/solution",
    ],
)
