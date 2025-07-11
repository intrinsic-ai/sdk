# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "printer",
    srcs = ["printer.go"],
    deps = ["@com_github_spf13_cobra//:go_default_library"],
)

go_library(
    name = "cobrautil",
    srcs = ["cobrautil.go"],
    deps = ["@com_github_spf13_cobra//:go_default_library"],
)

go_library(
    name = "grpc",
    srcs = ["grpc.go"],
    deps = [
        "//intrinsic/assets:baseclientutils",
        "//intrinsic/kubernetes/acl:identity",
        "//intrinsic/tools/inctl/auth",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
    ],
)

go_library(
    name = "templateutil",
    srcs = ["templateutil.go"],
)

go_library(
    name = "color",
    srcs = ["color.go"],
)

go_library(
    name = "viperutil",
    srcs = ["viperutil.go"],
    deps = [
        "@com_github_spf13_pflag//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
    ],
)

go_library(
    name = "orgutil",
    srcs = ["orgutil.go"],
    deps = [
        ":color",
        ":viperutil",
        "//intrinsic/config:environments",
        "//intrinsic/tools/inctl/auth",
        "@com_github_pkg_errors//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
    ],
)

go_test(
    name = "orgutil_test",
    srcs = ["orgutil_test.go"],
    library = ":orgutil",
    deps = [
        "//intrinsic/tools/inctl/auth",
        "//intrinsic/tools/inctl/auth:authtest",
        "@com_github_pkg_errors//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
    ],
)
