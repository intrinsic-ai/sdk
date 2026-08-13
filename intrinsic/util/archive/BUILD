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

# This package contains the tools to handle tar archives.

package(default_visibility = ["//visibility:private"])

go_library(
    name = "tartooling",
    srcs = ["tartooling.go"],
    importpath = "intrinsic/util/archive/tartooling",
    visibility = ["//intrinsic:intrinsic_internal"],
    deps = [
        "@com_github_google_safearchive//tar",
        "@com_github_pkg_errors//:go_default_library",
        "@org_golang_google_protobuf//proto",
    ],
)

go_test(
    name = "tartooling_test",
    srcs = ["tartooling_test.go"],
    embed = [":tartooling"],
    importpath = "intrinsic/util/archive/tartooling_test",
    deps = [
        "//intrinsic/util/proto/testing:diamond_a_go_proto",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_google_safearchive//tar",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
    ],
)

go_library(
    name = "untar",
    srcs = ["untar.go"],
    importpath = "intrinsic/util/archive/untar",
    visibility = ["//intrinsic:intrinsic_internal"],
    deps = ["@com_github_google_safearchive//tar"],
)
