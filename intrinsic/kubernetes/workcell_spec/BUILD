# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

exports_files(
    srcs = [
        "templates/workcell-cluster-service.yaml",
    ],
    visibility = ["//intrinsic/kubernetes/intrinsic_base:__pkg__"],
)

go_library(
    name = "imagetags",
    srcs = [
        "imagetags.go",
    ],
    deps = [
        "@com_github_pborman_uuid//:go_default_library",
        "@com_github_pkg_errors//:go_default_library",
    ],
)
