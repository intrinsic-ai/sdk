# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

go_library(
    name = "jwttesting",
    testonly = True,
    srcs = ["jwt_testing.go"],
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_jwt_jwt_v4//:go_default_library",
        "@com_github_pborman_uuid//:go_default_library",
    ],
)

go_test(
    name = "jwttesting_test",
    srcs = ["jwt_testing_test.go"],
    deps = [":jwttesting"],
)
