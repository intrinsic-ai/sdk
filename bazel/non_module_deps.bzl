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

    # Antlr as required by com_google_cel_cpp below. This is available as a module
    # (https://registry.bazel.build/modules/antlr4-cpp-runtime) but with a different repo name and
    # BUILD file, so not compatible with CEL unless we were to patch the refs to antlr in CEL.
    # Can be removed once CEL is available as a module
    # (see https://github.com/google/cel-cpp/issues/953).
    http_archive(
        name = "antlr4_runtimes",
        build_file_content = """
package(default_visibility = ["//visibility:public"])
cc_library(
    name = "cpp",
    srcs = glob(["runtime/Cpp/runtime/src/**/*.cpp"]),
    hdrs = glob(["runtime/Cpp/runtime/src/**/*.h"]),
    defines = ["ANTLR4CPP_USING_ABSEIL"],
    includes = ["runtime/Cpp/runtime/src"],
    deps = [
        "@com_google_absl//absl/base",
        "@com_google_absl//absl/base:core_headers",
        "@com_google_absl//absl/container:flat_hash_map",
        "@com_google_absl//absl/container:flat_hash_set",
        "@com_google_absl//absl/synchronization",
    ],
)
  """,
        sha256 = "365ff6aec0b1612fb964a763ca73748d80e0b3379cbdd9f82d86333eb8ae4638",
        strip_prefix = "antlr4-4.13.1",
        urls = ["https://github.com/antlr/antlr4/archive/refs/tags/4.13.1.zip"],
    )
    http_jar(
        name = "antlr4_jar",
        urls = ["https://www.antlr.org/download/antlr-4.13.1-complete.jar"],
        sha256 = "bc13a9c57a8dd7d5196888211e5ede657cb64a3ce968608697e4f668251a8487",
    )
    http_archive(
        name = "com_google_cel_cpp",
        url = "https://github.com/google/cel-cpp/archive/c4415027b89c0f4bce9db3f6d96e33bba52de87c.tar.gz",  # 2024-12-20
        strip_prefix = "cel-cpp-c4415027b89c0f4bce9db3f6d96e33bba52de87c",
        sha256 = "040ce76b4e0e7aef1c897c3dd3f85998d22407574df70697c60ed7845eea4e42",
    )

    OR_TOOLS_COMMIT = "b8e881fbde473a9e33e0dac475e498559eb0459d"  # v9.12
    http_archive(
        name = "or_tools",
        strip_prefix = "or-tools-%s" % OR_TOOLS_COMMIT,
        sha256 = "37bdd955b5224cc55cf37dea0fe4895204f21284247972f77dc47cbb504a76be",
        urls = ["https://github.com/google/or-tools/archive/%s.tar.gz" % OR_TOOLS_COMMIT],
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
