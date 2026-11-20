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

""" Intrinsic package geometry/storage. """

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "dummy_storage",
    srcs = ["dummy_storage.cc"],
    hdrs = ["dummy_storage.h"],
    deps = [
        ":geometry_deserializer",
        ":geometry_library",
        ":geometry_serializer",
        "//intrinsic/geometry/api:exact_geometry",
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/api:geometry_fingerprint",
        "//intrinsic/geometry/api:geometry_options",
        "//intrinsic/geometry/proto:geometry_storage_refs_cc_proto",
        "//intrinsic/geometry/proto/v1:geometry_storage_refs_cc_proto",
        "//intrinsic/geometry/proto/v1:material_cc_proto",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "geometry_deserializer",
    hdrs = ["geometry_deserializer.h"],
    deps = [
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/api:geometry_options",
        "//intrinsic/geometry/proto:geometry_storage_refs_cc_proto",
        "//intrinsic/geometry/proto/v1:geometry_storage_refs_cc_proto",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "geometry_serializer",
    hdrs = ["geometry_serializer.h"],
    deps = [
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/proto/v1:geometry_storage_refs_cc_proto",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "gzf_storage",
    srcs = ["gzf_storage.cc"],
    hdrs = ["gzf_storage.h"],
    copts = [
        "-fexceptions",
        "-frounding-math",
    ],
    deps = [
        ":geometry_deserializer",
        ":geometry_library",
        ":geometry_serializer",
        "//intrinsic/geometry/api:exact_geometry",
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/api:geometry_fingerprint",
        "//intrinsic/geometry/api:geometry_options",
        "//intrinsic/geometry/api:io",
        "//intrinsic/geometry/api:renderable",
        "//intrinsic/geometry/api:renderable_generation",
        "//intrinsic/geometry/proto:geometry_storage_refs_cc_proto",
        "//intrinsic/geometry/proto:lazy_exact_geometry_cc_proto",
        "//intrinsic/geometry/proto/v1:geometry_storage_refs_cc_proto",
        "//intrinsic/geometry/proto/v1:material_cc_proto",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/gzfile",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "in_memory_storage",
    srcs = ["in_memory_storage.cc"],
    hdrs = ["in_memory_storage.h"],
    deps = [
        ":geometry_deserializer",
        ":geometry_library",
        ":geometry_serializer",
        "//intrinsic/geometry/api:exact_geometry",
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/api:geometry_fingerprint",
        "//intrinsic/geometry/api:geometry_options",
        "//intrinsic/geometry/api:renderable",
        "//intrinsic/geometry/api:renderable_generation",
        "//intrinsic/geometry/proto:geometry_storage_refs_cc_proto",
        "//intrinsic/geometry/proto/v1:geometry_storage_refs_cc_proto",
        "//intrinsic/geometry/proto/v1:material_cc_proto",
        "//intrinsic/util/status:status_builder",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/hashing",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/synchronization",
    ],
)

cc_library(
    name = "geometry_library",
    hdrs = ["geometry_library.h"],
    deps = [
        ":geometry_deserializer",
        ":geometry_serializer",
    ],
)

cc_library(
    name = "geometry_storage_uri",
    srcs = ["geometry_storage_uri.cc"],
    hdrs = ["geometry_storage_uri.h"],
    visibility = [
        "//intrinsic/geometry/storage:__pkg__",
    ],
    deps = [
        "//intrinsic/geometry/proto/v1:geometry_storage_refs_cc_proto",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
    ],
)

cc_library(
    name = "kvstore_storage",
    srcs = ["kvstore_storage.cc"],
    hdrs = ["kvstore_storage.h"],
    deps = [
        ":geometry_library",
        ":geometry_storage_uri",
        "//intrinsic/geometry/api:io",
        "//intrinsic/geometry/api:renderable_generation",
        "//intrinsic/platform/pubsub:kvstore",
        "//intrinsic/util:hash",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/base:nullability",
        "@abseil-cpp//absl/log:die_if_null",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
        "@abseil-cpp//absl/time",
        "@com_google_protobuf//:protobuf",
    ],
)
