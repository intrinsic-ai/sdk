# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

# Application, process, product, and resource set describe all user-accessible configuration of a
# solution.

package(default_visibility = ["//visibility:public"])

go_library(
    name = "environments",
    srcs = ["environments.go"],
    visibility = [
        "//visibility:public",
    ],
)

go_test(
    name = "environments_test",
    srcs = ["environments_test.go"],
    library = ":environments",
)
