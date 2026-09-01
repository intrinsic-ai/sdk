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
load("@rules_cc//cc:cc_test.bzl", "cc_test")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "thread_options",
    srcs = ["thread_options.cc"],
    hdrs = ["thread_options.h"],
    deps = ["@abseil-cpp//absl/strings"],
)

cc_library(
    name = "thread_utils",
    srcs = [
        "thread_utils.cc",
    ],
    hdrs = [
        "rt_thread.h",
        "thread_utils.h",
    ],
    copts = ["-D_GNU_SOURCE"],
    deps = [
        ":stop_token",
        ":sysinfo",
        ":thread",
        ":thread_options",
        "//intrinsic/icon/testing:realtime_annotations",
        "//intrinsic/icon/utils:log",  # buildcleaner: keep
        "//intrinsic/icon/utils:realtime_guard",
        "//intrinsic/icon/utils:realtime_stack_trace",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/functional:any_invocable",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/synchronization",
    ],
)

cc_test(
    name = "thread_utils_test",
    srcs = ["thread_utils_test.cc"],
    deps = [
        ":stop_token",
        ":thread",
        ":thread_options",
        ":thread_utils",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@com_google_cel_cpp//internal:testing_no_main",
        "@com_google_googletest//:gtest_main",
    ],
)

cc_library(
    name = "thread",
    srcs = ["thread.cc"],
    hdrs = [
        "thread.h",
    ],
    deps = [
        ":stop_token",
        "//intrinsic/icon/utils:realtime_guard",
        "//intrinsic/util/status:status_macros",
    ],
)

cc_library(
    name = "stop_token",
    hdrs = ["stop_token.h"],
    deps = ["@abseil-cpp//absl/base:core_headers"],
)

cc_test(
    name = "stop_token_test",
    srcs = ["stop_token_test.cc"],
    deps = [
        ":stop_token",
        "@com_google_googletest//:gtest_main",
    ],
)

cc_library(
    name = "sysinfo",
    srcs = ["sysinfo.cc"],
    hdrs = ["sysinfo.h"],
    deps = [
        "@abseil-cpp//absl/time",
    ],
)

cc_library(
    name = "util",
    srcs = ["util.cc"],
    hdrs = ["util.h"],
    deps = [
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/container:flat_hash_set",
        "@abseil-cpp//absl/functional:any_invocable",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@com_googlesource_code_re2//:re2",
    ],
)

cc_test(
    name = "util_test",
    srcs = ["util_test.cc"],
    deps = [
        ":thread",
        ":util",
        "@abseil-cpp//absl/container:flat_hash_set",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:status_matchers",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@com_google_cel_cpp//internal:testing_no_main",
        "@com_google_googletest//:gtest_main",
        "@or_tools//ortools/base:file",
        "@or_tools//ortools/base:path",
    ],
)

cc_library(
    name = "lockstep",
    srcs = ["lockstep.cc"],
    hdrs = ["lockstep.h"],
    deps = [
        "//intrinsic/icon/interprocess:binary_futex",
        "//intrinsic/icon/utils:log",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_macro",
        "@abseil-cpp//absl/time",
    ],
)

cc_test(
    name = "lockstep_test",
    srcs = ["lockstep_test.cc"],
    deps = [
        ":lockstep",
        ":thread",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_matchers",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:status_matchers",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@com_google_cel_cpp//internal:testing_no_main",
        "@com_google_googletest//:gtest_main",
    ],
)

cc_library(
    name = "concurrent_queue",
    hdrs = ["concurrent_queue.h"],
    deps = [
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
    ],
)

cc_test(
    name = "concurrent_queue_test",
    srcs = ["concurrent_queue_test.cc"],
    deps = [
        ":concurrent_queue",
        ":thread",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:status_matchers",
        "@abseil-cpp//absl/time",
        "@com_google_cel_cpp//internal:testing_no_main",
        "@com_google_googletest//:gtest_main",
    ],
)

cc_library(
    name = "thread_pool",
    srcs = ["thread_pool.cc"],
    hdrs = ["thread_pool.h"],
    deps = [
        ":concurrent_queue",
        ":stop_token",
        ":thread",
        "@abseil-cpp//absl/functional:any_invocable",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/time",
    ],
)

cc_test(
    name = "thread_pool_test",
    srcs = ["thread_pool_test.cc"],
    deps = [
        ":stop_token",
        ":thread_pool",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:status_matchers",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@com_google_cel_cpp//internal:testing_no_main",
        "@com_google_googletest//:gtest_main",
    ],
)
