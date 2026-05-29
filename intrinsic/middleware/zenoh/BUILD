# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

cc_library(
    name = "imw_zenoh",
    srcs = [
        "imw_zenoh.cc",
        "imw_zenoh_data_callback_context.cc",
        "imw_zenoh_globals.cc",
        "imw_zenoh_publisher.cc",
        "imw_zenoh_queryable.cc",
        "imw_zenoh_queryable_context.cc",
        "imw_zenoh_subscription.cc",
    ],
    hdrs = [
        "imw_zenoh.h",
        "imw_zenoh_data_callback_context.h",
        "imw_zenoh_publisher.h",
        "imw_zenoh_query_context.h",
        "imw_zenoh_queryable.h",
        "imw_zenoh_queryable_context.h",
        "imw_zenoh_reply_context.h",
        "imw_zenoh_subscription.h",
    ],
    linkstatic = True,
    visibility = ["//visibility:public"],
    deps = [
        "//intrinsic/middleware:imw_interface",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/synchronization",
        "@com_github_nlohmann_json//:json",
        "@zenoh-c",
    ],
)

exports_files([
    "config/zenoh_router_config.json5",
    "config/zenoh_peer_config.json5",
])
