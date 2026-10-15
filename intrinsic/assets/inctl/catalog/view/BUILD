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

go_library(
    name = "view",
    srcs = ["view.go"],
    importpath = "intrinsic/assets/inctl/catalog/view",
    visibility = [
        "//intrinsic:internal_api_users",
        "//intrinsic/assets/inctl:__subpackages__",
    ],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:typeutils",
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_proto",
        "//intrinsic/assets/proto:asset_tag_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/tools/inctl/util:printer",
        "@org_golang_google_protobuf//encoding/protojson",
    ],
)

go_test(
    name = "view_test",
    srcs = ["view_test.go"],
    embed = [":view"],
    importpath = "intrinsic/assets/inctl/catalog/view_test",
    deps = [
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_proto",
        "//intrinsic/assets/catalog/proto/v1:release_metadata_go_proto",
        "//intrinsic/assets/proto:asset_tag_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:documentation_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/tools/inctl/util:printer",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@org_golang_google_protobuf//encoding/protojson",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)
