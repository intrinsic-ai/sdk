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

load("@rules_cc//cc:cc_test.bzl", "cc_test")
load("//intrinsic/assets/scene_objects/build_defs:scene_object.bzl", "intrinsic_scene_object")
load("//intrinsic/scene/build_defs:scene_object.bzl", "scene_object")
load("//intrinsic/scene/build_defs:sdf_scene_object.bzl", "sdf_scene_object")
load("//intrinsic/scene/build_defs:usd_scene_object.bzl", "usd_scene_object")

sdf_scene_object(
    name = "box",
    testonly = True,
    src = "box.sdf",
    visibility = [
        "//intrinsic/assets/scene_objects:__pkg__",
    ],
)

sdf_scene_object(
    name = "box_with_user_data",
    testonly = True,
    src = "box_with_user_data.sdf",
    deps = [
        "@com_google_protobuf//:struct_proto",
    ],
)

sdf_scene_object(
    name = "box_with_updates",
    testonly = True,
    src = "box.sdf",
    updates_pbtxts = ["box_updates.pbtxt"],
    visibility = ["//visibility:private"],
    deps = [
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:struct_proto",
        "@com_google_protobuf//:wrappers_proto",
    ],
)

cc_test(
    name = "box_with_updates_test",
    srcs = ["box_with_updates_test.cc"],
    data = [":box_with_updates"],
    deps = [
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/scene/util:scene_object_gzf",
        "//intrinsic/util/path_resolver",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/gzfile",
        "@abseil-cpp//absl/status:statusor",
        "@com_github_inazarenko_protobuf_matchers//protobuf-matchers",
        "@com_google_cel_cpp//internal:testing_no_main",
        "@com_google_googletest//:gtest_main",
        "@com_google_protobuf//:protobuf",
        "@or_tools//ortools/base:file",
    ],
)

cc_test(
    name = "box_with_user_data_test",
    srcs = ["box_with_user_data_test.cc"],
    data = [":box_with_user_data"],
    deps = [
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/scene/util:scene_object_gzf",
        "//intrinsic/util/path_resolver",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/gzfile",
        "@abseil-cpp//absl/status:statusor",
        "@com_github_inazarenko_protobuf_matchers//protobuf-matchers",
        "@com_google_cel_cpp//internal:testing_no_main",
        "@com_google_googletest//:gtest_main",
        "@com_google_protobuf//:protobuf",
        "@or_tools//ortools/base:file",
    ],
)

intrinsic_scene_object(
    name = "box_scene_object",
    testonly = True,
    manifest = "box_scene_object_manifest.textproto",
    scene_object = ":box",
    visibility = [
        "//intrinsic/assets/hardware_devices:__pkg__",
    ],
    deps = [
        "@com_google_protobuf//:struct_proto",
    ],
)

intrinsic_scene_object(
    name = "box_scene_object_type",
    testonly = True,
    manifest = "box_scene_object_resource_manifest.textproto",
    scene_object = ":box",
)

scene_object(
    name = "empty_object",
    testonly = True,
    src = "empty_object.pbtxt",
    visibility = ["//visibility:private"],
)

scene_object(
    name = "empty_object_with_updates",
    testonly = True,
    src = "empty_object.pbtxt",
    updates_pbtxts = ["box_updates.pbtxt"],
    visibility = ["//visibility:private"],
    deps = [
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:struct_proto",
        "@com_google_protobuf//:wrappers_proto",
    ],
)

intrinsic_scene_object(
    name = "empty_scene_object",
    testonly = True,
    manifest = "empty_scene_object_manifest.textproto",
    scene_object = ":empty_object",
    visibility = ["//visibility:private"],
)

intrinsic_scene_object(
    name = "empty_scene_object_type",
    testonly = True,
    manifest = "empty_scene_object_resource_manifest.textproto",
    scene_object = ":empty_object",
)

usd_scene_object(
    name = "basic_robot",
    testonly = True,
    src = "basic_robot.usda",
    visibility = ["//visibility:private"],
)

usd_scene_object(
    name = "basic_robot_with_updates",
    testonly = True,
    src = "basic_robot.usda",
    updates_pbtxts = ["basic_robot_updates.pbtxt"],
    visibility = ["//visibility:private"],
    deps = [
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:struct_proto",
        "@com_google_protobuf//:wrappers_proto",
    ],
)

cc_test(
    name = "usd_robot_with_updates_test",
    srcs = ["usd_robot_with_updates_test.cc"],
    data = [":basic_robot_with_updates"],
    deps = [
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/scene/util:scene_object_gzf",
        "//intrinsic/util/path_resolver",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/gzfile",
        "@abseil-cpp//absl/status:statusor",
        "@com_github_inazarenko_protobuf_matchers//protobuf-matchers",
        "@com_google_cel_cpp//internal:testing_no_main",
        "@com_google_googletest//:gtest_main",
        "@com_google_protobuf//:protobuf",
        "@or_tools//ortools/base:file",
    ],
)
