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

# This package holds wrappers that make ICON classes available to users of the RTCL C API
# (see parent directory).

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

cc_library(
    name = "string_wrapper",
    srcs = ["string_wrapper.cc"],
    hdrs = ["string_wrapper.h"],
    deps = [
        "//intrinsic/icon/control/c_api:icon_c_api",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "streaming_io_wrapper",
    hdrs = ["streaming_io_wrapper.h"],
    deps = [
        ":string_wrapper",
        "//intrinsic/icon/control/c_api:convert_c_realtime_status",
        "//intrinsic/icon/control/c_api:icon_c_api",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_or",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)
