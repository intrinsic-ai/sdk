# Copyright 2026 Intrinsic Innovation LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""
Module extension for non-module dependencies
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Module extension function signatures include a ctx variable, which should not be removed.
def _non_module_deps_impl(
        ctx):  # @unused
    http_archive(
        name = "intrinsic_llvm_sysroot",
        build_file_content = """
filegroup(
    name = "all_files",
    srcs = glob(["**"]),
    visibility = ["//visibility:public"]
)""",
        sha256 = "24d7e61ceb0a26a2002bd0a3e87dfbc8e12ec95456bd1cced7fc5ccd79c47ed8",
        urls = ["https://storage.googleapis.com/intrinsic-mirror/bazel/sysroot-2025-07-22-845e86b8.tar.zst"],
    )

    # Used for compiling Wasm modules for C++.
    http_archive(
        name = "wasi_sdk",
        build_file_content = """
exports_files(["bin/clang"])

filegroup(
    name = "compiler_deps",
    srcs = glob([
        "share/wasi-sysroot/include/wasm32-wasip1/**",
        "share/wasi-sysroot/lib/wasm32-wasip1/**",
        "share/wasi-sysroot/share/wasm32-wasip1/**",
        "bin/wasm-ld",
        "lib/clang/**",
        "lib/*.so*",
    ]),
    visibility = ["//visibility:public"]
)
""",
        sha256 = "c6c38aab56e5de88adf6c1ebc9c3ae8da72f88ec2b656fb024eda8d4167a0bc5",
        strip_prefix = "wasi-sdk-24.0-x86_64-linux",
        urls = ["https://github.com/WebAssembly/wasi-sdk/releases/download/wasi-sdk-24/wasi-sdk-24.0-x86_64-linux.tar.gz"],
    )

non_module_deps_ext = module_extension(implementation = _non_module_deps_impl)
