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

load("@rules_uv//uv:pip.bzl", "pip_compile")

exports_files(
    srcs = [
        ".bazelrc",
        ".bazelversion",
    ],
    visibility = ["//intrinsic/tools/inctl/cmd/bazel/templates:__subpackages__"],
)

exports_files(
    srcs = [
        "MODULE.bazel",
        "requirements.in",
        "requirements.txt",
    ],
    visibility = [
        "//:__pkg__",
        "//intrinsic/production/external:__pkg__",
    ],
)

pip_compile(
    name = "requirements_sdk",
    requirements_in = "requirements.in",
    requirements_txt = "requirements.txt",
    # Disable the automatically generated test target as it is slow, requires
    # network access and adds unnecessary dependencies on external services we
    # do not control.
    tags = ["manual"],
)
