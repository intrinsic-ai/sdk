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
    name = "start_forwarding",
    srcs = [
        "start_forwarding.go",
    ],
    importpath = "intrinsic/tools/inctl/cmd/pubsub/line_orchestration/start_forwarding",
    deps = [
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/platform/pubsub/connect/onprem/forwarding_service/proto:forwarding_service_go_proto",
        "//intrinsic/tools/inctl/cmd/pubsub:pubsub_cmd",
        "//intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_test(
    name = "start_forwarding_test",
    srcs = [
        "start_forwarding_test.go",
    ],
    embed = [":start_forwarding"],
    importpath = "intrinsic/tools/inctl/cmd/pubsub/line_orchestration/start_forwarding_test",
    deps = [
        "//intrinsic/assets/proto:asset_deployment_go_proto",
        "//intrinsic/assets/proto:installed_assets_go_proto",
        "//intrinsic/assets/proto/v1:asset_instances_go_proto",
        "//intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common",
        "//intrinsic/tools/inctl/cmd/pubsub/line_orchestration/testing",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
