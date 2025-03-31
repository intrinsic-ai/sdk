# Copyright 2023 Intrinsic Innovation LLC

load("@bazel_gazelle//:def.bzl", "gazelle")
load("@rules_python//python:pip.bzl", "compile_pip_requirements")

gazelle(name = "gazelle")

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
    visibility = ["//intrinsic/production/external:__pkg__"],
)

compile_pip_requirements(
    name = "requirements",
    src = "requirements.in",
    generate_hashes = True,
    requirements_txt = "requirements.txt",
    # Disable the automatically generated test target as it is slow, requires network access and
    # adds unnecessary dependencies on external services we do not control.
    tags = ["manual"],
)
