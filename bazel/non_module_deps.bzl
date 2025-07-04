# Copyright 2023 Intrinsic Innovation LLC

"""
Module extension for non-module dependencies
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file", "http_jar")

def _non_module_deps_impl(ctx):  # @unused
    # Sysroot and libc
    # How to upgrade:
    # - Find image in https://storage.googleapis.com/chrome-linux-sysroot/ for amd64 for
    #   a stable Linux (here: Debian bullseye), of this pick a current build.
    # - Verify the image contains expected /lib/x86_64-linux-gnu/libc* and defines correct
    #   __GLIBC_MINOR__ in /usr/include/features.h
    # - If system files are not found, add them in ../sysroot.BUILD.bazel
    http_archive(
        name = "com_googleapis_storage_chrome_linux_amd64_sysroot",
        build_file = Label("//intrinsic/production/external:sysroot.BUILD.bazel"),
        sha256 = "5df5be9357b425cdd70d92d4697d07e7d55d7a923f037c22dc80a78e85842d2c",
        urls = [
            # features.h defines GLIBC 2.31.
            "https://storage.googleapis.com/chrome-linux-sysroot/toolchain/4f611ec025be98214164d4bf9fbe8843f58533f7/debian_bullseye_amd64_sysroot.tar.xz",
        ],
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
