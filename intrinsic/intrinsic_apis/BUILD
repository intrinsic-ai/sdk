# Copyright 2023 Intrinsic Innovation LLC

"""This package is prepared to exclusively contain protobuf files migrated from the monorepo.

These APIs will eventually be moved out of this repository into their own repository.
"""

load("@aspect_rules_js//npm:defs.bzl", "npm_package")

package(default_visibility = ["//visibility:public"])

npm_package(
    name = "pkg",
    srcs = glob(["**/*"]),
    visibility = ["//visibility:public"],
)

alias(
    name = "intrinsic_apis",
    actual = ":pkg",
    visibility = ["//visibility:public"],
)
