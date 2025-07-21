# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "pubsub_c",
    srcs = ["pubsub_c.cc"],
    hdrs = ["pubsub_c.h"],
    visibility = ["//visibility:private"],
    deps = [
        "//intrinsic/platform/pubsub/adapters:pubsub_cc_proto",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@com_google_protobuf//:protobuf",
    ],
)

go_library(
    name = "pubsub",
    srcs = ["pubsub.go"],
    cdeps = [
        ":pubsub_c",
    ],
    cgo = 1,
    data = [
        "//intrinsic/platform/pubsub/zenoh_util:peer_config.json",
    ],
    deps = [
        ":kvstore",
        ":pubsubinterface",
        "//intrinsic/platform/pubsub/adapters:pubsub_go_proto",
        "//intrinsic/util/path_resolver:pathresolver",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//types/known/anypb",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

go_library(
    name = "pubsubinterface",
    srcs = ["pubsub_interface.go"],
    deps = [
        "//intrinsic/platform/pubsub/adapters:pubsub_go_proto",
        "@org_golang_google_protobuf//proto",
    ],
)

go_library(
    name = "kvstore",
    srcs = ["kvstore.go"],
    deps = [
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)
