# Copyright 2023 Intrinsic Innovation LLC

"""This file contains starlark rules for building shell targets."""

load("@bazel_skylib//lib:paths.bzl", "paths")
load("@rules_shell//shell:sh_binary.bzl", _sh_binary = "sh_binary")
load("@rules_shell//shell:sh_library.bzl", _sh_library = "sh_library")
load("@rules_shell//shell:sh_test.bzl", _sh_test = "sh_test")

sh_binary = _sh_binary
sh_library = _sh_library
sh_test = _sh_test
