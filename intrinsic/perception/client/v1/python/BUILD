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

# Python wrapper for public API.

load("@rules_python//python:defs.bzl", "py_library")

py_library(
    name = "image_utils",
    srcs = ["image_utils.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/perception/proto/v1:dimensions_py_pb2",
        "//intrinsic/perception/proto/v1:image_buffer_py_pb2",
        "@ai_intrinsic_sdks_pip_deps//numpy",
        "@ai_intrinsic_sdks_pip_deps//pillow",
    ],
)
