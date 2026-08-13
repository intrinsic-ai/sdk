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

""" Intrinsic package geometry/legacy/renderable_info. """

load("@rules_cc//cc:cc_library.bzl", "cc_library")

cc_library(
    name = "renderable_info",
    srcs = ["renderable_info.cc"],
    hdrs = ["renderable_info.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    visibility = [
        "//intrinsic/geometry:__subpackages__",
    ],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/internal/util:eigen_serialization",
        "//intrinsic/marshal:riegeli_coder",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@com_google_riegeli//riegeli/records:record_reader",
        "@com_google_riegeli//riegeli/records:record_writer",
    ],
)
