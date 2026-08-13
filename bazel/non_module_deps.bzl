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

    ################################
    # Google OSS replacement files #
    #      go/insrc-g3-to-oss      #
    ################################

non_module_deps_ext = module_extension(implementation = _non_module_deps_impl)
