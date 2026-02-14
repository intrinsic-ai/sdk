# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "inctldetails",
    srcs = ["inctldetails.go"],
    importpath = "intrinsic/tools/inctl/cmd/doctor/reporters/inctl_details/inctldetails",
    deps = [
        "//intrinsic/tools/inctl/cmd/doctor/api",
        "//intrinsic/tools/inctl/cmd/doctor/proto/v1:report_go_proto",
        "//intrinsic/tools/inctl/cmd/doctor/reporters/env_vars:envvars",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)
