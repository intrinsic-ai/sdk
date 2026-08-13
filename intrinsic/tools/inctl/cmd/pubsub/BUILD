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
    name = "pubsub",
    srcs = [
        "cmd_runner_base.go",
        "hub_service_create.go",
        "hub_service_delete.go",
        "pubsub.go",
        "service_deletion_utils.go",
        "service_installing_cmd_runner.go",
        "start_forwarding.go",
        "stop_forwarding.go",
    ],
    importpath = "intrinsic/tools/inctl/cmd/pubsub/pubsub",
    deps = [
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets:idutils",
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_proto",
        "//intrinsic/assets/proto:asset_deployment_go_proto",
        "//intrinsic/assets/proto:installed_assets_go_proto",
        "//intrinsic/assets/proto:view_go_proto",
        "//intrinsic/assets/proto/v1:asset_instances_go_proto",
        "//intrinsic/platform/pubsub/connect/onprem/forwarding_service/proto:forwarding_service_go_proto",
        "//intrinsic/platform/pubsub/connect/onprem/relay_router_service/proto:endpoint_spec_go_proto",
        "//intrinsic/platform/pubsub/connect/onprem/relay_router_service/proto:relay_router_service_go_proto",
        "//intrinsic/tools/inctl/auth",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:agents",
        "//intrinsic/tools/inctl/util:orgutil",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//credentials/insecure:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)

go_test(
    name = "pubsub_test",
    srcs = [
        "hub_service_create_test.go",
        "hub_service_test_utils.go",
        "pubsub_test.go",
        "service_deletion_utils_test.go",
        "service_installing_cmd_runner_test.go",
        "start_forwarding_test.go",
    ],
    embed = [":pubsub"],
    importpath = "intrinsic/tools/inctl/cmd/pubsub/pubsub_test",
    deps = [
        "//intrinsic/assets:cmdutils",
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_proto",
        "//intrinsic/assets/proto:asset_deployment_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:installed_assets_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/assets/proto/v1:asset_instances_go_proto",
        "//intrinsic/tools/inctl/cmd/pubsub/testing",
        "@com_github_spf13_viper//:go_default_library",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_grpc//test/bufconn:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)
