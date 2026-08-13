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

load("//bazel:go_macros.bzl", "go_binary", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

go_binary(
    name = "spinner_demo",
    srcs = ["spinner_demo.go"],
    deps = [
        ":util",
    ],
)

go_library(
    name = "buildozer",
    srcs = ["buildozer.go"],
    importpath = "intrinsic/tools/inctl/util/buildozer",
    deps = ["@com_github_bazelbuild_buildtools//edit:go_default_library"],
)

go_library(
    name = "printer",
    srcs = ["printer.go"],
    importpath = "intrinsic/tools/inctl/util/printer",
    deps = ["@com_github_spf13_cobra//:go_default_library"],
)

go_library(
    name = "cobrautil",
    srcs = ["cobrautil.go"],
    importpath = "intrinsic/tools/inctl/util/cobrautil",
    deps = ["@com_github_spf13_cobra//:go_default_library"],
)

go_library(
    name = "vmalias",
    srcs = ["vmalias.go"],
    importpath = "intrinsic/tools/inctl/util/vmalias",
)

go_library(
    name = "templateutil",
    srcs = ["templateutil.go"],
    importpath = "intrinsic/tools/inctl/util/templateutil",
)

go_library(
    name = "color",
    srcs = ["color.go"],
    importpath = "intrinsic/tools/inctl/util/color",
)

go_library(
    name = "util",
    srcs = ["spinner.go"],
    importpath = "intrinsic/tools/inctl/util",
    deps = [
        ":color",
    ],
)

go_test(
    name = "util_test",
    srcs = ["spinner_test.go"],
    embed = [":util"],
    importpath = "intrinsic/tools/inctl/util_test",
    deps = [
        "@com_github_stretchr_testify//assert:go_default_library",
    ],
)

go_library(
    name = "viperutil",
    srcs = ["viperutil.go"],
    importpath = "intrinsic/tools/inctl/util/viperutil",
    deps = [
        "@com_github_spf13_pflag//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
    ],
)

go_library(
    name = "orgutil",
    srcs = ["orgutil.go"],
    importpath = "intrinsic/tools/inctl/util/orgutil",
    deps = [
        ":color",
        ":viperutil",
        "//intrinsic/config:environments",
        "//intrinsic/tools/inctl/auth",
        "@com_github_pkg_errors//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
    ],
)

go_test(
    name = "orgutil_test",
    srcs = ["orgutil_test.go"],
    embed = [":orgutil"],
    importpath = "intrinsic/tools/inctl/util/orgutil_test",
    deps = [
        "//intrinsic/tools/inctl/auth",
        "//intrinsic/tools/inctl/auth:authtest",
        "@com_github_pkg_errors//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
    ],
)

go_library(
    name = "agents",
    srcs = ["agents.go"],
    importpath = "intrinsic/tools/inctl/util/agents",
    deps = [
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_test(
    name = "agents_test",
    srcs = ["agents_test.go"],
    embed = [":agents"],
    importpath = "intrinsic/tools/inctl/util/agents_test",
)
