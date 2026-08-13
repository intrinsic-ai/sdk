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

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//intrinsic/tools/inctl:__subpackages__"])

go_library(
    name = "cluster",
    srcs = [
        "cluster.go",
        "delete.go",
        "list.go",
        "mode.go",
        "poweroff.go",
        "upgrade.go",
        "ups.go",
    ],
    importpath = "intrinsic/tools/inctl/cmd/cluster/cluster",
    visibility = [
        "//intrinsic/tools/inctl:__subpackages__",
    ],
    deps = [
        "//intrinsic/assets:clientutils",
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/cloud/devicemanager/version",
        "//intrinsic/frontend/cloud/api/v1:clusterdiscovery_api_go_proto",
        "//intrinsic/frontend/cloud/api/v1:clustermanager_api_go_proto",
        "//intrinsic/frontend/cloud/api/v1alpha1:clustermanager_api_go_proto",
        "//intrinsic/kubernetes/inversion/v1:inversion_go_proto",
        "//intrinsic/kubernetes/workcellmode/proto:workcellmode_service_go_proto",
        "//intrinsic/skills/tools/skill/cmd:dialerutil",
        "//intrinsic/tools/inctl/auth",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:agents",
        "//intrinsic/tools/inctl/util:cobrautil",
        "//intrinsic/tools/inctl/util:orgutil",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_pkg_errors//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//credentials/insecure:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_golang_google_protobuf//types/known/fieldmaskpb",
    ],
)
