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
load("@rules_license//rules:license.bzl", "license")

package(
    default_applicable_licenses = [":license"],
    default_visibility = ["//visibility:public"],
)

exports_files(["LICENSE"])

license(
    name = "license",
    package_name = "jwt_verify_lib",
    license_kinds = [
        "@rules_license//licenses/spdx:Apache-2.0",
    ],
    license_text = "LICENSE",
    package_url = "https://github.com/google/jwt_verify_lib/tree/master/simple_lru_cache",
)

cc_library(
    name = "simple_lru_cache",
    hdrs = [
        "simple_lru_cache.h",
        "simple_lru_cache_inl.h",
    ],
    deps = [
        "@abseil-cpp//absl/container:flat_hash_map",
    ],
)
