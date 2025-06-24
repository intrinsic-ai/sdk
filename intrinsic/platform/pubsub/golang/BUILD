# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@rules_cc//cc:defs.bzl", "cc_library")
load("//bazel:go_macros.bzl", "go_library", "go_proto_library", "go_test")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "pubsub_c",
    srcs = ["pubsub_c.cc"],
    hdrs = ["pubsub_c.h"],
    visibility = ["//visibility:private"],
    deps = [
        "//intrinsic/platform/pubsub/adapters:pubsub_cc_proto",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
        "@com_google_absl//absl/base:core_headers",
        "@com_google_absl//absl/log:check",
        "@com_google_absl//absl/memory",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/status:statusor",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/synchronization",
        "@com_google_absl//absl/time",
        "@com_google_protobuf//:protobuf",
    ],
)

proto_library(
    name = "pubsub_test_proto",
    srcs = ["pubsub_test.proto"],
)

go_proto_library(
    name = "pubsub_test_go_proto",
    deps = [":pubsub_test_proto"],
)

go_test(
    name = "pubsub_test",
    size = "small",
    srcs = ["pubsub_test.go"],
    library = ":pubsub",
    # Test fails in asan configuration, this appears to be a problem in the
    # underlying C implementation.
    tags = [
        "noasan",
    ],
    deps = [
        ":pubsub_test_go_proto",
        ":pubsubinterface",
        "//intrinsic/icon/proto:part_status_go_proto",
        "//intrinsic/platform/pubsub/adapters:pubsub_go_proto",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//types/known/timestamppb",
        "@org_golang_google_protobuf//types/known/wrapperspb",
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
        "@org_golang_google_protobuf//types/known/anypb",
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

go_library(
    name = "fakekvstore",
    srcs = ["fakekvstore.go"],
    deps = [
        ":kvstore",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)

go_test(
    name = "fakekvstore_test",
    srcs = ["fakekvstore_test.go"],
    library = ":fakekvstore",
    deps = [
        ":kvstore",
        "@org_golang_google_protobuf//types/known/wrapperspb",
    ],
)
