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
    importpath = "intrinsic/tools/inbuild/cmd/skill/bundle",
    deps = [
        "//intrinsic/skills:skillbundle",
        "//intrinsic/skills:skillfix",
        "//intrinsic/skills:skillvalidate",
        "//intrinsic/tools/inbuild/util:skillmanifest",
        "@com_github_spf13_cobra//:go_default_library",
        "@org_golang_google_protobuf//reflect/protodesc:go_default_library",
    ],
)

go_library(
    name = "manifest",
    srcs = ["manifest.go"],
    importpath = "intrinsic/tools/inbuild/cmd/skill/manifest",
    deps = [
        "//intrinsic/tools/inbuild/util:skillmanifest",
        "//intrinsic/util/proto:protoio",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_library(
    name = "skill",
    srcs = ["skill.go"],
    importpath = "intrinsic/tools/inbuild/cmd/skill/skill",
    deps = [
        ":bundle",
        ":manifest",
        "//intrinsic/tools/inbuild/cmd/skill/generate",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_test(
    name = "bundle_test",
    srcs = ["bundle_test.go"],
    data = [
        "//intrinsic/tools/inbuild/cmd/skill/test_data:example_skill_augmented_test_descriptors",
        "//intrinsic/tools/inbuild/cmd/skill/test_data:example_skill_files",
        "//intrinsic/tools/inbuild/cmd/skill/test_data:example_skill_py_image.tar",
    ],
    embed = [":bundle"],
    importpath = "intrinsic/tools/inbuild/cmd/skill/bundle_test",
    deps = [
        "//intrinsic/skills:skillbundle",
        "//intrinsic/tools/inbuild/cmd/skill/test_data:example_skill_go_proto",
        "//intrinsic/util/testing:testio",
    ],
)

go_test(
    name = "manifest_test",
    srcs = ["manifest_test.go"],
    data = [
        "//intrinsic/tools/inbuild/cmd/skill/test_data:example_skill_files",
        "//intrinsic/tools/inbuild/cmd/skill/test_data:example_skill_proto",
    ],
    embed = [":manifest"],
    importpath = "intrinsic/tools/inbuild/cmd/skill/manifest_test",
    deps = [
        "//intrinsic/tools/inbuild/cmd/skill/test_data:example_skill_go_proto",
        "//intrinsic/util/testing:testio",
    ],
)
