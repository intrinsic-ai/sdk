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

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")
load("@rules_license//rules:license.bzl", "license")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

# Copied from https://github.com/google/sling/blob/master/LICENSE
exports_files(["LICENSE"])

license(
    name = "license",
    package_name = "sling",
    license_kinds = [
        "@rules_license//licenses/spdx:Apache-2.0",
    ],
    license_text = "LICENSE",
    package_url = "https://github.com/google/sling",
)

# Forked from https://github.com/google/sling/blob/master/sling/string/numbers.h
cc_library(
    name = "sort",
    srcs = ["sort.cc"],
    hdrs = ["sort.h"],
    applicable_licenses = [":license"],
    deps = ["@abseil-cpp//absl/strings"],
)

cc_test(
    name = "sort_test",
    srcs = ["sort_test.cc"],
    deps = [
        ":sort",
        "@abseil-cpp//absl/container:btree",
        "@com_google_googletest//:gtest_main",
    ],
)
