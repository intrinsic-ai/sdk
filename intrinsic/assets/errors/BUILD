# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//intrinsic:public_api_users"])

go_library(
    name = "report",
    srcs = ["report.go"],
    importpath = "intrinsic/assets/errors/report",
    deps = [
        "//intrinsic/util/status:extstatus",
    ],
)
