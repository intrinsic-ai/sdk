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
load(
    "//intrinsic/util/proto/build_defs:descriptor_set.bzl",
    "proto_source_code_info_transitive_descriptor_set",
)

package(default_visibility = ["//intrinsic/tools/inbuild:__subpackages__"])

proto_source_code_info_transitive_descriptor_set(
    name = "skill_services_provided_to_platform",
    visibility = ["//visibility:private"],
    deps = ["//intrinsic/skills/proto:skill_service_proto"],
)

go_library(
    name = "generate",
    srcs = ["generate.go"],
    importpath = "intrinsic/tools/inbuild/cmd/skill/generate/generate",
    deps = [
        ":config",
        ":entrypoint",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_library(
    name = "entrypoint",
    srcs = ["entrypoint.go"],
    embedsrcs = [":skill_services_provided_to_platform"],
    importpath = "intrinsic/tools/inbuild/cmd/skill/generate/entrypoint",
    deps = [
        "//intrinsic/skills:skillmanifest",
        "//intrinsic/skills/generator:gen",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "//intrinsic/util/proto:descriptor",
        "//intrinsic/util/proto:protoio",
        "@com_github_spf13_cobra//:go_default_library",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//types/descriptorpb:go_default_library",
    ],
)

go_test(
    name = "entrypoint_test",
    srcs = ["entrypoint_test.go"],
    data = [
        "//intrinsic/tools/inbuild/cmd/skill/test_data:example_skill_files",
        "//intrinsic/tools/inbuild/cmd/skill/test_data:example_skill_proto",
    ],
    embed = [":entrypoint"],
    importpath = "intrinsic/tools/inbuild/cmd/skill/generate/entrypoint_test",
    deps = [
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "//intrinsic/util/proto:protoio",
        "//intrinsic/util/testing:testio",
    ],
)

go_library(
    name = "config",
    srcs = ["config.go"],
    importpath = "intrinsic/tools/inbuild/cmd/skill/generate/config",
    deps = [
        "//intrinsic/skills/build_defs:skillserviceconfiggen",
        "//intrinsic/tools/inbuild/util:skillmanifest",
        "//intrinsic/util/proto:protoio",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_test(
    name = "config_test",
    srcs = ["config_test.go"],
    data = [
        "//intrinsic/tools/inbuild/cmd/skill/test_data:example_skill_files",
        "//intrinsic/tools/inbuild/cmd/skill/test_data:example_skill_proto",
    ],
    embed = [":config"],
    importpath = "intrinsic/tools/inbuild/cmd/skill/generate/config_test",
    deps = ["//intrinsic/util/testing:testio"],
)
