# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "platformlevelswitch",
    srcs = ["platformlevelswitch.go"] + select({
        "//conditions:default": ["platformlevelswitch_enterprise.go"],
    }),
    importpath = "intrinsic/assets/platformlevelswitch",
    deps = [
        "//intrinsic/kubernetes/acl:identity",
        "//intrinsic/tools/inctl/util:orgutil",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
    ],
)
