# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_binary")

package(default_visibility = ["//visibility:public"])

go_binary(
    name = "hwmodule",
    srcs = ["hwmodule.go"],
    deps = [
        "//intrinsic/icon/hal/tools/hwmodule/cmd",
        "//intrinsic/icon/hal/tools/hwmodule/cmd:start",
        "//intrinsic/icon/hal/tools/hwmodule/cmd:stop",
    ],
)
