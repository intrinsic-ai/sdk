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
load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

# intrinsic:skills_pubsub:strip_begin(b/224840414)

py_library(
    name = "skills_pub_topic",
    srcs = ["skills_pub_topic.py"],
)

cc_library(
    name = "skill_pub_topic",
    hdrs = ["skill_pub_topic.h"],
    deps = ["@com_google_protobuf//:protobuf"],
)

cc_library(
    name = "skill_pubsub",
    srcs = ["skill_pubsub.cc"],
    hdrs = ["skill_pubsub.h"],
    deps = [
        ":skill_pub_topic",
        "//intrinsic/platform/pubsub",
        "//intrinsic/platform/pubsub:publisher",
        "//intrinsic/platform/pubsub:subscription",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/container:flat_hash_set",
        "@abseil-cpp//absl/container:node_hash_map",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/strings:string_view",
        "@abseil-cpp//absl/synchronization",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "skill_pubsub_test",
    srcs = ["skill_pubsub_test.cc"],
    deps = [
        ":skill_pubsub",
        "//intrinsic/skills/testing:skill_proto_utils_test_cc_proto",
        "//intrinsic/skills/testing:test_skill_pubsub",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:status_matchers",
        "@abseil-cpp//absl/strings:string_view",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@com_github_inazarenko_protobuf_matchers//protobuf-matchers",
        "@com_google_cel_cpp//internal:testing_no_main",
        "@com_google_googletest//:gtest_main",
    ],
)

py_library(
    name = "skill_pubsub_py",
    srcs = ["skill_pubsub.py"],
    deps = [
        ":skills_pub_topic",
        "//intrinsic/platform/pubsub/python:pubsub",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_test(
    name = "skill_pubsub_py_test",
    size = "small",
    srcs = ["skill_pubsub_py_test.py"],
    deps = [
        ":skill_pubsub_py",
        ":skills_pub_topic",
        "//intrinsic/skills/testing:skill_proto_utils_test_cc_proto",
        "//intrinsic/skills/testing:skill_proto_utils_test_py_pb2",
        "//intrinsic/skills/testing:test_skill_pubsub_py",
        "//intrinsic/solutions/testing:compare",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)
# intrinsic:skills_pubsub:strip_end

go_library(
    name = "skillbundle",
    srcs = ["skillbundle.go"],
    importpath = "intrinsic/skills/skillbundle",
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        ":skillmanifest",
        ":skillvalidate",
        "//intrinsic/assets:imageutils",
        "//intrinsic/assets:ioutils",
        "//intrinsic/skills/proto:processed_skill_manifest_go_proto",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "//intrinsic/util/archive:tartooling",
        "@com_github_google_safearchive//tar",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protodesc:go_default_library",
        "@org_golang_google_protobuf//types/descriptorpb:go_default_library",
    ],
)

go_library(
    name = "skillmanifest",
    srcs = ["skillmanifest.go"],
    importpath = "intrinsic/skills/skillmanifest",
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "//intrinsic/util/proto:sourcecodeinfoview",
        "@org_golang_google_protobuf//types/descriptorpb:go_default_library",
    ],
)

go_library(
    name = "skillvalidate",
    srcs = ["skillvalidate.go"],
    importpath = "intrinsic/skills/skillvalidate",
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:interfaceutils",
        "//intrinsic/assets:metadatautils",
        "//intrinsic/assets/dependencies:platform",
        "//intrinsic/assets/errors:report",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/skills/proto:processed_skill_manifest_go_proto",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "@org_golang_google_protobuf//reflect/protodesc:go_default_library",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
    ],
)

go_library(
    name = "skillfix",
    srcs = ["skillfix.go"],
    importpath = "intrinsic/skills/skillfix",
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        "//intrinsic/skills/proto:processed_skill_manifest_go_proto",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
    ],
)
