# Copyright 2023 Intrinsic Innovation LLC

# Home of pubsub, a non-real-time Publish-Subscribe messaging
# for robotics systems.

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

exports_files(
    srcs = [
        "fds.yaml",
        "zenoh-router.yaml",
        "zenoh-router-with-storage.yaml",
    ],
    visibility = [
        "//intrinsic/kubernetes:__subpackages__",
        "//intrinsic/manipulation:__subpackages__",
        "//intrinsic/perception:__subpackages__",
        "//intrinsic/platform/pubsub:__subpackages__",
    ],
)

cc_library(
    name = "zenoh_pubsub_data",
    srcs = ["zenoh_pubsub_data.cc"],
    hdrs = ["zenoh_pubsub_data.h"],
    deps = [
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_config",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/base:no_destructor",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/synchronization",
    ],
)

# PubSub API library. The user-facing API providing access to
# publish-subscribe system.
cc_library(
    name = "pubsub",
    srcs = [
        "zenoh_pubsub.cc",
    ],
    hdrs = [
        "pubsub.h",
    ],
    copts = [
        "-fexceptions",
    ],
    features = ["-use_header_modules"],
    deps = [
        ":kvstore",
        ":publisher",
        ":pubsub_callbacks",
        ":queryable",
        ":subscription",
        ":topic_config",
        ":zenoh_publisher_data",
        ":zenoh_pubsub_data",
        ":zenoh_subscription_data",
        "//intrinsic/platform/pubsub/adapters:pubsub_cc_proto",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
        "//intrinsic/util/proto:type_url",
        "//intrinsic/util/status:status_conversion_rpc",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@abseil-cpp//absl/types:span",
        "@com_google_googleapis//google/rpc:status_cc_proto",
        "@com_google_protobuf//:protobuf",
    ],
)

# PubSub API library. The user-facing API providing access to
# publish-subscribe system.
cc_library(
    name = "kvstore",
    srcs = ["kvstore.cc"],
    hdrs = ["kvstore.h"],
    deps = [
        ":pubsub_callbacks",
        ":subscription",
        ":topic_config",
        ":zenoh_subscription_data",
        "//intrinsic/platform/pubsub/admin_set_grpc/v1:admin_set_cc_grpc",
        "//intrinsic/platform/pubsub/admin_set_grpc/v1:admin_set_cc_proto",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_helpers",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/strings:string_view",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_protobuf//:protobuf",
    ],
)

proto_library(
    name = "storage_location",
    srcs = ["storage_location.proto"],
)

py_proto_library(
    name = "storage_location_py_pb2",
    deps = [":storage_location"],
)

cc_proto_library(
    name = "storage_location_cc_proto",
    deps = [":storage_location"],
)

go_proto_library(
    name = "storage_location_go_proto",
    importpath = "intrinsic/platform/pubsub/storage_location_go_proto",
    protos = [":storage_location"],
)

cc_library(
    name = "zenoh_publisher_data",
    hdrs = ["zenoh_publisher_data.h"],
)

cc_library(
    name = "publisher",
    # We need a different publisher implementation also when we move to a
    # different pubsub implementation
    srcs = [
        "zenoh_publisher.cc",
    ],
    hdrs = [
        "publisher.h",
    ],
    copts = [
        "-fexceptions",
    ],
    features = ["-use_header_modules"],
    deps = [
        ":publisher_stats",
        ":zenoh_publisher_data",
        "//intrinsic/platform/pubsub/adapters:pubsub_cc_proto",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
        "//intrinsic/util:proto_time",
        "//intrinsic/util/status:ret_check",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:string_view",
        "@abseil-cpp//absl/time",
        "@com_google_protobuf//:protobuf",
        "@io_opencensus_cpp//opencensus/stats",
    ],
)

cc_library(
    name = "queryable",
    srcs =
        ["zenoh_queryable.cc"],
    hdrs = [
        "queryable.h",
    ],
    copts = [
        "-fexceptions",
    ],
    features = ["-use_header_modules"],
    deps = [
        "//intrinsic/platform/pubsub/adapters:pubsub_cc_proto",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:str_format",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "zenoh_subscription_data",
    hdrs = ["zenoh_subscription_data.h"],
    deps = [
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
    ],
)

cc_library(
    name = "subscription",
    srcs = [
        "zenoh_subscription.cc",
    ],
    hdrs = [
        "subscription.h",
    ],
    copts = [
        "-fexceptions",
    ],
    features = ["-use_header_modules"],
    deps = [
        ":zenoh_subscription_data",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "publisher_stats",
    srcs = ["publisher_stats.cc"],
    hdrs = ["publisher_stats.h"],
    deps = [
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/debugging:leak_check",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/synchronization",
    ],
)

cc_library(
    name = "zenoh_pubsub",
    srcs = ["zenoh_pubsub.cc"],
    deps = [
        ":kvstore",
        ":publisher",
        ":pubsub",
        ":queryable",
        ":subscription",
        ":zenoh_publisher_data",
        ":zenoh_pubsub_data",
        ":zenoh_subscription_data",
        "//intrinsic/platform/pubsub/adapters:pubsub_cc_proto",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_library(
    name = "zenoh_queryable",
    srcs = ["zenoh_queryable.cc"],
    deps = [
        ":queryable",
        "//intrinsic/platform/pubsub/adapters:pubsub_cc_proto",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings:str_format",
    ],
)

cc_library(
    name = "topic_config",
    srcs = ["topic_config.cc"],
    hdrs = ["topic_config.h"],
)

cc_library(
    name = "pubsub_callbacks",
    hdrs = ["pubsub_callbacks.h"],
    deps = [
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings:string_view",
    ],
)
