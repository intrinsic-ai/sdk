# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "checks",
    srcs = ["checks.go"],
    importpath = "intrinsic/tools/inctl/cmd/doctor/checks/checks",
    deps = [
        "//intrinsic/tools/inctl/cmd/doctor/api",
        "//intrinsic/tools/inctl/cmd/doctor/checks/cluster_check:clustercheck",
        "//intrinsic/tools/inctl/cmd/doctor/checks/inctl_auth_check:inctlauthcheck",
        "//intrinsic/tools/inctl/cmd/doctor/checks/solution_check:solutioncheck",
    ],
)
