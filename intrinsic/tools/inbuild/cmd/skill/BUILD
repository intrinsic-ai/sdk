# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//intrinsic/tools/inbuild:__subpackages__"])

go_library(
    name = "bundle",
    srcs = ["bundle.go"],
    deps = [
        "//intrinsic/assets:bundleio",
        "//intrinsic/tools/inbuild/util:skillmanifest",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_library(
    name = "skill",
    srcs = ["skill.go"],
    deps = [
        ":bundle",
        "//intrinsic/tools/inbuild/cmd/skill/generate",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_test(
    name = "bundle_test",
    srcs = ["bundle_test.go"],
    data = [
        "//intrinsic/tools/inbuild/cmd/skill/test_data:example_skill_files",
        "//intrinsic/tools/inbuild/cmd/skill/test_data:example_skill_proto",
        "//intrinsic/tools/inbuild/cmd/skill/test_data:example_skill_py_image.tar",
    ],
    library = ":bundle",
    deps = [
        "//intrinsic/assets:bundleio",
        "//intrinsic/tools/inbuild/cmd/skill/test_data:example_skill_go_proto",
        "//intrinsic/util/testing:testio",
    ],
)
