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

"""Malloc test is not implemented externally, only invoke regular cc_test."""

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
