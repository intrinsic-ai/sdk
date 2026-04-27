# Copyright 2023 Intrinsic Innovation LLC

load("@rules_uv//uv:pip.bzl", "pip_compile")

exports_files(
    srcs = [
        ".bazelrc",
        ".bazelversion",
    ],
    visibility = ["//intrinsic/tools/inctl/cmd/bazel/templates:__subpackages__"],
)

exports_files(
    srcs = [
        "MODULE.bazel",
        "requirements.in",
        "requirements.txt",
    ],
    visibility = [
        "//:__pkg__",
        "//intrinsic/production/external:__pkg__",
    ],
)

pip_compile(
    name = "requirements_sdk",
    requirements_in = "requirements.in",
    requirements_txt = "requirements.txt",
    # Disable the automatically generated test target as it is slow, requires
    # network access and adds unnecessary dependencies on external services we
    # do not control.
    tags = ["manual"],
)
