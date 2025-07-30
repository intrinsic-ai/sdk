# Copyright 2023 Intrinsic Innovation LLC

"""
Module extension for non-module dependencies
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file", "http_jar")

def _non_module_deps_impl(ctx):  # @unused
    # To update this, see https://github.com/intrinsic-ai/insrc/blob/main/bazel/sysroot/README.md.
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
    http_archive(
        name = "qt5",
        build_file = Label("//intrinsic/production/external:qt5/BUILD.bazel"),
        sha256 = "5df5be9357b425cdd70d92d4697d07e7d55d7a923f037c22dc80a78e85842d2c",
        urls = ["https://storage.googleapis.com/chrome-linux-sysroot/toolchain/4f611ec025be98214164d4bf9fbe8843f58533f7/debian_bullseye_amd64_sysroot.tar.xz"],
        patches = [Label("//intrinsic/production/external/patches:qt5.patch")],
    )

    http_archive(
        name = "xxd",
        build_file = Label("//intrinsic/production/external:xxd.BUILD.bazel"),
        sha256 = "a5cdcfcfeb13dc4deddcba461d40234dbf47e61941cb7170c9ebe147357bb62d",
        strip_prefix = "vim-9.1.0917/src/xxd",
        urls = ["https://github.com/vim/vim/archive/refs/tags/v9.1.0917.tar.gz"],
    )

    http_archive(
        name = "dumb-init",
        build_file = Label("//intrinsic/production/external:dumb_init.BUILD.bazel"),
        sha256 = "3eda470d8a4a89123f4516d26877a727c0945006c8830b7e3bad717a5f6efc4e",
        strip_prefix = "dumb-init-1.2.5",
        urls = ["https://github.com/Yelp/dumb-init/archive/refs/tags/v1.2.5.tar.gz"],
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
