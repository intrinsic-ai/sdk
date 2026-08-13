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

load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//intrinsic/tools/inbuild:__subpackages__"])

go_library(
    name = "root",
    srcs = ["root.go"],
    importpath = "intrinsic/tools/inbuild/cmd/root",
    deps = [
        "//intrinsic/production:intrinsic",
        "//intrinsic/tools/inbuild/cmd/data",
        "//intrinsic/tools/inbuild/cmd/httpjson",
        "//intrinsic/tools/inbuild/cmd/service",
        "//intrinsic/tools/inbuild/cmd/skill",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_test(
    name = "root_test",
    srcs = ["root_test.go"],
    embed = [":root"],
    importpath = "intrinsic/tools/inbuild/cmd/root_test",
)
