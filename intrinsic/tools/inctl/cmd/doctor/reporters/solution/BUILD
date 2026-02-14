# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "solution",
    srcs = ["solution.go"],
    importpath = "intrinsic/tools/inctl/cmd/doctor/reporters/solution/solution",
    deps = [
        "//intrinsic/tools/inctl/cmd/doctor/api",
        "//intrinsic/tools/inctl/cmd/doctor/proto/v1:report_go_proto",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)
