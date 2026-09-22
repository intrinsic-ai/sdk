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
    name = "line_id_utils",
    srcs = ["line_id_utils.go"],
    importpath = "intrinsic/platform/pubsub/connect/common/lineidutils",
)

go_test(
    name = "line_id_utils_test",
    srcs = ["line_id_utils_test.go"],
    embed = [":line_id_utils"],
    importpath = "intrinsic/platform/pubsub/connect/common/lineidutils_test",
)
