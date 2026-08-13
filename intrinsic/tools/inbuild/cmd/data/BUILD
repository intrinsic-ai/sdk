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
    name = "bundle",
    srcs = ["bundle.go"],
    importpath = "intrinsic/tools/inbuild/cmd/data/bundle",
    deps = [
        "//intrinsic/assets:referenceddata",
        "//intrinsic/assets/data:utils",
        "//intrinsic/assets/data/build_defs:datagen",
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/data/proto/v1:data_manifest_go_proto",
        "//intrinsic/assets/data/proto/v1:referenced_data_struct_go_proto",
        "//intrinsic/util/proto:descriptor",
        "//intrinsic/util/proto:registryutil",
        "@com_github_spf13_cobra//:go_default_library",
        "@org_golang_google_protobuf//types/descriptorpb:go_default_library",
    ],
)

go_test(
    name = "bundle_test",
    srcs = ["bundle_test.go"],
    data = ["//intrinsic/tools/inbuild/cmd/data/test_data:custom_data_fds"],
    embed = [":bundle"],
    importpath = "intrinsic/tools/inbuild/cmd/data/bundle_test",
    deps = [
        "//intrinsic/assets:referenceddata",
        "//intrinsic/assets/data:databundle",
        "//intrinsic/assets/data/proto/v1:data_manifest_go_proto",
        "//intrinsic/assets/data/proto/v1:referenced_data_go_proto",
        "//intrinsic/assets/data/proto/v1:referenced_data_struct_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/tools/inbuild/cmd/data/test_data:custom_data_go_proto",
        "//intrinsic/util/proto:protoio",
        "//intrinsic/util/testing:testio",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)

go_library(
    name = "data",
    srcs = ["data.go"],
    importpath = "intrinsic/tools/inbuild/cmd/data/data",
    deps = [
        ":bundle",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)
