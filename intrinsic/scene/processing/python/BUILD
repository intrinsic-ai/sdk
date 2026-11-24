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

# Copyright 2026 Intrinsic Innovation LLC. All rights reserved.

load("@pybind11_bazel//:build_defs.bzl", "pybind_extension")

pybind_extension(
    name = "scene_object_editing_py",
    srcs = ["scene_object_editing_py.cc"],
    visibility = ["//visibility:private"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/eigenmath/python:pose3_caster",
        "//intrinsic/scene/processing:joint_editing",
        "//intrinsic/scene/processing:orient_scene_object",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/world/proto:kinematics_component_cc_proto",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:string_view",
        "@pybind11_abseil//pybind11_abseil:absl_casters",
        "@pybind11_abseil//pybind11_abseil:status_casters",
        "@pybind11_protobuf//pybind11_protobuf:native_proto_caster",
    ],
)
