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

    ################################
    # Google OSS replacement files #
    #      go/insrc-g3-to-oss      #
    ################################

    XLS_COMMIT = "2e60753b05cb653cb166f4e74ebf6692c5ae393d"  # 2025-04-20
    http_file(
        name = "com_google_xls_strong_int_h",
        downloaded_file_path = "strong_int.h",
        urls = ["https://raw.githubusercontent.com/google/xls/%s/xls/common/strong_int.h" % XLS_COMMIT],
        sha256 = "8029a5dd05cb020997cfe80469abd3be0ec63044e8c1ae4da88982214186c608",
    )

non_module_deps_ext = module_extension(implementation = _non_module_deps_impl)
