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
    name = "pubsub_cmd",
    srcs = ["pubsub_cmd.go"],
    importpath = "intrinsic/tools/inctl/cmd/pubsub/pubsub_cmd",
    deps = [
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_library(
    name = "long_running_operation_utils",
    srcs = ["long_running_operation_utils.go"],
    importpath = "intrinsic/tools/inctl/cmd/pubsub/long_running_operation_utils",
    deps = [
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
    ],
)

go_library(
    name = "pubsub",
    srcs = [
        "pubsub.go",
    ],
    importpath = "intrinsic/tools/inctl/cmd/pubsub/pubsub",
    deps = [
        ":long_running_operation_utils",
        ":pubsub_cmd",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/cmd/pubsub/line_orchestration/hub_service_create",
        "//intrinsic/tools/inctl/cmd/pubsub/line_orchestration/hub_service_delete",
        "//intrinsic/tools/inctl/cmd/pubsub/line_orchestration/start_forwarding",
        "//intrinsic/tools/inctl/cmd/pubsub/line_orchestration/stop_forwarding",
    ],
)

go_test(
    name = "pubsub_test",
    srcs = [
        "pubsub_test.go",
    ],
    embed = [":pubsub"],
    importpath = "intrinsic/tools/inctl/cmd/pubsub/pubsub_test",
    deps = [
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
