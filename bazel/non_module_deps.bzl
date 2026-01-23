# Copyright 2023 Intrinsic Innovation LLC

"""
Module extension for non-module dependencies
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file", "http_jar")

def _non_module_deps_impl(ctx):  # @unused
    http_archive(
        name = "intrinsic_llvm_sysroot",
        sha256 = "24d7e61ceb0a26a2002bd0a3e87dfbc8e12ec95456bd1cced7fc5ccd79c47ed8",
        build_file_content = """
filegroup(
    name = "all_files",
    srcs = glob(["**"]),
    visibility = ["//visibility:public"]
)""",
        urls = ["https://storage.googleapis.com/intrinsic-mirror/bazel/sysroot-2025-07-22-845e86b8.tar.zst"],
    )

non_module_deps_ext = module_extension(implementation = _non_module_deps_impl)
