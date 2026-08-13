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

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "eigen_serialization",
    srcs = [
        "eigen_serialization.cc",
    ],
    hdrs = [
        "eigen_serialization.h",
    ],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/eigenmath:pose2",
        "//intrinsic/marshal:riegeli_coder",
        "//intrinsic/math:pose3",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_riegeli//riegeli/records:record_reader",
        "@com_google_riegeli//riegeli/records:record_writer",
    ],
)
