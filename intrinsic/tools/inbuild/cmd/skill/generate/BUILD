# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//intrinsic/tools/inbuild:__subpackages__"])

go_library(
    name = "generate",
    srcs = ["generate.go"],
    deps = [
        ":config",
        ":entrypoint",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_library(
    name = "entrypoint",
    srcs = ["entrypoint.go"],
    deps = [
        "//intrinsic/skills/generator:gen",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "//intrinsic/util/proto:protoio",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_test(
    name = "entrypoint_test",
    srcs = ["entrypoint_test.go"],
    data = [
        "//intrinsic/tools/inbuild/cmd/skill/test_data:example_skill_files",
    ],
    library = ":entrypoint",
    deps = ["//intrinsic/util/testing:testio"],
)

go_library(
    name = "config",
    srcs = ["config.go"],
    deps = [
        "//intrinsic/skills/build_defs:skillserviceconfiggen",
        "//intrinsic/tools/inbuild/util:skillmanifest",
        "//intrinsic/util/proto:protoio",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_test(
    name = "config_test",
    srcs = ["config_test.go"],
    data = [
        "//intrinsic/tools/inbuild/cmd/skill/test_data:example_skill_files",
        "//intrinsic/tools/inbuild/cmd/skill/test_data:example_skill_proto",
    ],
    library = ":config",
    deps = ["//intrinsic/util/testing:testio"],
)
