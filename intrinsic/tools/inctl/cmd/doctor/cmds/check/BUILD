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

package(default_visibility = ["//visibility:public"])

go_library(
    name = "check",
    srcs = [
        "check.go",
        "topological_sort.go",
    ],
    importpath = "intrinsic/tools/inctl/cmd/doctor/cmds/check/check",
    deps = [
        "//intrinsic/tools/inctl/cmd/doctor/api",
        "//intrinsic/tools/inctl/cmd/doctor/checks",
        "//intrinsic/tools/inctl/cmd/doctor/proto/v1:report_go_proto",
        "//intrinsic/tools/inctl/cmd/doctor/reporters",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_spf13_cobra//:go_default_library",
        "@org_golang_google_protobuf//encoding/protojson:go_default_library",
    ],
)

go_test(
    name = "check_test",
    srcs = ["topological_sort_test.go"],
    embed = [":check"],
    importpath = "intrinsic/tools/inctl/cmd/doctor/cmds/check/check_test",
    deps = ["@com_github_google_go_cmp//cmp:go_default_library"],
)
