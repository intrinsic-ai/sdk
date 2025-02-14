# Copyright 2023 Intrinsic Innovation LLC

load("@rules_jsonnet//jsonnet:jsonnet.bzl", "jsonnet_library", "jsonnet_to_json")

package(default_visibility = ["//visibility:public"])

jsonnet_library(
    name = "version_lib",
    srcs = ["versions.jsonnet"],
    visibility = ["//intrinsic:__subpackages__"],
)

jsonnet_to_json(
    name = "version_list",
    src = "versions.jsonnet",
    outs = ["versions.json"],
    visibility = ["//intrinsic:__subpackages__"],
)
