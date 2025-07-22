# Copyright 2023 Intrinsic Innovation LLC

"""Malloc test is not implemented externally, only invoke regular cc_test."""

load("@bazel_skylib//lib:new_sets.bzl", "sets")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

def cc_test_and_malloc_test(name, deps = [], local_defines = [], tags = [], **kwargs):
    cc_test(
        name = name,
        local_defines = local_defines,
        tags = tags,
        deps = deps + [
            "//intrinsic/util/testing:gtest_wrapper_main",
        ],
        **kwargs
    )
