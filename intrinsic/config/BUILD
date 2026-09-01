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
load("@rules_python//python:defs.bzl", "py_library", "py_test")
load("//bazel:go_macros.bzl", "go_library", "go_test")

# Application, process and resource set describe all user-accessible configuration of a
# solution.

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "environments_cc",
    srcs = ["environments.cc"],
    hdrs = ["environments.h"],
    visibility = [
        "//visibility:public",
    ],
    deps = [
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_test(
    name = "environments_test_cc",
    srcs = ["environments_test.cc"],
    deps = [
        ":environments_cc",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@com_google_googletest//:gtest_main",
    ],
)

py_library(
    name = "environments_py",
    srcs = ["environments.py"],
    visibility = [
        "//visibility:public",
    ],
)

py_test(
    name = "environments_test_py",
    srcs = ["environments_test.py"],
    main = "environments_test.py",
    deps = [
        ":environments_py",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

go_library(
    name = "environments",
    srcs = ["environments.go"],
    importpath = "intrinsic/config/environments",
    visibility = [
        "//visibility:public",
    ],
)

go_test(
    name = "environments_test",
    srcs = ["environments_test.go"],
    embed = [":environments"],
    importpath = "intrinsic/config/environments_test",
)

go_library(
    name = "operationmode",
    srcs = ["operation_mode.go"],
    importpath = "intrinsic/config/operationmode",
    deps = ["//intrinsic/config/proto:operation_mode_go_proto"],
)

go_test(
    name = "operationmode_test",
    srcs = ["operation_mode_test.go"],
    embed = [":operationmode"],
    importpath = "intrinsic/config/operationmode_test",
    deps = ["//intrinsic/config/proto:operation_mode_go_proto"],
)
