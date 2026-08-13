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

load("@rules_python//python:defs.bzl", "py_test")

package(default_visibility = ["//visibility:public"])

exports_files(
    srcs = [
        "intrinsic.bazelrc",
        "intrinsic.bazelversion",
    ],
    visibility = [
        "//intrinsic/tools/inctl/cmd/bazel/templates:__pkg__",
    ],
)

py_test(
    name = "grpc_synchronized_py_test",
    size = "small",
    srcs = ["grpc_synchronized_test.py"],
    args = [
        "--module-bazel",
        "$(location //:MODULE.bazel)",
        "--requirements-in",
        "$(location //:requirements.in)",
        "--requirements-txt",
        "$(location //:requirements.txt)",
    ],
    data = [
        "//:MODULE.bazel",
        "//:requirements.in",
        "//:requirements.txt",
    ],
    main = "grpc_synchronized_test.py",
)
