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
    name = "create",
    srcs = ["create.go"],
    embedsrcs = [
        "templates/BUILD_cc.template",  #keep
        "templates/BUILD_py.template",  #keep
        "templates/service_cc.template",  #keep
        "templates/service_manifest.template",  #keep
        "templates/config_proto.template",  #keep
        "templates/service_py.template",  #keep
        "templates/defaults.textproto.template",  #keep
    ],
    importpath = "intrinsic/assets/services/inctl/create/create",
    visibility = ["//intrinsic/assets/services/inctl:__pkg__"],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/tools/inctl/cmd:root",
        "//intrinsic/tools/inctl/util:buildozer",
        "//intrinsic/tools/inctl/util:printer",
        "//intrinsic/tools/inctl/util:templateutil",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_test(
    name = "create_test",
    srcs = ["create_test.go"],
    embed = [":create"],
    importpath = "intrinsic/assets/services/inctl/create/create_test",
    deps = [
        "@com_github_google_go_cmp//cmp:go_default_library",
    ],
)
