# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//intrinsic/tools/inbuild:__subpackages__"])

go_library(
    name = "generatemain",
    srcs = ["generatemain.go"],
    embedsrcs = [
        "templates/main.go.template",
    ],
    importpath = "intrinsic/tools/inbuild/cmd/httpjson/generatemain",
    deps = [
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_test(
    name = "generatemain_test",
    srcs = ["generatemain_test.go"],
    data = [
        "//intrinsic/tools/inbuild/cmd/httpjson/test_data",
    ],
    embed = [":generatemain"],
    importpath = "intrinsic/tools/inbuild/cmd/httpjson/generatemain_test",
    deps = ["//intrinsic/util/testing:testio"],
)

go_library(
    name = "httpjson",
    srcs = ["httpjson.go"],
    importpath = "intrinsic/tools/inbuild/cmd/httpjson/httpjson",
    deps = [
        ":generatemain",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)
