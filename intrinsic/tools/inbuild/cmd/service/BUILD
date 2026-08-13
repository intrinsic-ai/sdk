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
    importpath = "intrinsic/tools/inbuild/cmd/service/bundle",
    deps = [
        "//intrinsic/assets/services/build_defs:servicegen",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/util/proto:protoio",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_library(
    name = "service",
    srcs = ["service.go"],
    importpath = "intrinsic/tools/inbuild/cmd/service/service",
    deps = [
        ":bundle",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_test(
    name = "bundle_test",
    srcs = ["bundle_test.go"],
    data = [
        "//intrinsic/tools/inbuild/cmd/service/test_data:example_service_files",
        "//intrinsic/tools/inbuild/cmd/service/test_data:example_service_proto",
        "//intrinsic/tools/inbuild/cmd/service/test_data:example_service_py_image.tar",
    ],
    embed = [":bundle"],
    importpath = "intrinsic/tools/inbuild/cmd/service/bundle_test",
    deps = [
        "//intrinsic/assets/services:servicebundle",
        "//intrinsic/tools/inbuild/cmd/service/test_data:example_service_go_proto",
        "//intrinsic/util/path_resolver:pathresolver",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@org_golang_google_protobuf//encoding/prototext:go_default_library",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//testing/protocmp:go_default_library",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)
